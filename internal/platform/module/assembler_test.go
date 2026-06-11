package module

import (
	"errors"
	"reflect"
	"testing"
)

func TestAssembleReportsSatisfiedCapabilities(t *testing.T) {
	tests := []struct {
		name              string
		modules           []Module
		adapters          []string
		wantActiveModules []string
		wantCapability    CapabilityStatus
	}{
		{
			name: "required capability is satisfied",
			modules: []Module{
				staticModule{descriptor: Descriptor{
					Name: "orders",
					Kind: BusinessModule,
					Requires: []CapabilityRef{
						{Name: "mysql"},
					},
					EntryPoints: []EntryPoint{
						{Owner: "orders", Type: EntryPointHTTP, Name: "list orders"},
					},
				}},
				staticModule{descriptor: Descriptor{
					Name: "mysql",
					Kind: CapabilityModule,
					Provides: []Capability{
						{Name: "mysql", Status: CapabilityEnabled},
					},
				}},
			},
			adapters:          []string{EntryPointHTTP},
			wantActiveModules: []string{"orders", "mysql"},
			wantCapability:    CapabilityStatus{Name: "mysql", Status: CapabilityEnabled, Provider: "mysql"},
		},
		{
			name: "disabled optional capability is observable",
			modules: []Module{
				staticModule{descriptor: Descriptor{
					Name: "redis",
					Kind: CapabilityModule,
					Provides: []Capability{
						{Name: "redis", Status: CapabilityDisabled},
					},
				}},
			},
			adapters:          []string{EntryPointHTTP},
			wantActiveModules: []string{"redis"},
			wantCapability:    CapabilityStatus{Name: "redis", Status: CapabilityDisabled, Provider: "redis"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report, err := Assemble(tt.modules, WithEntryPointAdapters(tt.adapters...))
			if err != nil {
				t.Fatalf("Assemble() error = %v", err)
			}

			for _, moduleName := range tt.wantActiveModules {
				if !report.HasActiveModule(moduleName) {
					t.Fatalf("report.HasActiveModule(%q) = false, want true", moduleName)
				}
			}

			got, ok := report.Capability(tt.wantCapability.Name)
			if !ok {
				t.Fatalf("report.Capability(%q) ok = false, want true", tt.wantCapability.Name)
			}
			if got != tt.wantCapability {
				t.Fatalf("report.Capability(%q) = %+v, want %+v", tt.wantCapability.Name, got, tt.wantCapability)
			}
			if tt.wantCapability.Status == CapabilityEnabled {
				requirement, ok := report.Requirement("orders", tt.wantCapability.Name)
				if !ok {
					t.Fatalf("report.Requirement(%q, %q) ok = false, want true", "orders", tt.wantCapability.Name)
				}
				if !requirement.Satisfied || requirement.Provider != tt.wantCapability.Provider {
					t.Fatalf("requirement = %+v, want satisfied by %q", requirement, tt.wantCapability.Provider)
				}
			}
		})
	}
}

func TestAssembleUsesDefaultCapabilityProvider(t *testing.T) {
	report, err := Assemble([]Module{
		staticModule{descriptor: Descriptor{
			Name: "user",
			Kind: BusinessModule,
			Requires: []CapabilityRef{
				{Name: "user.storage"},
			},
			EntryPoints: []EntryPoint{
				{Owner: "user", Type: EntryPointHTTP, Name: "list users"},
			},
		}},
		staticModule{descriptor: Descriptor{
			Name: "user-storage-memory",
			Kind: CapabilityModule,
			Provides: []Capability{
				{Name: "user.storage", Provider: "memory", Status: CapabilityEnabled, Default: true},
			},
		}},
	}, WithEntryPointAdapters(EntryPointHTTP))
	if err != nil {
		t.Fatalf("Assemble() error = %v", err)
	}

	requirement, ok := report.Requirement("user", "user.storage")
	if !ok {
		t.Fatal(`report.Requirement("user", "user.storage") ok = false, want true`)
	}
	if !requirement.Satisfied || requirement.Provider != "memory" {
		t.Fatalf("requirement = %+v, want satisfied by provider memory", requirement)
	}

	capability, ok := report.Capability("user.storage")
	if !ok {
		t.Fatal(`report.Capability("user.storage") ok = false, want true`)
	}
	if capability.Provider != "memory" {
		t.Fatalf("capability.Provider = %q, want memory", capability.Provider)
	}
}

