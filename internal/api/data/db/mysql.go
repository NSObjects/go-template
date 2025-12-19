/*
 * Created by lintao on 2023/7/26 下午3:02
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package db

import (
	"fmt"
	"os"

	kitdb "github.com/NSObjects/go-kit/db"
	configs "github.com/NSObjects/go-kit/config"
	"gorm.io/gorm"
)

// NewMysql 创建MySQL连接，使用 go-kit/db 的基础功能
// 并添加业务特定的回调
func NewMysql(cfg configs.MysqlConfig) *gorm.DB {
	// 使用 go-kit/db 创建MySQL连接
	db, err := kitdb.NewMySQL(cfg, os.Stdout)
	if err != nil {
		panic(fmt.Errorf("mysql init: %w", err))
	}

	// 注册业务特定的回调
	err = db.Callback().Create().After("gorm:create").Register("role:menu_after_create", AfterCreate)
	if err != nil {
		panic(fmt.Errorf("register callback: %w", err))
	}

	return db
}
