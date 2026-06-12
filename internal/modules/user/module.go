package user

import (
	"github.com/NSObjects/go-template/internal/platform/module"
)

const (
	ModuleName        = "user"
	StorageCapability = "user.storage"
)

// Module declares the user business module for application assembly.
type Module struct {
	useCase UseCase
}

// New creates a user business module backed by the provided use case.
func New(useCase UseCase) Module {
	return Module{useCase: useCase}
}

// Descriptor returns the user module's assembly-time requirements and entry points.
func (m Module) Descriptor() module.Descriptor {
	return module.Descriptor{
		Name: ModuleName,
		Kind: module.BusinessModule,
		Requires: []module.CapabilityRef{
			{Name: StorageCapability},
		},
		EntryPoints: HTTPEntryPoints(ModuleName, m.useCase),
	}
}
