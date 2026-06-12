// Package userstorage provides storage adapters for the user business module.
package userstorage

import (
	usermemory "github.com/NSObjects/go-template/internal/capabilities/userstorage/memory"
	usermysql "github.com/NSObjects/go-template/internal/capabilities/userstorage/mysql"
	"github.com/NSObjects/go-template/internal/configs"
	"github.com/NSObjects/go-template/internal/platform/module"
)

// Modules contains the storage capability modules for the user business module.
type Modules struct {
	modules []module.Module
}

// New creates all user.storage provider modules.
func New(cfg configs.Config) Modules {
	cfg = configs.Normalize(cfg)
	return Modules{
		modules: []module.Module{
			usermemory.New(),
			usermysql.New(cfg.Mysql),
		},
	}
}

// Modules returns the provider modules that participate in app assembly.
func (p Modules) Modules() []module.Module {
	return append([]module.Module(nil), p.modules...)
}
