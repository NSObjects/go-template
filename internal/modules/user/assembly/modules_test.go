package assembly

import (
	"errors"
	"testing"

	"github.com/NSObjects/go-template/internal/configs"
	user "github.com/NSObjects/go-template/internal/modules/user"
	"github.com/NSObjects/go-template/internal/platform/module"
)

func TestModulesIncludesUserWithDefaultStorageProvider(t *testing.T) {
	modules, err := Modules(configs.Config{})
	if err != nil {
		t.Fatalf("Modules() error = %v", err)
	}

	report, err := module.Assemble(modules, module.WithEntryPointAdapters(module.EntryPointHTTP))
	if err != nil {
		t.Fatalf("module.Assemble() error = %v", err)
	}
	requirement, ok := report.Requirement(user.ModuleName, user.StorageCapability)
	if !ok {
		t.Fatalf("report.Requirement(%q, %q) ok = false, want true", user.ModuleName, user.StorageCapability)
	}
	if !requirement.Satisfied || requirement.Provider != "memory" {
		t.Fatalf("requirement = %+v, want satisfied by memory provider", requirement)
	}
}

func TestModulesSelectsConfiguredStorageProvider(t *testing.T) {
	cfg := configs.Config{
		Mysql: configs.MysqlConfig{
			Enabled:  true,
			Host:     "127.0.0.1",
			Port:     "3306",
			User:     "root",
			Database: "app",
		},
		Capabilities: configs.CapabilitiesConfig{
			Providers: map[string]string{
				user.StorageCapability: "mysql",
			},
		},
	}

	modules, err := Modules(cfg)
	if err != nil {
		t.Fatalf("Modules() error = %v", err)
	}

	report, err := module.Assemble(
		modules,
		module.WithEntryPointAdapters(module.EntryPointHTTP),
		module.WithCapabilityProviderSelections(cfg.Capabilities.Providers),
	)
	if err != nil {
		t.Fatalf("module.Assemble() error = %v", err)
	}
	requirement, ok := report.Requirement(user.ModuleName, user.StorageCapability)
	if !ok {
		t.Fatalf("report.Requirement(%q, %q) ok = false, want true", user.ModuleName, user.StorageCapability)
	}
	if !requirement.Satisfied || requirement.Provider != "mysql" {
		t.Fatalf("requirement = %+v, want satisfied by mysql provider", requirement)
	}
}

func TestModulesRejectsInvalidStorageProvider(t *testing.T) {
	cfg := configs.Config{
		Capabilities: configs.CapabilitiesConfig{
			Providers: map[string]string{
				user.StorageCapability: "unknown",
			},
		},
	}

	_, err := Modules(cfg)
	if err == nil {
		t.Fatal("Modules() error = nil, want unavailable capability provider error")
	}

	var unavailable *module.UnavailableCapabilityProviderError
	if !errors.As(err, &unavailable) {
		t.Fatalf("Modules() error = %T, want UnavailableCapabilityProviderError", err)
	}
	if unavailable.Module != user.ModuleName ||
		unavailable.Capability != user.StorageCapability ||
		unavailable.Provider != "unknown" {
		t.Fatalf("UnavailableCapabilityProviderError = %+v, want user/user.storage/unknown", unavailable)
	}
}
