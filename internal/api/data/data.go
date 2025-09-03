/*
 * Created by lintao on 2023/7/18 下午3:59
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package data

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/NSObjects/go-template/internal/api/data/db"
	"github.com/NSObjects/go-template/internal/api/data/query"
	"github.com/NSObjects/go-template/internal/configs"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// DataManager 统一的数据管理器，提供所有数据库组件的操作接口
type DataManager struct {
	// 数据源
	ds *db.DataSource

	// 查询接口
	Query *query.Query

	// 配置
	Config *configs.Config
}

// NewDataManager 创建统一的数据管理器
func NewDataManager(
	ds *db.DataSource,
	query *query.Query,
	cfg configs.Config,
) *DataManager {
	return &DataManager{
		ds:     ds,
		Query:  query,
		Config: &cfg,
	}
}

// Close 关闭所有数据库连接
func (dm *DataManager) Close() error {
	// 通过DataSource统一关闭所有连接
	return dm.ds.Stop(context.Background())
}

// Health 检查所有组件的健康状态
func (dm *DataManager) Health(ctx context.Context) map[string]error {
	health := make(map[string]error)

	// 通过DataSource获取组件状态
	status := dm.ds.GetComponentStatus(ctx)
	for component, status := range status {
		if status.Enabled {
			health[component] = status.Error
		}
	}

	return health
}

// ========== 统一的数据操作接口 ==========

// MySQL 获取MySQL数据库连接
func (dm *DataManager) MySQL() *gorm.DB {
	return dm.ds.Mysql
}

// Redis 获取Redis客户端
func (dm *DataManager) Redis() *redis.Client {
	return dm.ds.Redis
}

// Kafka 获取Kafka生产者
func (dm *DataManager) Kafka() sarama.SyncProducer {
	return dm.ds.Kafka
}

// MongoDB 获取MongoDB数据库
func (dm *DataManager) MongoDB() *mongo.Database {
	return dm.ds.Mongodb
}

// ========== 便捷操作方法 ==========

// MySQLWithContext 获取带上下文的MySQL连接
func (dm *DataManager) MySQLWithContext(ctx context.Context) *gorm.DB {
	if dm.ds.Mysql == nil {
		return nil
	}
	return dm.ds.Mysql.WithContext(ctx)
}

// RedisWithContext 获取带上下文的Redis客户端
func (dm *DataManager) RedisWithContext(ctx context.Context) *redis.Client {
	if dm.ds.Redis == nil {
		return nil
	}
	// Redis客户端本身已经支持context，直接返回
	return dm.ds.Redis
}

// SendKafkaMessage 发送Kafka消息
func (dm *DataManager) SendKafkaMessage(topic string, key, value []byte) error {
	if dm.ds.Kafka == nil {
		return fmt.Errorf("kafka producer not initialized")
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	}

	_, _, err := dm.ds.Kafka.SendMessage(msg)
	return err
}

// IsComponentEnabled 检查组件是否启用
func (dm *DataManager) IsComponentEnabled(component string) bool {
	switch component {
	case "mysql":
		return dm.ds.Mysql != nil
	case "redis":
		return dm.ds.Redis != nil
	case "kafka":
		return dm.ds.Kafka != nil
	case "mongodb":
		return dm.ds.Mongodb != nil
	default:
		return false
	}
}

var Model = fx.Options(
	fx.Provide(
		db.NewDataSource,
		NewDB,
		NewQuery,
		NewDataManager,
	),
)

func NewQuery(db *gorm.DB) *query.Query {
	// 使用生成的Query
	return query.Use(db)
}

// NewDB exposes the primary Gorm DB from the unified DataSource for DI consumers.
func NewDB(ds *db.DataSource) *gorm.DB {
	if ds == nil || ds.Mysql == nil {
		panic("mysql data source is not initialized")
	}
	return ds.Mysql
}
