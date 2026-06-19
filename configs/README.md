# 配置说明

项目默认只加载一个静态配置文件：

```go
cfg, err := configs.Load("configs/config.toml")
```

配置文件支持 TOML、YAML、JSON，格式由文件后缀识别。无后缀时按 TOML 解析，未知后缀会在启动时失败。未知配置字段会在启动时失败。环境变量可以覆盖同名配置项。

## 当前配置项

```toml
[app]
name = "go-template"
version = "dev"

[system]
port = ":9322"
level = 1 # 1=debug, 2=online

[log]
format = "console" # console, json
output = "stdout" # stdout, stderr
caller = false

[http]
recovery_disabled = false
request_context_disabled = false
request_log_disabled = false
gzip_disabled = false

[http.cors]
enabled = false
allow_origins = []
allow_methods = []
allow_headers = []
allow_credentials = false
expose_headers = []
max_age_seconds = 0

[jwt]
enabled = false
skip_paths = ["/api/health", "/api/info", "/api/ready"]

[mysql]
enabled = false
dsn = ""

[redis]
enabled = false
address = ""

[mongodb]
enabled = false
uri = ""
database = ""

[tracing]
enabled = false
endpoint = ""
protocol = "grpc" # grpc, http
insecure = true
```

HTTP middleware、MySQL、Redis、MongoDB、JWT、Jaeger tracing 都是 configured infrastructure resources。默认配置可直接运行 API；启用外部依赖后必须提供连接配置，并会进入 `/api/ready` 和 `/api/capabilities`。

业务代码放在 `internal/modules/<module>`，平台运行时代码放在 `internal/platform`。
业务 adapter 复用 boot 已经加载的资源：usecase 定义 usecase-owned outbound interface，adapter 使用配置好的 MySQL/Redis/MongoDB client 实现它，然后业务 module 用 `boot.NewModule`、`boot.Provide` 和 `boot.Route` 声明 adapter、usecase、handler 和 route。`product` 模块提供了 `adapters/mysql` 作为最小 GORM store 样板。

## 环境变量覆盖

```bash
export GO_TEMPLATE_APP_NAME=go-template
export GO_TEMPLATE_APP_VERSION=dev
export GO_TEMPLATE_SYSTEM_PORT=:9322
export GO_TEMPLATE_SYSTEM_LEVEL=1
export GO_TEMPLATE_LOG_FORMAT=console
export GO_TEMPLATE_LOG_OUTPUT=stdout
export GO_TEMPLATE_LOG_CALLER=false
export GO_TEMPLATE_HTTP_RECOVERY_DISABLED=false
export GO_TEMPLATE_HTTP_REQUEST_CONTEXT_DISABLED=false
export GO_TEMPLATE_HTTP_REQUEST_LOG_DISABLED=false
export GO_TEMPLATE_HTTP_GZIP_DISABLED=false
export GO_TEMPLATE_HTTP_CORS_ENABLED=false
export GO_TEMPLATE_HTTP_CORS_ALLOW_ORIGINS=
export GO_TEMPLATE_HTTP_CORS_ALLOW_METHODS=
export GO_TEMPLATE_HTTP_CORS_ALLOW_HEADERS=
export GO_TEMPLATE_HTTP_CORS_ALLOW_CREDENTIALS=false
export GO_TEMPLATE_HTTP_CORS_EXPOSE_HEADERS=
export GO_TEMPLATE_HTTP_CORS_MAX_AGE_SECONDS=0
export GO_TEMPLATE_JWT_ENABLED=false
export GO_TEMPLATE_JWT_SECRET=
export GO_TEMPLATE_MYSQL_ENABLED=false
export GO_TEMPLATE_MYSQL_DSN=
export GO_TEMPLATE_REDIS_ENABLED=false
export GO_TEMPLATE_REDIS_ADDRESS=
export GO_TEMPLATE_MONGODB_ENABLED=false
export GO_TEMPLATE_MONGODB_URI=
export GO_TEMPLATE_TRACING_ENABLED=false
export GO_TEMPLATE_TRACING_ENDPOINT=
export GO_TEMPLATE_TRACING_PROTOCOL=grpc
export GO_TEMPLATE_TRACING_INSECURE=true
```

JWT 默认不启用，且默认 secret 为空。业务启用 JWT 时必须显式提供 secret。
推荐通过 `GO_TEMPLATE_JWT_SECRET` 或 Secret 管理系统提供 JWT secret，不要把真实 secret 写入 ConfigMap。
`app.name` 和 `app.version` 会出现在 `/api/info` 响应中。
`system.level` 只接受 `1` 或 `2`。
`log.format` 只接受 `console` 或 `json`；`log.output` 只接受 `stdout` 或 `stderr`。
`log.caller` 打开后会在每条日志中加入 caller 字段。
`http.*_disabled` 可关闭模板默认安装的 HTTP middleware。
`http.cors.enabled=true` 时必须配置 `http.cors.allow_origins`。
`http.cors.allow_credentials=true` 时 `allow_origins` 不能包含 `*`。
环境变量中的列表项使用逗号分隔，例如 `GO_TEMPLATE_HTTP_CORS_ALLOW_ORIGINS=https://app.example.com,https://admin.example.com`。环境变量前缀由 `configs.EnvPrefix` 定义。
`/api/health` 是进程存活，`/api/ready` 只返回整体 readiness，`/api/capabilities` 返回 logging、MySQL、Redis、MongoDB、Jaeger tracing 的 disabled/available/unavailable 状态。

`docker-compose.yaml` 只用于本地开发，公开端口默认绑定 `127.0.0.1`。默认只启动 API；MySQL、Redis、MongoDB 和 Jaeger 在 `resources` profile 下，使用 `docker compose --profile resources up` 启动。Compose 端口和本地 MySQL 凭据可通过 `.env` 中的 `GO_TEMPLATE_API_PORT`、`GO_TEMPLATE_MYSQL_PORT`、`GO_TEMPLATE_REDIS_PORT`、`GO_TEMPLATE_MONGODB_PORT`、`GO_TEMPLATE_JAEGER_UI_PORT`、`GO_TEMPLATE_OTLP_GRPC_PORT`、`GO_TEMPLATE_OTLP_HTTP_PORT`、`GO_TEMPLATE_MYSQL_PASSWORD` 和 `GO_TEMPLATE_MYSQL_ROOT_PASSWORD` 覆盖。资源服务使用 Compose 命名卷保存本地开发数据；需要重置 MySQL、Redis 和 MongoDB 时，显式运行 `docker compose down -v`。
