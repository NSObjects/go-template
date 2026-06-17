package boot

import (
	"context"

	"github.com/samber/do/v2"

	productmemory "github.com/NSObjects/go-template/internal/modules/product/adapters/memory"
	producthttp "github.com/NSObjects/go-template/internal/modules/product/http"
	productusecase "github.com/NSObjects/go-template/internal/modules/product/usecase"
	salesorderusecase "github.com/NSObjects/go-template/internal/modules/salesorder/usecase"
	"github.com/NSObjects/go-template/internal/platform/apperr"
)

func productModule() Module {
	return NewModule("product",
		Provide(newProductStore),
		Provide(newProductUsecase),
		Provide(newProductHandler),
		Route(producthttp.Register),
	)
}

func newProductStore(do.Injector) (*productmemory.Store, error) {
	return productmemory.NewStore(), nil
}

func newProductUsecase(i do.Injector) (*productusecase.Usecase, error) {
	store, err := do.Invoke[*productmemory.Store](i)
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
