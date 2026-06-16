package middlewares

import (
	"errors"
	"strings"

	"github.com/NSObjects/go-template/internal/apperr"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

// JWTConfig controls JWT verification middleware.
type JWTConfig struct {
	SigningKey []byte
	SkipPaths  []string
	Enabled    bool
}

// DefaultJWTConfig returns disabled JWT verification defaults.
func DefaultJWTConfig() *JWTConfig {
	return &JWTConfig{
		SigningKey: nil,
		SkipPaths: []string{
			"/api/health",
			"/api/info",
		},
		Enabled: false,
	}
}

// JWT creates JWT verification middleware.
func JWT(config *JWTConfig) (echo.MiddlewareFunc, error) {
	if config == nil || !config.Enabled {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return next
		}, nil
	}
	if len(config.SigningKey) == 0 {
		return nil, errors.New("jwt signing key is required when jwt is enabled")
	}

	return echojwt.WithConfig(echojwt.Config{
		SigningKey: config.SigningKey,
		Skipper: func(c echo.Context) bool {
			path := c.Path()

			for _, skipPath := range config.SkipPaths {
				if path == skipPath ||
					(len(skipPath) > 0 && skipPath[len(skipPath)-1] == '*' &&
						len(path) >= len(skipPath)-1 &&
						strings.HasPrefix(path, skipPath[:len(skipPath)-1])) {
					return true
				}
			}
			return false
		},
		ErrorHandler: func(_ echo.Context, err error) error {
			return apperr.Wrap(err, apperr.ErrSignatureInvalid, "JWT signature is invalid")
		},
	}), nil
}

// CreateJWTConfig creates JWT middleware config from application config.
func CreateJWTConfig(secret string, skipPaths []string, enabled bool) *JWTConfig {
	return &JWTConfig{
		SigningKey: []byte(secret),
		SkipPaths:  skipPaths,
		Enabled:    enabled,
	}
}
