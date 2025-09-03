/*
 * Created by lintao on 2023/7/26 下午3:51
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package db

import (
	"context"

	"github.com/NSObjects/echo-admin/internal/configs"
	redis "github.com/redis/go-redis/v9"
)

func NewRedis(cfg configs.RedisConfig) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Password: cfg.Password, // no password set
		DB:       cfg.Database, // use default DB
	})
	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}
	return rdb
}
