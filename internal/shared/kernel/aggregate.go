package kernel

import "time"

// AggregateRoot 聚合根标记接口
// 聚合根是DDD中的核心概念，负责维护聚合内的一致性边界
type AggregateRoot interface {
	// GetUncommittedEvents 获取未提交的领域事件
	GetUncommittedEvents() []DomainEvent

	// MarkEventsAsCommitted 标记事件为已提交
	MarkEventsAsCommitted()

	// AddDomainEvent 添加领域事件
	AddDomainEvent(event DomainEvent)
}

// DomainEvent 领域事件接口
type DomainEvent interface {
	// EventID 事件唯一标识
	EventID() string

	// EventType 事件类型
	EventType() string

	// OccurredAt 事件发生时间
	OccurredAt() time.Time

	// AggregateID 聚合根ID
	AggregateID() string

	// Version 事件版本
	Version() int
}

// BaseAggregateRoot 聚合根基础实现
type BaseAggregateRoot struct {
	uncommittedEvents []DomainEvent
}

// GetUncommittedEvents 获取未提交的领域事件
func (a *BaseAggregateRoot) GetUncommittedEvents() []DomainEvent {
	return a.uncommittedEvents
}

// MarkEventsAsCommitted 标记事件为已提交
func (a *BaseAggregateRoot) MarkEventsAsCommitted() {
	a.uncommittedEvents = nil
}

// AddDomainEvent 添加领域事件
func (a *BaseAggregateRoot) AddDomainEvent(event DomainEvent) {
	a.uncommittedEvents = append(a.uncommittedEvents, event)
}
