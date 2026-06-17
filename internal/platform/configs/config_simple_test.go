package configs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigDefaults(t *testing.T) {
	cfg := Normalize(Config{})

	assert.Equal(t, DefaultAppName, cfg.App.Name)
	assert.Equal(t, DefaultAppVersion, cfg.App.Version)
	assert.Equal(t, DefaultPort, cfg.System.Port)
	assert.Equal(t, OnlineLevel, cfg.System.Level)

	assert.Equal(t, LogFormatJSON, cfg.Log.Format)
	assert.Equal(t, LogOutputStdout, cfg.Log.Output)
	assert.False(t, cfg.Log.Caller)

	assert.False(t, cfg.JWT.Enabled)
	assert.Equal(t, "", cfg.JWT.Secret)
	assert.Equal(t, []string{"/api/health", "/api/info", "/api/ready"}, cfg.JWT.SkipPaths)

	assert.False(t, cfg.MySQL.Enabled)
	assert.Equal(t, "", cfg.MySQL.DSN)
	assert.Equal(t, DefaultMySQLMaxOpenConns, cfg.MySQL.MaxOpenConns)
	assert.Equal(t, DefaultMySQLMaxIdleConns, cfg.MySQL.MaxIdleConns)
	assert.Equal(t, DefaultMySQLConnMaxLifetimeSeconds, cfg.MySQL.ConnMaxLifetimeSeconds)
	assert.Equal(t, DefaultCapabilityTimeoutSeconds, cfg.MySQL.PingTimeoutSeconds)

	assert.False(t, cfg.Redis.Enabled)
	assert.Equal(t, "", cfg.Redis.Address)
	assert.Equal(t, DefaultRedisDB, cfg.Redis.DB)
	assert.Equal(t, DefaultCapabilityTimeoutSeconds, cfg.Redis.PingTimeoutSeconds)

	assert.False(t, cfg.MongoDB.Enabled)
	assert.Equal(t, "", cfg.MongoDB.URI)
	assert.Equal(t, DefaultCapabilityTimeoutSeconds, cfg.MongoDB.ConnectTimeoutSeconds)
	assert.Equal(t, DefaultCapabilityTimeoutSeconds, cfg.MongoDB.PingTimeoutSeconds)

	assert.False(t, cfg.Tracing.Enabled)
	assert.Equal(t, "", cfg.Tracing.Endpoint)
	assert.Equal(t, DefaultTracingProtocol, cfg.Tracing.Protocol)
	assert.Equal(t, float64(0), cfg.Tracing.SampleRatio)
	assert.Equal(t, DefaultTracingShutdownTimeoutSeconds, cfg.Tracing.ShutdownTimeoutSeconds)
}

func TestAppConfig(t *testing.T) {
	cfg := AppConfig{
		Name:    "billing-api",
		Version: "2026.06.17",
	}

	assert.Equal(t, "billing-api", cfg.Name)
	assert.Equal(t, "2026.06.17", cfg.Version)
}

func TestSystemConfig(t *testing.T) {
	cfg := SystemConfig{
		Port: "8080",
	}

	assert.Equal(t, "8080", cfg.Port)
}

func TestLogConfig(t *testing.T) {
	cfg := LogConfig{
		Format: LogFormatConsole,
		Output: LogOutputStderr,
		Caller: true,
	}

	assert.Equal(t, LogFormatConsole, cfg.Format)
	assert.Equal(t, LogOutputStderr, cfg.Output)
	assert.True(t, cfg.Caller)
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

func TestCapabilityConfig(t *testing.T) {
	cfg := Config{
		MySQL: MySQLConfig{
			Enabled: true,
			DSN:     "user:pass@tcp(localhost:3306)/app?parseTime=true",
		},
		Redis: RedisConfig{
			Enabled: true,
			Address: "localhost:6379",
		},
		MongoDB: MongoDBConfig{
			Enabled: false,
		},
		Tracing: TracingConfig{
			Enabled: false,
		},
	}

	assert.True(t, cfg.MySQL.Enabled)
	assert.True(t, cfg.Redis.Enabled)
	assert.False(t, cfg.MongoDB.Enabled)
	assert.False(t, cfg.Tracing.Enabled)
}

func TestValidateRejectsEnabledJWTWithoutSecret(t *testing.T) {
	err := Validate(Config{
		JWT: JWTConfig{
			Enabled: true,
		},
	})

	assert.Error(t, err)
}

func TestValidateRejectsBlankAppIdentity(t *testing.T) {
	err := Validate(Config{
		App: AppConfig{
			Name:    " ",
			Version: "dev",
		},
	})

	assert.Error(t, err)

	err = Validate(Config{
		App: AppConfig{
			Name:    "go-template",
			Version: " ",
		},
	})

	assert.Error(t, err)
}
