package domain

import (
	"context"
)

// UserRepository 用户仓储接口（出站端口）
// 定义在领域层，由适配器层实现
type UserRepository interface {
	// Save 保存用户
	Save(ctx context.Context, user *User) error

	// FindByID 根据ID查找用户
	FindByID(ctx context.Context, id UserID) (*User, error)

	// FindByEmail 根据邮箱查找用户
	FindByEmail(ctx context.Context, email Email) (*User, error)

	// FindByUsername 根据用户名查找用户
	FindByUsername(ctx context.Context, username string) (*User, error)

	// List 分页查询用户列表
	List(ctx context.Context, offset, limit int) ([]*User, int64, error)

	// Delete 删除用户
	Delete(ctx context.Context, id UserID) error

	// Exists 检查用户是否存在
	Exists(ctx context.Context, id UserID) (bool, error)
}

// UserNotFoundError 用户未找到错误
type UserNotFoundError struct {
	UserID string
}

func (e *UserNotFoundError) Error() string {
	return "user not found: " + e.UserID
}

// UserAlreadyExistsError 用户已存在错误
type UserAlreadyExistsError struct {
	Email string
}

func (e *UserAlreadyExistsError) Error() string {
	return "user already exists with email: " + e.Email
}
