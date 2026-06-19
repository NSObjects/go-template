// Package middlewares contains server-owned Echo middleware adapters.
package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	echootel "github.com/labstack/echo-opentelemetry"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"

	"github.com/NSObjects/go-template/internal/platform/infrastructure/logging"
	"github.com/NSObjects/go-template/internal/platform/requestctx"
)

// MiddlewareConfig controls server-owned HTTP middleware.
type MiddlewareConfig struct {
	EnableRecovery bool

	EnableRequestContext bool

	EnableLogger bool

	EnableTracing bool

	TracingServiceName string

	EnableGzip bool

	EnableCORS bool

	CORS middleware.CORSConfig

	EnableJWT bool

	JWT *JWTConfig
}

// DefaultMiddlewareConfig returns the HTTP middleware defaults.
func DefaultMiddlewareConfig() *MiddlewareConfig {
	return &MiddlewareConfig{
		EnableRecovery:       true,
		EnableRequestContext: true,
		EnableLogger:         true,
		EnableGzip:           true,
		EnableCORS:           false,
		EnableJWT:            false,
		JWT:                  DefaultJWTConfig(),
	}
}

// ApplyMiddlewares installs server-owned middleware.
func ApplyMiddlewares(e *echo.Echo, config *MiddlewareConfig) error {
	if config == nil {
		config = DefaultMiddlewareConfig()
	}

	if config.EnableRecovery {
		e.Use(ErrorRecovery())
	}

	if config.EnableRequestContext {
		e.Use(RequestContext())
	}

	if config.EnableTracing {
		e.Use(echootel.NewMiddleware(config.TracingServiceName))
	}

	if config.EnableLogger {
		e.Use(requestLogger())
	}

	if config.EnableGzip {
		e.Use(middleware.Gzip())
	}

	if config.EnableCORS {
		corsConfig, err := normalizedCORSConfig(config.CORS)
		if err != nil {
			return err
		}
		e.Use(middleware.CORSWithConfig(corsConfig))
	}

	if config.EnableJWT && config.JWT != nil {
		jwtMiddleware, err := JWT(config.JWT)
		if err != nil {
			return err
		}
		e.Use(jwtMiddleware)
	}
	return nil
}

func requestLogger() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		BeforeNextFunc: func(c *echo.Context) {
			request := c.Request()
			ctx := contextWithActiveTrace(request.Context())
			logger := requestLoggerFromContext(ctx)
			c.SetRequest(request.WithContext(logger.WithContext(ctx)))
		},
		HandleError:  true,
		LogLatency:   true,
		LogMethod:    true,
		LogRoutePath: true,
		LogStatus:    true,
		LogURIPath:   true,
		LogValuesFunc: func(c *echo.Context, values middleware.RequestLoggerValues) error {
			logger := requestLoggerFromContext(c.Request().Context())
			status := responseStatus(c, values.Error)
			event := logger.Info()
			if status >= http.StatusInternalServerError {
				event = logger.Error()
			} else if status >= http.StatusBadRequest {
				event = logger.Warn()
			}
			if values.Error != nil {
				event = event.Err(values.Error)
			}
			if info, ok := requestctx.FromContext(c.Request().Context()); ok && info.UserID != "" {
				event = event.Str("user_id", info.UserID)
			}

			event.
				Str("method", values.Method).
				Str("path", requestLogPath(c, values)).
				Int("status", status).
				Dur("latency", values.Latency).
				Msg("HTTP request")

			return nil
		},
	})
}

func requestLoggerFromContext(ctx context.Context) zerolog.Logger {
	ctx = contextWithActiveTrace(ctx)
	info, _ := requestctx.FromContext(ctx)
	builder := logging.FromContext(ctx).With()
	if info.RequestID != "" {
		builder = builder.Str("request_id", info.RequestID)
	}
	if info.TraceID != "" {
		builder = builder.Str("trace_id", info.TraceID)
	}
	if info.SpanID != "" {
		builder = builder.Str("span_id", info.SpanID)
	}
	if info.UserID != "" {
		builder = builder.Str("user_id", info.UserID)
	}
	return builder.Logger()
}

func contextWithActiveTrace(ctx context.Context) context.Context {
	spanContext := trace.SpanContextFromContext(ctx)
	if !spanContext.IsValid() {
		return ctx
	}
	return requestctx.WithTraceSpan(ctx, spanContext.TraceID().String(), spanContext.SpanID().String())
}

func responseStatus(c *echo.Context, err error) int {
	_, status := echo.ResolveResponseStatus(c.Response(), err)
	return status
}

func requestLogPath(c *echo.Context, values middleware.RequestLoggerValues) string {
	if values.RoutePath != "" {
		return values.RoutePath
	}
	if path := c.Path(); path != "" {
		return path
	}
	if values.URIPath != "" {
		return values.URIPath
	}
	if c.Request() == nil || c.Request().URL == nil {
		return ""
	}
	return c.Request().URL.Path
}

func requestPath(c *echo.Context) string {
	if path := c.Path(); path != "" {
		return path
	}
	if c.Request() == nil || c.Request().URL == nil {
		return ""
	}
	return c.Request().URL.Path
}

func normalizedCORSConfig(config middleware.CORSConfig) (middleware.CORSConfig, error) {
	if len(config.AllowOrigins) == 0 {
		return middleware.CORSConfig{}, errors.New("cors allowed origins are required when cors is enabled")
	}
	if config.AllowCredentials && corsOriginsContainWildcard(config.AllowOrigins) {
		return middleware.CORSConfig{}, errors.New("cors allowed origins must not include wildcard when credentials are enabled")
	}
	if len(config.AllowHeaders) == 0 {
		config.AllowHeaders = []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		}
	}
	if len(config.AllowMethods) == 0 {
		config.AllowMethods = []string{
			http.MethodGet,
			http.MethodHead,
			http.MethodPut,
			http.MethodPatch,
			http.MethodPost,
			http.MethodDelete,
			http.MethodOptions,
		}
	}
	return config, nil
}

func corsOriginsContainWildcard(origins []string) bool {
	for _, origin := range origins {
		if strings.TrimSpace(origin) == "*" {
			return true
		}
	}
	return false
}
