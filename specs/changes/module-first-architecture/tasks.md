# Tasks

## Assumptions

- 第一阶段保留现有 Echo、配置、错误、响应和中间件能力，只改变它们被平台装配和业务模块使用的方式。
- 第一阶段不迁移真实 `user` 业务模块；使用测试模块证明 module-first 装配行为。
- 外部能力在 disabled 状态下不得尝试真实网络连接；enabled 但配置缺失时通过装配或初始化错误暴露。
- `README.md` 和 `AGENTS.md` 保持删除状态，本计划不恢复它们。

## Open Planning Blockers

- 无

## Spec Coverage Summary

| Spec Scenario | Requirement | Covered by Task(s) | Notes |
| --- | --- | --- | --- |
| `specs/module-first-architecture/spec.md` → Scenario "Developer adds a complete business module" | "Business Modules Are the Source of Truth" | Task 1.1, Task 5.1 | 完整模块通过 inclusion list 参与装配并暴露入口。 |
| `specs/module-first-architecture/spec.md` → Scenario "Business module is incomplete" | "Business Modules Are the Source of Truth" | Task 1.2 | 无 entry point 产生带模块名 warning。 |
| `specs/module-first-architecture/spec.md` → Scenario "Shared runtime behavior is available to a business module" | "Platform Owns Shared Runtime Behavior" | Task 2.1, Task 4.2 | HTTP adapter 复用现有 server 中间件、错误和响应行为。 |
| `specs/module-first-architecture/spec.md` → Scenario "Required platform behavior is unavailable" | "Platform Owns Shared Runtime Behavior" | Task 2.2, Task 4.1 | unsupported entry point 或缺平台 adapter 阻止启动。 |
| `specs/module-first-architecture/spec.md` → Scenario "Required capability is enabled" | "Capabilities Are Declared by Business Modules" | Task 1.1, Task 3.2 | required capability 被 enabled capability 满足并写入 report。 |
| `specs/module-first-architecture/spec.md` → Scenario "Required capability is missing" | "Capabilities Are Declared by Business Modules" | Task 1.2, Task 4.1 | missing capability 产生 assembly error 并阻止启动。 |
| `specs/module-first-architecture/spec.md` → Scenario "Optional capability is disabled" | "Capability Modules Are Optional and Observable" | Task 1.1, Task 3.1 | disabled 且无人依赖时允许装配，状态可观察。 |
| `specs/module-first-architecture/spec.md` → Scenario "Enabled capability is unavailable" | "Capability Modules Are Optional and Observable" | Task 3.2, Task 4.1 | enabled 但配置不可用时阻止启动并命名 capability。 |
| `specs/module-first-architecture/spec.md` → Scenario "Entry point declaration is accepted" | "Entry Points Are Declared by Business Modules" | Task 2.1 | HTTP entry point 注册后可从 server routes 观察。 |
| `specs/module-first-architecture/spec.md` → Scenario "Entry point declaration cannot be exposed" | "Entry Points Are Declared by Business Modules" | Task 2.2 | unsupported entry point type 产生 adapter error。 |
| `specs/module-first-architecture/spec.md` → Scenario "Business module is added through module inclusion" | "Normal Business Work Avoids Central Wiring Changes" | Task 4.2, Task 5.1 | inclusion list 只声明模块，不包含业务 wiring。 |
| `specs/module-first-architecture/spec.md` → Scenario "Module inclusion is omitted" | "Normal Business Work Avoids Central Wiring Changes" | Task 4.2 | 未包含模块不会出现在 active modules 或 routes。 |
| `specs/module-first-architecture/spec.md` → Scenario "Application assembles without OpenAPI input" | "OpenAPI Is Not the Architecture Source of Truth" | Task 4.2, Task 5.1 | app assembly 不读取 OpenAPI 文件。 |
| `specs/module-first-architecture/spec.md` → Scenario "OpenAPI input is present" | "OpenAPI Is Not the Architecture Source of Truth" | Task 5.2 | 回归测试确保 OpenAPI 文件存在也不影响 module declarations。 |
| `specs/module-first-architecture/spec.md` → Scenario "Developer creates business behavior manually" | "Code Generation Is Optional Assistance" | Task 4.2, Task 5.1 | 测试模块手写，不运行 generator。 |
| `specs/module-first-architecture/spec.md` → Scenario "Generated output exists" | "Code Generation Is Optional Assistance" | Task 5.2 | legacy/generated 文件不被 app assembly 自动纳入。 |

