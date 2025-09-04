# =============================================================================
# Go Template Project Makefile
# 整合了所有tools目录下的工具，提供完整的开发工作流
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
DEFAULT_OPENAPI := doc/openapi.yaml
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
	@go mod download && go mod vendor && git add . && git commit -m '$(m)'
	@echo "$(GREEN)[SUCCESS]$(NC) Code committed with message: $(m)"

# =============================================================================
# 代码质量工具
# =============================================================================

.PHONY: fmt vet lint test test-verbose test-coverage clean

# 格式化代码
fmt:
	@echo "$(BLUE)[INFO]$(NC) Formatting code..."
	@go fmt ./...
	@echo "$(GREEN)[SUCCESS]$(NC) Code formatting completed"

# 代码检查
vet:
	@echo "$(BLUE)[INFO]$(NC) Running go vet..."
	@go vet ./...
	@echo "$(GREEN)[SUCCESS]$(NC) Code vetting completed"

# 代码检查（使用golangci-lint）
lint:
	@echo "$(BLUE)[INFO]$(NC) Running linter..."
	@golangci-lint run || true
	@echo "$(GREEN)[SUCCESS]$(NC) Linting completed"

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

.PHONY: gen-code gen-module gen-module-tests gen-module-openapi gen-module-openapi-tests gen-module-route db-gen db-gen-table gen-all

# 生成错误码和文档
gen-code:
	@echo "$(BLUE)[INFO]$(NC) Generating error codes and documentation..."
	@go run tools/codegen/codegen.go -type=int -doc -output ./internal/code/error_code_generated.md ./internal/code
	@echo "$(GREEN)[SUCCESS]$(NC) Error code generation completed"

# 生成业务模块（基础模板）
gen-module:
	@if [ -z "$(NAME)" ]; then \
		echo "$(RED)[ERROR]$(NC) Usage: make gen-module NAME=module_name"; \
		exit 1; \
	fi
	@echo "$(BLUE)[INFO]$(NC) Generating module: $(NAME)"
	@go run tools/modgen/main.go --name=$(NAME) --force
	@echo "$(GREEN)[SUCCESS]$(NC) Module $(NAME) generation completed"

# 生成模块和测试用例
gen-module-tests:
	@if [ -z "$(NAME)" ]; then \
		echo "$(RED)[ERROR]$(NC) Usage: make gen-module-tests NAME=module_name"; \
		exit 1; \
	fi
	@echo "$(BLUE)[INFO]$(NC) Generating module with tests: $(NAME)"
	@go run tools/modgen/main.go --name=$(NAME) --tests --force
	@echo "$(GREEN)[SUCCESS]$(NC) Module $(NAME) with tests generation completed"

# 从OpenAPI文档生成模块
gen-module-openapi:
	@if [ -z "$(NAME)" ]; then \
		echo "$(RED)[ERROR]$(NC) Usage: make gen-module-openapi NAME=module_name [OPENAPI=doc/openapi.yaml]"; \
		exit 1; \
	fi
	@OPENAPI_FILE=$${OPENAPI:-$(DEFAULT_OPENAPI)}; \
	echo "$(BLUE)[INFO]$(NC) Generating module from OpenAPI: $(NAME) ($$OPENAPI_FILE)"; \
	go run tools/modgen/main.go --name=$(NAME) --openapi=$$OPENAPI_FILE --force
	@echo "$(GREEN)[SUCCESS]$(NC) Module $(NAME) generated from OpenAPI"

# 从OpenAPI生成模块和测试用例
gen-module-openapi-tests:
	@if [ -z "$(NAME)" ]; then \
		echo "$(RED)[ERROR]$(NC) Usage: make gen-module-openapi-tests NAME=module_name [OPENAPI=doc/openapi.yaml]"; \
		exit 1; \
	fi
	@OPENAPI_FILE=$${OPENAPI:-$(DEFAULT_OPENAPI)}; \
	echo "$(BLUE)[INFO]$(NC) Generating module from OpenAPI with tests: $(NAME) ($$OPENAPI_FILE)"; \
	go run tools/modgen/main.go --name=$(NAME) --openapi=$$OPENAPI_FILE --tests --force
	@echo "$(GREEN)[SUCCESS]$(NC) Module $(NAME) with tests generated from OpenAPI"

# 生成模块（指定路由）
gen-module-route:
	@if [ -z "$(NAME)" ] || [ -z "$(ROUTE)" ]; then \
		echo "$(RED)[ERROR]$(NC) Usage: make gen-module-route NAME=module_name ROUTE=/api/module_name"; \
		exit 1; \
	fi
	@echo "$(BLUE)[INFO]$(NC) Generating module with custom route: $(NAME) -> $(ROUTE)"
	@go run tools/modgen/main.go --name=$(NAME) --route=$(ROUTE) --force
	@echo "$(GREEN)[SUCCESS]$(NC) Module $(NAME) with route $(ROUTE) generation completed"

