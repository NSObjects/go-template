package boot

import (
	"github.com/samber/do/v2"

	customerusecase "github.com/NSObjects/go-template/internal/modules/customer/usecase"
	productusecase "github.com/NSObjects/go-template/internal/modules/product/usecase"
	salesordermemory "github.com/NSObjects/go-template/internal/modules/salesorder/adapters/memory"
	salesorderhttp "github.com/NSObjects/go-template/internal/modules/salesorder/http"
	salesorderusecase "github.com/NSObjects/go-template/internal/modules/salesorder/usecase"
)

func salesOrderModule() Module {
	return NewModule("sales-order",
		Provide(newSalesOrderStore),
		Provide(newSalesOrderUsecase),
		Provide(newSalesOrderHandler),
		Route(salesorderhttp.Register),
	)
}

func newSalesOrderStore(do.Injector) (*salesordermemory.Store, error) {
	return salesordermemory.NewStore(), nil
}

func newSalesOrderUsecase(i do.Injector) (*salesorderusecase.Usecase, error) {
	store, err := do.Invoke[*salesordermemory.Store](i)
	if err != nil {
		return nil, err
	}
	customers, err := do.Invoke[*customerusecase.Usecase](i)
	if err != nil {
		return nil, err
	}
	products, err := do.Invoke[*productusecase.Usecase](i)
	if err != nil {
		return nil, err
	}
	return salesorderusecase.New(
		store,
		salesOrderCustomerLookup{customers: customers},
		salesOrderProductLookup{products: products},
	), nil
}

func newSalesOrderHandler(i do.Injector) (*salesorderhttp.Handler, error) {
	usecase, err := do.Invoke[*salesorderusecase.Usecase](i)
	if err != nil {
		return nil, err
	}
	return salesorderhttp.New(usecase), nil
}
