package boot

import (
	"context"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/NSObjects/go-template/internal/platform/configs"
)

func TestBusinessModulesRegisterRoutes(t *testing.T) {
	app, err := NewApp(configs.Config{}, WithModules(BusinessModules()...))
	if err != nil {
		t.Fatalf("NewApp() error = %v", err)
	}
	t.Cleanup(func() {
		if err := app.close(context.Background()); err != nil {
			t.Fatalf("close() error = %v", err)
		}
	})

	assertRoute(t, app, echo.POST, "/api/customers")
	assertRoute(t, app, echo.GET, "/api/customers")
	assertRoute(t, app, echo.GET, "/api/customers/:id")
	assertRoute(t, app, echo.PATCH, "/api/customers/:id")
	assertRoute(t, app, echo.POST, "/api/products")
	assertRoute(t, app, echo.GET, "/api/products")
	assertRoute(t, app, echo.GET, "/api/products/:id")
	assertRoute(t, app, echo.PATCH, "/api/products/:id")
	assertRoute(t, app, echo.POST, "/api/sales-orders")
	assertRoute(t, app, echo.GET, "/api/sales-orders")
	assertRoute(t, app, echo.GET, "/api/sales-orders/:id")
}

func assertRoute(t *testing.T, app *App, method, path string) {
	t.Helper()
	for _, route := range app.Server().Echo().Routes() {
		if route.Method == method && route.Path == path {
			return
		}
	}
	t.Fatalf("route %s %s not registered", method, path)
}
