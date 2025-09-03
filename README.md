# Echo Admin - 企业级Go Web应用框架

[![Go Version](https://img.shields.io/badge/Go-1.24.0+-blue.svg)](https://golang.org)
[![Echo Version](https://img.shields.io/badge/Echo-v4.13.4-green.svg)](https://echo.labstack.com)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

## 📖 项目简介

Echo Admin 是一个基于 [Echo](https://echo.labstack.com) 框架构建的企业级Go Web应用框架，采用Clean Architecture架构设计，集成了完整的认证授权、数据库管理、代码生成等企业级功能。

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

### 环境要求

- Go 1.24.0+
- MySQL 8.0+
- Redis 6.0+
- Make

### 安装步骤

```bash
# 1. 克隆项目
git clone <repository-url>
cd echo-admin

# 2. 安装依赖
go mod download

# 3. 配置数据库
# 编辑 configs/config.toml 文件，配置数据库连接信息

# 4. 初始化数据库
make db-init

# 5. 启动服务
make run
```

## 📁 项目结构

```
echo-admin/
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
│   │   ├── main.go        # 主程序
│   │   ├── openapi_templates.go  # OpenAPI模板
│   │   ├── test_templates.go     # 测试模板
│   │   └── openapi_gen.go        # OpenAPI解析
│   └── encrypt.go         # 加密工具
├── scripts/               # 脚本文件
├── main.go                # 主入口文件
├── Makefile               # 构建脚本
└── README.md              # 项目说明
```

## 🔧 配置说明

### 主配置文件 (configs/config.toml)

```toml
# 应用配置
[app]
name = "echo-admin"
version = "1.0.0"
debug = true
port = 9322

# 数据库配置
[database]
driver = "mysql"
host = "localhost"
port = 3306
username = "root"
password = "password"
database = "echo_admin"
charset = "utf8mb4"
max_idle_conns = 10
max_open_conns = 100
conn_max_lifetime = "1h"

# Redis配置
[redis]
host = "localhost"
port = 6379
password = ""
database = 0
pool_size = 10

# JWT配置
[jwt]
secret = "your-secret-key"
expire = 3600
skip_paths = ["/api/health", "/api/login"]

# Casbin配置
[casbin]
model_path = "configs/rbac_model.conf"
policy_path = "configs/rbac_policy.csv"
```

## 🛠️ 开发指南

### 快速开发流程

1. **设计API接口** - 在 `doc/openapi.yaml` 中定义API规范
2. **生成代码** - 使用 `modgen` 工具生成完整的业务模块
3. **实现业务逻辑** - 在生成的 `biz` 层实现具体的业务逻辑
4. **自定义接口** - 在 `service` 层自定义HTTP接口
5. **注册路由** - 在 `RegisterRouter` 方法中注册API路由
6. **测试验证** - 运行生成的测试用例验证功能

### 代码生成工具

#### 基于OpenAPI3文档生成（推荐）

```bash
# 生成完整的API模块（包含测试用例）
go run tools/modgen/main.go --name=user --openapi=doc/openapi.yaml --tests --force

# 只生成代码，不生成测试
go run tools/modgen/main.go --name=user --openapi=doc/openapi.yaml --force
```

#### 基于默认模板生成

```bash
# 生成基础模块
go run tools/modgen/main.go --name=product --tests --force
```

### 生成的文件结构

```
internal/api/
├── biz/
│   ├── user.go           # 业务逻辑实现
│   └── user_test.go      # 业务逻辑测试
├── service/
│   ├── user.go           # HTTP控制器
│   ├── user_test.go      # 控制器测试
│   └── param/
│       └── user.go       # 请求参数结构
├── data/
│   └── model/
│       └── user.go       # 数据模型
└── code/
    └── user.go           # 错误码定义
```

### OpenAPI3文档规范

项目使用OpenAPI3规范定义API接口，支持以下特性：

- **参数验证**: 自动生成validator标签
- **错误码映射**: 自动生成错误码定义
- **测试用例**: 自动生成完整的测试用例
- **Mock支持**: 集成testify/mock框架
- **RESTful设计**: 根据HTTP方法自动设置参数位置

示例OpenAPI定义：

```yaml
paths:
  /api/users:
    get:
      operationId: find users
      parameters:
        - name: page
          in: query
          schema:
            type: integer
            minimum: 1
        - name: count
          in: query
          schema:
            type: integer
            minimum: 1
            maximum: 100
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

## 📚 API文档

### RESTful API设计规范

- **请求格式**: JSON
- **响应格式**: 统一JSON响应格式
- **状态码**: 标准HTTP状态码
- **认证**: JWT Bearer Token
- **权限**: 基于Casbin的RBAC权限控制
- **参数位置**: 根据RESTful标准自动设置（GET用query，POST/PUT用body等）

### 统一响应格式

```json
{
  "code": 200,
  "msg": "success",
  "data": {},
  "timestamp": 1640995200
}
```

### 错误码规范

```go
// 成功
const (
    Success = 200
)

// 客户端错误
const (
    ParamError = 400
    Unauthorized = 401
    Forbidden = 403
    NotFound = 404
)

// 服务端错误
const (
    InternalError = 500
    DBError = 502
)
```

## 🧪 测试

### 运行测试

```bash
# 运行所有测试
make test

# 运行详细测试
make test-verbose

# 运行特定模块测试
go test ./internal/api/service -v

# 运行特定测试用例
go test ./internal/api/service -v -run TestUserController_list

# 运行测试并生成覆盖率报告
make test-coverage
```

### 测试特性

- **自动生成测试用例**: 基于OpenAPI文档自动生成测试用例
- **Mock支持**: 集成testify/mock框架，支持依赖注入
- **参数验证测试**: 自动测试参数验证规则
- **HTTP测试**: 完整的HTTP请求测试
- **业务逻辑测试**: 独立的业务逻辑测试

### 测试覆盖率

当前项目测试覆盖率情况：

| 模块 | 覆盖率 | 状态 |
|------|--------|------|
| `internal/code` | 51.3% | ✅ 良好 |
| `internal/resp` | 86.1% | ✅ 优秀 |
| `internal/api/service` | 75.2% | ✅ 良好 |
| 其他模块 | 待完善 | 🔄 进行中 |

**总体覆盖率**: 65.8% (基于已测试模块)

## 🚀 部署

### 开发环境

```bash
# 启动开发服务器
make run

# 热重载开发
make dev
```

### Docker部署

1. **构建镜像**
```bash
docker build -t echo-admin .
```

2. **运行容器**
```bash
docker run -d \
  --name echo-admin \
  -p 9322:9322 \
  -v $(pwd)/configs:/app/configs \
  echo-admin
```

### 生产环境部署

1. **编译二进制文件**
```bash
make build
```

2. **使用systemd管理服务**
```bash
# 创建服务文件
sudo vim /etc/systemd/system/echo-admin.service

# 启动服务
sudo systemctl start echo-admin
sudo systemctl enable echo-admin
```

## 🔍 故障排除

### 常见问题

1. **数据库连接失败**
   - 检查数据库配置
   - 确认数据库服务是否启动
   - 验证用户名密码是否正确

2. **JWT认证失败**
   - 检查JWT密钥配置
   - 确认跳过路径配置
   - 验证令牌格式

3. **权限控制问题**
   - 检查Casbin模型配置
   - 确认策略文件路径
   - 验证用户角色配置

4. **代码生成问题**
   - 检查OpenAPI文档格式
   - 确认Go版本兼容性
   - 验证模板语法

### 日志查看

```bash
# 查看应用日志
tail -f logs/app.log

# 查看错误日志
tail -f logs/error.log

# 查看测试日志
go test -v ./internal/api/service
```

## 📋 开发命令

### Make命令

```bash
# 开发相关
make run          # 启动服务
make dev          # 开发模式（热重载）
make build        # 构建二进制文件
make clean        # 清理构建文件

# 测试相关
make test         # 运行所有测试
make test-verbose # 详细测试输出
make test-coverage # 测试覆盖率

# 数据库相关
make db-init      # 初始化数据库
make db-migrate   # 数据库迁移
make db-reset     # 重置数据库

# 代码生成
make gen-code     # 生成错误码文档
make gen-db       # 生成数据库模型
```

### 工具命令

```bash
# 模块生成
go run tools/modgen/main.go --name=user --openapi=doc/openapi.yaml --tests --force

# 加密工具
go run tools/encrypt.go --text="your-password"
```

## 🤝 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
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

## 📞 联系方式

- 项目维护者: [Your Name]
- 邮箱: [your.email@example.com]
- 项目链接: [https://github.com/your-username/echo-admin](https://github.com/your-username/echo-admin)

---

**注意**: 这是一个企业级框架，请在生产环境中使用前进行充分测试。建议在开发环境中先熟悉框架特性和代码生成工具的使用。