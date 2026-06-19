package server

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/NSObjects/go-template/internal/platform/configs"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	assert.NotNil(t, config)
	assert.Equal(t, ":9322", config.Port)
	assert.Equal(t, 30*time.Second, config.ReadTimeout)
	assert.Equal(t, 30*time.Second, config.WriteTimeout)
	assert.Equal(t, 120*time.Second, config.IdleTimeout)
	assert.Equal(t, 10*time.Second, config.ShutdownTimeout)
	assert.True(t, config.HideBanner)
}

func TestFromAppConfigUsesSystemPort(t *testing.T) {
	result := FromAppConfig(configs.Config{
		System: configs.SystemConfig{
			Port:  ":9090",
			Level: configs.DebugLevel,
		},
	})

	assertDefaultDurations(t, result)
	assert.Equal(t, ":9090", result.Port)
	assert.True(t, result.HideBanner)
}

func TestFromAppConfigOnlineMode(t *testing.T) {
	result := FromAppConfig(configs.Config{
		System: configs.SystemConfig{
			Port:  ":9322",
			Level: configs.OnlineLevel,
		},
	})

	assertDefaultDurations(t, result)
	assert.Equal(t, ":9322", result.Port)
	assert.True(t, result.HideBanner)
}

func TestFromAppConfigUsesDefaults(t *testing.T) {
	result := FromAppConfig(configs.Config{})

	assert.Equal(t, configs.DefaultPort, result.Port)
	assertDefaultDurations(t, result)
	assert.True(t, result.HideBanner)
}

func TestServerConfigFields(t *testing.T) {
	config := &Config{
		Port:            ":3000",
		ReadTimeout:     15 * time.Second,
		WriteTimeout:    15 * time.Second,
		IdleTimeout:     60 * time.Second,
		ShutdownTimeout: 5 * time.Second,
		HideBanner:      false,
	}

	assert.Equal(t, ":3000", config.Port)
	assert.Equal(t, 15*time.Second, config.ReadTimeout)
	assert.Equal(t, 15*time.Second, config.WriteTimeout)
	assert.Equal(t, 60*time.Second, config.IdleTimeout)
	assert.Equal(t, 5*time.Second, config.ShutdownTimeout)
	assert.False(t, config.HideBanner)
}

func assertDefaultDurations(t *testing.T, config *Config) {
	t.Helper()

	assert.Equal(t, 30*time.Second, config.ReadTimeout)
	assert.Equal(t, 30*time.Second, config.WriteTimeout)
	assert.Equal(t, 120*time.Second, config.IdleTimeout)
	assert.Equal(t, 10*time.Second, config.ShutdownTimeout)
}
