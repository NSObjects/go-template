package salesordermemory

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/NSObjects/go-template/internal/modules/salesorder/domain"
	"github.com/NSObjects/go-template/internal/modules/salesorder/usecase"
	"github.com/NSObjects/go-template/internal/platform/apperr"
	"github.com/NSObjects/go-template/internal/platform/pagination"
)

// Store is an in-memory sales order store for local development and tests.
type Store struct {
	mu     sync.RWMutex
	nextID int64
	orders map[int64]domain.Order
}

// NewStore creates an empty sales order memory store.
func NewStore() *Store {
	return &Store{nextID: 1, orders: map[int64]domain.Order{}}
}

// Create assigns an id and stores a sales order.
func (s *Store) Create(_ context.Context, order domain.Order) (domain.Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	created, err := domain.Restore(s.nextID, order.CustomerID(), order.Status(), order.Items(), now, now)
	if err != nil {
		return domain.Order{}, err
	}
	s.orders[created.ID()] = created
	s.nextID++
	return created, nil
}

// FindByID returns one stored sales order.
func (s *Store) FindByID(_ context.Context, id int64) (domain.Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[id]
	if !ok {
		return domain.Order{}, apperr.NewNotFound("sales order")
	}
	return order, nil
}

// List returns sales orders matching the validated filter.
func (s *Store) List(_ context.Context, filter usecase.ListFilter) ([]domain.Order, int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	all := make([]domain.Order, 0, len(s.orders))
	for _, order := range s.orders {
		if filter.CustomerID > 0 && order.CustomerID() != filter.CustomerID {
			continue
		}
		all = append(all, order)
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
	return append([]domain.Order(nil), all[start:end]...), total, nil
}
