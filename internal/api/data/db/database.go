/*
 * Created by lintao on 2023/7/26 下午3:02
 * Copyright © 2020-2024 LINTAO. All rights reserved.
 *
 */

package db

import (
	"fmt"
	"os"

	configs "github.com/NSObjects/go-kit/config"
	kitdb "github.com/NSObjects/go-kit/db"
	"gorm.io/gorm"
)

// NewDatabase 创建数据库连接，使用 go-kit/db 的通用工厂
// 支持 mysql, postgres, sqlite 驱动
func NewDatabase(cfg configs.DatabaseConfig) *gorm.DB {
	// 使用 go-kit/db 创建数据库连接
	db, err := kitdb.NewDatabase(cfg, os.Stdout)
	if err != nil {
		panic(fmt.Errorf("database init [%s]: %w", cfg.Driver, err))
	}

	// 注册业务特定的回调
	if err := db.Callback().Create().After("gorm:create").Register("role:menu_after_create", AfterCreate); err != nil {
		panic(fmt.Errorf("register callback: %w", err))
	}

	return db
}
