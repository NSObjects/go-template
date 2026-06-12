package app

import (
	"context"
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

func TestAppAssemblyAppliesCapabilitySelections(t *testing.T) {
	app, err := Assemble(Options{
		Modules: []module.Module{
			staticModule{descriptor: module.Descriptor{
				Name: "user",
				Kind: module.BusinessModule,
				Requires: []module.CapabilityRef{
					{Name: "user.storage"},
				},
				EntryPoints: []module.EntryPoint{
					{Owner: "user", Type: module.EntryPointHTTP, Name: "list users", Value: noopRoute("user")},
				},
			}},
			staticModule{descriptor: module.Descriptor{
				Name: "user-storage-memory",
				Kind: module.CapabilityModule,
				Provides: []module.Capability{
					{Name: "user.storage", Provider: "memory", Status: module.CapabilityEnabled, Default: true},
				},
			}},
			staticModule{descriptor: module.Descriptor{
				Name: "user-storage-mysql",
				Kind: module.CapabilityModule,
				Provides: []module.Capability{
					{Name: "user.storage", Provider: "mysql", Status: module.CapabilityEnabled},
				},
			}},
		},
		CapabilitySelections: []module.CapabilitySelection{
			{Capability: "user.storage", Provider: "mysql"},
		},
	})
	if err != nil {
		t.Fatalf("Assemble() error = %v", err)
	}

	requirement, ok := app.Report().Requirement("user", "user.storage")
	if !ok {
		t.Fatal(`Report().Requirement("user", "user.storage") ok = false, want true`)
	}
	if !requirement.Satisfied || requirement.Provider != "mysql" {
		t.Fatalf("requirement = %+v, want satisfied by provider mysql", requirement)
	}
}

func TestAppAssemblyAppliesConfiguredCapabilitySelections(t *testing.T) {
	app, err := Assemble(Options{
		Config: configs.Config{
			Capabilities: configs.CapabilitiesConfig{
				Providers: map[string]string{
					"user.storage": "mysql",
				},
			},
		},
		Modules: userStorageModulesForTest(),
	})
	if err != nil {
		t.Fatalf("Assemble() error = %v", err)
	}

	requirement, ok := app.Report().Requirement("user", "user.storage")
	if !ok {
		t.Fatal(`Report().Requirement("user", "user.storage") ok = false, want true`)
	}
	if !requirement.Satisfied || requirement.Provider != "mysql" {
		t.Fatalf("requirement = %+v, want satisfied by configured provider mysql", requirement)
	}
}

func TestAppAssemblyRejectsUnavailableCapabilitySelection(t *testing.T) {
	_, err := Assemble(Options{
		Modules: userStorageMemoryOnlyModulesForTest(),
		CapabilitySelections: []module.CapabilitySelection{
			{Capability: "user.storage", Provider: "unknown"},
		},
	})
	if err == nil {
		t.Fatal("Assemble() error = nil, want unavailable capability provider error")
	}

	var unavailable *module.UnavailableCapabilityProviderError
	if !errors.As(err, &unavailable) {
		t.Fatalf("Assemble() error = %T, want UnavailableCapabilityProviderError", err)
	}
	if unavailable.Module != "user" || unavailable.Capability != "user.storage" || unavailable.Provider != "unknown" {
		t.Fatalf("UnavailableCapabilityProviderError = %+v, want user/user.storage/unknown", unavailable)
	}
}

func TestAppAssemblyReportsProviderSelectionBeforeCapabilityLifecycle(t *testing.T) {
	_, err := Assemble(Options{
		Modules: []module.Module{
			staticModule{descriptor: module.Descriptor{
				Name: "user",
				Kind: module.BusinessModule,
				Requires: []module.CapabilityRef{
					{Name: "user.storage"},
				},
			}},
			lifecycleModule{
				staticModule: staticModule{descriptor: module.Descriptor{
					Name: "user-storage-mysql",
					Kind: module.CapabilityModule,
					Provides: []module.Capability{
						{Name: "user.storage", Provider: "mysql", Status: module.CapabilityUnavailable},
					},
				}},
				validateErr: errors.New("mysql config invalid"),
				startErr:    errors.New("mysql connection refused"),
			},
		},
		CapabilitySelections: []module.CapabilitySelection{
			{Capability: "user.storage", Provider: "mysql"},
		},
	})
	if err == nil {
		t.Fatal("Assemble() error = nil, want unavailable capability provider error")
	}

	var unavailable *module.UnavailableCapabilityProviderError
	if !errors.As(err, &unavailable) {
		t.Fatalf("Assemble() error = %T, want UnavailableCapabilityProviderError", err)
	}
	if unavailable.Module != "user" || unavailable.Capability != "user.storage" || unavailable.Provider != "mysql" {
		t.Fatalf("UnavailableCapabilityProviderError = %+v, want user/user.storage/mysql", unavailable)
	}
}

func TestAppAssemblySkipsLifecycleForUnselectedCapabilityProvider(t *testing.T) {
	_, err := Assemble(Options{
		Modules: []module.Module{
			staticModule{descriptor: module.Descriptor{
				Name: "user",
				Kind: module.BusinessModule,
				Requires: []module.CapabilityRef{
					{Name: "user.storage"},
				},
				EntryPoints: []module.EntryPoint{
					{Owner: "user", Type: module.EntryPointHTTP, Name: "list users", Value: noopRoute("user")},
				},
			}},
			staticModule{descriptor: module.Descriptor{
				Name: "user-storage-memory",
				Kind: module.CapabilityModule,
				Provides: []module.Capability{
					{Name: "user.storage", Provider: "memory", Status: module.CapabilityEnabled, Default: true},
				},
			}},
			lifecycleModule{
				staticModule: staticModule{descriptor: module.Descriptor{
					Name: "user-storage-mysql",
					Kind: module.CapabilityModule,
					Provides: []module.Capability{
						{Name: "user.storage", Provider: "mysql", Status: module.CapabilityEnabled},
					},
				}},
				validateErr: errors.New("unselected mysql validate should not run"),
				startErr:    errors.New("unselected mysql start should not run"),
			},
		},
	})
	if err != nil {
		t.Fatalf("Assemble() error = %v, want nil for unselected capability lifecycle errors", err)
	}
}

func TestAppAssemblyRunsLifecycleForSelectedCapabilityProvider(t *testing.T) {
	validateErr := errors.New("selected mysql validate failed")
	_, err := Assemble(Options{
		Modules: []module.Module{
			staticModule{descriptor: module.Descriptor{
				Name: "user",
				Kind: module.BusinessModule,
				Requires: []module.CapabilityRef{
					{Name: "user.storage"},
				},
				EntryPoints: []module.EntryPoint{
					{Owner: "user", Type: module.EntryPointHTTP, Name: "list users", Value: noopRoute("user")},
				},
			}},
			staticModule{descriptor: module.Descriptor{
				Name: "user-storage-memory",
				Kind: module.CapabilityModule,
				Provides: []module.Capability{
					{Name: "user.storage", Provider: "memory", Status: module.CapabilityEnabled, Default: true},
				},
			}},
			lifecycleModule{
				staticModule: staticModule{descriptor: module.Descriptor{
					Name: "user-storage-mysql",
					Kind: module.CapabilityModule,
					Provides: []module.Capability{
						{Name: "user.storage", Provider: "mysql", Status: module.CapabilityEnabled},
					},
				}},
				validateErr: validateErr,
			},
		},
		CapabilitySelections: []module.CapabilitySelection{
			{Capability: "user.storage", Provider: "mysql"},
		},
	})
	if !errors.Is(err, validateErr) {
		t.Fatalf("Assemble() error = %v, want selected provider validate error", err)
	}
}

func TestAppAssemblyRunsSharedCapabilityProviderLifecycleOnce(t *testing.T) {
	var validateCalls int
	var startCalls int
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
			staticModule{descriptor: module.Descriptor{
				Name: "users",
				Kind: module.BusinessModule,
				Requires: []module.CapabilityRef{
					{Name: "mysql"},
				},
				EntryPoints: []module.EntryPoint{
					{Owner: "users", Type: module.EntryPointHTTP, Name: "list users", Value: noopRoute("users")},
				},
			}},
			lifecycleModule{
				staticModule: staticModule{descriptor: module.Descriptor{
					Name: "mysql",
					Kind: module.CapabilityModule,
					Provides: []module.Capability{
						{Name: "mysql", Status: module.CapabilityEnabled},
					},
				}},
				validateCalls: &validateCalls,
				startCalls:    &startCalls,
			},
		},
	})
	if err != nil {
		t.Fatalf("Assemble() error = %v", err)
	}
	if validateCalls != 1 {
		t.Fatalf("validateCalls = %d, want 1 for shared provider", validateCalls)
	}
	if startCalls != 1 {
		t.Fatalf("startCalls = %d, want 1 for shared provider", startCalls)
	}
}

