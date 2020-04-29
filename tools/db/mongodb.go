/*
 *
 * mongodb.go
 * db
 *
 * Created by lintao on 2020/4/27 10:28 上午
 * Copyright © 2020-2020 LINTAO. All rights reserved.
 *
 */

package db

import (
	"go-template/tools/configs"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoClient() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoUrl()))
	if err != nil {
		panic(err)
	}
	return client
}

func mongoUrl() string {
	uri := "mongodb://"
	if configs.Mgo.Password != "" && configs.Mgo.User != "" {
		uri += configs.Mgo.User + ":" + configs.Mgo.Password + "@"
	}
	return uri + configs.Mgo.Host + ":" + configs.Mgo.Port
}
