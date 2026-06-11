// Package app declares the modules included in this application.
package app

import (
	"github.com/NSObjects/go-template/internal/capabilities/kafka"
	"github.com/NSObjects/go-template/internal/capabilities/mongodb"
	"github.com/NSObjects/go-template/internal/capabilities/mysql"
	"github.com/NSObjects/go-template/internal/capabilities/redis"
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
	userModules, err := userassembly.Modules(cfg)
	if err != nil {
		return ModulesResult{}, err
	}

	modules := []module.Module{
		mysql.New(cfg.Mysql),
		redis.New(cfg.Redis),
		mongodb.New(cfg.Mongodb),
		kafka.New(cfg.Kafka),
	}
	modules = append(modules, userModules...)

	return ModulesResult{Modules: modules}, nil
}
