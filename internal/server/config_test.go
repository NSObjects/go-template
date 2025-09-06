/*
 * Server Configuration Tests
 * 服务器配置测试用例
 */

package server

import (
	"testing"
	"time"

	"github.com/NSObjects/go-template/internal/configs"
	"github.com/stretchr/testify/assert"
)

func TestDefaultServerConfig(t *testing.T) {
	config := DefaultServerConfig()

	assert.NotNil(t, config)
	assert.Equal(t, ":8080", config.Port)
	assert.Equal(t, 30*time.Second, config.ReadTimeout)
	assert.Equal(t, 30*time.Second, config.WriteTimeout)
	assert.Equal(t, 120*time.Second, config.IdleTimeout)
	assert.Equal(t, 10*time.Second, config.ShutdownTimeout)
	assert.True(t, config.HideBanner)
	assert.False(t, config.Debug)
}

func TestFromAppConfig(t *testing.T) {
	tests := []struct {
		name     string
		config   configs.Config
		expected *ServerConfig
	}{
		{
			name: "debug mode",
			config: configs.Config{
				System: configs.SystemConfig{
					Port:  ":9090",
					Level: 1, // debug
				},
			},
			expected: &ServerConfig{
				Port:            ":9090",
				ReadTimeout:     30 * time.Second,
				WriteTimeout:    30 * time.Second,
				IdleTimeout:     120 * time.Second,
				ShutdownTimeout: 10 * time.Second,
				HideBanner:      true,
				Debug:           true,
			},
		},
		{
			name: "production mode",
			config: configs.Config{
				System: configs.SystemConfig{
					Port:  ":8080",
					Level: 2, // online
				},
			},
			expected: &ServerConfig{
				Port:            ":8080",
				ReadTimeout:     30 * time.Second,
				WriteTimeout:    30 * time.Second,
				IdleTimeout:     120 * time.Second,
				ShutdownTimeout: 10 * time.Second,
				HideBanner:      true,
				Debug:           false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FromAppConfig(tt.config)

			assert.Equal(t, tt.expected.Port, result.Port)
			assert.Equal(t, tt.expected.ReadTimeout, result.ReadTimeout)
			assert.Equal(t, tt.expected.WriteTimeout, result.WriteTimeout)
			assert.Equal(t, tt.expected.IdleTimeout, result.IdleTimeout)
			assert.Equal(t, tt.expected.ShutdownTimeout, result.ShutdownTimeout)
			assert.Equal(t, tt.expected.HideBanner, result.HideBanner)
			assert.Equal(t, tt.expected.Debug, result.Debug)
		})
	}
}

func TestServerConfigFields(t *testing.T) {
	config := &ServerConfig{
		Port:            ":3000",
		ReadTimeout:     15 * time.Second,
		WriteTimeout:    15 * time.Second,
		IdleTimeout:     60 * time.Second,
		ShutdownTimeout: 5 * time.Second,
		HideBanner:      false,
		Debug:           true,
	}

	assert.Equal(t, ":3000", config.Port)
	assert.Equal(t, 15*time.Second, config.ReadTimeout)
	assert.Equal(t, 15*time.Second, config.WriteTimeout)
	assert.Equal(t, 60*time.Second, config.IdleTimeout)
	assert.Equal(t, 5*time.Second, config.ShutdownTimeout)
	assert.False(t, config.HideBanner)
	assert.True(t, config.Debug)
}
