# Repository Development Notes

This project scaffolds API layers (service, biz, data) from OpenAPI descriptions. When working under the repository, follow the
principles below so regenerated code and manual extensions continue to compose cleanly.

## Layer Responsibilities

- **Service layer (`internal/api/service`)**
    - Hosts transport handlers and request/response translation.
    - Perform syntactic validation, bind request parameters, and convert biz responses into the unified `resp` envelope.
    - Do **not** embed business rules or persistence logic. Service functions should remain thin, primarily orchestrating calls to
      biz interfaces and forwarding typed errors.
- **Biz layer (`internal/api/biz`)**
    - Encapsulates domain workflows and business invariants while staying storage-agnostic.
    - Coordinate repository interfaces, handle branching logic (e.g., rate limits, token rotation), and return rich error values
      from the `code` package so higher layers can map them consistently.
    - Keep side-effects limited to calling repositories or emitting domain events; prefer pure functions for validation helpers.
- **Data layer (`internal/api/data`)**
    - Implements concrete persistence adapters and external integrations only.
    - Expose clear interfaces that can be mocked; avoid leaking transport concepts (HTTP status, envelopes) into this layer.
    - Keep `DataManager` focused on infrastructure wiring (DB connections, caches, configs). Business services should receive
      repositories or configuration explicitly rather than depending on `DataManager` directly.

## Coding Conventions

- Favor composition over shared state. Inject dependencies via constructors to keep components testable.
- Prefer Go interfaces that describe behaviors needed by the caller rather than concrete types.
- Organize files by concern: new business logic belongs beside the generated scaffolding, while shared helpers should live in dedicated packages (`internal/resp`, `internal/code`, etc.).
- Use context-aware functions (`ctx context.Context`) for operations that may need cancellation or tracing.
- Handle secrets and credentials via configuration structs—avoid hard-coding values in source files.

## Error Handling & Responses

- Wrap domain failures with the helpers defined in `internal/code`. These errors will be translated by the global
  `ErrorHandler` middleware; service handlers should simply return them.
- Successful responses must go through `resp.SuccessJSON` (or related helpers) to guarantee the
  `{ "msg": "string", "code": 0, "data": { ... } }` contract.
- Add new business error codes to the appropriate `internal/code/*` files and update documentation/tests when doing so.
- Keep HTTP status usage within the allowed set (`200, 201, 400, 401, 403, 404, 500`) unless the specification explicitly
  expands it.

## Testing Expectations

- Write table-driven tests for new logic and prefer standard library assertions (`if`, `t.Helper`, etc.).
- Service layer tests should exercise request binding and response envelopes.
- Biz layer tests should mock repositories to cover edge cases (lockouts, rotations, etc.).
- Run `go test ./...` (or a targeted subset) before submitting changes; include any non-standard flags used in the test output.

## Working With Generated Code

- The OpenAPI specification is the source of truth for generated artifacts. Avoid editing generated files directly—extend
  behavior in new files or protected regions to survive regeneration.
- When updating templates or tooling that emits these layers, regenerate from the spec and ensure the repository still builds
  and passes tests prior to merging.

