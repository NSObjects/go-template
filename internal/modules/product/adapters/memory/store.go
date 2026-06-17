package productmemory

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/NSObjects/go-template/internal/modules/product/domain"
	"github.com/NSObjects/go-template/internal/modules/product/usecase"
	"github.com/NSObjects/go-template/internal/platform/apperr"
	"github.com/NSObjects/go-template/internal/platform/pagination"
)

// Store is an in-memory product store for local development and tests.
type Store struct {
	mu       sync.RWMutex
	nextID   int64
	products map[int64]domain.Product
	skus     map[string]int64
}

// NewStore creates an empty product memory store.
func NewStore() *Store {
	return &Store{nextID: 1, products: map[int64]domain.Product{}, skus: map[string]int64{}}
}

// Create assigns an id and stores a product.
func (s *Store) Create(_ context.Context, product domain.Product) (domain.Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.skus[product.SKU()]; ok {
		return domain.Product{}, apperr.NewConflict("product sku already exists")
	}
	now := time.Now().UTC()
	created, err := domain.Restore(s.nextID, product.SKU(), product.Name(), product.PriceCents(), product.Active(), now, now)
	if err != nil {
		return domain.Product{}, err
	}
	s.products[created.ID()] = created
	s.skus[created.SKU()] = created.ID()
	s.nextID++
	return created, nil
}

// FindByID returns one stored product.
func (s *Store) FindByID(_ context.Context, id int64) (domain.Product, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	product, ok := s.products[id]
	if !ok {
		return domain.Product{}, apperr.NewNotFound("product")
	}
	return product, nil
}

// List returns products matching the validated filter.
func (s *Store) List(_ context.Context, filter usecase.ListFilter) ([]domain.Product, int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	all := make([]domain.Product, 0, len(s.products))
	query := strings.ToLower(filter.Query)
	for _, product := range s.products {
		if filter.ActiveOnly && !product.Active() {
			continue
		}
		if query != "" && !strings.Contains(strings.ToLower(product.SKU()+" "+product.Name()), query) {
			continue
		}
		all = append(all, product)
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
	return append([]domain.Product(nil), all[start:end]...), total, nil
}

// Update replaces a stored product while preserving SKU and creation time.
func (s *Store) Update(_ context.Context, product domain.Product) (domain.Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, ok := s.products[product.ID()]
	if !ok {
		return domain.Product{}, apperr.NewNotFound("product")
	}
	updated, err := domain.Restore(product.ID(), existing.SKU(), product.Name(), product.PriceCents(), product.Active(), existing.CreatedAt(), time.Now().UTC())
	if err != nil {
		return domain.Product{}, err
	}
	s.products[updated.ID()] = updated
	return updated, nil
}
