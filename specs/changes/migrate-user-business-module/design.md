# Design: Migrate User Business Module

## Technical Approach

新增 `internal/modules/user` 作为第一个真实业务模块。`user` 模块自己声明 HTTP entry points 和抽象 storage capability requirement；它只依赖 `user.storage` 这个能力，不命名 MySQL、MongoDB 或其他具体数据库。

同时扩展平台 capability 装配模型，让一个能力可以由不同 provider 满足。应用通过配置或默认规则选择 `user.storage` provider，业务模块声明和 HTTP entry points 在 provider 切换时保持不变。OpenAPI 和旧 generated wiring 不参与装配。

## Requirements Mapping

| Requirement | Technical Approach | Files / Modules |
|---|---|---|
| User Module Inclusion | 在应用模块列表中显式追加 `user` business module；单独测试 custom assembly 省略 user 时不会暴露 user。 | `internal/app/modules.go`, `internal/app/modules_test.go`, `internal/modules/user/module.go` |
| User Entry Points Are Declared by the Module | `user.Descriptor()` 返回 5 个 HTTP entry points；平台 HTTP adapter 将它们注册成 `/api/users` 相关路由。 | `internal/modules/user/module.go`, `internal/modules/user/routes.go`, `internal/modules/user/module_test.go`, `internal/modules/user/http_test.go` |
| User Storage Capability Is Declared | `user` module 声明 required capability `user.storage`；装配阶段匹配 selected/default provider，并在 report 中记录 provider。 | `internal/modules/user/module.go`, `internal/platform/module/*`, `internal/modules/user/storage/*` |
| User Storage Provider Is Switchable | 平台 capability selection 支持同一能力的 provider 切换；切换 provider 不改变 user module descriptor 或 user HTTP route descriptors。 | `internal/platform/module/*`, `internal/platform/app/app.go`, `internal/app/modules.go`, `internal/modules/user/storage/*` |
| Existing User HTTP Behavior Remains Available | 模块路由 handler 复用现有 `service.UserController`，保持 binding、validation、response envelope 和 error handler 路径。 | `internal/modules/user/routes.go`, `internal/api/service/user.go`, `internal/modules/user/http_test.go` |
| OpenAPI And Generated Output Do Not Drive User Assembly | user assembly 只来自 `internal/modules/user` descriptor 和 app inclusion list；legacy/generated files 不会自动激活 user。 | `internal/modules/user/module.go`, `internal/app/modules.go`, `internal/platform/app/app_test.go` |

## Existing Patterns Observed

- Naming: Go 包名短小单数或能力名，例如 `mysql`、`redis`、`module`、`http`；新业务包使用 `user`，user storage provider 使用 `storage/memory`、`storage/mysql`。
- File structure: Phase 1 已建立 `internal/platform/*`、`internal/capabilities/*`、`internal/app/modules.go`；业务模块应落在新的 `internal/modules/<name>`，业务特定 provider 放在该模块下。
- State/data flow: 当前 app assembly 接收 `[]module.Module`，module descriptor 提供 requirements 和 entry points，HTTP route 由 `platform/http` 注册到 Echo。Phase 2 会扩展 assembly option，让 capability requirement 可以按 provider 选择。
- Error handling: Service handler 直接返回错误，平台 Echo error handler 统一转换；成功响应走 `resp` envelope。
- Testing: 现有测试偏 Go 标准库断言，平台模块测试使用 fake/static module；新增测试应沿用这个风格。
- Configuration: `configs/config.toml` 默认外部 MySQL disabled。调整后默认 user storage 使用本地 default provider，因此默认 app 不因 MySQL disabled 而失败；显式选择 MySQL provider 时才要求 MySQL 可用。provider 选择使用 `configs.Config.Capabilities.Providers`，从 `[capabilities.providers]` 读取 capability 到 provider 的映射。

## Architecture Decisions

### Decision: Put User Ownership In `internal/modules/user`

Decision type: Compose

Chose a new business module package over adding more wiring to `internal/api/service` because:

- 目标是让业务 module 成为 source of truth，而不是继续让 service/biz/data layer 决定应用组成。
- `internal/modules/user` 可以把 required capabilities 和 entry points 放在一个可复制的业务模块示例里。

Alternatives considered:

- Keep `service.Register` as the user assembly path - rejected because it preserves the old layer-first wiring path.
- Move all user code out of `internal/api/*` immediately - rejected because本阶段目标是迁移装配边界，不是大规模重写业务实现。

Trade-offs:

- Gain: 新业务模块有清晰模板，未来开发者看 `internal/modules/user` 就能理解装配方式。
- Cost: 迁移期会同时存在 legacy user implementation 和 module wrapper。
- Risk: 开发者仍误以为 legacy generated files 是入口。Mitigation: tests assert user only becomes active through explicit module inclusion.