# 生成数据库模型和查询
db-gen:
	@echo "$(BLUE)[INFO]$(NC) Generating database models and queries..."
	@$(shell go env GOPATH)/bin/gentool \
		-dsn=$(DEFAULT_DSN) \
		-outPath="./internal/api/data/query" \
		-modelPkgName="model" \
		-fieldWithIndexTag \
		-fieldWithTypeTag
	@echo "$(GREEN)[SUCCESS]$(NC) Database generation completed"

# 生成指定表的模型和查询方法
db-gen-table:
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

# 生成Dynamic SQL查询方法
db-gen-dynamic:
	@echo "$(BLUE)[INFO]$(NC) Generating Dynamic SQL queries..."
	@go run tools/dynamic-sql-gen/basic.go
	@echo "$(GREEN)[SUCCESS]$(NC) Dynamic SQL generation completed"

# 完整生成（数据库 + 错误码）
gen-all: db-gen gen-code
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
dev-full: clean dev-check gen-all
	@echo "$(GREEN)[SUCCESS]$(NC) Full development workflow completed"

# =============================================================================
# 清理和维护
# =============================================================================

# 清理生成的文件
clean:
	@echo "$(BLUE)[INFO]$(NC) Cleaning generated files..."
	@rm -f coverage.out coverage.html
	@rm -f internal/code/error_code_generated.md
	@rm -rf $(BIN_DIR)
	@echo "$(GREEN)[SUCCESS]$(NC) Clean completed"

# 深度清理
clean-all: clean
	@echo "$(BLUE)[INFO]$(NC) Deep cleaning..."
	@go clean -cache
	@go clean -modcache
	@echo "$(GREEN)[SUCCESS]$(NC) Deep clean completed"

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
	@echo "  $(GREEN)push$(NC)               - 提交代码 (需要设置 m=commit_message)"
	@echo ""
	@echo "$(YELLOW)代码质量:$(NC)"
	@echo "  $(GREEN)fmt$(NC)                - 格式化代码"
	@echo "  $(GREEN)vet$(NC)                - 运行go vet检查"
	@echo "  $(GREEN)lint$(NC)               - 运行golangci-lint检查"
	@echo "  $(GREEN)test$(NC)               - 运行所有测试"
	@echo "  $(GREEN)test-verbose$(NC)       - 运行详细测试"
	@echo "  $(GREEN)test-coverage$(NC)      - 生成测试覆盖率报告"
	@echo ""
	@echo "$(YELLOW)代码生成:$(NC)"
	@echo "  $(GREEN)gen-code$(NC)                    - 生成错误码和文档"
	@echo "  $(GREEN)gen-module$(NC)                  - 生成业务模块 (NAME=module_name)"
	@echo "  $(GREEN)gen-module-tests$(NC)            - 生成模块和测试用例 (NAME=module_name)"
	@echo "  $(GREEN)gen-module-openapi$(NC)          - 从OpenAPI生成模块 (NAME=name [OPENAPI=path])"
	@echo "  $(GREEN)gen-module-openapi-tests$(NC)    - 从OpenAPI生成模块和测试 (NAME=name [OPENAPI=path])"
	@echo "  $(GREEN)gen-module-route$(NC)            - 生成模块（指定路由） (NAME=name ROUTE=path)"
	@echo "  $(GREEN)db-gen$(NC)                      - 生成数据库模型和查询"
	@echo "  $(GREEN)db-gen-table$(NC)                - 生成指定表模型 (TABLE=table_name)"
	@echo "  $(GREEN)db-gen-dynamic$(NC)              - 生成Dynamic SQL查询"
	@echo "  $(GREEN)gen-all$(NC)                     - 生成所有代码（数据库+错误码）"
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
	@echo "$(YELLOW)环境变量:$(NC)"
	@echo "  $(GREEN)NAME$(NC)               - 模块名 (用于gen-module)"
	@echo "  $(GREEN)ROUTE$(NC)              - 路由路径 (用于gen-module-route)"
	@echo "  $(GREEN)OPENAPI$(NC)            - OpenAPI文档路径 (默认: doc/openapi.yaml)"
	@echo "  $(GREEN)TABLE$(NC)              - 表名 (用于db-gen-table)"
	@echo "  $(GREEN)m$(NC)                  - 提交消息 (用于push)"
	@echo ""
	@echo "$(YELLOW)示例用法:$(NC)"
	@echo "  make gen-module NAME=user"
	@echo "  make gen-module-tests NAME=product"
	@echo "  make gen-module-openapi NAME=article"
	@echo "  make gen-module-openapi-tests NAME=user"
	@echo "  make gen-module-route NAME=order ROUTE=/api/v1/orders"
	@echo "  make db-gen-table TABLE=users"
	@echo "  make push m=\"feat: add user module\""
	@echo "  make dev-full"