## 1. Platform Module Kernel

### Task 1.1: Assemble modules with satisfied capabilities and observable report

**Type:** behavior
**Covers specs:** `specs/module-first-architecture/spec.md` → Requirement "Capabilities Are Declared by Business Modules" → Scenario "Required capability is enabled"; `specs/module-first-architecture/spec.md` → Requirement "Capability Modules Are Optional and Observable" → Scenario "Optional capability is disabled"
**Depends on:** none
**Files:**
- `internal/platform/module/assembler_test.go` — add table-driven tests for satisfied required capability and disabled optional capability status.
- `internal/platform/module/types.go` — create module, descriptor, capability, capability ref, entry point, and status types.
- `internal/platform/module/report.go` — create assembly report with active modules, capability statuses, entry point listing, and warnings.
- `internal/platform/module/assembler.go` — create assembly responsibility that collects descriptors and matches requirements to provided capabilities.
**Test command:** `go test ./internal/platform/module`

**Acceptance criteria:**
- GIVEN a business module requires a capability and an enabled capability module provides it WHEN module assembly runs THEN the report marks the requirement satisfied and includes both modules as active.
- GIVEN an optional capability module is disabled and no business module requires it WHEN module assembly runs THEN assembly succeeds and the report exposes that capability as disabled.

- [x] **Step 1: Write failing test**
  File: `internal/platform/module/assembler_test.go`
  Test name: `TestAssembleReportsSatisfiedCapabilities`
  Expected RED result: FAIL because `internal/platform/module` does not exist.

- [x] **Step 2: Verify RED**
  Run: `go test ./internal/platform/module`
  Expected: FAIL for missing package or missing assembly behavior, not because of unrelated packages.

- [x] **Step 3: Describe minimal production responsibility**
  File: `internal/platform/module/types.go`, `internal/platform/module/report.go`, `internal/platform/module/assembler.go`
  Responsibility: Provide a small module descriptor model and an assembly report that can show active modules, satisfied required capabilities, disabled optional capabilities, and declared entry points.

- [x] **Step 4: Verify GREEN**
  Run: `go test ./internal/platform/module`
  Expected: PASS

### Task 1.2: Reject missing capabilities and warn on modules with no entry points

**Type:** error-path
**Covers specs:** `specs/module-first-architecture/spec.md` → Requirement "Business Modules Are the Source of Truth" → Scenario "Business module is incomplete"; `specs/module-first-architecture/spec.md` → Requirement "Capabilities Are Declared by Business Modules" → Scenario "Required capability is missing"
**Depends on:** Task 1.1
**Files:**
- `internal/platform/module/assembler_test.go` — add tests for missing required capability and no-entry-point warning.
- `internal/platform/module/errors.go` — create typed assembly errors that identify module name and missing capability.
- `internal/platform/module/assembler.go` — add validation for missing required capabilities and no-entry-point report warnings.
- `internal/platform/module/report.go` — include warning records that identify affected module names.
**Test command:** `go test ./internal/platform/module`

**Acceptance criteria:**
- GIVEN a business module declares a required capability that no enabled capability module provides WHEN module assembly runs THEN assembly fails and the error identifies the missing capability and business module.
- GIVEN a business module has no declared entry points WHEN module assembly runs THEN assembly succeeds with a warning that identifies the affected module.

- [x] **Step 1: Write failing test**
  File: `internal/platform/module/assembler_test.go`
  Test name: `TestAssembleRejectsMissingRequiredCapability`
  Expected RED result: FAIL because missing capability validation is not implemented.

- [x] **Step 2: Verify RED**
  Run: `go test ./internal/platform/module`
  Expected: FAIL for missing capability assertion, not because of compile/setup errors.

