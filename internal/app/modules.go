// Package app declares the modules included in this application.
package app

import (
	databasemysql "github.com/NSObjects/go-template/internal/capabilities/database/mysql"
	"github.com/NSObjects/go-template/internal/configs"
	user "github.com/NSObjects/go-template/internal/modules/user"
	usergorm "github.com/NSObjects/go-template/internal/modules/user/repository/gorm"
	usermemory "github.com/NSObjects/go-template/internal/modules/user/repository/memory"
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

	database, err := databaseCapability(cfg)
	if err != nil {
		return ModulesResult{}, err
	}
	repository := userRepository(database.provider)
	userModule := user.New(user.NewUseCase(repository))
	modules := append([]module.Module(nil), database.modules...)
	modules = append(modules, userModule)
	return ModulesResult{
		Modules:              modules,
		CapabilitySelections: database.selections,
	}, nil
}

type databaseRuntime struct {
	provider   databasemysql.DBProvider
	modules    []module.Module
	selections []module.CapabilitySelection
}

func databaseCapability(cfg configs.Config) (databaseRuntime, error) {
	databaseProvider := cfg.Capabilities.Providers[databasemysql.Capability]
	if databaseProvider == "" {
		return databaseRuntime{}, nil
	}

	databaseModules := []module.Module{
		requiredCapabilityModule{
			name:       "app-database",
			capability: databasemysql.Capability,
		},
		databasemysql.New(cfg.Mysql),
	}
	resolvedDB, err := module.ResolveCapabilityFromModules[databasemysql.DBProvider](
		databaseModules[1:],
		user.ModuleName,
		databasemysql.Capability,
		module.WithCapabilityProviderSelections(cfg.Capabilities.Providers),
	)
	if err != nil {
		return databaseRuntime{}, err
	}
	return databaseRuntime{
		provider: resolvedDB.Value,
		modules:  databaseModules,
		selections: []module.CapabilitySelection{
			{Capability: databasemysql.Capability, Provider: resolvedDB.Provider},
		},
	}, nil
}

func userRepository(database databasemysql.DBProvider) user.Repository {
	if database == nil {
		return usermemory.NewRepository()
	}
	return usergorm.New(database)
}

type requiredCapabilityModule struct {
	name       string
	capability string
}

func (m requiredCapabilityModule) Descriptor() module.Descriptor {
	return module.Descriptor{
		Name: m.name,
		Kind: module.PlatformModule,
		Requires: []module.CapabilityRef{
			{Name: m.capability},
		},
	}
}
