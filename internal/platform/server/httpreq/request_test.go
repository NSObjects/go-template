package httpreq

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
)

func TestPathIDParsesPositiveID(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/customers/12", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPathValues(echo.PathValues{{Name: "id", Value: "12"}})

	id, err := PathID(c, "id", "customer")
	if err != nil {
		t.Fatalf("PathID() error = %v", err)
	}
	if id != 12 {
		t.Fatalf("PathID() = %d, want 12", id)
	}
}

func TestPaginationRejectsInvalidPage(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/customers?page=bad", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	_, _, err := Pagination(c, 20)
	if err == nil {
		t.Fatal("Pagination() error = nil, want invalid page error")
	}
}

func TestQueryBoolRejectsInvalidValue(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/products?active_only=maybe", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	_, err := QueryBool(c, "active_only", false)
	if err == nil {
		t.Fatal("QueryBool() error = nil, want invalid bool error")
	}
}