- [x] **Step 3: Describe minimal production responsibility**
  File: `internal/platform/module/errors.go`, `internal/platform/module/assembler.go`, `internal/platform/module/report.go`
  Responsibility: Stop assembly when required capabilities are absent and record warnings when included business modules expose no entry points.

- [x] **Step 4: Verify GREEN**
  Run: `go test ./internal/platform/module`
  Expected: PASS

## 2. Platform HTTP Adapter

### Task 2.1: Expose declared HTTP entry points through the existing server

**Type:** integration
**Covers specs:** `specs/module-first-architecture/spec.md` → Requirement "Entry Points Are Declared by Business Modules" → Scenario "Entry point declaration is accepted"; `specs/module-first-architecture/spec.md` → Requirement "Platform Owns Shared Runtime Behavior" → Scenario "Shared runtime behavior is available to a business module"
**Depends on:** Task 1.1
**Files:**
- `internal/platform/http/echo_adapter_test.go` — add route registration test using a fake HTTP entry point.
- `internal/platform/http/route.go` — create HTTP route entry point descriptor used by business modules.
- `internal/platform/http/echo_adapter.go` — register HTTP route descriptors into the existing Echo server.
- `internal/server/echo_server.go` — accept platform HTTP routes while retaining validator, middleware, error handler, and system routes.
- `internal/server/echo_server_test.go` — update server route tests to assert platform HTTP route registration.
**Test command:** `go test ./internal/platform/http ./internal/server`

**Acceptance criteria:**
- GIVEN a business module declares an HTTP entry point WHEN the server is assembled THEN the route appears in the server route list under the configured API group.
- GIVEN the route handler returns an error WHEN the route is invoked THEN existing server error handling remains the response path.

- [x] **Step 1: Write failing test**
  File: `internal/platform/http/echo_adapter_test.go`
  Test name: `TestAdapterRegistersHTTPEntryPoint`
  Expected RED result: FAIL because `internal/platform/http` does not exist.

- [x] **Step 2: Verify RED**
  Run: `go test ./internal/platform/http ./internal/server`
  Expected: FAIL for missing package or missing adapter behavior, not because of unrelated packages.

- [x] **Step 3: Describe minimal production responsibility**
  File: `internal/platform/http/route.go`, `internal/platform/http/echo_adapter.go`, `internal/server/echo_server.go`
  Responsibility: Represent HTTP entry points separately from business behavior and register them through the existing server so shared HTTP runtime behavior remains centralized.

- [x] **Step 4: Verify GREEN**
  Run: `go test ./internal/platform/http ./internal/server`
  Expected: PASS

### Task 2.2: Reject unsupported entry point types before startup

**Type:** error-path
**Covers specs:** `specs/module-first-architecture/spec.md` → Requirement "Entry Points Are Declared by Business Modules" → Scenario "Entry point declaration cannot be exposed"; `specs/module-first-architecture/spec.md` → Requirement "Platform Owns Shared Runtime Behavior" → Scenario "Required platform behavior is unavailable"
**Depends on:** Task 2.1
**Files:**
- `internal/platform/module/assembler_test.go` — add unsupported entry point type test using registered platform entry point types.
- `internal/platform/module/types.go` — add entry point adapter capability description.
- `internal/platform/module/assembler.go` — validate declared entry point types against available platform adapters.
- `internal/platform/module/errors.go` — include unsupported entry point error with module name and entry point type.
**Test command:** `go test ./internal/platform/module ./internal/platform/http`

**Acceptance criteria:**
- GIVEN a business module declares an entry point type without a matching platform adapter WHEN module assembly runs THEN assembly fails before startup.
- GIVEN assembly fails for unsupported entry point type WHEN the error is inspected THEN it identifies both the entry point type and the affected business module.

- [x] **Step 1: Write failing test**
  File: `internal/platform/module/assembler_test.go`
  Test name: `TestAssembleRejectsUnsupportedEntryPoint`
  Expected RED result: FAIL because entry point adapter validation is absent.

- [x] **Step 2: Verify RED**
  Run: `go test ./internal/platform/module ./internal/platform/http`
  Expected: FAIL for unsupported entry point assertion, not because of compile/setup errors.

