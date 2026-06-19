package salesorderhttp_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"

	salesordermemory "github.com/NSObjects/go-template/internal/modules/salesorder/adapters/memory"
	salesorderhttp "github.com/NSObjects/go-template/internal/modules/salesorder/http"
	salesorderusecase "github.com/NSObjects/go-template/internal/modules/salesorder/usecase"
	"github.com/NSObjects/go-template/internal/platform/server/middlewares"
)

func TestHandlerCreateOrder(t *testing.T) {
	e := newTestEcho()
	handler := salesorderhttp.New(salesorderusecase.New(
		salesordermemory.NewStore(),
		customerLookup{existing: map[int64]bool{7: true}},
		productLookup{products: map[int64]salesorderusecase.ProductSnapshot{11: {Exists: true, Active: true}}},
	))
	salesorderhttp.Register(e.Group("/api"), handler)

	create := doJSON(t, e, http.MethodPost, "/api/sales-orders", `{"customer_id":7,"items":[{"product_id":11,"quantity":2,"unit_price_cents":500}]}`)
	if create.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want %d: %s", create.Code, http.StatusCreated, create.Body.String())
	}
	body := decodeOrderResponse(t, create)
	if body.Data.CustomerID != 7 {
		t.Fatalf("customer_id = %d, want 7", body.Data.CustomerID)
	}
	if body.Data.TotalPriceCents != 1000 {
		t.Fatalf("total_price_cents = %d, want 1000", body.Data.TotalPriceCents)
	}
	if len(body.Data.Items) != 1 || body.Data.Items[0].LineTotalCents != 1000 {
		t.Fatalf("items = %+v, want one line total 1000", body.Data.Items)
	}
}

func TestHandlerCreateOrderRejectsInactiveProduct(t *testing.T) {
	e := newTestEcho()
	handler := salesorderhttp.New(salesorderusecase.New(
		salesordermemory.NewStore(),
		customerLookup{existing: map[int64]bool{7: true}},
		productLookup{products: map[int64]salesorderusecase.ProductSnapshot{11: {Exists: true, Active: false}}},
	))
	salesorderhttp.Register(e.Group("/api"), handler)

	create := doJSON(t, e, http.MethodPost, "/api/sales-orders", `{"customer_id":7,"items":[{"product_id":11,"quantity":2,"unit_price_cents":500}]}`)
	if create.Code != http.StatusConflict {
		t.Fatalf("create status = %d, want %d: %s", create.Code, http.StatusConflict, create.Body.String())
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

type orderResponse struct {
	Data struct {
		CustomerID      int64 `json:"customer_id"`
		TotalPriceCents int64 `json:"total_price_cents"`
		Items           []struct {
			LineTotalCents int64 `json:"line_total_cents"`
		} `json:"items"`
	} `json:"data"`
}

func decodeOrderResponse(t *testing.T, rec *httptest.ResponseRecorder) orderResponse {
	t.Helper()
	var body orderResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v\n%s", err, rec.Body.String())
	}
	return body
}

type customerLookup struct {
	existing map[int64]bool
}

func (l customerLookup) CustomerExists(_ context.Context, id int64) (bool, error) {
	return l.existing[id], nil
}

type productLookup struct {
	products map[int64]salesorderusecase.ProductSnapshot
}

func (l productLookup) FindProduct(_ context.Context, id int64) (salesorderusecase.ProductSnapshot, error) {
	return l.products[id], nil
}
