package middlewares

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/NSObjects/echo-admin/internal/log"
	"github.com/NSObjects/echo-admin/internal/resp"
	"github.com/labstack/echo/v4"
	"github.com/marmotedu/errors"
)

// ErrorHandler 增强的错误处理器
func ErrorHandler(err error, c echo.Context) {
	// 记录错误开始时间
	start := time.Now()

	// 处理不同类型的错误
	switch e := err.(type) {
	case *echo.HTTPError:
		handleHTTPError(e, c)
	case *ValidationError:
		handleValidationError(e, c)
	default:
		handleGenericError(err, c)
	}

	// 记录处理时间
	duration := time.Since(start)
	log.Debug("Error handled",
		slog.Duration("duration", duration),
		slog.String("method", c.Request().Method),
		slog.String("uri", c.Request().RequestURI),
	)
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

// handleHTTPError 处理HTTP错误
func handleHTTPError(err *echo.HTTPError, c echo.Context) {
	// 将Echo HTTP错误转换为业务错误
	var bizErr error
	switch err.Code {
	case http.StatusBadRequest:
		bizErr = errors.WithCode(100400, "%s", err.Message.(string))
	case http.StatusUnauthorized:
		bizErr = errors.WithCode(100401, "%s", err.Message.(string))
	case http.StatusForbidden:
		bizErr = errors.WithCode(100403, "%s", err.Message.(string))
	case http.StatusNotFound:
		bizErr = errors.WithCode(100404, "%s", err.Message.(string))
	default:
		bizErr = errors.WithCode(100500, "%s", err.Message.(string))
	}

	// 返回标准错误响应
	_ = resp.APIError(bizErr, c)
}

// handleValidationError 处理验证错误
func handleValidationError(err *ValidationError, c echo.Context) {
	// 记录验证错误
	log.Warn("Validation Error",
		slog.String("field", err.Field),
		slog.String("message", err.Message),
		slog.Any("value", err.Value),
		slog.String("method", c.Request().Method),
		slog.String("uri", c.Request().RequestURI),
	)

	// 创建业务错误
	bizErr := errors.WithCode(100400, "validation failed") // 使用验证错误码
	_ = resp.APIError(bizErr, c)
}

// handleGenericError 处理通用错误
func handleGenericError(err error, c echo.Context) {
	// 记录通用错误
	log.Error("Generic Error",
		slog.String("error", err.Error()),
		slog.String("method", c.Request().Method),
		slog.String("uri", c.Request().RequestURI),
	)

	// 返回标准错误响应
	_ = resp.APIError(err, c)
}

// ErrorRecovery 错误恢复中间件
func ErrorRecovery() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					// 记录panic信息
					log.Error("Panic recovered",
						slog.Any("panic", r),
						slog.String("method", c.Request().Method),
						slog.String("uri", c.Request().RequestURI),
					)

					// 创建内部服务器错误
					err := errors.WithCode(100500, "internal server error")
					_ = resp.APIError(err, c)
				}
			}()

			return next(c)
		}
	}
}