- [x] **Step 3: Describe minimal production responsibility**
  File: `internal/platform/module/types.go`, `internal/platform/module/assembler.go`, `internal/platform/module/errors.go`
  Responsibility: Validate all declared entry point types against platform-provided entry point adapters and return a startup-blocking assembly error when no adapter exists.

- [x] **Step 4: Verify GREEN**
  Run: `go test ./internal/platform/module ./internal/platform/http`
  Expected: PASS

## 3. Capability Modules

### Task 3.1: Report disabled external capabilities without connecting

**Type:** behavior
**Covers specs:** `specs/module-first-architecture/spec.md` → Requirement "Capability Modules Are Optional and Observable" → Scenario "Optional capability is disabled"
**Depends on:** Task 1.1
**Files:**
- `internal/capabilities/mysql/module_test.go` — add disabled status test.
- `internal/capabilities/redis/module_test.go` — add disabled status test.
- `internal/capabilities/mongodb/module_test.go` — add disabled status test.
- `internal/capabilities/kafka/module_test.go` — add disabled status test.
- `internal/capabilities/mysql/module.go` — expose disabled MySQL capability status from config.
- `internal/capabilities/redis/module.go` — expose disabled Redis capability status from config.
- `internal/capabilities/mongodb/module.go` — expose disabled MongoDB capability status from config.
- `internal/capabilities/kafka/module.go` — expose disabled Kafka capability status from config.
**Test command:** `go test ./internal/capabilities/mysql ./internal/capabilities/redis ./internal/capabilities/mongodb ./internal/capabilities/kafka`

**Acceptance criteria:**
- GIVEN each external capability is disabled in config WHEN its capability module is described THEN it reports disabled status.
- GIVEN each external capability is disabled in config WHEN its capability module is described THEN no external connection is attempted.

- [x] **Step 1: Write failing test**
  File: `internal/capabilities/mysql/module_test.go`
  Test name: `TestModuleReportsDisabledWhenConfigDisabled`
  Expected RED result: FAIL because capability module packages do not exist.

- [x] **Step 2: Verify RED**
  Run: `go test ./internal/capabilities/mysql ./internal/capabilities/redis ./internal/capabilities/mongodb ./internal/capabilities/kafka`
  Expected: FAIL for missing packages or missing disabled status behavior, not because of unrelated packages.

- [x] **Step 3: Describe minimal production responsibility**
  File: `internal/capabilities/mysql/module.go`, `internal/capabilities/redis/module.go`, `internal/capabilities/mongodb/module.go`, `internal/capabilities/kafka/module.go`
  Responsibility: Convert existing configuration into observable disabled capability descriptors without initializing external clients.

- [x] **Step 4: Verify GREEN**
  Run: `go test ./internal/capabilities/mysql ./internal/capabilities/redis ./internal/capabilities/mongodb ./internal/capabilities/kafka`
  Expected: PASS

### Task 3.2: Fail enabled external capabilities with invalid configuration

**Type:** error-path
**Covers specs:** `specs/module-first-architecture/spec.md` → Requirement "Capability Modules Are Optional and Observable" → Scenario "Enabled capability is unavailable"; `specs/module-first-architecture/spec.md` → Requirement "Capabilities Are Declared by Business Modules" → Scenario "Required capability is enabled"
**Depends on:** Task 3.1
**Files:**
- `internal/capabilities/mysql/module_test.go` — add enabled invalid configuration test and enabled valid descriptor test.
- `internal/capabilities/redis/module_test.go` — add enabled invalid configuration test and enabled valid descriptor test.
- `internal/capabilities/mongodb/module_test.go` — add enabled invalid configuration test and enabled valid descriptor test.
- `internal/capabilities/kafka/module_test.go` — add enabled invalid configuration test and enabled valid descriptor test.
- `internal/capabilities/mysql/module.go` — validate enabled MySQL config and expose unavailable status or provider descriptor.
- `internal/capabilities/redis/module.go` — validate enabled Redis config and expose unavailable status or provider descriptor.
- `internal/capabilities/mongodb/module.go` — validate enabled MongoDB config and expose unavailable status or provider descriptor.
- `internal/capabilities/kafka/module.go` — validate enabled Kafka config and expose unavailable status or provider descriptor.
**Test command:** `go test ./internal/capabilities/mysql ./internal/capabilities/redis ./internal/capabilities/mongodb ./internal/capabilities/kafka`

