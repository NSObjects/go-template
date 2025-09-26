package project

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// 嵌入的模板内容
var templates = map[string]string{
	"go.mod.tmpl": `module {{.ModulePath}}

go 1.24.0

require (
	github.com/IBM/sarama v1.46.0
	github.com/casbin/casbin/v2 v2.122.0
	github.com/casbin/gorm-adapter/v3 v3.36.0
	github.com/fsnotify/fsnotify v1.9.0
	github.com/go-playground/validator/v10 v10.27.0
	github.com/go-sql-driver/mysql v1.9.3
	github.com/golang-jwt/jwt/v5 v5.3.0
	github.com/google/uuid v1.6.0
	github.com/hashicorp/consul/api v1.32.1
	github.com/labstack/echo-contrib v0.17.4
	github.com/labstack/echo-jwt/v4 v4.3.1
	github.com/labstack/echo/v4 v4.13.4
	github.com/lmittmann/tint v1.1.2
	github.com/marmotedu/errors v1.0.2
	github.com/novalagung/gubrak v1.0.0
	github.com/prometheus/client_golang v1.22.0
	github.com/redis/go-redis/v9 v9.14.0
	github.com/samber/lo v1.51.0
	github.com/spf13/cast v1.10.0
	github.com/spf13/cobra v1.10.1
	github.com/spf13/viper v1.21.0
	github.com/stretchr/testify v1.11.1
	go.etcd.io/etcd/client/v3 v3.6.4
	go.mongodb.org/mongo-driver v1.17.4
	go.uber.org/fx v1.24.0
	golang.org/x/crypto v0.42.0
	golang.org/x/tools v0.37.0
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/driver/mysql v1.6.0
	gorm.io/gen v0.3.27
	gorm.io/gorm v1.30.5
	gorm.io/plugin/dbresolver v1.6.2
)`,

	"main.go.tmpl": `package main

import (
	toolcmd "{{.ModulePath}}/cmd"
)

func main() {
	toolcmd.Execute()
}`,

	"README.md.tmpl": `# {{.ProjectName}}

{{.ProjectName}} 是一个基于 Go 的 Web 应用项目，使用 Echo 框架构建。

## 功能特性

- 🚀 基于 Echo 框架的高性能 Web 服务
- 🔐 JWT 认证和授权
- 🗄️ 数据库支持 (MySQL, PostgreSQL, SQLite)
- 📊 监控和指标收集
- 🔧 配置管理和热重载
- 📝 自动 API 文档生成
- 🧪 完整的测试覆盖

## 快速开始

### 环境要求

- Go 1.24+
- MySQL 8.0+ (或其他支持的数据库)
- Redis (可选)

### 安装依赖

` + "```bash" + `
go mod tidy
` + "```" + `

### 配置

1. 复制配置文件：
` + "```bash" + `
cp env.example .env
` + "```" + `

2. 编辑 .env 文件，配置数据库连接等信息

### 运行

` + "```bash" + `
# 开发模式
make dev

# 生产模式
make build
./{{.ProjectName}}
` + "```" + `

## 项目结构

` + "```" + `
.
├── cmd/                    # 命令行工具
├── configs/               # 配置文件
├── internal/              # 内部包
│   ├── api/              # API 层
│   ├── cache/            # 缓存
│   ├── code/             # 错误码
│   ├── configs/          # 配置管理
│   ├── health/           # 健康检查
│   ├── log/              # 日志
│   ├── metrics/          # 监控指标
│   ├── middleware/       # 中间件
│   ├── resp/             # 响应处理
│   ├── server/           # 服务器
│   ├── utils/            # 工具函数
│   └── validator/        # 验证器
├── k8s/                  # Kubernetes 配置
├── scripts/              # 脚本
└── sql/                  # SQL 文件
` + "```" + `

## 开发指南

### 添加新的 API

1. 在 internal/api/biz/ 中定义业务逻辑
2. 在 internal/api/service/ 中实现服务层
3. 在 internal/api/data/ 中实现数据层
4. 在 cmd/ 中注册路由

### 数据库迁移

` + "```bash" + `
# 生成迁移文件
make migrate-create name=create_users_table

# 运行迁移
make migrate-up
` + "```" + `

## 部署

### Docker

` + "```bash" + `
# 构建镜像
docker build -t {{.ProjectName}} .

# 运行容器
docker run -p 8080:8080 {{.ProjectName}}
` + "```" + `

### Kubernetes

` + "```bash" + `
kubectl apply -f k8s/
` + "```" + `

## 许可证

MIT License`,

	"Makefile.tmpl": `# {{.ProjectName}} Makefile

.PHONY: help build run test clean dev dev-setup lint fmt vet

# 默认目标
help: ## 显示帮助信息
	@echo "可用的命令:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

# 构建相关
build: ## 构建应用
	@echo "Building {{.ProjectName}}..."
	@go build -o bin/{{.ProjectName}} .

run: build ## 构建并运行应用
	@echo "Running {{.ProjectName}}..."
	@./bin/{{.ProjectName}}

# 开发相关
dev: ## 开发模式运行
	@echo "Starting {{.ProjectName}} in development mode..."
	@air

dev-setup: ## 设置开发环境
	@echo "Setting up development environment..."
	@go mod download
	@go mod tidy
	@echo "Development environment setup complete!"

# 测试相关
test: ## 运行测试
	@echo "Running tests..."
	@go test -v ./...

test-coverage: ## 运行测试并生成覆盖率报告
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

# 代码质量
lint: ## 运行 linter
	@echo "Running linter..."
	@golangci-lint run

fmt: ## 格式化代码
	@echo "Formatting code..."
	@go fmt ./...

vet: ## 运行 go vet
	@echo "Running go vet..."
	@go vet ./...

# 清理
clean: ## 清理构建文件
	@echo "Cleaning up..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html

# 数据库相关
migrate-up: ## 运行数据库迁移
	@echo "Running database migrations..."
	@# 在这里添加迁移命令

migrate-down: ## 回滚数据库迁移
	@echo "Rolling back database migrations..."
	@# 在这里添加回滚命令

# Docker 相关
docker-build: ## 构建 Docker 镜像
	@echo "Building Docker image..."
	@docker build -t {{.ProjectName}} .

docker-run: ## 运行 Docker 容器
	@echo "Running Docker container..."
	@docker run -p 8080:8080 {{.ProjectName}}

# 安装工具
install-tools: ## 安装开发工具
	@echo "Installing development tools..."
	@go install github.com/cosmtrek/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`,

	"cmd/fx.go.tmpl": `/*
 * Created by lintao on 2023/7/27 上午10:04
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package cmd

import (
	"context"

	"log/slog"

	"{{.ModulePath}}/internal/api/biz"
	"{{.ModulePath}}/internal/api/data"
	"{{.ModulePath}}/internal/api/service"
	"{{.ModulePath}}/internal/configs"
	"{{.ModulePath}}/internal/log"
	"{{.ModulePath}}/internal/server"

	"go.uber.org/fx"
)

func Run(cfg string) {
	fx.New(
		fx.Module("config", fx.Provide(func() (configs.Config, *configs.Store) {
			merged, store := configs.Bootstrap(cfg)
			return merged, store
		})),
		fx.Module("log", fx.Provide(func(cfg configs.Config) log.Logger {
			return log.NewLogger(cfg)
		})),
		fx.Module("data", data.Model, data.CasbinModule),
		fx.Module("biz", biz.Model),
		fx.Module("service", service.Model),
		fx.Module("server", fx.Provide(server.NewEchoServer)),
		fx.Invoke(func(lifecycle fx.Lifecycle, s *server.EchoServer, cfg configs.Config, logger log.Logger) {
			// 测试日志输出
			logger.Info("Application starting", slog.String("port", cfg.System.Port))

			lifecycle.Append(
				fx.Hook{
					OnStart: func(context.Context) error {
						logger.Info("Server starting", slog.String("port", cfg.System.Port))
						go s.Run(cfg.System.Port)
						return nil
					},
					OnStop: func(context.Context) error {
						logger.Info("Server stopping")
						return nil
					},
				})
		}),
	).Run()
}`,

	"cmd/root.go.tmpl": `// Package cmd /*
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "",
	Short: "A brief description of your application",
	Long: ` + "`" + `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.` + "`" + `,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		Run(cfgFile)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "configs/config.toml",
		"config file (default is $HOME/.echo-admin.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}`,

	"internal/server/echo_server.go.tmpl": `/*
 * Echo HTTP Server
 * 基于Echo框架的HTTP服务器实现
 *
 * Created by lintao on 2023/7/26
 * Copyright © 2020-2024 LINTAO. All rights reserved.
 */

package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"{{.ModulePath}}/internal/api/service"
	"{{.ModulePath}}/internal/configs"
	"{{.ModulePath}}/internal/resp"
	"{{.ModulePath}}/internal/server/middlewares"
	"github.com/casbin/casbin/v2"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/marmotedu/errors"
	"go.uber.org/fx"
)

// EchoServer Echo HTTP服务器
type EchoServer struct {
	server  *echo.Echo
	config  *ServerConfig
	routers []service.RegisterRouter
	cfg     configs.Config
	store   *configs.Store
}

// Server 获取Echo实例
func (s *EchoServer) Server() *echo.Echo {
	return s.server
}

// Params 依赖注入参数
type Params struct {
	fx.In

	Routes   []service.RegisterRouter ` + "`" + `group:"routes"` + "`" + `
	Enforcer *casbin.Enforcer
	Cfg      configs.Config
	Store    *configs.Store
}

// NewEchoServer 创建Echo服务器实例
func NewEchoServer(p Params) *EchoServer {
	s := &EchoServer{
		server:  echo.New(),
		config:  FromAppConfig(p.Cfg),
		routers: p.Routes,
		cfg:     p.Cfg,
		store:   p.Store,
	}

	// 配置服务器
	s.setupServer()
	s.loadMiddleware(p.Enforcer)
	s.registerRouter()

	return s
}

// setupServer 配置服务器基础设置
func (s *EchoServer) setupServer() {
	// 设置验证器
	s.server.Validator = &middlewares.CustomValidator{Validator: validator.New()}

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
func (s *EchoServer) loadMiddleware(enforce *casbin.Enforcer) {
	// 创建中间件配置
	config := s.createMiddlewareConfig()

	// 应用基础中间件
	middlewares.ApplyMiddlewares(s.server, config)

	// 应用Casbin中间件
	middlewares.ApplyCasbinMiddleware(s.server, enforce, config.Casbin)
}

// createMiddlewareConfig 创建中间件配置
func (s *EchoServer) createMiddlewareConfig() *middlewares.MiddlewareConfig {
	cur := s.store.Current()

	// 创建JWT配置 - 禁用JWT用于演示
	jwtConfig := middlewares.CreateJWTConfig(
		cur.JWT.Secret,
		cur.JWT.SkipPaths,
		false, // 禁用JWT
	)

	// 调试日志
	fmt.Printf("DEBUG: JWT Config - Enabled: %v, SkipPaths: %v\n", jwtConfig.Enabled, jwtConfig.SkipPaths)

	// 创建Casbin配置
	casbinConfig := middlewares.CreateCasbinConfig(
		false, // 默认禁用Casbin
		[]string{
			"/api/health",
			"/api/info",
			"/api/login",
			"/api/users",
		},
		[]string{"root", "admin"},
	)

	return &middlewares.MiddlewareConfig{
		EnableRecovery: true,
		EnableLogger:   true,
		EnableGzip:     true,
		EnableCORS:     true,
		EnableJWT:      jwtConfig.Enabled,
		EnableCasbin:   casbinConfig.Enabled,
		LoggerFormat:   "method=${method}, uri=${uri}, status=${status}, latency=${latency_human}\n",
		JWT:            jwtConfig,
		Casbin:         casbinConfig,
	}
}

// registerRouter 注册路由
func (s *EchoServer) registerRouter() {
	// 创建API路由组
	apiGroup := s.server.Group("/api")

	// 注册业务路由
	for _, router := range s.routers {
		router.RegisterRouter(apiGroup)
	}

	// 注册系统路由
	s.registerSystemRoutes(apiGroup)
}

// registerSystemRoutes 注册系统路由
func (s *EchoServer) registerSystemRoutes(g *echo.Group) {
	// 健康检查
	g.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// 路由信息
	g.GET("/routes", func(c echo.Context) error {
		return resp.ListDataResponse(s.server.Routes(), int64(len(s.server.Routes())), c)
	})

	// 系统信息
	g.GET("/info", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"name":    "{{.ProjectName}}",
			"version": "1.0.0",
			"time":    time.Now().Format(time.RFC3339),
		})
	})
}

// Run 启动服务器
func (s *EchoServer) Run(port string) {
	if port == "" {
		port = s.config.Port
	}

	// 启动服务器
	go func() {
		s.server.Logger.Infof("Starting server on %s", port)
		if err := s.server.Start(port); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.server.Logger.Fatal("Failed to start server", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	s.server.Logger.Info("Shutting down server...")

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		s.server.Logger.Fatal("Server forced to shutdown", err)
	}

	s.server.Logger.Info("Server exited")
}`,

	"internal/configs/config.go.tmpl": `/*
 * Created by lintao on 2023/7/18 下午3:56
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package configs

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Environment string

const (
	ENVIRONMENT Environment = "RUN_ENVIRONMENT"
)

type EnvironmentType int

const (
	Dev EnvironmentType = iota
	Docker
	Test
)

var runContext = map[string]EnvironmentType{
	"":       Dev,
	"docker": Docker,
	"test":   Test,
}

type Level int8

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel Level = iota + 1
	// ProsecutionLevel InfoLevel is the default logging priority.
	ProsecutionLevel
)

var (
	Mysql  MysqlConfig
	System SystemConfig
	Log    LogConfig
	Mgo    Mongodb
	JWT    JWTConfig
)

type Config struct {
	Mysql   MysqlConfig        ` + "`" + `mapstructure:"mysql"` + "`" + `
	System  SystemConfig       ` + "`" + `mapstructure:"system"` + "`" + `
	Log     LogConfig          ` + "`" + `mapstructure:"log"` + "`" + `
	Mongodb Mongodb            ` + "`" + `mapstructure:"mongodb"` + "`" + `
	Redis   RedisConfig        ` + "`" + `mapstructure:"redis"` + "`" + `
	JWT     JWTConfig          ` + "`" + `mapstructure:"jwt"` + "`" + `
	CORS    CORSConfig         ` + "`" + `mapstructure:"cors"` + "`" + `
	Casbin  CasbinConfig       ` + "`" + `mapstructure:"casbin"` + "`" + `
	Kafka   KafkaConfig        ` + "`" + `mapstructure:"kafka"` + "`" + `
	Etcd    EtcdClientConfig   ` + "`" + `mapstructure:"etcd"` + "`" + `
	Consul  ConsulClientConfig ` + "`" + `mapstructure:"consul"` + "`" + `
}

type SystemConfig struct {
	Port  string ` + "`" + `mapstructure:"port"` + "`" + `
	Level Level  ` + "`" + `mapstructure:"level"` + "`" + `
	Env   string ` + "`" + `mapstructure:"env"` + "`" + `
}

type RedisConfig struct {
	Host     string ` + "`" + `mapstructure:"host"` + "`" + `
	Port     string ` + "`" + `mapstructure:"port"` + "`" + `
	Password string ` + "`" + `mapstructure:"password"` + "`" + `
	Database int    ` + "`" + `mapstructure:"database"` + "`" + `
}

type LogConfig struct {
	Level  string ` + "`" + `mapstructure:"level"` + "`" + `
	Format string ` + "`" + `mapstructure:"format"` + "`" + `

	Console ConsoleSinkConfig ` + "`" + `mapstructure:"console"` + "`" + `
	File    FileSinkConfig    ` + "`" + `mapstructure:"file"` + "`" + `

	Elasticsearch ElasticsearchSinkConfig ` + "`" + `mapstructure:"elasticsearch"` + "`" + `
	Loki          LokiSinkConfig          ` + "`" + `mapstructure:"loki"` + "`" + `
}

type ConsoleSinkConfig struct {
	Format string ` + "`" + `mapstructure:"format"` + "`" + `
	Output string ` + "`" + `mapstructure:"output"` + "`" + `
}

type FileSinkConfig struct {
	Filename   string ` + "`" + `mapstructure:"filename"` + "`" + `
	MaxSize    int    ` + "`" + `mapstructure:"max_size"` + "`" + `
	MaxBackups int    ` + "`" + `mapstructure:"max_backups"` + "`" + `
	MaxAge     int    ` + "`" + `mapstructure:"max_age"` + "`" + `
	Compress   bool   ` + "`" + `mapstructure:"compress"` + "`" + `
	Format     string ` + "`" + `mapstructure:"format"` + "`" + `
}

type ElasticsearchSinkConfig struct {
	URL     string        ` + "`" + `mapstructure:"url"` + "`" + `
	Index   string        ` + "`" + `mapstructure:"index"` + "`" + `
	Timeout time.Duration ` + "`" + `mapstructure:"timeout"` + "`" + `
}

type LokiSinkConfig struct {
	URL     string            ` + "`" + `mapstructure:"url"` + "`" + `
	Labels  map[string]string ` + "`" + `mapstructure:"labels"` + "`" + `
	Timeout time.Duration     ` + "`" + `mapstructure:"timeout"` + "`" + `
}

type MysqlConfig struct {
	DockerHost   string ` + "`" + `mapstructure:"docker_host"` + "`" + `
	Host         string ` + "`" + `mapstructure:"host"` + "`" + `
	Port         string ` + "`" + `mapstructure:"port"` + "`" + `
	User         string ` + "`" + `mapstructure:"user"` + "`" + `
	Password     string ` + "`" + `mapstructure:"password"` + "`" + `
	MaxOpenConns int    ` + "`" + `mapstructure:"max_open_conns"` + "`" + `
	MaxIdleConns int    ` + "`" + `mapstructure:"max_idle_conns"` + "`" + `
	Database     string ` + "`" + `mapstructure:"database"` + "`" + `
}

type JWTConfig struct {
	Secret    string   ` + "`" + `mapstructure:"secret"` + "`" + `
	Expire    int      ` + "`" + `mapstructure:"expire"` + "`" + `
	SkipPaths []string ` + "`" + `mapstructure:"skip_paths"` + "`" + `
}

type Mongodb struct {
	Host     string ` + "`" + `mapstructure:"host"` + "`" + `
	Port     string ` + "`" + `mapstructure:"port"` + "`" + `
	User     string ` + "`" + `mapstructure:"user"` + "`" + `
	Password string ` + "`" + `mapstructure:"password"` + "`" + `
	DataBase string ` + "`" + `mapstructure:"database"` + "`" + `
}

type CORSConfig struct {
	AllowOrigins     []string ` + "`" + `mapstructure:"allow_origins"` + "`" + `
	AllowHeaders     []string ` + "`" + `mapstructure:"allow_headers"` + "`" + `
	AllowMethods     []string ` + "`" + `mapstructure:"allow_methods"` + "`" + `
	AllowCredentials bool     ` + "`" + `mapstructure:"allow_credentials"` + "`" + `
}

type CasbinConfig struct {
	Model     string ` + "`" + `mapstructure:"model"` + "`" + `
	ModelFile string ` + "`" + `mapstructure:"model_file"` + "`" + `
}

type KafkaConfig struct {
	Brokers  []string ` + "`" + `mapstructure:"brokers"` + "`" + `
	ClientID string   ` + "`" + `mapstructure:"client_id"` + "`" + `
	Topic    string   ` + "`" + `mapstructure:"topic"` + "`" + `
}

type EtcdClientConfig struct {
	Endpoints          []string ` + "`" + `mapstructure:"endpoints"` + "`" + `
	Key                string   ` + "`" + `mapstructure:"key"` + "`" + `
	Format             string   ` + "`" + `mapstructure:"format"` + "`" + `
	Username           string   ` + "`" + `mapstructure:"username"` + "`" + `
	Password           string   ` + "`" + `mapstructure:"password"` + "`" + `
	DialTimeoutSeconds int      ` + "`" + `mapstructure:"dial_timeout_seconds"` + "`" + `
}

type ConsulClientConfig struct {
	Address string ` + "`" + `mapstructure:"address"` + "`" + `
	Token   string ` + "`" + `mapstructure:"token"` + "`" + `
	Key     string ` + "`" + `mapstructure:"key"` + "`" + `
	Format  string ` + "`" + `mapstructure:"format"` + "`" + `
}

func InitConfig(configPath string) (err error) {
	// 设置配置文件
	viper.SetConfigFile(configPath)
	viper.SetConfigType("toml")

	// 读取配置文件
	if err = viper.ReadInConfig(); err != nil {
		return
	}

	var c Config
	if err = viper.Unmarshal(&c); err != nil {
		fmt.Println(err)
		return
	}

	Mysql = c.Mysql
	System = c.System
	Log = c.Log
	Mgo = c.Mongodb
	JWT = c.JWT
	return
}

func NewCfg(p string) Config {
	return NewCfgFrom(FileSource{Path: p})
}

// Source 抽象：配置来源（本地文件、远程配置中心等）。
// 实现方直接返回完整的 Config，便于从任意介质填充（文件/etcd/http 等）。
type Source interface {
	Load(ctx context.Context) (Config, error)
}

// WatchableSource 可选：支持热更新的配置源
// Watch 应启动监听并在变更时回调返回新的 Config
type WatchableSource interface {
	Source
	Watch(ctx context.Context, onChange func(Config)) error
}

// NewCfgFrom 通过自定义 Source 加载 Config。
func NewCfgFrom(src Source) Config {
	c, err := src.Load(context.Background())
	if err != nil {
		panic(err)
	}
	return c
}

func RunEnvironment() EnvironmentType {
	return runContext[os.Getenv(string(ENVIRONMENT))]
}`,

	"internal/log/logger.go.tmpl": `package log

import (
	"log/slog"
	"os"

	"{{.ModulePath}}/internal/configs"
)

// Logger 日志接口
type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Fatal(msg string, args ...any)
}

// logger 日志实现
type logger struct {
	*slog.Logger
}

// NewLogger 创建新的日志器
func NewLogger(cfg configs.Config) Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	var handler slog.Handler
	switch cfg.Log.Format {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, opts)
	default:
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return &logger{
		Logger: slog.New(handler),
	}
}

// Debug 调试日志
func (l *logger) Debug(msg string, args ...any) {
	l.Logger.Debug(msg, args...)
}

// Info 信息日志
func (l *logger) Info(msg string, args ...any) {
	l.Logger.Info(msg, args...)
}

// Warn 警告日志
func (l *logger) Warn(msg string, args ...any) {
	l.Logger.Warn(msg, args...)
}

// Error 错误日志
func (l *logger) Error(msg string, args ...any) {
	l.Logger.Error(msg, args...)
}

// Fatal 致命错误日志
func (l *logger) Fatal(msg string, args ...any) {
	l.Logger.Error(msg, args...)
	os.Exit(1)
}`,

	"internal/resp/response.go.tmpl": `package resp

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Response 统一响应结构
type Response struct {
	Code    int         ` + "`" + `json:"code"` + "`" + `
	Message string      ` + "`" + `json:"message"` + "`" + `
	Data    interface{} ` + "`" + `json:"data,omitempty"` + "`" + `
}

// ListResponse 列表响应结构
type ListResponse struct {
	Code    int         ` + "`" + `json:"code"` + "`" + `
	Message string      ` + "`" + `json:"message"` + "`" + `
	Data    interface{} ` + "`" + `json:"data"` + "`" + `
	Total   int64       ` + "`" + `json:"total"` + "`" + `
}

// SuccessResponse 成功响应
func SuccessResponse(data interface{}, c echo.Context) error {
	return c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

// ErrorResponse 错误响应
func ErrorResponse(code int, message string, c echo.Context) error {
	return c.JSON(code, Response{
		Code:    code,
		Message: message,
	})
}

// ListDataResponse 列表数据响应
func ListDataResponse(data interface{}, total int64, c echo.Context) error {
	return c.JSON(http.StatusOK, ListResponse{
		Code:    200,
		Message: "success",
		Data:    data,
		Total:   total,
	})
}`,

	"internal/code/code.go.tmpl": `package code

import (
	"github.com/marmotedu/errors"
)

// 定义错误码
const (
	// 通用错误码
	Success                = 200
	ErrInvalidParams       = 400
	ErrUnauthorized        = 401
	ErrForbidden          = 403
	ErrNotFound           = 404
	ErrInternalServer     = 500
	ErrServiceUnavailable = 503

	// 业务错误码
	ErrUserNotFound     = 10001
	ErrUserAlreadyExist = 10002
	ErrInvalidPassword  = 10003
	ErrTokenExpired     = 10004
	ErrTokenInvalid     = 10005
)

// 错误码映射
var codeMap = map[int]string{
	Success:                "success",
	ErrInvalidParams:       "invalid params",
	ErrUnauthorized:        "unauthorized",
	ErrForbidden:          "forbidden",
	ErrNotFound:           "not found",
	ErrInternalServer:     "internal server error",
	ErrServiceUnavailable: "service unavailable",
	ErrUserNotFound:       "user not found",
	ErrUserAlreadyExist:   "user already exist",
	ErrInvalidPassword:    "invalid password",
	ErrTokenExpired:       "token expired",
	ErrTokenInvalid:       "token invalid",
}

// GetMessage 获取错误信息
func GetMessage(code int) string {
	if msg, ok := codeMap[code]; ok {
		return msg
	}
	return "unknown error"
}

// NewError 创建错误
func NewError(code int) error {
	return errors.New(GetMessage(code))
}

// NewErrorWithMessage 创建带自定义消息的错误
func NewErrorWithMessage(code int, message string) error {
	return errors.New(message)
}`,

	"internal/api/biz/biz.go.tmpl": `package biz

import (
	"go.uber.org/fx"
)

// Model 业务逻辑模块
var Model = fx.Module("biz",
	fx.Provide(
		// 在这里添加业务逻辑提供者
	),
)

// 在这里添加业务逻辑相关的结构体和函数`,

	"internal/api/service/service.go.tmpl": `package service

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

// RegisterRouter 路由注册接口
type RegisterRouter interface {
	RegisterRouter(g *echo.Group)
}

// Model 服务层模块
var Model = fx.Module("service",
	fx.Provide(
		// 在这里添加服务层提供者
	),
)

// 在这里添加服务层相关的结构体和函数`,

	"internal/api/data/data.go.tmpl": `package data

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"{{.ModulePath}}/internal/api/data/db"
	"{{.ModulePath}}/internal/api/data/query"
	"{{.ModulePath}}/internal/configs"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// DataManager 统一的数据管理器，提供所有数据库组件的操作接口
type DataManager struct {
	// 数据源
	ds *db.DataSource

	// 查询接口
	Query *query.Query

	// 配置
	Config *configs.Config
}

// NewDataManager 创建统一的数据管理器
func NewDataManager(
	ds *db.DataSource,
	query *query.Query,
	cfg configs.Config,
) *DataManager {
	return &DataManager{
		ds:     ds,
		Query:  query,
		Config: &cfg,
	}
}

// Close 关闭所有数据库连接
func (dm *DataManager) Close() error {
	// 通过DataSource统一关闭所有连接
	return dm.ds.Stop(context.Background())
}

// Health 检查所有组件的健康状态
func (dm *DataManager) Health(ctx context.Context) map[string]error {
	health := make(map[string]error)

	// 通过DataSource获取组件状态
	status := dm.ds.GetComponentStatus(ctx)
	for component, status := range status {
		if status.Enabled {
			health[component] = status.Error
		}
	}

	return health
}

// ========== 统一的数据操作接口 ==========

// MySQL 获取MySQL数据库连接
func (dm *DataManager) MySQL() *gorm.DB {
	return dm.ds.Mysql
}

// Redis 获取Redis客户端
func (dm *DataManager) Redis() *redis.Client {
	return dm.ds.Redis
}

// Kafka 获取Kafka生产者
func (dm *DataManager) Kafka() sarama.SyncProducer {
	return dm.ds.Kafka
}

// MongoDB 获取MongoDB数据库
func (dm *DataManager) MongoDB() *mongo.Database {
	return dm.ds.Mongodb
}

// ========== 便捷操作方法 ==========

// MySQLWithContext 获取带上下文的MySQL连接
func (dm *DataManager) MySQLWithContext(ctx context.Context) *gorm.DB {
	if dm.ds.Mysql == nil {
		return nil
	}
	return dm.ds.Mysql.WithContext(ctx)
}

// RedisWithContext 获取带上下文的Redis客户端
func (dm *DataManager) RedisWithContext(ctx context.Context) *redis.Client {
	if dm.ds.Redis == nil {
		return nil
	}
	// Redis客户端本身已经支持context，直接返回
	return dm.ds.Redis
}

// SendKafkaMessage 发送Kafka消息
func (dm *DataManager) SendKafkaMessage(topic string, key, value []byte) error {
	if dm.ds.Kafka == nil {
		return fmt.Errorf("kafka producer not initialized")
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	}

	_, _, err := dm.ds.Kafka.SendMessage(msg)
	return err
}

// IsComponentEnabled 检查组件是否启用
func (dm *DataManager) IsComponentEnabled(component string) bool {
	switch component {
	case "mysql":
		return dm.ds.Mysql != nil
	case "redis":
		return dm.ds.Redis != nil
	case "kafka":
		return dm.ds.Kafka != nil
	case "mongodb":
		return dm.ds.Mongodb != nil
	default:
		return false
	}
}

var Model = fx.Options(
	fx.Provide(
		db.NewDataSource,
		NewDB,
		NewQuery,
		NewDataManager,
	),
)

func NewQuery(db *gorm.DB) *query.Query {
	// 使用生成的Query
	return query.Use(db)
}

// NewDB exposes the primary Gorm DB from the unified DataSource for DI consumers.
func NewDB(ds *db.DataSource) *gorm.DB {
	if ds == nil || ds.Mysql == nil {
		panic("mysql data source is not initialized")
	}
	return ds.Mysql
}`,

	"configs/config.toml.tmpl": `[system]
port = ":8080"
level = 1
env = "dev"

[log]
level = "debug"
format = "text"

[log.console]
format = "text"
output = "stdout"

[log.file]
filename = "logs/app.log"
max_size = 100
max_backups = 3
max_age = 7
compress = true
format = "json"

[mysql]
docker_host = "localhost"
host = "localhost"
port = "3306"
user = "root"
password = "password"
database = "{{.PackageName}}"
max_open_conns = 100
max_idle_conns = 10

[redis]
host = "localhost"
port = "6379"
password = ""
database = 0

[jwt]
secret = "your-secret-key"
expire = 3600
skip_paths = ["/api/health", "/api/info", "/api/login"]

[cors]
allow_origins = ["*"]
allow_headers = ["*"]
allow_methods = ["*"]
allow_credentials = true

[casbin]
model = "rbac"
model_file = "configs/rbac_model.conf"

[kafka]
brokers = ["localhost:9092"]
client_id = "{{.PackageName}}"
topic = "{{.PackageName}}-events"

[etcd]
endpoints = ["localhost:2379"]
key = "/{{.PackageName}}/config"
format = "toml"
dial_timeout_seconds = 5

[consul]
address = "localhost:8500"
token = ""
key = "{{.PackageName}}/config"
format = "toml"`,

	"env.example.tmpl": `# {{.ProjectName}} Environment Variables

# 系统配置
RUN_ENVIRONMENT=dev
PORT=8080

# 数据库配置
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_USER=root
MYSQL_PASSWORD=password
MYSQL_DATABASE={{.PackageName}}

# Redis配置
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DATABASE=0

# JWT配置
JWT_SECRET=your-secret-key
JWT_EXPIRE=3600

# 日志配置
LOG_LEVEL=debug
LOG_FORMAT=text`,

	"cmd/gen.go.tmpl": `package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate code from templates",
	Long: ` + "`" + `Generate code from templates using the modgen tool.
This command helps generate boilerplate code for your Go project.` + "`" + `,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gen called")
		// 这里可以调用代码生成逻辑
	},
}

func init() {
	rootCmd.AddCommand(genCmd)
}`,

	"internal/server/config.go.tmpl": `package server

import (
	"time"

	"{{.ModulePath}}/internal/configs"
)

// ServerConfig 服务器配置
type ServerConfig struct {
	Port             string
	Debug            bool
	HideBanner       bool
	ReadTimeout      time.Duration
	WriteTimeout     time.Duration
	IdleTimeout      time.Duration
	ShutdownTimeout  time.Duration
}

// FromAppConfig 从应用配置创建服务器配置
func FromAppConfig(cfg configs.Config) *ServerConfig {
	return &ServerConfig{
		Port:             cfg.System.Port,
		Debug:            cfg.System.Level == configs.DebugLevel,
		HideBanner:       false,
		ReadTimeout:      30 * time.Second,
		WriteTimeout:     30 * time.Second,
		IdleTimeout:      120 * time.Second,
		ShutdownTimeout:  30 * time.Second,
	}
}`,

	"internal/server/middlewares/middleware.go.tmpl": `package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// MiddlewareConfig 中间件配置
type MiddlewareConfig struct {
	EnableRecovery bool
	EnableLogger   bool
	EnableGzip     bool
	EnableCORS     bool
	EnableJWT      bool
	EnableCasbin   bool
	LoggerFormat   string
	JWT            *JWTConfig
	Casbin         *CasbinConfig
}

// JWTConfig JWT配置
type JWTConfig struct {
	Enabled   bool
	Secret    string
	SkipPaths []string
}

// CasbinConfig Casbin配置
type CasbinConfig struct {
	Enabled      bool
	SkipPaths    []string
	AdminRoles   []string
}

// ApplyMiddlewares 应用中间件
func ApplyMiddlewares(e *echo.Echo, config *MiddlewareConfig) {
	if config.EnableRecovery {
		e.Use(Recovery())
	}
	if config.EnableLogger {
		e.Use(Logger())
	}
	if config.EnableGzip {
		e.Use(Gzip())
	}
	if config.EnableCORS {
		e.Use(CORS())
	}
}

// Recovery 恢复中间件
func Recovery() echo.MiddlewareFunc {
	return middleware.Recover()
}

// Logger 日志中间件
func Logger() echo.MiddlewareFunc {
	return middleware.Logger()
}

// Gzip 压缩中间件
func Gzip() echo.MiddlewareFunc {
	return middleware.Gzip()
}

// CORS 跨域中间件
func CORS() echo.MiddlewareFunc {
	return middleware.CORS()
}`,

	"internal/server/middlewares/config.go.tmpl": `package middlewares

import (
	"github.com/go-playground/validator/v10"
)

// CustomValidator 自定义验证器
type CustomValidator struct {
	Validator *validator.Validate
}

// Validate 验证结构体
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

// CreateJWTConfig 创建JWT配置
func CreateJWTConfig(secret string, skipPaths []string, enabled bool) *JWTConfig {
	return &JWTConfig{
		Enabled:   enabled,
		Secret:    secret,
		SkipPaths: skipPaths,
	}
}

// CreateCasbinConfig 创建Casbin配置
func CreateCasbinConfig(enabled bool, skipPaths, adminRoles []string) *CasbinConfig {
	return &CasbinConfig{
		Enabled:    enabled,
		SkipPaths:  skipPaths,
		AdminRoles: adminRoles,
	}
}`,

	"internal/server/middlewares/error.go.tmpl": `package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// ErrorHandler 错误处理器
func ErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	message := "Internal Server Error"

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		message = he.Message.(string)
	}

	c.JSON(code, map[string]interface{}{
		"code":    code,
		"message": message,
	})
}`,

	"internal/server/middlewares/jwt.go.tmpl": `package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo-jwt/v4"
)

// JWT JWT中间件
func JWT(config *JWTConfig) echo.MiddlewareFunc {
	if !config.Enabled {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return next
		}
	}

	return echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.Secret),
		Skipper: func(c echo.Context) bool {
			path := c.Request().URL.Path
			for _, skipPath := range config.SkipPaths {
				if path == skipPath {
					return true
				}
			}
			return false
		},
	})
}`,

	"internal/server/middlewares/casbin.go.tmpl": `package middlewares

import (
	"github.com/casbin/casbin/v2"
	"github.com/labstack/echo/v4"
)

// ApplyCasbinMiddleware 应用Casbin中间件
func ApplyCasbinMiddleware(e *echo.Echo, enforcer *casbin.Enforcer, config *CasbinConfig) {
	if !config.Enabled || enforcer == nil {
		return
	}

	e.Use(Casbin(enforcer, config))
}

// Casbin Casbin权限控制中间件
func Casbin(enforcer *casbin.Enforcer, config *CasbinConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 检查路径是否在跳过列表中
			path := c.Request().URL.Path
			for _, skipPath := range config.SkipPaths {
				if path == skipPath {
					return next(c)
				}
			}

			// 这里应该从JWT token中获取用户角色
			// 简化示例，实际应该从认证信息中获取
			role := "guest"

			// 检查权限
			allowed, err := enforcer.Enforce(role, path, c.Request().Method)
			if err != nil {
				return echo.NewHTTPError(500, "Permission check failed")
			}

			if !allowed {
				return echo.NewHTTPError(403, "Access denied")
			}

			return next(c)
		}
	}
}`,

	"internal/configs/bootstrap.go.tmpl": `package configs

import (
	"context"
	"fmt"
)

// Bootstrap 启动配置
func Bootstrap(path string) (Config, *Store) {
	base := NewCfg(path)
	merged := base
	ctx := context.Background()

	// 增量合并：etcd
	if len(base.Etcd.Endpoints) > 0 && base.Etcd.Key != "" {
		etcdSource := EtcdSource{
			Endpoints:          base.Etcd.Endpoints,
			Key:                base.Etcd.Key,
			Format:             base.Etcd.Format,
			Username:           base.Etcd.Username,
			Password:           base.Etcd.Password,
			DialTimeoutSeconds: base.Etcd.DialTimeoutSeconds,
		}
		if etcdCfg, err := etcdSource.Load(ctx); err == nil {
			merged = Merge(merged, etcdCfg)
		}
	}
	// 增量合并：consul
	if base.Consul.Address != "" && base.Consul.Key != "" {
		consulSource := ConsulSource{
			Address: base.Consul.Address,
			Token:   base.Consul.Token,
			Key:     base.Consul.Key,
			Format:  base.Consul.Format,
		}
		if consulCfg, err := consulSource.Load(ctx); err == nil {
			merged = Merge(merged, consulCfg)
		}
	}

	store := NewStore(merged)
	// 文件热更新（作为默认入口）
	_ = FileSource{Path: path}.Watch(ctx, func(nc Config) {
		store.Update(Merge(store.Current(), nc))
	})
	// etcd 热更新（如果配置了）
	if len(base.Etcd.Endpoints) > 0 && base.Etcd.Key != "" {
		etcdSource := EtcdSource{
			Endpoints:          base.Etcd.Endpoints,
			Key:                base.Etcd.Key,
			Format:             base.Etcd.Format,
			Username:           base.Etcd.Username,
			Password:           base.Etcd.Password,
			DialTimeoutSeconds: base.Etcd.DialTimeoutSeconds,
		}
		_ = etcdSource.Watch(ctx, func(nc Config) {
			store.Update(Merge(store.Current(), nc))
		})
	}
	// consul 热更新（如果配置了）
	if base.Consul.Address != "" && base.Consul.Key != "" {
		consulSource := ConsulSource{
			Address: base.Consul.Address,
			Token:   base.Consul.Token,
			Key:     base.Consul.Key,
			Format:  base.Consul.Format,
		}
		_ = consulSource.Watch(ctx, func(nc Config) {
			store.Update(Merge(store.Current(), nc))
		})
	}

	return merged, store
}`,

	"internal/configs/store.go.tmpl": `package configs

import (
	"sync"
)

// Store 配置存储
type Store struct {
	mu   sync.RWMutex
	cfg  Config
}

// NewStore 创建配置存储
func NewStore(cfg Config) *Store {
	return &Store{
		cfg: cfg,
	}
}

// Current 获取当前配置
func (s *Store) Current() Config {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cfg
}

// Update 更新配置
func (s *Store) Update(cfg Config) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cfg = cfg
}`,

	"internal/configs/file_source.go.tmpl": `package configs

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// FileSource 文件配置源
type FileSource struct {
	Path string
}

// Load 加载配置
func (f FileSource) Load(ctx context.Context) (Config, error) {
	// 检查文件是否存在
	if _, err := os.Stat(f.Path); os.IsNotExist(err) {
		return Config{}, fmt.Errorf("config file not found: %s", f.Path)
	}

	// 设置配置文件
	viper.SetConfigFile(f.Path)
	viper.SetConfigType("toml")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("failed to read config file: %w", err)
	}

	// 解析配置
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return config, nil
}`,

	"internal/configs/merge.go.tmpl": `package configs

// MergeConfigs 合并配置
func MergeConfigs(configs ...Config) Config {
	if len(configs) == 0 {
		return Config{}
	}

	result := configs[0]
	for i := 1; i < len(configs); i++ {
		// 这里可以实现配置合并逻辑
		// 简化示例，实际应该根据具体需求实现
	}
	
	return result
}`,

	"internal/configs/hot_reload.go.tmpl": `package configs

import (
	"context"

	"github.com/fsnotify/fsnotify"
)

// HotReload 热重载配置
func HotReload(ctx context.Context, store *Store, configPath string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	// 监听配置文件
	if err := watcher.Add(configPath); err != nil {
		return err
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				// 配置文件被修改，重新加载
				source := FileSource{Path: configPath}
				if newConfig, err := source.Load(ctx); err == nil {
					store.Update(newConfig)
				}
			}
		case err := <-watcher.Errors:
			return err
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}`,

	"internal/configs/etcd_source.go.tmpl": `package configs

import (
	"bytes"
	"context"
	"time"

	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// EtcdSource 从 etcd 读取配置，支持 json/yaml/toml 三种格式。
type EtcdSource struct {
	Endpoints          []string
	Key                string
	Format             string // json|yaml|toml（默认 toml）
	Username           string
	Password           string
	DialTimeoutSeconds int
}

func (e EtcdSource) Load(ctx context.Context) (Config, error) {
	dialTimeout := 5 * time.Second
	if e.DialTimeoutSeconds > 0 {
		dialTimeout = time.Duration(e.DialTimeoutSeconds) * time.Second
	}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   e.Endpoints,
		Username:    e.Username,
		Password:    e.Password,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		return Config{}, err
	}
	defer cli.Close()

	resp, err := cli.Get(ctx, e.Key)
	if err != nil {
		return Config{}, err
	}
	if len(resp.Kvs) == 0 {
		return Config{}, nil
	}

	v := viper.New()
	v.SetConfigType(e.Format)
	if err := v.ReadConfig(bytes.NewReader(resp.Kvs[0].Value)); err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (e EtcdSource) Watch(ctx context.Context, onChange func(Config)) error {
	dialTimeout := 5 * time.Second
	if e.DialTimeoutSeconds > 0 {
		dialTimeout = time.Duration(e.DialTimeoutSeconds) * time.Second
	}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   e.Endpoints,
		Username:    e.Username,
		Password:    e.Password,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		return err
	}
	defer cli.Close()

	watchCh := cli.Watch(ctx, e.Key)
	for watchResp := range watchCh {
		for _, event := range watchResp.Events {
			if event.Type == clientv3.EventTypePut {
				cfg, err := e.Load(ctx)
				if err == nil {
					onChange(cfg)
				}
			}
		}
	}
	return nil
}`,

	"internal/configs/consul_source.go.tmpl": `package configs

import (
	"bytes"
	"context"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
)

// ConsulSource 从 Consul 读取配置
type ConsulSource struct {
	Address string
	Token   string
	Key     string
	Format  string // json|yaml|toml（默认 toml）
}

func (c ConsulSource) Load(ctx context.Context) (Config, error) {
	client, err := api.NewClient(&api.Config{
		Address: c.Address,
		Token:   c.Token,
	})
	if err != nil {
		return Config{}, err
	}

	kv, _, err := client.KV().Get(c.Key, nil)
	if err != nil {
		return Config{}, err
	}
	if kv == nil {
		return Config{}, nil
	}

	v := viper.New()
	v.SetConfigType(c.Format)
	if err := v.ReadConfig(bytes.NewReader(kv.Value)); err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c ConsulSource) Watch(ctx context.Context, onChange func(Config)) error {
	client, err := api.NewClient(&api.Config{
		Address: c.Address,
		Token:   c.Token,
	})
	if err != nil {
		return err
	}

	queryOptions := &api.QueryOptions{
		WaitTime: 10 * time.Second,
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			kv, meta, err := client.KV().Get(c.Key, queryOptions)
			if err != nil {
				continue
			}
			if kv != nil {
				cfg, err := c.Load(ctx)
				if err == nil {
					onChange(cfg)
				}
			}
			queryOptions.WaitIndex = meta.LastIndex
		}
	}
}`,

	"configs/rbac_model.conf.tmpl": `[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act`,

	"internal/log/factory.go.tmpl": `package log

import (
	"{{.ModulePath}}/internal/configs"
)

// Factory 日志工厂
type Factory struct{}

// NewFactory 创建日志工厂
func NewFactory() *Factory {
	return &Factory{}
}

// CreateLogger 创建日志器
func (f *Factory) CreateLogger(cfg configs.Config) Logger {
	return NewLogger(cfg)
}`,

	"internal/log/console_sink.go.tmpl": `package log

import (
	"os"
)

// ConsoleSink 控制台日志输出
type ConsoleSink struct {
	format string
	output string
}

// NewConsoleSink 创建控制台日志输出
func NewConsoleSink(format, output string) *ConsoleSink {
	return &ConsoleSink{
		format: format,
		output: output,
	}
}

// Write 写入日志
func (c *ConsoleSink) Write(p []byte) (n int, err error) {
	return os.Stdout.Write(p)
}`,

	"internal/log/file_sink.go.tmpl": `package log

import (
	"os"
	"path/filepath"
)

// FileSink 文件日志输出
type FileSink struct {
	filename string
	file     *os.File
}

// NewFileSink 创建文件日志输出
func NewFileSink(filename string) *FileSink {
	return &FileSink{
		filename: filename,
	}
}

// Write 写入日志
func (f *FileSink) Write(p []byte) (n int, err error) {
	if f.file == nil {
		// 确保目录存在
		dir := filepath.Dir(f.filename)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return 0, err
		}
		
		// 打开文件
		f.file, err = os.OpenFile(f.filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return 0, err
		}
	}
	
	return f.file.Write(p)
}`,

	"internal/log/global.go.tmpl": `package log

import (
	"sync"
)

var (
	globalLogger Logger
	once         sync.Once
)

// SetGlobalLogger 设置全局日志器
func SetGlobalLogger(logger Logger) {
	once.Do(func() {
		globalLogger = logger
	})
}

// GetGlobalLogger 获取全局日志器
func GetGlobalLogger() Logger {
	return globalLogger
}`,

	"internal/log/slog.go.tmpl": `package log

import (
	"log/slog"
)

// SlogLogger slog日志器包装
type SlogLogger struct {
	*slog.Logger
}

// NewSlogLogger 创建slog日志器
func NewSlogLogger() *SlogLogger {
	return &SlogLogger{
		Logger: slog.Default(),
	}
}`,

	"internal/code/base.go.tmpl": `package code

// BaseError 基础错误接口
type BaseError interface {
	error
	Code() int
	Message() string
}

// Error 错误结构
type Error struct {
	code    int
	message string
}

// New 创建新错误
func New(code int, message string) *Error {
	return &Error{
		code:    code,
		message: message,
	}
}

// Error 实现error接口
func (e *Error) Error() string {
	return e.message
}

// Code 获取错误码
func (e *Error) Code() int {
	return e.code
}

// Message 获取错误消息
func (e *Error) Message() string {
	return e.message
}`,

	"internal/code/errors.go.tmpl": `package code

import (
	"fmt"
)

// ErrorWithCode 带错误码的错误
type ErrorWithCode struct {
	Code    int    ` + "`" + `json:"code"` + "`" + `
	Message string ` + "`" + `json:"message"` + "`" + `
}

// Error 实现error接口
func (e *ErrorWithCode) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// NewErrorWithCode 创建带错误码的错误
func NewErrorWithCode(code int, message string) *ErrorWithCode {
	return &ErrorWithCode{
		Code:    code,
		Message: message,
	}
}`,

	"internal/code/http_status.go.tmpl": `package code

import "net/http"

// HTTPStatusMap HTTP状态码映射
var HTTPStatusMap = map[int]int{
	Success:                http.StatusOK,
	ErrInvalidParams:       http.StatusBadRequest,
	ErrUnauthorized:        http.StatusUnauthorized,
	ErrForbidden:          http.StatusForbidden,
	ErrNotFound:           http.StatusNotFound,
	ErrInternalServer:     http.StatusInternalServerError,
	ErrServiceUnavailable: http.StatusServiceUnavailable,
	ErrUserNotFound:       http.StatusNotFound,
	ErrUserAlreadyExist:   http.StatusConflict,
	ErrInvalidPassword:    http.StatusUnauthorized,
	ErrTokenExpired:       http.StatusUnauthorized,
	ErrTokenInvalid:       http.StatusUnauthorized,
}

// GetHTTPStatus 获取HTTP状态码
func GetHTTPStatus(code int) int {
	if status, ok := HTTPStatusMap[code]; ok {
		return status
	}
	return http.StatusInternalServerError
}`,

	"internal/code/error_types.go.tmpl": `package code

// ErrorType 错误类型
type ErrorType int

const (
	// ErrorTypeSystem 系统错误
	ErrorTypeSystem ErrorType = iota
	// ErrorTypeBusiness 业务错误
	ErrorTypeBusiness
	// ErrorTypeValidation 验证错误
	ErrorTypeValidation
	// ErrorTypeAuth 认证错误
	ErrorTypeAuth
)

// GetErrorType 获取错误类型
func GetErrorType(code int) ErrorType {
	switch {
	case code >= 10000:
		return ErrorTypeBusiness
	case code >= 1000:
		return ErrorTypeValidation
	case code >= 400 && code < 500:
		return ErrorTypeAuth
	default:
		return ErrorTypeSystem
	}
}`,

	"internal/api/data/casbin.go.tmpl": `package data

import (
	"go.uber.org/fx"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

// CasbinModule Casbin权限控制模块
var CasbinModule = fx.Module("casbin",
	fx.Provide(NewCasbinEnforcer),
)

// NewCasbinEnforcer 创建Casbin执行器
func NewCasbinEnforcer(db *gorm.DB) (*casbin.Enforcer, error) {
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}

	enforcer, err := casbin.NewEnforcer("configs/rbac_model.conf", adapter)
	if err != nil {
		return nil, err
	}

	return enforcer, nil
}`,

	"internal/api/data/jwt.go.tmpl": `package data

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims JWT声明
type JWTClaims struct {
	UserID   string ` + "`" + `json:"user_id"` + "`" + `
	Username string ` + "`" + `json:"username"` + "`" + `
	Role     string ` + "`" + `json:"role"` + "`" + `
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT令牌
func GenerateToken(userID, username, role, secret string, expireHours int) (string, error) {
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseToken 解析JWT令牌
func ParseToken(tokenString, secret string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}`,

	"internal/api/data/db/db.go.tmpl": `package db

import (
	"context"

	"github.com/IBM/sarama"
	"{{.ModulePath}}/internal/configs"
	_ "github.com/go-sql-driver/mysql"
	redis "github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// DataSource 统一的数据源，包含所有数据库组件
type DataSource struct {
	Mysql   *gorm.DB
	Mongodb *mongo.Database
	Redis   *redis.Client
	Kafka   sarama.SyncProducer
}

// ComponentStatus 组件状态
type ComponentStatus struct {
	Enabled bool
	Error   error
}

// NewDataSource 根据配置初始化数据源
func NewDataSource(lc fx.Lifecycle, cfg configs.Config) *DataSource {
	ds := &DataSource{}

	// 初始化MySQL
	if cfg.Mysql.Host != "" {
		ds.Mysql = NewMysql(cfg.Mysql)
	}

	// 初始化MongoDB
	if cfg.Mongodb.Host != "" {
		ds.Mongodb = MongoClient(cfg.Mongodb)
	}

	// 初始化Redis
	if cfg.Redis.Host != "" {
		ds.Redis = NewRedis(cfg.Redis)
	}

	// 初始化Kafka
	if len(cfg.Kafka.Brokers) > 0 {
		producer, err := NewKafkaProducer(cfg.Kafka)
		if err != nil {
			panic(err)
		}
		ds.Kafka = producer
	}

	// 注册生命周期钩子
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return ds.start(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return ds.Stop(ctx)
		},
	})

	return ds
}

// start 启动所有组件
func (ds *DataSource) start(ctx context.Context) error {
	// 检查MySQL连接
	if ds.Mysql != nil {
		if sqlDB, err := ds.Mysql.DB(); err == nil {
			if err := sqlDB.PingContext(ctx); err != nil {
				return err
			}
		}
	}

	// 检查Redis连接
	if ds.Redis != nil {
		if err := ds.Redis.Ping(ctx).Err(); err != nil {
			return err
		}
	}

	// Kafka连接检查（可选）
	if ds.Kafka != nil {
		// 发送一条空消息作为连通性检查（可选）
		// 忽略错误以避免启动硬失败，也可改为严格校验
	}

	return nil
}

// Stop 停止所有组件
func (ds *DataSource) Stop(ctx context.Context) error {
	// 关闭MySQL连接
	if ds.Mysql != nil {
		if sqlDB, err := ds.Mysql.DB(); err == nil {
			_ = sqlDB.Close()
		}
	}

	// 关闭Redis连接
	if ds.Redis != nil {
		_ = ds.Redis.Close()
	}

	// 关闭Kafka连接
	if ds.Kafka != nil {
		_ = ds.Kafka.Close()
	}

	// MongoDB连接由客户端管理
	return nil
}

// GetComponentStatus 获取所有组件状态
func (ds *DataSource) GetComponentStatus(ctx context.Context) map[string]ComponentStatus {
	status := make(map[string]ComponentStatus)

	// MySQL状态
	if ds.Mysql != nil {
		if sqlDB, err := ds.Mysql.DB(); err == nil {
			status["mysql"] = ComponentStatus{
				Enabled: true,
				Error:   sqlDB.PingContext(ctx),
			}
		} else {
			status["mysql"] = ComponentStatus{
				Enabled: true,
				Error:   err,
			}
		}
	} else {
		status["mysql"] = ComponentStatus{Enabled: false}
	}

	// Redis状态
	if ds.Redis != nil {
		status["redis"] = ComponentStatus{
			Enabled: true,
			Error:   ds.Redis.Ping(ctx).Err(),
		}
	} else {
		status["redis"] = ComponentStatus{Enabled: false}
	}

	// Kafka状态
	if ds.Kafka != nil {
		status["kafka"] = ComponentStatus{
			Enabled: true,
			Error:   nil, // Kafka状态检查比较复杂，这里简化处理
		}
	} else {
		status["kafka"] = ComponentStatus{Enabled: false}
	}

	// MongoDB状态
	if ds.Mongodb != nil {
		status["mongodb"] = ComponentStatus{
			Enabled: true,
			Error:   nil, // MongoDB状态检查比较复杂，这里简化处理
		}
	} else {
		status["mongodb"] = ComponentStatus{Enabled: false}
	}

	return status
}`,

	"internal/api/data/db/mysql.go.tmpl": `package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"database/sql"

	"{{.ModulePath}}/internal/configs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewMysql(cfg configs.MysqlConfig) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,          // Don't include params in the SQL log
			Colorful:                  false,         // Disable color
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		panic(err)
	}

	err = db.Callback().Create().After("gorm:create").Register("role:menu_after_create", AfterCreate)
	if err != nil {
		panic(err)
	}
	// set connection pool on the underlying *sql.DB
	sqlDB, err := db.DB()
	if err == nil {
		configureSQLPool(sqlDB, cfg)
	}
	return db
}

func configureSQLPool(sqlDB *sql.DB, cfg configs.MysqlConfig) {
	if cfg.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	// 5 minutes default if not set elsewhere; safe choice
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
}

func AfterCreate(db *gorm.DB) {
	// 在这里添加创建后的回调逻辑
}`,

	"internal/api/data/db/redis.go.tmpl": `package db

import (
	"context"

	"{{.ModulePath}}/internal/configs"
	redis "github.com/redis/go-redis/v9"
)

func NewRedis(cfg configs.RedisConfig) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Password: cfg.Password, // no password set
		DB:       cfg.Database, // use default DB
	})
	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}
	return rdb
}`,

	"internal/api/data/db/mongodb.go.tmpl": `package db

import (
	"context"

	"{{.ModulePath}}/internal/configs"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoClient(cfg configs.Mongodb) *mongo.Database {

	uri := "mongodb://"
	if cfg.Password != "" && cfg.User != "" {
		uri += cfg.User + ":" + cfg.Password + "@"
	}
	uri += cfg.Host + ":" + cfg.Port

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	if err := client.Connect(context.Background()); err != nil {
		panic(err)
	}
	return client.Database(cfg.DataBase)
}`,

	"internal/api/data/db/kafka.go.tmpl": `package db

import (
	"github.com/IBM/sarama"
	"{{.ModulePath}}/internal/configs"
)

func NewKafkaProducer(cfg configs.KafkaConfig) (sarama.SyncProducer, error) {
	sc := sarama.NewConfig()
	sc.ClientID = cfg.ClientID
	sc.Producer.RequiredAcks = sarama.WaitForAll
	sc.Producer.Retry.Max = 3
	sc.Producer.Return.Successes = true
	return sarama.NewSyncProducer(cfg.Brokers, sc)
}`,

	"internal/api/data/db/callback.go.tmpl": `package db

import (
	"gorm.io/gorm"
)

// Callback 数据库回调
type Callback struct{}

// NewCallback 创建回调
func NewCallback() *Callback {
	return &Callback{}
}

// BeforeCreate 创建前回调
func (c *Callback) BeforeCreate(tx *gorm.DB) error {
	// 在这里添加创建前的逻辑
	return nil
}

// AfterCreate 创建后回调
func (c *Callback) AfterCreate(tx *gorm.DB) error {
	// 在这里添加创建后的逻辑
	return nil
}

// BeforeUpdate 更新前回调
func (c *Callback) BeforeUpdate(tx *gorm.DB) error {
	// 在这里添加更新前的逻辑
	return nil
}

// AfterUpdate 更新后回调
func (c *Callback) AfterUpdate(tx *gorm.DB) error {
	// 在这里添加更新后的逻辑
	return nil
}`,

	"internal/api/data/model/user.go.tmpl": `package model

import (
	"time"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           ` + "`" + `gorm:"primarykey" json:"id"` + "`" + `
	Username  string         ` + "`" + `gorm:"uniqueIndex;size:50;not null" json:"username"` + "`" + `
	Email     string         ` + "`" + `gorm:"uniqueIndex;size:100;not null" json:"email"` + "`" + `
	Password  string         ` + "`" + `gorm:"size:255;not null" json:"-"` + "`" + `
	Role      string         ` + "`" + `gorm:"size:20;default:user" json:"role"` + "`" + `
	Status    int            ` + "`" + `gorm:"default:1" json:"status"` + "`" + `
	CreatedAt time.Time      ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt time.Time      ` + "`" + `json:"updated_at"` + "`" + `
	DeletedAt gorm.DeletedAt ` + "`" + `gorm:"index" json:"-"` + "`" + `
}

// TableName 表名
func (User) TableName() string {
	return "users"
}`,

	"internal/api/data/query/gen.go.tmpl": `package query

import (
	"gorm.io/gen"
	"gorm.io/gorm"
)

// Generator 代码生成器
type Generator struct {
	db *gorm.DB
}

// NewGenerator 创建生成器
func NewGenerator(db *gorm.DB) *Generator {
	return &Generator{db: db}
}

// Generate 生成代码
func (g *Generator) Generate() error {
	gen := gen.NewGenerator(gen.Config{
		OutPath: "./internal/api/data/query",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	gen.UseDB(g.db)

	// 生成所有表的模型
	gen.ApplyBasic(gen.GenerateAllTable()...)

	gen.Execute()

	return nil
}`,

	"internal/utils/context.go.tmpl": `package utils

import (
	"context"
	"time"
)

// ContextKey 上下文键类型
type ContextKey string

const (
	// UserIDKey 用户ID键
	UserIDKey ContextKey = "user_id"
	// UsernameKey 用户名键
	UsernameKey ContextKey = "username"
	// RoleKey 角色键
	RoleKey ContextKey = "role"
)

// WithUserID 设置用户ID到上下文
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// GetUserID 从上下文获取用户ID
func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}

// WithUsername 设置用户名到上下文
func WithUsername(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, UsernameKey, username)
}

// GetUsername 从上下文获取用户名
func GetUsername(ctx context.Context) (string, bool) {
	username, ok := ctx.Value(UsernameKey).(string)
	return username, ok
}

// WithRole 设置角色到上下文
func WithRole(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, RoleKey, role)
}

// GetRole 从上下文获取角色
func GetRole(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(RoleKey).(string)
	return role, ok
}

// WithTimeout 设置超时上下文
func WithTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, timeout)
}`,

	"internal/utils/encrypt.go.tmpl": `package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// MD5 MD5加密
func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// SHA256 SHA256加密
func SHA256(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// EncryptPassword 加密密码
func EncryptPassword(password, salt string) string {
	return SHA256(password + salt)
}

// VerifyPassword 验证密码
func VerifyPassword(password, salt, hashedPassword string) bool {
	return EncryptPassword(password, salt) == hashedPassword
}

// GenerateSalt 生成盐值
func GenerateSalt() string {
	return fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%d", time.Now().UnixNano()))))
}`,

	"internal/utils/validator.go.tmpl": `package utils

import (
	"regexp"
	"strings"
)

// IsValidEmail 验证邮箱格式
func IsValidEmail(email string) bool {
	pattern := ` + "`" + `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$` + "`" + `
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

// IsValidPhone 验证手机号格式
func IsValidPhone(phone string) bool {
	pattern := ` + "`" + `^1[3-9]\d{9}$` + "`" + `
	matched, _ := regexp.MatchString(pattern, phone)
	return matched
}

// IsValidUsername 验证用户名格式
func IsValidUsername(username string) bool {
	if len(username) < 3 || len(username) > 20 {
		return false
	}
	pattern := ` + "`" + `^[a-zA-Z0-9_]+$` + "`" + `
	matched, _ := regexp.MatchString(pattern, username)
	return matched
}

// IsValidPassword 验证密码强度
func IsValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	// 至少包含字母和数字
	hasLetter := regexp.MustCompile(` + "`" + `[a-zA-Z]` + "`" + `).MatchString(password)
	hasNumber := regexp.MustCompile(` + "`" + `[0-9]` + "`" + `).MatchString(password)
	return hasLetter && hasNumber
}

// TrimSpace 去除字符串两端空格
func TrimSpace(s string) string {
	return strings.TrimSpace(s)
}`,

	"internal/health/checker.go.tmpl": `package health

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Checker 健康检查器
type Checker struct {
	checks map[string]CheckFunc
}

// CheckFunc 检查函数类型
type CheckFunc func(ctx context.Context) error

// NewChecker 创建健康检查器
func NewChecker() *Checker {
	return &Checker{
		checks: make(map[string]CheckFunc),
	}
}

// AddCheck 添加检查项
func (c *Checker) AddCheck(name string, check CheckFunc) {
	c.checks[name] = check
}

// Check 执行所有检查
func (c *Checker) Check(ctx context.Context) map[string]error {
	results := make(map[string]error)
	
	for name, check := range c.checks {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		err := check(ctx)
		cancel()
		
		results[name] = err
	}
	
	return results
}

// IsHealthy 检查是否健康
func (c *Checker) IsHealthy(ctx context.Context) bool {
	results := c.Check(ctx)
	for _, err := range results {
		if err != nil {
			return false
		}
	}
	return true
}

// HTTPHandler HTTP健康检查处理器
func (c *Checker) HTTPHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	if c.IsHealthy(ctx) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("Service Unavailable"))
	}
}`,

	"internal/metrics/prometheus.go.tmpl": `package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics 指标收集器
type Metrics struct {
	HTTPRequestsTotal    *prometheus.CounterVec
	HTTPRequestDuration  *prometheus.HistogramVec
	HTTPRequestsInFlight prometheus.Gauge
}

// NewMetrics 创建指标收集器
func NewMetrics() *Metrics {
	return &Metrics{
		HTTPRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "endpoint", "status"},
		),
		HTTPRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "endpoint"},
		),
		HTTPRequestsInFlight: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "http_requests_in_flight",
				Help: "Current number of HTTP requests being processed",
			},
		),
	}
}

// RecordRequest 记录请求
func (m *Metrics) RecordRequest(method, endpoint, status string, duration float64) {
	m.HTTPRequestsTotal.WithLabelValues(method, endpoint, status).Inc()
	m.HTTPRequestDuration.WithLabelValues(method, endpoint).Observe(duration)
}

// IncRequestsInFlight 增加进行中的请求数
func (m *Metrics) IncRequestsInFlight() {
	m.HTTPRequestsInFlight.Inc()
}

// DecRequestsInFlight 减少进行中的请求数
func (m *Metrics) DecRequestsInFlight() {
	m.HTTPRequestsInFlight.Dec()
}`,

	"internal/middleware/rate_limit.go.tmpl": `package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RateLimiter 限流器
type RateLimiter struct {
	client *redis.Client
}

// NewRateLimiter 创建限流器
func NewRateLimiter(client *redis.Client) *RateLimiter {
	return &RateLimiter{
		client: client,
	}
}

// Allow 检查是否允许请求
func (r *RateLimiter) Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	pipe := r.client.Pipeline()
	
	// 使用滑动窗口算法
	now := time.Now()
	windowStart := now.Add(-window)
	
	// 删除过期的记录
	pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart.Unix()))
	
	// 添加当前请求
	pipe.ZAdd(ctx, key, redis.Z{
		Score:  float64(now.Unix()),
		Member: now.UnixNano(),
	})
	
	// 获取当前窗口内的请求数
	pipe.ZCard(ctx, key)
	
	// 设置过期时间
	pipe.Expire(ctx, key, window)
	
	results, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}
	
	count := results[2].(*redis.IntCmd).Val()
	return count <= int64(limit), nil
}`,

	"internal/validator/custom.go.tmpl": `package validator

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

// CustomValidator 自定义验证器
type CustomValidator struct {
	validator *validator.Validate
}

// NewCustomValidator 创建自定义验证器
func NewCustomValidator() *CustomValidator {
	v := validator.New()
	
	// 注册自定义验证规则
	v.RegisterValidation("email", validateEmail)
	v.RegisterValidation("phone", validatePhone)
	v.RegisterValidation("username", validateUsername)
	v.RegisterValidation("password", validatePassword)
	
	return &CustomValidator{
		validator: v,
	}
}

// Validate 验证结构体
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// validateEmail 验证邮箱
func validateEmail(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	pattern := ` + "`" + `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$` + "`" + `
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

// validatePhone 验证手机号
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	pattern := ` + "`" + `^1[3-9]\d{9}$` + "`" + `
	matched, _ := regexp.MatchString(pattern, phone)
	return matched
}

// validateUsername 验证用户名
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	if len(username) < 3 || len(username) > 20 {
		return false
	}
	pattern := ` + "`" + `^[a-zA-Z0-9_]+$` + "`" + `
	matched, _ := regexp.MatchString(pattern, username)
	return matched
}

// validatePassword 验证密码
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 8 {
		return false
	}
	// 至少包含字母和数字
	hasLetter := regexp.MustCompile(` + "`" + `[a-zA-Z]` + "`" + `).MatchString(password)
	hasNumber := regexp.MustCompile(` + "`" + `[0-9]` + "`" + `).MatchString(password)
	return hasLetter && hasNumber
}`,

	"internal/docs/swagger.go.tmpl": `package docs

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// SwaggerConfig Swagger配置
type SwaggerConfig struct {
	Title       string
	Description string
	Version     string
	Host        string
	BasePath    string
}

// NewSwaggerConfig 创建Swagger配置
func NewSwaggerConfig() *SwaggerConfig {
	return &SwaggerConfig{
		Title:       "{{.ProjectName}} API",
		Description: "{{.ProjectName}} API Documentation",
		Version:     "1.0.0",
		Host:        "localhost:8080",
		BasePath:    "/api",
	}
}

// SetupSwagger 设置Swagger
func SetupSwagger(e *echo.Echo, config *SwaggerConfig) {
	// 这里可以集成swagger-ui
	// 简化示例，实际应该使用swagger相关的库
	e.GET("/swagger/*", func(c echo.Context) error {
		return c.String(200, "Swagger documentation")
	})
}`,

	"internal/cache/redis.go.tmpl": `package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

// Cache Redis缓存
type Cache struct {
	client *redis.Client
}

// NewCache 创建缓存
func NewCache(client *redis.Client) *Cache {
	return &Cache{
		client: client,
	}
}

// Set 设置缓存
func (c *Cache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, expiration).Err()
}

// Get 获取缓存
func (c *Cache) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), dest)
}

// Delete 删除缓存
func (c *Cache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// Exists 检查缓存是否存在
func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	count, err := c.client.Exists(ctx, key).Result()
	return count > 0, err
}`,
}