### Decision: Separate Capability Name From Provider Name

Decision type: Build

Chose abstract capability requirements with provider selection over direct database capability requirements because:

- 业务模块应该表达“需要用户存储”，而不是表达“需要 MySQL”。
- 用户明确指出数据库选择应该方便切换，这是架构基础能力，不应散落在业务模块里。

Alternatives considered:

- Let `user` require `mysql` directly - rejected because it couples business behavior to one concrete database engine.
- Let each business module read config and instantiate its preferred database - rejected because it hides provider selection inside business code.

Trade-offs:

- Gain: 切换 provider 不需要改 user module declaration、HTTP routes 或 business use case。
- Cost: `internal/platform/module` 需要从 “capability name only” 扩展到 “capability name + provider identity + selection rule”。
- Risk: Provider selection 过度抽象。Mitigation: 只支持必需的字段：capability name、provider name、default marker、selected provider validation。

### Decision: Provide A Default User Storage Provider

Decision type: Compose

Chose a default local user storage provider plus optional MySQL provider over requiring MySQL for the migrated user example because:

- 项目目标是让支持的组件开箱即用，默认启动不应因为示例模块缺外部数据库而失败。
- MySQL disabled 是当前默认配置；默认 provider 可以证明 module-first path 和 HTTP 行为，不依赖外部服务。

Alternatives considered:

- Require MySQL whenever `user` is included - rejected because it makes storage engine choice a business-module concern and hurts out-of-box experience.
- Implement every future database provider now - rejected because this change only needs architecture boundary and one concrete external adapter path.

Trade-offs:

- Gain: 默认装配可用；显式选择 MySQL 时仍可 fail fast。
- Cost: 需要一个本地 provider 实现现有 user repository contract。
- Risk: 本地 provider 被误认为生产存储。Mitigation: name and tests frame it as the default development provider; provider selection remains explicit in assembly report.

### Decision: Reuse Existing User Controller Behind Module Routes

Decision type: Adopt

Chose to reuse `service.UserController` over duplicating HTTP handler logic because:

- 现有 handler 已覆盖 binding、validation、path param parsing、response envelope 和 error propagation。
- 本阶段不能引入新的 user business behavior。

Alternatives considered:

- Reimplement handlers directly in `internal/modules/user` - rejected because it duplicates request/response behavior and raises regression risk.
- Keep `RegisterRouter` as the runtime route adapter - rejected because module entry points need explicit route descriptors before server startup.

Trade-offs:

- Gain: 保持现有 HTTP 行为，迁移集中在装配边界和 provider 选择。
- Cost: `internal/modules/user` 暂时依赖 legacy service/biz/data packages。
- Risk: `service.NewUserController` 当前返回 interface，不方便暴露 handler methods。Mitigation: change it to return `*UserController` while keeping `*UserController` implementing `RegisterRouter`.

### Decision: Keep OpenAPI Out Of The Assembly Path

Decision type: Build

Chose handwritten module descriptors over generator output because:

- 用户明确要求不要回到旧 README/AGENTS/OpenAPI 思维。
- Phase 1 platform assembly already proves descriptors are enough to expose routes.

Alternatives considered:

- Generate the `user` module from OpenAPI - rejected because it makes OpenAPI authoritative again.
- Delete all legacy generated user files now - rejected because the migration can prove the new boundary without a risky broad cleanup.

Trade-offs:

- Gain: The new architecture is visible in ordinary Go code.
- Cost: Some legacy files and comments may remain until later cleanup.
- Risk: Old files create visual noise. Mitigation: tests cover that legacy files do not activate user unless `internal/modules/user` is included.

## Module and Interface Design

