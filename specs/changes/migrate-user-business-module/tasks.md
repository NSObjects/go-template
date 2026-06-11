# Tasks

## Assumptions

- 用户已确认数据库选择应由架构基础能力解决，业务模块不得绑定具体数据库引擎。
- `user` 业务模块声明抽象 capability `user.storage`；具体 provider 由平台装配和应用配置选择。
- 默认 `user.storage` provider 为 `memory`，用于开箱即用和无外部数据库测试。
- `mysql` 是 `user.storage` 的可选 provider，不是 `user` 业务模块的 requirement 名称。
- 本变更不恢复 `README.md` 或 `AGENTS.md`，不让 OpenAPI 或生成器驱动装配。

## Open Planning Blockers

- 无

## Spec Coverage Summary

| Spec Scenario | Requirement | Covered by Task(s) | Notes |
| --- | --- | --- | --- |
| `specs/user-business-module/spec.md` → Scenario "User module is included" | "User Module Inclusion" | Task 4.1 | 应用模块列表显式包含 user business module。 |
| `specs/user-business-module/spec.md` → Scenario "User module is not included in a custom assembly" | "User Module Inclusion" | Task 5.4 | custom assembly 省略 user 时不暴露 user routes。 |
| `specs/user-business-module/spec.md` → Scenario "User routes are exposed from module declaration" | "User Entry Points Are Declared by the Module" | Task 3.1, Task 5.1 | routes 来自 user module descriptor。 |
| `specs/user-business-module/spec.md` → Scenario "User entry point adapter is unavailable" | "User Entry Points Are Declared by the Module" | Task 5.3 | 缺 HTTP adapter 时启动前失败并指向 user/http。 |
| `specs/user-business-module/spec.md` → Scenario "User storage capability is available" | "User Storage Capability Is Declared" | Task 1.1, Task 3.3 | `user.storage` 由 enabled provider 满足并写入 report。 |
| `specs/user-business-module/spec.md` → Scenario "User storage capability is missing" | "User Storage Capability Is Declared" | Task 1.3 | 无 enabled `user.storage` provider 时启动前失败。 |
| `specs/user-business-module/spec.md` → Scenario "User storage provider is switched" | "User Storage Provider Is Switchable" | Task 1.2, Task 3.4, Task 4.2 | provider 切换不改变 user module routes。 |
| `specs/user-business-module/spec.md` → Scenario "User storage provider is not explicitly selected" | "User Storage Provider Is Switchable" | Task 1.1, Task 2.1, Task 4.1 | 未显式选择时使用 default provider。 |
| `specs/user-business-module/spec.md` → Scenario "User storage provider selection is invalid" | "User Storage Provider Is Switchable" | Task 1.3, Task 2.2, Task 4.2 | invalid provider 错误包含 user、capability 和 provider。 |
| `specs/user-business-module/spec.md` → Scenario "Existing user route is reachable" | "Existing User HTTP Behavior Remains Available" | Task 3.2, Task 5.1 | 请求经 module-first path 到达现有 user 行为。 |
| `specs/user-business-module/spec.md` → Scenario "Existing user route receives invalid input" | "Existing User HTTP Behavior Remains Available" | Task 5.2 | invalid input 通过平台错误 envelope 返回。 |
| `specs/user-business-module/spec.md` → Scenario "User module assembles without OpenAPI input" | "OpenAPI And Generated Output Do Not Drive User Assembly" | Task 5.1 | assembly 不读取 OpenAPI 输入。 |
| `specs/user-business-module/spec.md` → Scenario "Legacy generated files exist" | "OpenAPI And Generated Output Do Not Drive User Assembly" | Task 5.4 | legacy/generated files 不会自动激活 user。 |

## 1. Platform Capability Provider Selection

### Task 1.1: Resolve default capability provider

