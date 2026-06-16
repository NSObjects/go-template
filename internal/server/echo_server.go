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
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/NSObjects/go-template/internal/configs"
	"github.com/NSObjects/go-template/internal/server/middlewares"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// Server owns the Echo HTTP server lifecycle and system routes.
type Server struct {
	server *echo.Echo
	api    *echo.Group
	config *ServerConfig
	cfg    configs.Config
}

// Echo returns the underlying Echo instance for HTTP adapter tests and
// framework-level integration.
func (s *Server) Echo() *echo.Echo {
	return s.server
}

// API returns the root API route group used by boot to register business routes.
func (s *Server) API() *echo.Group {
	return s.api
}

// New creates an Echo-backed HTTP server.
func New(cfg configs.Config) *Server {
	s := &Server{
		server: echo.New(),
		config: FromAppConfig(cfg),
		cfg:    cfg,
	}

	// 配置服务器
	s.setupServer()
	s.loadMiddleware()
	s.registerRouter()

	return s
}

// setupServer 配置服务器基础设置
func (s *Server) setupServer() {
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
func (s *Server) loadMiddleware() {
	// 创建中间件配置
	config := s.createMiddlewareConfig()

	// 应用基础中间件
	middlewares.ApplyMiddlewares(s.server, config)
}

// createMiddlewareConfig 创建中间件配置
func (s *Server) createMiddlewareConfig() *middlewares.MiddlewareConfig {
	// 默认不启用 JWT；业务接入认证时显式开启。
	jwtConfig := middlewares.CreateJWTConfig(
		s.cfg.JWT.Secret,
		s.cfg.JWT.SkipPaths,
		s.cfg.JWT.Enabled,
	)

	return &middlewares.MiddlewareConfig{
		EnableRecovery: true,
		EnableLogger:   true,
		EnableGzip:     true,
		EnableCORS:     false,
		EnableJWT:      jwtConfig.Enabled,
		LoggerFormat:   "method=${method}, uri=${uri}, status=${status}, latency=${latency_human}\n",
		JWT:            jwtConfig,
	}
}

// registerRouter 注册路由
func (s *Server) registerRouter() {
	// 创建API路由组
	s.api = s.server.Group("/api")

	// 注册系统路由
	s.registerSystemRoutes(s.api)
}

// registerSystemRoutes 注册系统路由
func (s *Server) registerSystemRoutes(g *echo.Group) {
	// 健康检查
	g.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// 系统信息
	g.GET("/info", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"name":    "go-template",
			"version": "1.0.0",
			"time":    time.Now().Format(time.RFC3339),
		})
	})
}

// Run starts the HTTP server and blocks until ctx is cancelled or startup fails.
func (s *Server) Run(ctx context.Context) error {
	if ctx == nil {
		return errors.New("server run: nil context")
	}

	port := s.config.Port
	if port == "" {
		port = DefaultServerConfig().Port
	}

	errCh := make(chan error, 1)
	go func() {
		s.server.Logger.Infof("Starting server on %s", port)
		err := s.server.Start(port)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
	}

	s.server.Logger.Info("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
	defer cancel()

	if err := s.server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown server: %w", err)
	}

	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case <-shutdownCtx.Done():
		return fmt.Errorf("wait for server shutdown: %w", shutdownCtx.Err())
	}

	s.server.Logger.Info("Server exited")
	return nil
}