func TestAppStopReleasesStartedCapabilityProvidersInReverseOrder(t *testing.T) {
	var stops []string
	app, err := Assemble(Options{
		Modules: []module.Module{
			staticModule{descriptor: module.Descriptor{
				Name: "orders",
				Kind: module.BusinessModule,
				Requires: []module.CapabilityRef{
					{Name: "mysql"},
					{Name: "redis"},
				},
				EntryPoints: []module.EntryPoint{
					{Owner: "orders", Type: module.EntryPointHTTP, Name: "list orders", Value: noopRoute("orders")},
				},
			}},
			lifecycleModule{
				staticModule: staticModule{descriptor: module.Descriptor{
					Name: "mysql",
					Kind: module.CapabilityModule,
					Provides: []module.Capability{
						{Name: "mysql", Status: module.CapabilityEnabled},
					},
				}},
				stopName:  "mysql",
				stopOrder: &stops,
			},
			lifecycleModule{
				staticModule: staticModule{descriptor: module.Descriptor{
					Name: "redis",
					Kind: module.CapabilityModule,
					Provides: []module.Capability{
						{Name: "redis", Status: module.CapabilityEnabled},
					},
				}},
				stopName:  "redis",
				stopOrder: &stops,
			},
		},
	})
	if err != nil {
		t.Fatalf("Assemble() error = %v", err)
	}

	if err := app.Stop(context.Background()); err != nil {
		t.Fatalf("Stop() error = %v", err)
	}

	if len(stops) != 2 || stops[0] != "redis" || stops[1] != "mysql" {
		t.Fatalf("stop order = %v, want [redis mysql]", stops)
	}
}

