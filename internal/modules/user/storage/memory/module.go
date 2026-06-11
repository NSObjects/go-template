package memory

import (
	"github.com/NSObjects/go-template/internal/api/biz"
	user "github.com/NSObjects/go-template/internal/modules/user"
	"github.com/NSObjects/go-template/internal/platform/module"
)

const (
	ModuleName   = "user-storage-memory"
	ProviderName = "memory"
)

// Module exposes the in-memory user storage provider.
type Module struct {
	repository biz.UserRepository
}

// New creates the default in-memory user storage provider module.
func New() Module {
	return Module{repository: NewRepository()}
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
				Default:  true,
				Value:    m.Repository(),
			},
		},
	}
}

// Repository returns the repository implementation used by the user module.
func (m Module) Repository() biz.UserRepository {
	return m.repository
}
