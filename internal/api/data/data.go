/*
 * Created by lintao on 2023/7/18 下午3:59
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package data

import (
	"fmt"

	"github.com/NSObjects/go-template/internal/api/data/db"
	"github.com/NSObjects/go-template/internal/configs"
	"go.uber.org/fx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Model = fx.Options(
	fx.Provide(db.NewDataSource, NewMysql),
)

func NewMysql(cfg configs.Config) *gorm.DB {
	if cfg.Mysql.Host == "" {
		panic("mysql config is empty")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		cfg.Mysql.User, cfg.Mysql.Password, cfg.Mysql.Host, cfg.Mysql.Port, cfg.Mysql.Database)

	open, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		panic(err)
	}

	err = open.Callback().Create().After("gorm:after_create").Register("role:menu_after_create", AfterCreate)

	if err != nil {
		panic(err)
	}

	return open
}

func AfterCreate(db *gorm.DB) {
	if db.Error == nil &&
		db.Statement.Schema != nil &&
		!db.Statement.SkipHooks &&
		(db.Statement.Schema.AfterCreate || db.Statement.Schema.AfterSave) {
		fmt.Println("BeforeCreate", db.Statement.Schema.Name)
	}
}
