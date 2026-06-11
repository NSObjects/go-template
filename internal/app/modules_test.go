package app

import (
	"testing"

	"github.com/NSObjects/go-template/internal/configs"
	user "github.com/NSObjects/go-template/internal/modules/user"
	"github.com/NSObjects/go-template/internal/platform/app"
	"github.com/NSObjects/go-template/internal/platform/module"
)

func TestModulesIncludesUserWithDefaultStorageProvider(t *testing.T) {
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
	requirement, ok := report.Requirement(user.ModuleName, user.StorageCapability)
	if !ok {
		t.Fatalf("report.Requirement(%q, %q) ok = false, want true", user.ModuleName, user.StorageCapability)
	}
	if !requirement.Satisfied || requirement.Provider != "memory" {
		t.Fatalf("requirement = %+v, want satisfied by memory provider", requirement)
	}
	assembled, err := app.Assemble(app.Options{
		Config:  configs.Config{},
		Modules: result.Modules,
	})
	if err != nil {
		t.Fatalf("app.Assemble() error = %v", err)
	}
	if !assembled.Report().HasActiveModule(user.ModuleName) {
		t.Fatal("assembled app does not report active user module")
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
