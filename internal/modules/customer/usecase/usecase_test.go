package usecase

import (
	"context"
	"errors"
	"math"
	"testing"
	"time"

	"github.com/NSObjects/go-template/internal/modules/customer/domain"
	"github.com/NSObjects/go-template/internal/platform/apperr"
)

func TestCustomerUsecaseCreatesAndListsCustomers(t *testing.T) {
	uc := New(newStore())

	created, err := uc.Create(context.Background(), CreateInput{Name: "Acme", Email: "Owner@Example.com"})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if created.ID == 0 {
		t.Fatal("Create() ID = 0, want assigned id")
	}
	if created.Email != "owner@example.com" {
		t.Fatalf("Create() Email = %q, want normalized email", created.Email)
	}

	list, err := uc.List(context.Background(), ListInput{Query: "acme"})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if list.Total != 1 || len(list.Items) != 1 {
		t.Fatalf("List() = total %d items %d, want 1 item", list.Total, len(list.Items))
	}
}

func TestCustomerUsecaseRejectsOverflowingPagination(t *testing.T) {
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
	nextID    int64
	customers []domain.Customer
}

func newStore() *store {
	return &store{nextID: 1}
}

func (s *store) Create(_ context.Context, customer domain.Customer) (domain.Customer, error) {
	now := time.Now().UTC()
	created, err := domain.Restore(s.nextID, customer.Name(), customer.Email(), customer.Status(), now, now)
	if err != nil {
		return domain.Customer{}, err
	}
	s.nextID++
	s.customers = append(s.customers, created)
	return created, nil
}

func (s *store) FindByID(_ context.Context, id int64) (domain.Customer, error) {
	for _, customer := range s.customers {
		if customer.ID() == id {
			return customer, nil
		}
	}
	return domain.Customer{}, apperr.NewNotFound("customer")
}

func (s *store) List(_ context.Context, filter ListFilter) ([]domain.Customer, int, error) {
	out := make([]domain.Customer, 0, len(s.customers))
	for _, customer := range s.customers {
		if filter.Query != "" && customer.Name() != "Acme" {
			continue
		}
		out = append(out, customer)
	}
	return out, len(out), nil
}

func (s *store) Update(_ context.Context, customer domain.Customer) (domain.Customer, error) {
	for i, existing := range s.customers {
		if existing.ID() == customer.ID() {
			s.customers[i] = customer
			return customer, nil
		}
	}
	return domain.Customer{}, apperr.NewNotFound("customer")
}
