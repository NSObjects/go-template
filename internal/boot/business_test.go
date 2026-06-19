package boot

import (
	"context"
	"net/http"
	"testing"

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

	assertRoute(t, app, http.MethodPost, "/api/customers")
	assertRoute(t, app, http.MethodGet, "/api/customers")
	assertRoute(t, app, http.MethodGet, "/api/customers/:id")
	assertRoute(t, app, http.MethodPatch, "/api/customers/:id")
	assertRoute(t, app, http.MethodPost, "/api/products")
	assertRoute(t, app, http.MethodGet, "/api/products")
	assertRoute(t, app, http.MethodGet, "/api/products/:id")
	assertRoute(t, app, http.MethodPatch, "/api/products/:id")
	assertRoute(t, app, http.MethodPost, "/api/sales-orders")
	assertRoute(t, app, http.MethodGet, "/api/sales-orders")
	assertRoute(t, app, http.MethodGet, "/api/sales-orders/:id")
}

func assertRoute(t *testing.T, app *App, method, path string) {
	t.Helper()
	for _, route := range app.Server().Echo().Router().Routes() {
		if route.Method == method && route.Path == path {
			return
		}
	}
	t.Fatalf("route %s %s not registered", method, path)
}
