package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/NSObjects/go-template/internal/modules/product/domain"
	"github.com/NSObjects/go-template/internal/platform/apperr"
	"github.com/NSObjects/go-template/internal/platform/pagination"
)

const (
	defaultPageSize = 20
	maxPageSize     = 100
)

// Store persists products for the product usecase.
type Store interface {
	Create(context.Context, domain.Product) (domain.Product, error)
	FindByID(context.Context, int64) (domain.Product, error)
	List(context.Context, ListFilter) ([]domain.Product, int, error)
	Update(context.Context, domain.Product) (domain.Product, error)
}

// Usecase coordinates product application workflows.
type Usecase struct {
	store Store
}

// New creates a product usecase with its required store.
func New(store Store) *Usecase {
	return &Usecase{store: store}
}

// CreateInput carries product creation fields from an adapter.
type CreateInput struct {
	SKU        string
	Name       string
	PriceCents int64
}

// UpdateInput carries optional product fields for partial updates.
type UpdateInput struct {
	ID         int64
	Name       *string
	PriceCents *int64
	Active     *bool
}

// ListInput carries product list filters from an adapter.
type ListInput struct {
	Page       int
	PageSize   int
	Query      string
	ActiveOnly bool
}

// ListFilter is the validated store-facing product list filter.
type ListFilter struct {
	Offset     int
	Limit      int
	Query      string
	ActiveOnly bool
	Page       int
	PageSize   int
}

// Product is the adapter-facing product DTO returned by the usecase.
type Product struct {
	ID         int64     `json:"id"`
	SKU        string    `json:"sku"`
	Name       string    `json:"name"`
	PriceCents int64     `json:"price_cents"`
	Active     bool      `json:"active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// ListOutput is a paginated product list result.
type ListOutput struct {
	Items    []Product
	Page     int
	PageSize int
	Total    int
}

// Create validates and persists a new product.
func (u *Usecase) Create(ctx context.Context, input CreateInput) (Product, error) {
	if err := u.ready(); err != nil {
		return Product{}, err
	}
	product, err := domain.New(input.SKU, input.Name, input.PriceCents)
	if err != nil {
		return Product{}, mapDomainError(err)
	}
	created, err := u.store.Create(ctx, product)
	if err != nil {
		return Product{}, err
	}
	return fromDomain(created), nil
}

// Get returns one product by id.
func (u *Usecase) Get(ctx context.Context, id int64) (Product, error) {
	if err := u.ready(); err != nil {
		return Product{}, err
	}
	if id <= 0 {
		return Product{}, apperr.NewBadRequest("invalid product id")
	}
	product, err := u.store.FindByID(ctx, id)
	if err != nil {
		return Product{}, err
	}
	return fromDomain(product), nil
}

// List returns products matching the requested filters.
func (u *Usecase) List(ctx context.Context, input ListInput) (ListOutput, error) {
	if err := u.ready(); err != nil {
		return ListOutput{}, err
	}
	filter, err := normalizeListInput(input)
	if err != nil {
		return ListOutput{}, err
	}
	products, total, err := u.store.List(ctx, filter)
	if err != nil {
		return ListOutput{}, err
	}
	items := make([]Product, 0, len(products))
	for _, product := range products {
		items = append(items, fromDomain(product))
	}
	return ListOutput{Items: items, Page: filter.Page, PageSize: filter.PageSize, Total: total}, nil
}

// Update validates and persists mutable product fields.
func (u *Usecase) Update(ctx context.Context, input UpdateInput) (Product, error) {
	if err := u.ready(); err != nil {
		return Product{}, err
	}
	if input.ID <= 0 {
		return Product{}, apperr.NewBadRequest("invalid product id")
	}
	existing, err := u.store.FindByID(ctx, input.ID)
	if err != nil {
		return Product{}, err
	}
	name := existing.Name()
	if input.Name != nil {
		name = *input.Name
	}
	priceCents := existing.PriceCents()
	if input.PriceCents != nil {
		priceCents = *input.PriceCents
	}
	active := existing.Active()
	if input.Active != nil {
		active = *input.Active
	}
	updated, err := existing.Update(name, priceCents, active)
	if err != nil {
		return Product{}, mapDomainError(err)
	}
	saved, err := u.store.Update(ctx, updated)
	if err != nil {
		return Product{}, err
	}
	return fromDomain(saved), nil
}

func (u *Usecase) ready() error {
	if u == nil || u.store == nil {
		return apperr.New(apperr.ErrInternalServer, "product store is not configured")
	}
	return nil
}

func normalizeListInput(input ListInput) (ListFilter, error) {
	window, err := pagination.Normalize(input.Page, input.PageSize, pagination.Options{
		DefaultPageSize: defaultPageSize,
		MaxPageSize:     maxPageSize,
	})
	if err != nil {
		return ListFilter{}, apperr.NewBadRequest("invalid pagination")
	}
	return ListFilter{
		Offset:     window.Offset,
		Limit:      window.Limit,
		Query:      strings.TrimSpace(input.Query),
		ActiveOnly: input.ActiveOnly,
		Page:       window.Page,
		PageSize:   window.PageSize,
	}, nil
}

func fromDomain(product domain.Product) Product {
	return Product{
		ID:         product.ID(),
		SKU:        product.SKU(),
		Name:       product.Name(),
		PriceCents: product.PriceCents(),
		Active:     product.Active(),
		CreatedAt:  product.CreatedAt(),
		UpdatedAt:  product.UpdatedAt(),
	}
}

func mapDomainError(err error) error {
	switch {
	case errors.Is(err, domain.ErrInvalidProductID):
		return apperr.NewBadRequest("invalid product id")
	case errors.Is(err, domain.ErrInvalidSKU):
		return apperr.NewBadRequest("invalid product sku")
	case errors.Is(err, domain.ErrInvalidName):
		return apperr.NewBadRequest("invalid product name")
	case errors.Is(err, domain.ErrInvalidPrice):
		return apperr.NewBadRequest("invalid product price")
	default:
		return err
	}
}
