# Proposal: Migrate User Business Module

## Intent

The module-first architecture now exists as a platform foundation. The next step is to prove it with one real business example so developers can see how business behavior lives in a business module instead of being spread across legacy layers or generated artifacts.

This change migrates the existing user example into an explicit business module that declares its entry points and required capabilities through the new module-first assembly path.

## Scope

### In-scope behavior

- Introduce a `user` business module as the first real module-first example.
- The `user` module declares its HTTP entry points through module declarations.
- The `user` module declares its storage capability requirement through module declarations without binding the business module to a concrete database engine.
- The application module list includes the `user` module explicitly.
- Existing user HTTP behavior remains available through the application after migration.
- Legacy user code may be reused behind the module, but it must not drive application assembly directly.
- Tests demonstrate that the `user` module can be included, assembled, and expose user routes without OpenAPI or generator input.

### Out-of-scope behavior

- No new user business features.
- No large rewrite of user domain rules.
- No OpenAPI-first or generator-first workflow.
- No migration of every existing legacy layer module.
- No removal of all legacy API files unless they become unused by this change.
- No change to response envelope, error envelope, or existing route semantics unless needed to keep current behavior working.
- No database schema migration.
- No complete implementation of every possible storage provider; this change only needs enough provider selection behavior to prove the architecture boundary.

## User Workflow

Before this change, the example user behavior exists in legacy service, biz, and data packages, and application startup no longer includes it through the new module-first path.

With this change, a developer can inspect one business module and see how user behavior declares capabilities and HTTP entry points. The application includes the module through the explicit module list, the platform exposes its entry points, and storage provider choice remains outside the business module.

After this change, future business modules have a concrete local example to copy that does not rely on OpenAPI or generated application wiring.

## Boundary Decisions

- The migrated `user` module is the only business module in scope.
- The module-first application list is the only runtime inclusion path for the migrated `user` behavior.
- Existing legacy implementation may be reused as an adapter behind the `user` module to keep the change small.
- OpenAPI documents and generator outputs remain non-authoritative.
- Route behavior should stay compatible with the current user example where practical.

## Definitions

- User business module - The module-first representation of the existing user example behavior.
- Legacy user implementation - Existing user code under the old API layer that may be reused behind the new module declaration.
- Module inclusion - Explicitly adding the `user` module to the application module list.

## Open Questions

- None blocking for proposal approval.
- Non-blocking: whether a later change should delete or archive old user layer files after the module-first replacement is accepted.

## Approach

Move the user example behind a module declaration that owns its entry point and capability declarations. Keep behavior compatibility by reusing existing user logic where appropriate, while ensuring the new application assembly path only sees the module declaration and not OpenAPI or generator output.

## Observable Success Criteria

- The application module list explicitly includes the `user` business module.
- The assembled application report lists `user` as an active business module.
- The assembled application exposes user HTTP routes from the `user` module declaration.
- User route tests pass through the new module-first path.
- OpenAPI and generator output are not required for the `user` module to assemble.