**Type:** behavior
**Covers specs:** `specs/user-business-module/spec.md` → Requirement "User Storage Capability Is Declared" → Scenario "User storage capability is available"; `specs/user-business-module/spec.md` → Requirement "User Storage Provider Is Switchable" → Scenario "User storage provider is not explicitly selected"
**Depends on:** none
**Files:**
- `internal/platform/module/assembler_test.go` — add tests for default provider resolution and report provider identity.
- `internal/platform/module/types.go` — add provider identity and default marker to capability declarations, plus capability selection type.
- `internal/platform/module/report.go` — preserve selected provider identity in capability and requirement report records.
- `internal/platform/module/assembler.go` — resolve an enabled default provider when no explicit selection is supplied.
**Test command:** `go test ./internal/platform/module`

**Acceptance criteria:**
- GIVEN a business module requires `user.storage` and an enabled default provider exists WHEN module assembly runs THEN the requirement is satisfied and the report identifies provider `memory`.
- GIVEN no concrete provider is explicitly selected WHEN module assembly runs with a default provider THEN assembly succeeds and the selected provider is observable in the report.

- [x] **Step 1: Write failing test**
  File: `internal/platform/module/assembler_test.go`
  Test name: `TestAssembleUsesDefaultCapabilityProvider`
  Expected RED result: FAIL because provider identity and default provider resolution are not implemented.

- [x] **Step 2: Verify RED**
  Run: `go test ./internal/platform/module`
  Expected: FAIL for missing default provider behavior, not because of compile/setup errors.

- [x] **Step 3: Describe minimal production responsibility**
  File: `internal/platform/module/types.go`, `internal/platform/module/report.go`, `internal/platform/module/assembler.go`
  Responsibility: Represent provider identity on capability declarations and choose one enabled default provider to satisfy an abstract requirement when no explicit selection is supplied.

- [x] **Step 4: Verify GREEN**
  Run: `go test ./internal/platform/module`
  Expected: PASS

### Task 1.2: Resolve explicitly selected capability provider

**Type:** behavior
**Covers specs:** `specs/user-business-module/spec.md` → Requirement "User Storage Provider Is Switchable" → Scenario "User storage provider is switched"
**Depends on:** Task 1.1
**Files:**
- `internal/platform/module/assembler_test.go` — add test for selecting a non-default provider.
- `internal/platform/module/assembler.go` — apply explicit capability selections before default provider resolution.
**Test command:** `go test ./internal/platform/module`

**Acceptance criteria:**
- GIVEN `user.storage` has enabled providers `memory` and `mysql` WHEN assembly selects provider `mysql` THEN the user requirement is satisfied by provider `mysql`.
- GIVEN the selected provider changes from `memory` to `mysql` WHEN assembly runs again THEN the business module requirement name remains `user.storage` and only the reported provider changes.

- [x] **Step 1: Write failing test**
  File: `internal/platform/module/assembler_test.go`
  Test name: `TestAssembleUsesSelectedCapabilityProvider`
  Expected RED result: FAIL because explicit provider selection is not applied.

- [x] **Step 2: Verify RED**
  Run: `go test ./internal/platform/module`
  Expected: FAIL for selected provider assertion, not because of compile/setup errors.

- [x] **Step 3: Describe minimal production responsibility**
  File: `internal/platform/module/assembler.go`
  Responsibility: Match a requirement to the enabled provider named by the capability selection and record that provider in the requirement report.

- [x] **Step 4: Verify GREEN**
  Run: `go test ./internal/platform/module`
  Expected: PASS

### Task 1.3: Reject missing and unavailable selected providers

**Type:** error-path
**Covers specs:** `specs/user-business-module/spec.md` → Requirement "User Storage Capability Is Declared" → Scenario "User storage capability is missing"; `specs/user-business-module/spec.md` → Requirement "User Storage Provider Is Switchable" → Scenario "User storage provider selection is invalid"
**Depends on:** Task 1.2
**Files:**
- `internal/platform/module/assembler_test.go` — add tests for missing `user.storage` provider and unavailable selected provider.
- `internal/platform/module/errors.go` — add typed provider selection error including module, capability, and provider.
- `internal/platform/module/assembler.go` — reject unsatisfied abstract capabilities and selected providers that cannot satisfy the requirement.
**Test command:** `go test ./internal/platform/module`

