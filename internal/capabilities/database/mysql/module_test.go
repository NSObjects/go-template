package mysql

import (
	"context"
	"errors"
	"testing"

	"github.com/NSObjects/go-template/internal/configs"
	"github.com/NSObjects/go-template/internal/platform/module"
	"gorm.io/gorm"
)

func TestModuleReportsDatabaseProviderStatus(t *testing.T) {
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
			cfg:  validMysqlConfig(),
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
			if capability.Name != Capability {
				t.Fatalf("capability.Name = %q, want %q", capability.Name, Capability)
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
			if _, ok := capability.Value.(DBProvider); !ok {
				t.Fatalf("capability.Value = %T, want DBProvider", capability.Value)
			}
		})
	}
}

func TestModuleDefaultCanBeEnabled(t *testing.T) {
	mod := New(validMysqlConfig(), WithDefault(true))
	capability := mod.Descriptor().Provides[0]

	if !capability.Default {
		t.Fatal("capability.Default = false, want true")
	}
}

func TestModuleExposesStableDBProvider(t *testing.T) {
	mod := New(configs.MysqlConfig{})
	provider := mod.Provider()

	if mod.Provider() != provider {
		t.Fatal("Provider() returned a different instance, want stable provider owned by module")
	}
	if mod.Descriptor().Provides[0].Value != provider {
		t.Fatal("Descriptor() capability value is not the module-owned provider")
	}
}

func TestModuleStartInitializesDatabaseRuntimeOnce(t *testing.T) {
	mod := New(validMysqlConfig())
	db := &gorm.DB{}
	var factoryCalls int
	mod.provider.runtimeFactory = func(_ context.Context, cfg configs.MysqlConfig) (*runtime, error) {
		factoryCalls++
		if cfg != validMysqlConfig() {
			t.Fatalf("runtime factory cfg = %+v, want valid mysql config", cfg)
		}
		return &runtime{db: db}, nil
	}

	if err := mod.Start(context.Background()); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	if err := mod.Start(context.Background()); err != nil {
		t.Fatalf("second Start() error = %v", err)
	}
	gotDB, err := mod.Provider().DB(context.Background())
	if err != nil {
		t.Fatalf("Provider().DB() error = %v", err)
	}
	if gotDB != db {
		t.Fatal("Provider().DB() returned a different database")
	}
	if factoryCalls != 1 {
		t.Fatalf("runtime factory calls = %d, want 1", factoryCalls)
	}
}

func TestModuleStopReleasesDatabaseRuntime(t *testing.T) {
	mod := New(validMysqlConfig())
	var shutdownCalls int
	mod.provider.runtimeFactory = func(context.Context, configs.MysqlConfig) (*runtime, error) {
		return &runtime{
			db: &gorm.DB{},
			shutdown: func(context.Context) error {
				shutdownCalls++
				return nil
			},
		}, nil
	}

	if err := mod.Start(context.Background()); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	if err := mod.Stop(context.Background()); err != nil {
		t.Fatalf("Stop() error = %v", err)
	}
	if err := mod.Stop(context.Background()); err != nil {
		t.Fatalf("second Stop() error = %v", err)
	}
	if shutdownCalls != 1 {
		t.Fatalf("shutdown calls = %d, want 1", shutdownCalls)
	}
}

func TestModuleStartCanRetryAfterRuntimeFailure(t *testing.T) {
	mod := New(validMysqlConfig())
	startErr := errors.New("mysql unavailable")
	var factoryCalls int
	mod.provider.runtimeFactory = func(context.Context, configs.MysqlConfig) (*runtime, error) {
		factoryCalls++
		if factoryCalls == 1 {
			return nil, startErr
		}
		return &runtime{db: &gorm.DB{}}, nil
	}

	if err := mod.Start(context.Background()); !errors.Is(err, startErr) {
		t.Fatalf("first Start() error = %v, want start error", err)
	}
	if err := mod.Start(context.Background()); err != nil {
		t.Fatalf("second Start() error = %v, want retry success", err)
	}
	if factoryCalls != 2 {
		t.Fatalf("runtime factory calls = %d, want 2", factoryCalls)
	}
}

func TestValidateIgnoresDisabledConfig(t *testing.T) {
	mod := New(configs.MysqlConfig{Enabled: false})

	if err := mod.Validate(); err != nil {
		t.Fatalf("Validate() error = %v, want nil for disabled config", err)
	}
}

func TestValidateRejectsEnabledInvalidConfig(t *testing.T) {
	mod := New(configs.MysqlConfig{Enabled: true})

	err := mod.Validate()
	if err == nil {
		t.Fatal("Validate() error = nil, want invalid config error")
	}
	if !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("Validate() error = %v, want ErrInvalidConfig", err)
	}
}

func validMysqlConfig() configs.MysqlConfig {
	return configs.MysqlConfig{
		Enabled:  true,
		Host:     "127.0.0.1",
		Port:     "3306",
		User:     "root",
		Database: "app",
	}
}
