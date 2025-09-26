# go-template - 企业级 API 脚手架

[![Go Version](https://img.shields.io/badge/Go-1.24.0%2B-blue.svg)](https://golang.org)
[![Echo Version](https://img.shields.io/badge/Echo-v4.13.4-green.svg)](https://echo.labstack.com)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

## 📖 项目简介

go-template 是一个基于 Echo + Fx 的服务模板，围绕 Clean Architecture 构建，提供完备的认证授权、数据库访问、可观测性与工程化工具链，帮助团队快速落地生产可用的 API 服务。

## ✨ 核心特性

- **分层架构**：业务、数据、接口职责清晰，方便扩展与测试
- **依赖注入**：通过 Uber Fx 组织组件与生命周期
- **数据访问**：集成 GORM、数据库迁移与查询封装
- **认证与权限**：内置 JWT 登录流程与 Casbin RBAC 模型
- **工程化工具**：Makefile、代码生成器、lint/test 流水线即开即用
- **可观测性**：结构化日志、健康检查、Prometheus 指标支持

## 🚀 快速开始

### 环境准备

- Go 1.24+
- MySQL、Redis 等可选基础设施（可在配置中关闭）

### 初始化项目

```bash
# 1. 安装 CLI（可在任何目录执行）
go install github.com/NSObjects/go-template/muban@latest

# 2. 使用 CLI 生成项目（无需预先下载模板仓库）
muban new -m github.com/acme/awesome-api -o ../awesome-api

# 3. 进入新项目目录并启动服务
cd ../awesome-api
make dev-setup
make run
```

### muban CLI 命令一览

#### `muban new`

使用模板生成一个全新的项目骨架。

- `-m, --module`：新项目的 Go Module 路径（必填）
- `-o, --output`：生成项目的目标目录，默认使用模块名
- `-n, --name`：项目展示名称，用于 README、LICENSE 等
- `-f, --force`：目标目录存在时覆盖

```bash
# 最常见的项目初始化
muban new -m github.com/acme/awesome-api -o ../awesome-api

# 指定展示名称并覆盖已存在目录
muban new -m github.com/acme/awesome-api -n "Awesome API" -f
```

#### `muban new module`

在现有仓库内生成业务模块脚手架，可选基于 OpenAPI 自动生成 service/biz/data 代码。提供 `--openapi` 时：

- 未指定 `--name` 会一次性生成 OpenAPI 中的所有模块
- 指定 `--name` 则只生成对应模块

未提供 `--openapi` 时会使用默认模板生成单个模块。

- `-n, --name`：模块名，例如 user、article
- `--route`：自定义基础路由前缀（默认根据模块名推导）
- `--openapi`：OpenAPI3 文档路径，用于自动生成 handler 和 DTO
- `--tests`：同时生成 Table-Driven 风格的测试用例
- `-f, --force`：覆盖已有文件

```bash
# 使用默认模板生成 user 模块
muban new module --name=user

# 基于 OpenAPI 生成 article 模块并附带测试
muban new module --name=article --openapi=doc/openapi.yaml --tests

# 基于 OpenAPI 一次性生成所有模块
muban new module --openapi=doc/openapi.yaml
```

#### `muban codegen`

根据错误码常量生成字符串方法或 Markdown 文档，帮助维护错误码体系。

- `-t, --type`：需要处理的常量类型列表（必填，可逗号分隔）
- `-o, --output`：输出文件路径
- `--doc`：生成 Markdown 文档而非 Go 代码
- `--trimprefix`：去除常量公共前缀
- `--tags`：指定编译标签

```bash
# 为 ErrCode 生成字符串方法
muban codegen -t ErrCode -o internal/pkg/errors/code_string.go

# 生成错误码 Markdown 文档
muban codegen -t ErrCode --doc -o doc/error-code.md
```

#### `muban dynamicsql`

读取配置中的数据库连接，使用 GORM Gen 生成动态 SQL 查询接口。

- `--config`：配置文件路径（默认 `configs/config.toml`）

```bash
# 使用默认配置生成基础查询接口
muban dynamicsql

# 指定自定义配置文件
muban dynamicsql --config=configs/config.local.toml
```

### 常用 Makefile 命令

| 命令 | 说明 |
| --- | --- |
| `make dev` | 启动热加载开发环境 |
| `make test` | 执行单元测试 |
| `make lint` | 运行 golangci-lint |
| `make db-gen` | 根据数据库生成 GORM 代码 |
| `make gen-code` | 生成错误码和文档 |

## 🧰 生成新项目

使用 `muban` CLI 可以把模板复制成新的仓库：

```bash
# 在目标目录生成新项目
muban new -m github.com/acme/awesome-api -o ../awesome-api

# 自定义展示名称或覆盖目录
muban new -m github.com/acme/awesome-api \
  --name="Awesome API" \
  -o ../awesome-api \
  -f
```

如果你正在本仓库中开发 CLI，也可以直接运行源码：

```bash
go run ./muban -- new -m github.com/acme/awesome-api -o ../awesome-api
```

或者通过 Makefile 包装：

```bash
make init-project MODULE=github.com/acme/awesome-api OUTPUT=../awesome-api
```

## 📁 项目结构

```
go-template/
├── cmd/                # 命令入口与 FX 组合
├── configs/            # 配置文件与示例
├── doc/                # OpenAPI 等规范文件
├── internal/           # 业务代码 (api/biz/data/service)
├── scripts/            # 开发脚本
├── sql/                # 数据库迁移
├── muban/              # 项目 CLI 与代码生成器
└── Makefile            # 常用任务
```

### API 分层说明

`internal/api` 目录按照 Clean Architecture 分为三层：

- **service 层**：暴露 HTTP/RPC 接口，负责参数绑定、校验、错误转换以及中间件集成，只依赖 biz 层并通过响应模型向外输出数据。
- **biz 层**：实现核心业务用例，组织领域服务、事务控制和跨模块协作，不直接依赖底层技术细节，而是通过接口与 data 层交互。
- **data 层**：封装数据库、缓存和第三方服务访问，提供 biz 层所需的仓储实现，并统一处理连接池、重试、监控等基础能力。

## 🛠️ 开发建议

- 使用 `make gen-module` 家族命令快速生成模块骨架
- 将 `make lint` 和 `make test` 集成到 CI/CD 中
- 结合 `doc/openapi.yaml` 与代码生成器保持接口与文档一致

## 📄 许可证

MIT License
