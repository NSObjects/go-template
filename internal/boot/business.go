package boot

import (
	"context"

	"github.com/samber/do/v2"
	"gorm.io/gorm"

	customermemory "github.com/NSObjects/go-template/internal/modules/customer/adapters/memory"
	customerhttp "github.com/NSObjects/go-template/internal/modules/customer/http"
	customerusecase "github.com/NSObjects/go-template/internal/modules/customer/usecase"
	productmemory "github.com/NSObjects/go-template/internal/modules/product/adapters/memory"
	productmysql "github.com/NSObjects/go-template/internal/modules/product/adapters/mysql"
	producthttp "github.com/NSObjects/go-template/internal/modules/product/http"
	productusecase "github.com/NSObjects/go-template/internal/modules/product/usecase"
	salesordermemory "github.com/NSObjects/go-template/internal/modules/salesorder/adapters/memory"
	salesorderhttp "github.com/NSObjects/go-template/internal/modules/salesorder/http"
	salesorderusecase "github.com/NSObjects/go-template/internal/modules/salesorder/usecase"
	"github.com/NSObjects/go-template/internal/platform/apperr"
	"github.com/NSObjects/go-template/internal/platform/configs"
)

// BusinessModules returns the business modules installed by the default runtime.
func BusinessModules() []Module {
	return []Module{
		customerModule(),
		productModule(),
		salesOrderModule(),
	}
}

func customerModule() Module {
	return NewModule("customer",
		Provide(newCustomerStore),
		Provide(newCustomerUsecase),
		Provide(newCustomerHandler),
		Route(customerhttp.Register),
	)
}

func newCustomerStore(do.Injector) (*customermemory.Store, error) {
	return customermemory.NewStore(), nil
}

func newCustomerUsecase(i do.Injector) (*customerusecase.Usecase, error) {
	store, err := do.Invoke[*customermemory.Store](i)
	if err != nil {
		return nil, err
	}
	return customerusecase.New(store), nil
}

func newCustomerHandler(i do.Injector) (*customerhttp.Handler, error) {
	usecase, err := do.Invoke[*customerusecase.Usecase](i)
	if err != nil {
		return nil, err
	}
	return customerhttp.New(usecase), nil
}

func productModule() Module {
	return NewModule("product",
		Provide(newProductStore),
		Provide(newProductUsecase),
		Provide(newProductHandler),
		Route(producthttp.Register),
	)
}

func newProductStore(i do.Injector) (productusecase.Store, error) {
	cfg, err := do.Invoke[configs.Config](i)
	if err != nil {
		return nil, err
	}
	if !cfg.MySQL.Enabled {
		return productmemory.NewStore(), nil
	}
	db, err := do.Invoke[*gorm.DB](i)
	if err != nil {
		return nil, err
	}
	return productmysql.NewStore(db)
}

func newProductUsecase(i do.Injector) (*productusecase.Usecase, error) {
	store, err := do.Invoke[productusecase.Store](i)
	if err != nil {
		return nil, err
	}
	return productusecase.New(store), nil
}

func newProductHandler(i do.Injector) (*producthttp.Handler, error) {
	usecase, err := do.Invoke[*productusecase.Usecase](i)
	if err != nil {
		return nil, err
	}
	return producthttp.New(usecase), nil
}

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

type salesOrderCustomerLookup struct {
	customers *customerusecase.Usecase
}

func (l salesOrderCustomerLookup) CustomerExists(ctx context.Context, id int64) (bool, error) {
	if _, err := l.customers.Get(ctx, id); err != nil {
		if appErr, ok := apperr.Parse(err); ok && appErr.Code() == apperr.ErrNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

type salesOrderProductLookup struct {
	products *productusecase.Usecase
}

func (l salesOrderProductLookup) FindProduct(ctx context.Context, id int64) (salesorderusecase.ProductSnapshot, error) {
	product, err := l.products.Get(ctx, id)
	if err != nil {
		if appErr, ok := apperr.Parse(err); ok && appErr.Code() == apperr.ErrNotFound {
			return salesorderusecase.ProductSnapshot{Exists: false}, nil
		}
		return salesorderusecase.ProductSnapshot{}, err
	}
	return salesorderusecase.ProductSnapshot{Exists: true, Active: product.Active}, nil
}
