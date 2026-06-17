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

func TestCommonRequestErrorCodes(t *testing.T) {
	tests := []struct {
		name string
		got  int
		want int
	}{
		{name: "bad request", got: ErrBadRequest, want: 100400},
		{name: "unauthorized", got: ErrUnauthorized, want: 100401},
		{name: "forbidden", got: ErrForbidden, want: 100403},
		{name: "not found", got: ErrNotFound, want: 100404},
		{name: "method not allowed", got: ErrMethodNotAllowed, want: 100405},
		{name: "conflict", got: ErrConflict, want: 100409},
		{name: "internal", got: ErrInternalServer, want: 100500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("code = %d, want %d", tt.got, tt.want)
			}
		})
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

func TestNewInfoConflict(t *testing.T) {
	err := NewConflict("order already paid")
	info := NewInfo(err)

	if info.Kind != KindConflict {
		t.Fatalf("Kind = %q, want %q", info.Kind, KindConflict)
	}
	if info.Category != CategoryBusiness {
		t.Fatalf("Category = %q, want %q", info.Category, CategoryBusiness)
	}
	if info.Code != ErrConflict {
		t.Fatalf("Code = %d, want %d", info.Code, ErrConflict)
	}
	if info.Message != "order already paid" {
		t.Fatalf("Message = %q, want order already paid", info.Message)
	}
}

func TestNewInfoMethodNotAllowed(t *testing.T) {
	err := NewMethodNotAllowed("")
	info := NewInfo(err)

	if info.Kind != KindMethodNotAllowed {
		t.Fatalf("Kind = %q, want %q", info.Kind, KindMethodNotAllowed)
	}
	if info.Code != ErrMethodNotAllowed {
		t.Fatalf("Code = %d, want %d", info.Code, ErrMethodNotAllowed)
	}
	if info.Message != "Method not allowed" {
		t.Fatalf("Message = %q, want Method not allowed", info.Message)
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
	if err := Wrap(nil, ErrForbidden, "access denied"); err != nil {
		t.Fatalf("Wrap(nil) = %v, want nil", err)
	}
	if err := WrapBadRequest(nil, "invalid id"); err != nil {
		t.Fatalf("WrapBadRequest(nil) = %v, want nil", err)
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

	err = WrapConflict(errors.New("duplicate key"), "resource already exists")
	if err == nil {
		t.Fatal("WrapConflict() error = nil, want app error")
	}
	appErr, ok = Parse(err)
	if !ok {
		t.Fatal("Parse(WrapConflict()) ok = false, want true")
	}
	if appErr.Code() != ErrConflict {
		t.Fatalf("Code = %d, want %d", appErr.Code(), ErrConflict)
	}
}
