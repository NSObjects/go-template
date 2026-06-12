package userstorage

import (
	"testing"

	"github.com/NSObjects/go-template/internal/configs"
	user "github.com/NSObjects/go-template/internal/modules/user"
	"github.com/NSObjects/go-template/internal/platform/module"
)

func TestNewProvidesUserStorageCapabilityModules(t *testing.T) {
	providers := New(configs.Config{})

	report, err := module.Assemble(providers.Modules())
	if err != nil {
		t.Fatalf("module.Assemble() error = %v", err)
	}

	memory, ok := report.Capability(user.StorageCapability)
	if !ok {
		t.Fatalf("report.Capability(%q) ok = false, want true", user.StorageCapability)
	}
	if memory.Provider != "memory" || !memory.Default {
		t.Fatalf("memory capability = %+v, want default memory provider", memory)
	}
}

func TestModulesReturnsCopy(t *testing.T) {
	providers := New(configs.Config{})

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