**Acceptance criteria:**
- GIVEN a business module requires `user.storage` and no enabled provider exists WHEN module assembly runs THEN assembly fails and the error identifies module `user` and capability `user.storage`.
- GIVEN provider `unknown` is selected for `user.storage` WHEN no enabled provider named `unknown` is available THEN assembly fails and the error identifies module `user`, capability `user.storage`, and provider `unknown`.

- [x] **Step 1: Write failing test**
  File: `internal/platform/module/assembler_test.go`
  Test name: `TestAssembleRejectsUnavailableSelectedCapabilityProvider`
  Expected RED result: FAIL because unavailable selected provider errors do not include provider identity.

- [x] **Step 2: Verify RED**
  Run: `go test ./internal/platform/module`
  Expected: FAIL for provider selection error assertion, not because of compile/setup errors.

- [x] **Step 3: Describe minimal production responsibility**
  File: `internal/platform/module/errors.go`, `internal/platform/module/assembler.go`
  Responsibility: Return startup-blocking typed errors when an abstract capability has no enabled provider or when a selected provider is absent, disabled, or unavailable.

- [x] **Step 4: Verify GREEN**
  Run: `go test ./internal/platform/module`
  Expected: PASS

## 2. Configuration And Platform App Wiring

### Task 2.1: Load capability provider selections from config

**Type:** infrastructure
**Covers specs:** Infrastructure — enables Task 4.1 and Task 4.2 for `specs/user-business-module/spec.md` → Requirement "User Storage Provider Is Switchable"
**Depends on:** Task 1.1
**Files:**
- `internal/configs/capabilities_test.go` — add tests for reading capability provider selections and omitted selections.
- `internal/configs/config.go` — add config shape for capability provider selection.
- `configs/config.toml` — add an empty provider selection section showing default behavior without selecting MySQL.
**Test command:** `go test ./internal/configs`

**Acceptance criteria:**
- GIVEN config contains a provider selection for `user.storage` WHEN config is loaded THEN the loaded config exposes provider `mysql` for that capability.
- GIVEN config omits a provider selection for `user.storage` WHEN config is loaded THEN the loaded config preserves an empty selection so default provider behavior can apply.

- [x] **Step 1: Create infrastructure artifact**
  File: `internal/configs/config.go`, `configs/config.toml`
  Responsibility: Expose capability provider selections from configuration while allowing omitted selections to represent default provider behavior.

- [x] **Step 2: Verify infrastructure**
  Run: `go test ./internal/configs`
  Expected: PASS

### Task 2.2: Pass capability selections through platform app assembly

**Type:** integration
**Covers specs:** `specs/user-business-module/spec.md` → Requirement "User Storage Provider Is Switchable" → Scenario "User storage provider is switched"; `specs/user-business-module/spec.md` → Requirement "User Storage Provider Is Switchable" → Scenario "User storage provider selection is invalid"
**Depends on:** Task 1.3
**Files:**
- `internal/platform/app/app_test.go` — add tests proving app assembly applies provider selections and surfaces provider selection errors.
- `internal/platform/app/app.go` — pass capability selections from app options into module assembly.
**Test command:** `go test ./internal/platform/app`

**Acceptance criteria:**
- GIVEN platform app options select provider `mysql` for `user.storage` WHEN the app is assembled THEN the report marks the user requirement satisfied by provider `mysql`.
- GIVEN platform app options select provider `unknown` for `user.storage` WHEN the app is assembled THEN assembly fails before server construction and identifies the unavailable provider.

- [x] **Step 1: Write failing test**
  File: `internal/platform/app/app_test.go`
  Test name: `TestAppAssemblyAppliesCapabilitySelections`
  Expected RED result: FAIL because platform app options do not pass provider selections into module assembly.

