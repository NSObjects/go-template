# Proposal: Module-First Architecture

## Intent

The project should help developers focus on business development. The architecture should make business modules the source of truth, while platform capabilities are available by default and can be composed without repeatedly wiring infrastructure code.

This change moves the project away from document-driven or generator-driven scaffolding. Business code should not be organized around OpenAPI artifacts, generated layers, or legacy template assumptions.

## Scope

### In-scope behavior

- Define the project architecture around business modules as the primary unit of development.
- Define platform responsibilities for application startup, configuration, logging, error handling, response handling, lifecycle management, health checks, and observability.
- Define capability modules for optional infrastructure such as relational storage, cache, document storage, messaging, authentication, rate limiting, and metrics.
- Define how a business module declares the platform capabilities it needs and the entry points it exposes.
- Establish that business developers should not need to edit central wiring files for normal feature work.
- Establish that OpenAPI is not the source of truth for the architecture.
- Establish that code generation is optional assistance, not the primary development workflow or architectural dependency.

### Out-of-scope behavior

- No OpenAPI-first workflow.
- No generator-first workflow.
- No requirement to preserve the old README or AGENTS architecture language.
- No large-scale rewrite of all existing business code in the first change.
- No decision yet about concrete package names, framework APIs, or exact interface definitions.
- No decision yet about whether existing example modules stay, move, or get replaced.

## User Workflow

Before this change, a developer adding business behavior must understand framework wiring, layered placement, infrastructure setup, and generated-code expectations before the business module feels complete.

With this change, a developer starts from a business module, declares what capabilities it needs, implements business behavior, and exposes entry points. The platform composes the module with configured capabilities and handles shared runtime concerns.

After this change, normal business development should feel like adding or changing one module, not editing the application skeleton.

## Boundary Decisions

- The business module is the primary development unit.
- The platform owns application composition and shared runtime behavior.
- Capability modules own reusable infrastructure behavior and lifecycle integration.
- Business modules may depend on capability interfaces, but platform and capability modules should not depend on business-specific behavior.
- OpenAPI and code generation may exist later as adapters or helpers, but they must not define the core architecture.
- The first architecture change should define the target shape and migration path before rewriting existing modules.

## Definitions

- Business module - A cohesive unit of business behavior that declares needed capabilities and exposed entry points.
- Platform module - Shared runtime behavior that composes the application and handles non-business concerns.
- Capability module - Reusable infrastructure capability that can be enabled, configured, checked, and used by business modules.
- Entry point - A way for external input to reach business behavior, such as HTTP routes, scheduled jobs, message subscriptions, or command handlers.
- Source of truth - The thing developers primarily edit and reason about when changing behavior.

## Open Questions

- None blocking for proposal approval.
- Non-blocking: whether the first implementation slice should introduce the module mechanism beside the existing code or move one example business module immediately.
- Non-blocking: which capability should be used as the first proof point.

## Approach

Create a module-first architecture direction where the platform composes business modules and capability modules. Business modules declare needs and entry points; capabilities provide reusable infrastructure behavior; platform concerns stay outside business code. The first implementation should be incremental and should prove the model with a small vertical slice before broader migration.

## Observable Success Criteria

- A reviewer can describe the architecture without referring to OpenAPI, generated layers, README, or AGENTS.
- A new business feature can be described as adding or changing a business module.
- Shared runtime concerns are clearly owned by the platform rather than repeated in business modules.
- Optional infrastructure capabilities are described as reusable modules with lifecycle and health behavior.
- The next specification can define observable behavior for module registration, capability declaration, and platform composition without relying on legacy scaffolding assumptions.
