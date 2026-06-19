package producthttp_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"

	productmemory "github.com/NSObjects/go-template/internal/modules/product/adapters/memory"
	producthttp "github.com/NSObjects/go-template/internal/modules/product/http"
	productusecase "github.com/NSObjects/go-template/internal/modules/product/usecase"
	"github.com/NSObjects/go-template/internal/platform/server/middlewares"
)

func TestHandlerPatchPreservesOmittedFields(t *testing.T) {
	e := newTestEcho()
	handler := producthttp.New(productusecase.New(productmemory.NewStore()))
	producthttp.Register(e.Group("/api"), handler)

	create := doJSON(t, e, http.MethodPost, "/api/products", `{"sku":"sku-1","name":"Starter","price_cents":1999}`)
	if create.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want %d: %s", create.Code, http.StatusCreated, create.Body.String())
	}

	patch := doJSON(t, e, http.MethodPatch, "/api/products/1", `{"name":"Starter Plus"}`)
	if patch.Code != http.StatusOK {
		t.Fatalf("patch status = %d, want %d: %s", patch.Code, http.StatusOK, patch.Body.String())
	}
	body := decodeProductResponse(t, patch)
	if body.Data.Name != "Starter Plus" {
		t.Fatalf("patched name = %q, want Starter Plus", body.Data.Name)
	}
	if body.Data.PriceCents != 1999 {
		t.Fatalf("patched price_cents = %d, want preserved 1999", body.Data.PriceCents)
	}
	if !body.Data.Active {
		t.Fatal("patched active = false, want preserved true")
	}
}

func TestHandlerPatchRejectsEmptyBody(t *testing.T) {
	e := newTestEcho()
	handler := producthttp.New(productusecase.New(productmemory.NewStore()))
	producthttp.Register(e.Group("/api"), handler)

	create := doJSON(t, e, http.MethodPost, "/api/products", `{"sku":"sku-1","name":"Starter","price_cents":1999}`)
	if create.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want %d: %s", create.Code, http.StatusCreated, create.Body.String())
	}

	patch := doJSON(t, e, http.MethodPatch, "/api/products/1", `{}`)
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

type productResponse struct {
	Data struct {
		Name       string `json:"name"`
		PriceCents int64  `json:"price_cents"`
		Active     bool   `json:"active"`
	} `json:"data"`
}

func decodeProductResponse(t *testing.T, rec *httptest.ResponseRecorder) productResponse {
	t.Helper()
	var body productResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v\n%s", err, rec.Body.String())
	}
	return body
}
