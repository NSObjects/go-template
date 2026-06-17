package pagination

import (
	"errors"
	"math"
	"testing"
)

func TestNormalizeAppliesDefaults(t *testing.T) {
	got, err := Normalize(0, 0, Options{DefaultPageSize: 20, MaxPageSize: 100})
	if err != nil {
		t.Fatalf("Normalize() error = %v", err)
	}
	if got.Page != 1 || got.PageSize != 20 || got.Offset != 0 || got.Limit != 20 {
		t.Fatalf("Normalize() = %#v, want first page with default size", got)
	}
}

func TestNormalizeRejectsOverflowingOffset(t *testing.T) {
	_, err := Normalize(math.MaxInt, 100, Options{DefaultPageSize: 20, MaxPageSize: 100})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("Normalize() error = %v, want ErrInvalid", err)
	}
}

func TestNormalizeRejectsInvalidPolicy(t *testing.T) {
	_, err := Normalize(1, 0, Options{DefaultPageSize: 0, MaxPageSize: 100})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("Normalize() error = %v, want ErrInvalid", err)
	}
}

func TestBoundsClampsOverflowingLimit(t *testing.T) {
	start, end, ok, err := Bounds(5, 4, math.MaxInt)
	if err != nil {
		t.Fatalf("Bounds() error = %v", err)
	}
	if !ok || start != 4 || end != 5 {
		t.Fatalf("Bounds() = %d %d %v, want 4 5 true", start, end, ok)
	}
}

func TestBoundsRejectsInvalidInput(t *testing.T) {
	_, _, _, err := Bounds(5, -1, 10)
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("Bounds() error = %v, want ErrInvalid", err)
	}
}
