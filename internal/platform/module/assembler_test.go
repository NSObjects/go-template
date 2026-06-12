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
						{Name: "s3"},
					},
					EntryPoints: []EntryPoint{
						{Owner: "orders", Type: EntryPointHTTP, Name: "list orders"},
					},
				}},
				staticModule{descriptor: Descriptor{
					Name: "s3",
					Kind: CapabilityModule,
					Provides: []Capability{
						{Name: "s3", Status: CapabilityEnabled},
					},
				}},
			},
			adapters:          []string{EntryPointHTTP},
			wantActiveModules: []string{"orders", "s3"},
			wantCapability:    CapabilityStatus{Name: "s3", Status: CapabilityEnabled, Provider: "s3"},
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
			Name: "documents",
			Kind: BusinessModule,
			Requires: []CapabilityRef{
				{Name: "document.store"},
			},
			EntryPoints: []EntryPoint{
				{Owner: "documents", Type: EntryPointHTTP, Name: "list documents"},
			},
		}},
		staticModule{descriptor: Descriptor{
			Name: "document-store-memory",
			Kind: CapabilityModule,
			Provides: []Capability{
				{Name: "document.store", Provider: "memory", Status: CapabilityEnabled, Default: true},
			},
		}},
	}, WithEntryPointAdapters(EntryPointHTTP))
	if err != nil {
		t.Fatalf("Assemble() error = %v", err)
	}

	requirement, ok := report.Requirement("documents", "document.store")
	if !ok {
		t.Fatal(`report.Requirement("documents", "document.store") ok = false, want true`)
	}
	if !requirement.Satisfied || requirement.Provider != "memory" {
		t.Fatalf("requirement = %+v, want satisfied by provider memory", requirement)
	}

	capability, ok := report.Capability("document.store")
	if !ok {
		t.Fatal(`report.Capability("document.store") ok = false, want true`)
	}
	if capability.Provider != "memory" {
		t.Fatalf("capability.Provider = %q, want memory", capability.Provider)
	}
}

func TestAssembleUsesSelectedCapabilityProvider(t *testing.T) {
	memoryRepository := testRepository{name: "memory"}
	s3Repository := testRepository{name: "s3"}
	report, err := Assemble([]Module{
		staticModule{descriptor: Descriptor{
			Name: "documents",
			Kind: BusinessModule,
			Requires: []CapabilityRef{
				{Name: "document.store"},
			},
			EntryPoints: []EntryPoint{
				{Owner: "documents", Type: EntryPointHTTP, Name: "list documents"},
			},
		}},
		staticModule{descriptor: Descriptor{
			Name: "document-store-memory",
			Kind: CapabilityModule,
			Provides: []Capability{
				{Name: "document.store", Provider: "memory", Status: CapabilityEnabled, Default: true, Value: memoryRepository},
			},
		}},
		staticModule{descriptor: Descriptor{
			Name: "document-store-s3",
			Kind: CapabilityModule,
			Provides: []Capability{
				{Name: "document.store", Provider: "s3", Status: CapabilityEnabled, Value: s3Repository},
			},
		}},
	}, WithEntryPointAdapters(EntryPointHTTP), WithCapabilitySelections(CapabilitySelection{
		Capability: "document.store",
		Provider:   "s3",
	}))
	if err != nil {
		t.Fatalf("Assemble() error = %v", err)
	}

	requirement, ok := report.Requirement("documents", "document.store")
	if !ok {
		t.Fatal(`report.Requirement("documents", "document.store") ok = false, want true`)
	}
	if !requirement.Satisfied || requirement.Provider != "s3" {
		t.Fatalf("requirement = %+v, want satisfied by provider s3", requirement)
	}
	if _, ok := reflect.TypeOf(CapabilityStatus{}).FieldByName("Value"); ok {
		t.Fatal("CapabilityStatus exposes Value field, want observable report metadata only")
	}
}