- [x] **Step 2: Verify RED**
  Run: `go test ./internal/platform/app`
  Expected: FAIL for missing provider selection propagation, not because of compile/setup errors.

- [x] **Step 3: Describe minimal production responsibility**
  File: `internal/platform/app/app.go`
  Responsibility: Accept capability selections in app assembly options and pass them to platform module assembly before routes are adapted.

- [x] **Step 4: Verify GREEN**
  Run: `go test ./internal/platform/app`
  Expected: PASS

## 3. User Module And Storage Providers

### Task 3.1: Declare user module routes and abstract storage requirement

**Type:** behavior
**Covers specs:** `specs/user-business-module/spec.md` → Requirement "User Entry Points Are Declared by the Module" → Scenario "User routes are exposed from module declaration"; `specs/user-business-module/spec.md` → Requirement "User Storage Capability Is Declared" → Scenario "User storage capability is available"
**Depends on:** Task 1.1
**Files:**
- `internal/modules/user/module_test.go` — add tests for user descriptor kind, route declarations, and abstract storage requirement.
- `internal/modules/user/module.go` — create user business module descriptor.
- `internal/modules/user/routes.go` — create route descriptors for existing user HTTP operations.
- `internal/api/service/user.go` — return concrete user controller from constructor while keeping router compatibility.
**Test command:** `go test ./internal/modules/user ./internal/api/service`

**Acceptance criteria:**
- GIVEN the user module descriptor is inspected WHEN its requirements are read THEN it requires `user.storage` and does not name a concrete database provider.
- GIVEN the user module descriptor is inspected WHEN its entry points are read THEN it exposes the existing user HTTP route set owned by module `user`.

- [x] **Step 1: Write failing test**
  File: `internal/modules/user/module_test.go`
  Test name: `TestDescriptorDeclaresUserStorageAndHTTPRoutes`
  Expected RED result: FAIL because `internal/modules/user` does not exist and user routes are not declared by a module.

- [x] **Step 2: Verify RED**
  Run: `go test ./internal/modules/user ./internal/api/service`
  Expected: FAIL for missing user module package or missing descriptor behavior, not because of unrelated packages.

- [x] **Step 3: Describe minimal production responsibility**
  File: `internal/modules/user/module.go`, `internal/modules/user/routes.go`, `internal/api/service/user.go`
  Responsibility: Wrap existing user controller behavior as module-owned HTTP route descriptors and declare `user.storage` as the only user storage requirement.

- [x] **Step 4: Verify GREEN**
  Run: `go test ./internal/modules/user ./internal/api/service`
  Expected: PASS

### Task 3.2: Provide default in-memory user repository

**Type:** behavior
**Covers specs:** `specs/user-business-module/spec.md` → Requirement "User Storage Provider Is Switchable" → Scenario "User storage provider is not explicitly selected"; `specs/user-business-module/spec.md` → Requirement "Existing User HTTP Behavior Remains Available" → Scenario "Existing user route is reachable"
**Depends on:** Task 3.1
**Files:**
- `internal/modules/user/storage/memory/repository_test.go` — add tests for the repository behavior used by existing user use cases.
- `internal/modules/user/storage/memory/repository.go` — create in-memory implementation of the existing user repository interface.
**Test command:** `go test ./internal/modules/user/storage/memory`

**Acceptance criteria:**
- GIVEN a new in-memory user repository WHEN a user is created and listed THEN the list contains the created user data in the existing response data shape.
- GIVEN a missing user id WHEN the repository is asked for that user THEN it returns an error that can flow through existing user behavior.

- [x] **Step 1: Write failing test**
  File: `internal/modules/user/storage/memory/repository_test.go`
  Test name: `TestRepositorySupportsUserLifecycle`
  Expected RED result: FAIL because the default in-memory user repository does not exist.

