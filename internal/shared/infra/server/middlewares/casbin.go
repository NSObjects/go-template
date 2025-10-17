/*
 * Casbin Middleware
 * 基于Casbin的权限控制中间件
 *
 * Created by lintao on 2024/1/4
 * Copyright © 2020-2024 LINTAO. All rights reserved.
 */

package middlewares

import (
	"github.com/NSObjects/go-template/internal/pkg/code"
	"github.com/NSObjects/go-template/internal/pkg/utils"
	"github.com/casbin/casbin/v2"
	"github.com/golang-jwt/jwt/v5"
	casbin_mw "github.com/labstack/echo-contrib/casbin"
	"github.com/labstack/echo/v4"
	"github.com/marmotedu/errors"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

// CasbinConfig Casbin中间件配置
type CasbinConfig struct {
	// 是否启用
	Enabled bool
	// 跳过路径
	SkipPaths []string
	// 管理员用户
	AdminUsers []string
}

// DefaultCasbinConfig 默认Casbin配置
func DefaultCasbinConfig() *CasbinConfig {
	return &CasbinConfig{
		Enabled: false,
		SkipPaths: []string{
			"/api/health",
			"/api/info",
			"/api/login",
			"/api/users",
		},
		AdminUsers: []string{"root", "admin"},
	}
}

// Casbin Casbin权限控制中间件
func Casbin(enforce *casbin.Enforcer, config *CasbinConfig) echo.MiddlewareFunc {
	if !config.Enabled || enforce == nil {
		// 如果Casbin未启用或enforcer为空，返回空中间件
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return next
		}
	}

	return casbin_mw.MiddlewareWithConfig(casbin_mw.Config{
		Enforcer: enforce,
		Skipper: func(c echo.Context) bool {
			path := c.Path()
			return lo.Contains(config.SkipPaths, path)
		},
		ErrorHandler: func(c echo.Context, internal error, proposedStatus int) error {
			return errors.WrapC(internal, code.ErrPermissionDenied, "权限不足")
		},
		UserGetter: func(c echo.Context) (string, error) {
			token, ok := c.Get("user").(*jwt.Token)
			if !ok {
				return "", errors.WrapC(errors.New("token is nil"), code.ErrSignatureInvalid, "JWT签名无效")
			}
			if token == nil {
				return "", nil
			}

			user, ok := token.Claims.(*utils.JwtCustomClaims)
			if !ok {
				return "", errors.WrapC(errors.New("invalid token claims type"), code.ErrSignatureInvalid, "JWT签名无效")
			}
			if user == nil {
				return "", nil
			}
			if user.Admin {
				return "root", nil
			}

			return cast.ToString(user.ID), nil
		},
		EnforceHandler: func(c echo.Context, user string) (bool, error) {
			// 检查是否为管理员用户
			if lo.Contains(config.AdminUsers, user) {
				return true, nil
			}

			// 获取请求路径和方法
			path := c.Path()
			method := c.Request().Method

			// 使用Casbin进行权限检查
			allowed, err := enforce.Enforce(user, path, method)
			if err != nil {
				return false, errors.WrapC(err, code.ErrPermissionDenied, "权限检查失败")
			}

			return allowed, nil
		},
	})
}

// CreateCasbinConfig 从应用配置创建Casbin配置
func CreateCasbinConfig(enabled bool, skipPaths []string, adminUsers []string) *CasbinConfig {
	return &CasbinConfig{
		Enabled:    enabled,
		SkipPaths:  skipPaths,
		AdminUsers: adminUsers,
	}
}
