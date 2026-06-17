package domain

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	maxSKULength  = 64
	maxNameLength = 160
)

var (
	// ErrInvalidProductID reports a non-positive persisted product id.
	ErrInvalidProductID = errors.New("invalid product id")
	// ErrInvalidSKU reports a blank or overlong product SKU.
	ErrInvalidSKU = errors.New("invalid product sku")
	// ErrInvalidName reports a blank or overlong product name.
	ErrInvalidName = errors.New("invalid product name")
	// ErrInvalidPrice reports a negative product price.
	ErrInvalidPrice = errors.New("invalid product price")
)

// Product is the domain product aggregate.
type Product struct {
	id         int64
	sku        string
	name       string
	priceCents int64
	active     bool
	createdAt  time.Time
	updatedAt  time.Time
}

// New creates an active product and normalizes SKU and name.
func New(sku, name string, priceCents int64) (Product, error) {
	sku, err := cleanSKU(sku)
	if err != nil {
		return Product{}, err
	}
	name, err = cleanName(name)
	if err != nil {
		return Product{}, err
	}
	if priceCents < 0 {
		return Product{}, ErrInvalidPrice
	}
	return Product{sku: sku, name: name, priceCents: priceCents, active: true}, nil
}

// Restore rebuilds a persisted product while preserving timestamps.
func Restore(id int64, sku, name string, priceCents int64, active bool, createdAt, updatedAt time.Time) (Product, error) {
	if id <= 0 {
		return Product{}, ErrInvalidProductID
	}
	product, err := New(sku, name, priceCents)
	if err != nil {
		return Product{}, err
	}
	product.id = id
	product.active = active
	product.createdAt = createdAt
	product.updatedAt = updatedAt
	return product, nil
}

// Update returns a product with validated mutable fields changed.
func (p Product) Update(name string, priceCents int64, active bool) (Product, error) {
	name, err := cleanName(name)
	if err != nil {
		return Product{}, err
	}
	if priceCents < 0 {
		return Product{}, ErrInvalidPrice
	}
	p.name = name
	p.priceCents = priceCents
	p.active = active
	return p, nil
}

// ID returns the persisted product id.
func (p Product) ID() int64 {
	return p.id
}

// SKU returns the normalized product SKU.
func (p Product) SKU() string {
	return p.sku
}

// Name returns the normalized product name.
func (p Product) Name() string {
	return p.name
}

// PriceCents returns the product price in cents.
func (p Product) PriceCents() int64 {
	return p.priceCents
}

// Active reports whether the product can be ordered.
func (p Product) Active() bool {
	return p.active
}

// CreatedAt returns the persisted creation time.
func (p Product) CreatedAt() time.Time {
	return p.createdAt
}

// UpdatedAt returns the persisted last update time.
func (p Product) UpdatedAt() time.Time {
	return p.updatedAt
}

func cleanSKU(value string) (string, error) {
	value = strings.ToUpper(strings.TrimSpace(value))
	if value == "" || utf8.RuneCountInString(value) > maxSKULength {
		return "", ErrInvalidSKU
	}
	return value, nil
}

func cleanName(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" || utf8.RuneCountInString(value) > maxNameLength {
		return "", ErrInvalidName
	}
	return value, nil
}