func TestAppStopIsIdempotent(t *testing.T) {
	var stops []string
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
			lifecycleModule{
				staticModule: staticModule{descriptor: module.Descriptor{
					Name: "mysql",
					Kind: module.CapabilityModule,
					Provides: []module.Capability{
						{Name: "mysql", Status: module.CapabilityEnabled},
					},
				}},
				stopName:  "mysql",
				stopOrder: &stops,
			},
		},
	})
	if err != nil {
		t.Fatalf("Assemble() error = %v", err)
	}

	if err := app.Stop(context.Background()); err != nil {
		t.Fatalf("Stop() error = %v", err)
	}
	if err := app.Stop(context.Background()); err != nil {
		t.Fatalf("second Stop() error = %v", err)
	}

	if len(stops) != 1 || stops[0] != "mysql" {
		t.Fatalf("stop order after repeated Stop = %v, want [mysql]", stops)
	}
}

func TestAppAssemblyStopsStartedProvidersWhenLaterStartFails(t *testing.T) {
	var stops []string
	startErr := errors.New("redis start failed")
	_, err := Assemble(Options{
		Modules: []module.Module{
			staticModule{descriptor: module.Descriptor{
				Name: "orders",
				Kind: module.BusinessModule,
				Requires: []module.CapabilityRef{
					{Name: "mysql"},
					{Name: "redis"},
				},
				EntryPoints: []module.EntryPoint{
					{Owner: "orders", Type: module.EntryPointHTTP, Name: "list orders", Value: noopRoute("orders")},
				},
			}},
			lifecycleModule{
				staticModule: staticModule{descriptor: module.Descriptor{
					Name: "mysql",
					Kind: module.CapabilityModule,
					Provides: []module.Capability{
						{Name: "mysql", Status: module.CapabilityEnabled},
					},
				}},
				stopName:  "mysql",
				stopOrder: &stops,
			},
			lifecycleModule{
				staticModule: staticModule{descriptor: module.Descriptor{
					Name: "redis",
					Kind: module.CapabilityModule,
					Provides: []module.Capability{
						{Name: "redis", Status: module.CapabilityEnabled},
					},
				}},
				startErr: startErr,
				stopName: "redis",
			},
		},
	})
	if !errors.Is(err, startErr) {
		t.Fatalf("Assemble() error = %v, want redis start error", err)
	}

	if len(stops) != 1 || stops[0] != "mysql" {
		t.Fatalf("stop order after failed startup = %v, want [mysql]", stops)
	}
}

