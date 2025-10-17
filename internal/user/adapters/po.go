package adapters

import (
	"errors"
	"time"

	"github.com/NSObjects/go-template/internal/user/domain"
)

// UserPO 用户持久化对象
type UserPO struct {
	ID        string    `gorm:"column:id;type:varchar(50);primaryKey" json:"id"`
	Username  string    `gorm:"column:username;type:varchar(50);index:username_idx,priority:1" json:"username"`
	Email     string    `gorm:"column:email;type:varchar(100);index:email_idx,priority:1" json:"email"`
	BirthDate time.Time `gorm:"column:birth_date;type:date" json:"birth_date"`
	Status    int32     `gorm:"column:status;type:tinyint;default:1" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;index:created_at_idx,priority:1;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName 指定表名
func (UserPO) TableName() string {
	return "users"
}

// UserMapper 用户PO与Entity转换器
type UserMapper struct{}

// ToPO 将领域实体转换为持久化对象
func (m UserMapper) ToPO(user *domain.User) *UserPO {
	return &UserPO{
		ID:        user.ID().Value(),
		Username:  user.Username(),
		Email:     user.Email().Value(),
		BirthDate: user.BirthDate().Value(),
		Status:    int32(m.statusToInt(user.Status())),
		CreatedAt: user.CreatedAt(),
		UpdatedAt: user.UpdatedAt(),
	}
}

// ToEntity 将持久化对象转换为领域实体
func (m UserMapper) ToEntity(po *UserPO) (*domain.User, error) {
	// 暂时返回错误，需要实现FromPersistence构造函数
	// 由于NewUser会发布事件，我们需要一个不发布事件的构造函数
	return nil, errors.New("not implemented - need FromPersistence constructor")
}

// statusToInt 将状态枚举转换为整数
func (m UserMapper) statusToInt(status domain.UserStatus) int {
	switch status {
	case domain.UserStatusActive:
		return 1
	case domain.UserStatusInactive:
		return 0
	case domain.UserStatusSuspended:
		return 2
	default:
		return 0
	}
}

// intToStatus 将整数转换为状态枚举
func (m UserMapper) intToStatus(status int32) domain.UserStatus {
	switch status {
	case 1:
		return domain.UserStatusActive
	case 0:
		return domain.UserStatusInactive
	case 2:
		return domain.UserStatusSuspended
	default:
		return domain.UserStatusInactive
	}
}
