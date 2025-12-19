/*
 * Created by lintao on 2023/7/26 下午3:51
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package db

import (
	"context"
	"fmt"

	kitdb "github.com/NSObjects/go-kit/db"
	configs "github.com/NSObjects/go-kit/config"
	redis "github.com/redis/go-redis/v9"
)

// NewRedis 创建Redis连接，使用 go-kit/db 的基础功能
func NewRedis(cfg configs.RedisConfig) *redis.Client {
	// 使用 go-kit/db 创建Redis连接
	rdb := kitdb.NewRedis(cfg)

	// 验证连接
	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		panic(fmt.Errorf("redis ping: %w", err))
	}

	return rdb
}
