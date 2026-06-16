/*
 * JWT Middleware
 * JWT认证中间件
 *
 * Created by lintao on 2024/1/4
 * Copyright © 2020-2024 LINTAO. All rights reserved.
 */

package middlewares

import (
	"strings"

	"github.com/NSObjects/go-template/internal/apperr"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

// JWTConfig JWT中间件配置
type JWTConfig struct {
	// 签名密钥
	SigningKey []byte
	// 跳过路径
	SkipPaths []string
	// 是否启用
	Enabled bool
}

// DefaultJWTConfig 默认JWT配置
func DefaultJWTConfig() *JWTConfig {
	return &JWTConfig{
		SigningKey: nil,
		SkipPaths: []string{
			"/api/health",
			"/api/info",
			"/api/login",
		},
		Enabled: false,
	}
}

// JWT JWT认证中间件
func JWT(config *JWTConfig) echo.MiddlewareFunc {
	if config == nil || !config.Enabled {
		// 如果 JWT 未启用，返回空中间件
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return next
		}
	}
	if len(config.SigningKey) == 0 {
		panic("jwt signing key is required when JWT is enabled")
	}

	return echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(JwtCustomClaims)
		},
		SigningKey: config.SigningKey,
		Skipper: func(c echo.Context) bool {
			path := c.Path()

			for _, skipPath := range config.SkipPaths {
				// 支持精确匹配和前缀匹配
				if path == skipPath ||
					(len(skipPath) > 0 && skipPath[len(skipPath)-1] == '*' &&
						len(path) >= len(skipPath)-1 &&
						strings.HasPrefix(path, skipPath[:len(skipPath)-1])) {
					return true
				}
			}
			return false
		},
		ErrorHandler: func(c echo.Context, err error) error {
			return apperr.Wrap(err, apperr.ErrSignatureInvalid, "JWT signature is invalid")
		},
	})
}

// CreateJWTConfig 从应用配置创建JWT配置
func CreateJWTConfig(secret string, skipPaths []string, enabled bool) *JWTConfig {
	return &JWTConfig{
		SigningKey: []byte(secret),
		SkipPaths:  skipPaths,
		Enabled:    enabled,
	}
}
