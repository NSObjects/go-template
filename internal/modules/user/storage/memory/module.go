package memory

import (
	user "github.com/NSObjects/go-template/internal/modules/user"
	"github.com/NSObjects/go-template/internal/platform/module"
)

const (
	ModuleName   = "user-storage-memory"
	ProviderName = "memory"
)

// Module exposes the in-memory user storage provider.
type Module struct {
	repository      user.Repository
	defaultProvider bool
}

// Option customizes the memory user storage provider.
type Option func(*Module)

// WithDefault controls whether memory is the default user.storage provider.
func WithDefault(defaultProvider bool) Option {
	return func(m *Module) {
		m.defaultProvider = defaultProvider
	}
}

// New creates the default in-memory user storage provider module.
func New(options ...Option) Module {
	m := Module{
		repository:      NewRepository(),
		defaultProvider: true,
	}
	for _, option := range options {
		option(&m)
	}
	return m
}

// Descriptor reports memory as the enabled default provider for user.storage.
func (m Module) Descriptor() module.Descriptor {
	return module.Descriptor{
		Name: ModuleName,
		Kind: module.CapabilityModule,
		Provides: []module.Capability{
			{
				Name:     user.StorageCapability,
				Provider: ProviderName,
				Status:   module.CapabilityEnabled,
				Default:  m.defaultProvider,
				Value:    m.Repository(),
			},
		},
	}
}

// Repository returns the repository implementation used by the user module.
func (m Module) Repository() user.Repository {
	return m.repository
}
