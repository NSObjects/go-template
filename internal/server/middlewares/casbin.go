/*
 * Casbin Middleware
 * 基于Casbin的权限控制中间件
 *
 * Created by lintao on 2024/1/4
 * Copyright © 2020-2024 LINTAO. All rights reserved.
 */

package middlewares

import (
	"net/http"

	"github.com/NSObjects/go-template/internal/code"
	"github.com/NSObjects/go-template/internal/utils"
	"github.com/casbin/casbin/v3"
	"github.com/golang-jwt/jwt/v5"
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

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if lo.Contains(config.SkipPaths, c.Path()) {
				return next(c)
			}

			user, err := casbinUser(c)
			if err != nil {
				return errors.WrapC(err, code.ErrPermissionDenied, "权限不足")
			}

			allowed, err := casbinAllowed(enforce, config, c, user)
			if err != nil {
				return errors.WrapC(err, code.ErrPermissionDenied, "权限检查失败")
			}
			if !allowed {
				return errors.WrapC(errors.New(http.StatusText(http.StatusForbidden)), code.ErrPermissionDenied, "权限不足")
			}

			return next(c)
		}
	}
}

func casbinUser(c echo.Context) (string, error) {
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
}

func casbinAllowed(enforce *casbin.Enforcer, config *CasbinConfig, c echo.Context, user string) (bool, error) {
	if lo.Contains(config.AdminUsers, user) {
		return true, nil
	}

	return enforce.Enforce(user, c.Path(), c.Request().Method)
}

// CreateCasbinConfig 从应用配置创建Casbin配置
func CreateCasbinConfig(enabled bool, skipPaths []string, adminUsers []string) *CasbinConfig {
	return &CasbinConfig{
		Enabled:    enabled,
		SkipPaths:  skipPaths,
		AdminUsers: adminUsers,
	}
}
