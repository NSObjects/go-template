# Design: Module-First Architecture

## Technical Approach

引入一个很小的模块装配内核：业务模块声明能力需求和入口，能力模块声明可提供的能力与健康状态，平台模块负责把入口暴露到运行时并托管共享行为。第一阶段采用增量落地：先建立模块内核、应用装配、HTTP 入口适配和能力状态模型，再逐步把现有业务迁移到 `internal/modules/*`。

现有配置、日志、错误、响应、中间件和 HTTP server 代码作为平台资产复用，但不再让业务模块直接参与中心化 wiring。OpenAPI 和生成器不参与新架构的装配路径。

## Requirements Mapping

| Requirement | Technical Approach | Files / Modules |
|---|---|---|
| Business Modules Are the Source of Truth | 新增业务模块目录约定和模块声明接口；业务模块自己声明 entry points 和 required capabilities。 | `internal/platform/module/*`, `internal/modules/*`, `internal/app/modules.go` |
| Platform Owns Shared Runtime Behavior | 把配置、日志、错误、响应、HTTP 中间件、生命周期、健康聚合放入 platform/app 装配流程。 | `internal/platform/app/*`, `internal/platform/http/*`, `internal/server/echo_server.go` |
| Capabilities Are Declared by Business Modules | 模块描述中包含 required capability refs；装配阶段统一校验并生成 assembly report。 | `internal/platform/module/types.go`, `internal/platform/module/assembler.go`, `internal/platform/module/report.go` |
| Capability Modules Are Optional and Observable | 能力模块报告 enabled、disabled、healthy、unavailable 状态；禁用且无人依赖时允许启动。 | `internal/platform/module/capability.go`, `internal/capabilities/*/module.go` |
| Entry Points Are Declared by Business Modules | entry point 使用 type + owner + payload 描述；platform adapter 负责暴露对应入口。 | `internal/platform/module/entrypoint.go`, `internal/platform/http/route.go`, `internal/platform/http/echo_adapter.go` |
| Normal Business Work Avoids Central Wiring Changes | 应用只维护一个模块 inclusion list；业务模块内部完成自身声明和适配。 | `internal/app/modules.go`, `internal/modules/*/module.go` |
| OpenAPI Is Not the Architecture Source of Truth | 新 platform/module/app 包不读取 OpenAPI；OpenAPI 只能作为未来可选 adapter 的输入。 | `internal/platform/*`, `internal/app/*` |
| Code Generation Is Optional Assistance | 生成产物只有在被包装为业务模块并加入 module list 后才参与装配。 | `internal/app/modules.go`, `internal/modules/*` |

## Existing Patterns Observed

- Naming: 现有 Go 包使用短小单数名，例如 `configs`、`server`、`resp`、`code`；业务示例使用 `user`。
- File structure: 当前主要是 layer-based：`internal/api/service`、`internal/api/biz`、`internal/api/data`。新架构需要转向 module-based：`internal/modules/<name>`。
- State/data flow: 当前由 `cmd/app.go` 通过 `samber/do` 中心化注册配置、日志、db、data、biz、service、server。新设计保留 DI 可用性，但把中心化注册收敛到 platform app。
- Error handling: 现有 `internal/code` 和 `internal/server/middlewares/error.go` 已经提供统一错误表达和 HTTP 错误转换，应复用。
- Response handling: 现有 `internal/resp` 提供统一响应 envelope，应作为 HTTP 平台行为复用。
- Testing: 现有测试使用 Go 标准测试和少量 `testify`；新模块核心优先使用 table-driven tests。
- Configuration: 现有 `internal/configs` 已支持 file/source/store/hot reload，应复用为 platform 配置能力。

## Architecture Decisions

### Decision: Build A Small Module Kernel

Decision type: Build

Chose a project-owned module kernel over extending the current service/biz/data registration because:

- 当前注册方式把业务、数据、路由分散到多个中心文件，业务模块不是一等开发单元。
- 规格要求 OpenAPI 和生成器不能成为架构源头，现有分层注册无法表达 capability 和 entry point 声明。

Alternatives considered:

- Keep service/biz/data as the main architecture — rejected because it preserves layer-first thinking.
- Use reflection or filesystem discovery for modules — rejected because Go 项目中显式 inclusion list 更可审查、更稳定。

Trade-offs:

- Gain: 业务模块、能力模块、平台模块拥有统一装配语言。
- Cost: 第一阶段会同时存在 legacy layer code 和新 module code。
- Risk: 新内核过大。Mitigation: `internal/platform/module` 只负责描述、校验和报告，不承载 HTTP、DB 或业务实现。

### Decision: Use Explicit Application Module List

Decision type: Compose

Chose an explicit module list over automatic discovery because:

- Go 没有稳定、透明的包级自动发现机制。
- 显式列表让“哪些模块被启用”成为可审查状态，也满足规格里的 module inclusion 行为。

Alternatives considered:

- Blank imports for self-registration — rejected because side effects make装配顺序和依赖更难追踪。
- Config-only dynamic loading — rejected for first phase because it hides compile-time dependencies and raises operational complexity.

Trade-offs:

- Gain: 新业务只需要加入 module list，模块内部自带声明。
- Cost: 新模块仍需一处 inclusion 操作。
- Risk: 用户可能把 inclusion list 当成 wiring 文件继续堆逻辑。Mitigation: `internal/app/modules.go` 只允许返回模块列表，不包含能力解析或路由注册逻辑。

### Decision: Keep Platform HTTP As Adapter

Decision type: Extend

Chose to extend the existing Echo server through a platform HTTP adapter because:

- 现有中间件、验证器、错误处理和响应 envelope 可以直接复用。
- 规格要求共享运行时行为由平台负责，HTTP adapter 正好集中这些行为。

Alternatives considered:

- Replace Echo immediately — rejected because没有行为收益，迁移成本高。
- Let business modules register Echo routes directly — rejected because会把平台 HTTP 细节泄漏到业务模块。

Trade-offs:

- Gain: 保留现有稳定能力，同时让业务模块只声明 HTTP entry points。
- Cost: 第一阶段需要一层 adapter 把模块 entry point 转成 Echo route。
- Risk: adapter 变成大文件。Mitigation: route type、adapter、tests 分文件放在 `internal/platform/http`。

### Decision: Model Capabilities As Observable Runtime Providers

Decision type: Compose

Chose capability providers with explicit status over exposing raw clients directly because:

- 规格要求 capability 可观察为 enabled、disabled、healthy、unavailable。
- 业务模块应声明需要什么，而不是检查底层连接是否存在。

Alternatives considered:

- Keep `DataManager` as the only data capability — rejected as long-term shape because它是宽接口，容易让业务绕过模块声明。
- Split every existing infra package before module kernel — rejected for first phase because会扩大改动面。

Trade-offs:

- Gain: capability status、健康检查和缺失依赖可以统一报告。
- Cost: 现有 MySQL/Redis/Mongo/Kafka 需要逐步从 `DataManager` 迁移到 capability modules。
- Risk: 迁移期重复能力入口。Mitigation: first phase 允许 compatibility adapter，但新业务只使用 capability declarations。

## Module and Interface Design

