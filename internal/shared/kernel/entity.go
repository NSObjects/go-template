package kernel

import (
	"time"
)

// EntityID 实体ID接口
type EntityID interface {
	Value() string
	Equals(other EntityID) bool
}

// BaseEntity 基础实体，包含通用字段
type BaseEntity struct {
	id        EntityID
	createdAt time.Time
	updatedAt time.Time
}

// NewBaseEntity 创建基础实体
func NewBaseEntity(id EntityID) BaseEntity {
	now := time.Now()
	return BaseEntity{
		id:        id,
		createdAt: now,
		updatedAt: now,
	}
}

// ID 获取实体ID
func (e *BaseEntity) ID() EntityID {
	return e.id
}

// CreatedAt 获取创建时间
func (e *BaseEntity) CreatedAt() time.Time {
	return e.createdAt
}

// UpdatedAt 获取更新时间
func (e *BaseEntity) UpdatedAt() time.Time {
	return e.updatedAt
}

// SetUpdatedAt 设置更新时间
func (e *BaseEntity) SetUpdatedAt(t time.Time) {
	e.updatedAt = t
}

// MarkUpdated 标记为已更新
func (e *BaseEntity) MarkUpdated() {
	e.updatedAt = time.Now()
}
