package memory

import (
	"testing"

	user "github.com/NSObjects/go-template/internal/modules/user"
	"github.com/NSObjects/go-template/internal/platform/module"
)

func TestModuleProvidesDefaultUserStorageCapability(t *testing.T) {
	mod := New()
	descriptor := mod.Descriptor()

	if descriptor.Name != ModuleName {
		t.Fatalf("descriptor.Name = %q, want %q", descriptor.Name, ModuleName)
	}
	if descriptor.Kind != module.CapabilityModule {
		t.Fatalf("descriptor.Kind = %q, want capability", descriptor.Kind)
	}
	if len(descriptor.Provides) != 1 {
		t.Fatalf("len(descriptor.Provides) = %d, want 1", len(descriptor.Provides))
	}

	capability := descriptor.Provides[0]
	if capability.Name != user.StorageCapability {
		t.Fatalf("capability.Name = %q, want %q", capability.Name, user.StorageCapability)
	}
	if capability.Provider != ProviderName {
		t.Fatalf("capability.Provider = %q, want %q", capability.Provider, ProviderName)
	}
	if capability.Status != module.CapabilityEnabled {
		t.Fatalf("capability.Status = %q, want enabled", capability.Status)
	}
	if !capability.Default {
		t.Fatal("capability.Default = false, want true")
	}
	if _, ok := capability.Value.(user.Repository); !ok {
		t.Fatalf("capability.Value = %T, want user.Repository", capability.Value)
	}
	if mod.Repository() == nil {
		t.Fatal("Repository() = nil, want default repository")
	}
}

func TestModuleDefaultCanBeDisabled(t *testing.T) {
	mod := New(WithDefault(false))
	capability := mod.Descriptor().Provides[0]

	if capability.Default {
		t.Fatal("capability.Default = true, want false")
	}
}
