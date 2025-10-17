package app

import (
	"errors"
	"time"

	"github.com/NSObjects/go-template/internal/user/domain"
)

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	UserID    string `json:"user_id" validate:"required"`
	Username  string `json:"username" validate:"required,min=3,max=20"`
	Email     string `json:"email" validate:"required,email"`
	BirthDate string `json:"birth_date" validate:"required"` // YYYY-MM-DD格式
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Username string `json:"username,omitempty" validate:"omitempty,min=3,max=20"`
	Email    string `json:"email,omitempty" validate:"omitempty,email"`
}

// ListUsersRequest 用户列表请求
type ListUsersRequest struct {
	Page   int `json:"page" validate:"min=1"`
	Size   int `json:"size" validate:"min=1,max=100"`
	Offset int // 内部计算字段
}

// CalculateOffset 计算偏移量
func (r *ListUsersRequest) CalculateOffset() {
	if r.Page <= 0 {
		r.Page = 1
	}
	if r.Size <= 0 {
		r.Size = 10
	}
	r.Offset = (r.Page - 1) * r.Size
}

// UserDTO 用户数据传输对象
type UserDTO struct {
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Age       int       `json:"age"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ListUsersResponse 用户列表响应
type ListUsersResponse struct {
	Users []*UserDTO `json:"users"`
	Total int64      `json:"total"`
	Page  int        `json:"page"`
	Size  int        `json:"size"`
}

// UserAssembler 用户DTO转换器
type UserAssembler struct{}

// ToDTO 将领域实体转换为DTO
func (a UserAssembler) ToDTO(user *domain.User) *UserDTO {
	return &UserDTO{
		UserID:    user.ID().Value(),
		Username:  user.Username(),
		Email:     user.Email().Value(),
		Age:       user.Age(),
		Status:    a.statusToString(user.Status()),
		CreatedAt: user.CreatedAt(),
		UpdatedAt: user.UpdatedAt(),
	}
}

// ToEntity 将DTO转换为领域实体（用于更新场景）
func (a UserAssembler) ToEntity(dto *UserDTO) (*domain.User, error) {
	// 注意：这里需要从生日计算年龄，但DTO中没有生日字段
	// 实际场景中可能需要额外的字段或从数据库重新加载
	// 通常更新场景下会先加载现有实体，然后应用变更
	return nil, errors.New("not implemented")
}

// statusToString 将状态枚举转换为字符串
func (a UserAssembler) statusToString(status domain.UserStatus) string {
	switch status {
	case domain.UserStatusActive:
		return "active"
	case domain.UserStatusInactive:
		return "inactive"
	case domain.UserStatusSuspended:
		return "suspended"
	default:
		return "unknown"
	}
}

// stringToStatus 将字符串转换为状态枚举
func (a UserAssembler) stringToStatus(status string) domain.UserStatus {
	switch status {
	case "active":
		return domain.UserStatusActive
	case "inactive":
		return domain.UserStatusInactive
	case "suspended":
		return domain.UserStatusSuspended
	default:
		return domain.UserStatusInactive
	}
}
