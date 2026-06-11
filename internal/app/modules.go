// Package app declares the modules included in this application.
package app

import (
	"github.com/NSObjects/go-template/internal/api/biz"
	"github.com/NSObjects/go-template/internal/capabilities/kafka"
	"github.com/NSObjects/go-template/internal/capabilities/mongodb"
	"github.com/NSObjects/go-template/internal/capabilities/mysql"
	"github.com/NSObjects/go-template/internal/capabilities/redis"
	"github.com/NSObjects/go-template/internal/configs"
	user "github.com/NSObjects/go-template/internal/modules/user"
	usermemory "github.com/NSObjects/go-template/internal/modules/user/storage/memory"
	usermysql "github.com/NSObjects/go-template/internal/modules/user/storage/mysql"
	"github.com/NSObjects/go-template/internal/platform/module"
)

// ModulesResult contains the explicit modules and capability selections for assembly.
type ModulesResult struct {
	Modules              []module.Module
	CapabilitySelections []module.CapabilitySelection
}

// Modules returns the capability and business modules explicitly included in the app.
func Modules(cfg configs.Config) (ModulesResult, error) {
	memoryProvider := usermemory.New()
	mysqlProvider := usermysql.New(cfg.Mysql)
	userRepository := userRepositoryFor(cfg, memoryProvider, mysqlProvider)
	userUseCase := biz.NewUserHandler(userRepository)

	return ModulesResult{
		Modules: []module.Module{
			mysql.New(cfg.Mysql),
			redis.New(cfg.Redis),
			mongodb.New(cfg.Mongodb),
			kafka.New(cfg.Kafka),
			memoryProvider,
			mysqlProvider,
			user.New(userUseCase),
		},
		CapabilitySelections: capabilitySelections(cfg),
	}, nil
}

func capabilitySelections(cfg configs.Config) []module.CapabilitySelection {
	selections := make([]module.CapabilitySelection, 0, len(cfg.Capabilities.Providers))
	for capability, provider := range cfg.Capabilities.Providers {
		selections = append(selections, module.CapabilitySelection{
			Capability: capability,
			Provider:   provider,
		})
	}
	return selections
}

func userRepositoryFor(
	cfg configs.Config,
	memoryProvider usermemory.Module,
	mysqlProvider usermysql.Module,
) biz.UserRepository {
	if cfg.Capabilities.Providers[user.StorageCapability] == usermysql.ProviderName {
		return mysqlProvider.Repository()
	}
	return memoryProvider.Repository()
}