- [x] **Step 2: Verify RED**
  Run: `go test ./internal/modules/user/storage/memory`
  Expected: FAIL for missing repository package or behavior, not because of compile/setup errors.

- [x] **Step 3: Describe minimal production responsibility**
  File: `internal/modules/user/storage/memory/repository.go`
  Responsibility: Provide an in-memory repository implementing existing user repository operations so the default provider can support user HTTP behavior without external services.

- [x] **Step 4: Verify GREEN**
  Run: `go test ./internal/modules/user/storage/memory`
  Expected: PASS

### Task 3.3: Declare memory as default user storage provider

**Type:** behavior
**Covers specs:** `specs/user-business-module/spec.md` → Requirement "User Storage Capability Is Declared" → Scenario "User storage capability is available"; `specs/user-business-module/spec.md` → Requirement "User Storage Provider Is Switchable" → Scenario "User storage provider is not explicitly selected"
**Depends on:** Task 3.2
**Files:**
- `internal/modules/user/storage/memory/module_test.go` — add tests for provider descriptor name, capability, status, and default marker.
- `internal/modules/user/storage/memory/module.go` — expose memory provider as a capability module for `user.storage`.
**Test command:** `go test ./internal/modules/user/storage/memory`

**Acceptance criteria:**
- GIVEN the memory provider module is inspected WHEN its descriptor is read THEN it provides enabled capability `user.storage` with provider `memory`.
- GIVEN no concrete user storage provider is selected WHEN memory provider is included THEN it is eligible as the default provider.

- [x] **Step 1: Write failing test**
  File: `internal/modules/user/storage/memory/module_test.go`
  Test name: `TestModuleProvidesDefaultUserStorageCapability`
  Expected RED result: FAIL because memory provider module descriptor does not exist.

- [x] **Step 2: Verify RED**
  Run: `go test ./internal/modules/user/storage/memory`
  Expected: FAIL for missing provider descriptor behavior, not because of compile/setup errors.

- [x] **Step 3: Describe minimal production responsibility**
  File: `internal/modules/user/storage/memory/module.go`
  Responsibility: Describe memory as the enabled default provider for `user.storage` and expose the repository needed by the user module builder.

- [x] **Step 4: Verify GREEN**
  Run: `go test ./internal/modules/user/storage/memory`
  Expected: PASS

### Task 3.4: Declare MySQL as optional user storage provider

**Type:** behavior
**Covers specs:** `specs/user-business-module/spec.md` → Requirement "User Storage Provider Is Switchable" → Scenario "User storage provider is switched"; `specs/user-business-module/spec.md` → Requirement "User Storage Provider Is Switchable" → Scenario "User storage provider selection is invalid"
**Depends on:** Task 3.1
**Files:**
- `internal/modules/user/storage/mysql/module_test.go` — add tests for enabled, disabled, and invalid-config provider descriptors.
- `internal/modules/user/storage/mysql/module.go` — expose MySQL as an optional provider for `user.storage`.
**Test command:** `go test ./internal/modules/user/storage/mysql`

**Acceptance criteria:**
- GIVEN MySQL user storage config is enabled and valid WHEN the provider descriptor is read THEN it provides capability `user.storage` with provider `mysql`.
- GIVEN MySQL user storage config is disabled or invalid WHEN provider selection asks for `mysql` THEN that provider cannot satisfy `user.storage`.

- [x] **Step 1: Write failing test**
  File: `internal/modules/user/storage/mysql/module_test.go`
  Test name: `TestModuleReportsUserStorageProviderStatus`
  Expected RED result: FAIL because MySQL user storage provider descriptor does not exist.

- [x] **Step 2: Verify RED**
  Run: `go test ./internal/modules/user/storage/mysql`
  Expected: FAIL for missing MySQL provider behavior, not because of compile/setup errors.

- [x] **Step 3: Describe minimal production responsibility**
  File: `internal/modules/user/storage/mysql/module.go`
  Responsibility: Describe MySQL as an optional provider for `user.storage` and adapt existing user data repository construction without changing the user module requirement name.