func TestAssembleRejectsAmbiguousCapabilityProvider(t *testing.T) {
	_, err := Assemble([]Module{
		staticModule{descriptor: Descriptor{
			Name: "documents",
			Kind: BusinessModule,
			Requires: []CapabilityRef{
				{Name: "document.store"},
			},
		}},
		staticModule{descriptor: Descriptor{
			Name: "document-store-memory",
			Kind: CapabilityModule,
			Provides: []Capability{
				{Name: "document.store", Provider: "memory", Status: CapabilityEnabled},
			},
		}},
		staticModule{descriptor: Descriptor{
			Name: "document-store-s3",
			Kind: CapabilityModule,
			Provides: []Capability{
				{Name: "document.store", Provider: "s3", Status: CapabilityEnabled},
			},
		}},
	})
	if err == nil {
		t.Fatal("Assemble() error = nil, want ambiguous capability provider error")
	}

	var ambiguous *AmbiguousCapabilityProviderError
	if !errors.As(err, &ambiguous) {
		t.Fatalf("Assemble() error = %T, want AmbiguousCapabilityProviderError", err)
	}
	if ambiguous.Module != "documents" || ambiguous.Capability != "document.store" {
		t.Fatalf("AmbiguousCapabilityProviderError = %+v, want documents/document.store", ambiguous)
	}
	if !reflect.DeepEqual(ambiguous.Providers, []string{"memory", "s3"}) {
		t.Fatalf("AmbiguousCapabilityProviderError.Providers = %+v, want [memory s3]", ambiguous.Providers)
	}
}

func TestAssembleRejectsMultipleDefaultCapabilityProviders(t *testing.T) {
	_, err := Assemble([]Module{
		staticModule{descriptor: Descriptor{
			Name: "documents",
			Kind: BusinessModule,
			Requires: []CapabilityRef{
				{Name: "document.store"},
			},
		}},
		staticModule{descriptor: Descriptor{
			Name: "document-store-memory",
			Kind: CapabilityModule,
			Provides: []Capability{
				{Name: "document.store", Provider: "memory", Status: CapabilityEnabled, Default: true},
			},
		}},
		staticModule{descriptor: Descriptor{
			Name: "document-store-s3",
			Kind: CapabilityModule,
			Provides: []Capability{
				{Name: "document.store", Provider: "s3", Status: CapabilityEnabled, Default: true},
			},
		}},
	})
	if err == nil {
		t.Fatal("Assemble() error = nil, want ambiguous capability provider error")
	}

	var ambiguous *AmbiguousCapabilityProviderError
	if !errors.As(err, &ambiguous) {
		t.Fatalf("Assemble() error = %T, want AmbiguousCapabilityProviderError", err)
	}
	if ambiguous.Module != "documents" || ambiguous.Capability != "document.store" {
		t.Fatalf("AmbiguousCapabilityProviderError = %+v, want documents/document.store", ambiguous)
	}
	if !reflect.DeepEqual(ambiguous.Providers, []string{"memory", "s3"}) {
		t.Fatalf("AmbiguousCapabilityProviderError.Providers = %+v, want [memory s3]", ambiguous.Providers)
	}
}

func TestResolveCapabilityValueFromModulesUsesProviderSelectionMap(t *testing.T) {
	memoryRepository := testRepository{name: "memory"}
	s3Repository := testRepository{name: "s3"}

	resolved, err := ResolveCapabilityValueFromModules[testRepository](
		[]Module{
			staticModule{descriptor: Descriptor{
				Name: "document-store-memory",
				Kind: CapabilityModule,
				Provides: []Capability{
					{Name: "document.store", Provider: "memory", Status: CapabilityEnabled, Default: true, Value: memoryRepository},
				},
			}},
			staticModule{descriptor: Descriptor{
				Name: "document-store-s3",
				Kind: CapabilityModule,
				Provides: []Capability{
					{Name: "document.store", Provider: "s3", Status: CapabilityEnabled, Value: s3Repository},
				},
			}},
		},
		"documents",
		"document.store",
		WithCapabilityProviderSelections(map[string]string{"document.store": "s3"}),
	)
	if err != nil {
		t.Fatalf("ResolveCapabilityValueFromModules() error = %v", err)
	}
	if resolved != s3Repository {
		t.Fatalf("resolved repository = %+v, want %+v", resolved, s3Repository)
	}
}

