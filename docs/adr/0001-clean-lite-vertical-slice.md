# ADR 0001: Clean-lite Vertical Slice Architecture

## Status

Accepted

## Date

2026-06-16

## Context

The previous architecture grew into a runtime module framework with platform modules, capability providers, module manifests, and generated example business code. That made ordinary development harder because adding a feature required understanding framework-shaped infrastructure before understanding the business capability.

The project still wants the benefits of Clean Architecture:

- business rules independent from frameworks;
- database and HTTP as details;
- testable usecases;
- explicit process startup.

The project does not want ceremony that makes the template feel like a framework generator.

## Decision

Use a Clean-lite vertical slice architecture.

Each business capability lives under one business module:

```text
internal/<business>/
  domain/
  usecase/
  http/
  mysql/
```

The module may omit packages it does not need yet. The usecase package owns outbound interfaces. HTTP and MySQL are adapters. The boot module wires concrete adapters and usecases explicitly.

The boot module is the only intentionally dirty module. It may import configuration, server, infrastructure, and business adapters.

The server module owns Echo and HTTP lifecycle concerns. It must not import business modules.

## Consequences

Feature development becomes predictable:

1. Write or update the usecase.
2. Add only the outbound interface the usecase needs.
3. Implement HTTP and persistence adapters.
4. Wire the pieces in boot.
5. Test the usecase without HTTP or DB.

The architecture gives up framework-independent HTTP routing. This is acceptable because Echo is an outer detail, and replacing it would only require rewriting HTTP adapters and server setup, not business rules.

The architecture also avoids generic provider systems. Shared external clients can still be created by boot or infrastructure packages, but business behavior stays inside business modules.

## Rejected Options

### Runtime module registry

Rejected because it hides wiring and makes feature development depend on manifest conventions.

### Capability provider system

Rejected because one concrete adapter does not justify a generic provider seam.

### Top-level horizontal layers

Rejected because they scatter one business capability across the repository and reduce locality.

### Handler-first CRUD

Rejected because it lets HTTP and DB details absorb business rules. It is fast at the first endpoint and expensive after the first real workflow.
