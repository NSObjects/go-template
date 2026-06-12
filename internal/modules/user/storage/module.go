// Package storage owns the user.storage provider set and selected repository.
package storage

import (
	"strings"

	"github.com/NSObjects/go-template/internal/configs"
	user "github.com/NSObjects/go-template/internal/modules/user"
	usermemory "github.com/NSObjects/go-template/internal/modules/user/storage/memory"
	usermysql "github.com/NSObjects/go-template/internal/modules/user/storage/mysql"
	"github.com/NSObjects/go-template/internal/platform/module"
)

// Providers contains the user storage provider modules and selected repository.
type Providers struct {
	modules    []module.Module
	repository user.Repository
}

// New creates all user.storage providers and resolves the selected repository.
func New(cfg configs.Config) (Providers, error) {
	selectedProvider := selectedUserStorageProvider(cfg)
	modules := []module.Module{
		usermemory.New(usermemory.WithDefault(selectedProvider == "" || selectedProvider == usermemory.ProviderName)),
		usermysql.New(cfg.Mysql, usermysql.WithDefault(selectedProvider == usermysql.ProviderName)),
	}
	repository, err := module.ResolveCapabilityValueFromModules[user.Repository](
		modules,
		user.ModuleName,
		user.StorageCapability,
		module.WithCapabilityProviderSelections(userStorageSelections(cfg)),
	)
	if err != nil {
		return Providers{}, err
	}
	return Providers{
		modules:    modules,
		repository: repository,
	}, nil
}

func userStorageSelections(cfg configs.Config) map[string]string {
	selections := make(map[string]string, len(cfg.Capabilities.Providers)+1)
	for capability, provider := range cfg.Capabilities.Providers {
		selections[capability] = provider
	}
	if provider := strings.TrimSpace(cfg.User.Storage.Provider); provider != "" {
		selections[user.StorageCapability] = provider
	}
	return selections
}

func selectedUserStorageProvider(cfg configs.Config) string {
	if provider := strings.TrimSpace(cfg.User.Storage.Provider); provider != "" {
		return provider
	}
	return strings.TrimSpace(cfg.Capabilities.Providers[user.StorageCapability])
}

// Modules returns the provider modules that participate in app assembly.
func (p Providers) Modules() []module.Module {
	return append([]module.Module(nil), p.modules...)
}

// Repository returns the repository selected for the user business module.
func (p Providers) Repository() user.Repository {
	return p.repository
}
