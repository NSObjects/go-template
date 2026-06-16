# Architecture

## Summary

This project uses a Clean-lite architecture: Clean Architecture dependency rules with Go-friendly package layout and minimal ceremony.

The goal is not to create a reusable framework inside the app. The goal is to make feature work easy:

- a developer can find a feature in one business module;
- business logic can be tested without Echo, DB, Redis, or config;
- adapters can change without rewriting usecases;
- startup wiring is explicit and boring;
- abstractions are introduced only when they buy locality or test leverage.

## Design Constraints

1. Dependencies point inward.
2. Business modules do not import framework or driver details.
3. The boot module is the only place where all concrete pieces are assembled.
4. Route registration is explicit.
5. Business modules are vertical slices, not scattered horizontal layers.
6. Interfaces are owned by the caller that needs them.
7. A module starts compact and splits only when it becomes easier to maintain.

## Target Layout

```text
cmd/
  root.go
  app.go

internal/
  boot/
    app.go

  server/
    server.go
    middlewares/
    httpresp/

  configs/

  apperr/
  requestctx/

  <business>/
    domain/
    usecase/
    http/
    mysql/
```

Small business modules may omit folders they do not need yet. For example, a simple read-only feature may start with only:

```text
internal/report/
  usecase/
  http/
  mysql/
```

Add `domain/` when the module has real business concepts or invariants. Add `memory/`, `queue/`, `client/`, or other concrete adapters only when those details exist. Split `internal/boot` into helper files such as `routes.go` or `resources.go` only after real route or resource wiring makes `app.go` hard to read.

## Dependency Direction

```text
cmd
  -> boot

boot
  -> configs
  -> server
  -> infrastructure packages
  -> business/http
  -> business/mysql
  -> business/usecase

business/http
  -> business/usecase
  -> server/httpresp

business/mysql
  -> business/usecase
  -> business/domain

business/usecase
  -> business/domain
  -> apperr
  -> requestctx

business/domain
  -> standard library
```

Forbidden imports:

- `domain` must not import Echo, GORM, Redis, configs, server, http response helpers, or global loggers.
- `usecase` must not import Echo, GORM, Redis, configs, server, or concrete adapter packages.
- `http` adapters must not import persistence adapters, domain packages, driver details, or configs.
- `mysql` adapters must not import HTTP adapters, server packages, or configs.
- `server` must not import business modules.
- `configs` must not import business modules.
- `boot` may import everything needed for wiring.

## Boot Module

`internal/boot` is the composition root. It owns process-level wiring:

```go
cfg, err := configs.Load(configPath)
srv, err := server.New(cfg)
if err != nil {
	return err
}

db := mysqlinfra.Open(databaseConfig)
orders := ordermysql.NewStore(db)
orderUC := orderusecase.New(orders)
orderHTTP := orderhttp.New(orderUC)

orderhttp.Register(srv.API().Group("/orders"), orderHTTP)

return srv.Run(ctx)
```

This code is allowed to be concrete. It should be easy to read and easy to change. Avoid hiding it behind a module registry or provider system.

## Server Module

`internal/server` owns Echo details:

- Echo instance creation;
- middleware setup;
- validation adapter;
- system routes such as health and info;
- HTTP error rendering;
- graceful shutdown.

It should expose a small interface:

```go
type Server struct { ... }

func New(cfg configs.Config) *Server
func (s *Server) API() *echo.Group
func (s *Server) Run(ctx context.Context) error
```

The server module may expose Echo groups because Echo is the chosen HTTP framework. Clean-lite protects the business modules, not the outer HTTP adapter.

## Business Module Shape

Use one vertical module per business capability:

```text
internal/order/
  domain/
    order.go
    money.go
  usecase/
    create.go
    get.go
    ports.go
    errors.go
  http/
    handler.go
    dto.go
    routes.go
  mysql/
    store.go
    model.go
```

The common request path is:

```text
HTTP request
  -> order/http
  -> order/usecase
  -> order/domain
  -> order/usecase outbound interface
  -> order/mysql
  -> database
```

### Domain

Domain code contains business facts that should remain true no matter how the app is delivered or stored.

Good domain code:

- validates invariants;
- names business concepts;
- computes state transitions;
- contains no web or database details.

Do not create domain packages for plain CRUD records with no business rule yet.

### Usecase

Usecases perform application workflows. They accept context plus command/query input and return output plus error.

```go
type Store interface {
	Save(ctx context.Context, order domain.Order) error
	Find(ctx context.Context, id domain.OrderID) (domain.Order, error)
}

type Usecase struct {
	store Store
}

func New(store Store) *Usecase {
	return &Usecase{store: store}
}
```

Prefer one concrete `Usecase` type per small module. Split into focused usecase types only when the method set becomes hard to understand or test.

Do not create inbound interfaces by default. Handlers can depend on `*usecase.Usecase` until there is a real reason to hide it.

### HTTP Adapter

HTTP adapters translate HTTP into usecase calls:

- bind and validate request DTOs;
- extract authenticated identity and request metadata;
- call usecase methods;
- map usecase output to response DTOs;
- render errors through `server/httpresp`.

HTTP adapters should not contain business decisions.

### MySQL Adapter

MySQL adapters implement usecase outbound interfaces:

- own DB row models;
- own SQL or GORM details;
- translate between rows and domain values;
- wrap driver errors into `apperr` or module errors.

DB row models must not cross into usecase or domain packages.

