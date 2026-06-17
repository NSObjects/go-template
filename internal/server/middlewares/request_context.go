package middlewares

import (
	"github.com/NSObjects/go-template/internal/requestctx"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const (
	headerRequestID = "X-Request-ID"
	headerTraceID   = "X-Trace-ID"
	headerSpanID    = "X-Span-ID"
)

// RequestContext stores request-scoped metadata in the standard context.
func RequestContext() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			request := c.Request()
			requestID := request.Header.Get(headerRequestID)
			if requestID == "" {
				requestID = uuid.NewString()
			}
			c.Response().Header().Set(headerRequestID, requestID)

			info := requestctx.Info{
				TraceID:   request.Header.Get(headerTraceID),
				SpanID:    request.Header.Get(headerSpanID),
				RequestID: requestID,
			}
			c.SetRequest(request.WithContext(requestctx.WithInfo(request.Context(), info)))
			return next(c)
		}
	}
}
