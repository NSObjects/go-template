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

// ModulesResult contains the explicit modules for assembly.
type ModulesResult struct {
	Modules []module.Module
}

// Modules returns the capability and business modules explicitly included in the app.
func Modules(cfg configs.Config) (ModulesResult, error) {
	memoryProvider := usermemory.New()
	mysqlProvider := usermysql.New(cfg.Mysql)
	storageProviders := []module.Module{
		memoryProvider,
		mysqlProvider,
	}
	userRepository, err := module.ResolveCapabilityValueFromModules[biz.UserRepository](
		storageProviders,
		user.ModuleName,
		user.StorageCapability,
		module.WithCapabilityProviderSelections(cfg.Capabilities.Providers),
	)
	if err != nil {
		return ModulesResult{}, err
	}
	userUseCase := biz.NewUserHandler(userRepository)

	return ModulesResult{
		Modules: []module.Module{
			mysql.New(cfg.Mysql),
			redis.New(cfg.Redis),
			mongodb.New(cfg.Mongodb),
			kafka.New(cfg.Kafka),
			storageProviders[0],
			storageProviders[1],
			user.New(userUseCase),
		},
	}, nil
}
