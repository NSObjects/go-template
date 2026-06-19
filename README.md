# go-template

一个 module-first 的 Go API 模板。业务代码放在 `internal/modules/<module>`，
平台运行时代码放在 `internal/platform`，启动装配集中在 `internal/boot`。

## 常用命令

```bash
make run
make verify
make test
golangci-lint run
```

`make verify` 会运行 `go test ./... -count=1`、`go vet ./...`、
`go build ./...`、`docker compose config` 和 `git diff --check`。

## 新增业务模块

新增一个模块时，按这个顺序写最顺：

1. `internal/modules/<module>/domain`：领域对象、构造函数、领域错误。
2. `internal/modules/<module>/usecase`：用例输入输出和 usecase-owned outbound interface。
3. `internal/modules/<module>/adapters/memory`：本地开发和单测用的 store。
4. `internal/modules/<module>/http`：Echo handler 和 `Register(group, handler)`。
5. `internal/boot/business.go`：用 `NewModule`、`Provide`、`Route` 装配 store、usecase、handler 和 route。
6. 同步补 usecase/http/boot 的行为测试。

跨模块依赖不要直接 import 对方 store。由 consumer usecase 定义小接口，再在
`internal/boot` 写 adapter bridge。`salesorder` 对 `product` 的 lookup 就是这个模式。

## 真实存储样板

默认配置使用 memory store，开箱可跑。需要真实 MySQL store 时，参考：

- `internal/modules/product/adapters/mysql`
- `internal/boot/business.go` 里的 `newProductStore`

`product` 模块在 `mysql.enabled=true` 时使用 MySQL/GORM store，否则使用 memory store。
这给后续模块提供一个最小可抄路径：usecase 仍只依赖自己的 `Store` interface，
adapter 通过 boot 复用已经配置好的 `*gorm.DB`。

## 本地 Compose

默认 `docker compose up` 只启动 API。MySQL、Redis、MongoDB 和 Jaeger 是可选资源：

```bash
cp env.example .env
docker compose --profile resources up
```

如果把 `.env` 里的 `GO_TEMPLATE_MYSQL_ENABLED=true` 打开，API 容器会读取同一个
`.env`，并使用 Compose service hostname 连接 MySQL。直接在宿主机运行应用时，把
DSN 里的 `mysql:3306` 改成 `127.0.0.1:<端口>`。

环境变量前缀集中在 `configs.EnvPrefix`。新项目改名时，优先从那里、`go.mod`、
Docker/K8s 名称和 `.golangci.yml` 的 `local-prefixes` 开始改。