- [x] **Step 4: Verify GREEN**
  Run: `go test ./internal/modules/user/storage/mysql`
  Expected: PASS

## 4. Application Module Inclusion

### Task 4.1: Include user module with default storage provider

**Type:** integration
**Covers specs:** `specs/user-business-module/spec.md` → Requirement "User Module Inclusion" → Scenario "User module is included"; `specs/user-business-module/spec.md` → Requirement "User Storage Provider Is Switchable" → Scenario "User storage provider is not explicitly selected"
**Depends on:** Task 2.1, Task 3.3
**Files:**
- `internal/app/modules_test.go` — add tests for included user module and default storage provider selection.
- `internal/app/modules.go` — return app module result containing capability modules, memory provider, user module, and capability selections.
- `cmd/app.go` — use the app module result when assembling the platform app.
**Test command:** `go test ./internal/app ./cmd`

**Acceptance criteria:**
- GIVEN application modules are built from default config WHEN the module list is inspected THEN it contains user as a business module.
- GIVEN no provider is selected for `user.storage` WHEN application modules are assembled THEN default provider `memory` satisfies the user storage requirement.

- [x] **Step 1: Write failing test**
  File: `internal/app/modules_test.go`
  Test name: `TestModulesIncludesUserWithDefaultStorageProvider`
  Expected RED result: FAIL because application modules do not include the user business module or default user storage provider.

- [x] **Step 2: Verify RED**
  Run: `go test ./internal/app ./cmd`
  Expected: FAIL for missing application module inclusion behavior, not because of compile/setup errors.

- [x] **Step 3: Describe minimal production responsibility**
  File: `internal/app/modules.go`, `cmd/app.go`
  Responsibility: Build the explicit app module list with user, user storage providers, low-level capabilities, and capability selections, then pass that result into platform app assembly.

- [x] **Step 4: Verify GREEN**
  Run: `go test ./internal/app ./cmd`
  Expected: PASS

### Task 4.2: Switch user storage provider from application config

**Type:** integration
**Covers specs:** `specs/user-business-module/spec.md` → Requirement "User Storage Provider Is Switchable" → Scenario "User storage provider is switched"; `specs/user-business-module/spec.md` → Requirement "User Storage Provider Is Switchable" → Scenario "User storage provider selection is invalid"
**Depends on:** Task 3.4, Task 4.1
**Files:**
- `internal/app/modules_test.go` — add tests for configured provider switch and invalid provider selection.
- `internal/app/modules.go` — apply configured provider selection for `user.storage`.
**Test command:** `go test ./internal/app`

**Acceptance criteria:**
- GIVEN config selects provider `mysql` for `user.storage` WHEN application modules are assembled THEN user remains active and the report identifies provider `mysql`.
- GIVEN config selects provider `unknown` for `user.storage` WHEN application assembly runs THEN startup fails before serving and identifies the unavailable selected provider.

- [x] **Step 1: Write failing test**
  File: `internal/app/modules_test.go`
  Test name: `TestModulesSelectsConfiguredUserStorageProvider`
  Expected RED result: FAIL because configured provider selection is not applied by the app module builder.

- [x] **Step 2: Verify RED**
  Run: `go test ./internal/app`
  Expected: FAIL for missing configured provider selection behavior, not because of compile/setup errors.

- [x] **Step 3: Describe minimal production responsibility**
  File: `internal/app/modules.go`
  Responsibility: Convert configured capability provider selections into platform module selections while keeping the user module descriptor and route descriptors unchanged.

- [x] **Step 4: Verify GREEN**
  Run: `go test ./internal/app`
  Expected: PASS

## 5. HTTP And Assembly Regression

### Task 5.1: Expose existing user route through module-first path

