package mysql

import (
	"testing"

	"github.com/NSObjects/go-template/internal/api/biz"
	"github.com/NSObjects/go-template/internal/configs"
	user "github.com/NSObjects/go-template/internal/modules/user"
	"github.com/NSObjects/go-template/internal/platform/module"
)

func TestModuleReportsUserStorageProviderStatus(t *testing.T) {
	tests := []struct {
		name string
		cfg  configs.MysqlConfig
		want module.CapabilityState
	}{
		{
			name: "disabled config cannot satisfy selected mysql provider",
			cfg:  configs.MysqlConfig{Enabled: false},
			want: module.CapabilityDisabled,
		},
		{
			name: "enabled invalid config is unavailable",
			cfg:  configs.MysqlConfig{Enabled: true},
			want: module.CapabilityUnavailable,
		},
		{
			name: "enabled valid config can satisfy selected mysql provider",
			cfg: configs.MysqlConfig{
				Enabled:  true,
				Host:     "127.0.0.1",
				Port:     "3306",
				User:     "root",
				Database: "app",
			},
			want: module.CapabilityEnabled,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mod := New(tt.cfg)
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
			if capability.Status != tt.want {
				t.Fatalf("capability.Status = %q, want %q", capability.Status, tt.want)
			}
			if capability.Default {
				t.Fatal("capability.Default = true, want false")
			}
		})
	}
}

func TestModuleExposesUserRepository(t *testing.T) {
	var _ biz.UserRepository = New(configs.MysqlConfig{}).Repository()
}
