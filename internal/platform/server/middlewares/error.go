package middlewares

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/labstack/echo/v5"

	"github.com/NSObjects/go-template/internal/platform/apperr"
	"github.com/NSObjects/go-template/internal/platform/infrastructure/logging"
	"github.com/NSObjects/go-template/internal/platform/server/httpresp"
)

// ErrorHandler 增强的错误处理器
func ErrorHandler(c *echo.Context, err error) {
	if response, unwrapErr := echo.UnwrapResponse(c.Response()); unwrapErr == nil && response.Committed {
		return
	}

	normalized := normalizeError(err)
	info := apperr.NewInfo(normalized)
	logAPIError(c, info)
	if renderErr := httpresp.APIError(c, normalized); renderErr != nil {
		logging.FromContext(c.Request().Context()).
			Error().
			Err(renderErr).
			Str("request_id", httpresp.RequestID(c)).
			Msg("API error render failed")
	}
}

func normalizeError(err error) error {
	if err == nil {
		return apperr.New(apperr.ErrInternalServer, "internal server error")
	}

	if appErr, ok := apperr.Parse(err); ok && appErr != nil {
		return err
	}

	if status := echo.StatusCode(err); status != 0 {
		return normalizeHTTPStatus(status, err)
	}

	return apperr.WrapInternal(err, "internal server error")
}

func normalizeHTTPStatus(status int, err error) error {
	switch status {
	case http.StatusBadRequest:
		return apperr.WrapBadRequest(err, "")
	case http.StatusUnauthorized:
		return apperr.WrapUnauthorized(err, "")
	case http.StatusForbidden:
		return apperr.WrapForbidden(err, "")
	case http.StatusNotFound:
		return apperr.WrapNotFound(err, "")
	case http.StatusMethodNotAllowed:
		return apperr.WrapMethodNotAllowed(err, "")
	case http.StatusConflict:
		return apperr.WrapConflict(err, "")
	default:
		if status >= http.StatusBadRequest && status < http.StatusInternalServerError {
			return apperr.WrapBadRequest(err, "")
		}
		return apperr.WrapInternal(err, "")
	}
}

func logAPIError(c *echo.Context, info apperr.Info) {
	logger := logging.FromContext(c.Request().Context())
	event := logger.
		Warn().
		Int("code", info.Code).
		Str("message", info.Message).
		Str("category", string(info.Category)).
		Str("request_id", httpresp.RequestID(c)).
		Str("method", c.Request().Method).
		Str("path", requestPath(c)).
		Str("user_agent", c.Request().UserAgent())
	if info.Detail != "" {
		event = event.Str("detail", info.Detail)
	}

	if info.IsInternal() {
		event = logger.
			Error().
			Int("code", info.Code).
			Str("message", info.Message).
			Str("category", string(info.Category)).
			Str("request_id", httpresp.RequestID(c)).
			Str("method", c.Request().Method).
			Str("path", requestPath(c)).
			Str("user_agent", c.Request().UserAgent())
		if info.Detail != "" {
			event = event.Str("detail", info.Detail)
		}
		event.Msg("API internal error")
		return
	}
	event.Msg("API business error")
}

// ErrorRecovery 错误恢复中间件
func ErrorRecovery() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					err := apperr.WrapInternal(
						fmt.Errorf("panic recovered: %v\n%s", r, debug.Stack()),
						"internal server error",
					)
					ErrorHandler(c, err)
				}
			}()

			return next(c)
		}
	}
}
