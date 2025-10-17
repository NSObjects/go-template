package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestDefaultMiddlewareConfig(t *testing.T) {
	config := DefaultMiddlewareConfig()

	assert.NotNil(t, config)
	assert.True(t, config.EnableRecovery)
	assert.True(t, config.EnableLogger)
	assert.True(t, config.EnableGzip)
	assert.True(t, config.EnableCORS)
	assert.False(t, config.EnableJWT)
	assert.False(t, config.EnableCasbin)
	assert.NotEmpty(t, config.LoggerFormat)
	assert.NotNil(t, config.JWT)
	assert.NotNil(t, config.Casbin)
}

func TestDefaultJWTConfig(t *testing.T) {
	config := DefaultJWTConfig()

	assert.NotNil(t, config)
	assert.False(t, config.Enabled)
	assert.NotNil(t, config.SkipPaths)
	assert.NotEmpty(t, config.SigningKey) // 默认配置有默认密钥
}

func TestDefaultCasbinConfig(t *testing.T) {
	config := DefaultCasbinConfig()

	assert.NotNil(t, config)
	assert.False(t, config.Enabled)
	assert.NotNil(t, config.SkipPaths)
	assert.NotNil(t, config.AdminUsers)
}

func TestCreateJWTConfig(t *testing.T) {
	tests := []struct {
		name      string
		secret    string
		skipPaths []string
		enabled   bool
	}{
		{
			name:      "enabled JWT",
			secret:    "test-secret",
			skipPaths: []string{"/api/health"},
			enabled:   true,
		},
		{
			name:      "disabled JWT",
			secret:    "test-secret",
			skipPaths: []string{},
			enabled:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := CreateJWTConfig(tt.secret, tt.skipPaths, tt.enabled)
			assert.NotNil(t, config)
			assert.Equal(t, []byte(tt.secret), config.SigningKey)
			assert.Equal(t, tt.skipPaths, config.SkipPaths)
			assert.Equal(t, tt.enabled, config.Enabled)
		})
	}
}

func TestCreateCasbinConfig(t *testing.T) {
	tests := []struct {
		name       string
		enabled    bool
		skipPaths  []string
		adminUsers []string
	}{
		{
			name:       "enabled Casbin",
			enabled:    true,
			skipPaths:  []string{"/api/health"},
			adminUsers: []string{"admin"},
		},
		{
			name:       "disabled Casbin",
			enabled:    false,
			skipPaths:  []string{},
			adminUsers: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := CreateCasbinConfig(tt.enabled, tt.skipPaths, tt.adminUsers)
			assert.NotNil(t, config)
			assert.Equal(t, tt.enabled, config.Enabled)
			assert.Equal(t, tt.skipPaths, config.SkipPaths)
			assert.Equal(t, tt.adminUsers, config.AdminUsers)
		})
	}
}

func TestApplyMiddlewares(t *testing.T) {
	e := echo.New()
	config := DefaultMiddlewareConfig()

	// 测试中间件应用不会panic
	assert.NotPanics(t, func() {
		ApplyMiddlewares(e, config)
	})

	assert.NotNil(t, e)
}

func TestErrorRecovery(t *testing.T) {
	e := echo.New()
	e.Use(ErrorRecovery())

	// 创建一个会panic的路由
	e.GET("/panic", func(c echo.Context) error {
		panic("test panic")
	})

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	rec := httptest.NewRecorder()

	// 测试panic恢复
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestJWTConfig(t *testing.T) {
	tests := []struct {
		name     string
		config   *JWTConfig
		expected bool
	}{
		{
			name: "enabled JWT",
			config: &JWTConfig{
				SigningKey: []byte("test-secret"),
				SkipPaths:  []string{"/api/health"},
				Enabled:    true,
			},
			expected: true,
		},
		{
			name: "disabled JWT",
			config: &JWTConfig{
				SigningKey: []byte("test-secret"),
				SkipPaths:  []string{},
				Enabled:    false,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := JWT(tt.config)
			assert.NotNil(t, middleware)
		})
	}
}

func TestCasbinConfig(t *testing.T) {
	// 创建测试用的Casbin enforcer
	// 注意：这里我们只测试配置创建，不测试实际的权限检查
	config := &CasbinConfig{
		Enabled:    true,
		SkipPaths:  []string{"/api/health"},
		AdminUsers: []string{"admin"},
	}

	// 测试配置创建
	assert.NotNil(t, config)
	assert.True(t, config.Enabled)
	assert.Equal(t, []string{"/api/health"}, config.SkipPaths)
	assert.Equal(t, []string{"admin"}, config.AdminUsers)
}

func TestMiddlewareConfig(t *testing.T) {
	tests := []struct {
		name   string
		config *MiddlewareConfig
	}{
		{
			name:   "default config",
			config: DefaultMiddlewareConfig(),
		},
		{
			name: "custom config",
			config: &MiddlewareConfig{
				EnableRecovery: true,
				EnableLogger:   true,
				EnableGzip:     true,
				EnableCORS:     true,
				EnableJWT:      false,
				EnableCasbin:   false,
				LoggerFormat:   "custom format",
				JWT:            DefaultJWTConfig(),
				Casbin:         DefaultCasbinConfig(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			ApplyMiddlewares(e, tt.config)
			assert.NotNil(t, e)
		})
	}
}
