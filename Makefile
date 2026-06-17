# =============================================================================
# Go Template Project Makefile
# 提供业务开发常用的构建、运行、测试和检查命令
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
LOAD_TEST_URL ?= http://localhost:9322/api/health

# =============================================================================
# 基础命令
# =============================================================================

.PHONY: build run tidy

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

# 整理依赖
tidy:
	@echo "$(BLUE)[INFO]$(NC) Tidying dependencies..."
	@go mod tidy
	@echo "$(GREEN)[SUCCESS]$(NC) Dependencies tidied"

# =============================================================================
# 代码质量工具
# =============================================================================

.PHONY: fmt vet lint lint-strict lint-fast lint-fix lint-dir lint-report install-lint test test-verbose test-coverage clean clean-all

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
	@golangci-lint run
	@echo "$(GREEN)[SUCCESS]$(NC) Linting completed"

# 严格代码检查（失败时退出）
lint-strict:
	@echo "$(BLUE)[INFO]$(NC) Running strict linter..."
	@golangci-lint run
	@echo "$(GREEN)[SUCCESS]$(NC) Strict linting completed"

# 快速代码检查（只运行快速linter）
lint-fast:
	@echo "$(BLUE)[INFO]$(NC) Running fast linter..."
	@golangci-lint run --fast-only
	@echo "$(GREEN)[SUCCESS]$(NC) Fast linting completed"

# 修复可自动修复的问题
lint-fix:
	@echo "$(BLUE)[INFO]$(NC) Running linter with auto-fix..."
	@golangci-lint run --fix
	@echo "$(GREEN)[SUCCESS]$(NC) Linting with auto-fix completed"

# 检查特定目录
lint-dir:
	@if [ -z "$(DIR)" ]; then \
		echo "$(RED)[ERROR]$(NC) Usage: make lint-dir DIR=./internal/platform/server"; \
		exit 1; \
	fi
	@echo "$(BLUE)[INFO]$(NC) Running linter on directory: $(DIR)"
	@golangci-lint run $(DIR)
	@echo "$(GREEN)[SUCCESS]$(NC) Directory linting completed"

# 生成lint报告
lint-report:
	@echo "$(BLUE)[INFO]$(NC) Generating lint report..."
	@golangci-lint run --output.checkstyle.path=golangci-report.xml
	@echo "$(GREEN)[SUCCESS]$(NC) Lint report generated: golangci-report.xml"

# 安装golangci-lint
install-lint:
	@echo "$(BLUE)[INFO]$(NC) Installing golangci-lint..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v2.11.4
	@echo "$(GREEN)[SUCCESS]$(NC) golangci-lint installed"

# 运行测试
test:
	@echo "$(BLUE)[INFO]$(NC) Running tests..."
	@go test -race $(shell go list ./...)
	@echo "$(GREEN)[SUCCESS]$(NC) Tests completed"

# 详细测试输出
test-verbose:
	@echo "$(BLUE)[INFO]$(NC) Running tests with verbose output..."
	@go test -v -race $(shell go list ./...)
	@echo "$(GREEN)[SUCCESS]$(NC) Verbose tests completed"

# 生成测试覆盖率报告
test-coverage:
	@echo "$(BLUE)[INFO]$(NC) Generating test coverage report..."
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)[SUCCESS]$(NC) Coverage report generated: coverage.html"

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

.PHONY: docker-build docker-run docker-stop docker-clean docker-up docker-down verify

# 构建Docker镜像
docker-build:
	@echo "$(BLUE)[INFO]$(NC) Building Docker image..."
	@docker build -t go-template:latest .
	@echo "$(GREEN)[SUCCESS]$(NC) Docker image built: go-template:latest"

# 运行Docker容器
docker-run:
	@echo "$(BLUE)[INFO]$(NC) Starting Docker container..."
	@docker compose up -d
	@echo "$(GREEN)[SUCCESS]$(NC) Docker container started"

# 停止Docker容器
docker-stop:
	@echo "$(BLUE)[INFO]$(NC) Stopping Docker container..."
	@docker compose down
	@echo "$(GREEN)[SUCCESS]$(NC) Docker container stopped"

docker-up: docker-run

docker-down: docker-stop

verify:
	@go test ./... -count=1
	@go vet ./...
	@go build ./...
	@docker compose config >/dev/null
	@git diff --check

# 清理Docker资源
docker-clean: docker-stop
	@echo "$(BLUE)[INFO]$(NC) Cleaning Docker resources..."
	@if docker image inspect go-template:latest >/dev/null 2>&1; then \
		docker image rm go-template:latest; \
	fi
	@echo "$(GREEN)[SUCCESS]$(NC) Docker resources cleaned"

.PHONY: security-scan

# 安全扫描
security-scan:
	@echo "$(BLUE)[INFO]$(NC) Running security scan..."
	@if ! command -v gosec >/dev/null 2>&1; then \
		echo "$(RED)[ERROR]$(NC) gosec is not installed. Install: go install github.com/securego/gosec/v2/cmd/gosec@latest"; \
		exit 1; \
	fi
	@gosec ./...
	@echo "$(GREEN)[SUCCESS]$(NC) Security scan completed"

# =============================================================================
# 性能测试相关命令
# =============================================================================

.PHONY: bench load-test

# 性能基准测试
bench:
	@echo "$(BLUE)[INFO]$(NC) Running performance benchmarks..."
	@go test -bench=. -benchmem ./...
	@echo "$(GREEN)[SUCCESS]$(NC) Performance benchmarks completed"

# 负载测试
load-test:
	@echo "$(BLUE)[INFO]$(NC) Running load tests..."
	@if ! command -v hey >/dev/null 2>&1; then \
		echo "$(RED)[ERROR]$(NC) hey is not installed. Install: go install github.com/rakyll/hey@latest"; \
		exit 1; \
	fi
	@hey -n 1000 -c 10 $(LOAD_TEST_URL)
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
	@echo "  $(GREEN)tidy$(NC)               - 整理Go模块依赖"
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
	@echo "  $(GREEN)DIR$(NC)                - 目录路径 (用于lint-dir)"
	@echo ""
	@echo "$(YELLOW)示例用法:$(NC)"
	@echo "  make lint-dir DIR=./internal/platform/server"
	@echo "  make lint-fix"
	@echo "  make dev-full"
