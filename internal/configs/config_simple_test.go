package configs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigDefaults(t *testing.T) {
	cfg := Normalize(Config{})

	assert.Equal(t, DefaultPort, cfg.System.Port)
	assert.Equal(t, DefaultEnv, cfg.System.Env)
	assert.Equal(t, OnlineLevel, cfg.System.Level)

	assert.False(t, cfg.JWT.Enabled)
	assert.Equal(t, "", cfg.JWT.Secret)
	assert.Equal(t, []string{"/api/health", "/api/info"}, cfg.JWT.SkipPaths)
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
		SkipPaths: []string{"/api/health"},
	}

	assert.True(t, cfg.Enabled)
	assert.Equal(t, "test-secret", cfg.Secret)
	assert.Equal(t, []string{"/api/health"}, cfg.SkipPaths)
}

func TestValidateRejectsEnabledJWTWithoutSecret(t *testing.T) {
	err := Validate(Config{
		JWT: JWTConfig{
			Enabled: true,
		},
	})

	assert.Error(t, err)
}
