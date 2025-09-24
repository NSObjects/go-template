/*
 * Generated from OpenAPI3 document
 * Module: User
 */

package biz

import (
	"context"

	"github.com/NSObjects/go-template/internal/api/data"
	"github.com/NSObjects/go-template/internal/api/data/model"
	"github.com/NSObjects/go-template/internal/api/service/param"
	"github.com/NSObjects/go-template/internal/code"
	"gorm.io/gorm"
)

// UserUseCase 业务逻辑接口
type UserUseCase interface {

	// ListUsers 获取用户列表
	ListUsers(ctx context.Context, req param.UserListUsersRequest) ([]param.UserListItem, int64, error)

	// Create 创建用户
	Create(ctx context.Context, req param.UserCreateRequest) (*param.UserData, error)

	// GetByID 获取用户详情
	GetByID(ctx context.Context, id int64) (*param.UserData, error)

	// Update 更新用户
	Update(ctx context.Context, id int64, req param.UserUpdateRequest) (*param.UserData, error)

	// Delete 删除用户
	Delete(ctx context.Context, id int64) error
}

// UserHandler 业务逻辑处理器
type UserHandler struct {
	dataManager *data.DataManager
}

// NewUserHandler 创建业务逻辑处理器
func NewUserHandler(dataManager *data.DataManager) UserUseCase {
	return &UserHandler{
		dataManager: dataManager,
	}
}

// TODO: 实现业务逻辑方法

func (h *UserHandler) ListUsers(ctx context.Context, req param.UserListUsersRequest) ([]param.UserListItem, int64, error) {
	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}

	db := h.dataManager.MySQLWithContext(ctx)

	// 查询总数
	var total int64
	if err := db.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	var users []model.User
	offset := (req.Page - 1) * req.Size
	if err := db.Offset(offset).Limit(req.Size).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	// 转换为响应格式
	var responses []param.UserListItem
	for _, user := range users {
		responses = append(responses, param.UserListItem{
			Id:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Age:       user.Age,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	return responses, total, nil
}

func (h *UserHandler) Create(ctx context.Context, req param.UserCreateRequest) (*param.UserData, error) {
	db := h.dataManager.MySQLWithContext(ctx)

	// 检查用户名是否已存在
	var existingUser model.User
	if err := db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return nil, code.NewError(code.ErrUserAlreadyExists, "用户名已存在")
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// 检查邮箱是否已存在
	if err := db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, code.NewError(code.ErrUserAlreadyExists, "邮箱已存在")
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// 创建新用户
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Age:      req.Age,
	}

	if err := db.Create(user).Error; err != nil {
		return nil, err
	}

	// 转换为响应格式
	return &param.UserData{
		Id:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Age:       user.Age,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (h *UserHandler) GetByID(ctx context.Context, id int64) (*param.UserData, error) {
	db := h.dataManager.MySQLWithContext(ctx)

	var user model.User
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, code.NewError(code.ErrUserNotFound, "用户不存在")
		}
		return nil, err
	}

	// 转换为响应格式
	return &param.UserData{
		Id:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Age:       user.Age,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (h *UserHandler) Update(ctx context.Context, id int64, req param.UserUpdateRequest) (*param.UserData, error) {
	db := h.dataManager.MySQLWithContext(ctx)

	// 检查用户是否存在
	var user model.User
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, code.NewError(code.ErrUserNotFound, "用户不存在")
		}
		return nil, err
	}

	// 检查用户名是否与其他用户冲突
	if req.Username != "" && req.Username != user.Username {
		var existingUser model.User
		if err := db.Where("username = ? AND id != ?", req.Username, id).First(&existingUser).Error; err == nil {
			return nil, code.NewError(code.ErrUserAlreadyExists, "用户名已存在")
		} else if err != gorm.ErrRecordNotFound {
			return nil, err
		}
	}

	// 检查邮箱是否与其他用户冲突
	if req.Email != "" && req.Email != user.Email {
		var existingUser model.User
		if err := db.Where("email = ? AND id != ?", req.Email, id).First(&existingUser).Error; err == nil {
			return nil, code.NewError(code.ErrUserAlreadyExists, "邮箱已存在")
		} else if err != gorm.ErrRecordNotFound {
			return nil, err
		}
	}

	// 更新用户信息
	updates := make(map[string]interface{})
	if req.Username != "" {
		updates["username"] = req.Username
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Age > 0 {
		updates["age"] = req.Age
	}

	if err := db.Model(&user).Updates(updates).Error; err != nil {
		return nil, err
	}

	// 重新查询用户信息
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	// 转换为响应格式
	return &param.UserData{
		Id:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Age:       user.Age,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (h *UserHandler) Delete(ctx context.Context, id int64) error {
	db := h.dataManager.MySQLWithContext(ctx)

	// 检查用户是否存在
	var user model.User
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return code.NewError(code.ErrUserNotFound, "用户不存在")
		}
		return err
	}

	// 删除用户
	if err := db.Delete(&user).Error; err != nil {
		return err
	}

	return nil
}