**Type:** integration
**Covers specs:** `specs/user-business-module/spec.md` → Requirement "Existing User HTTP Behavior Remains Available" → Scenario "Existing user route is reachable"; `specs/user-business-module/spec.md` → Requirement "OpenAPI And Generated Output Do Not Drive User Assembly" → Scenario "User module assembles without OpenAPI input"
**Depends on:** Task 4.1
**Files:**
- `internal/modules/user/http_test.go` — add server-level test invoking an existing user route through module-first assembly.
- `internal/modules/user/routes.go` — adjust route descriptors only if the integration test exposes a route registration mismatch.
**Test command:** `go test ./internal/modules/user`

**Acceptance criteria:**
- GIVEN the application is assembled with user and default storage provider WHEN `GET /api/users` is invoked through the server THEN the response uses the existing success envelope format.
- GIVEN no OpenAPI input is provided to assembly WHEN the user module is included THEN user HTTP entry points are available from module declarations.

- [x] **Step 1: Write failing test**
  File: `internal/modules/user/http_test.go`
  Test name: `TestExistingUserRouteReachableThroughModuleFirstPath`
  Expected RED result: FAIL because user routes are not yet verified through the assembled server path.

- [x] **Step 2: Verify RED**
  Run: `go test ./internal/modules/user`
  Expected: FAIL for route reachability assertion, not because of compile/setup errors.

- [x] **Step 3: Describe minimal production responsibility**
  File: `internal/modules/user/routes.go`
  Responsibility: Ensure module-declared user HTTP routes register through the platform server and reach existing user behavior without OpenAPI input.

- [x] **Step 4: Verify GREEN**
  Run: `go test ./internal/modules/user`
  Expected: PASS

### Task 5.2: Preserve platform error envelope for invalid user input

**Type:** error-path
**Covers specs:** `specs/user-business-module/spec.md` → Requirement "Existing User HTTP Behavior Remains Available" → Scenario "Existing user route receives invalid input"
**Depends on:** Task 5.1
**Files:**
- `internal/modules/user/http_test.go` — add server-level invalid input test for existing user route.
- `internal/modules/user/routes.go` — adjust route descriptors only if invalid input bypasses platform error handling.
**Test command:** `go test ./internal/modules/user`

**Acceptance criteria:**
- GIVEN user routes are registered through module-first assembly WHEN `POST /api/users` receives invalid input THEN the response uses the existing error envelope format.
- GIVEN invalid input is rejected WHEN the request is handled THEN user behavior does not bypass the platform error handler.

- [x] **Step 1: Write failing test**
  File: `internal/modules/user/http_test.go`
  Test name: `TestInvalidUserInputUsesPlatformErrorEnvelope`
  Expected RED result: FAIL because invalid input behavior is not yet covered through the module-first server path.

- [x] **Step 2: Verify RED**
  Run: `go test ./internal/modules/user`
  Expected: FAIL for invalid input envelope assertion, not because of compile/setup errors.

- [x] **Step 3: Describe minimal production responsibility**
  File: `internal/modules/user/routes.go`
  Responsibility: Preserve existing service error return behavior so platform middleware converts invalid user input into the existing error envelope.

- [x] **Step 4: Verify GREEN**
  Run: `go test ./internal/modules/user`
  Expected: PASS

### Task 5.3: Reject user HTTP routes when HTTP adapter is unavailable

**Type:** error-path
**Covers specs:** `specs/user-business-module/spec.md` → Requirement "User Entry Points Are Declared by the Module" → Scenario "User entry point adapter is unavailable"
**Depends on:** Task 3.1
**Files:**
- `internal/platform/app/app_test.go` — add user-specific unsupported entry point adapter test.
- `internal/platform/app/app.go` — adjust adapter validation only if the test shows user module errors are not surfaced.
**Test command:** `go test ./internal/platform/app`

**Acceptance criteria:**
- GIVEN the user module declares HTTP entry points and the platform app has no HTTP adapter WHEN assembly runs THEN assembly fails before serving.
- GIVEN assembly fails for unsupported user HTTP entry points WHEN the error is inspected THEN it identifies module `user` and entry point type `http`.

