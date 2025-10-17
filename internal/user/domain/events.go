package domain

import (
	"time"

	"github.com/NSObjects/go-template/internal/shared/event"
)

// UserCreated 用户创建事件
type UserCreated struct {
	event.BaseEvent
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// EventID 实现DomainEvent接口
func (e *UserCreated) EventID() string {
	return e.BaseEvent.EventID()
}

// EventType 实现DomainEvent接口
func (e *UserCreated) EventType() string {
	return e.BaseEvent.EventType()
}

// OccurredAt 实现DomainEvent接口
func (e *UserCreated) OccurredAt() time.Time {
	return e.BaseEvent.OccurredAt()
}

// AggregateID 实现DomainEvent接口
func (e *UserCreated) AggregateID() string {
	return e.BaseEvent.AggregateID()
}

// Version 实现DomainEvent接口
func (e *UserCreated) Version() int {
	return e.BaseEvent.Version()
}

// UserEmailChanged 用户邮箱变更事件
type UserEmailChanged struct {
	event.BaseEvent
	UserID   string `json:"user_id"`
	OldEmail string `json:"old_email"`
	NewEmail string `json:"new_email"`
}

// EventID 实现DomainEvent接口
func (e *UserEmailChanged) EventID() string {
	return e.BaseEvent.EventID()
}

// EventType 实现DomainEvent接口
func (e *UserEmailChanged) EventType() string {
	return e.BaseEvent.EventType()
}

// OccurredAt 实现DomainEvent接口
func (e *UserEmailChanged) OccurredAt() time.Time {
	return e.BaseEvent.OccurredAt()
}

// AggregateID 实现DomainEvent接口
func (e *UserEmailChanged) AggregateID() string {
	return e.BaseEvent.AggregateID()
}

// Version 实现DomainEvent接口
func (e *UserEmailChanged) Version() int {
	return e.BaseEvent.Version()
}

// UserSuspended 用户暂停事件
type UserSuspended struct {
	event.BaseEvent
	UserID   string `json:"user_id"`
	Username string `json:"username"`
}

// EventID 实现DomainEvent接口
func (e *UserSuspended) EventID() string {
	return e.BaseEvent.EventID()
}

// EventType 实现DomainEvent接口
func (e *UserSuspended) EventType() string {
	return e.BaseEvent.EventType()
}

// OccurredAt 实现DomainEvent接口
func (e *UserSuspended) OccurredAt() time.Time {
	return e.BaseEvent.OccurredAt()
}

// AggregateID 实现DomainEvent接口
func (e *UserSuspended) AggregateID() string {
	return e.BaseEvent.AggregateID()
}

// Version 实现DomainEvent接口
func (e *UserSuspended) Version() int {
	return e.BaseEvent.Version()
}
