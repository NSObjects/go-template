# go-template

`go-template` 是一个 module-first 的 Go API 模板。它不追求生成大量框架代码，而是提供一条后续业务 API 开发可以持续复用的路径：业务模块独立表达领域、用例和适配器，平台层统一处理 HTTP runtime、配置、日志、认证、观测和基础设施资源。

当前项目使用 Go 1.26，HTTP runtime 基于 Echo v5，依赖装配使用 `samber/do`。默认配置不依赖外部服务，可以直接用 memory store 启动；MySQL、Redis、MongoDB 和 tracing 都是可选 capability。

## 架构边界

仓库统一使用这些架构词：

- `business module`：业务模块，位于 `internal/modules/<module>`。
- `platform`：平台运行时代码，位于 `internal/platform`。
- `capability`：可复用基础设施能力，例如 MySQL、Redis、MongoDB、tracing、logging。
- `composition root`：启动装配层，位于 `internal/boot`。

核心目录：

```text
.
├── cmd/                         # CLI 入口，默认读取 configs/config.toml
├── configs/                     # 静态配置示例和配置说明
├── internal/boot/               # composition root，负责装配资源、模块和路由
├── internal/modules/            # business modules
│   ├── customer/
│   ├── product/
│   └── salesorder/
├── internal/platform/           # HTTP runtime、配置、错误、请求响应、基础设施能力
├── k8s/                         # Kubernetes 示例
├── docker-compose.yaml          # 本地开发 Compose
├── Dockerfile                   # scratch runtime image
└── Makefile                     # 构建、测试、验证命令
```

不要把本项目改回 `service/biz/data` 分层心智。业务路由只在 `internal/boot` 通过 `NewModule`、`Provide`、`Route` 显式装配，`internal/platform/server` 不 import 业务模块。

## 快速开始

本地直接运行：

```bash
make run
```

等价于：

```bash
go run main.go --config configs/config.toml
```

默认监听 `:9322`。启动后可访问：

```bash
curl http://127.0.0.1:9322/api/health
curl http://127.0.0.1:9322/api/info
curl http://127.0.0.1:9322/api/ready
curl http://127.0.0.1:9322/api/capabilities
```

常用开发命令：

```bash
make build
make test
make lint
make verify
```

`make verify` 会运行 `go test ./... -count=1`、`go vet ./...`、`go build ./...`、`docker compose config` 和 `git diff --check`。`make test` 会运行 race detector。

## 配置

默认配置文件是 `configs/config.toml`。配置支持 TOML、YAML、JSON，未知配置字段会导致启动失败。环境变量可以覆盖同名配置项，前缀由 `configs.EnvPrefix` 定义，当前是 `GO_TEMPLATE_`。

最小配置可直接启动 API。可选 capability 默认关闭：

- `mysql.enabled=false`
- `redis.enabled=false`
- `mongodb.enabled=false`
- `tracing.enabled=false`
- `jwt.enabled=false`
- `http.cors.enabled=false`

启用外部资源后，它们会进入 `/api/ready` 和 `/api/capabilities`。`/api/ready` 只返回整体 readiness；详细 capability 状态由 `/api/capabilities` 返回。

更完整的配置说明见 `configs/README.md`。

## 本地 Compose

默认 `docker compose up` 只启动 API，公开端口默认绑定 `127.0.0.1`：

```bash
docker compose up
```

MySQL、Redis、MongoDB 和 Jaeger 在 `resources` profile 下：

```bash
cp env.example .env
docker compose --profile resources up
```

`.env` 只给 Docker Compose 做本地开发端口和凭据替换；应用本身读取真实进程环境变量，不会自动加载 `.env` 文件。

如果在 API 容器内启用 MySQL，可以使用 Compose service hostname：

```text
go_template:go_template_dev_password@tcp(mysql:3306)/go_template?charset=utf8mb4&parseTime=true&loc=Local
```

如果在宿主机直接运行应用并连接 Compose MySQL，把 DSN 里的 `mysql:3306` 改为 `127.0.0.1:<GO_TEMPLATE_MYSQL_PORT>`。

资源服务使用 Compose 命名卷保存本地开发数据，普通 `docker compose down` 不会删除
MySQL、Redis 和 MongoDB 数据。需要重置本地资源状态时，显式运行
`docker compose down -v`。

## Kubernetes

`k8s/deployment.yaml` 使用镜像内置的 `/configs/config.toml` 作为静态配置基线，通过 `go-template-runtime` ConfigMap 里的 `GO_TEMPLATE_*` 环境变量覆盖运行环境差异。

当前示例把 JWT secret、MySQL DSN、Redis password 和 MongoDB URI 从 `go-template-secrets` 读取，不写入 ConfigMap。正式环境发布前应把 `go-template:latest` 替换为带不可变 tag 的镜像地址。

## 业务模块开发

新增业务模块时，按这个顺序写：