| Unit | Responsibility | Public Interface | Dependencies | Failure Behavior |
|---|---|---|---|---|
| `internal/platform/module` provider selection | Match capability requirements to enabled providers, honoring explicit selections and defaults. | `Capability{Name, Provider, Status, Default}`, `CapabilitySelection{Capability, Provider}`, `WithCapabilitySelections(...)`. | 标准库。 | Missing capability, unavailable selected provider, or ambiguous default provider returns typed assembly errors. |
| `internal/platform/app` | Pass provider selections into module assembly and expose selected provider in report. | `Options.CapabilitySelections`, `App.Report()`. | `platform/module`, `platform/http`, `server`. | Provider selection errors block startup before routes are served. |
| `internal/modules/user.Module` | Own the user business module declaration. | `New(repo biz.UserRepository) Module`, `Descriptor() module.Descriptor`. | `platform/module`, `platform/http`, existing user use case. | Missing `user.storage` provider is reported by platform assembly. |
| `internal/modules/user.routes` | Convert user controller methods into HTTP route descriptors. | `Routes(owner string, useCase biz.UserUseCase) []platformhttp.Route`. | `service.UserController`, `net/http`, `platform/http`. | Invalid route descriptor causes platform HTTP registration failure before serving. |
| `internal/modules/user/storage/memory` | Default local provider for `user.storage`, allowing user routes to run without external database. | Repository constructor and provider module descriptor. | `biz.UserRepository`, `param`. | In-memory repository errors return through existing user behavior and platform error handler. |
| `internal/modules/user/storage/mysql` | Optional MySQL-backed provider for `user.storage`, reusing existing data repository. | Provider module descriptor and repository constructor from config. | existing `data.NewUserRepository`, `db.NewDataManager`, MySQL config. | Disabled or unavailable MySQL provider cannot satisfy selected `user.storage=mysql`. |
| `service.UserController` | Existing HTTP behavior for user requests. | `ListUsers`, `Create`, `GetByID`, `Update`, `Delete`; `RegisterRouter` remains compatible. | `biz.UserUseCase`, `resp`, `code`, `utils`. | Binding/validation errors return to platform error handler. |
| `internal/app.Modules` | Explicitly include capability modules, selected/default user storage provider, and user business module. | `Modules(cfg configs.Config) (ModulesResult, error)` where `ModulesResult` contains modules and capability selections. | capability modules, `internal/modules/user`, user storage providers. | Invalid provider selection returns startup error before platform assembly. |

Interface sketch:

```go
// Pseudocode only.
const UserStorageCapability = "user.storage"

type Capability struct {
    Name     string
    Provider string
    Status   CapabilityState
    Default  bool
}

type CapabilitySelection struct {
    Capability string
    Provider   string
}

type UserModule struct {
    useCase biz.UserUseCase
}

func (m UserModule) Descriptor() module.Descriptor {
    return module.Descriptor{
        Name: "user",
        Kind: module.BusinessModule,
        Requires: []module.CapabilityRef{{Name: UserStorageCapability}},
        EntryPoints: userHTTPEntryPoints(m.useCase),
    }
}
```

## Data Flow

Default provider:

```text
cmd/app.go
  -> configs.BootstrapE
  -> internal/app builds module list and capability selections
      -> user storage provider defaults to memory when no provider is selected
      -> user module receives repository from selected/default provider
  -> platform/app.Assemble
      -> validates user requires user.storage
      -> resolves user.storage to provider memory
      -> records provider in assembly report
      -> extracts user-owned HTTP routes
  -> server.NewEchoServer
      -> registers /api/users routes from module descriptors
```

Explicit provider switch:

```text
config selects user.storage=mysql
  -> internal/app includes mysql user-storage provider
  -> platform/module resolves user.storage to provider mysql
  -> user module descriptor and HTTP entry points are unchanged
  -> report.Requirement("user", "user.storage").Provider == "mysql"
```

Invalid provider:

```text
config selects user.storage=unknown
  -> platform/module cannot find enabled provider unknown for user.storage
  -> startup fails before accepting requests
  -> error identifies module user, capability user.storage, provider unknown
```

Custom assembly omission:

```text
custom []module.Module without user
  -> platform/app.Assemble
  -> report.HasActiveModule("user") == false
  -> no /api/users routes
```

## File Changes

