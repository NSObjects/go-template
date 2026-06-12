// Package assembly wires the user module with its local capability providers.
package assembly

import (
	"github.com/NSObjects/go-template/internal/configs"
	user "github.com/NSObjects/go-template/internal/modules/user"
	"github.com/NSObjects/go-template/internal/modules/user/storage"
	"github.com/NSObjects/go-template/internal/platform/module"
)

// Modules returns user storage providers and the user business module.
func Modules(cfg configs.Config) ([]module.Module, error) {
	storageProviders, err := storage.New(cfg)
	if err != nil {
		return nil, err
	}
	userUseCase := user.NewUseCase(storageProviders.Repository())

	modules := storageProviders.Modules()
	modules = append(modules, user.New(userUseCase))
	return modules, nil
}