1. `internal/modules/<module>/domain`
2. `internal/modules/<module>/usecase`
3. `internal/modules/<module>/adapters/memory`
4. `internal/modules/<module>/http`
5. `internal/boot/business.go`
6. 行为测试

推荐结构：

```text
internal/modules/order/
├── adapters/
│   └── memory/
│       └── store.go
├── domain/
│   └── order.go
├── http/
│   ├── handler.go
│   └── handler_test.go
└── usecase/
    ├── usecase.go
    └── usecase_test.go
```

职责边界：

- `domain` 放领域对象、构造函数、领域错误和不可变业务规则。
- `usecase` 放输入输出、业务流程和 usecase-owned outbound interface。
- `adapters/<adapter>` 实现 usecase 定义的 interface。
- `http` 只做请求解析、validator 校验、DTO 转换、调用 usecase、统一响应。
- `internal/boot` 负责选择 adapter、装配依赖、挂载路由。

HTTP adapter 应复用 `internal/platform/server/httpreq` 和 `internal/platform/server/httpresp`。错误语义优先使用 `internal/platform/apperr`。

## Store 和外部资源

默认开发路径应先提供 `adapters/memory`，保证本地开发和测试开箱可跑。真实外部资源由 `internal/boot` 统一打开，业务 adapter 复用 boot 已经加载的资源。

`product` 是当前真实存储样板：

- `internal/modules/product/adapters/memory`
- `internal/modules/product/adapters/mysql`
- `internal/boot/business.go` 中的 `newProductStore`

当 `mysql.enabled=false` 时使用 memory store；当 `mysql.enabled=true` 时，boot 注入配置好的 `*gorm.DB`，并选择 MySQL adapter。usecase 只依赖自己的 `Store` interface，不直接依赖 GORM。

## 跨模块依赖

跨模块依赖不要直接 import 对方 store 或 adapter。由消费方 usecase 定义小接口，再在 `internal/boot` 写 bridge。

当前示例是 `salesorder` 依赖 `customer` 和 `product`：

- `salesorder` 用 `CustomerExists(ctx, id)` 检查客户存在性。
- `salesorder` 用 `FindProduct(ctx, id)` 返回 `ProductSnapshot{Exists, Active}`，区分“商品不存在”和“商品存在但不可用”。
- bridge 位于 `internal/boot/business.go`，业务模块之间不直接互相持有 adapter。

这个模式让跨模块语义由消费方拥有，避免把重要业务状态压成一个裸 `bool` 或泄漏到平台层。

## HTTP API

系统路由：

| Method | Path | 用途 |
| --- | --- | --- |
| `GET` | `/api/health` | 进程存活检查 |
| `GET` | `/api/info` | 应用名称、版本和当前时间 |
| `GET` | `/api/ready` | 整体 readiness |
| `GET` | `/api/capabilities` | capability 状态详情 |

示例业务路由：

| Module | Method | Path |
| --- | --- | --- |
| customer | `POST` | `/api/customers` |
| customer | `GET` | `/api/customers` |
| customer | `GET` | `/api/customers/:id` |
| customer | `PATCH` | `/api/customers/:id` |
| product | `POST` | `/api/products` |
| product | `GET` | `/api/products` |
| product | `GET` | `/api/products/:id` |
| product | `PATCH` | `/api/products/:id` |
| salesorder | `POST` | `/api/sales-orders` |
| salesorder | `GET` | `/api/sales-orders` |
| salesorder | `GET` | `/api/sales-orders/:id` |

## 测试和验证

行为变化必须有测试。新增业务模块至少补：

- domain 或 usecase 的核心业务测试；
- HTTP adapter 成功路径测试；
- 一个关键失败路径测试；
- boot 装配测试，确保模块能被默认 runtime 挂载。

现有可参考测试：

- `internal/modules/product/http/handler_test.go`
- `internal/modules/customer/http/handler_test.go`
- `internal/modules/salesorder/http/handler_test.go`
- `internal/boot/business_test.go`
- `internal/platform/server/server_test.go`

提交前优先运行：

```bash
make test
make lint
make build
make verify
```

如果修改了依赖，必须运行：

```bash
go mod tidy
```

## 改名清单

基于模板创建新项目时，至少同步修改：

- `go.mod` module path；
- `.golangci.yml` 里的 `local-prefixes`；
- `configs.EnvPrefix`；
- `configs/config.toml` 和 `env.example` 的应用名、环境变量前缀；
- Docker Compose service/image 名称；
- Dockerfile 中保留的默认环境变量；
- K8s 示例里的镜像、名称和环境变量；
- README 中的项目名和示例命令。

不要同时保留旧前缀和新前缀。新项目阶段应直接统一命名，避免留下长期兼容分支。

## 进一步阅读

- `AGENTS.md`：仓库级开发约束。
- `configs/README.md`：配置项和环境变量覆盖规则。
- `internal/boot/README.md`：composition root 和模块装配说明。
- `internal/platform/server/README.md`：HTTP runtime 边界说明。
- `internal/platform/server/middlewares/README.md`：middleware 行为说明。
