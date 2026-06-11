package mysql

import (
	"testing"

	"github.com/NSObjects/go-template/internal/configs"
	"github.com/NSObjects/go-template/internal/platform/module"
)

func TestModuleReportsDisabledWhenConfigDisabled(t *testing.T) {
	descriptor := New(configs.MysqlConfig{Enabled: false}).Descriptor()
	assertCapabilityStatus(t, descriptor, module.CapabilityDisabled)
}

func TestModuleRejectsEnabledInvalidConfig(t *testing.T) {
	mod := New(configs.MysqlConfig{Enabled: true})
	if err := mod.Validate(); err == nil {
		t.Fatal("Validate() error = nil, want invalid config error")
	}
	assertCapabilityStatus(t, mod.Descriptor(), module.CapabilityUnavailable)
}

func TestModuleReportsEnabledWithValidConfig(t *testing.T) {
	mod := New(configs.MysqlConfig{
		Enabled:  true,
		Host:     "127.0.0.1",
		Port:     "3306",
		User:     "root",
		Database: "app",
	})
	if err := mod.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	assertCapabilityStatus(t, mod.Descriptor(), module.CapabilityEnabled)
}

func assertCapabilityStatus(t *testing.T, descriptor module.Descriptor, want module.CapabilityState) {
	t.Helper()
	if descriptor.Name != Name {
		t.Fatalf("descriptor.Name = %q, want %q", descriptor.Name, Name)
	}
	if len(descriptor.Provides) != 1 {
		t.Fatalf("len(descriptor.Provides) = %d, want 1", len(descriptor.Provides))
	}
	if got := descriptor.Provides[0].Status; got != want {
		t.Fatalf("capability status = %q, want %q", got, want)
	}
}
