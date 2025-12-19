/*
 * Created by lintao on 2023/7/18 下午3:56
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package db

import (
	"context"
	"fmt"

	configs "github.com/NSObjects/go-kit/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoClient(cfg configs.MongoConfig) *mongo.Database {
	uri := "mongodb://"
	if cfg.Password != "" && cfg.User != "" {
		uri += cfg.User + ":" + cfg.Password + "@"
	}
	uri += fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	if err := client.Connect(context.Background()); err != nil {
		panic(err)
	}
	return client.Database(cfg.Database)
}
