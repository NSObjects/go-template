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
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NSObjects/go-template/internal/api/service"
	"github.com/NSObjects/go-template/internal/configs"
	"github.com/NSObjects/go-template/internal/resp"
	"github.com/NSObjects/go-template/internal/server/middlewares"
	"github.com/casbin/casbin/v2"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/marmotedu/errors"
	"go.uber.org/fx"
)

// EchoServer Echo HTTP服务器
type EchoServer struct {
	server  *echo.Echo
	config  *ServerConfig
	routers []service.RegisterRouter
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

	Routes   []service.RegisterRouter `group:"routes"`
	Enforcer *casbin.Enforcer
	Cfg      configs.Config
	Store    *configs.Store
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
	s.loadMiddleware(p.Enforcer)
	s.registerRouter()

	return s
}

// setupServer 配置服务器基础设置
func (s *EchoServer) setupServer() {
	// 设置验证器
	s.server.Validator = &middlewares.Validator{Validator: validator.New()}

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
func (s *EchoServer) loadMiddleware(enforce *casbin.Enforcer) {
	// 创建中间件配置
	config := s.createMiddlewareConfig()

	// 应用基础中间件
	middlewares.ApplyMiddlewares(s.server, config)

	// 应用Casbin中间件
	middlewares.ApplyCasbinMiddleware(s.server, enforce, config.Casbin)
}

// createMiddlewareConfig 创建中间件配置
func (s *EchoServer) createMiddlewareConfig() *middlewares.MiddlewareConfig {
	cur := s.store.Current()

	// 创建JWT配置
	jwtConfig := middlewares.CreateJWTConfig(
		cur.JWT.Secret,
		cur.JWT.SkipPaths,
		len(cur.JWT.SkipPaths) > 0, // 如果有跳过路径，则启用JWT
	)

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

	return &middlewares.MiddlewareConfig{
		EnableRecovery: true,
		EnableLogger:   true,
		EnableGzip:     true,
		EnableCORS:     true,
		EnableJWT:      jwtConfig.Enabled,
		EnableCasbin:   casbinConfig.Enabled,
		LoggerFormat:   "method=${method}, uri=${uri}, status=${status}, latency=${latency_human}\n",
		JWT:            jwtConfig,
		Casbin:         casbinConfig,
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
		return resp.ListDataResponse(s.server.Routes(), int64(len(s.server.Routes())), c)
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