func TestResolveCapabilityValueFromModulesRejectsMissingRuntimeValue(t *testing.T) {
	_, err := ResolveCapabilityValueFromModules[testRepository](
		[]Module{
			staticModule{descriptor: Descriptor{
				Name: "document-store-memory",
				Kind: CapabilityModule,
				Provides: []Capability{
					{Name: "document.store", Provider: "memory", Status: CapabilityEnabled, Default: true},
				},
			}},
		},
		"documents",
		"document.store",
	)
	if err == nil {
		t.Fatal("ResolveCapabilityValueFromModules() error = nil, want missing capability value error")
	}

	var missingValue *MissingCapabilityValueError
	if !errors.As(err, &missingValue) {
		t.Fatalf("ResolveCapabilityValueFromModules() error = %T, want MissingCapabilityValueError", err)
	}
	if missingValue.Module != "documents" || missingValue.Capability != "document.store" || missingValue.Provider != "memory" {
		t.Fatalf("MissingCapabilityValueError = %+v, want documents/document.store/memory", missingValue)
	}
}

func TestAssembleRuntimeReportsSelectedCapabilityModuleIndexes(t *testing.T) {
	result, err := AssembleRuntime([]Module{
		staticModule{descriptor: Descriptor{
			Name: "documents",
			Kind: BusinessModule,
			Requires: []CapabilityRef{
				{Name: "document.store"},
			},
			EntryPoints: []EntryPoint{
				{Owner: "documents", Type: EntryPointHTTP, Name: "list documents"},
			},
		}},
		staticModule{descriptor: Descriptor{
			Name: "document-store-memory",
			Kind: CapabilityModule,
			Provides: []Capability{
				{Name: "document.store", Provider: "memory", Status: CapabilityEnabled, Default: true},
			},
		}},
		staticModule{descriptor: Descriptor{
			Name: "document-store-s3",
			Kind: CapabilityModule,
			Provides: []Capability{
				{Name: "document.store", Provider: "s3", Status: CapabilityEnabled},
			},
		}},
	}, WithEntryPointAdapters(EntryPointHTTP), WithCapabilitySelections(CapabilitySelection{
		Capability: "document.store",
		Provider:   "s3",
	}))
	if err != nil {
		t.Fatalf("AssembleRuntime() error = %v", err)
	}

	selected := result.SelectedCapabilityModuleIndexes()
	if len(selected) != 1 {
		t.Fatalf("len(SelectedCapabilityModuleIndexes()) = %d, want 1", len(selected))
	}
	if selected[0] != 2 {
		t.Fatalf("SelectedCapabilityModuleIndexes()[0] = %d, want s3 provider index 2", selected[0])
	}
}

