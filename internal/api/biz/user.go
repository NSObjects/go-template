/*
 * Generated from OpenAPI3 document
 * Module: User
 */

package biz

import (
	"context"

	"github.com/NSObjects/go-template/internal/api/service/param"
	"github.com/NSObjects/go-template/internal/code"
)

// UserRepository 数据访问接口 - 符合依赖注入原则
type UserRepository interface {

	// ListUsers 查询User列表
	ListUsers(ctx context.Context, req param.UserListUsersRequest) ([]param.UserListItem, int64, error)

	// Create 查询User详情
	Create(ctx context.Context, req param.UserCreateRequest) error

	// GetByID 查询User详情
	GetByID(ctx context.Context, id int64) (param.UserData, error)

	// Update 查询User详情
	Update(ctx context.Context, id int64, req param.UserUpdateRequest) error

	// Delete 执行User操作
	Delete(ctx context.Context, id int64) error
}

// UserUseCase 业务逻辑接口
type UserUseCase interface {
	// ListUsers 获取用户列表
	ListUsers(ctx context.Context, req param.UserListUsersRequest) ([]param.UserListItem, int64, error)

	// Create 创建用户
	Create(ctx context.Context, req param.UserCreateRequest) error

	// GetByID 获取用户详情
	GetByID(ctx context.Context, id int64) (param.UserData, error)

	// Update 更新用户
	Update(ctx context.Context, id int64, req param.UserUpdateRequest) error

	// Delete 删除用户
	Delete(ctx context.Context, id int64) error
}

// UserHandler 业务逻辑处理器 - 通过依赖注入获取Repository
type UserHandler struct {
	repo UserRepository
}

// NewUserHandler 创建业务逻辑处理器 - 符合依赖注入原则
func NewUserHandler(repo UserRepository) UserUseCase {
	return &UserHandler{
		repo: repo,
	}
}

func (h *UserHandler) ListUsers(ctx context.Context, req param.UserListUsersRequest) ([]param.UserListItem, int64, error) {
	list, total, err := h.repo.ListUsers(ctx, req)
	if err != nil {
		// 使用internal/code包包装错误 - 符合错误处理规范
		return nil, 0, code.WrapDatabaseError(err, "查询User列表失败")
	}
	return list, total, nil

}

func (h *UserHandler) Create(ctx context.Context, req param.UserCreateRequest) error {

	err := h.repo.Create(ctx, req)
	if err != nil {
		return code.WrapDatabaseError(err, "查询User详情失败")
	}
	return nil

}

func (h *UserHandler) GetByID(ctx context.Context, id int64) (param.UserData, error) {

	result, err := h.repo.GetByID(ctx, id)
	if err != nil {

		return result, code.WrapDatabaseError(err, "查询User详情失败")
	}
	return result, nil

}

func (h *UserHandler) Update(ctx context.Context, id int64, req param.UserUpdateRequest) error {

	err := h.repo.Update(ctx, id, req)
	if err != nil {
		return code.WrapDatabaseError(err, "查询User详情失败")
	}
	return nil

}

func (h *UserHandler) Delete(ctx context.Context, id int64) error {
	err := h.repo.Delete(ctx, id)
	if err != nil {
		return code.WrapDatabaseError(err, "删除User失败")
	}
	return nil

}
