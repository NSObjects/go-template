package apperr

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestLookup(t *testing.T) {
	def, ok := Lookup(ErrDatabase)
	if !ok {
		t.Fatal("Lookup(ErrDatabase) ok = false, want true")
	}
	if def.Code != ErrDatabase {
		t.Fatalf("Code = %d, want %d", def.Code, ErrDatabase)
	}
	if def.Kind != KindInternal {
		t.Fatalf("Kind = %q, want %q", def.Kind, KindInternal)
	}
	if def.Category != CategoryDatabase {
		t.Fatalf("Category = %q, want %q", def.Category, CategoryDatabase)
	}
	if def.Message != "Database error" {
		t.Fatalf("Message = %q, want Database error", def.Message)
	}

	if _, ok := Lookup(-1); ok {
		t.Fatal("Lookup(-1) ok = true, want false")
	}
}

func TestNewInfoValidation(t *testing.T) {
	err := NewValidation("email", "invalid format")
	info := NewInfo(err)

	if info.Kind != KindValidation {
		t.Fatalf("Kind = %q, want %q", info.Kind, KindValidation)
	}
	if info.Category != CategoryValidation {
		t.Fatalf("Category = %q, want %q", info.Category, CategoryValidation)
	}
	if info.Code != ErrValidation {
		t.Fatalf("Code = %d, want %d", info.Code, ErrValidation)
	}
	if info.Message != "invalid format" {
		t.Fatalf("Message = %q, want invalid format", info.Message)
	}
	if !strings.Contains(info.Detail, "validation failed for field email") {
		t.Fatalf("Detail = %q, want validation detail", info.Detail)
	}
}

func TestNewInfoInternalHidesDetail(t *testing.T) {
	err := WrapDatabase(errors.New("connection refused"), "query user")
	info := NewInfo(err)

	if info.Kind != KindInternal {
		t.Fatalf("Kind = %q, want %q", info.Kind, KindInternal)
	}
	if info.Category != CategoryDatabase {
		t.Fatalf("Category = %q, want %q", info.Category, CategoryDatabase)
	}
	if info.Code != ErrDatabase {
		t.Fatalf("Code = %d, want %d", info.Code, ErrDatabase)
	}
	if info.Message != "Database error" {
		t.Fatalf("Message = %q, want Database error", info.Message)
	}
	if !strings.Contains(info.Detail, "database query user failed") {
		t.Fatalf("Detail = %q, want operation detail", info.Detail)
	}
}

func TestNewInfoPlainErrorUsesUnknownInternal(t *testing.T) {
	info := NewInfo(errors.New("unexpected"))

	if info.Kind != KindInternal {
		t.Fatalf("Kind = %q, want %q", info.Kind, KindInternal)
	}
	if info.Category != CategorySystem {
		t.Fatalf("Category = %q, want %q", info.Category, CategorySystem)
	}
	if info.Code != ErrUnknown {
		t.Fatalf("Code = %d, want %d", info.Code, ErrUnknown)
	}
	if info.Message != "Internal server error" {
		t.Fatalf("Message = %q, want Internal server error", info.Message)
	}
	if info.Detail == "" {
		t.Fatal("Detail is empty, want diagnostic detail")
	}
}

func TestNewInfoUnregisteredAppErrorFallsBackToUnknown(t *testing.T) {
	err := New(999999, "should not leak")
	info := NewInfo(err)

	if info.Kind != KindInternal {
		t.Fatalf("Kind = %q, want %q", info.Kind, KindInternal)
	}
	if info.Code != ErrUnknown {
		t.Fatalf("Code = %d, want %d", info.Code, ErrUnknown)
	}
	if info.Message != "Internal server error" {
		t.Fatalf("Message = %q, want Internal server error", info.Message)
	}
	if !strings.Contains(info.Detail, "should not leak") {
		t.Fatalf("Detail = %q, want original detail", info.Detail)
	}
}

func TestParse(t *testing.T) {
	t.Run("business error keeps public message and detail", func(t *testing.T) {
		err := WrapBadRequest(errors.New("raw parse failure"), "invalid id")

		appErr, ok := Parse(err)
		if !ok {
			t.Fatal("Parse() ok = false, want true")
		}
		if appErr.Code() != ErrBadRequest {
			t.Fatalf("Code = %d, want %d", appErr.Code(), ErrBadRequest)
		}
		if appErr.Message() != "invalid id" {
			t.Fatalf("Message = %q, want invalid id", appErr.Message())
		}
		if !strings.Contains(appErr.Detail(), "raw parse failure") {
			t.Fatalf("Detail = %q, want cause detail", appErr.Detail())
		}
	})

	t.Run("fmt wrapped app error is parsed", func(t *testing.T) {
		err := fmt.Errorf("outer: %w", NewNotFound("order"))

		appErr, ok := Parse(err)
		if !ok {
			t.Fatal("Parse() ok = false, want true")
		}
		if appErr.Code() != ErrNotFound {
			t.Fatalf("Code = %d, want %d", appErr.Code(), ErrNotFound)
		}
		if appErr.Message() != "order not found" {
			t.Fatalf("Message = %q, want order not found", appErr.Message())
		}
	})
}

func TestWrapHelpers(t *testing.T) {
	if err := WrapDatabase(nil, "query"); err != nil {
		t.Fatalf("WrapDatabase(nil) = %v, want nil", err)
	}

	err := Wrap(errors.New("boom"), ErrForbidden, "access denied")
	if err == nil {
		t.Fatal("Wrap() error = nil, want app error")
	}
	appErr, ok := Parse(err)
	if !ok {
		t.Fatal("Parse(Wrap()) ok = false, want true")
	}
	if appErr.Code() != ErrForbidden {
		t.Fatalf("Code = %d, want %d", appErr.Code(), ErrForbidden)
	}
}