## Error Design

Use two layers:

```text
internal/apperr      # framework-free app errors and codes
internal/server/httpresp  # HTTP status and JSON rendering
```

`apperr` should not know HTTP. It should express:

- code;
- safe client message;
- internal detail;
- kind such as validation, not found, permission, conflict, internal;
- wrapped cause.

`server/httpresp` maps `apperr.Kind` to HTTP status and response JSON.

Business modules may define module-specific errors near the usecase or domain that raises them.

## Request Context

Use `context.Context` for cancellation and deadlines.

Use `internal/requestctx` for framework-free request metadata:

- request ID;
- trace ID;
- user ID;
- roles or permissions when needed.

Echo extraction lives in the HTTP adapter or server middleware. Usecases receive ordinary `context.Context` and read request metadata through `requestctx`.

## External Resources

Create external clients in boot, business-specific adapters, or narrowly named infrastructure packages. Pass concrete adapters into usecases through small outbound interfaces.

Do not add generic capability providers.

Examples:

- MySQL client factory belongs in an infrastructure package, a business adapter, or a boot helper.
- Cache or queue client factories belong in a business adapter or a narrowly named infrastructure package when real cache or queue-backed behavior exists.
- Business-specific cache or queue behavior belongs in the business module's adapter.

## Transactions

Do not introduce a global transaction abstraction before a usecase needs it.

When a usecase must coordinate multiple stores atomically, define the smallest transaction interface in that usecase package:

```go
type TxRunner interface {
	WithinTx(ctx context.Context, fn func(ctx context.Context, tx Stores) error) error
}
```

The MySQL adapter implements it. Other modules should not depend on that transaction shape unless they share the same usecase need.

## Development Flow

### Add a new endpoint

1. Create or open the business module.
2. Add a usecase method and its input/output types.
3. Define outbound interfaces only for external behavior the usecase needs.
4. Implement an HTTP handler.
5. Implement a MySQL or memory adapter if persistence is needed.
6. Wire the concrete pieces in `internal/boot/app.go`, or in a boot helper file once real wiring pressure exists.
7. Add usecase tests first, then adapter tests for HTTP or DB behavior.

### Add a new external system

1. Keep the client factory outside the business module unless it is business-specific.
2. Define a usecase-owned outbound interface.
3. Implement the adapter next to the business module if behavior is business-specific.
4. Wire it in boot.

### Split a module

Split only when locality improves:

- `usecase.go` has unrelated workflows;
- domain rules are mixed with request/response mapping;
- adapter details make usecase tests hard;
- multiple adapters now satisfy the same outbound interface.

Do not split only because a diagram has another layer.

## Testing Strategy

Domain tests:

- pure Go;
- no Echo;
- no DB;
- no config.

Usecase tests:

- fake outbound interfaces in test files;
- no Echo;
- no real DB unless the usecase is explicitly testing integration behavior.

HTTP adapter tests:

- Echo recorder is allowed;
- fake or real usecase depending on test value;
- verify request mapping, response shape, and error rendering.

MySQL adapter tests:

- integration tests when SQL behavior matters;
- no business workflow assertions beyond persistence contract.

Boot tests:

- smoke-level only;
- verify wiring does not panic and routes are registered.

## What Not To Build

Do not rebuild these patterns:

- runtime module registry;
- module manifest;
- generic provider or capability system;
- framework-independent router abstraction;
- generated repository code before there is stable schema pressure;
- global shared `utils` package that mixes framework and pure helpers;
- business packages importing config directly;
- DB models used as domain entities.

## Migration From Current Skeleton

The cleaned skeleton has already started this migration:

- `internal/boot` owns startup wiring.
- Server startup uses `Run(ctx) error`.
- Echo response helpers live in `internal/server/httpresp`.
- Request metadata lives in framework-free `internal/requestctx`.
- Import boundaries are enforced by `internal/archtest`.
- The old global `internal/utils` package has been removed.
- Framework-free errors live in `internal/apperr`, while HTTP status mapping lives in `internal/server/httpresp`.
- The old mixed-responsibility `internal/code` package has been removed.
- The old unused generic cache, metrics, rate-limit, validator, database, log, and demo user packages have been removed.
- Config loading is static and explicit through `configs.Load(path)`; the old config center, hot reload, merge, and store abstractions have been removed.

The next business feature should use the vertical module shape. Do not add a fake sample business module just to demonstrate the layout.

## Import Guard Target

The repository enforces these checks through `internal/archtest`:

```text
internal/*/domain must not import:
  github.com/labstack/echo
  gorm.io
  github.com/redis/go-redis
  internal/configs
  internal/log
  internal/server
  internal/server/httpresp
  internal/code

internal/*/usecase must not import:
  github.com/labstack/echo
  gorm.io
  github.com/redis/go-redis
  internal/configs
  internal/log
  internal/server
  internal/server/httpresp
  internal/code
  internal/*/http
  internal/*/mysql

internal/*/http must not import:
  gorm.io
  github.com/redis/go-redis
  internal/configs
  internal/log
  internal/code
  internal/*/domain
  internal/*/mysql

internal/*/mysql must not import:
  github.com/labstack/echo
  internal/configs
  internal/log
  internal/server
  internal/server/httpresp
  internal/code
  internal/*/http

internal/server and internal/configs must not import:
  internal/<business>
```

The guard matters because package layout is only useful when the dependency rule is mechanically protected.
