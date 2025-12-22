# Project Context

## Purpose

go-template 是一个基于 Echo + Fx 的企业级服务模板，围绕 Clean Architecture 构建。项目旨在提供一个开箱即用的 API 脚手架，包含：

- **分层架构**：Service / Biz / Data 三层职责清晰，方便扩展与测试
- **依赖注入**：通过 Uber Fx 组织组件与生命周期管理
- **统一基础设施**：使用 go-kit 提供统一的错误码、中间件、配置热加载、结构化日志
- **数据访问**：集成 GORM Gen，自动生成类型安全的查询代码
- **认证与权限**：内置 JWT 登录流程与 Casbin RBAC 模型
- **工程化工具**：完备的 Makefile 工作流，代码生成器即开即用

## Tech Stack

### 核心框架
- **Go 1.24+**：编程语言
- **Echo v4**：HTTP Web 框架
- **Uber Fx**：依赖注入框架
- **Cobra**：命令行工具框架

### 数据存储
- **GORM + GORM Gen**：ORM 和类型安全的查询生成器
- **MySQL**：关系型数据库
- **Redis**：缓存和会话存储
- **MongoDB**：文档数据库（可选）
- **Kafka**：消息队列（可选）

### 基础设施库
- **github.com/NSObjects/go-kit**：核心基础库
  - `go-kit/code`：错误码管理
  - `go-kit/errors`：错误包装
  - `go-kit/resp`：统一响应格式
  - `go-kit/middleware`：中间件（JWT、Casbin、Recovery、Logger）
  - `go-kit/config`：配置管理（支持热加载）
  - `go-kit/log`：结构化日志
  - `go-kit/validator`：请求验证

### 认证与授权
- **JWT (golang-jwt/jwt/v5)**：身份认证
- **Casbin v2**：RBAC 权限控制

### 可观测性
- **OpenTelemetry**：分布式追踪
  - `go.opentelemetry.io/otel`：OpenTelemetry 核心库
  - `go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho`：Echo 中间件集成
  - 支持 OTLP 和 stdout 导出器
  - 可配置采样率

### 开发工具
- **golangci-lint**：代码静态分析
- **go-playground/validator**：数据验证

## Project Conventions

### Code Style

#### 格式化
- 使用 `gofmt` 和 `goimports` 进行代码格式化
- 运行 `make fmt` 自动格式化代码
- 提交前必须通过 `make lint` 检查

#### 命名约定
- **包名**：小写，简短，有意义（如 `biz`, `data`, `service`）
- **接口名**：使用描述性名称，通常以 `-er` 结尾（如 `UserRepository`, `UserUseCase`）
- **结构体**：使用驼峰命名（如 `UserHandler`, `UserController`）
- **函数/方法**：使用驼峰命名，公开方法首字母大写
- **常量**：使用全大写，单词间用下划线分隔（如 `ErrCodeUserNotFound`）

#### 错误处理
- 所有错误必须使用 `go-kit/code` 或 `go-kit/errors` 包装
- 使用语义化的错误包装函数：
  - `code.WrapDatabaseError()`：数据库错误
  - `code.WrapValidationError()`：验证错误
  - `code.WrapNotFoundError()`：资源未找到
  - `errors.WrapCode()`：自定义错误码
- Service 层直接返回错误，由 ErrorHandler 中间件统一处理

#### 响应格式
- 所有成功响应必须使用 `go-kit/resp` 包的方法：
  - `resp.ListDataResponse()`：列表数据
  - `resp.OneDataResponse()`：单条数据
  - `resp.OperateSuccess()`：操作成功
  - `resp.SuccessJSON()`：自定义成功响应

### Architecture Patterns

#### Clean Architecture 三层架构

**Service 层** (`internal/api/service`)
- **职责**：HTTP 传输层，负责请求绑定、参数验证和响应格式化
- **禁止**：包含业务逻辑或数据访问逻辑
- **依赖**：只依赖 Biz 层

**Biz 层** (`internal/api/biz`)
- **职责**：业务逻辑层，封装领域工作流和业务规则
- **原则**：通过 Repository 接口访问数据（依赖倒置）
- **禁止**：直接依赖 DataManager 或数据库连接
- **依赖**：定义 Repository 接口，依赖 Data 层实现

**Data 层** (`internal/api/data`)
- **职责**：数据访问层，实现 Repository 接口
- **实现**：使用 GORM Gen 进行类型安全的数据库操作
- **依赖**：依赖 `db.DataManager` 获取数据库连接

#### 依赖注入模式
- 使用 Uber Fx 进行依赖注入
- 通过构造函数注入依赖，避免全局变量
- 接口定义在调用方（Biz 层），实现在被调用方（Data 层）

