package mysql

import (
	"context"
	"errors"
	"testing"

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
			if _, ok := capability.Value.(user.Repository); !ok {
				t.Fatalf("capability.Value = %T, want user.Repository", capability.Value)
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

func TestModuleExposesUserRepository(t *testing.T) {
	var _ user.Repository = New(configs.MysqlConfig{}).Repository()
}

func TestModuleExposesStableUserRepository(t *testing.T) {
	mod := New(configs.MysqlConfig{})
	repository := mod.Repository()

	if mod.Repository() != repository {
		t.Fatal("Repository() returned a different instance, want stable repository owned by module")
	}
	if mod.Descriptor().Provides[0].Value != repository {
		t.Fatal("Descriptor() capability value is not the module-owned repository")
	}
}

func TestModuleStartInitializesRepositoryRuntimeOnce(t *testing.T) {
	mod := New(validMysqlConfig())
	repository := &fakeUserRepository{}
	var factoryCalls int
	mod.repository.runtimeFactory = func(_ context.Context, cfg configs.MysqlConfig) (*storageRuntime, error) {
		factoryCalls++
		if cfg != validMysqlConfig() {
			t.Fatalf("runtime factory cfg = %+v, want valid mysql config", cfg)
		}
		return &storageRuntime{repository: repository}, nil
	}

	if err := mod.Start(context.Background()); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	if err := mod.Start(context.Background()); err != nil {
		t.Fatalf("second Start() error = %v", err)
	}
	if _, _, err := mod.Repository().ListUsers(context.Background(), user.ListUsersRequest{}); err != nil {
		t.Fatalf("Repository().ListUsers() error = %v", err)
	}

	if factoryCalls != 1 {
		t.Fatalf("runtime factory calls = %d, want 1", factoryCalls)
	}
	if repository.listUsersCalls != 1 {
		t.Fatalf("fake repository ListUsers calls = %d, want 1", repository.listUsersCalls)
	}
}

func TestModuleStopReleasesRepositoryRuntime(t *testing.T) {
	mod := New(validMysqlConfig())
	var shutdownCalls int
	mod.repository.runtimeFactory = func(context.Context, configs.MysqlConfig) (*storageRuntime, error) {
		return &storageRuntime{
			repository: &fakeUserRepository{},
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
	mod.repository.runtimeFactory = func(context.Context, configs.MysqlConfig) (*storageRuntime, error) {
		factoryCalls++
		if factoryCalls == 1 {
			return nil, startErr
		}
		return &storageRuntime{repository: &fakeUserRepository{}}, nil
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

type fakeUserRepository struct {
	listUsersCalls int
}

func (r *fakeUserRepository) ListUsers(context.Context, user.ListUsersRequest) ([]user.ListItem, int64, error) {
	r.listUsersCalls++
	return nil, 0, nil
}

func (r *fakeUserRepository) Create(context.Context, user.CreateRequest) error {
	return nil
}

func (r *fakeUserRepository) GetByID(context.Context, int64) (user.Data, error) {
	return user.Data{}, nil
}

func (r *fakeUserRepository) Update(context.Context, int64, user.UpdateRequest) error {
	return nil
}

func (r *fakeUserRepository) Delete(context.Context, int64) error {
	return nil
}
