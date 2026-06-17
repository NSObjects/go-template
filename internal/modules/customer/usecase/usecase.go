package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/NSObjects/go-template/internal/modules/customer/domain"
	"github.com/NSObjects/go-template/internal/platform/apperr"
	"github.com/NSObjects/go-template/internal/platform/pagination"
)

const (
	defaultPageSize = 20
	maxPageSize     = 100
)

// Store persists customers for the customer usecase.
type Store interface {
	Create(context.Context, domain.Customer) (domain.Customer, error)
	FindByID(context.Context, int64) (domain.Customer, error)
	List(context.Context, ListFilter) ([]domain.Customer, int, error)
	Update(context.Context, domain.Customer) (domain.Customer, error)
}

// Usecase coordinates customer application workflows.
type Usecase struct {
	store Store
}

// New creates a customer usecase with its required store.
func New(store Store) *Usecase {
	return &Usecase{store: store}
}

// CreateInput carries customer creation fields from an adapter.
type CreateInput struct {
	Name  string
	Email string
}

// UpdateInput carries optional customer fields for partial updates.
type UpdateInput struct {
	ID     int64
	Name   *string
	Email  *string
	Status *string
}

// ListInput carries customer list filters from an adapter.
type ListInput struct {
	Page     int
	PageSize int
	Query    string
	Status   string
}

// ListFilter is the validated store-facing customer list filter.
type ListFilter struct {
	Offset   int
	Limit    int
	Query    string
	Status   string
	Page     int
	PageSize int
}

// Customer is the adapter-facing customer DTO returned by the usecase.
type Customer struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ListOutput is a paginated customer list result.
type ListOutput struct {
	Items    []Customer
	Page     int
	PageSize int
	Total    int
}

// Create validates and persists a new customer.
func (u *Usecase) Create(ctx context.Context, input CreateInput) (Customer, error) {
	if err := u.ready(); err != nil {
		return Customer{}, err
	}
	customer, err := domain.New(input.Name, input.Email)
	if err != nil {
		return Customer{}, mapDomainError(err)
	}
	created, err := u.store.Create(ctx, customer)
	if err != nil {
		return Customer{}, err
	}
	return fromDomain(created), nil
}

// Get returns one customer by id.
func (u *Usecase) Get(ctx context.Context, id int64) (Customer, error) {
	if err := u.ready(); err != nil {
		return Customer{}, err
	}
	if id <= 0 {
		return Customer{}, apperr.NewBadRequest("invalid customer id")
	}
	customer, err := u.store.FindByID(ctx, id)
	if err != nil {
		return Customer{}, err
	}
	return fromDomain(customer), nil
}

// List returns customers matching the requested filters.
func (u *Usecase) List(ctx context.Context, input ListInput) (ListOutput, error) {
	if err := u.ready(); err != nil {
		return ListOutput{}, err
	}
	filter, err := normalizeListInput(input)
	if err != nil {
		return ListOutput{}, err
	}
	customers, total, err := u.store.List(ctx, filter)
	if err != nil {
		return ListOutput{}, err
	}
	items := make([]Customer, 0, len(customers))
	for _, customer := range customers {
		items = append(items, fromDomain(customer))
	}
	return ListOutput{Items: items, Page: filter.Page, PageSize: filter.PageSize, Total: total}, nil
}

// Update validates and persists mutable customer fields.
func (u *Usecase) Update(ctx context.Context, input UpdateInput) (Customer, error) {
	if err := u.ready(); err != nil {
		return Customer{}, err
	}
	if input.ID <= 0 {
		return Customer{}, apperr.NewBadRequest("invalid customer id")
	}
	existing, err := u.store.FindByID(ctx, input.ID)
	if err != nil {
		return Customer{}, err
	}
	name := existing.Name()
	if input.Name != nil {
		name = *input.Name
	}
	email := existing.Email()
	if input.Email != nil {
		email = *input.Email
	}
	status := existing.Status()
	if input.Status != nil {
		status = *input.Status
	}
	updated, err := existing.Update(name, email, status)
	if err != nil {
		return Customer{}, mapDomainError(err)
	}
	saved, err := u.store.Update(ctx, updated)
	if err != nil {
		return Customer{}, err
	}
	return fromDomain(saved), nil
}

func (u *Usecase) ready() error {
	if u == nil || u.store == nil {
		return apperr.New(apperr.ErrInternalServer, "customer store is not configured")
	}
	return nil
}

func normalizeListInput(input ListInput) (ListFilter, error) {
	status := strings.TrimSpace(input.Status)
	window, err := pagination.Normalize(input.Page, input.PageSize, pagination.Options{
		DefaultPageSize: defaultPageSize,
		MaxPageSize:     maxPageSize,
	})
	if err != nil {
		return ListFilter{}, apperr.NewBadRequest("invalid pagination")
	}
	if status != "" && status != domain.StatusActive && status != domain.StatusDisabled {
		return ListFilter{}, apperr.NewBadRequest("invalid customer status")
	}
	return ListFilter{
		Offset:   window.Offset,
		Limit:    window.Limit,
		Query:    strings.TrimSpace(input.Query),
		Status:   status,
		Page:     window.Page,
		PageSize: window.PageSize,
	}, nil
}

func fromDomain(customer domain.Customer) Customer {
	return Customer{
		ID:        customer.ID(),
		Name:      customer.Name(),
		Email:     customer.Email(),
		Status:    customer.Status(),
		CreatedAt: customer.CreatedAt(),
		UpdatedAt: customer.UpdatedAt(),
	}
}

func mapDomainError(err error) error {
	switch {
	case errors.Is(err, domain.ErrInvalidCustomerID):
		return apperr.NewBadRequest("invalid customer id")
	case errors.Is(err, domain.ErrInvalidName):
		return apperr.NewBadRequest("invalid customer name")
	case errors.Is(err, domain.ErrInvalidEmail):
		return apperr.NewBadRequest("invalid customer email")
	case errors.Is(err, domain.ErrInvalidStatus):
		return apperr.NewBadRequest("invalid customer status")
	default:
		return err
	}
}