func TestAssembleRuntimeReportsDefaultCapabilityModuleIndex(t *testing.T) {
	result, err := AssembleRuntime([]Module{
		staticModule{descriptor: Descriptor{
			Name: "documents",
			Kind: BusinessModule,
			Requires: []CapabilityRef{
				{Name: "document.store"},
			},
			EntryPoints: []EntryPoint{
				{Owner: "documents", Type: EntryPointHTTP, Name: "list documents"},
			},
		}},
		staticModule{descriptor: Descriptor{
			Name: "document-store-memory",
			Kind: CapabilityModule,
			Provides: []Capability{
				{Name: "document.store", Provider: "memory", Status: CapabilityEnabled, Default: true},
			},
		}},
		staticModule{descriptor: Descriptor{
			Name: "document-store-s3",
			Kind: CapabilityModule,
			Provides: []Capability{
				{Name: "document.store", Provider: "s3", Status: CapabilityEnabled},
			},
		}},
	}, WithEntryPointAdapters(EntryPointHTTP))
	if err != nil {
		t.Fatalf("AssembleRuntime() error = %v", err)
	}

	selected := result.SelectedCapabilityModuleIndexes()
	if len(selected) != 1 {
		t.Fatalf("len(SelectedCapabilityModuleIndexes()) = %d, want 1", len(selected))
	}
	if selected[0] != 1 {
		t.Fatalf("SelectedCapabilityModuleIndexes()[0] = %d, want memory provider index 1", selected[0])
	}
}

func TestAssembleRuntimeSelectedCapabilityModuleIndexesReturnsCopy(t *testing.T) {
	result, err := AssembleRuntime([]Module{
		staticModule{descriptor: Descriptor{
			Name: "documents",
			Kind: BusinessModule,
			Requires: []CapabilityRef{
				{Name: "document.store"},
			},
		}},
		staticModule{descriptor: Descriptor{
			Name: "document-store-memory",
			Kind: CapabilityModule,
			Provides: []Capability{
				{Name: "document.store", Provider: "memory", Status: CapabilityEnabled, Default: true},
			},
		}},
	})
	if err != nil {
		t.Fatalf("AssembleRuntime() error = %v", err)
	}

	selected := result.SelectedCapabilityModuleIndexes()
	selected[0] = 99

	selectedAgain := result.SelectedCapabilityModuleIndexes()
	if len(selectedAgain) != 1 {
		t.Fatalf("len(SelectedCapabilityModuleIndexes()) = %d, want 1", len(selectedAgain))
	}
	if selectedAgain[0] != 1 {
		t.Fatalf("SelectedCapabilityModuleIndexes()[0] = %d after caller mutation, want 1", selectedAgain[0])
	}
}

func TestAssembleRuntimeReportsSharedCapabilityModuleIndexOnce(t *testing.T) {
	result, err := AssembleRuntime([]Module{
		staticModule{descriptor: Descriptor{
			Name: "orders",
			Kind: BusinessModule,
			Requires: []CapabilityRef{
				{Name: "s3"},
			},
			EntryPoints: []EntryPoint{
				{Owner: "orders", Type: EntryPointHTTP, Name: "list orders"},
			},
		}},
		staticModule{descriptor: Descriptor{
			Name: "users",
			Kind: BusinessModule,
			Requires: []CapabilityRef{
				{Name: "s3"},
			},
			EntryPoints: []EntryPoint{
				{Owner: "users", Type: EntryPointHTTP, Name: "list documents"},
			},
		}},
		staticModule{descriptor: Descriptor{
			Name: "s3",
			Kind: CapabilityModule,
			Provides: []Capability{
				{Name: "s3", Status: CapabilityEnabled},
			},
		}},
	}, WithEntryPointAdapters(EntryPointHTTP))
	if err != nil {
		t.Fatalf("AssembleRuntime() error = %v", err)
	}

	selected := result.SelectedCapabilityModuleIndexes()
	if len(selected) != 1 {
		t.Fatalf("len(SelectedCapabilityModuleIndexes()) = %d, want 1 shared provider", len(selected))
	}
	if selected[0] != 2 {
		t.Fatalf("SelectedCapabilityModuleIndexes()[0] = %d, want shared s3 provider index 2", selected[0])
	}
}