func TestAssemblyIgnoresLegacyGeneratedFiles(t *testing.T) {
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

func TestUserModuleFailsWhenHTTPAdapterUnavailable(t *testing.T) {
	_, err := Assemble(Options{
		Modules: []module.Module{
			staticModule{descriptor: module.Descriptor{
				Name: "user",
				Kind: module.BusinessModule,
				EntryPoints: []module.EntryPoint{
					{Owner: "user", Type: module.EntryPointHTTP, Name: "list users", Value: noopRoute("user")},
				},
			}},
		},
		EntryPointAdapters: []string{"schedule"},
	})
	if err == nil {
		t.Fatal("Assemble() error = nil, want unsupported entry point error")
	}

	var unsupported *module.UnsupportedEntryPointError
	if !errors.As(err, &unsupported) {
		t.Fatalf("Assemble() error = %T, want UnsupportedEntryPointError", err)
	}
	if unsupported.Module != "user" || unsupported.EntryPointType != module.EntryPointHTTP {
		t.Fatalf("UnsupportedEntryPointError = %+v, want user/http", unsupported)
	}
}

type staticModule struct {
	descriptor module.Descriptor
}

func (m staticModule) Descriptor() module.Descriptor {
	return m.descriptor
}

func userStorageModulesForTest() []module.Module {
	return []module.Module{
		staticModule{descriptor: module.Descriptor{
			Name: "user",
			Kind: module.BusinessModule,
			Requires: []module.CapabilityRef{
				{Name: "user.storage"},
			},
			EntryPoints: []module.EntryPoint{
				{Owner: "user", Type: module.EntryPointHTTP, Name: "list users", Value: noopRoute("user")},
			},
		}},
		staticModule{descriptor: module.Descriptor{
			Name: "user-storage-memory",
			Kind: module.CapabilityModule,
			Provides: []module.Capability{
				{Name: "user.storage", Provider: "memory", Status: module.CapabilityEnabled, Default: true},
			},
		}},
		staticModule{descriptor: module.Descriptor{
			Name: "user-storage-mysql",
			Kind: module.CapabilityModule,
			Provides: []module.Capability{
				{Name: "user.storage", Provider: "mysql", Status: module.CapabilityEnabled},
			},
		}},
	}
}

func userStorageMemoryOnlyModulesForTest() []module.Module {
	return userStorageModulesForTest()[:2]
}

type lifecycleModule struct {
	staticModule
	validateErr   error
	startErr      error
	validateCalls *int
	startCalls    *int
	stopName      string
	stopOrder     *[]string
}

func (m lifecycleModule) Validate() error {
	if m.validateCalls != nil {
		*m.validateCalls++
	}
	return m.validateErr
}

func (m lifecycleModule) Start(_ context.Context) error {
	if m.startCalls != nil {
		*m.startCalls++
	}
	return m.startErr
}

func (m lifecycleModule) Stop(_ context.Context) error {
	if m.stopOrder != nil {
		*m.stopOrder = append(*m.stopOrder, m.stopName)
	}
	return nil
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
