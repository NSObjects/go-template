package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/NSObjects/go-template/internal/code"
	"github.com/NSObjects/go-template/internal/server/middlewares"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func TestBindAndValidateBindsJSON(t *testing.T) {
	var req struct {
		Name string `json:"name" validate:"required"`
	}
	ctx := newRequestTestContext(http.MethodPost, "/", `{"name":"lintao"}`)

	if err := BindAndValidate(ctx, &req); err != nil {
		t.Fatalf("BindAndValidate() error = %v", err)
	}
	if req.Name != "lintao" {
		t.Fatalf("Name = %q, want lintao", req.Name)
	}
}

func TestBindAndValidateWrapsValidationError(t *testing.T) {
	var req struct {
		Name string `json:"name" validate:"required"`
	}
	ctx := newRequestTestContext(http.MethodPost, "/", `{}`)

	err := BindAndValidate(ctx, &req)
	if err == nil {
		t.Fatal("BindAndValidate() error = nil, want validation error")
	}
	appErr, ok := code.ParseError(err)
	if !ok {
		t.Fatalf("BindAndValidate() error = %T, want AppError", err)
	}
	if appErr.Code() != code.ErrValidation {
		t.Fatalf("BindAndValidate() error code = %d, want %d", appErr.Code(), code.ErrValidation)
	}
}

func TestPathInt64ParsesNamedParameter(t *testing.T) {
	ctx := newRequestTestContext(http.MethodGet, "/users/42", "")
	ctx.SetParamNames("id")
	ctx.SetParamValues("42")

	got, err := PathInt64(ctx, "id")
	if err != nil {
		t.Fatalf("PathInt64() error = %v", err)
	}
	if got != 42 {
		t.Fatalf("PathInt64() = %d, want 42", got)
	}
}

func newRequestTestContext(method, target, body string) echo.Context {
	e := echo.New()
	e.Validator = &middlewares.Validator{Validator: validator.New()}
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	if body != "" {
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	}
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec)
}