- Modify: `internal/platform/module/types.go` - add provider identity/default marker to capability declarations and define capability selections; expected net change under 80 lines.
- Modify: `internal/platform/module/report.go` - preserve selected provider in requirement status and capability status.
- Modify: `internal/platform/module/assembler.go` - resolve requirements by selected provider, default provider, or single enabled provider; reject unavailable selected provider.
- Modify: `internal/platform/module/errors.go` - add typed provider selection errors that identify module, capability, and selected provider.
- Modify: `internal/platform/module/assembler_test.go` - cover default provider selection, explicit provider switch, invalid selected provider, and unchanged requirement behavior.
- Modify: `internal/platform/app/app.go` - pass capability selections into module assembly.
- Create: `internal/modules/user/module.go` - user module descriptor, name constants, `user.storage` requirement; expected under 140 lines.
- Create: `internal/modules/user/routes.go` - HTTP route descriptors mapping to existing user controller methods; expected under 140 lines.
- Create: `internal/modules/user/storage/memory/module.go` - default provider module for `user.storage=memory`; expected under 180 lines.
- Create: `internal/modules/user/storage/memory/repository.go` - in-memory implementation of existing `biz.UserRepository`; expected under 220 lines, split if it grows beyond that.
- Create: `internal/modules/user/storage/mysql/module.go` - optional provider module for `user.storage=mysql`, adapting existing data repository; expected under 180 lines.
- Create: `internal/modules/user/module_test.go` - verifies module kind, `user.storage` requirement, route descriptors, and no concrete database name in module declaration.
- Create: `internal/modules/user/storage/memory/repository_test.go` - verifies the default provider supports existing user repository operations without external database.
- Create: `internal/modules/user/http_test.go` - assembles a test app with selected/default storage provider, invokes user routes through Echo, and checks envelopes/error handling.
- Modify: `internal/app/modules.go` - include provider modules and user module explicitly; return `ModulesResult` with module list and capability selections.
- Create: `internal/app/modules_test.go` - verifies default provider selection, explicit provider switch, and user module inclusion.
- Modify: `cmd/app.go` - handle the new app module builder return shape before platform assembly.
- Modify: `internal/api/service/user.go` - make `NewUserController` return `*UserController` while preserving `RegisterRouter` compatibility.
- Modify: `internal/api/service/user_test.go` - keep legacy handler behavior tests passing after constructor signature change if needed.
- Modify: `internal/platform/app/app_test.go` - add regression tests for custom assembly omitting user and for generated files not activating user.
- Modify: `internal/configs/config.go` and `configs/config.toml` - add `Capabilities.Providers` config loaded from `[capabilities.providers]`, with default behavior for omitted `user.storage`.
- Do not modify: `doc/example_openapi.yaml` for this change; it remains non-authoritative and unread by assembly.
- Do not recreate: `README.md` or `AGENTS.md`.

## Test Strategy

- Unit tests:
  - `internal/platform/module` verifies selected provider, default provider, missing provider, unavailable selected provider, and report provider identity.
  - `internal/modules/user` verifies descriptor name, kind, abstract `user.storage` requirement, route count, route owner, HTTP methods, and paths.
  - `internal/modules/user/storage/memory` verifies CRUD-like repository behavior for existing user operations without external services.
  - `internal/app` verifies default user storage provider, explicit provider switch, and user module inclusion.
- Integration tests:
  - Assemble user with default storage provider, build the server, invoke `GET /api/users`, and assert existing success envelope.
  - Assemble user with a different enabled storage provider selection in tests and assert user routes remain unchanged while the report provider changes.
  - Invoke invalid user input through the server and assert the platform error handler produces the existing error envelope.
- Error-path tests:
  - Assemble user without any enabled `user.storage` provider and assert a missing capability error identifies `user` and `user.storage`.
  - Select an unavailable provider and assert the error identifies `user`, `user.storage`, and the selected provider.
  - Assemble user with unsupported HTTP entry point adapters and assert `UnsupportedEntryPointError` identifies `user` and `http`.
- Regression tests:
  - Existing `internal/api/service`, `internal/api/biz`, and `internal/api/data` tests continue to pass.
  - Existing platform/module and platform/http tests continue to prove generic assembly behavior.
- Full verification:
  - Run `go test ./...` after implementation.

## Risks and Mitigations

| Risk | Impact | Mitigation |
|---|---|---|
| Provider selection becomes too generic too early. | Business developers face a new mini-framework. | Keep provider selection fields minimal and test only capability/provider/default behavior needed by user storage. |
| Default in-memory provider is mistaken for production storage. | Users may deploy with non-persistent data unexpectedly. | Report selected provider clearly and keep provider selection config explicit for non-default deployments. |
| MySQL provider and existing MySQL capability drift. | The architecture has two ways to talk about MySQL. | Treat `mysql` as a low-level infrastructure capability and `user.storage=mysql` as a business storage provider adapter; keep names distinct in reports. |
| Legacy generated files still exist. | Developers may still look at old files first. | Make `internal/modules/user` the only activation path and cover this with tests. |
| User module depends on legacy service/biz/data packages. | Module-first purity is not complete. | Keep dependency one-way and local to user module; later change can move implementation fully under `internal/modules/user`. |
| Constructor or route changes break old service tests. | Existing behavior regression. | Keep `RegisterRouter` interface compatibility and run existing service tests. |

## Self-Review

- [x] 每个 Requirement 都已映射到设计方案。
- [x] 所有会变更的文件都列在 File Changes 中，且路径精确。
- [x] 新文件遵循现有 Go 包命名和 `internal/modules/<name>` 方向。
- [x] 每个模块职责单一、接口清晰、依赖可见。
- [x] 没有计划新增 300+ 行 god file；route、descriptor、provider 和 repository 已拆分。
- [x] 每个 Architecture Decision 都包含备选方案和取舍。
- [x] 测试策略覆盖 happy path、默认 provider、provider 切换、缺 capability、invalid provider、unsupported adapter、invalid input 和 OpenAPI 非权威。
