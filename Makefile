# =============================================================================
# Go Template Project Makefile
# 提供业务开发常用的构建、运行、测试和数据库模型生成命令
# =============================================================================

# 默认目标
.DEFAULT_GOAL := help

# 颜色定义
BLUE := \033[0;34m
GREEN := \033[0;32m
YELLOW := \033[1;33m
RED := \033[0;31m
NC := \033[0m # No Color

# 项目配置
BIN_DIR := bin
APP_NAME := app
DEFAULT_DSN := "root:12345678@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"

# =============================================================================
# 基础命令
# =============================================================================

.PHONY: build run tidy push

# 构建应用
build:
	@echo "$(BLUE)[INFO]$(NC) Building application..."
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/$(APP_NAME) main.go
	@echo "$(GREEN)[SUCCESS]$(NC) Build completed: $(BIN_DIR)/$(APP_NAME)"

# 运行应用
run:
	@echo "$(BLUE)[INFO]$(NC) Starting application..."
	@go run main.go --config configs/config.toml

# 运行应用（开发模式）
run-dev:
	@echo "$(BLUE)[INFO]$(NC) Starting application in development mode..."
	@export RUN_ENVIRONMENT=dev && go run main.go --config configs/config.toml

# 运行应用（测试模式）
run-test:
	@echo "$(BLUE)[INFO]$(NC) Starting application in test mode..."
	@export RUN_ENVIRONMENT=test && go run main.go --config configs/config.toml

# 整理依赖
tidy:
	@echo "$(BLUE)[INFO]$(NC) Tidying dependencies..."
	@go mod tidy
	@echo "$(GREEN)[SUCCESS]$(NC) Dependencies tidied"

# 提交代码
push:
	@if [ -z "$(m)" ]; then \
		echo "$(RED)[ERROR]$(NC) Usage: make push m=\"commit_message\""; \
		exit 1; \
	fi
	@echo "$(BLUE)[INFO]$(NC) Preparing to push..."
	@go mod download && go mod vendor && git add . && git commit -m '$(m)' && git push
	@echo "$(GREEN)[SUCCESS]$(NC) Code committed with message: $(m)"

# =============================================================================
# 代码质量工具
# =============================================================================

.PHONY: fmt vet lint test test-verbose test-coverage clean

# 格式化代码
fmt:
	@echo "$(BLUE)[INFO]$(NC) Formatting code..."
	@gofmt -s -w .
	@command -v goimports >/dev/null 2>&1 && goimports -w . || echo "goimports 未安装，跳过 (安装: go install golang.org/x/tools/cmd/goimports@latest)"
	@echo "$(GREEN)[SUCCESS]$(NC) Code formatting completed"

# 代码检查
vet:
	@echo "$(BLUE)[INFO]$(NC) Running go vet..."
	@go vet ./...
	@echo "$(GREEN)[SUCCESS]$(NC) Code vetting completed"

# 代码检查（使用golangci-lint）
lint:
	@echo "$(BLUE)[INFO]$(NC) Running linter..."
	@golangci-lint run --skip-dirs=vendor,internal/api/data/query --skip-files='.*\.gen\.go$$' || true
	@echo "$(GREEN)[SUCCESS]$(NC) Linting completed"

# 严格代码检查（失败时退出）
lint-strict:
	@echo "$(BLUE)[INFO]$(NC) Running strict linter..."
	@golangci-lint run --skip-dirs=vendor,internal/api/data/query --skip-files='.*\.gen\.go$$'
	@echo "$(GREEN)[SUCCESS]$(NC) Strict linting completed"

# 快速代码检查（只运行快速linter）
lint-fast:
	@echo "$(BLUE)[INFO]$(NC) Running fast linter..."
	@golangci-lint run --fast-only --skip-dirs=vendor,internal/api/data/query --skip-files='.*\.gen\.go$$'
	@echo "$(GREEN)[SUCCESS]$(NC) Fast linting completed"

# 修复可自动修复的问题
lint-fix:
	@echo "$(BLUE)[INFO]$(NC) Running linter with auto-fix..."
	@golangci-lint run --fix --skip-dirs=vendor,internal/api/data/query --skip-files='.*\.gen\.go$$'
	@echo "$(GREEN)[SUCCESS]$(NC) Linting with auto-fix completed"

# 检查特定目录
lint-dir:
	@if [ -z "$(DIR)" ]; then \
		echo "$(RED)[ERROR]$(NC) Usage: make lint-dir DIR=./internal/api"; \
		exit 1; \
	fi
	@echo "$(BLUE)[INFO]$(NC) Running linter on directory: $(DIR)"
	@golangci-lint run --skip-dirs=vendor,internal/api/data/query --skip-files='.*\.gen\.go$$' $(DIR)
	@echo "$(GREEN)[SUCCESS]$(NC) Directory linting completed"

# 生成lint报告
lint-report:
	@echo "$(BLUE)[INFO]$(NC) Generating lint report..."
	@golangci-lint run --output.checkstyle.path=golangci-report.xml --skip-dirs=vendor,internal/api/data/query --skip-files='.*\.gen\.go$$' || true
	@echo "$(GREEN)[SUCCESS]$(NC) Lint report generated: golangci-report.xml"

# 安装golangci-lint
install-lint:
	@echo "$(BLUE)[INFO]$(NC) Installing golangci-lint..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v2.4.0
	@echo "$(GREEN)[SUCCESS]$(NC) golangci-lint installed"

# 运行测试
test:
	@echo "$(BLUE)[INFO]$(NC) Running tests..."
	@export RUN_ENVIRONMENT=test && go test -race $(shell go list ./...)
	@echo "$(GREEN)[SUCCESS]$(NC) Tests completed"

# 详细测试输出
test-verbose:
	@echo "$(BLUE)[INFO]$(NC) Running tests with verbose output..."
	@export RUN_ENVIRONMENT=test && go test -v -race $(shell go list ./...)
	@echo "$(GREEN)[SUCCESS]$(NC) Verbose tests completed"

# 生成测试覆盖率报告
test-coverage:
	@echo "$(BLUE)[INFO]$(NC) Generating test coverage report..."
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)[SUCCESS]$(NC) Coverage report generated: coverage.html"

# =============================================================================
# 代码生成工具
# =============================================================================

.PHONY: ensure-gentool db-gen db-gen-table gen-all

ensure-gentool:
	@command -v gentool >/dev/null 2>&1 || { \
		echo "$(YELLOW)[WARNING]$(NC) gentool 未安装，正在安装 gorm.io/gen/tools/gentool@v0.3.28"; \
		go install gorm.io/gen/tools/gentool@v0.3.28; \
	}

# 生成数据库模型和查询
db-gen: ensure-gentool
	@echo "$(BLUE)[INFO]$(NC) Generating database models and queries..."
	@$(shell go env GOPATH)/bin/gentool \
		-dsn=$(DEFAULT_DSN) \
		-outPath="./internal/api/data/query" \
		-modelPkgName="model" \
		-fieldWithIndexTag \
		-fieldWithTypeTag
	@echo "$(GREEN)[SUCCESS]$(NC) Database generation completed"

# 生成指定表的模型和查询方法
db-gen-table: ensure-gentool
	@if [ -z "$(TABLE)" ]; then \
		echo "$(RED)[ERROR]$(NC) Usage: make db-gen-table TABLE=table_name"; \
		exit 1; \
	fi
	@echo "$(BLUE)[INFO]$(NC) Generating model for table: $(TABLE)"
	@$(shell go env GOPATH)/bin/gentool \
		-dsn=$(DEFAULT_DSN) \
		-outPath="./internal/api/data/query" \
		-modelPkgName="model" \
		-fieldWithIndexTag \
		-fieldWithTypeTag \
		-tables="$(TABLE)"
	@echo "$(GREEN)[SUCCESS]$(NC) Table $(TABLE) generation completed"

# 完整生成（数据库模型和查询）
gen-all: db-gen
	@echo "$(GREEN)[SUCCESS]$(NC) All generation completed"

# =============================================================================
# 开发工作流
# =============================================================================

.PHONY: dev-setup dev-check dev-full

# 开发环境设置
dev-setup: tidy
	@echo "$(BLUE)[INFO]$(NC) Setting up development environment..."
	@go mod download
	@echo "$(GREEN)[SUCCESS]$(NC) Development environment ready"

# 开发检查（格式化、检查、测试）
dev-check: fmt vet lint test
	@echo "$(GREEN)[SUCCESS]$(NC) Development check completed"

# 完整开发流程
dev-full: clean dev-check
	@echo "$(GREEN)[SUCCESS]$(NC) Full development workflow completed"

# =============================================================================
# 清理和维护
# =============================================================================

# 清理生成的文件
clean:
	@echo "$(BLUE)[INFO]$(NC) Cleaning generated files..."
	@rm -f coverage.out coverage.html
	@rm -rf $(BIN_DIR)
	@echo "$(GREEN)[SUCCESS]$(NC) Clean completed"

# 深度清理
clean-all: clean
	@echo "$(BLUE)[INFO]$(NC) Deep cleaning..."
	@go clean -cache
	@go clean -modcache
	@echo "$(GREEN)[SUCCESS]$(NC) Deep clean completed"

# =============================================================================
# Docker 相关命令
# =============================================================================

.PHONY: docker-build docker-run docker-stop docker-clean

# 构建Docker镜像
docker-build:
	@echo "$(BLUE)[INFO]$(NC) Building Docker image..."
	@docker build -t go-template:latest .
	@echo "$(GREEN)[SUCCESS]$(NC) Docker image built: go-template:latest"

# 运行Docker容器
docker-run:
	@echo "$(BLUE)[INFO]$(NC) Starting Docker container..."
	@docker-compose up -d
	@echo "$(GREEN)[SUCCESS]$(NC) Docker container started"

# 停止Docker容器
docker-stop:
	@echo "$(BLUE)[INFO]$(NC) Stopping Docker container..."
	@docker-compose down
	@echo "$(GREEN)[SUCCESS]$(NC) Docker container stopped"

# 清理Docker资源
docker-clean: docker-stop
	@echo "$(BLUE)[INFO]$(NC) Cleaning Docker resources..."
	@docker system prune -f
	@echo "$(GREEN)[SUCCESS]$(NC) Docker resources cleaned"

.PHONY: security-scan

# 安全扫描
security-scan:
	@echo "$(BLUE)[INFO]$(NC) Running security scan..."
	@gosec ./...
	@echo "$(GREEN)[SUCCESS]$(NC) Security scan completed"

# =============================================================================
# 性能测试相关命令
# =============================================================================

.PHONY: bench-load-test

# 性能基准测试
bench:
	@echo "$(BLUE)[INFO]$(NC) Running performance benchmarks..."
	@go test -bench=. -benchmem ./...
	@echo "$(GREEN)[SUCCESS]$(NC) Performance benchmarks completed"

# 负载测试
load-test:
	@echo "$(BLUE)[INFO]$(NC) Running load tests..."
	@echo "Please install hey: go install github.com/rakyll/hey@latest"
	@hey -n 1000 -c 10 http://localhost:9322/api/users
	@echo "$(GREEN)[SUCCESS]$(NC) Load tests completed"

# =============================================================================
# 帮助信息
# =============================================================================

.PHONY: help

help:
	@echo "$(BLUE)Go Template Project - Available Commands$(NC)"
	@echo ""
	@echo "$(YELLOW)基础命令:$(NC)"
	@echo "  $(GREEN)build$(NC)              - 构建应用程序"
	@echo "  $(GREEN)run$(NC)                - 运行应用程序"
	@echo "  $(GREEN)run-dev$(NC)            - 运行应用程序（开发模式）"
	@echo "  $(GREEN)run-test$(NC)           - 运行应用程序（测试模式）"
	@echo "  $(GREEN)tidy$(NC)               - 整理Go模块依赖"
	@echo "  $(GREEN)push$(NC)               - 提交代码 (需要设置 m=commit_message)"
	@echo ""
	@echo "$(YELLOW)代码质量:$(NC)"
	@echo "  $(GREEN)fmt$(NC)                - 格式化代码"
	@echo "  $(GREEN)vet$(NC)                - 运行go vet检查"
	@echo "  $(GREEN)lint$(NC)               - 运行golangci-lint检查"
	@echo "  $(GREEN)lint-strict$(NC)        - 严格代码检查（失败时退出）"
	@echo "  $(GREEN)lint-fast$(NC)          - 快速代码检查"
	@echo "  $(GREEN)lint-fix$(NC)           - 自动修复可修复的问题"
	@echo "  $(GREEN)lint-dir$(NC)           - 检查特定目录 (DIR=./path)"
	@echo "  $(GREEN)lint-report$(NC)        - 生成lint报告"
	@echo "  $(GREEN)install-lint$(NC)       - 安装golangci-lint"
	@echo "  $(GREEN)test$(NC)               - 运行所有测试"
	@echo "  $(GREEN)test-verbose$(NC)       - 运行详细测试"
	@echo "  $(GREEN)test-coverage$(NC)      - 生成测试覆盖率报告"
	@echo ""
	@echo "$(YELLOW)代码生成:$(NC)"
	@echo "  $(GREEN)db-gen$(NC)                      - 生成数据库模型和查询"
	@echo "  $(GREEN)db-gen-table$(NC)                - 生成指定表模型 (TABLE=table_name)"
	@echo "  $(GREEN)gen-all$(NC)                     - 生成所有数据库代码"
	@echo ""
	@echo "$(YELLOW)开发工作流:$(NC)"
	@echo "  $(GREEN)dev-setup$(NC)          - 设置开发环境"
	@echo "  $(GREEN)dev-check$(NC)          - 运行开发检查"
	@echo "  $(GREEN)dev-full$(NC)           - 完整开发流程"
	@echo ""
	@echo "$(YELLOW)维护工具:$(NC)"
	@echo "  $(GREEN)clean$(NC)              - 清理生成的文件"
	@echo "  $(GREEN)clean-all$(NC)          - 深度清理"
	@echo "  $(GREEN)help$(NC)               - 显示此帮助信息"
	@echo ""
	@echo "$(YELLOW)Docker 命令:$(NC)"
	@echo "  $(GREEN)docker-build$(NC)       - 构建Docker镜像"
	@echo "  $(GREEN)docker-run$(NC)         - 运行Docker容器"
	@echo "  $(GREEN)docker-stop$(NC)        - 停止Docker容器"
	@echo "  $(GREEN)docker-clean$(NC)       - 清理Docker资源"
	@echo ""
	@echo "$(YELLOW)环境变量:$(NC)"
	@echo "  $(GREEN)TABLE$(NC)              - 表名 (用于db-gen-table)"
	@echo "  $(GREEN)DIR$(NC)                - 目录路径 (用于lint-dir)"
	@echo "  $(GREEN)m$(NC)                  - 提交消息 (用于push)"
	@echo ""
	@echo "$(YELLOW)示例用法:$(NC)"
	@echo "  make db-gen-table TABLE=users"
	@echo "  make lint-dir DIR=./internal/api"
	@echo "  make lint-fix"
	@echo "  make push m=\"feat: add user module\""
	@echo "  make dev-full"