func TestAssembleUsesSelectedCapabilityProvider(t *testing.T) {
	memoryRepository := testRepository{name: "memory"}
	mysqlRepository := testRepository{name: "mysql"}
	report, err := Assemble([]Module{
		staticModule{descriptor: Descriptor{
			Name: "user",
			Kind: BusinessModule,
			Requires: []CapabilityRef{
				{Name: "user.storage"},
			},
			EntryPoints: []EntryPoint{
				{Owner: "user", Type: EntryPointHTTP, Name: "list users"},
			},
		}},
		staticModule{descriptor: Descriptor{
			Name: "user-storage-memory",
			Kind: CapabilityModule,
			Provides: []Capability{
				{Name: "user.storage", Provider: "memory", Status: CapabilityEnabled, Default: true, Value: memoryRepository},
			},
		}},
		staticModule{descriptor: Descriptor{
			Name: "user-storage-mysql",
			Kind: CapabilityModule,
			Provides: []Capability{
				{Name: "user.storage", Provider: "mysql", Status: CapabilityEnabled, Value: mysqlRepository},
			},
		}},
	}, WithEntryPointAdapters(EntryPointHTTP), WithCapabilitySelections(CapabilitySelection{
		Capability: "user.storage",
		Provider:   "mysql",
	}))
	if err != nil {
		t.Fatalf("Assemble() error = %v", err)
	}

	requirement, ok := report.Requirement("user", "user.storage")
	if !ok {
		t.Fatal(`report.Requirement("user", "user.storage") ok = false, want true`)
	}
	if !requirement.Satisfied || requirement.Provider != "mysql" {
		t.Fatalf("requirement = %+v, want satisfied by provider mysql", requirement)
	}
	if _, ok := reflect.TypeOf(CapabilityStatus{}).FieldByName("Value"); ok {
		t.Fatal("CapabilityStatus exposes Value field, want observable report metadata only")
	}
}

func TestResolveCapabilityValueFromModulesUsesProviderSelectionMap(t *testing.T) {
	memoryRepository := testRepository{name: "memory"}
	mysqlRepository := testRepository{name: "mysql"}

	resolved, err := ResolveCapabilityValueFromModules[testRepository](
		[]Module{
			staticModule{descriptor: Descriptor{
				Name: "user-storage-memory",
				Kind: CapabilityModule,
				Provides: []Capability{
					{Name: "user.storage", Provider: "memory", Status: CapabilityEnabled, Default: true, Value: memoryRepository},
				},
			}},
			staticModule{descriptor: Descriptor{
				Name: "user-storage-mysql",
				Kind: CapabilityModule,
				Provides: []Capability{
					{Name: "user.storage", Provider: "mysql", Status: CapabilityEnabled, Value: mysqlRepository},
				},
			}},
		},
		"user",
		"user.storage",
		WithCapabilityProviderSelections(map[string]string{"user.storage": "mysql"}),
	)
	if err != nil {
		t.Fatalf("ResolveCapabilityValueFromModules() error = %v", err)
	}
	if resolved != mysqlRepository {
		t.Fatalf("resolved repository = %+v, want %+v", resolved, mysqlRepository)
	}
}

func TestResolveCapabilityValueFromModulesRejectsMissingRuntimeValue(t *testing.T) {
	_, err := ResolveCapabilityValueFromModules[testRepository](
		[]Module{
			staticModule{descriptor: Descriptor{
				Name: "user-storage-memory",
				Kind: CapabilityModule,
				Provides: []Capability{
					{Name: "user.storage", Provider: "memory", Status: CapabilityEnabled, Default: true},
				},
			}},
		},
		"user",
		"user.storage",
	)
	if err == nil {
		t.Fatal("ResolveCapabilityValueFromModules() error = nil, want missing capability value error")
	}

	var missingValue *MissingCapabilityValueError
	if !errors.As(err, &missingValue) {
		t.Fatalf("ResolveCapabilityValueFromModules() error = %T, want MissingCapabilityValueError", err)
	}
	if missingValue.Module != "user" || missingValue.Capability != "user.storage" || missingValue.Provider != "memory" {
		t.Fatalf("MissingCapabilityValueError = %+v, want user/user.storage/memory", missingValue)
	}
}

func TestAssembleUsesDefaultCapabilityValue(t *testing.T) {
	memoryRepository := testRepository{name: "memory"}
	mysqlRepository := testRepository{name: "mysql"}

	resolved, err := ResolveCapabilityValueFromModules[testRepository](
		[]Module{
			staticModule{descriptor: Descriptor{
				Name: "user-storage-memory",
				Kind: CapabilityModule,
				Provides: []Capability{
					{Name: "user.storage", Provider: "memory", Status: CapabilityEnabled, Default: true, Value: memoryRepository},
				},
			}},
			staticModule{descriptor: Descriptor{
				Name: "user-storage-mysql",
				Kind: CapabilityModule,
				Provides: []Capability{
					{Name: "user.storage", Provider: "mysql", Status: CapabilityEnabled, Value: mysqlRepository},
				},
			}},
		},
		"user",
		"user.storage",
	)
	if err != nil {
		t.Fatalf("ResolveCapabilityValueFromModules() error = %v", err)
	}
	if resolved != memoryRepository {
		t.Fatalf("resolved repository = %+v, want %+v", resolved, memoryRepository)
	}
}

