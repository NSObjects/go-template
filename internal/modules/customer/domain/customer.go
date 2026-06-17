package domain

import (
	"errors"
	"net/mail"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	// StatusActive marks a customer that can place orders.
	StatusActive = "active"
	// StatusDisabled marks a customer blocked from new business actions.
	StatusDisabled = "disabled"

	maxNameLength = 120
)

var (
	// ErrInvalidCustomerID reports a non-positive persisted customer id.
	ErrInvalidCustomerID = errors.New("invalid customer id")
	// ErrInvalidName reports a blank or overlong customer name.
	ErrInvalidName = errors.New("invalid customer name")
	// ErrInvalidEmail reports an empty or unparsable email address.
	ErrInvalidEmail = errors.New("invalid customer email")
	// ErrInvalidStatus reports a status outside the customer lifecycle.
	ErrInvalidStatus = errors.New("invalid customer status")
)

// Customer is the domain customer aggregate.
type Customer struct {
	id        int64
	name      string
	email     string
	status    string
	createdAt time.Time
	updatedAt time.Time
}

// New creates an active customer and normalizes name and email.
func New(name, email string) (Customer, error) {
	name, err := cleanName(name)
	if err != nil {
		return Customer{}, err
	}
	email, err = cleanEmail(email)
	if err != nil {
		return Customer{}, err
	}
	return Customer{name: name, email: email, status: StatusActive}, nil
}

// Restore rebuilds a persisted customer while preserving timestamps.
func Restore(id int64, name, email, status string, createdAt, updatedAt time.Time) (Customer, error) {
	if id <= 0 {
		return Customer{}, ErrInvalidCustomerID
	}
	customer, err := New(name, email)
	if err != nil {
		return Customer{}, err
	}
	if !validStatus(status) {
		return Customer{}, ErrInvalidStatus
	}
	customer.id = id
	customer.status = status
	customer.createdAt = createdAt
	customer.updatedAt = updatedAt
	return customer, nil
}

// Update returns a customer with validated mutable fields changed.
func (c Customer) Update(name, email, status string) (Customer, error) {
	name, err := cleanName(name)
	if err != nil {
		return Customer{}, err
	}
	email, err = cleanEmail(email)
	if err != nil {
		return Customer{}, err
	}
	if !validStatus(status) {
		return Customer{}, ErrInvalidStatus
	}
	c.name = name
	c.email = email
	c.status = status
	return c, nil
}

// ID returns the persisted customer id.
func (c Customer) ID() int64 {
	return c.id
}

// Name returns the normalized customer name.
func (c Customer) Name() string {
	return c.name
}

// Email returns the normalized customer email.
func (c Customer) Email() string {
	return c.email
}

// Status returns the customer lifecycle status.
func (c Customer) Status() string {
	return c.status
}

// CreatedAt returns the persisted creation time.
func (c Customer) CreatedAt() time.Time {
	return c.createdAt
}

// UpdatedAt returns the persisted last update time.
func (c Customer) UpdatedAt() time.Time {
	return c.updatedAt
}

func cleanName(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" || utf8.RuneCountInString(value) > maxNameLength {
		return "", ErrInvalidName
	}
	return value, nil
}

func cleanEmail(value string) (string, error) {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" {
		return "", ErrInvalidEmail
	}
	if _, err := mail.ParseAddress(value); err != nil {
		return "", ErrInvalidEmail
	}
	return value, nil
}

func validStatus(status string) bool {
	switch status {
	case StatusActive, StatusDisabled:
		return true
	default:
		return false
	}
}
