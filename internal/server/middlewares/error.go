package middlewares

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/NSObjects/go-template/internal/apperr"
	"github.com/NSObjects/go-template/internal/server/httpresp"
	"github.com/labstack/echo/v4"
)

// ErrorHandler 增强的错误处理器
func ErrorHandler(err error, c echo.Context) {
	normalized := normalizeError(err)
	info := apperr.NewInfo(normalized)
	logAPIError(c, info)
	_ = httpresp.APIError(c, normalized)
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
		return apperr.New(apperr.ErrInternalServer, "internal server error")
	}

	if _, ok := apperr.Parse(err); ok {
		return err
	}

	var httpErr *echo.HTTPError
	if errors.As(err, &httpErr) {
		return normalizeHTTPError(httpErr)
	}

	var validationErr *ValidationError
	if errors.As(err, &validationErr) {
		return apperr.NewValidation(validationErr.Field, validationErr.Message)
	}

	return apperr.WrapInternal(err, "internal server error")
}

func normalizeHTTPError(err *echo.HTTPError) error {
	message := extractErrorMessage(err.Message)
	switch err.Code {
	case http.StatusBadRequest:
		return apperr.WrapBadRequest(nil, message)
	case http.StatusUnauthorized:
		return apperr.WrapUnauthorized(nil, message)
	case http.StatusForbidden:
		return apperr.WrapForbidden(nil, message)
	case http.StatusNotFound:
		return apperr.WrapNotFound(nil, message)
	default:
		return apperr.WrapInternal(err, message)
	}
}

func logAPIError(c echo.Context, info apperr.Info) {
	fields := []slog.Attr{
		slog.Int("code", info.Code),
		slog.String("message", info.Message),
		slog.String("category", string(info.Category)),
		slog.String("request_id", httpresp.RequestID(c)),
		slog.String("method", c.Request().Method),
		slog.String("uri", c.Request().RequestURI),
		slog.String("user_agent", c.Request().UserAgent()),
	}
	if info.Detail != "" {
		fields = append(fields, slog.String("detail", info.Detail))
	}

	if info.IsInternal() {
		slog.LogAttrs(c.Request().Context(), slog.LevelError, "API internal error", fields...)
		return
	}
	slog.LogAttrs(c.Request().Context(), slog.LevelWarn, "API business error", fields...)
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
					err := apperr.WrapInternal(fmt.Errorf("panic recovered: %v", r), "internal server error")
					ErrorHandler(err, c)
				}
			}()

			return next(c)
		}
	}
}
