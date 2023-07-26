/*
 * Created by lintao on 2023/7/18 下午3:56
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package db

import (
	"github.com/NSObjects/go-template/internal/configs"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoClient(cfg configs.Mongodb) *mongo.Database {

	uri := "mongodb://"
	if cfg.Password != "" && cfg.User != "" {
		uri += cfg.User + ":" + cfg.Password + "@"
	}
	uri += cfg.Host + ":" + cfg.Port

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	return client.Database(cfg.DataBase)
}
