package middlewares

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/NSObjects/go-template/internal/code"
	"github.com/NSObjects/go-template/internal/log"
	"github.com/NSObjects/go-template/internal/resp"
	"github.com/labstack/echo/v4"
)

// ErrorHandler 增强的错误处理器
func ErrorHandler(err error, c echo.Context) {
	normalized := normalizeError(err)
	info := code.NewErrorInfo(normalized)
	logAPIError(c, info)
	_ = resp.APIError(c, normalized)
}

// ValidationError 验证错误
type ValidationError struct {
	Field   string      `json:"field"`
	Message string      `json:"message"`
	Value   interface{} `json:"value,omitempty"`
}

func (e *ValidationError) Error() string {
	return e.Message
}

func normalizeError(err error) error {
	if err == nil {
		return code.NewError(code.ErrInternalServer, "internal server error")
	}

	if _, ok := code.ParseRegisteredCoder(err); ok {
		return err
	}

	switch e := err.(type) {
	case *echo.HTTPError:
		return normalizeHTTPError(e)
	case *ValidationError:
		return code.NewValidationError(e.Field, e.Message)
	default:
		return code.WrapInternalServerError(err, "internal server error")
	}
}

func normalizeHTTPError(err *echo.HTTPError) error {
	message := extractErrorMessage(err.Message)
	switch err.Code {
	case http.StatusBadRequest:
		return code.WrapBadRequestError(nil, message)
	case http.StatusUnauthorized:
		return code.WrapUnauthorizedError(nil, message)
	case http.StatusForbidden:
		return code.WrapForbiddenError(nil, message)
	case http.StatusNotFound:
		return code.WrapNotFoundError(nil, message)
	default:
		return code.WrapInternalServerError(err, message)
	}
}

func logAPIError(c echo.Context, info code.ErrorInfo) {
	fields := []slog.Attr{
		slog.Int("code", info.Code),
		slog.String("message", info.Message),
		slog.String("category", string(info.Category)),
		slog.String("request_id", resp.RequestID(c)),
		slog.String("method", c.Request().Method),
		slog.String("uri", c.Request().RequestURI),
		slog.String("user_agent", c.Request().UserAgent()),
	}
	if info.Details != "" {
		fields = append(fields, slog.String("details", info.Details))
	}

	if info.IsInternal() {
		log.Error("API internal error", fields...)
		return
	}
	log.Warn("API business error", fields...)
}

// extractErrorMessage 将 Echo 错误消息转换为字符串
func extractErrorMessage(message interface{}) string {
	switch v := message.(type) {
	case string:
		return v
	case error:
		return v.Error()
	default:
		return fmt.Sprint(v)
	}
}

// ErrorRecovery 错误恢复中间件
func ErrorRecovery() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					err := code.WrapInternalServerError(fmt.Errorf("panic recovered: %v", r), "internal server error")
					ErrorHandler(err, c)
				}
			}()

			return next(c)
		}
	}
}
