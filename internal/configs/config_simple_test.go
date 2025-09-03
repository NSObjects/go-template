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

	// 测试数据库配置默认值
	assert.Equal(t, "", cfg.Mysql.Host)
	assert.Equal(t, "", cfg.Mysql.Port)
	assert.Equal(t, "", cfg.Mysql.Database)

	// 测试日志配置默认值
	assert.Equal(t, "", cfg.Log.Level)
	assert.Equal(t, "", cfg.Log.Format)

	// 测试JWT配置默认值
	assert.Equal(t, "", cfg.JWT.Secret)
	assert.Equal(t, 0, cfg.JWT.Expire)

	// 测试CORS配置默认值
	assert.Nil(t, cfg.CORS.AllowOrigins)
	assert.Nil(t, cfg.CORS.AllowHeaders)
	assert.Nil(t, cfg.CORS.AllowMethods)
	assert.False(t, cfg.CORS.AllowCredentials)
}

func TestSystemConfig(t *testing.T) {
	cfg := SystemConfig{
		Port: "8080",
		Env:  "test",
	}

	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, "test", cfg.Env)
}

func TestMysqlConfig(t *testing.T) {
	cfg := MysqlConfig{
		Host:         "localhost",
		Port:         "3306",
		Database:     "testdb",
		User:         "testuser",
		Password:     "testpass",
		MaxIdleConns: 10,
		MaxOpenConns: 100,
	}

	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, "3306", cfg.Port)
	assert.Equal(t, "testdb", cfg.Database)
	assert.Equal(t, "testuser", cfg.User)
	assert.Equal(t, "testpass", cfg.Password)
	assert.Equal(t, 10, cfg.MaxIdleConns)
	assert.Equal(t, 100, cfg.MaxOpenConns)
}

func TestLogConfig(t *testing.T) {
	cfg := LogConfig{
		Level:  "info",
		Format: "json",
		Console: ConsoleSinkConfig{
			Format: "color",
			Output: "stdout",
		},
		File: FileSinkConfig{
			Filename:   "app.log",
			MaxSize:    100,
			MaxBackups: 3,
			MaxAge:     7,
			Compress:   true,
			Format:     "json",
		},
	}

	assert.Equal(t, "info", cfg.Level)
	assert.Equal(t, "json", cfg.Format)
	assert.Equal(t, "color", cfg.Console.Format)
	assert.Equal(t, "stdout", cfg.Console.Output)
	assert.Equal(t, "app.log", cfg.File.Filename)
	assert.Equal(t, 100, cfg.File.MaxSize)
	assert.Equal(t, 3, cfg.File.MaxBackups)
	assert.Equal(t, 7, cfg.File.MaxAge)
	assert.True(t, cfg.File.Compress)
	assert.Equal(t, "json", cfg.File.Format)
}

func TestJWTConfig(t *testing.T) {
	cfg := JWTConfig{
		Secret:    "test-secret",
		Expire:    3600,
		SkipPaths: []string{"/api/health", "/api/login"},
	}

	assert.Equal(t, "test-secret", cfg.Secret)
	assert.Equal(t, 3600, cfg.Expire)
	assert.Equal(t, []string{"/api/health", "/api/login"}, cfg.SkipPaths)
}

func TestCORSConfig(t *testing.T) {
	cfg := CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowCredentials: true,
	}

	assert.Equal(t, []string{"http://localhost:3000"}, cfg.AllowOrigins)
	assert.Equal(t, []string{"Content-Type", "Authorization"}, cfg.AllowHeaders)
	assert.Equal(t, []string{"GET", "POST", "PUT", "DELETE"}, cfg.AllowMethods)
	assert.True(t, cfg.AllowCredentials)
}

func TestCasbinConfig(t *testing.T) {
	cfg := CasbinConfig{
		Model:     "rbac",
		ModelFile: "rbac_model.conf",
	}

	assert.Equal(t, "rbac", cfg.Model)
	assert.Equal(t, "rbac_model.conf", cfg.ModelFile)
}

func TestKafkaConfig(t *testing.T) {
	cfg := KafkaConfig{
		ClientID: "echo-admin",
		Topic:    "user-events",
		Brokers:  []string{"localhost:9092"},
	}

	assert.Equal(t, "echo-admin", cfg.ClientID)
	assert.Equal(t, "user-events", cfg.Topic)
	assert.Equal(t, []string{"localhost:9092"}, cfg.Brokers)
}

func TestMongodbConfig(t *testing.T) {
	cfg := Mongodb{
		Host: "localhost",
		Port: "27017",
	}

	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, "27017", cfg.Port)
}

func TestElasticsearchSinkConfig(t *testing.T) {
	cfg := ElasticsearchSinkConfig{
		Index: "echo-admin-logs",
	}

	assert.Equal(t, "echo-admin-logs", cfg.Index)
}

func TestLokiSinkConfig(t *testing.T) {
	cfg := LokiSinkConfig{
		Labels: map[string]string{"app": "echo-admin"},
	}

	assert.Equal(t, map[string]string{"app": "echo-admin"}, cfg.Labels)
}