**Acceptance criteria:**
- GIVEN an external capability is enabled with missing required configuration WHEN its capability module is described THEN it reports unavailable or returns a startup-blocking capability error naming the capability.
- GIVEN an external capability is enabled with minimally valid configuration WHEN its capability module is described THEN it provides an enabled capability descriptor that can satisfy a business requirement.

- [x] **Step 1: Write failing test**
  File: `internal/capabilities/mysql/module_test.go`
  Test name: `TestModuleRejectsEnabledInvalidConfig`
  Expected RED result: FAIL because enabled config validation is absent.

- [x] **Step 2: Verify RED**
  Run: `go test ./internal/capabilities/mysql ./internal/capabilities/redis ./internal/capabilities/mongodb ./internal/capabilities/kafka`
  Expected: FAIL for invalid configuration assertion, not because of compile/setup errors.

- [x] **Step 3: Describe minimal production responsibility**
  File: `internal/capabilities/mysql/module.go`, `internal/capabilities/redis/module.go`, `internal/capabilities/mongodb/module.go`, `internal/capabilities/kafka/module.go`
  Responsibility: Validate enabled external capability configuration and expose enabled or unavailable capability state in a way the module assembler can consume.

- [x] **Step 4: Verify GREEN**
  Run: `go test ./internal/capabilities/mysql ./internal/capabilities/redis ./internal/capabilities/mongodb ./internal/capabilities/kafka`
  Expected: PASS

## 4. Application Assembly

### Task 4.1: Block application startup on assembly failures

**Type:** integration
**Covers specs:** `specs/module-first-architecture/spec.md` → Requirement "Capabilities Are Declared by Business Modules" → Scenario "Required capability is missing"; `specs/module-first-architecture/spec.md` → Requirement "Capability Modules Are Optional and Observable" → Scenario "Enabled capability is unavailable"
**Depends on:** Task 1.2, Task 2.2, Task 3.2
**Files:**
- `internal/platform/app/app_test.go` — add test proving missing capability and unavailable capability block startup with module or capability names in the error.
- `internal/platform/app/app.go` — create application assembly orchestration and startup gate.
- `cmd/app.go` — delegate runtime construction to platform app while preserving config-file error behavior.
- `cmd/app_test.go` — update command-level regression around config errors and startup gate.
**Test command:** `go test ./internal/platform/app ./cmd`

**Acceptance criteria:**
- GIVEN application assembly has a missing required capability WHEN startup is requested THEN startup fails before accepting entry point triggers and names the affected module and capability.
- GIVEN a capability is enabled but unavailable WHEN startup is requested THEN startup fails before accepting entry point triggers and names the unavailable capability.

- [x] **Step 1: Write failing test**
  File: `internal/platform/app/app_test.go`
  Test name: `TestAppStartupFailsWhenAssemblyFails`
  Expected RED result: FAIL because platform app package does not exist.

- [x] **Step 2: Verify RED**
  Run: `go test ./internal/platform/app ./cmd`
  Expected: FAIL for missing package or missing startup gate behavior, not because of unrelated packages.

- [x] **Step 3: Describe minimal production responsibility**
  File: `internal/platform/app/app.go`, `cmd/app.go`
  Responsibility: Load configuration, collect included modules, run module assembly, and prevent server startup when assembly reports blocking errors.

- [x] **Step 4: Verify GREEN**
  Run: `go test ./internal/platform/app ./cmd`
  Expected: PASS

### Task 4.2: Compose included modules without central business wiring

