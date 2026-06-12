package storage

import (
	"errors"
	"testing"

	"github.com/NSObjects/go-template/internal/configs"
	user "github.com/NSObjects/go-template/internal/modules/user"
	"github.com/NSObjects/go-template/internal/platform/module"
)

func TestNewProvidesDefaultUserStorage(t *testing.T) {
	providers, err := New(configs.Config{})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if providers.Repository() == nil {
		t.Fatal("Repository() = nil, want selected user repository")
	}

	report := assembleUserWithStorage(t, providers, configs.Config{})
	requirement, ok := report.Requirement(user.ModuleName, user.StorageCapability)
	if !ok {
		t.Fatalf("report.Requirement(%q, %q) ok = false, want true", user.ModuleName, user.StorageCapability)
	}
	if !requirement.Satisfied || requirement.Provider != "memory" {
		t.Fatalf("requirement = %+v, want satisfied by memory provider", requirement)
	}
}

func TestNewSelectsConfiguredUserStorageProvider(t *testing.T) {
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

	providers, err := New(cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if providers.Repository() == nil {
		t.Fatal("Repository() = nil, want selected user repository")
	}

	report := assembleUserWithStorage(t, providers, cfg)
	requirement, ok := report.Requirement(user.ModuleName, user.StorageCapability)
	if !ok {
		t.Fatalf("report.Requirement(%q, %q) ok = false, want true", user.ModuleName, user.StorageCapability)
	}
	if !requirement.Satisfied || requirement.Provider != "mysql" {
		t.Fatalf("requirement = %+v, want satisfied by mysql provider", requirement)
	}
}

func TestNewRejectsUnknownUserStorageProvider(t *testing.T) {
	cfg := configs.Config{
		Capabilities: configs.CapabilitiesConfig{
			Providers: map[string]string{
				user.StorageCapability: "unknown",
			},
		},
	}

	_, err := New(cfg)
	if err == nil {
		t.Fatal("New() error = nil, want unavailable capability provider error")
	}

	var unavailable *module.UnavailableCapabilityProviderError
	if !errors.As(err, &unavailable) {
		t.Fatalf("New() error = %T, want UnavailableCapabilityProviderError", err)
	}
	if unavailable.Module != user.ModuleName ||
		unavailable.Capability != user.StorageCapability ||
		unavailable.Provider != "unknown" {
		t.Fatalf("UnavailableCapabilityProviderError = %+v, want user/user.storage/unknown", unavailable)
	}
}

func TestProvidersModulesReturnsCopy(t *testing.T) {
	providers, err := New(configs.Config{})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	modules := providers.Modules()
	if len(modules) == 0 {
		t.Fatal("len(Modules()) = 0, want provider modules")
	}
	modules[0] = nil

	again := providers.Modules()
	if again[0] == nil {
		t.Fatal("Modules() leaked internal slice, want defensive copy")
	}
}

func assembleUserWithStorage(t *testing.T, providers Providers, cfg configs.Config) module.Report {
	t.Helper()

	userModule := user.New(user.NewUseCase(providers.Repository()))
	modules := providers.Modules()
	modules = append(modules, userModule)
	report, err := module.Assemble(
		modules,
		module.WithEntryPointAdapters(module.EntryPointHTTP),
		module.WithCapabilityProviderSelections(cfg.Capabilities.Providers),
	)
	if err != nil {
		t.Fatalf("module.Assemble() error = %v", err)
	}
	return report
}
