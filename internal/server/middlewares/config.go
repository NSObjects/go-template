/*
 * Middleware Configuration
 * 中间件配置管理
 *
 * Created by lintao on 2024/1/4
 * Copyright © 2020-2024 LINTAO. All rights reserved.
 */

package middlewares

import (
	"github.com/casbin/casbin/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// MiddlewareConfig 中间件配置
type MiddlewareConfig struct {
	// 是否启用错误恢复
	EnableRecovery bool
	// 是否启用请求日志
	EnableLogger bool
	// 是否启用压缩
	EnableGzip bool
	// 是否启用CORS
	EnableCORS bool
	// 是否启用JWT
	EnableJWT bool
	// 是否启用Casbin
	EnableCasbin bool
	// 日志格式
	LoggerFormat string
	// JWT配置
	JWT *JWTConfig
	// Casbin配置
	Casbin *CasbinConfig
}

// DefaultMiddlewareConfig 默认中间件配置
func DefaultMiddlewareConfig() *MiddlewareConfig {
	return &MiddlewareConfig{
		EnableRecovery: true,
		EnableLogger:   true,
		EnableGzip:     true,
		EnableCORS:     true,
		EnableJWT:      false,
		EnableCasbin:   false,
		LoggerFormat:   "method=${method}, uri=${uri}, status=${status}, latency=${latency_human}\n",
		JWT:            DefaultJWTConfig(),
		Casbin:         DefaultCasbinConfig(),
	}
}

// ApplyMiddlewares 应用中间件
func ApplyMiddlewares(e *echo.Echo, config *MiddlewareConfig) {
	if config == nil {
		config = DefaultMiddlewareConfig()
	}

	// 错误恢复中间件
	if config.EnableRecovery {
		e.Use(ErrorRecovery())
	}

	// 请求日志中间件
	if config.EnableLogger {
		e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: config.LoggerFormat,
		}))
	}

	// 压缩中间件
	if config.EnableGzip {
		e.Use(middleware.Gzip())
	}

	// CORS中间件
	if config.EnableCORS {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowHeaders: []string{
				echo.HeaderOrigin,
				echo.HeaderContentType,
				echo.HeaderAccept,
				echo.HeaderAuthorization,
			},
			AllowMethods: []string{
				echo.GET,
				echo.HEAD,
				echo.PUT,
				echo.PATCH,
				echo.POST,
				echo.DELETE,
				echo.OPTIONS,
			},
		}))
	}

	// JWT中间件
	if config.EnableJWT && config.JWT != nil {
		e.Use(JWT(config.JWT))
	}
}

// ApplyCasbinMiddleware 应用Casbin中间件
func ApplyCasbinMiddleware(e *echo.Echo, enforce interface{}, config *CasbinConfig) {
	if config != nil && config.Enabled {
		if enforcer, ok := enforce.(*casbin.Enforcer); ok {
			e.Use(Casbin(enforcer, config))
		}
	}
}
