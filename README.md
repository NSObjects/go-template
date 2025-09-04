# go-template - API快速开发框架

[![Go Version](https://img.shields.io/badge/Go-1.24.0+-blue.svg)](https://golang.org)
[![Echo Version](https://img.shields.io/badge/Echo-v4.13.4-green.svg)](https://echo.labstack.com)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

## 📖 项目简介

go-template 是一个基于 [Echo](https://echo.labstack.com) 框架构建的企业级Go Web应用框架，采用Clean Architecture架构设计，集成了完整的认证授权、数据库管理、智能代码生成等企业级功能。

### ✨ 核心特性

- 🏗️ **Clean Architecture**: 清晰的分层架构，易于维护和扩展
- 🔐 **JWT认证**: 完整的JWT令牌认证系统
- 🛡️ **Casbin权限控制**: 基于RBAC的细粒度权限管理
- 🗄️ **多数据库支持**: MySQL、Redis、MongoDB、Kafka
- 🚀 **智能代码生成**: 基于OpenAPI3文档自动生成CRUD代码和API接口
- 📦 **依赖注入**: 基于Uber FX的依赖注入框架
- 🔧 **配置管理**: 灵活的配置管理系统
- 📊 **日志系统**: 结构化日志记录
- 🧪 **测试支持**: 完整的测试框架和Mock支持
- 🌐 **RESTful API**: 标准化的RESTful API设计

## 🚀 快速开始

```bash
# 1. 克隆项目
git clone git@github.com:NSObjects/go-template.git
cd go-template

# 2. 安装依赖
go mod download

# 3. 启动服务
make run
```

## 📁 项目结构

```
go-template/
├── cmd/                    # 命令行工具和入口
│   ├── fx.go              # FX依赖注入配置
│   ├── gen.go             # 代码生成工具
│   └── root.go            # 根命令
├── configs/               # 配置文件
│   ├── config.toml        # 主配置文件
│   └── rbac_model.conf    # Casbin RBAC模型
├── doc/                   # 文档
│   └── openapi.yaml       # OpenAPI3规范文档
├── internal/              # 内部包
│   ├── api/               # API层
│   │   ├── biz/           # 业务逻辑层
│   │   ├── data/          # 数据访问层
│   │   │   ├── model/     # 数据模型
│   │   │   └── query/     # 查询接口
│   │   └── service/       # 服务层（控制器）
│   │       └── param/     # 请求参数结构
│   ├── code/              # 错误码定义
│   ├── configs/           # 配置管理
│   ├── log/               # 日志系统
│   ├── resp/              # 响应处理
│   ├── server/            # HTTP服务器
│   │   └── middlewares/   # 中间件
│   └── utils/             # 工具函数
├── tools/                 # 开发工具
│   ├── modgen/            # 模块生成器
│   ├── dynamic-sql-gen/   # 动态SQL生成器
│   └── codegen/           # 错误码生成器
├── examples/              # 使用示例
├── main.go                # 主入口文件
├── Makefile               # 构建脚本
└── README.md              # 项目说明
```

## 🔄 开发流程

### 1. 快速生成API模块

#### 方法一：使用默认模板生成

```bash
# 生成用户模块（包含完整的CRUD操作）
go run tools/modgen/main.go --name=user --force

# 生成文章模块
go run tools/modgen/main.go --name=article --force

# 生成模块并包含测试用例
go run tools/modgen/main.go --name=product --tests --force
```

#### 方法二：基于OpenAPI3文档生成（推荐）

```bash
# 从OpenAPI3文档生成用户模块
go run tools/modgen/main.go --name=user --openapi=doc/openapi.yaml --force

# 生成模块并包含测试用例
go run tools/modgen/main.go --name=user --openapi=doc/openapi.yaml --tests --force
```

#### 生成的文件结构

每个模块会生成以下文件：

```
internal/
├── api/
│   ├── biz/
│   │   ├── {name}.go          # 业务逻辑层
│   │   └── {name}_test.go     # 业务逻辑测试用例（--tests时生成）
│   ├── service/
│   │   ├── {name}.go          # 服务层（控制器）
│   │   ├── {name}_test.go     # 服务层测试用例（--tests时生成）
│   │   └── param/
│   │       └── {name}.go      # 参数结构体
│   └── data/
│       └── model/
│           └── {name}.go      # 数据模型
└── code/
    └── {name}.go              # 错误码定义
```

### 2. 快速生成数据库映射

#### 方法一：生成所有表的模型和查询方法

```bash
# 生成所有数据库模型和查询方法
make db-gen

# 生成指定表的模型
make db-gen-table TABLE=users
```

#### 方法二：生成Dynamic SQL查询（推荐）

```bash
# 生成Dynamic SQL查询方法
make db-gen-dynamic

# 生成完整的数据库代码（模型 + Dynamic SQL）
make db-gen-full
```

#### Dynamic SQL特性

- **通用方法**: 所有表都生成相同的查询方法
- **类型安全**: 所有生成的代码都是类型安全的，编译时检查
- **模板表达式**: 支持 if/else, where, set, for 等高级功能
- **占位符**: `@@table` 自动替换为表名，`@param` 绑定参数

#### 生成的查询接口

```go
// 通用查询接口 - 适用于所有模型的基础CRUD操作
type ICommonQuery interface {
    GetByID(id uint) (gen.T, error)
    GetByIDs(ids []uint) ([]gen.T, error)
    CountRecords() (int64, error)
    Exists(id uint) (bool, error)
    DeleteByID(id uint) error
    DeleteByIDs(ids []uint) error
}

// 分页查询接口 - 适用于需要分页的模型
type IPaginationQuery interface {
    GetPage(offset, limit int, orderBy string) ([]gen.T, error)
    GetPageWithCondition(condition string, offset, limit int, orderBy string) ([]gen.T, error)
}

// 搜索查询接口 - 适用于需要搜索功能的模型
type ISearchQuery interface {
    Search(field, keyword string) ([]gen.T, error)
    SearchMultiple(field1, field2, keyword string) ([]gen.T, error)
}
```

### 3. 完整的开发流程

1. **设计API接口** - 在 `doc/openapi.yaml` 中定义API规范
2. **生成API模块** - 使用 `modgen` 工具生成完整的业务模块
3. **生成数据库映射** - 使用 `make db-gen-dynamic` 生成Dynamic SQL查询
4. **实现业务逻辑** - 在生成的 `biz` 层实现具体的业务逻辑
5. **运行测试** - 使用 `make test` 验证功能
6. **启动服务** - 使用 `make run` 启动服务

## 🛠️ 开发命令

### 代码生成

```bash
# 生成API模块
go run tools/modgen/main.go --name=user --force

# 生成数据库映射
make db-gen-dynamic

# 生成错误码文档
make gen-code
```

### 测试和运行

```bash
# 运行测试
make test

# 启动服务
make run

# 构建应用
make build
```

### 代码质量

```bash
# 格式化代码
make fmt

# 代码检查
make vet

# 清理生成的文件
make clean
```

## 📚 使用示例

### 创建用户管理模块

#### 1. 定义OpenAPI规范

在 `doc/openapi.yaml` 中添加用户相关接口：

```yaml
paths:
  /users:
    get:
      operationId: find users
      parameters:
        - name: page
          in: query
          schema:
            type: integer
        - name: count
          in: query
          schema:
            type: integer
    post:
      operationId: createUser
      requestBody:
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/User'
components:
  schemas:
    User:
      type: object
      required:
        - name
        - email
      properties:
        name:
          type: string
          minLength: 2
          maxLength: 50
        email:
          type: string
          format: email
```

#### 2. 生成API模块

```bash
# 生成用户模块
go run tools/modgen/main.go --name=user --openapi=doc/openapi.yaml --tests --force
```

#### 3. 生成数据库映射

```bash
# 生成Dynamic SQL查询
make db-gen-dynamic
```

#### 4. 实现业务逻辑

编辑 `internal/api/biz/user.go`：

```go
func (h *UserHandler) List(ctx context.Context, p param.UserParam) ([]param.UserResponse, int64, error) {
    // 使用Dynamic SQL查询
    users, err := h.q.User.GetPage(p.Offset(), p.Limit(), "id ASC")
    if err != nil {
        return nil, 0, code.WrapDatabaseError(err, "query user list")
    }
    
    total, err := h.q.User.CountRecords()
    if err != nil {
        return nil, 0, code.WrapDatabaseError(err, "count users")
    }
    
    return convertToResponses(users), total, nil
}
```

### 使用Dynamic SQL查询

```go
// 基础查询
users, err := q.User.GetByID(1)
users, err := q.User.GetByIDs([]uint{1, 2, 3})
count, err := q.User.CountRecords()

// 分页查询
users, err := q.User.GetPage(0, 10, "id ASC")
users, err := q.User.GetPageWithCondition("status = 1", 0, 10, "id ASC")

// 搜索查询
users, err := q.User.Search("name", "admin")
users, err := q.User.SearchMultiple("name", "email", "admin")

// 高级查询
users, err := q.User.FilterWithCondition("status = 1 AND created_at > '2023-01-01'")
users, err := q.User.GetByField("email", "admin@example.com")
```

## 🔧 核心特性

### 智能代码生成

- **OpenAPI3支持**: 从OpenAPI3文档自动生成完整的API模块
- **默认模板**: 快速生成标准的CRUD API模块
- **测试用例**: 自动生成业务逻辑和服务层测试用例
- **依赖注入**: 自动注册到fx.Options

### Dynamic SQL查询

- **通用方法**: 所有表都生成相同的查询方法
- **类型安全**: 编译时检查，避免运行时错误
- **模板表达式**: 支持if/else, where, set, for等高级功能
- **占位符**: `@@table`自动替换为表名，`@param`绑定参数

### 统一响应格式

```json
{
  "code": 200,
  "msg": "success",
  "data": {},
  "timestamp": 1640995200
}
```

## 🚀 部署

### 开发环境

```bash
# 启动开发服务器
make run

# 构建应用
make build
```

### Docker部署

```bash
# 构建镜像
docker build -t go-template .

# 运行容器
docker run -d \
  --name go-template \
  -p 9322:9322 \
  -v $(pwd)/configs:/app/configs \
  go-template
```

## 🔍 故障排除

### 常见问题

1. **代码生成失败**
   - 检查OpenAPI文档格式是否正确
   - 确认Go版本兼容性
   - 使用 `--force` 参数覆盖已存在文件

2. **数据库连接失败**
   - 检查数据库配置
   - 确认数据库服务是否启动
   - 验证用户名密码是否正确

3. **测试失败**
   - 检查依赖是否正确安装
   - 确认数据库连接正常
   - 查看详细错误信息

### 日志查看

```bash
# 查看应用日志
tail -f logs/app.log

# 查看错误日志
tail -f logs/error.log

# 查看测试日志
go test -v ./internal/api/service
```

## 🤝 贡献指南

1. Fork 项目
2. 创建特性分支
3. 提交更改
4. 推送到分支
5. 打开 Pull Request

### 代码规范

- 遵循Go官方代码规范
- 使用gofmt格式化代码
- 编写完整的测试用例
- 更新相关文档

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🙏 致谢

- [Echo](https://echo.labstack.com) - 高性能Go Web框架
- [GORM](https://gorm.io) - Go ORM库
- [Casbin](https://casbin.org) - 权限控制库
- [Uber FX](https://uber-go.github.io/fx/) - 依赖注入框架
- [Testify](https://github.com/stretchr/testify) - 测试框架