- [x] **Step 1: Write failing test**
  File: `internal/platform/app/app_test.go`
  Test name: `TestUserModuleFailsWhenHTTPAdapterUnavailable`
  Expected RED result: FAIL because user-specific unsupported adapter behavior is not covered.

- [x] **Step 2: Verify RED**
  Run: `go test ./internal/platform/app`
  Expected: FAIL for unsupported adapter assertion, not because of compile/setup errors.

- [x] **Step 3: Describe minimal production responsibility**
  File: `internal/platform/app/app.go`
  Responsibility: Surface platform module adapter validation errors with the affected user module and unsupported entry point type.

- [x] **Step 4: Verify GREEN**
  Run: `go test ./internal/platform/app`
  Expected: PASS

### Task 5.4: Keep legacy and generated user files non-authoritative

**Type:** regression
**Covers specs:** `specs/user-business-module/spec.md` → Requirement "User Module Inclusion" → Scenario "User module is not included in a custom assembly"; `specs/user-business-module/spec.md` → Requirement "OpenAPI And Generated Output Do Not Drive User Assembly" → Scenario "Legacy generated files exist"
**Depends on:** Task 4.1
**Files:**
- `internal/platform/app/app_test.go` — add regression tests proving legacy user files do not activate user without module inclusion.
- `internal/app/modules_test.go` — add regression assertion that explicit app inclusion, not legacy files, makes user active.
**Test command:** `go test ./internal/platform/app ./internal/app`

**Acceptance criteria:**
- GIVEN a custom application assembly omits the user module WHEN assembly runs THEN the report does not list user as active and no user routes are exposed.
- GIVEN legacy generated or layered user files exist WHEN app assembly uses an explicit module list without user THEN generated output does not activate user or override user module declarations.

- [x] **Step 1: Write failing test**
  File: `internal/platform/app/app_test.go`
  Test name: `TestLegacyUserFilesDoNotActivateUserModule`
  Expected RED result: FAIL if legacy user files can activate user without explicit module inclusion.

- [x] **Step 2: Verify RED**
  Run: `go test ./internal/platform/app ./internal/app`
  Expected: FAIL for legacy activation assertion if the old path is still authoritative, not because of compile/setup errors.

- [x] **Step 3: Describe minimal production responsibility**
  File: `internal/platform/app/app_test.go`, `internal/app/modules_test.go`
  Responsibility: Lock module-first assembly to explicit module inclusion and prevent generated or layered user files from becoming an implicit runtime source.

- [x] **Step 4: Verify GREEN**
  Run: `go test ./internal/platform/app ./internal/app`
  Expected: PASS

## Final Verification

- Run `go test ./...` after all tasks are completed.
- Confirm every task checkbox is checked before reporting the change ready.

## Self-Review

- [x] **Spec coverage:** Every Spec Scenario appears in `Spec Coverage Summary` and maps to at least one Task.
- [x] **No unmapped behavior:** Every non-infrastructure task has a concrete Spec Scenario mapping.
- [x] **Infrastructure justified:** Every infrastructure task names the behavior task or scenario it enables.
- [x] **No placeholders:** The plan contains no `TBD`, `TODO`, vague validation/error handling, or inferred steps.
- [x] **No implementation content:** The plan contains no code blocks, diffs, pseudocode, function bodies, imports, or concrete control-flow steps.
- [x] **Test-first:** Every behavior, integration, error-path, and regression task starts with a failing test.
- [x] **Concrete commands:** Every task has an exact test, compile, lint, or smoke-test command.
- [x] **Concrete files:** Every task lists exact file paths and responsibilities.
- [x] **Granularity:** No task has more than 5 steps, more than 3 production files, or more than 2 Spec Scenarios.
- [x] **Dependency order:** Tasks are ordered so the project compiles and relevant tests pass after each task.
- [x] **Consistency:** Later tasks reference files, types, and responsibilities created by earlier tasks accurately.
