package boot

import (
	"context"

	"github.com/samber/do/v2"

	customermemory "github.com/NSObjects/go-template/internal/modules/customer/adapters/memory"
	customerhttp "github.com/NSObjects/go-template/internal/modules/customer/http"
	customerusecase "github.com/NSObjects/go-template/internal/modules/customer/usecase"
	"github.com/NSObjects/go-template/internal/platform/apperr"
)

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
