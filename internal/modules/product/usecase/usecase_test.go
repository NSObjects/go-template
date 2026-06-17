package usecase

import (
	"context"
	"errors"
	"math"
	"testing"
	"time"

	"github.com/NSObjects/go-template/internal/modules/product/domain"
	"github.com/NSObjects/go-template/internal/platform/apperr"
)

func TestProductUsecaseCreatesAndListsProducts(t *testing.T) {
	uc := New(newStore())

	created, err := uc.Create(context.Background(), CreateInput{SKU: "sku-1", Name: "Starter", PriceCents: 1999})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if created.SKU != "SKU-1" {
		t.Fatalf("Create() SKU = %q, want normalized SKU-1", created.SKU)
	}

	list, err := uc.List(context.Background(), ListInput{Query: "starter", ActiveOnly: true})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if list.Total != 1 || len(list.Items) != 1 {
		t.Fatalf("List() = total %d items %d, want 1 item", list.Total, len(list.Items))
	}
}

func TestProductUsecaseRejectsOverflowingPagination(t *testing.T) {
	uc := New(newStore())

	_, err := uc.List(context.Background(), ListInput{Page: math.MaxInt, PageSize: 100})
	if err == nil {
		t.Fatal("List() error = nil, want pagination error")
	}
	var appErr *apperr.Error
	if !errors.As(err, &appErr) || appErr.Code() != apperr.ErrBadRequest {
		t.Fatalf("List() error = %v, want bad request app error", err)
	}
}

type store struct {
	nextID   int64
	products []domain.Product
}

func newStore() *store {
	return &store{nextID: 1}
}

func (s *store) Create(_ context.Context, product domain.Product) (domain.Product, error) {
	now := time.Now().UTC()
	created, err := domain.Restore(s.nextID, product.SKU(), product.Name(), product.PriceCents(), product.Active(), now, now)
	if err != nil {
		return domain.Product{}, err
	}
	s.nextID++
	s.products = append(s.products, created)
	return created, nil
}

func (s *store) FindByID(_ context.Context, id int64) (domain.Product, error) {
	for _, product := range s.products {
		if product.ID() == id {
			return product, nil
		}
	}
	return domain.Product{}, apperr.NewNotFound("product")
}

func (s *store) List(_ context.Context, filter ListFilter) ([]domain.Product, int, error) {
	out := make([]domain.Product, 0, len(s.products))
	for _, product := range s.products {
		if filter.ActiveOnly && !product.Active() {
			continue
		}
		if filter.Query != "" && product.Name() != "Starter" {
			continue
		}
		out = append(out, product)
	}
	return out, len(out), nil
}

func (s *store) Update(_ context.Context, product domain.Product) (domain.Product, error) {
	for i, existing := range s.products {
		if existing.ID() == product.ID() {
			s.products[i] = product
			return product, nil
		}
	}
	return domain.Product{}, apperr.NewNotFound("product")
}
