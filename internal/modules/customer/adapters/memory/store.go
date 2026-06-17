package customermemory

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/NSObjects/go-template/internal/modules/customer/domain"
	"github.com/NSObjects/go-template/internal/modules/customer/usecase"
	"github.com/NSObjects/go-template/internal/platform/apperr"
	"github.com/NSObjects/go-template/internal/platform/pagination"
)

// Store is an in-memory customer store for local development and tests.
type Store struct {
	mu        sync.RWMutex
	nextID    int64
	customers map[int64]domain.Customer
}

// NewStore creates an empty customer memory store.
func NewStore() *Store {
	return &Store{nextID: 1, customers: map[int64]domain.Customer{}}
}

// Create assigns an id and stores a customer.
func (s *Store) Create(_ context.Context, customer domain.Customer) (domain.Customer, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	created, err := domain.Restore(s.nextID, customer.Name(), customer.Email(), customer.Status(), now, now)
	if err != nil {
		return domain.Customer{}, err
	}
	s.customers[created.ID()] = created
	s.nextID++
	return created, nil
}

// FindByID returns one stored customer.
func (s *Store) FindByID(_ context.Context, id int64) (domain.Customer, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	customer, ok := s.customers[id]
	if !ok {
		return domain.Customer{}, apperr.NewNotFound("customer")
	}
	return customer, nil
}

// List returns customers matching the validated filter.
func (s *Store) List(_ context.Context, filter usecase.ListFilter) ([]domain.Customer, int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	all := make([]domain.Customer, 0, len(s.customers))
	query := strings.ToLower(filter.Query)
	for _, customer := range s.customers {
		if filter.Status != "" && customer.Status() != filter.Status {
			continue
		}
		if query != "" && !strings.Contains(strings.ToLower(customer.Name()+" "+customer.Email()), query) {
			continue
		}
		all = append(all, customer)
	}
	sort.Slice(all, func(i, j int) bool { return all[i].ID() > all[j].ID() })
	total := len(all)
	start, end, ok, err := pagination.Bounds(total, filter.Offset, filter.Limit)
	if err != nil {
		return nil, total, apperr.NewBadRequest("invalid pagination")
	}
	if !ok {
		return nil, total, nil
	}
	return append([]domain.Customer(nil), all[start:end]...), total, nil
}

// Update replaces a stored customer while preserving creation time.
func (s *Store) Update(_ context.Context, customer domain.Customer) (domain.Customer, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, ok := s.customers[customer.ID()]
	if !ok {
		return domain.Customer{}, apperr.NewNotFound("customer")
	}
	updated, err := domain.Restore(customer.ID(), customer.Name(), customer.Email(), customer.Status(), existing.CreatedAt(), time.Now().UTC())
	if err != nil {
		return domain.Customer{}, err
	}
	s.customers[updated.ID()] = updated
	return updated, nil
}