func TestAssembleUsesDefaultCapabilityValue(t *testing.T) {
	memoryRepository := testRepository{name: "memory"}
	s3Repository := testRepository{name: "s3"}

	resolved, err := ResolveCapabilityValueFromModules[testRepository](
		[]Module{
			staticModule{descriptor: Descriptor{
				Name: "document-store-memory",
				Kind: CapabilityModule,
				Provides: []Capability{
					{Name: "document.store", Provider: "memory", Status: CapabilityEnabled, Default: true, Value: memoryRepository},
				},
			}},
			staticModule{descriptor: Descriptor{
				Name: "document-store-s3",
				Kind: CapabilityModule,
				Provides: []Capability{
					{Name: "document.store", Provider: "s3", Status: CapabilityEnabled, Value: s3Repository},
				},
			}},
		},
		"documents",
		"document.store",
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
				{Name: "s3"},
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
	if missing.Module != "orders" || missing.Capability != "s3" {
		t.Fatalf("MissingCapabilityError = %+v, want module orders capability s3", missing)
	}
}

func TestAssembleRejectsUnavailableSelectedCapabilityProvider(t *testing.T) {
	_, err := Assemble([]Module{
		staticModule{descriptor: Descriptor{
			Name: "documents",
			Kind: BusinessModule,
			Requires: []CapabilityRef{
				{Name: "document.store"},
			},
			EntryPoints: []EntryPoint{
				{Owner: "documents", Type: EntryPointHTTP, Name: "list documents"},
			},
		}},
		staticModule{descriptor: Descriptor{
			Name: "document-store-memory",
			Kind: CapabilityModule,
			Provides: []Capability{
				{Name: "document.store", Provider: "memory", Status: CapabilityEnabled, Default: true},
			},
		}},
	}, WithEntryPointAdapters(EntryPointHTTP), WithCapabilitySelections(CapabilitySelection{
		Capability: "document.store",
		Provider:   "unknown",
	}))
	if err == nil {
		t.Fatal("Assemble() error = nil, want unavailable selected provider error")
	}

	var unavailable *UnavailableCapabilityProviderError
	if !errors.As(err, &unavailable) {
		t.Fatalf("Assemble() error = %T, want UnavailableCapabilityProviderError", err)
	}
	if unavailable.Module != "documents" || unavailable.Capability != "document.store" || unavailable.Provider != "unknown" {
		t.Fatalf("UnavailableCapabilityProviderError = %+v, want documents/document.store/unknown", unavailable)
	}
}

func TestAssembleRejectsDuplicateCapabilityProvider(t *testing.T) {
	_, err := Assemble([]Module{
		staticModule{descriptor: Descriptor{
			Name: "document-store-memory",
			Kind: CapabilityModule,
			Provides: []Capability{
				{Name: "document.store", Provider: "memory", Status: CapabilityEnabled, Default: true},
			},
		}},
		staticModule{descriptor: Descriptor{
			Name: "another-document-memory-store",
			Kind: CapabilityModule,
			Provides: []Capability{
				{Name: "document.store", Provider: "memory", Status: CapabilityEnabled},
			},
		}},
	})
	if err == nil {
		t.Fatal("Assemble() error = nil, want duplicate capability provider error")
	}

	var duplicate *DuplicateCapabilityProviderError
	if !errors.As(err, &duplicate) {
		t.Fatalf("Assemble() error = %T, want DuplicateCapabilityProviderError", err)
	}
	if duplicate.Capability != "document.store" || duplicate.Provider != "memory" {
		t.Fatalf("DuplicateCapabilityProviderError = %+v, want document.store/memory", duplicate)
	}
	if duplicate.FirstModule != "document-store-memory" || duplicate.SecondModule != "another-document-memory-store" {
		t.Fatalf("DuplicateCapabilityProviderError modules = %+v, want document-store-memory/another-document-memory-store", duplicate)
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
				{Name: "s3"},
			},
			EntryPoints: []EntryPoint{
				{Owner: "orders", Type: EntryPointHTTP, Name: "list orders"},
			},
		}},
		staticModule{descriptor: Descriptor{
			Name: "s3",
			Kind: CapabilityModule,
			Provides: []Capability{
				{Name: "s3", Status: CapabilityEnabled},
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
