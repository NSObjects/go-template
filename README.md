# go-template - 企业级 API 脚手架

[![Go Version](https://img.shields.io/badge/Go-1.24.0%2B-blue.svg)](https://golang.org)
[![Echo Version](https://img.shields.io/badge/Echo-v4.13.4-green.svg)](https://echo.labstack.com)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

## 📖 项目简介

go-template 是一个基于 Echo + Fx 的企业级服务模板，围绕 Clean Architecture 构建。项目使用 [go-kit](https://github.com/NSObjects/go-kit) 作为核心基础库，提供统一的错误处理、中间件、配置管理、日志和响应格式化等能力。

## ✨ 核心特性

- **分层架构**：Service / Biz / Data 三层职责清晰，方便扩展与测试
- **依赖注入**：通过 Uber Fx 组织组件与生命周期管理
- **go-kit 集成**：统一的错误码、中间件、配置热加载、结构化日志
- **数据访问**：集成 GORM Gen，自动生成类型安全的查询代码
- **认证与权限**：内置 JWT 登录流程与 Casbin RBAC 模型
- **工程化工具**：完备的 Makefile 工作流，代码生成器即开即用
- **可观测性**：结构化日志、健康检查、Prometheus 指标支持

## 🚀 快速开始

### 环境准备

- Go 1.24+
- MySQL、Redis 等可选基础设施（可在配置中关闭）
- golangci-lint（可通过 `make install-lint` 安装）

### 初始化项目

```bash
# 1. 安装 muban CLI
go install github.com/NSObjects/go-template/muban@latest

# 2. 使用 CLI 生成项目
muban new -m github.com/acme/awesome-api -o ../awesome-api

# 3. 进入新项目目录
cd ../awesome-api

# 4. 设置开发环境
make dev-setup

# 5. 启动服务
make run
```

### 在本仓库开发

```bash
# 设置开发环境
make dev-setup

# 启动开发服务器
make run-dev

# 运行完整检查（格式化 + lint + 测试）
make dev-check
```

## 🛠️ Makefile 开发工作流

Makefile 整合了所有开发命令，输入 `make help` 查看完整命令列表。

### 基础命令

| 命令 | 说明 |
|------|------|
| `make build` | 构建应用程序到 `bin/app` |
| `make run` | 启动应用程序 |
| `make run-dev` | 开发模式启动（设置 `RUN_ENVIRONMENT=dev`） |
| `make run-test` | 测试模式启动（设置 `RUN_ENVIRONMENT=test`） |
| `make tidy` | 整理 Go 模块依赖 |
| `make push m="message"` | 快速提交代码（vendor + commit + push） |

### 代码质量

| 命令 | 说明 |
|------|------|
| `make fmt` | 格式化代码（gofmt + goimports） |
| `make vet` | 运行 go vet 检查 |
| `make lint` | 运行 golangci-lint 检查 |
| `make lint-strict` | 严格检查（失败时退出） |
| `make lint-fast` | 快速检查（只运行快速 linter） |
| `make lint-fix` | 自动修复可修复的问题 |
| `make lint-dir DIR=./path` | 检查特定目录 |
| `make lint-report` | 生成 XML 格式的 lint 报告 |
| `make install-lint` | 安装 golangci-lint v2.4.0 |

### 测试命令

| 命令 | 说明 |
|------|------|
| `make test` | 运行所有测试（带 race 检测） |
| `make test-verbose` | 详细测试输出 |
| `make test-coverage` | 生成测试覆盖率报告（HTML） |
| `make bench` | 运行性能基准测试 |

### 代码生成

| 命令 | 说明 | 示例 |
|------|------|------|
| `make init-project` | 基于模板生成新项目 | `make init-project MODULE=github.com/acme/demo` |
| `make gen-code` | 生成错误码和文档 | `make gen-code` |
| `make gen-module` | 生成业务模块（基础模板） | `make gen-module NAME=user` |
| `make gen-module-openapi` | 从 OpenAPI 生成模块 | `make gen-module-openapi NAME=article OPENAPI=doc/openapi.yaml` |
| `make gen-module-route` | 生成模块（指定路由） | `make gen-module-route NAME=order ROUTE=/api/v1/orders` |
| `make gen-all-modules` | 从 OpenAPI 生成所有模块 | `make gen-all-modules OPENAPI=doc/openapi.yaml` |
| `make db-gen` | 生成数据库模型和查询（GORM Gen） | `make db-gen` |
| `make db-gen-table` | 生成指定表的模型 | `make db-gen-table TABLE=users` |
| `make db-gen-dynamic` | 生成 Dynamic SQL 查询 | `make db-gen-dynamic` |
| `make gen-all` | 生成所有代码（数据库 + 错误码） | `make gen-all` |

### 开发工作流

| 命令 | 说明 |
|------|------|
| `make dev-setup` | 设置开发环境（tidy + download） |
| `make dev-check` | 运行开发检查（fmt + vet + lint + test） |
| `make dev-full` | 完整开发流程（clean + check + gen-all） |

### Docker 命令

| 命令 | 说明 |
|------|------|
| `make docker-build` | 构建 Docker 镜像 |
| `make docker-run` | 启动 Docker 容器（docker-compose） |
| `make docker-stop` | 停止 Docker 容器 |
| `make docker-clean` | 清理 Docker 资源 |

### 清理命令

| 命令 | 说明 |
|------|------|
| `make clean` | 清理生成的文件（bin、coverage） |
| `make clean-all` | 深度清理（包括 Go cache） |

## 📁 项目结构

```
go-template/
├── cmd/                          # 命令入口
│   ├── fx.go                     # Fx 应用启动 & 模块组装
│   ├── gen.go                    # 代码生成命令
│   └── root.go                   # Cobra 根命令
├── configs/                      # 配置文件
│   └── config.toml               # 应用配置
├── doc/                          # 文档
│   └── openapi.yaml              # OpenAPI 规范
├── internal/                     # 私有应用代码
│   ├── api/                      # API 分层
│   │   ├── biz/                  # 业务逻辑层
│   │   │   ├── biz.go            # Fx 模块注册
│   │   │   └── user.go           # User 业务逻辑
│   │   ├── data/                 # 数据访问层
│   │   │   ├── data.go           # Fx 模块注册
│   │   │   ├── user.go           # User Repository 实现
│   │   │   ├── db/               # 数据库基础设施
│   │   │   ├── model/            # GORM 模型（自动生成）
│   │   │   └── query/            # GORM Gen 查询（自动生成）
│   │   └── service/              # HTTP Handler 层
│   │       ├── service.go        # Fx 模块注册
│   │       ├── user.go           # User Controller
│   │       └── param/            # 请求/响应 DTO
│   ├── code/                     # 应用错误码
│   ├── configs/                  # 配置类型定义
│   ├── model/                    # 共享数据模型
│   ├── pkg/                      # 内部工具包
│   │   └── casbin/               # Casbin 授权配置
│   ├── server/                   # HTTP 服务器
│   │   ├── echo_server.go        # Echo 服务器配置
│   │   ├── config.go             # 服务器配置
│   │   └── middlewares/          # 中间件
│   └── types/                    # 共享类型
├── muban/                        # 代码生成工具
│   ├── codegen/                  # 错误码生成器
│   ├── modgen/                   # 模块生成器
│   ├── project/                  # 项目模板生成器
│   └── newcmd/                   # 命令生成器
├── scripts/                      # 开发脚本
├── main.go                       # 应用入口
├── Makefile                      # 开发工作流
└── docker-compose.yaml           # Docker 编排
```

## 🏗️ API 分层说明

`internal/api` 目录按照 Clean Architecture 分为三层：

### Service 层 (`internal/api/service`)
- 暴露 HTTP 接口，负责参数绑定和校验
- 使用 `go-kit/resp` 统一响应格式
- 只依赖 Biz 层，不包含业务逻辑

### Biz 层 (`internal/api/biz`)
- 实现核心业务用例和领域逻辑
- 定义 Repository 接口（依赖倒置）
- 使用 `go-kit/code` 包装业务错误

### Data 层 (`internal/api/data`)
- 实现 Repository 接口
- 封装数据库访问（GORM Gen）
- 处理数据转换和持久化

## 🔧 go-kit 集成

项目使用 [go-kit](https://github.com/NSObjects/go-kit) 提供统一的基础能力：

```go
import (
    "github.com/NSObjects/go-kit/code"        // 错误码
    "github.com/NSObjects/go-kit/errors"      // 错误包装
    "github.com/NSObjects/go-kit/resp"        // 响应格式化
    "github.com/NSObjects/go-kit/middleware"  // 中间件
    "github.com/NSObjects/go-kit/config"      // 配置管理
    "github.com/NSObjects/go-kit/log"         // 结构化日志
    "github.com/NSObjects/go-kit/validator"   // 验证器
)
```

## 📋 muban CLI 命令

### `muban new`

使用模板生成一个全新的项目骨架：

```bash
# 基础用法
muban new -m github.com/acme/awesome-api -o ../awesome-api

# 指定展示名称并覆盖已存在目录
muban new -m github.com/acme/awesome-api -n "Awesome API" -f
```

### `muban module`

在现有仓库内生成业务模块脚手架：

```bash
# 使用默认模板生成 user 模块
muban module --name=user

# 基于 OpenAPI 生成 article 模块
muban module --name=article --openapi=doc/openapi.yaml

# 从 OpenAPI 一次性生成所有模块
muban module --openapi=doc/openapi.yaml
```

### `muban codegen`

根据错误码常量生成字符串方法或 Markdown 文档：

```bash
# 为 ErrCode 生成字符串方法
muban codegen -t ErrCode -o internal/code/code_string.go

# 生成错误码 Markdown 文档
muban codegen -t ErrCode --doc -o doc/error-code.md
```

### `muban dynamicsql`

使用 GORM Gen 生成动态 SQL 查询接口：

```bash
# 使用默认配置
muban dynamicsql

# 指定配置文件
muban dynamicsql --config=configs/config.local.toml
```

## 🔄 典型开发流程

### 新功能开发

```bash
# 1. 确保环境就绪
make dev-setup

# 2. 生成新模块骨架
make gen-module NAME=product

# 3. 编写业务代码...

# 4. 运行检查
make dev-check

# 5. 提交代码
make push m="feat: add product module"
```

### 数据库变更

```bash
# 1. 修改数据库表结构...

# 2. 重新生成模型和查询代码
make db-gen

# 或只生成特定表
make db-gen-table TABLE=products

# 3. 更新 Repository 实现
```

### CI/CD 集成

```bash
# 在 CI 中使用严格检查
make lint-strict
make test

# 生成覆盖率报告
make test-coverage
```

## 🔒 配置说明

配置文件位于 `configs/config.toml`，支持热加载。主要配置项：

```toml
[system]
port = ":9322"
level = "debug"

[jwt]
enabled = true
secret = "your-secret-key"
skip_paths = ["/api/login", "/api/health"]

[casbin]
enabled = true
skip_paths = ["/api/login", "/api/health"]
admin_users = ["root", "admin"]

[mysql]
enabled = true
host = "127.0.0.1"
port = "3306"
database = "test"
username = "root"
password = "12345678"
```

## 📚 更多文档

- [AGENTS.md](./AGENTS.md) - 详细的开发规范和代码约定
- [internal/server/README.md](./internal/server/README.md) - 服务器配置说明
- [internal/server/middlewares/README.md](./internal/server/middlewares/README.md) - 中间件说明

## 📄 许可证

MIT License
