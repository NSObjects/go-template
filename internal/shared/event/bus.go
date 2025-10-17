package event

import "context"

// EventHandler 事件处理器接口
type EventHandler interface {
	Handle(ctx context.Context, event Event) error
	EventType() string
}

// EventBus 事件总线接口
type EventBus interface {
	// Publish 发布事件
	Publish(ctx context.Context, event Event) error

	// Subscribe 订阅事件
	Subscribe(handler EventHandler)

	// Start 启动事件总线
	Start(ctx context.Context) error

	// Stop 停止事件总线
	Stop() error
}

// EventStore 事件存储接口
type EventStore interface {
	// Save 保存事件
	Save(ctx context.Context, events []Event) error

	// Load 加载聚合根的所有事件
	Load(ctx context.Context, aggregateID string) ([]Event, error)

	// LoadFromVersion 从指定版本开始加载事件
	LoadFromVersion(ctx context.Context, aggregateID string, fromVersion int) ([]Event, error)
}
