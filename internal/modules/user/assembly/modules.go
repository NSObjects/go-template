// Package assembly wires the user module with its local capability providers.
package assembly

import (
	"github.com/NSObjects/go-template/internal/api/biz"
	"github.com/NSObjects/go-template/internal/configs"
	user "github.com/NSObjects/go-template/internal/modules/user"
	usermemory "github.com/NSObjects/go-template/internal/modules/user/storage/memory"
	usermysql "github.com/NSObjects/go-template/internal/modules/user/storage/mysql"
	"github.com/NSObjects/go-template/internal/platform/module"
)

// Modules returns user storage providers and the user business module.
func Modules(cfg configs.Config) ([]module.Module, error) {
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
		return nil, err
	}
	userUseCase := biz.NewUserHandler(userRepository)

	return []module.Module{
		storageProviders[0],
		storageProviders[1],
		user.New(userUseCase),
	}, nil
}
