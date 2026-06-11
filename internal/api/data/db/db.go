/*
 * Created by lintao on 2023/7/18 下午3:56
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/IBM/sarama"
	"github.com/NSObjects/go-template/internal/api/data/query"
	"github.com/NSObjects/go-template/internal/configs"
	"github.com/redis/go-redis/v9"
	"github.com/samber/do/v2"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
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

// NewDataManager 创建统一的数据管理器，仅初始化配置显式启用的组件
func NewDataManager(cfg configs.Config) (*DataManager, error) {
	dm := &DataManager{
		Config: &cfg,
	}
	initialized := false
	defer func() {
		if !initialized {
			_ = dm.Shutdown(context.Background())
		}
	}()

	if err := validateEnabledComponents(cfg); err != nil {
		return nil, err
	}

	// 初始化MySQL
	if cfg.Mysql.Enabled {
		mysqlDB, err := NewMysql(cfg.Mysql)
		if err != nil {
			return nil, fmt.Errorf("initialize mysql: %w", err)
		}
		dm.Mysql = mysqlDB
	}

	// 初始化MongoDB
	if cfg.Mongodb.Enabled {
		mongoDB, err := MongoClient(cfg.Mongodb)
		if err != nil {
			return nil, fmt.Errorf("initialize mongodb: %w", err)
		}
		dm.Mongodb = mongoDB
	}

	// 初始化Redis
	if cfg.Redis.Enabled {
		dm.Redis = NewRedis(cfg.Redis)
	}

	// 初始化Kafka
	if cfg.Kafka.Enabled {
		producer, err := NewKafkaProducer(cfg.Kafka)
		if err != nil {
			return nil, fmt.Errorf("initialize kafka: %w", err)
		}
		dm.Kafka = producer
	}

	// 初始化Query
	if dm.Mysql != nil {
		dm.Query = query.Use(dm.Mysql)
	}

	initialized = true
	return dm, nil
}

func validateEnabledComponents(cfg configs.Config) error {
	if cfg.Mysql.Enabled {
		if isBlank(cfg.Mysql.Host) || isBlank(cfg.Mysql.Port) || isBlank(cfg.Mysql.User) || isBlank(cfg.Mysql.Database) {
			return fmt.Errorf("mysql enabled but host, port, user, or database is empty")
		}
	}
	if cfg.Redis.Enabled {
		if isBlank(cfg.Redis.Host) || isBlank(cfg.Redis.Port) {
			return fmt.Errorf("redis enabled but host or port is empty")
		}
	}
	if cfg.Mongodb.Enabled {
		if isBlank(cfg.Mongodb.Host) || isBlank(cfg.Mongodb.Port) || isBlank(cfg.Mongodb.DataBase) {
			return fmt.Errorf("mongodb enabled but host, port, or database is empty")
		}
	}
	if cfg.Kafka.Enabled && len(cfg.Kafka.Brokers) == 0 {
		return fmt.Errorf("kafka enabled but brokers is empty")
	}
	return nil
}

func isBlank(value string) bool {
	return strings.TrimSpace(value) == ""
}

// Start 启动并验证所有已启用组件。
func (dm *DataManager) Start(ctx context.Context) error {
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

	// 检查MongoDB连接
	if dm.Mongodb != nil {
		if err := dm.Mongodb.Client().Ping(ctx, readpref.Primary()); err != nil {
			return err
		}
	}

	return nil
}

// Shutdown 停止所有组件。
func (dm *DataManager) Shutdown(ctx context.Context) error {
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

	if dm.Mongodb != nil {
		_ = dm.Mongodb.Client().Disconnect(ctx)
	}

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
		health["mongodb"] = dm.Mongodb.Client().Ping(ctx, readpref.Primary())
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

// Register 注册数据库相关依赖。
func Register(i do.Injector) {
	do.Provide(i, func(i do.Injector) (*DataManager, error) {
		cfg, err := do.Invoke[configs.Config](i)
		if err != nil {
			return nil, err
		}
		return NewDataManager(cfg)
	})
	do.Provide(i, func(i do.Injector) (*gorm.DB, error) {
		dm, err := do.Invoke[*DataManager](i)
		if err != nil {
			return nil, err
		}
		return NewDB(dm), nil
	})
	do.Provide(i, func(i do.Injector) (*query.Query, error) {
		dm, err := do.Invoke[*DataManager](i)
		if err != nil {
			return nil, err
		}
		return NewQuery(dm), nil
	})
}
