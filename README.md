# go-template

Go HTTP backend template using a Clean-lite vertical slice architecture.

This repository is a starter for business APIs, not a framework generator. Keep business behavior inside vertical business modules, keep wiring explicit in `internal/boot`, and avoid adding default infrastructure that no real feature uses yet.

## Current Shape

```text
cmd/                  CLI entrypoint
internal/boot/        composition root
internal/server/      Echo runtime, middleware, system routes, shutdown
internal/configs/     static config loading and validation
internal/apperr/      framework-free application errors
internal/requestctx/  framework-free request metadata
```

Future business APIs should live under `internal/<business>/`:

```text
internal/<business>/
  usecase/
  http/
  mysql/     # only when persistence is real
  domain/    # only when business invariants are real
```

## Add An API

1. Put workflow input, output, and outbound interfaces in `internal/<business>/usecase`.
2. Add `domain/` only when there are real entities, value objects, or invariants.
3. Put Echo handlers and HTTP DTOs in `internal/<business>/http`.
4. Put DB row models and store code in `internal/<business>/mysql` only when persistence exists.
5. Wire concrete adapters, usecases, and routes in `internal/boot`.
6. Test usecases first, then HTTP or DB adapters.

## Boundaries

- `domain` imports only the Go standard library.
- `usecase` does not import Echo, GORM, Redis, config, server, HTTP response helpers, or concrete adapters.
- Business modules do not import other business modules by default.
- `server` and `configs` do not import business modules.
- Request IDs and trace IDs can come from request metadata; authenticated user identity must come from a real authentication boundary, not client-supplied metadata headers.

## Verify

Run these before handing off Go changes:

```bash
go test ./... -count=1
go vet ./...
go build ./...
git diff --check
```

For Docker or Compose changes:

```bash
docker compose config
docker build --progress=plain --target final -t go-template-arch-check .
```

See `CONTEXT.md`, `docs/architecture.md`, and `docs/adr/` for the architecture decisions behind this template.
