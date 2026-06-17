# Project Context

## Purpose

`go-template` is a Go HTTP backend template. It should help a new project start with clear dependency direction, predictable wiring, and a short path for adding business features.

It is not a framework generator. The template should avoid runtime module systems, automatic route registries, generic capability providers, and other abstractions that make ordinary feature work harder to follow.

## Architecture Vocabulary

### Clean-lite

The project uses a lightweight Clean Architecture style. The important rule is dependency direction: business rules stay inside, frameworks and external systems stay outside.

Clean-lite does not require every feature to start with many folders or interfaces. A module may start small and split only when the extra structure improves locality.

### Boot Module

The boot module is the composition root. It is allowed to know about configuration, Echo, databases, caches, repositories, usecases, and route registration.

The boot module is intentionally dirty so the business modules can stay clean.

### Business Module

A business module owns one business capability, such as order intake, account management, or payment routing.

A business module should keep its domain rules, usecases, and adapters close together so developers can understand the full feature without jumping across the whole repository.

### Domain Package

The domain package contains entities, value objects, and business invariants. It imports only the Go standard library.

Create a domain package only when there are real domain rules or durable business concepts.

### Usecase Package

The usecase package contains application workflows. It coordinates domain objects and outbound ports such as stores, transaction runners, clocks, ID generators, or external clients.

Outbound interfaces live in the usecase package because the usecase is the caller that owns the need.

### Adapter Package

An adapter package connects a usecase to an external detail.

Typical adapters are HTTP handlers, MySQL repositories, Redis caches, message publishers, and third-party clients. Adapters may import frameworks and driver libraries.

### Server Package

The server package owns Echo setup, middleware, system routes, HTTP error rendering, and graceful shutdown. It does not know business rules.

### Infrastructure Package

Infrastructure packages create concrete clients for external systems such as Redis, MySQL, object storage, and observability backends.

Infrastructure packages should be called from boot or adapters, not from domain packages.

## Design Biases

- Prefer explicit wiring over automatic discovery.
- Prefer a concrete type until there are multiple implementations or a useful test seam.
- Prefer small outbound interfaces defined by the caller.
- Prefer vertical business modules over top-level horizontal layers.
- Prefer deleting shallow modules over preserving architecture-shaped ceremony.
- Prefer tests through the same interface that production callers use.
