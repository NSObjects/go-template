package boot

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/samber/do/v2"

	"github.com/NSObjects/go-template/internal/platform/configs"
	"github.com/NSObjects/go-template/internal/platform/infrastructure/resources"
)

func TestRunReturnsConfigLoadError(t *testing.T) {
	err := Run(t.TempDir() + "/missing.toml")
	if err == nil {
		t.Fatal("Run() error = nil, want config load error")
	}
	if !strings.Contains(err.Error(), "load config") {
		t.Fatalf("Run() error = %q, want load config context", err.Error())
	}
}

func TestNewAppAssemblesServer(t *testing.T) {
	app, err := NewApp(configs.Config{})
	if err != nil {
		t.Fatalf("NewApp() error = %v", err)
	}
	if app.Server() == nil {
		t.Fatal("Server() = nil, want assembled server")
	}
}

func TestNewAppInstallsModule(t *testing.T) {
	app, err := NewApp(configs.Config{}, WithModules(NewModule(
		"ping",
		Provide(func(do.Injector) (*testModuleHandler, error) {
			return &testModuleHandler{}, nil
		}),
		Route(func(group *echo.Group, handler *testModuleHandler) {
			group.GET("/module-ping", handler.Handle)
		}),
	)))
	if err != nil {
		t.Fatalf("NewApp() error = %v", err)
	}
	t.Cleanup(func() {
		if err := app.close(context.Background()); err != nil {
			t.Fatalf("close() error = %v", err)
		}
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/module-ping", nil)
	app.Server().Echo().ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("module route status = %d, want %d", rec.Code, http.StatusNoContent)
	}
}

func TestNewAppRejectsDuplicateModules(t *testing.T) {
	app, err := NewApp(configs.Config{}, WithModules(
		testModule("payments"),
		testModule("payments"),
	))
	if err == nil {
		t.Fatal("NewApp() error = nil, want duplicate module error")
	}
	if app != nil {
		t.Fatalf("NewApp() app = %#v, want nil", app)
	}
	if !strings.Contains(err.Error(), "payments") {
		t.Fatalf("NewApp() error = %q, want module name", err)
	}
}

func TestNewAppRejectsRouteWithoutRegistrar(t *testing.T) {
	app, err := NewApp(configs.Config{}, WithModules(NewModule(
		"broken",
		Provide(func(do.Injector) (*testModuleHandler, error) {
			return &testModuleHandler{}, nil
		}),
		Route[*testModuleHandler](nil),
	)))
	if err == nil {
		t.Fatal("NewApp() error = nil, want missing registrar error")
	}
	if app != nil {
		t.Fatalf("NewApp() app = %#v, want nil", app)
	}
	if !strings.Contains(err.Error(), "registrar") {
		t.Fatalf("NewApp() error = %q, want registrar context", err)
	}
}

func TestNewAppReturnsModuleProviderRegistrationError(t *testing.T) {
	app, err := NewApp(configs.Config{}, WithModules(NewModule(
		"broken",
		Provide(func(do.Injector) (*testModuleHandler, error) {
			return &testModuleHandler{}, nil
		}),
		Provide(func(do.Injector) (*testModuleHandler, error) {
			return &testModuleHandler{}, nil
		}),
		Route(func(*echo.Group, *testModuleHandler) {}),
	)))
	if err == nil {
		t.Fatal("NewApp() error = nil, want provider registration error")
	}
	if app != nil {
		t.Fatalf("NewApp() app = %#v, want nil", app)
	}
	if !strings.Contains(err.Error(), "register module broken") {
		t.Fatalf("NewApp() error = %q, want module registration context", err)
	}
	if !strings.Contains(err.Error(), "already been declared") {
		t.Fatalf("NewApp() error = %q, want do duplicate-provider context", err)
	}
}

func TestNewAppClosesInjectorWhenModuleRegistrationFails(t *testing.T) {
	resource := &testShutdownResource{}
	app, err := NewApp(configs.Config{}, WithModules(NewModule(
		"broken",
		Use(func(i do.Injector) {
			do.Provide(i, func(do.Injector) (*testShutdownResource, error) {
				return resource, nil
			})
			if _, err := do.Invoke[*testShutdownResource](i); err != nil {
				panic(err)
			}
			panic(errors.New("registration failed"))
		}),
		Route(func(*echo.Group, *testModuleHandler) {}),
	)))
	if err == nil {
		t.Fatal("NewApp() error = nil, want registration error")
	}
	if app != nil {
		t.Fatalf("NewApp() app = %#v, want nil", app)
	}
	if !resource.closed {
		t.Fatal("registration failure left invoked resource open")
	}
}

func TestNewAppReturnsModuleRouteResolutionError(t *testing.T) {
	app, err := NewApp(configs.Config{}, WithModules(NewModule(
		"broken",
		Route(func(*echo.Group, *testModuleHandler) {}),
	)))
	if err == nil {
		t.Fatal("NewApp() error = nil, want route handler resolution error")
	}
	if app != nil {
		t.Fatalf("NewApp() app = %#v, want nil", app)
	}
	if !strings.Contains(err.Error(), "mount module broken") {
		t.Fatalf("NewApp() error = %q, want route mount context", err)
	}
}

func TestNewAppReturnsConfigError(t *testing.T) {
	app, err := NewApp(configs.Config{
		JWT: configs.JWTConfig{
			Enabled: true,
		},
	})

	if err == nil {
		t.Fatal("NewApp() error = nil, want config error")
	}
	if app != nil {
		t.Fatalf("NewApp() app = %#v, want nil", app)
	}
}

func TestNewAppReturnsEnabledCapabilityConfigError(t *testing.T) {
	app, err := NewApp(configs.Config{
		MongoDB: configs.MongoDBConfig{
			Enabled: true,
		},
	})

	if err == nil {
		t.Fatal("NewApp() error = nil, want MongoDB config error")
	}
	if app != nil {
		t.Fatalf("NewApp() app = %#v, want nil", app)
	}
	if !strings.Contains(err.Error(), resources.CapabilityMongoDB) {
		t.Fatalf("NewApp() error = %q, want MongoDB capability", err)
	}
}

func TestAppCloseUsesInjectorShutdownWhenAvailable(t *testing.T) {
	app := &App{
		injector: do.New(func(i do.Injector) {
			do.Provide(i, func(do.Injector) (*testShutdownResource, error) {
				return &testShutdownResource{}, nil
			})
		}),
	}
	resource, err := do.Invoke[*testShutdownResource](app.injector)
	if err != nil {
		t.Fatalf("Invoke[*testShutdownResource]() error = %v", err)
	}

	if err := app.close(context.Background()); err != nil {
		t.Fatalf("close() error = %v", err)
	}
	if !resource.closed {
		t.Fatal("resource closed = false, want injector shutdown")
	}
}

type testShutdownResource struct {
	closed bool
}

func (r *testShutdownResource) Shutdown(context.Context) error {
	r.closed = true
	return nil
}

func TestAppCloseFallsBackToResourcesWithoutInjector(t *testing.T) {
	app := &App{
		resources: &Resources{
			status: resources.New(resources.Component{
				Name: resources.CapabilityTracing,
				Close: func(context.Context) error {
					return errors.New("flush failed")
				},
			}),
		},
	}

	err := app.close(context.Background())
	if err == nil {
		t.Fatal("close() error = nil, want resource close failure")
	}
	if !strings.Contains(err.Error(), resources.CapabilityTracing) {
		t.Fatalf("close() error = %q, want tracing capability", err)
	}
}

func TestAppRunRejectsNilServer(t *testing.T) {
	err := (&App{}).Run(context.Background())
	if err == nil {
		t.Fatal("Run() error = nil, want nil server error")
	}
}

type testModuleHandler struct{}

func (h *testModuleHandler) Handle(c *echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

func testModule(name string) Module {
	return NewModule(name, Route(func(*echo.Group, *testModuleHandler) {}))
}
