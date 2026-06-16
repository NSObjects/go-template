/*
 * Echo Server Tests
 * Echo服务器测试用例
 */

package server

import (
	"context"
	"testing"
	"time"

	"github.com/NSObjects/go-template/internal/configs"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestServerEcho(t *testing.T) {
	server := New(configs.Config{})

	assert.NotNil(t, server)
	assert.NotNil(t, server.Echo())
	assert.IsType(t, &echo.Echo{}, server.Echo())
}

func TestServerSetupServer(t *testing.T) {
	server := &Server{
		server: echo.New(),
		config: DefaultServerConfig(),
	}

	server.setupServer()

	// 验证服务器配置
	assert.NotNil(t, server.server.Validator)
	assert.NotNil(t, server.server.HTTPErrorHandler)
	assert.Equal(t, server.config.HideBanner, server.server.HideBanner)
	assert.Equal(t, server.config.Debug, server.server.Debug)
	assert.Equal(t, server.config.ReadTimeout, server.server.Server.ReadTimeout)
	assert.Equal(t, server.config.WriteTimeout, server.server.Server.WriteTimeout)
	assert.Equal(t, server.config.IdleTimeout, server.server.Server.IdleTimeout)
}

func TestServerCreateMiddlewareConfig(t *testing.T) {
	server := &Server{
		server: echo.New(),
		config: DefaultServerConfig(),
		cfg:    configs.Config{},
	}

	config := server.createMiddlewareConfig()

	assert.NotNil(t, config)
	assert.True(t, config.EnableRecovery)
	assert.True(t, config.EnableLogger)
	assert.True(t, config.EnableGzip)
	assert.False(t, config.EnableCORS)
	assert.False(t, config.EnableJWT) // 默认情况下JWT应该被禁用
	assert.NotNil(t, config.JWT)
}

func TestServerCreateMiddlewareConfigEnablesJWTFromConfig(t *testing.T) {
	server := &Server{
		server: echo.New(),
		config: DefaultServerConfig(),
		cfg: configs.Config{
			JWT: configs.JWTConfig{
				Enabled:   true,
				Secret:    "test-secret",
				SkipPaths: []string{"/api/health"},
			},
		},
	}

	config := server.createMiddlewareConfig()

	assert.True(t, config.EnableJWT)
	assert.NotNil(t, config.JWT)
	assert.True(t, config.JWT.Enabled)
	assert.Equal(t, []byte("test-secret"), config.JWT.SigningKey)
	assert.Equal(t, []string{"/api/health"}, config.JWT.SkipPaths)
}

func TestServerRegisterSystemRoutes(t *testing.T) {
	server := &Server{
		server: echo.New(),
		config: DefaultServerConfig(),
	}

	// 创建测试路由组
	apiGroup := server.server.Group("/api")
	server.registerSystemRoutes(apiGroup)

	// 验证路由已注册
	routes := server.server.Routes()
	assert.NotEmpty(t, routes)

	// 验证至少包含系统路由
	hasHealthRoute := false
	hasInfoRoute := false
	hasRoutesRoute := false

	for _, route := range routes {
		if route.Path == "/api/health" && route.Method == "GET" {
			hasHealthRoute = true
		}
		if route.Path == "/api/routes" && route.Method == "GET" {
			hasRoutesRoute = true
		}
		if route.Path == "/api/info" && route.Method == "GET" {
			hasInfoRoute = true
		}
	}

	assert.True(t, hasHealthRoute, "Health route should be registered")
	assert.True(t, hasInfoRoute, "Info route should be registered")
	assert.False(t, hasRoutesRoute, "Routes introspection route should not be registered by default")
}

func TestServerRunReturnsStartupError(t *testing.T) {
	server := &Server{
		server: echo.New(),
		config: &ServerConfig{
			Port:            "invalid-address",
			ReadTimeout:     1 * time.Second,
			WriteTimeout:    1 * time.Second,
			IdleTimeout:     1 * time.Second,
			ShutdownTimeout: 1 * time.Second,
		},
	}

	err := server.Run(context.Background())
	assert.Error(t, err)
}

func TestServerRunRejectsNilContext(t *testing.T) {
	server := New(configs.Config{})

	err := server.Run(nil)
	if err == nil {
		t.Fatal("Run(nil) error = nil, want nil context error")
	}
	assert.Contains(t, err.Error(), "nil context")
}

func TestServerAPIGroupRegistersBusinessRoutes(t *testing.T) {
	server := New(configs.Config{})

	server.API().GET("/ping", func(c echo.Context) error {
		return c.NoContent(204)
	})

	hasPingRoute := false
	for _, route := range server.Echo().Routes() {
		if route.Method == "GET" && route.Path == "/api/ping" {
			hasPingRoute = true
			break
		}
	}
	assert.True(t, hasPingRoute, "API group should register routes under /api")
}

func TestServerNew(t *testing.T) {
	// 创建模拟参数
	cfg := configs.Config{
		System: configs.SystemConfig{
			Port:  ":8080",
			Level: 1,
		},
		JWT: configs.JWTConfig{
			Enabled:   true,
			Secret:    "test-secret",
			SkipPaths: []string{"/api/health"},
		},
	}

	server := New(cfg)

	assert.NotNil(t, server)
	assert.NotNil(t, server.server)
	assert.NotNil(t, server.config)
	assert.NotNil(t, server.api)
	assert.Equal(t, cfg, server.cfg)
}

func TestServerSystemRoutes(t *testing.T) {
	server := &Server{
		server: echo.New(),
		config: DefaultServerConfig(),
	}

	// 注册系统路由
	apiGroup := server.server.Group("/api")
	server.registerSystemRoutes(apiGroup)

	// 测试路由是否已注册
	routes := server.server.Routes()
	assert.NotEmpty(t, routes)

	// 验证系统路由存在
	hasSystemRoutes := false
	hasRoutesRoute := false
	for _, route := range routes {
		if route.Path == "/api/health" || route.Path == "/api/info" {
			hasSystemRoutes = true
		}
		if route.Path == "/api/routes" {
			hasRoutesRoute = true
		}
	}
	assert.True(t, hasSystemRoutes, "System routes should be registered")
	assert.False(t, hasRoutesRoute, "Routes introspection route should not be registered by default")
}

func TestServerConfig(t *testing.T) {
	config := DefaultServerConfig()

	// 测试配置字段
	assert.Equal(t, ":8080", config.Port)
	assert.Equal(t, 30*time.Second, config.ReadTimeout)
	assert.Equal(t, 30*time.Second, config.WriteTimeout)
	assert.Equal(t, 120*time.Second, config.IdleTimeout)
	assert.Equal(t, 10*time.Second, config.ShutdownTimeout)
	assert.True(t, config.HideBanner)
	assert.False(t, config.Debug)

	// 测试修改配置
	config.Port = ":9090"
	config.Debug = true
	config.HideBanner = false

	assert.Equal(t, ":9090", config.Port)
	assert.True(t, config.Debug)
	assert.False(t, config.HideBanner)
}
