package domain

import (
	"errors"
	"math"
	"time"
)

// StatusCreated marks an order accepted by the system but not yet advanced.
const StatusCreated = "created"

var (
	// ErrInvalidOrderID reports a non-positive persisted order id.
	ErrInvalidOrderID = errors.New("invalid order id")
	// ErrInvalidCustomerID reports a non-positive customer id.
	ErrInvalidCustomerID = errors.New("invalid customer id")
	// ErrInvalidProductID reports a non-positive product id.
	ErrInvalidProductID = errors.New("invalid product id")
	// ErrInvalidQuantity reports a non-positive item quantity.
	ErrInvalidQuantity = errors.New("invalid quantity")
	// ErrInvalidUnitPrice reports a negative item unit price.
	ErrInvalidUnitPrice = errors.New("invalid unit price")
	// ErrInvalidTotal reports an overflowing line or order total.
	ErrInvalidTotal = errors.New("invalid order total")
	// ErrEmptyItems reports an order without items.
	ErrEmptyItems = errors.New("empty order items")
	// ErrInvalidStatus reports an order status outside the domain lifecycle.
	ErrInvalidStatus = errors.New("invalid order status")
)

// Item is an immutable sales order line item.
type Item struct {
	productID      int64
	quantity       int
	unitPriceCents int64
}

// NewItem creates an order item and rejects price arithmetic overflow.
func NewItem(productID int64, quantity int, unitPriceCents int64) (Item, error) {
	if productID <= 0 {
		return Item{}, ErrInvalidProductID
	}
	if quantity <= 0 {
		return Item{}, ErrInvalidQuantity
	}
	if unitPriceCents < 0 {
		return Item{}, ErrInvalidUnitPrice
	}
	if unitPriceCents > 0 && int64(quantity) > math.MaxInt64/unitPriceCents {
		return Item{}, ErrInvalidTotal
	}
	return Item{productID: productID, quantity: quantity, unitPriceCents: unitPriceCents}, nil
}

// ProductID returns the ordered product id.
func (i Item) ProductID() int64 {
	return i.productID
}

// Quantity returns the ordered unit count.
func (i Item) Quantity() int {
	return i.quantity
}

// UnitPriceCents returns the item unit price in cents.
func (i Item) UnitPriceCents() int64 {
	return i.unitPriceCents
}

// LineTotalCents returns quantity multiplied by unit price in cents.
func (i Item) LineTotalCents() int64 {
	return int64(i.quantity) * i.unitPriceCents
}

// Order is the domain sales order aggregate.
type Order struct {
	id              int64
	customerID      int64
	status          string
	items           []Item
	totalPriceCents int64
	createdAt       time.Time
	updatedAt       time.Time
}

// New creates a sales order and computes a checked total.
func New(customerID int64, items []Item) (Order, error) {
	if customerID <= 0 {
		return Order{}, ErrInvalidCustomerID
	}
	if len(items) == 0 {
		return Order{}, ErrEmptyItems
	}
	copied := append([]Item(nil), items...)
	totalPriceCents, err := total(copied)
	if err != nil {
		return Order{}, err
	}
	return Order{
		customerID:      customerID,
		status:          StatusCreated,
		items:           copied,
		totalPriceCents: totalPriceCents,
	}, nil
}

// Restore rebuilds a persisted sales order while preserving timestamps.
func Restore(id, customerID int64, status string, items []Item, createdAt, updatedAt time.Time) (Order, error) {
	if id <= 0 {
		return Order{}, ErrInvalidOrderID
	}
	if status != StatusCreated {
		return Order{}, ErrInvalidStatus
	}
	order, err := New(customerID, items)
	if err != nil {
		return Order{}, err
	}
	order.id = id
	order.status = status
	order.createdAt = createdAt
	order.updatedAt = updatedAt
	return order, nil
}

// ID returns the persisted order id.
func (o Order) ID() int64 {
	return o.id
}

// CustomerID returns the customer that owns the order.
func (o Order) CustomerID() int64 {
	return o.customerID
}

// Status returns the order lifecycle status.
func (o Order) Status() string {
	return o.status
}

// Items returns a copy of the order line items.
func (o Order) Items() []Item {
	return append([]Item(nil), o.items...)
}

// TotalPriceCents returns the checked order total in cents.
func (o Order) TotalPriceCents() int64 {
	return o.totalPriceCents
}

// CreatedAt returns the persisted creation time.
func (o Order) CreatedAt() time.Time {
	return o.createdAt
}

// UpdatedAt returns the persisted last update time.
func (o Order) UpdatedAt() time.Time {
	return o.updatedAt
}

func total(items []Item) (int64, error) {
	var sum int64
	for _, item := range items {
		lineTotal := item.LineTotalCents()
		if math.MaxInt64-sum < lineTotal {
			return 0, ErrInvalidTotal
		}
		sum += lineTotal
	}
	return sum, nil
}
