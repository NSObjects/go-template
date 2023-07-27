/*
 * Created by lintao on 2023/7/26 下午3:02
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package db

import (
	"fmt"
	"github.com/NSObjects/go-template/internal/api/data/model"

	"github.com/NSObjects/go-template/internal/configs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMysql(cfg configs.MysqlConfig) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		fmt.Println(err)
	}
	return db
}
