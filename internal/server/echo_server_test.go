/*
 * Echo Server Tests
 * Echo服务器测试用例
 */

package server

import (
	"testing"
	"time"

	"github.com/NSObjects/go-template/internal/api/service"
	"github.com/NSObjects/go-template/internal/configs"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRegisterRouter 模拟路由注册器
type MockRegisterRouter struct {
	mock.Mock
}

func (m *MockRegisterRouter) RegisterRouter(s *echo.Group, middlewareFunc ...echo.MiddlewareFunc) {
	m.Called(s, middlewareFunc)
}

func TestEchoServer_Server(t *testing.T) {
	// 创建模拟参数
	params := Params{
		Routes:   []service.RegisterRouter{},
		Enforcer: nil,
		Cfg:      configs.Config{},
		Store:    &configs.Store{},
	}

	server := NewEchoServer(params)

	assert.NotNil(t, server)
	assert.NotNil(t, server.Server())
	assert.IsType(t, &echo.Echo{}, server.Server())
}

func TestEchoServer_setupServer(t *testing.T) {
	server := &EchoServer{
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

func TestEchoServer_createMiddlewareConfig(t *testing.T) {
	store := &configs.Store{}

	server := &EchoServer{
		server: echo.New(),
		config: DefaultServerConfig(),
		store:  store,
	}

	config := server.createMiddlewareConfig()

	assert.NotNil(t, config)
	assert.True(t, config.EnableRecovery)
	assert.True(t, config.EnableLogger)
	assert.True(t, config.EnableGzip)
	assert.True(t, config.EnableCORS)
	assert.False(t, config.EnableJWT) // 默认情况下JWT应该被禁用
	assert.False(t, config.EnableCasbin)
	assert.NotNil(t, config.JWT)
	assert.NotNil(t, config.Casbin)
}

func TestEchoServer_registerSystemRoutes(t *testing.T) {
	server := &EchoServer{
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
	hasRoutesRoute := false
	hasInfoRoute := false

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
	assert.True(t, hasRoutesRoute, "Routes route should be registered")
	assert.True(t, hasInfoRoute, "Info route should be registered")
}

func TestEchoServer_registerRouter(t *testing.T) {
	// 创建模拟路由注册器
	mockRouter := new(MockRegisterRouter)
	mockRouter.On("RegisterRouter", mock.AnythingOfType("*echo.Group"), mock.AnythingOfType("[]echo.MiddlewareFunc")).Return()

	server := &EchoServer{
		server:  echo.New(),
		config:  DefaultServerConfig(),
		routers: []service.RegisterRouter{mockRouter},
	}

	server.registerRouter()

	// 验证路由注册器被调用
	mockRouter.AssertExpectations(t)
}

func TestEchoServer_Run(t *testing.T) {
	server := &EchoServer{
		server: echo.New(),
		config: &ServerConfig{
			Port:            ":0", // 使用随机端口
			ReadTimeout:     1 * time.Second,
			WriteTimeout:    1 * time.Second,
			IdleTimeout:     1 * time.Second,
			ShutdownTimeout: 1 * time.Second,
		},
	}

	// 在goroutine中运行服务器，避免阻塞测试
	go func() {
		server.Run(":0")
	}()

	// 等待服务器启动
	time.Sleep(100 * time.Millisecond)

	// 验证服务器配置
	assert.NotNil(t, server.server)
}

func TestEchoServer_NewEchoServer(t *testing.T) {
	// 创建模拟参数
	cfg := configs.Config{
		System: configs.SystemConfig{
			Port:  ":8080",
			Level: 1,
		},
		JWT: configs.JWTConfig{
			Secret:    "test-secret",
			SkipPaths: []string{"/api/health"},
		},
	}

	store := &configs.Store{}
	// 注意：这里需要根据实际的Store实现来设置配置
	// store.Set(cfg)

	params := Params{
		Routes:   []service.RegisterRouter{},
		Enforcer: nil,
		Cfg:      cfg,
		Store:    store,
	}

	server := NewEchoServer(params)

	assert.NotNil(t, server)
	assert.NotNil(t, server.server)
	assert.NotNil(t, server.config)
	assert.Equal(t, cfg, server.cfg)
	assert.Equal(t, store, server.store)
}

func TestEchoServer_SystemRoutes(t *testing.T) {
	server := &EchoServer{
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
	for _, route := range routes {
		if route.Path == "/api/health" || route.Path == "/api/routes" || route.Path == "/api/info" {
			hasSystemRoutes = true
			break
		}
	}
	assert.True(t, hasSystemRoutes, "System routes should be registered")
}

func TestEchoServer_Config(t *testing.T) {
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