| Unit | Responsibility | Public Interface | Dependencies | Failure Behavior |
|---|---|---|---|---|
| `internal/platform/module` | 描述模块、能力、入口并执行装配校验。 | Module descriptor, capability refs, entry point descriptors, assembly report. | 标准库。 | 缺 required capability、unsupported entry point 时返回 assembly error；无 entry point 时产生 warning。 |
| `internal/platform/app` | 应用装配入口，组合 config、logger、platform modules、capability modules 和 business modules。 | App options, assemble result, run lifecycle. | `configs`, `log`, `server`, `platform/module`, `platform/http`。 | 装配失败时应用不启动并返回带模块名的错误。 |
| `internal/platform/http` | 将 HTTP entry points 暴露到现有 HTTP server，并应用共享 HTTP 行为。 | HTTP route descriptor, adapter registration. | `server`, `resp`, `code`, middlewares。 | 不支持的 route declaration 返回 adapter error；handler error 走统一错误处理。 |
| `internal/capabilities/mysql` | 提供 MySQL capability，报告 enabled/disabled/healthy/unavailable。 | Capability module descriptor and provider. | 现有 MySQL config/constructor。 | enabled 但连接失败时标记 unavailable 并阻止启动。 |
| `internal/capabilities/redis` | 提供 Redis capability，报告 enabled/disabled/healthy/unavailable。 | Capability module descriptor and provider. | 现有 Redis config/constructor。 | enabled 但连接失败时标记 unavailable 并阻止启动。 |
| `internal/capabilities/mongodb` | 提供 MongoDB capability，报告 enabled/disabled/healthy/unavailable。 | Capability module descriptor and provider. | 现有 Mongo config/constructor。 | enabled 但连接失败时标记 unavailable 并阻止启动。 |
| `internal/capabilities/kafka` | 提供 Kafka capability，报告 enabled/disabled/healthy/unavailable。 | Capability module descriptor and provider. | 现有 Kafka config/constructor。 | enabled 但 producer 创建失败时标记 unavailable 并阻止启动。 |
| `internal/app/modules.go` | 声明当前应用启用哪些业务模块和 capability modules。 | A module list. | business modules, capability modules。 | 未列入的业务模块不参与装配，entry points 不暴露。 |
| `internal/modules/<name>` | 承载一个业务模块的行为、entry point adapter 和 capability declarations。 | Module declaration plus module-local business interfaces. | platform module contracts and selected capabilities. | required capability 缺失时由装配阶段失败；业务错误通过平台错误处理返回。 |

Interface sketch:

```go
// Pseudocode only.
type Module interface {
    Descriptor() Descriptor
    Register(Context) error
}

type Descriptor struct {
    Name        string
    Kind        ModuleKind
    Requires    []CapabilityRef
    Provides    []Capability
    EntryPoints []EntryPoint
}

type EntryPoint struct {
    Owner string
    Type  string
    Name  string
    Value any
}
```

## Data Flow

```text
cmd
  -> platform/app loads config and logger
  -> internal/app/modules.go returns included modules
  -> platform/module assembles descriptors
      -> validates required capabilities
      -> validates entry point adapters
      -> creates assembly report
  -> platform/http registers HTTP entry points with server
  -> platform/app starts lifecycle and health checks
  -> incoming trigger reaches platform adapter
  -> platform shared behavior wraps the trigger
  -> business module behavior runs
```

Capability status flow:

```text
config
  -> capability module decides enabled or disabled
  -> enabled capability initializes provider
  -> provider reports healthy or unavailable
  -> assembly matches business module requirements
  -> report exposes module requirements and capability status
```

## File Changes

- Create: `internal/platform/module/types.go` — module, descriptor, capability ref, capability status, entry point types; expected under 160 lines.
- Create: `internal/platform/module/assembler.go` — descriptor collection, required capability validation, unsupported entry point validation; expected under 200 lines.
- Create: `internal/platform/module/report.go` — assembly report, warnings, active modules, capability statuses, entry point listing; expected under 140 lines.
- Create: `internal/platform/module/errors.go` — typed assembly errors with module name, capability name, entry point type; expected under 120 lines.
- Create: `internal/platform/module/assembler_test.go` — table-driven coverage for missing capability, disabled optional capability, no entry point warning, unsupported entry point.
- Create: `internal/platform/http/route.go` — HTTP route entry point descriptor independent from business use cases; expected under 120 lines.
- Create: `internal/platform/http/echo_adapter.go` — converts HTTP route descriptors into existing server routes; expected under 180 lines.
- Create: `internal/platform/http/echo_adapter_test.go` — verifies route registration and unsupported entry point rejection.
- Create: `internal/platform/app/app.go` — application assembly orchestration and lifecycle handoff; expected under 200 lines.
- Create: `internal/platform/app/app_test.go` — verifies assembly failure prevents startup and report contains module/capability details.
- Create: `internal/app/modules.go` — application module inclusion list; expected under 80 lines.
- Create: `internal/capabilities/mysql/module.go` — MySQL capability provider wrapper around existing config and constructor; expected under 180 lines.
- Create: `internal/capabilities/redis/module.go` — Redis capability provider wrapper; expected under 160 lines.
- Create: `internal/capabilities/mongodb/module.go` — MongoDB capability provider wrapper; expected under 160 lines.
- Create: `internal/capabilities/kafka/module.go` — Kafka capability provider wrapper; expected under 180 lines.
- Create: `internal/capabilities/*/module_test.go` — disabled/enabled/misconfigured capability status tests.
- Modify: `cmd/app.go` — delegate application assembly and run to `internal/platform/app`, keeping command-line behavior intact.
- Modify: `internal/server/echo_server.go` — accept platform HTTP route registrations instead of legacy service router list, while retaining middleware, validator, error handling, and system routes.
- Modify: `internal/server/echo_server_test.go` — update tests from legacy router registration to platform HTTP route registration.
- Modify: `internal/api/service/service.go`, `internal/api/biz/biz.go`, `internal/api/data/data.go` — mark legacy registration path as transitional or remove from app startup once replacement app assembly is active.
- Do not modify: `doc/example_openapi.yaml` in this change; it is ignored by the new architecture unless a future optional adapter includes it.
- Do not recreate: `README.md` or `AGENTS.md`.

