/*
 * Created by lintao on 2023/7/18 下午3:56
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package db

import (
	"github.com/5xxxx/pie"
	"github.com/5xxxx/pie/driver"
	"github.com/NSObjects/go-template/internal/configs"

	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoClient(cfg configs.Mongodb) driver.Client {

	uri := "mongodb://"
	if cfg.Password != "" && cfg.User != "" {
		uri += cfg.User + ":" + cfg.Password + "@"
	}
	uri += cfg.Host + ":" + cfg.Port

	newClient, err := pie.NewClient(cfg.DataBase, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	return newClient
}
