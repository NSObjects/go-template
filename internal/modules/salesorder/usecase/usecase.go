package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/NSObjects/go-template/internal/modules/salesorder/domain"
	"github.com/NSObjects/go-template/internal/platform/apperr"
	"github.com/NSObjects/go-template/internal/platform/pagination"
)

const (
	defaultPageSize = 20
	maxPageSize     = 100
)

// Store persists sales orders for the sales order usecase.
type Store interface {
	Create(context.Context, domain.Order) (domain.Order, error)
	FindByID(context.Context, int64) (domain.Order, error)
	List(context.Context, ListFilter) ([]domain.Order, int, error)
}

// CustomerLookup checks customer existence without importing the customer module.
type CustomerLookup interface {
	CustomerExists(context.Context, int64) (bool, error)
}

// ProductSnapshot is the product state required before an order item can be accepted.
type ProductSnapshot struct {
	Exists bool
	Active bool
}

// ProductLookup supplies product state without coupling sales orders to product storage.
type ProductLookup interface {
	FindProduct(context.Context, int64) (ProductSnapshot, error)
}

// Usecase coordinates sales order application workflows.
type Usecase struct {
	store     Store
	customers CustomerLookup
	products  ProductLookup
}

// New creates a sales order usecase with its required collaborators.
func New(store Store, customers CustomerLookup, products ProductLookup) *Usecase {
	return &Usecase{store: store, customers: customers, products: products}
}

// CreateInput carries order creation fields from an adapter.
type CreateInput struct {
	CustomerID int64
	Items      []CreateItemInput
}

// CreateItemInput carries one order item from an adapter.
type CreateItemInput struct {
	ProductID      int64
	Quantity       int
	UnitPriceCents int64
}

// ListInput carries sales order list filters from an adapter.
type ListInput struct {
	Page       int
	PageSize   int
	CustomerID int64
}

// ListFilter is the validated store-facing sales order list filter.
type ListFilter struct {
	Offset     int
	Limit      int
	CustomerID int64
	Page       int
	PageSize   int
}