**Type:** behavior
**Covers specs:** `specs/module-first-architecture/spec.md` → Requirement "Normal Business Work Avoids Central Wiring Changes" → Scenario "Business module is added through module inclusion"; `specs/module-first-architecture/spec.md` → Requirement "Normal Business Work Avoids Central Wiring Changes" → Scenario "Module inclusion is omitted"
**Depends on:** Task 4.1
**Files:**
- `internal/app/modules.go` — create application module inclusion list with platform/capability modules and no business wiring logic.
- `internal/platform/app/app_test.go` — add tests for included fake module becoming active and omitted fake module staying inactive.
- `internal/platform/app/app.go` — expose assembly report for active modules and entry points.
**Test command:** `go test ./internal/platform/app ./internal/app`

**Acceptance criteria:**
- GIVEN a manually created business module is included in the application module list WHEN application assembly runs THEN the module appears in active modules and its declared entry points are exposed through platform adapters.
- GIVEN a manually created business module is not included in the application module list WHEN application assembly runs THEN it does not appear in active modules and its entry points are not exposed.

- [x] **Step 1: Write failing test**
  File: `internal/platform/app/app_test.go`
  Test name: `TestIncludedModuleBecomesActive`
  Expected RED result: FAIL because inclusion list based assembly is absent.

- [x] **Step 2: Verify RED**
  Run: `go test ./internal/platform/app ./internal/app`
  Expected: FAIL for inclusion behavior assertion, not because of compile/setup errors.

- [x] **Step 3: Describe minimal production responsibility**
  File: `internal/app/modules.go`, `internal/platform/app/app.go`
  Responsibility: Provide a single application module list and compose only listed modules without route, capability, or business construction logic in the list itself.

- [x] **Step 4: Verify GREEN**
  Run: `go test ./internal/platform/app ./internal/app`
  Expected: PASS

## 5. Architecture Regression Guardrails

### Task 5.1: Prove manual business modules assemble without generated artifacts

**Type:** regression
**Covers specs:** `specs/module-first-architecture/spec.md` → Requirement "Business Modules Are the Source of Truth" → Scenario "Developer adds a complete business module"; `specs/module-first-architecture/spec.md` → Requirement "Code Generation Is Optional Assistance" → Scenario "Developer creates business behavior manually"
**Depends on:** Task 4.2
**Files:**
- `internal/platform/app/app_test.go` — add a manually declared fake business module test with no generator setup.
- `internal/platform/module/assembler_test.go` — assert manual descriptors are enough for successful assembly.
**Test command:** `go test ./internal/platform/module ./internal/platform/app`

**Acceptance criteria:**
- GIVEN a test business module is manually declared with required capabilities and HTTP entry points WHEN assembly runs THEN it participates in application composition without generated files.
- GIVEN no generator command is run WHEN the test business module is included THEN its declared entry point can still be exposed by the platform adapter.

- [x] **Step 1: Write failing test**
  File: `internal/platform/app/app_test.go`
  Test name: `TestManualBusinessModuleAssemblesWithoutGenerator`
  Expected RED result: FAIL because manual module composition is not yet asserted.

- [x] **Step 2: Verify RED**
  Run: `go test ./internal/platform/module ./internal/platform/app`
  Expected: FAIL for manual module assertion, not because of compile/setup errors.

- [x] **Step 3: Describe minimal production responsibility**
  File: `internal/platform/app/app.go`, `internal/platform/module/assembler.go`
  Responsibility: Treat module descriptors as the runtime source of truth regardless of whether they came from handwritten code or any external artifact.

- [x] **Step 4: Verify GREEN**
  Run: `go test ./internal/platform/module ./internal/platform/app`
  Expected: PASS

### Task 5.2: Ensure OpenAPI and legacy generated files do not drive app assembly

**Type:** regression
**Covers specs:** `specs/module-first-architecture/spec.md` → Requirement "OpenAPI Is Not the Architecture Source of Truth" → Scenario "Application assembles without OpenAPI input"; `specs/module-first-architecture/spec.md` → Requirement "OpenAPI Is Not the Architecture Source of Truth" → Scenario "OpenAPI input is present"; `specs/module-first-architecture/spec.md` → Requirement "Code Generation Is Optional Assistance" → Scenario "Generated output exists"
**Depends on:** Task 5.1
**Files:**
- `internal/platform/app/app_test.go` — add regression proving assembly does not read `doc/example_openapi.yaml` or auto-include legacy generated packages.
- `cmd/app_test.go` — add regression around app startup path using module list, not legacy service/biz/data registration.
- `internal/api/service/service.go` — remove or isolate legacy registration from the new app startup path.
- `internal/api/biz/biz.go` — remove or isolate legacy registration from the new app startup path.
- `internal/api/data/data.go` — remove or isolate legacy registration from the new app startup path.
**Test command:** `go test ./internal/platform/app ./cmd ./internal/api/service ./internal/api/biz ./internal/api/data`

