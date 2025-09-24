package model

import (
	"time"
)

// User 用户模型
type User struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"type:varchar(50);not null;uniqueIndex" json:"username"`
	Email     string    `gorm:"type:varchar(100);not null;uniqueIndex" json:"email"`
	Age       int       `gorm:"type:int" json:"age"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}
