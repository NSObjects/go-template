package http

import (
	"errors"
	stdhttp "net/http"
	"testing"

	"github.com/NSObjects/go-template/internal/platform/module"
	"github.com/labstack/echo/v4"
)

func TestAdapterRegistersHTTPEntryPoint(t *testing.T) {
	e := echo.New()
	group := e.Group("/api")
	route := Route{
		Owner:   "orders",
		Name:    "list orders",
		Method:  stdhttp.MethodGet,
		Path:    "/orders",
		Handler: func(c echo.Context) error { return c.NoContent(stdhttp.StatusNoContent) },
	}

	routes, err := RoutesFromEntryPoints([]module.EntryPoint{EntryPoint("orders", route)})
	if err != nil {
		t.Fatalf("RoutesFromEntryPoints() error = %v", err)
	}
	if err := RegisterRoutes(group, routes); err != nil {
		t.Fatalf("RegisterRoutes() error = %v", err)
	}

	got := e.Routes()
	if len(got) != 1 {
		t.Fatalf("len(routes) = %d, want 1", len(got))
	}
	if got[0].Method != stdhttp.MethodGet || got[0].Path != "/api/orders" || got[0].Name != "list orders" {
		t.Fatalf("route = %+v, want GET /api/orders named list orders", got[0])
	}
}

func TestRoutesFromEntryPointsRejectsUnsupportedType(t *testing.T) {
	_, err := RoutesFromEntryPoints([]module.EntryPoint{
		{Owner: "orders", Type: "schedule", Name: "close stale orders"},
	})
	if err == nil {
		t.Fatal("RoutesFromEntryPoints() error = nil, want unsupported type error")
	}
	if !errors.Is(err, ErrUnsupportedEntryPoint) && err.Error() == "" {
		t.Fatalf("RoutesFromEntryPoints() error = %v, want non-empty unsupported type error", err)
	}
}
