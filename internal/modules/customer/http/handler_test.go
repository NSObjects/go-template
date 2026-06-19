package customerhttp_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"

	customermemory "github.com/NSObjects/go-template/internal/modules/customer/adapters/memory"
	customerhttp "github.com/NSObjects/go-template/internal/modules/customer/http"
	customerusecase "github.com/NSObjects/go-template/internal/modules/customer/usecase"
	"github.com/NSObjects/go-template/internal/platform/server/middlewares"
)

func TestHandlerPatchPreservesOmittedFields(t *testing.T) {
	e := newTestEcho()
	handler := customerhttp.New(customerusecase.New(customermemory.NewStore()))
	customerhttp.Register(e.Group("/api"), handler)

	create := doJSON(t, e, http.MethodPost, "/api/customers", `{"name":"Acme","email":"Owner@Example.com"}`)
	if create.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want %d: %s", create.Code, http.StatusCreated, create.Body.String())
	}

	patch := doJSON(t, e, http.MethodPatch, "/api/customers/1", `{"status":"disabled"}`)
	if patch.Code != http.StatusOK {
		t.Fatalf("patch status = %d, want %d: %s", patch.Code, http.StatusOK, patch.Body.String())
	}
	body := decodeCustomerResponse(t, patch)
	if body.Data.Name != "Acme" {
		t.Fatalf("patched name = %q, want preserved Acme", body.Data.Name)
	}
	if body.Data.Email != "owner@example.com" {
		t.Fatalf("patched email = %q, want preserved normalized owner@example.com", body.Data.Email)
	}
	if body.Data.Status != "disabled" {
		t.Fatalf("patched status = %q, want disabled", body.Data.Status)
	}
}

func TestHandlerPatchRejectsEmptyBody(t *testing.T) {
	e := newTestEcho()
	handler := customerhttp.New(customerusecase.New(customermemory.NewStore()))
	customerhttp.Register(e.Group("/api"), handler)

	create := doJSON(t, e, http.MethodPost, "/api/customers", `{"name":"Acme","email":"owner@example.com"}`)
	if create.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want %d: %s", create.Code, http.StatusCreated, create.Body.String())
	}

	patch := doJSON(t, e, http.MethodPatch, "/api/customers/1", `{}`)
	if patch.Code != http.StatusBadRequest {
		t.Fatalf("patch status = %d, want %d: %s", patch.Code, http.StatusBadRequest, patch.Body.String())
	}
}

func newTestEcho() *echo.Echo {
	e := echo.New()
	e.Validator = &middlewares.Validator{Validator: validator.New()}
	e.HTTPErrorHandler = middlewares.ErrorHandler
	return e
}

func doJSON(t *testing.T, e *echo.Echo, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec
}

type customerResponse struct {
	Data struct {
		Name   string `json:"name"`
		Email  string `json:"email"`
		Status string `json:"status"`
	} `json:"data"`
}

func decodeCustomerResponse(t *testing.T, rec *httptest.ResponseRecorder) customerResponse {
	t.Helper()
	var body customerResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v\n%s", err, rec.Body.String())
	}
	return body
}
