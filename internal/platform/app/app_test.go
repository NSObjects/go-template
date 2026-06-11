package app

import (
	"errors"
	"testing"

	"github.com/NSObjects/go-template/internal/configs"
	platformhttp "github.com/NSObjects/go-template/internal/platform/http"
	"github.com/NSObjects/go-template/internal/platform/module"
	"github.com/labstack/echo/v4"
)

func TestAppStartupFailsWhenAssemblyFails(t *testing.T) {
	_, err := Assemble(Options{
		Modules: []module.Module{
			staticModule{descriptor: module.Descriptor{
				Name: "orders",
				Kind: module.BusinessModule,
				Requires: []module.CapabilityRef{
					{Name: "mysql"},
				},
				EntryPoints: []module.EntryPoint{
					{Owner: "orders", Type: module.EntryPointHTTP, Name: "list orders", Value: noopRoute("orders")},
				},
			}},
		},
	})
	if err == nil {
		t.Fatal("Assemble() error = nil, want missing capability error")
	}

	var missing *module.MissingCapabilityError
	if !errors.As(err, &missing) {
		t.Fatalf("Assemble() error = %T, want MissingCapabilityError", err)
	}
	if missing.Module != "orders" || missing.Capability != "mysql" {
		t.Fatalf("MissingCapabilityError = %+v, want orders/mysql", missing)
	}
}

func TestIncludedModuleBecomesActive(t *testing.T) {
	included := staticModule{descriptor: module.Descriptor{
		Name: "orders",
		Kind: module.BusinessModule,
		EntryPoints: []module.EntryPoint{
			{Owner: "orders", Type: module.EntryPointHTTP, Name: "list orders", Value: noopRoute("orders")},
		},
	}}
	app, err := Assemble(Options{
		Config:  configs.Config{},
		Modules: []module.Module{included},
	})
	if err != nil {
		t.Fatalf("Assemble() error = %v", err)
	}

	if !app.Report().HasActiveModule("orders") {
		t.Fatal("orders module is not active")
	}
	if len(app.Routes()) != 1 {
		t.Fatalf("len(app.Routes()) = %d, want 1", len(app.Routes()))
	}
}

func TestOmittedModuleStaysInactive(t *testing.T) {
	app, err := Assemble(Options{Config: configs.Config{}})
	if err != nil {
		t.Fatalf("Assemble() error = %v", err)
	}

	if app.Report().HasActiveModule("orders") {
		t.Fatal("orders module is active, want omitted module inactive")
	}
	if len(app.Routes()) != 0 {
		t.Fatalf("len(app.Routes()) = %d, want 0", len(app.Routes()))
	}
}

func TestManualBusinessModuleAssemblesWithoutGenerator(t *testing.T) {
	app, err := Assemble(Options{
		Modules: []module.Module{
			staticModule{descriptor: module.Descriptor{
				Name: "orders",
				Kind: module.BusinessModule,
				Requires: []module.CapabilityRef{
					{Name: "mysql"},
				},
				EntryPoints: []module.EntryPoint{
					{Owner: "orders", Type: module.EntryPointHTTP, Name: "list orders", Value: noopRoute("orders")},
				},
			}},
			staticModule{descriptor: module.Descriptor{
				Name: "mysql",
				Kind: module.CapabilityModule,
				Provides: []module.Capability{
					{Name: "mysql", Status: module.CapabilityEnabled},
				},
			}},
		},
	})
	if err != nil {
		t.Fatalf("Assemble() error = %v", err)
	}

	if !app.Report().HasActiveModule("orders") {
		t.Fatal("manual orders module is not active")
	}
	if len(app.Routes()) != 1 {
		t.Fatalf("len(app.Routes()) = %d, want 1", len(app.Routes()))
	}
}

func TestAssemblyIgnoresOpenAPIAndLegacyGeneratedFiles(t *testing.T) {
	app, err := Assemble(Options{Config: configs.Config{}})
	if err != nil {
		t.Fatalf("Assemble() error = %v", err)
	}

	if app.Report().HasActiveModule("user") {
		t.Fatal("legacy user module is active without explicit module inclusion")
	}
	if len(app.Routes()) != 0 {
		t.Fatalf("len(app.Routes()) = %d, want 0", len(app.Routes()))
	}
}

type staticModule struct {
	descriptor module.Descriptor
}

func (m staticModule) Descriptor() module.Descriptor {
	return m.descriptor
}

func noopRoute(owner string) platformhttp.Route {
	return platformhttp.Route{
		Owner:   owner,
		Name:    "noop",
		Method:  "GET",
		Path:    "/noop",
		Handler: func(c echo.Context) error { return c.NoContent(204) },
	}
}
