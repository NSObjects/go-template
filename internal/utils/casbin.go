package utils

import (
	"github.com/NSObjects/go-template/internal/configs"
	"github.com/casbin/casbin/v3"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

// NewCasbinEnforcer 创建Casbin权限控制器
func NewCasbinEnforcer(db *gorm.DB, cfg configs.Config) (*casbin.Enforcer, error) {
	// 使用GORM适配器
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}

	// 创建enforcer
	enforcer, err := casbin.NewEnforcer("configs/rbac_model.conf", adapter)
	if err != nil {
		return nil, err
	}

	// 加载策略
	err = enforcer.LoadPolicy()
	if err != nil {
		return nil, err
	}

	return enforcer, nil
}
