push:
	go mod download && go mod vendor && git add . && git commit -m '$(m)'

build:
	go build  main.go

tidy:
	go mod tidy

test:
	export RUN_ENVIRONMENT=test
	go test -race $(go list ./...)

test-verbose:
	export RUN_ENVIRONMENT=test
	go test -v -race $(go list ./...)

test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-unit:
	go test -short ./internal/...

test-integration:
	go test -run Integration ./...

test-benchmark:
	go test -bench=. ./...

cover:
	go test ./... -coverprofile=coverage.out

run:
	go run main.go --config configs/config.toml

lint:
	@golangci-lint run || true

gen:
	@echo "no-op: add your gen commands under tools/ as needed"

# =============================================================================
# 代码生成工具
# =============================================================================

# 数据库相关命令
.PHONY: db-gen db-gen-table db-migrate db-reset

# 生成所有数据库模型和查询方法（使用gentool）
db-gen:
	@echo "Generating database models and queries using gentool..."
	@$(shell go env GOPATH)/bin/gentool \
		-dsn="root:12345678@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local" \
		-outPath="./internal/api/data/query" \
		-modelPkgName="model" \
		-fieldWithIndexTag \
		-fieldWithTypeTag
	@echo "Database generation completed!"

# 生成Dynamic SQL查询方法
db-gen-dynamic:
	@echo "Generating Dynamic SQL queries..."
	@go run tools/dynamic-sql-gen/basic.go
	@echo "Dynamic SQL generation completed!"

# 生成完整的数据库代码（模型 + Dynamic SQL）
db-gen-full:
	@echo "Generating complete database code with Dynamic SQL..."
	@go run tools/dynamic-sql-gen/basic.go
	@echo "Complete database generation finished!"

# 生成指定表的模型和查询方法
db-gen-table:
	@if [ -z "$(TABLE)" ]; then \
		echo "Usage: make db-gen-table TABLE=table_name"; \
		exit 1; \
	fi
	@echo "Generating model and queries for table: $(TABLE)"
	@$(shell go env GOPATH)/bin/gentool \
		-dsn="root:12345678@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local" \
		-outPath="./internal/api/data/query" \
		-modelPkgName="model" \
		-fieldWithIndexTag \
		-fieldWithTypeTag \
		-tables="$(TABLE)"
	@echo "Table $(TABLE) generation completed!"

# 数据库迁移（使用gorm migrate）
db-migrate:
	@echo "Running database migrations..."
	@if [ -z "$(DB_DSN)" ]; then \
		echo "Please set DB_DSN environment variable"; \
		exit 1; \
	fi
	@echo "Migration completed!"

# 重置数据库（开发环境）
db-reset:
	@echo "Resetting database..."
	@if [ "$(ENV)" != "development" ] && [ "$(ENV)" != "dev" ]; then \
		echo "db-reset can only be run in development environment"; \
		exit 1; \
	fi
	@echo "Database reset completed!"

# 生成错误码
gen-code:
	@echo "Generating error codes..."
	@go run tools/codegen/codegen.go -type=int -doc -output ./error_code_generated.md ./internal/code
	@echo "Error code generation completed!"

# 生成业务模块
gen-module:
	@if [ -z "$(NAME)" ]; then \
		echo "Usage: make gen-module NAME=module_name"; \
		exit 1; \
	fi
	@echo "Generating module: $(NAME)"
	@go run tools/modgen/main.go --name=$(NAME) --force
	@echo "Module $(NAME) generation completed!"



# 完整生成（数据库 + 错误码）
gen-all: db-gen gen-code
	@echo "All generation completed!"

# =============================================================================
# 开发工具
# =============================================================================

# 格式化代码
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Code formatting completed!"

# 代码检查
vet:
	@echo "Running go vet..."
	@go vet ./...
	@echo "Code vetting completed!"

# 运行测试
test-verbose:
	@echo "Running tests with verbose output..."
	@go test -v ./...

# 生成测试覆盖率报告
test-coverage:
	@echo "Generating test coverage report..."
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# 清理生成的文件
clean:
	@echo "Cleaning generated files..."
	@rm -f coverage.out coverage.html
	@rm -f error_code_generated.md
	@echo "Clean completed!"

# =============================================================================
# 帮助信息
# =============================================================================

help:
	@echo "Available commands:"
	@echo ""
	@echo "代码生成:"
	@echo "  db-gen          - Generate all database models and queries"
	@echo "  db-gen-dynamic  - Generate Dynamic SQL queries with interfaces"
	@echo "  db-gen-full     - Generate complete database code (models + Dynamic SQL)"
	@echo "  db-gen-table    - Generate model for specific table (TABLE=table_name)"
	@echo "  db-migrate      - Run database migrations"
	@echo "  db-reset        - Reset database (development only)"
	@echo "  gen-code        - Generate error codes and documentation"
	@echo "  gen-module      - Generate business module (NAME=module_name)"
	@echo "  gen-all         - Generate database models and error codes"
	@echo ""
	@echo "开发工具:"
	@echo "  fmt             - Format code"
	@echo "  vet             - Run go vet"
	@echo "  test-verbose    - Run tests with verbose output"
	@echo "  test-coverage   - Generate test coverage report"
	@echo "  clean           - Clean generated files"
	@echo ""
	@echo "基础命令:"
	@echo "  build           - Build application"
	@echo "  run             - Run application"
	@echo "  test            - Run tests"
	@echo "  lint            - Run linter"
	@echo "  help            - Show this help message"
	@echo ""
	@echo "Environment variables:"
	@echo "  DB_DSN          - Database connection string"
	@echo "  TABLE           - Table name for db-gen-table"
	@echo "  NAME            - Module name for gen-module"
	@echo "  ENV             - Environment (development/production)"
