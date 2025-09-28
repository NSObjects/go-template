/*
 * JWT Middleware
 * JWT认证中间件
 *
 * Created by lintao on 2024/1/4
 * Copyright © 2020-2024 LINTAO. All rights reserved.
 */

package middlewares

import (
	"fmt"
	"strings"

	"github.com/NSObjects/go-template/internal/code"
	"github.com/NSObjects/go-template/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/marmotedu/errors"
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
		SigningKey: []byte("default-secret"),
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
	if !config.Enabled {
		// 如果JWT未启用，返回空中间件
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return next
		}
	}

	return echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(utils.JwtCustomClaims)
		},
		SigningKey: config.SigningKey,
		Skipper: func(c echo.Context) bool {
			path := c.Path()

			// 调试日志
			fmt.Printf("DEBUG: Request path: %s, Skip paths: %v\n", path, config.SkipPaths)

			for _, skipPath := range config.SkipPaths {
				// 支持精确匹配和前缀匹配
				if path == skipPath ||
					(len(skipPath) > 0 && skipPath[len(skipPath)-1] == '*' &&
						len(path) >= len(skipPath)-1 &&
						strings.HasPrefix(path, skipPath[:len(skipPath)-1])) {
					fmt.Printf("DEBUG: Skipping JWT for path: %s\n", path)
					return true
				}
			}
			return false
		},
		ErrorHandler: func(c echo.Context, err error) error {
			return errors.WrapC(err, code.ErrSignatureInvalid, "JWT签名无效")
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
