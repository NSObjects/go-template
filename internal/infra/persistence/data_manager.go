package persistence

import (
	"context"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// DataManager 数据管理器
type DataManager struct {
	mysql *gorm.DB
	redis *redis.Client
}

// NewDataManager 创建数据管理器
func NewDataManager(mysql *gorm.DB, redis *redis.Client) *DataManager {
	return &DataManager{
		mysql: mysql,
		redis: redis,
	}
}

// MySQL 获取MySQL连接
func (dm *DataManager) MySQL() *gorm.DB {
	return dm.mysql
}

// MySQLWithContext 获取带上下文的MySQL连接
func (dm *DataManager) MySQLWithContext(ctx context.Context) *gorm.DB {
	return dm.mysql.WithContext(ctx)
}

// Redis 获取Redis客户端
func (dm *DataManager) Redis() *redis.Client {
	return dm.redis
}

// Model 数据管理器模块
var Model = fx.Options(
	fx.Provide(NewDataManager),
)
