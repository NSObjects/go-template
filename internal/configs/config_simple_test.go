package configs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigDefaults(t *testing.T) {
	// 测试默认配置结构
	cfg := Config{}

	// 测试系统配置默认值
	assert.Equal(t, "", cfg.System.Port)
	assert.Equal(t, "", cfg.System.Env)

	// 测试JWT配置默认值
	assert.False(t, cfg.JWT.Enabled)
	assert.Equal(t, "", cfg.JWT.Secret)
	assert.Equal(t, 0, cfg.JWT.Expire)
}

func TestSystemConfig(t *testing.T) {
	cfg := SystemConfig{
		Port: "8080",
		Env:  "test",
	}

	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, "test", cfg.Env)
}

func TestJWTConfig(t *testing.T) {
	cfg := JWTConfig{
		Enabled:   true,
		Secret:    "test-secret",
		Expire:    3600,
		SkipPaths: []string{"/api/health", "/api/login"},
	}

	assert.True(t, cfg.Enabled)
	assert.Equal(t, "test-secret", cfg.Secret)
	assert.Equal(t, 3600, cfg.Expire)
	assert.Equal(t, []string{"/api/health", "/api/login"}, cfg.SkipPaths)
}
