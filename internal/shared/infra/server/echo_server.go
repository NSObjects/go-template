/*
 * Echo HTTP Server
 * 基于Echo框架的HTTP服务器实现
 *
 * Created by lintao on 2023/7/26
 * Copyright © 2020-2024 LINTAO. All rights reserved.
 */

package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NSObjects/go-template/internal/configs"
	"github.com/NSObjects/go-template/internal/pkg/validator"
	"github.com/NSObjects/go-template/internal/shared/infra/server/middlewares"
	"github.com/NSObjects/go-template/internal/shared/ports/resp"
	"github.com/NSObjects/go-template/internal/user/adapters"
	"github.com/casbin/casbin/v2"
	"github.com/labstack/echo/v4"
	"github.com/marmotedu/errors"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

// EchoServer Echo HTTP服务器
type EchoServer struct {
	server  *echo.Echo
	config  *ServerConfig
	routers []adapters.RegisterRouter
	cfg     configs.Config
	store   *configs.Store
}

// Server 获取Echo实例
func (s *EchoServer) Server() *echo.Echo {
	return s.server
}

// Params 依赖注入参数
type Params struct {
	fx.In

	Routes      []adapters.RegisterRouter `group:"routes"`
	Enforcer    *casbin.Enforcer
	RateLimiter *middlewares.RateLimiter
	Cfg         configs.Config
	Store       *configs.Store
}

// NewEchoServer 创建Echo服务器实例
func NewEchoServer(p Params) *EchoServer {
	s := &EchoServer{
		server:  echo.New(),
		config:  FromAppConfig(p.Cfg),
		routers: p.Routes,
		cfg:     p.Cfg,
		store:   p.Store,
	}

	// 配置服务器
	s.setupServer()
	s.loadMiddleware(p.Enforcer, p.RateLimiter)
	s.registerRouter()

	return s
}

// NewRateLimiter 创建限流器
func NewRateLimiter(redis *redis.Client) *middlewares.RateLimiter {
	return middlewares.NewRateLimiter(redis)
}

// setupServer 配置服务器基础设置
func (s *EchoServer) setupServer() {
	// 设置验证器 - 使用统一的验证器实现
	customValidator := validator.NewCustomValidator()
	s.server.Validator = customValidator

	// 设置错误处理器
	s.server.HTTPErrorHandler = middlewares.ErrorHandler

	// 应用服务器配置
	s.server.HideBanner = s.config.HideBanner
	s.server.Debug = s.config.Debug

	// 设置超时
	s.server.Server.ReadTimeout = s.config.ReadTimeout
	s.server.Server.WriteTimeout = s.config.WriteTimeout
	s.server.Server.IdleTimeout = s.config.IdleTimeout
}

// loadMiddleware 加载中间件
func (s *EchoServer) loadMiddleware(enforce *casbin.Enforcer, rateLimiter *middlewares.RateLimiter) {
	// 创建中间件配置
	config := s.createMiddlewareConfig()

	// 应用基础中间件
	middlewares.ApplyMiddlewares(s.server, config)

	// 应用Casbin中间件
	middlewares.ApplyCasbinMiddleware(s.server, enforce, config.Casbin)

	// 应用限流中间件
	middlewares.ApplyRateLimitMiddleware(s.server, rateLimiter, config.RateLimit)
}

// createMiddlewareConfig 创建中间件配置
func (s *EchoServer) createMiddlewareConfig() *middlewares.MiddlewareConfig {
	cur := s.store.Current()

	// 创建JWT配置 - 禁用JWT用于演示
	jwtConfig := middlewares.CreateJWTConfig(
		cur.JWT.Secret,
		cur.JWT.SkipPaths,
		false, // 禁用JWT
	)

	// 调试日志
	fmt.Printf("DEBUG: JWT Config - Enabled: %v, SkipPaths: %v\n", jwtConfig.Enabled, jwtConfig.SkipPaths)

	// 创建Casbin配置
	casbinConfig := middlewares.CreateCasbinConfig(
		false, // 默认禁用Casbin
		[]string{
			"/api/health",
			"/api/info",
			"/api/login",
			"/api/users",
		},
		[]string{"root", "admin"},
	)

	// 创建限流配置
	rateLimitConfig := middlewares.DefaultRateLimitConfig()

	return &middlewares.MiddlewareConfig{
		EnableRecovery:  true,
		EnableLogger:    true,
		EnableGzip:      true,
		EnableCORS:      true,
		EnableJWT:       jwtConfig.Enabled,
		EnableCasbin:    casbinConfig.Enabled,
		EnableRateLimit: rateLimitConfig.Enabled,
		LoggerFormat:    "method=${method}, uri=${uri}, status=${status}, latency=${latency_human}\n",
		JWT:             jwtConfig,
		Casbin:          casbinConfig,
		RateLimit:       rateLimitConfig,
	}
}

// registerRouter 注册路由
func (s *EchoServer) registerRouter() {
	// 创建API路由组
	apiGroup := s.server.Group("/api")

	// 注册业务路由
	for _, router := range s.routers {
		router.RegisterRouter(apiGroup)
	}

	// 注册系统路由
	s.registerSystemRoutes(apiGroup)
}

// registerSystemRoutes 注册系统路由
func (s *EchoServer) registerSystemRoutes(g *echo.Group) {
	// 健康检查
	g.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// 路由信息
	g.GET("/routes", func(c echo.Context) error {
		return resp.ListDataResponse(c, s.server.Routes(), int64(len(s.server.Routes())))
	})

	// 系统信息
	g.GET("/info", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"name":    "echo-admin",
			"version": "1.0.0",
			"time":    time.Now().Format(time.RFC3339),
		})
	})
}

// Run 启动服务器
func (s *EchoServer) Run(port string) {
	if port == "" {
		port = s.config.Port
	}

	// 启动服务器
	go func() {
		s.server.Logger.Infof("Starting server on %s", port)
		if err := s.server.Start(port); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.server.Logger.Fatal("Failed to start server", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	s.server.Logger.Info("Shutting down server...")

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		s.server.Logger.Fatal("Server forced to shutdown", err)
	}

	s.server.Logger.Info("Server exited")
}
