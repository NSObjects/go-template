// Package middlewares contains server-owned Echo middleware adapters.
package middlewares

import (
	"errors"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// MiddlewareConfig controls server-owned HTTP middleware.
type MiddlewareConfig struct {
	EnableRecovery bool

	EnableRequestContext bool

	EnableLogger bool

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
		LogMethod:  true,
		LogURI:     true,
		LogStatus:  true,
		LogLatency: true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			c.Logger().Printf(
				"method=%s, uri=%s, status=%d, latency=%s\n",
				values.Method,
				values.URI,
				values.Status,
				values.Latency.Round(time.Microsecond),
			)
			return nil
		},
	})
}

func normalizedCORSConfig(config middleware.CORSConfig) (middleware.CORSConfig, error) {
	if len(config.AllowOrigins) == 0 {
		return middleware.CORSConfig{}, errors.New("cors allowed origins are required when cors is enabled")
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
			echo.GET,
			echo.HEAD,
			echo.PUT,
			echo.PATCH,
			echo.POST,
			echo.DELETE,
			echo.OPTIONS,
		}
	}
	return config, nil
}
