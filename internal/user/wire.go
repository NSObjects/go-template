package user

import (
	"context"
	"fmt"
	"time"

	"github.com/NSObjects/go-template/internal/shared/event"
	"github.com/NSObjects/go-template/internal/user/adapters"
	"github.com/NSObjects/go-template/internal/user/app"
	"github.com/google/uuid"
	"go.uber.org/fx"
)

// UserModule User上下文模块
var UserModule = fx.Options(
	// 应用服务
	fx.Provide(app.NewUserService),

	// 适配器
	fx.Provide(adapters.NewUserRepository),
	fx.Provide(adapters.NewGormTxManager),
	fx.Provide(AsRoute(adapters.NewUserController)),

	// 事件总线（简单实现）
	fx.Provide(NewSimpleEventBus),

	// 其他出站端口
	fx.Provide(NewSimpleCachePort),
	fx.Provide(NewSimpleIdGenPort),
	fx.Provide(NewSimpleClockPort),
)

// AsRoute 将控制器标记为路由
func AsRoute(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(adapters.RegisterRouter)),
		fx.ResultTags(`group:"routes"`),
	)
}

// SimpleEventBus 简单事件总线实现
type SimpleEventBus struct{}

// NewSimpleEventBus 创建简单事件总线
func NewSimpleEventBus() app.EventBus {
	return &SimpleEventBus{}
}

// Publish 发布事件
func (b *SimpleEventBus) Publish(ctx context.Context, event event.Event) error {
	// 简单实现：直接打印日志
	// 实际项目中应该集成Kafka、RabbitMQ等消息队列
	fmt.Printf("Publishing event: %s, ID: %s\n", event.EventType(), event.EventID())
	return nil
}

// SimpleCachePort 简单缓存端口实现
type SimpleCachePort struct{}

// NewSimpleCachePort 创建简单缓存端口
func NewSimpleCachePort() app.CachePort {
	return &SimpleCachePort{}
}

// Get 获取缓存
func (c *SimpleCachePort) Get(ctx context.Context, key string) (string, error) {
	// 简单实现：返回空
	return "", nil
}

// Set 设置缓存
func (c *SimpleCachePort) Set(ctx context.Context, key, value string, ttl int) error {
	// 简单实现：直接返回
	return nil
}

// Delete 删除缓存
func (c *SimpleCachePort) Delete(ctx context.Context, key string) error {
	// 简单实现：直接返回
	return nil
}

// SimpleIdGenPort 简单ID生成器端口实现
type SimpleIdGenPort struct{}

// NewSimpleIdGenPort 创建简单ID生成器端口
func NewSimpleIdGenPort() app.IdGenPort {
	return &SimpleIdGenPort{}
}

// GenerateID 生成唯一ID
func (g *SimpleIdGenPort) GenerateID() (string, error) {
	// 简单实现：使用UUID
	return uuid.New().String(), nil
}

// SimpleClockPort 简单时钟端口实现
type SimpleClockPort struct{}

// NewSimpleClockPort 创建简单时钟端口
func NewSimpleClockPort() app.ClockPort {
	return &SimpleClockPort{}
}

// Now 获取当前时间
func (c *SimpleClockPort) Now() time.Time {
	return time.Now()
}
