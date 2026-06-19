# Boot

`internal/boot` is the framework composition root. It owns process startup,
configured infrastructure resources, `samber/do` dependency injection, and
business route mounting.

Business modules should declare what they provide and what they mount:

```go
func orderModule() boot.Module {
	return boot.NewModule("order",
		boot.Provide(newOrderStore),
		boot.Provide(newOrderUsecase),
		boot.Provide(newOrderHandler),
		boot.Route(orderhttp.Register),
	)
}
```

Provider functions are ordinary `do.Provider[T]` functions:

```go
func newOrderUsecase(i do.Injector) (*orderusecase.Usecase, error) {
	store, err := do.Invoke[*ordermemory.Store](i)
	if err != nil {
		return nil, err
	}
	return orderusecase.New(store), nil
}
```

When a module has both local and real storage adapters, keep the decision in
boot. `product` uses `internal/modules/product/adapters/mysql` when
`mysql.enabled=true` and falls back to the memory store otherwise.

Business code lives under `internal/modules/<module>`. Platform runtime code
lives under `internal/platform`. Keep business adapters under
`internal/modules/<module>/adapters/<adapter>` so the `internal` root does not
grow one directory per business concern.

Boot installs configured infrastructure resources once. Business adapters reuse
those configured infrastructure resources through usecase-owned outbound
interfaces instead of repeating MySQL, Redis, MongoDB, tracing, logging,
readiness, and shutdown code inside each module.

Keep this package concrete. It is allowed to import adapters, infrastructure,
server, and configs so usecase and domain packages stay clean.
