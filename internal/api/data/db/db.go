/*
 * Created by lintao on 2023/7/18 下午3:56
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package db

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/NSObjects/go-template/internal/configs"
	_ "github.com/go-sql-driver/mysql"
	redis "github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

//var Model = fx.Options(
//	fx.Provide(NewDataSource),
//)

// DataSource 统一的数据源，包含所有数据库组件
type DataSource struct {
	Mysql   *gorm.DB
	Mongodb *mongo.Database
	Redis   *redis.Client
	Kafka   sarama.SyncProducer
}

// ComponentStatus 组件状态
type ComponentStatus struct {
	Enabled bool
	Error   error
}

// NewDataSource 根据配置初始化数据源
func NewDataSource(lc fx.Lifecycle, cfg configs.Config) *DataSource {
	ds := &DataSource{}

	// 初始化MySQL
	if cfg.Mysql.Host != "" {
		ds.Mysql = NewMysql(cfg.Mysql)
	}

	// 初始化MongoDB
	if cfg.Mongodb.Host != "" {
		ds.Mongodb = MongoClient(cfg.Mongodb)
	}

	// 初始化Redis
	if cfg.Redis.Host != "" {
		ds.Redis = NewRedis(cfg.Redis)
	}

	// 初始化Kafka
	if len(cfg.Kafka.Brokers) > 0 {
		producer, err := NewKafkaProducer(cfg.Kafka)
		if err != nil {
			panic(err)
		}
		ds.Kafka = producer
	}

	// 注册生命周期钩子
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return ds.start(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return ds.Stop(ctx)
		},
	})

	return ds
}

// start 启动所有组件
func (ds *DataSource) start(ctx context.Context) error {
	// 检查MySQL连接
	if ds.Mysql != nil {
		if sqlDB, err := ds.Mysql.DB(); err == nil {
			if err := sqlDB.PingContext(ctx); err != nil {
				return err
			}
		}
	}

	// 检查Redis连接
	if ds.Redis != nil {
		if err := ds.Redis.Ping(ctx).Err(); err != nil {
			return err
		}
	}

	// Kafka连接检查（可选）
	if ds.Kafka != nil {
		// 发送一条空消息作为连通性检查（可选）
		// 忽略错误以避免启动硬失败，也可改为严格校验
	}

	return nil
}

// Stop 停止所有组件
func (ds *DataSource) Stop(ctx context.Context) error {
	// 关闭MySQL连接
	if ds.Mysql != nil {
		if sqlDB, err := ds.Mysql.DB(); err == nil {
			_ = sqlDB.Close()
		}
	}

	// 关闭Redis连接
	if ds.Redis != nil {
		_ = ds.Redis.Close()
	}

	// 关闭Kafka连接
	if ds.Kafka != nil {
		_ = ds.Kafka.Close()
	}

	// MongoDB连接由客户端管理
	return nil
}

// GetComponentStatus 获取所有组件状态
func (ds *DataSource) GetComponentStatus(ctx context.Context) map[string]ComponentStatus {
	status := make(map[string]ComponentStatus)

	// MySQL状态
	if ds.Mysql != nil {
		if sqlDB, err := ds.Mysql.DB(); err == nil {
			status["mysql"] = ComponentStatus{
				Enabled: true,
				Error:   sqlDB.PingContext(ctx),
			}
		} else {
			status["mysql"] = ComponentStatus{
				Enabled: true,
				Error:   err,
			}
		}
	} else {
		status["mysql"] = ComponentStatus{Enabled: false}
	}

	// Redis状态
	if ds.Redis != nil {
		status["redis"] = ComponentStatus{
			Enabled: true,
			Error:   ds.Redis.Ping(ctx).Err(),
		}
	} else {
		status["redis"] = ComponentStatus{Enabled: false}
	}

	// Kafka状态
	if ds.Kafka != nil {
		status["kafka"] = ComponentStatus{
			Enabled: true,
			Error:   nil, // Kafka状态检查比较复杂，这里简化处理
		}
	} else {
		status["kafka"] = ComponentStatus{Enabled: false}
	}

	// MongoDB状态
	if ds.Mongodb != nil {
		status["mongodb"] = ComponentStatus{
			Enabled: true,
			Error:   nil, // MongoDB状态检查比较复杂，这里简化处理
		}
	} else {
		status["mongodb"] = ComponentStatus{Enabled: false}
	}

	return status
}
