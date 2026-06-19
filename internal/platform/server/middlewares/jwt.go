package middlewares

import (
	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"

	"github.com/NSObjects/go-template/internal/platform/apperr"
	"github.com/NSObjects/go-template/internal/platform/infrastructure/logging"
	"github.com/NSObjects/go-template/internal/platform/requestctx"
)

const jwtContextKey = "user"

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
			"/api/ready",
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
		ContextKey: jwtContextKey,
		Skipper: func(c *echo.Context) bool {
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
		ErrorHandler: func(_ *echo.Context, err error) error {
			return apperr.WrapUnauthorized(err, "")
		},
		SuccessHandler: storeJWTSubject,
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

func storeJWTSubject(c *echo.Context) error {
	token, ok := c.Get(jwtContextKey).(*jwt.Token)
	if !ok || token == nil || token.Claims == nil {
		return nil
	}

	subject, err := token.Claims.GetSubject()
	if err != nil || subject == "" {
		return nil
	}

	request := c.Request()
	if request == nil {
		return nil
	}
	ctx := requestctx.WithUserID(request.Context(), subject)
	logger := logging.FromContext(ctx).With().Str("user_id", subject).Logger()
	c.SetRequest(request.WithContext(logger.WithContext(ctx)))
	return nil
}
