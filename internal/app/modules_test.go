package app

import (
	"testing"

	databasemysql "github.com/NSObjects/go-template/internal/capabilities/database/mysql"
	"github.com/NSObjects/go-template/internal/configs"
	user "github.com/NSObjects/go-template/internal/modules/user"
	"github.com/NSObjects/go-template/internal/platform/app"
	"github.com/NSObjects/go-template/internal/platform/module"
)

func TestModulesIncludesUserWithDefaultMemoryRepository(t *testing.T) {
	result, err := Modules(configs.Config{})
	if err != nil {
		t.Fatalf("Modules() error = %v", err)
	}

	report, err := module.Assemble(result.Modules, module.WithEntryPointAdapters(module.EntryPointHTTP))
	if err != nil {
		t.Fatalf("module.Assemble() error = %v", err)
	}
	if !report.HasActiveModule(user.ModuleName) {
		t.Fatal("user module is not active")
	}
	if report.HasActiveModule(databasemysql.ModuleName) {
		t.Fatal("database module is active for default memory repository")
	}
	assembled, err := app.Assemble(app.Options{
		Config:               configs.Config{},
		Modules:              result.Modules,
		CapabilitySelections: result.CapabilitySelections,
	})
	if err != nil {
		t.Fatalf("app.Assemble() error = %v", err)
	}
	if !assembled.Report().HasActiveModule(user.ModuleName) {
		t.Fatal("assembled app does not report active user module")
	}
}

func TestModulesDoNotIncludeUnusedGenericDatabaseCapabilities(t *testing.T) {
	result, err := Modules(configs.Config{})
	if err != nil {
		t.Fatalf("Modules() error = %v", err)
	}

	report, err := module.Assemble(result.Modules, module.WithEntryPointAdapters(module.EntryPointHTTP))
	if err != nil {
		t.Fatalf("module.Assemble() error = %v", err)
	}

	for _, name := range []string{"mysql", "redis", "mongodb", "kafka"} {
		if report.HasActiveModule(name) {
			t.Fatalf("%s capability module is active without a business requirement", name)
		}
	}
}

func TestModulesSelectsUserGormRepositoryWithDatabaseCapability(t *testing.T) {
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
				databasemysql.Capability: "mysql",
			},
		},
	}

	result, err := Modules(cfg)
	if err != nil {
		t.Fatalf("Modules() error = %v", err)
	}

	report, err := assembleAppReport(t, result)
	if err != nil {
		t.Fatalf("assembleAppReport() error = %v", err)
	}
	requirement, ok := report.Requirement("app-database", databasemysql.Capability)
	if !ok {
		t.Fatalf("report.Requirement(%q, %q) ok = false, want true", "app-database", databasemysql.Capability)
	}
	if !requirement.Satisfied || requirement.Provider != "mysql" {
		t.Fatalf("requirement = %+v, want satisfied by mysql database provider", requirement)
	}
}

func TestModulesResultCarriesResolvedCapabilitySelection(t *testing.T) {
	cfg := configs.Config{}

	result, err := Modules(cfg)
	if err != nil {
		t.Fatalf("Modules() error = %v", err)
	}

	assembled, err := app.Assemble(app.Options{
		Config: configs.Config{
			Capabilities: configs.CapabilitiesConfig{
				Providers: map[string]string{
					databasemysql.Capability: "mysql",
				},
			},
		},
		Modules:              result.Modules,
		CapabilitySelections: result.CapabilitySelections,
	})
	if err != nil {
		t.Fatalf("app.Assemble() error = %v", err)
	}

	if assembled.Report().HasActiveModule(databasemysql.ModuleName) {
		t.Fatal("database module became active despite composition root selecting memory repository")
	}
}

func TestModulesResultCarriesSelectedDatabaseProvider(t *testing.T) {
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
				databasemysql.Capability: "mysql",
			},
		},
	}

	result, err := Modules(cfg)
	if err != nil {
		t.Fatalf("Modules() error = %v", err)
	}
	if len(result.CapabilitySelections) != 1 {
		t.Fatalf("len(CapabilitySelections) = %d, want 1", len(result.CapabilitySelections))
	}
	if result.CapabilitySelections[0] != (module.CapabilitySelection{Capability: databasemysql.Capability, Provider: "mysql"}) {
		t.Fatalf("CapabilitySelections[0] = %+v, want database.gorm/mysql", result.CapabilitySelections[0])
	}
}

func TestExplicitAppModulesActivateUserNotLegacyFiles(t *testing.T) {
	result, err := Modules(configs.Config{})
	if err != nil {
		t.Fatalf("Modules() error = %v", err)
	}
	report, err := module.Assemble(result.Modules, module.WithEntryPointAdapters(module.EntryPointHTTP))
	if err != nil {
		t.Fatalf("module.Assemble() explicit app modules error = %v", err)
	}
	if !report.HasActiveModule(user.ModuleName) {
		t.Fatal("explicit app modules did not activate user module")
	}

	customReport, err := module.Assemble(nil, module.WithEntryPointAdapters(module.EntryPointHTTP))
	if err != nil {
		t.Fatalf("module.Assemble() custom empty modules error = %v", err)
	}
	if customReport.HasActiveModule(user.ModuleName) {
		t.Fatal("custom assembly activated user from legacy files without explicit module inclusion")
	}
	if len(customReport.EntryPoints) != 0 {
		t.Fatalf("len(customReport.EntryPoints) = %d, want 0", len(customReport.EntryPoints))
	}
}

func assembleAppReport(t *testing.T, result ModulesResult) (module.Report, error) {
	t.Helper()

	return module.Assemble(
		result.Modules,
		module.WithEntryPointAdapters(module.EntryPointHTTP),
		module.WithCapabilitySelections(result.CapabilitySelections...),
	)
}
