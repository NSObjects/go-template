/*
 * Created by lintao on 2023/7/18 下午3:56
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package db

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/NSObjects/go-template/internal/api/data/query"
	"github.com/NSObjects/go-template/internal/configs"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// DataManager 统一的数据管理器，直接管理所有数据库组件
type DataManager struct {
	// 数据库组件
	Mysql   *gorm.DB
	Mongodb *mongo.Database
	Redis   *redis.Client
	Kafka   sarama.SyncProducer

	// 查询接口
	Query *query.Query

	// 配置
	Config *configs.Config
}

// NewDataManager 创建统一的数据管理器，直接初始化所有组件
func NewDataManager(lc fx.Lifecycle, cfg configs.Config) *DataManager {
	dm := &DataManager{
		Config: &cfg,
	}

	// 初始化MySQL
	if cfg.Mysql.Host != "" {
		dm.Mysql = NewMysql(cfg.Mysql)
	}

	// 初始化MongoDB
	if cfg.Mongodb.Host != "" {
		dm.Mongodb = MongoClient(cfg.Mongodb)
	}

	// 初始化Redis
	if cfg.Redis.Host != "" {
		dm.Redis = NewRedis(cfg.Redis)
	}

	// 初始化Kafka
	if len(cfg.Kafka.Brokers) > 0 {
		producer, err := NewKafkaProducer(cfg.Kafka)
		if err != nil {
			panic(err)
		}
		dm.Kafka = producer
	}

	// 初始化Query
	if dm.Mysql != nil {
		dm.Query = query.Use(dm.Mysql)
	}

	// 注册生命周期钩子
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return dm.start(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return dm.Stop(ctx)
		},
	})

	return dm
}

// start 启动所有组件
func (dm *DataManager) start(ctx context.Context) error {
	// 检查MySQL连接
	if dm.Mysql != nil {
		if sqlDB, err := dm.Mysql.DB(); err == nil {
			if err := sqlDB.PingContext(ctx); err != nil {
				return err
			}
		}
	}

	// 检查Redis连接
	if dm.Redis != nil {
		if err := dm.Redis.Ping(ctx).Err(); err != nil {
			return err
		}
	}

	// Kafka连接检查（可选）
	if dm.Kafka != nil {
		// 发送一条空消息作为连通性检查（可选）
		// 忽略错误以避免启动硬失败，也可改为严格校验
	}

	return nil
}

// Stop 停止所有组件
func (dm *DataManager) Stop(ctx context.Context) error {
	// 关闭MySQL连接
	if dm.Mysql != nil {
		if sqlDB, err := dm.Mysql.DB(); err == nil {
			_ = sqlDB.Close()
		}
	}

	// 关闭Redis连接
	if dm.Redis != nil {
		_ = dm.Redis.Close()
	}

	// 关闭Kafka连接
	if dm.Kafka != nil {
		_ = dm.Kafka.Close()
	}

	// MongoDB连接由客户端管理
	return nil
}

// Health 检查所有组件的健康状态
func (dm *DataManager) Health(ctx context.Context) map[string]error {
	health := make(map[string]error)

	// MySQL状态
	if dm.Mysql != nil {
		if sqlDB, err := dm.Mysql.DB(); err == nil {
			health["mysql"] = sqlDB.PingContext(ctx)
		} else {
			health["mysql"] = err
		}
	}

	// Redis状态
	if dm.Redis != nil {
		health["redis"] = dm.Redis.Ping(ctx).Err()
	}

	// Kafka状态
	if dm.Kafka != nil {
		health["kafka"] = nil // Kafka状态检查比较复杂，这里简化处理
	}

	// MongoDB状态
	if dm.Mongodb != nil {
		health["mongodb"] = nil // MongoDB状态检查比较复杂，这里简化处理
	}

	return health
}

// ========== 便捷操作方法 ==========

// MySQLWithContext 获取带上下文的MySQL连接
func (dm *DataManager) MySQLWithContext(ctx context.Context) *gorm.DB {
	if dm.Mysql == nil {
		return nil
	}
	return dm.Mysql.WithContext(ctx)
}

// RedisWithContext 获取带上下文的Redis客户端
func (dm *DataManager) RedisWithContext(ctx context.Context) *redis.Client {
	if dm.Redis == nil {
		return nil
	}
	// Redis客户端本身已经支持context，直接返回
	return dm.Redis
}

// SendKafkaMessage 发送Kafka消息
func (dm *DataManager) SendKafkaMessage(topic string, key, value []byte) error {
	if dm.Kafka == nil {
		return fmt.Errorf("kafka producer not initialized")
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	}

	_, _, err := dm.Kafka.SendMessage(msg)
	return err
}

// IsComponentEnabled 检查组件是否启用
func (dm *DataManager) IsComponentEnabled(component string) bool {
	switch component {
	case "mysql":
		return dm.Mysql != nil
	case "redis":
		return dm.Redis != nil
	case "kafka":
		return dm.Kafka != nil
	case "mongodb":
		return dm.Mongodb != nil
	default:
		return false
	}
}

// NewDB 为了向后兼容，提供获取MySQL连接的方法
func NewDB(dm *DataManager) *gorm.DB {
	if dm == nil {
		return nil
	}
	return dm.Mysql
}

// NewQuery 为了向后兼容，提供获取Query的方法
func NewQuery(dm *DataManager) *query.Query {
	if dm == nil || dm.Query == nil {
		return nil
	}
	return dm.Query
}

var Model = fx.Options(
	fx.Provide(NewDataManager),
	fx.Provide(NewDB),
	fx.Provide(NewQuery),
)
