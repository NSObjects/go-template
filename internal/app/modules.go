// Package app declares the modules included in this application.
package app

import (
	"github.com/NSObjects/go-template/internal/configs"
	userassembly "github.com/NSObjects/go-template/internal/modules/user/assembly"
	"github.com/NSObjects/go-template/internal/platform/module"
)

// ModulesResult contains the explicit modules for assembly.
type ModulesResult struct {
	Modules []module.Module
}

// Modules returns the capability and business modules explicitly included in the app.
func Modules(cfg configs.Config) (ModulesResult, error) {
	cfg = configs.Normalize(cfg)
	userModules, err := userassembly.Modules(cfg)
	if err != nil {
		return ModulesResult{}, err
	}

	return ModulesResult{Modules: userModules}, nil
}