## Test Strategy

- Unit tests:
  - `internal/platform/module` validates capability matching, disabled optional capability status, missing required capability errors, unsupported entry point errors, and no-entry-point warnings.
  - `internal/capabilities/*` validates disabled status without external connection attempts and misconfigured enabled status failures.
- Integration tests:
  - `internal/platform/http` registers a fake HTTP entry point and verifies the route appears in the server route list.
  - `internal/platform/app` assembles fake business and capability modules and verifies startup is blocked on missing requirements.
- Regression tests:
  - Existing config loading, error response, response envelope, and server middleware tests stay active.
  - Existing application command tests are updated to assert platform app assembly rather than direct service/biz/data wiring.
- Manual verification:
  - Start with all external capabilities disabled and no module requiring them; application starts and reports disabled statuses.
  - Add a fake business module requiring a disabled capability; application fails before accepting requests and names the missing capability and module.

## Risks and Mitigations

| Risk | Impact | Mitigation |
|---|---|---|
| New module kernel becomes another framework inside the project. | Business developers must learn too much before writing business behavior. | Keep kernel descriptors small; move HTTP, lifecycle, and capabilities into separate packages. |
| Legacy layer code and new module code coexist too long. | Developers see two architectural styles. | Mark legacy registration as transitional and migrate one business module after foundation lands. |
| Capability wrappers still depend on old `DataManager` internals. | New architecture leaks old data access shape. | First phase may reuse constructors, but new business modules depend on capability declarations, not `DataManager`. |
| Explicit module list becomes a dumping ground. | Central wiring returns under a new name. | Restrict `internal/app/modules.go` to inclusion only; no route, capability, or business construction logic. |
| HTTP adapter leaks Echo details into business behavior. | Business modules become transport-coupled. | Business use cases stay framework-free; only module-local HTTP adapter can touch platform HTTP types. |

## Migration Notes

- First implementation should build the platform/module/capability foundation and keep existing business behavior stable.
- After foundation is tested, migrate one existing business area into `internal/modules/<name>` as the first real business module proof.
- OpenAPI files may remain in the repository as documentation or future optional adapter input, but they must not be read by the platform app assembly.
- Generated or legacy files only participate in runtime behavior when a module explicitly includes them.

## Self-Review

- [x] 每个 Requirement 都已映射到设计方案。
- [x] 所有会变更的文件都列在 File Changes 中，且路径精确。
- [x] 新文件遵循 Go 项目现有短包名和 `internal/*` 组织方式。
- [x] 每个模块职责单一、接口清晰、依赖可见。
- [x] 没有计划新增 300+ 行 god file；预计超过 160 行的文件已有拆分。
- [x] 每个 Architecture Decision 都包含备选方案和取舍。
- [x] 测试策略覆盖 happy path、边界和错误场景。
