/*
 * Created by lintao on 2023/7/18 下午3:56
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package db

import (
	"github.com/NSObjects/go-template/internal/configs"
	_ "github.com/go-sql-driver/mysql"
	redis "github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var Model = fx.Options(
	fx.Provide(NewDataSource),
)

// DataSource 在使用多个db的项目中在DataSource结构体中增加Engine即可
type DataSource struct {
	Mysql   *gorm.DB
	Mongodb *mongo.Database
	Redis   *redis.Client
}

func NewDataSource(cfg configs.Config) *DataSource {
	var dataSource DataSource
	if cfg.Mysql.Host != "" {
		dataSource.Mysql = NewMysql(cfg.Mysql)
	}
	if cfg.Mongodb.Host != "" {
		dataSource.Mongodb = MongoClient(cfg.Mongodb)
	}

	if cfg.Redis.Host != "" {
		dataSource.Redis = NewRedis(cfg.Redis)
	}
	return &dataSource
}