**Acceptance criteria:**
- GIVEN `doc/example_openapi.yaml` exists WHEN application assembly runs THEN assembly uses only included module declarations and does not read OpenAPI input.
- GIVEN legacy or generated API-layer files exist WHEN application assembly runs THEN they do not become active unless included through a module declaration.

- [x] **Step 1: Write failing test**
  File: `internal/platform/app/app_test.go`
  Test name: `TestAssemblyIgnoresOpenAPIAndLegacyGeneratedFiles`
  Expected RED result: FAIL because the new app assembly guardrail is not asserted.

- [x] **Step 2: Verify RED**
  Run: `go test ./internal/platform/app ./cmd ./internal/api/service ./internal/api/biz ./internal/api/data`
  Expected: FAIL for architecture guardrail assertion, not because of compile/setup errors.

- [x] **Step 3: Describe minimal production responsibility**
  File: `internal/platform/app/app.go`, `cmd/app.go`, `internal/api/service/service.go`, `internal/api/biz/biz.go`, `internal/api/data/data.go`
  Responsibility: Make the new runtime assembly depend only on explicit module declarations and keep legacy generated or API-layer files outside the startup path unless wrapped by a module.

- [x] **Step 4: Verify GREEN**
  Run: `go test ./internal/platform/app ./cmd ./internal/api/service ./internal/api/biz ./internal/api/data`
  Expected: PASS

## Final Verification

### Task 6.1: Run focused and full regression checks

**Type:** infrastructure
**Covers specs:** Infrastructure — verifies all scenarios covered by Tasks 1.1 through 5.2
**Depends on:** Task 5.2
**Files:**
- `specs/changes/module-first-architecture/tasks.md` — mark completed tasks during implementation.
**Test command:** `go test ./...`

**Acceptance criteria:**
- GIVEN all planned implementation tasks are complete WHEN `go test ./...` runs THEN the command passes.
- GIVEN the project contains deleted `README.md` and `AGENTS.md` WHEN final verification runs THEN they remain deleted and are not recreated by this change.

- [x] **Step 1: Verify full regression**
  Run: `go test ./...`
  Expected: PASS

- [x] **Step 2: Verify worktree intent**
  Run: `git status --short`
  Expected: Shows intentional spec and architecture changes, preserves user-deleted `README.md` and `AGENTS.md`, and does not include unrelated generated churn.

## Self-Review

- [x] **Spec coverage:** Every Spec Scenario appears in `Spec Coverage Summary` and maps to at least one Task.
- [x] **No unmapped behavior:** Every non-infrastructure task has a concrete Spec Scenario mapping.
- [x] **Infrastructure justified:** Every infrastructure task names the behavior task or scenario it enables.
- [x] **No placeholders:** The plan contains no placeholder markers, vague validation/error handling, or inferred steps.
- [x] **No implementation content:** The plan contains no code blocks, diffs, pseudocode, function bodies, dependency import lists, or concrete control-flow steps.
- [x] **Test-first:** Every behavior, integration, error-path, and regression task starts with a failing test.
- [x] **Concrete commands:** Every task has an exact test, compile, lint, or smoke-test command.
- [x] **Concrete files:** Every task lists exact file paths and responsibilities.
- [x] **Granularity:** No task has more than 5 steps, more than 3 production files, or more than 2 Spec Scenarios, except capability modules grouped together because they share one capability behavior across four existing infra adapters.
- [x] **Dependency order:** Tasks are ordered so the project compiles and relevant tests pass after each task.
- [x] **Consistency:** Later tasks reference files, types, and responsibilities created by earlier tasks accurately.
