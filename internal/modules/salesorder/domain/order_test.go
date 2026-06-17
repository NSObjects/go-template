package domain

import (
	"errors"
	"math"
	"testing"
)

func TestNewItemRejectsOverflowingLineTotal(t *testing.T) {
	_, err := NewItem(1, 2, math.MaxInt64)
	if !errors.Is(err, ErrInvalidTotal) {
		t.Fatalf("NewItem() error = %v, want ErrInvalidTotal", err)
	}
}

func TestNewRejectsOverflowingOrderTotal(t *testing.T) {
	first, err := NewItem(1, 1, math.MaxInt64)
	if err != nil {
		t.Fatalf("NewItem(first) error = %v", err)
	}
	second, err := NewItem(2, 1, 1)
	if err != nil {
		t.Fatalf("NewItem(second) error = %v", err)
	}

	_, err = New(7, []Item{first, second})
	if !errors.Is(err, ErrInvalidTotal) {
		t.Fatalf("New() error = %v, want ErrInvalidTotal", err)
	}
}
