// Package app declares the modules included in this application.
package app

import (
	userstorage "github.com/NSObjects/go-template/internal/capabilities/userstorage"
	"github.com/NSObjects/go-template/internal/configs"
	user "github.com/NSObjects/go-template/internal/modules/user"
	"github.com/NSObjects/go-template/internal/platform/module"
)

// ModulesResult contains the explicit modules for assembly.
type ModulesResult struct {
	Modules              []module.Module
	CapabilitySelections []module.CapabilitySelection
}

// Modules returns the capability and business modules explicitly included in the app.
func Modules(cfg configs.Config) (ModulesResult, error) {
	cfg = configs.Normalize(cfg)
	userStorage := userstorage.New(cfg)
	capabilityModules := userStorage.Modules()
	resolvedStorage, err := module.ResolveCapabilityFromModules[user.Repository](
		capabilityModules,
		user.ModuleName,
		user.StorageCapability,
		module.WithCapabilityProviderSelections(cfg.Capabilities.Providers),
	)
	if err != nil {
		return ModulesResult{}, err
	}

	userModule := user.New(user.NewUseCase(resolvedStorage.Value))
	modules := append([]module.Module(nil), capabilityModules...)
	modules = append(modules, userModule)
	return ModulesResult{
		Modules: modules,
		CapabilitySelections: []module.CapabilitySelection{
			{
				Capability: user.StorageCapability,
				Provider:   resolvedStorage.Provider,
			},
		},
	}, nil
}
