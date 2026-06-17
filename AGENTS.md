# AGENTS.md

本仓库是 Go HTTP backend template，采用 Clean-lite vertical slice 架构。做任何改动前，先理解当前代码和文档，不许走捷径。

## 最高优先级

- 不要参考旧 spec；旧 spec 已经失效。
- 当前事实来源是代码、`CONTEXT.md`、`docs/architecture.md` 和 `docs/adr/`。
- 不要恢复旧的 runtime module registry、capability provider、platform module、generated demo user、全局 utils、旧 `internal/code` 或旧横向能力包。
- 每一步都要想清楚：先查真实代码和真实 diff，再下结论。
- 不要为了“架构完整”添加假业务模块、空壳接口或未来可能用到的 provider。

## 架构形状

业务代码按一个业务能力一个竖切模块放在 `internal/<business>/`：

```text
internal/<business>/
  domain/
  usecase/
  http/
  mysql/
```

小模块可以省略不需要的目录。没有真实业务规则时，不要创建 `domain/`。没有持久化时，不要创建 `mysql/`。需要 cache、queue、client、memory adapter 时，在业务模块内或窄命名 infrastructure 包内按真实需求添加。

依赖方向必须向内：

```text
http adapter -> usecase -> domain
mysql adapter -> usecase/domain
boot -> concrete adapters/usecases/server/config
```

`internal/boot` 是 composition root，可以脏，可以接具体实现。业务模块必须干净。

## 业务代码规则

- 先写 `usecase`，再接 HTTP 和 DB。
- `usecase` 拥有 outbound interface，因为 usecase 是调用方。
- HTTP handler 可以直接依赖 `*usecase.Usecase`；不要默认创建 inbound interface。
- DB row model 只能留在 adapter 包，不能进入 `usecase` 或 `domain`。
- HTTP DTO 只能留在 `http` adapter，不能进入 `usecase` 或 `domain`。
- 业务错误用 `internal/apperr` 或模块内错误包装；HTTP 状态映射只在 `internal/server/httpresp`。
- 请求元数据使用 `internal/requestctx`，不要把 Echo context 传进 usecase。

## 禁止导入

`domain` 只允许导入 Go 标准库。

`usecase` 禁止导入：

- Echo 或其他 HTTP framework
- GORM、Redis、Kafka、具体 driver
- `internal/configs`
- `internal/server`
- `internal/server/httpresp`
- 同业务模块的 `http` 或 `mysql` adapter
- 其他业务模块

这些规则由 `internal/archtest` 保护。改架构时必须同步更新守卫测试。

## Server 和配置边界

- `internal/server` 只负责 Echo、middleware、system routes、HTTP error rendering、graceful shutdown。
- `internal/server` 不得 import 业务模块。
- `internal/configs` 不得 import 业务模块。
- 配置文件只保留当前真实运行需要的配置。不要重新加入默认 MySQL、Redis、Kafka、CORS 等假能力配置。
- JWT 默认不启用，且不能有默认 secret。业务启用 JWT 时必须显式提供 secret。
- CORS 默认不启用。业务需要跨域时再显式配置。

## 新增业务的步骤

1. 在 `internal/<business>/usecase` 写输入、输出、流程和 outbound interface。
2. 有真实业务规则时，在 `internal/<business>/domain` 写实体、值对象和不变量。
3. 在 `internal/<business>/http` 写 Echo handler、DTO 和 `Register`。
4. 需要持久化时，在 `internal/<business>/mysql` 写 store 和 row model。
5. 只在 `internal/boot` 组装 concrete adapter、usecase 和 route。
6. 先测 usecase，再测 adapter；不要只靠 HTTP 测试覆盖业务逻辑。

## 测试和验证

改 Go 代码后至少运行：

```bash
go test ./... -count=1
go vet ./...
go build ./...
git diff --check
```

改 Docker 或 Compose 后额外运行：

```bash
docker compose config
docker build --progress=plain --target final -t go-template-arch-check .
```

如果构建了临时镜像或启动了临时容器，验证后清理掉。

## 工作方式

- 优先删除浅层、未接入、误导性的代码，而不是包一层兼容壳。
- 不要把“将来可能需要”写进默认模板。
- 不要引入全局共享 `utils`。
- 不要把业务流程藏进 middleware、config、server 或 boot。
- 不要保留开发期 debug 输出、默认密码、默认 secret 或示例账号。
- 改动完成前必须重新扫描旧概念残留，并用真实命令验证。