#### Repository 模式
- Biz 层定义 Repository 接口
- Data 层实现 Repository 接口
- 便于单元测试（可轻松 mock Repository）

### Testing Strategy

#### 测试类型
- **单元测试**：测试单个函数或方法的逻辑
- **集成测试**：测试多个组件协作（如 Repository 与数据库）
- **表驱动测试**：使用 table-driven tests 覆盖多种场景

#### 测试要求
- Service 层测试：验证请求绑定和响应格式
- Biz 层测试：mock Repository，测试业务逻辑和边界情况
- Data 层测试：使用测试数据库进行集成测试
- 提交前运行 `make test` 确保所有测试通过

#### 测试覆盖率
- 运行 `make test-coverage` 生成覆盖率报告
- 关键业务逻辑应达到 80%+ 覆盖率

### Git Workflow

#### 提交规范
- 使用 `make push m="message"` 快速提交（自动执行 vendor + commit + push）
- 提交信息使用约定格式：
  - `feat:`：新功能
  - `fix:`：修复 bug
  - `docs:`：文档更新
  - `refactor:`：重构
  - `test:`：测试相关
  - `chore:`：构建/工具相关

#### 分支策略
- `main`：主分支，稳定版本
- `develop`：开发分支（如适用）
- 功能分支：`feature/xxx`
- 修复分支：`fix/xxx`

#### 代码审查
- 提交前运行 `make dev-check`（包含 fmt + vet + lint + test）
- 确保代码通过 golangci-lint 检查

## Domain Context

### 用户管理模块
项目包含一个示例的 User 模块，展示了完整的三层架构实现：
- **Service 层**：`internal/api/service/user.go` - HTTP 处理器
- **Biz 层**：`internal/api/biz/user.go` - 业务逻辑
- **Data 层**：`internal/api/data/user.go` - 数据访问

### 配置管理
- 配置文件：`configs/config.toml`
- 支持环境变量覆盖
- 使用 `go-kit/config` 实现配置热加载
- 配置结构定义在 `internal/configs/config.go`

### 代码生成
项目包含多个代码生成工具（`muban/` 目录）：
- **codegen**：错误码生成器
- **modgen**：业务模块生成器（支持 OpenAPI）
- **project**：项目模板生成器
- **newcmd**：命令生成器

### 数据库模型生成
- 使用 GORM Gen 自动生成类型安全的查询代码
- 模型文件：`internal/api/data/model/*.gen.go`（自动生成，禁止手动编辑）
- 查询文件：`internal/api/data/query/*.gen.go`（自动生成，禁止手动编辑）

## Important Constraints

### 技术约束
- **Go 版本**：必须使用 Go 1.24 或更高版本
- **自动生成文件**：以下文件由工具自动生成，禁止手动编辑：
  - `internal/api/data/model/*.gen.go`
  - `internal/api/data/query/*.gen.go`
  - `internal/code/*_generated.go`
- **依赖注入**：必须通过 Fx 进行依赖注入，禁止使用全局变量
- **错误处理**：必须使用 `go-kit/code` 包装错误，禁止直接返回原始错误

### 架构约束
- **分层依赖**：Service → Biz → Data，禁止跨层或反向依赖
- **接口定义**：Repository 接口必须在 Biz 层定义，Data 层实现
- **业务逻辑**：业务规则必须在 Biz 层，禁止在 Service 或 Data 层处理业务逻辑
- **数据访问**：所有数据库操作必须通过 Repository 接口，禁止在 Biz 层直接访问数据库

### 代码质量约束
- **Lint 检查**：所有代码必须通过 `make lint` 检查
- **测试要求**：新功能必须包含测试用例
- **格式化**：提交前必须运行 `make fmt`

## External Dependencies

### 核心依赖
- **github.com/NSObjects/go-kit**：核心基础库，提供统一的基础设施能力
- **go.uber.org/fx**：依赖注入框架，管理组件生命周期
- **github.com/labstack/echo/v4**：HTTP Web 框架

### 数据存储服务
- **MySQL**：主数据库，存储业务数据
- **Redis**：缓存和会话存储（可选）
- **MongoDB**：文档数据库（可选）
- **Kafka**：消息队列（可选）

### 外部服务集成
- **JWT**：用于身份认证，配置在 `configs/config.toml` 的 `[jwt]` 部分
- **Casbin**：用于 RBAC 权限控制，模型文件在 `configs/rbac_model.conf`
- **OpenTelemetry Collector**：分布式追踪收集器（可选），配置在 `configs/config.toml` 的 `[otel]` 部分

### 开发工具依赖
- **golangci-lint**：代码静态分析工具
- **GORM Gen**：数据库查询代码生成器