// Order is the adapter-facing sales order DTO returned by the usecase.
type Order struct {
	ID              int64     `json:"id"`
	CustomerID      int64     `json:"customer_id"`
	Status          string    `json:"status"`
	Items           []Item    `json:"items"`
	TotalPriceCents int64     `json:"total_price_cents"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Item is the adapter-facing sales order item DTO.
type Item struct {
	ProductID      int64 `json:"product_id"`
	Quantity       int   `json:"quantity"`
	UnitPriceCents int64 `json:"unit_price_cents"`
	LineTotalCents int64 `json:"line_total_cents"`
}

// ListOutput is a paginated sales order list result.
type ListOutput struct {
	Items    []Order
	Page     int
	PageSize int
	Total    int
}

// Create validates references, builds an order, and persists it.
func (u *Usecase) Create(ctx context.Context, input CreateInput) (Order, error) {
	if err := u.ready(); err != nil {
		return Order{}, err
	}
	if input.CustomerID <= 0 {
		return Order{}, mapDomainError(domain.ErrInvalidCustomerID)
	}
	if len(input.Items) == 0 {
		return Order{}, mapDomainError(domain.ErrEmptyItems)
	}
	items := make([]domain.Item, 0, len(input.Items))
	for _, inputItem := range input.Items {
		item, err := domain.NewItem(inputItem.ProductID, inputItem.Quantity, inputItem.UnitPriceCents)
		if err != nil {
			return Order{}, mapDomainError(err)
		}
		if err := u.validateProduct(ctx, item.ProductID()); err != nil {
			return Order{}, err
		}
		items = append(items, item)
	}
	if err := u.validateCustomer(ctx, input.CustomerID); err != nil {
		return Order{}, err
	}
	order, err := domain.New(input.CustomerID, items)
	if err != nil {
		return Order{}, mapDomainError(err)
	}
	created, err := u.store.Create(ctx, order)
	if err != nil {
		return Order{}, err
	}
	return fromDomain(created), nil
}

// Get returns one sales order by id.
func (u *Usecase) Get(ctx context.Context, id int64) (Order, error) {
	if err := u.ready(); err != nil {
		return Order{}, err
	}
	if id <= 0 {
		return Order{}, apperr.NewBadRequest("invalid order id")
	}
	order, err := u.store.FindByID(ctx, id)
	if err != nil {
		return Order{}, err
	}
	return fromDomain(order), nil
}

// List returns sales orders matching the requested filters.
func (u *Usecase) List(ctx context.Context, input ListInput) (ListOutput, error) {
	if err := u.ready(); err != nil {
		return ListOutput{}, err
	}
	filter, err := normalizeListInput(input)
	if err != nil {
		return ListOutput{}, err
	}
	orders, total, err := u.store.List(ctx, filter)
	if err != nil {
		return ListOutput{}, err
	}
	items := make([]Order, 0, len(orders))
	for _, order := range orders {
		items = append(items, fromDomain(order))
	}
	return ListOutput{Items: items, Page: filter.Page, PageSize: filter.PageSize, Total: total}, nil
}

func (u *Usecase) ready() error {
	if u == nil || u.store == nil || u.customers == nil || u.products == nil {
		return apperr.New(apperr.ErrInternalServer, "sales order dependencies are not configured")
	}
	return nil
}

func (u *Usecase) validateCustomer(ctx context.Context, customerID int64) error {
	exists, err := u.customers.CustomerExists(ctx, customerID)
	if err != nil {
		return err
	}
	if !exists {
		return apperr.NewNotFound("customer")
	}
	return nil
}

func (u *Usecase) validateProduct(ctx context.Context, productID int64) error {
	product, err := u.products.FindProduct(ctx, productID)
	if err != nil {
		return err
	}
	if !product.Exists {
		return apperr.NewNotFound("product")
	}
	if !product.Active {
		return apperr.NewConflict("product is not available")
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
	if input.CustomerID < 0 {
		return ListFilter{}, apperr.NewBadRequest("invalid customer id")
	}
	return ListFilter{
		Offset:     window.Offset,
		Limit:      window.Limit,
		CustomerID: input.CustomerID,
		Page:       window.Page,
		PageSize:   window.PageSize,
	}, nil
}

func fromDomain(order domain.Order) Order {
	items := order.Items()
	out := make([]Item, 0, len(items))
	for _, item := range items {
		out = append(out, Item{
			ProductID:      item.ProductID(),
			Quantity:       item.Quantity(),
			UnitPriceCents: item.UnitPriceCents(),
			LineTotalCents: item.LineTotalCents(),
		})
	}
	return Order{
		ID:              order.ID(),
		CustomerID:      order.CustomerID(),
		Status:          order.Status(),
		Items:           out,
		TotalPriceCents: order.TotalPriceCents(),
		CreatedAt:       order.CreatedAt(),
		UpdatedAt:       order.UpdatedAt(),
	}
}

func mapDomainError(err error) error {
	switch {
	case errors.Is(err, domain.ErrInvalidOrderID):
		return apperr.NewBadRequest("invalid order id")
	case errors.Is(err, domain.ErrInvalidCustomerID):
		return apperr.NewBadRequest("invalid customer id")
	case errors.Is(err, domain.ErrInvalidProductID):
		return apperr.NewBadRequest("invalid product id")
	case errors.Is(err, domain.ErrInvalidQuantity):
		return apperr.NewBadRequest("invalid quantity")
	case errors.Is(err, domain.ErrInvalidUnitPrice):
		return apperr.NewBadRequest("invalid unit price")
	case errors.Is(err, domain.ErrInvalidTotal):
		return apperr.NewBadRequest("invalid order total")
	case errors.Is(err, domain.ErrEmptyItems):
		return apperr.NewBadRequest("empty order items")
	case errors.Is(err, domain.ErrInvalidStatus):
		return apperr.NewBadRequest("invalid order status")
	default:
		return err
	}
}
