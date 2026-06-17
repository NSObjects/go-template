package usecase

import (
	"context"
	"errors"
	"math"
	"testing"
	"time"

	"github.com/NSObjects/go-template/internal/modules/salesorder/domain"
	"github.com/NSObjects/go-template/internal/platform/apperr"
)

func TestSalesOrderUsecaseCreatesAndListsOrders(t *testing.T) {
	uc := New(newStore(), customerLookup{existing: map[int64]bool{7: true}}, productLookup{
		products: map[int64]ProductSnapshot{11: {Exists: true, Active: true}},
	})

	created, err := uc.Create(context.Background(), CreateInput{
		CustomerID: 7,
		Items: []CreateItemInput{{
			ProductID:      11,
			Quantity:       2,
			UnitPriceCents: 500,
		}},
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if created.TotalPriceCents != 1000 {
		t.Fatalf("Create() TotalPriceCents = %d, want 1000", created.TotalPriceCents)
	}

	list, err := uc.List(context.Background(), ListInput{CustomerID: 7})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if list.Total != 1 || len(list.Items) != 1 {
		t.Fatalf("List() = total %d items %d, want 1 item", list.Total, len(list.Items))
	}
}

func TestSalesOrderUsecaseRejectsMissingCustomer(t *testing.T) {
	uc := New(newStore(), customerLookup{}, productLookup{
		products: map[int64]ProductSnapshot{11: {Exists: true, Active: true}},
	})

	_, err := uc.Create(context.Background(), CreateInput{
		CustomerID: 7,
		Items: []CreateItemInput{{
			ProductID:      11,
			Quantity:       2,
			UnitPriceCents: 500,
		}},
	})
	if err == nil {
		t.Fatal("Create() error = nil, want missing customer error")
	}
	var appErr *apperr.Error
	if !errors.As(err, &appErr) || appErr.Code() != apperr.ErrNotFound {
		t.Fatalf("Create() error = %v, want not found app error", err)
	}
}

func TestSalesOrderUsecaseRejectsMissingProduct(t *testing.T) {
	uc := New(newStore(), customerLookup{existing: map[int64]bool{7: true}}, productLookup{})

	_, err := uc.Create(context.Background(), CreateInput{
		CustomerID: 7,
		Items: []CreateItemInput{{
			ProductID:      11,
			Quantity:       2,
			UnitPriceCents: 500,
		}},
	})
	if err == nil {
		t.Fatal("Create() error = nil, want missing product error")
	}
	var appErr *apperr.Error
	if !errors.As(err, &appErr) || appErr.Code() != apperr.ErrNotFound {
		t.Fatalf("Create() error = %v, want not found app error", err)
	}
}

func TestSalesOrderUsecaseRejectsInactiveProduct(t *testing.T) {
	uc := New(newStore(), customerLookup{existing: map[int64]bool{7: true}}, productLookup{
		products: map[int64]ProductSnapshot{11: {Exists: true, Active: false}},
	})

	_, err := uc.Create(context.Background(), CreateInput{
		CustomerID: 7,
		Items: []CreateItemInput{{
			ProductID:      11,
			Quantity:       2,
			UnitPriceCents: 500,
		}},
	})
	if err == nil {
		t.Fatal("Create() error = nil, want inactive product error")
	}
	var appErr *apperr.Error
	if !errors.As(err, &appErr) || appErr.Code() != apperr.ErrConflict {
		t.Fatalf("Create() error = %v, want conflict app error", err)
	}
}

func TestSalesOrderUsecaseRejectsOverflowingPagination(t *testing.T) {
	uc := New(newStore(), customerLookup{}, productLookup{})

	_, err := uc.List(context.Background(), ListInput{Page: math.MaxInt, PageSize: 100})
	if err == nil {
		t.Fatal("List() error = nil, want pagination error")
	}
	var appErr *apperr.Error
	if !errors.As(err, &appErr) || appErr.Code() != apperr.ErrBadRequest {
		t.Fatalf("List() error = %v, want bad request app error", err)
	}
}

func TestSalesOrderUsecaseRejectsOverflowingTotal(t *testing.T) {
	uc := New(newStore(), customerLookup{existing: map[int64]bool{7: true}}, productLookup{
		products: map[int64]ProductSnapshot{
			11: {Exists: true, Active: true},
			12: {Exists: true, Active: true},
		},
	})

	_, err := uc.Create(context.Background(), CreateInput{
		CustomerID: 7,
		Items: []CreateItemInput{
			{ProductID: 11, Quantity: 1, UnitPriceCents: math.MaxInt64},
			{ProductID: 12, Quantity: 1, UnitPriceCents: 1},
		},
	})
	if err == nil {
		t.Fatal("Create() error = nil, want overflowing total error")
	}
	var appErr *apperr.Error
	if !errors.As(err, &appErr) || appErr.Code() != apperr.ErrBadRequest {
		t.Fatalf("Create() error = %v, want bad request app error", err)
	}
}

type store struct {
	nextID int64
	orders []domain.Order
}

func newStore() *store {
	return &store{nextID: 1}
}

func (s *store) Create(_ context.Context, order domain.Order) (domain.Order, error) {
	now := time.Now().UTC()
	created, err := domain.Restore(s.nextID, order.CustomerID(), order.Status(), order.Items(), now, now)
	if err != nil {
		return domain.Order{}, err
	}
	s.nextID++
	s.orders = append(s.orders, created)
	return created, nil
}

func (s *store) FindByID(_ context.Context, id int64) (domain.Order, error) {
	for _, order := range s.orders {
		if order.ID() == id {
			return order, nil
		}
	}
	return domain.Order{}, apperr.NewNotFound("sales order")
}

func (s *store) List(_ context.Context, filter ListFilter) ([]domain.Order, int, error) {
	out := make([]domain.Order, 0, len(s.orders))
	for _, order := range s.orders {
		if filter.CustomerID > 0 && order.CustomerID() != filter.CustomerID {
			continue
		}
		out = append(out, order)
	}
	return out, len(out), nil
}

type customerLookup struct {
	existing map[int64]bool
}

func (l customerLookup) CustomerExists(_ context.Context, id int64) (bool, error) {
	return l.existing[id], nil
}

type productLookup struct {
	products map[int64]ProductSnapshot
}

func (l productLookup) FindProduct(_ context.Context, id int64) (ProductSnapshot, error) {
	return l.products[id], nil
}
