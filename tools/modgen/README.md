# 🚀 Echo Admin 代码生成工具

## 功能特性

- ✅ **默认模板生成**: 快速生成标准的CRUD API模块
- ✅ **OpenAPI3支持**: 从OpenAPI3文档生成API模块
- ✅ **测试用例生成**: 自动生成业务逻辑和服务层测试用例
- ✅ **自动依赖注入**: 自动注册到fx.Options
- ✅ **完整文件生成**: 生成biz、service、param、model、code文件
- ✅ **彩色输出**: 友好的命令行界面

## 使用方法

### 1. 默认模板生成

```bash
# 生成用户模块
go run tools/modgen/main.go --name=user

# 生成文章模块
go run tools/modgen/main.go --name=article

# 指定路由前缀
go run tools/modgen/main.go --name=user --route=/api/users

# 强制覆盖已存在文件
go run tools/modgen/main.go --name=user --force

# 生成模块和测试用例
go run tools/modgen/main.go --name=user --tests --force
```

### 2. 从OpenAPI3文档生成

```bash
# 从OpenAPI3文档生成用户模块
go run tools/modgen/main.go --name=user --openapi=doc/openapi.yaml

# 从OpenAPI3文档生成文章模块
go run tools/modgen/main.go --name=article --openapi=doc/openapi.yaml --force

# 从OpenAPI3文档生成模块和测试用例
go run tools/modgen/main.go --name=user --openapi=doc/openapi.yaml --tests --force
```

## 参数说明

| 参数 | 说明 | 示例 |
|------|------|------|
| `--name` | 模块名（必需） | `--name=user` |
| `--route` | 基础路由前缀 | `--route=/api/users` |
| `--openapi` | OpenAPI3文档路径 | `--openapi=doc/openapi.yaml` |
| `--tests` | 生成测试用例 | `--tests` |
| `--force` | 强制覆盖已存在文件 | `--force` |

## 生成的文件

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

## OpenAPI3文档格式

工具支持标准的OpenAPI3格式，包括：

### 基本结构
```yaml
openapi: 3.0.0
info:
  title: API文档
  version: 1.0.0
paths:
  /users:
    get:
      summary: 查询用户
      operationId: findUsers
      parameters:
        - name: page
          in: query
          schema:
            type: integer
      responses:
        "200":
          description: 成功
    post:
      summary: 创建用户
      operationId: createUser
      requestBody:
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/User'
      responses:
        "200":
          description: 成功
components:
  schemas:
    User:
      type: object
      properties:
        id:
          type: integer
        name:
          type: string
        email:
          type: string
```

### 支持的HTTP方法
- `GET` - 查询操作
- `POST` - 创建操作
- `PUT` - 更新操作
- `DELETE` - 删除操作
- `PATCH` - 部分更新操作

### 支持的参数类型
- `query` - 查询参数
- `path` - 路径参数
- `header` - 请求头参数
- `cookie` - Cookie参数

### 支持的数据类型
- `string` - 字符串
- `integer` - 整数
- `number` - 数字
- `boolean` - 布尔值
- `array` - 数组
- `object` - 对象

## 生成代码示例

### 默认模板生成

**业务逻辑层 (biz/user.go)**
```go
package biz

import (
	"context"
	"github.com/NSObjects/echo-admin/internal/api/service/param"
)

// UserUseCase 业务逻辑接口
type UserUseCase interface {
	List(ctx context.Context, p param.UserParam) ([]param.UserResponse, int64, error)
	Create(ctx context.Context, b param.UserBody) (*param.UserResponse, error)
	Update(ctx context.Context, id int64, b param.UserBody) (*param.UserResponse, error)
	Delete(ctx context.Context, id int64) error
	Detail(ctx context.Context, id int64) (*param.UserResponse, error)
}

// UserHandler 业务逻辑处理器
type UserHandler struct {
	// TODO: 注入依赖
}

// NewUserHandler 创建业务逻辑处理器
func NewUserHandler() UserUseCase {
	return &UserHandler{}
}
```

**服务层 (service/user.go)**
```go
package service

import (
	"github.com/NSObjects/echo-admin/internal/api/biz"
	"github.com/NSObjects/echo-admin/internal/api/service/param"
	"github.com/NSObjects/echo-admin/internal/resp"
	"github.com/labstack/echo/v4"
)

type UserController struct {
	user biz.UserUseCase
}

func NewUserController(h biz.UserUseCase) RegisterRouter {
	return &UserController{user: h}
}

func (c *UserController) RegisterRouter(g *echo.Group, m ...echo.MiddlewareFunc) {
	g.GET("/users", c.list).Name = "列表示例"
	g.POST("/users", c.create).Name = "创建示例"
	g.GET("/users/:id", c.detail).Name = "详情示例"
	g.PUT("/users/:id", c.update).Name = "更新示例"
	g.DELETE("/users/:id", c.remove).Name = "删除示例"
}
```

### OpenAPI3生成

**参数结构体 (param/user.go)**
```go
package param

import "time"

// UserCreateRequest 请求结构体
type UserCreateRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

// UserResponse 响应结构体
type UserResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// UserParam 查询参数结构体
type UserParam struct {
	Page  int    `json:"page" form:"page" query:"page"`
	Count int    `json:"count" form:"count" query:"count"`
	Name  string `json:"name" form:"name" query:"name"`
	Email string `json:"email" form:"email" query:"email"`
}
```

## 使用流程

### 1. 生成模块
```bash
# 使用默认模板
go run tools/modgen/main.go --name=user

# 或从OpenAPI3文档生成
go run tools/modgen/main.go --name=user --openapi=doc/openapi.yaml
```

### 2. 实现业务逻辑
编辑 `internal/api/biz/user.go`，实现具体的业务逻辑：

```go
func (h *UserHandler) List(ctx context.Context, p param.UserParam) ([]param.UserResponse, int64, error) {
	// 实现查询逻辑
	return users, total, nil
}

func (h *UserHandler) Create(ctx context.Context, b param.UserBody) (*param.UserResponse, error) {
	// 实现创建逻辑
	return user, nil
}
```

### 3. 配置依赖注入
工具会自动注册到fx.Options，如果没有自动注册，请手动添加：

```go
// internal/api/biz/biz.go
fx.Options(
	// ... 其他模块
	fx.Provide(NewUserHandler),
)

// internal/api/service/service.go
fx.Options(
	// ... 其他模块
	fx.Provide(AsRoute(NewUserController)),
)
```

### 4. 运行测试
```bash
# 运行测试
go test ./internal/api/...

# 启动服务
make run
```

## 注意事项

1. **文件覆盖**: 使用 `--force` 参数会覆盖已存在的文件
2. **OpenAPI3格式**: 确保OpenAPI3文档格式正确
3. **依赖注入**: 检查自动注册是否成功
4. **业务逻辑**: 生成的是模板代码，需要实现具体业务逻辑
5. **参数验证**: 生成的参数结构体包含验证标签，需要配置验证器

## 故障排除

### 1. 编译错误
```bash
# 检查Go模块
go mod tidy

# 检查依赖
go mod download
```

### 2. OpenAPI3解析错误
- 检查YAML/JSON格式是否正确
- 确保文件路径正确
- 检查OpenAPI3版本是否支持

### 3. 文件生成失败
- 检查目标目录是否存在
- 检查文件权限
- 使用 `--force` 参数覆盖

## 贡献指南

1. Fork 项目
2. 创建功能分支
3. 提交更改
4. 推送到分支
5. 创建 Pull Request

## 许可证

MIT License
