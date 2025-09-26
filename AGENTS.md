# Repository Development Notes

This project scaffolds API layers (service, biz, data) from OpenAPI descriptions. When adding or modifying functionality under
`internal/api`, keep the following in mind:

## Layer Responsibilities

- **Service layer (`internal/api/service`)** – Hosts transport-level handlers and request/response translation. Service methods
  should focus on validating inbound payloads, orchestrating calls to the biz layer, and mapping errors into transport-friendly
  responses. Avoid embedding business rules or persistence details here so the HTTP/RPC surface stays thin and repeatable.
- **Biz layer (`internal/api/biz`)** – Implements domain orchestration. Biz components coordinate multiple repositories,
  enforce business invariants, and express use cases in a storage-agnostic way. Keep side-effects limited to invoking data-layer
  interfaces, and prefer returning typed domain errors that can be interpreted by upper layers.
- **Data layer (`internal/api/data`)** – Owns persistence concerns and third-party integrations. Implement concrete adapters for
  databases, caches, or external services here. Ensure data-layer functions expose clean contracts (interfaces, DTOs) so they can
  be mocked easily during testing and to keep regeneration safe.

## Working With Generated Code

- **OpenAPI as the source of truth** – The base structs and interfaces for the service, biz, and data layers are generated
  automatically from the OpenAPI specification. Avoid manual edits to generated files that would be overwritten during
  regeneration.
- **Extending generated layers** – Place custom logic in separate files or clearly marked extension points so AI-assisted code
  completion can build atop the generated foundations without merge conflicts.
- **Consistent layer responsibilities** – Ensure changes respect the layering contract documented above to keep the generated
  code predictable for downstream automation.

When updating templates or tooling that emits these layers, verify that regeneration from the OpenAPI definition still produces
compilable artifacts before merging.