func TestAssembleRejectsMissingRequiredCapability(t *testing.T) {
	_, err := Assemble([]Module{
		staticModule{descriptor: Descriptor{
			Name: "orders",
			Kind: BusinessModule,
			Requires: []CapabilityRef{
				{Name: "mysql"},
			},
			EntryPoints: []EntryPoint{
				{Owner: "orders", Type: EntryPointHTTP, Name: "list orders"},
			},
		}},
	}, WithEntryPointAdapters(EntryPointHTTP))
	if err == nil {
		t.Fatal("Assemble() error = nil, want missing capability error")
	}

	var missing *MissingCapabilityError
	if !errors.As(err, &missing) {
		t.Fatalf("Assemble() error = %T, want MissingCapabilityError", err)
	}
	if missing.Module != "orders" || missing.Capability != "mysql" {
		t.Fatalf("MissingCapabilityError = %+v, want module orders capability mysql", missing)
	}
}

func TestAssembleRejectsUnavailableSelectedCapabilityProvider(t *testing.T) {
	_, err := Assemble([]Module{
		staticModule{descriptor: Descriptor{
			Name: "user",
			Kind: BusinessModule,
			Requires: []CapabilityRef{
				{Name: "user.storage"},
			},
			EntryPoints: []EntryPoint{
				{Owner: "user", Type: EntryPointHTTP, Name: "list users"},
			},
		}},
		staticModule{descriptor: Descriptor{
			Name: "user-storage-memory",
			Kind: CapabilityModule,
			Provides: []Capability{
				{Name: "user.storage", Provider: "memory", Status: CapabilityEnabled, Default: true},
			},
		}},
	}, WithEntryPointAdapters(EntryPointHTTP), WithCapabilitySelections(CapabilitySelection{
		Capability: "user.storage",
		Provider:   "unknown",
	}))
	if err == nil {
		t.Fatal("Assemble() error = nil, want unavailable selected provider error")
	}

	var unavailable *UnavailableCapabilityProviderError
	if !errors.As(err, &unavailable) {
		t.Fatalf("Assemble() error = %T, want UnavailableCapabilityProviderError", err)
	}
	if unavailable.Module != "user" || unavailable.Capability != "user.storage" || unavailable.Provider != "unknown" {
		t.Fatalf("UnavailableCapabilityProviderError = %+v, want user/user.storage/unknown", unavailable)
	}
}

func TestAssembleReportsBusinessModuleWithNoEntryPoints(t *testing.T) {
	report, err := Assemble([]Module{
		staticModule{descriptor: Descriptor{
			Name: "orders",
			Kind: BusinessModule,
		}},
	})
	if err != nil {
		t.Fatalf("Assemble() error = %v", err)
	}

	if len(report.Warnings) != 1 {
		t.Fatalf("len(report.Warnings) = %d, want 1", len(report.Warnings))
	}
	got := report.Warnings[0]
	if got.Module != "orders" || got.Code != WarningNoEntryPoints {
		t.Fatalf("warning = %+v, want no entry points warning for orders", got)
	}
}

func TestAssembleRejectsUnsupportedEntryPoint(t *testing.T) {
	_, err := Assemble([]Module{
		staticModule{descriptor: Descriptor{
			Name: "orders",
			Kind: BusinessModule,
			EntryPoints: []EntryPoint{
				{Owner: "orders", Type: "schedule", Name: "close stale orders"},
			},
		}},
	}, WithEntryPointAdapters(EntryPointHTTP))
	if err == nil {
		t.Fatal("Assemble() error = nil, want unsupported entry point error")
	}

	var unsupported *UnsupportedEntryPointError
	if !errors.As(err, &unsupported) {
		t.Fatalf("Assemble() error = %T, want UnsupportedEntryPointError", err)
	}
	if unsupported.Module != "orders" || unsupported.EntryPointType != "schedule" {
		t.Fatalf("UnsupportedEntryPointError = %+v, want module orders type schedule", unsupported)
	}
}

func TestAssembleManualDescriptors(t *testing.T) {
	report, err := Assemble([]Module{
		staticModule{descriptor: Descriptor{
			Name: "orders",
			Kind: BusinessModule,
			Requires: []CapabilityRef{
				{Name: "mysql"},
			},
			EntryPoints: []EntryPoint{
				{Owner: "orders", Type: EntryPointHTTP, Name: "list orders"},
			},
		}},
		staticModule{descriptor: Descriptor{
			Name: "mysql",
			Kind: CapabilityModule,
			Provides: []Capability{
				{Name: "mysql", Status: CapabilityEnabled},
			},
		}},
	}, WithEntryPointAdapters(EntryPointHTTP))
	if err != nil {
		t.Fatalf("Assemble() error = %v", err)
	}

	if !report.HasActiveModule("orders") {
		t.Fatal("manual business module is not active")
	}
	if len(report.EntryPoints) != 1 {
		t.Fatalf("len(report.EntryPoints) = %d, want 1", len(report.EntryPoints))
	}
}

type staticModule struct {
	descriptor Descriptor
}

func (m staticModule) Descriptor() Descriptor {
	return m.descriptor
}

type testRepository struct {
	name string
}