// TemplateData 模板数据
type TemplateData struct {
	ModulePath  string
	ProjectName string
	PackageName string
}

// TemplateGenerator 模板生成器
type TemplateGenerator struct {
	outputDir string
	data      TemplateData
}

// NewTemplateGenerator 创建模板生成器
func NewTemplateGenerator(outputDir string, modulePath, projectName string) *TemplateGenerator {
	return &TemplateGenerator{
		outputDir: outputDir,
		data: TemplateData{
			ModulePath:  modulePath,
			ProjectName: projectName,
			PackageName: filepath.Base(modulePath),
		},
	}
}

// Generate 生成项目文件
func (g *TemplateGenerator) Generate() error {
	// 创建输出目录
	if err := os.MkdirAll(g.outputDir, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	// 定义要生成的文件
	files := map[string]string{
		"go.mod":                         "go.mod.tmpl",
		"main.go":                        "main.go.tmpl",
		"README.md":                      "README.md.tmpl",
		"Makefile":                       "Makefile.tmpl",
		"cmd/fx.go":                      "cmd/fx.go.tmpl",
		"cmd/root.go":                    "cmd/root.go.tmpl",
		"cmd/gen.go":                     "cmd/gen.go.tmpl",
		"internal/server/echo_server.go": "internal/server/echo_server.go.tmpl",
		"internal/server/config.go":      "internal/server/config.go.tmpl",
		"internal/server/middlewares/middleware.go": "internal/server/middlewares/middleware.go.tmpl",
		"internal/server/middlewares/config.go":     "internal/server/middlewares/config.go.tmpl",
		"internal/server/middlewares/error.go":      "internal/server/middlewares/error.go.tmpl",
		"internal/server/middlewares/jwt.go":        "internal/server/middlewares/jwt.go.tmpl",
		"internal/server/middlewares/casbin.go":     "internal/server/middlewares/casbin.go.tmpl",
		"internal/configs/config.go":                "internal/configs/config.go.tmpl",
		"internal/configs/bootstrap.go":             "internal/configs/bootstrap.go.tmpl",
		"internal/configs/store.go":                 "internal/configs/store.go.tmpl",
		"internal/configs/file_source.go":           "internal/configs/file_source.go.tmpl",
		"internal/configs/merge.go":                 "internal/configs/merge.go.tmpl",
		"internal/configs/hot_reload.go":            "internal/configs/hot_reload.go.tmpl",
		"internal/configs/etcd_source.go":           "internal/configs/etcd_source.go.tmpl",
		"internal/configs/consul_source.go":         "internal/configs/consul_source.go.tmpl",
		"internal/log/logger.go":                    "internal/log/logger.go.tmpl",
		"internal/log/factory.go":                   "internal/log/factory.go.tmpl",
		"internal/log/console_sink.go":              "internal/log/console_sink.go.tmpl",
		"internal/log/file_sink.go":                 "internal/log/file_sink.go.tmpl",
		"internal/log/global.go":                    "internal/log/global.go.tmpl",
		"internal/log/slog.go":                      "internal/log/slog.go.tmpl",
		"internal/resp/response.go":                 "internal/resp/response.go.tmpl",
		"internal/code/code.go":                     "internal/code/code.go.tmpl",
		"internal/code/base.go":                     "internal/code/base.go.tmpl",
		"internal/code/errors.go":                   "internal/code/errors.go.tmpl",
		"internal/code/http_status.go":              "internal/code/http_status.go.tmpl",
		"internal/code/error_types.go":              "internal/code/error_types.go.tmpl",
		"internal/api/biz/biz.go":                   "internal/api/biz/biz.go.tmpl",
		"internal/api/service/service.go":           "internal/api/service/service.go.tmpl",
		"internal/api/data/data.go":                 "internal/api/data/data.go.tmpl",
		"internal/api/data/casbin.go":               "internal/api/data/casbin.go.tmpl",
		"internal/api/data/jwt.go":                  "internal/api/data/jwt.go.tmpl",
		"internal/api/data/db/db.go":                "internal/api/data/db/db.go.tmpl",
		"internal/api/data/db/mysql.go":             "internal/api/data/db/mysql.go.tmpl",
		"internal/api/data/db/redis.go":             "internal/api/data/db/redis.go.tmpl",
		"internal/api/data/db/mongodb.go":           "internal/api/data/db/mongodb.go.tmpl",
		"internal/api/data/db/kafka.go":             "internal/api/data/db/kafka.go.tmpl",
		"internal/api/data/db/callback.go":          "internal/api/data/db/callback.go.tmpl",
		"internal/api/data/model/user.go":           "internal/api/data/model/user.go.tmpl",
		"internal/api/data/query/gen.go":            "internal/api/data/query/gen.go.tmpl",
		"internal/utils/context.go":                 "internal/utils/context.go.tmpl",
		"internal/utils/encrypt.go":                 "internal/utils/encrypt.go.tmpl",
		"internal/utils/validator.go":               "internal/utils/validator.go.tmpl",
		"internal/health/checker.go":                "internal/health/checker.go.tmpl",
		"internal/metrics/prometheus.go":            "internal/metrics/prometheus.go.tmpl",
		"internal/middleware/rate_limit.go":         "internal/middleware/rate_limit.go.tmpl",
		"internal/validator/custom.go":              "internal/validator/custom.go.tmpl",
		"internal/docs/swagger.go":                  "internal/docs/swagger.go.tmpl",
		"internal/cache/redis.go":                   "internal/cache/redis.go.tmpl",
		"configs/config.toml":                       "configs/config.toml.tmpl",
		"configs/rbac_model.conf":                   "configs/rbac_model.conf.tmpl",
		"env.example":                               "env.example.tmpl",
	}

	// 先创建目录结构
	if err := g.createDirectories(); err != nil {
		return fmt.Errorf("创建目录结构失败: %w", err)
	}

	// 生成每个文件
	for outputFile, templateFile := range files {
		if err := g.generateFile(outputFile, templateFile); err != nil {
			return fmt.Errorf("生成文件 %s 失败: %w", outputFile, err)
		}
	}

	return nil
}

// generateFile 生成单个文件
func (g *TemplateGenerator) generateFile(outputFile, templateFile string) error {
	// 从嵌入的模板中获取内容
	templateContent, exists := templates[templateFile]
	if !exists {
		return fmt.Errorf("模板文件 %s 不存在", templateFile)
	}

	tmpl, err := template.New(templateFile).Parse(templateContent)
	if err != nil {
		return fmt.Errorf("解析模板文件失败: %w", err)
	}

	// 创建输出文件
	outputPath := filepath.Join(g.outputDir, outputFile)
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("创建输出文件失败: %w", err)
	}
	defer file.Close()

	// 执行模板
	if err := tmpl.Execute(file, g.data); err != nil {
		return fmt.Errorf("执行模板失败: %w", err)
	}

	fmt.Printf("✅ 生成文件: %s\n", outputPath)
	return nil
}

// createDirectories 创建目录结构
func (g *TemplateGenerator) createDirectories() error {
	dirs := []string{
		"cmd",
		"configs",
		"internal/api/biz",
		"internal/api/data",
		"internal/api/data/db",
		"internal/api/data/model",
		"internal/api/data/query",
		"internal/api/service",
		"internal/cache",
		"internal/code",
		"internal/configs",
		"internal/health",
		"internal/log",
		"internal/metrics",
		"internal/middleware",
		"internal/resp",
		"internal/server",
		"internal/server/middlewares",
		"internal/utils",
		"internal/validator",
		"internal/docs",
		"k8s",
		"scripts",
		"sql",
	}

	for _, dir := range dirs {
		dirPath := filepath.Join(g.outputDir, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("创建目录 %s 失败: %w", dir, err)
		}
		fmt.Printf("✅ 创建目录: %s\n", dirPath)
	}

	return nil
}
