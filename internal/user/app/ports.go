package app

import (
	"context"
	"time"

	"github.com/NSObjects/go-template/internal/shared/event"
)

// TransactionManager 事务管理器接口（出站端口）
type TransactionManager interface {
	// ExecuteInTransaction 在事务中执行操作
	ExecuteInTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// EventBus 事件总线接口（出站端口）
type EventBus interface {
	// Publish 发布事件
	Publish(ctx context.Context, event event.Event) error
}

// CachePort 缓存端口（出站端口）
type CachePort interface {
	// Get 获取缓存
	Get(ctx context.Context, key string) (string, error)

	// Set 设置缓存
	Set(ctx context.Context, key, value string, ttl int) error

	// Delete 删除缓存
	Delete(ctx context.Context, key string) error
}

// IdGenPort ID生成器端口（出站端口）
type IdGenPort interface {
	// GenerateID 生成唯一ID
	GenerateID() (string, error)
}

// ClockPort 时钟端口（出站端口，主要用于测试）
type ClockPort interface {
	// Now 获取当前时间
	Now() time.Time
}
