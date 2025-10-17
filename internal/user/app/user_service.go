package app

import (
	"context"
	"time"

	"github.com/NSObjects/go-template/internal/shared/kernel"
	"github.com/NSObjects/go-template/internal/user/domain"
)

// UserService 用户应用服务接口（入站端口）
type UserService interface {
	// CreateUser 创建用户
	CreateUser(ctx context.Context, req CreateUserRequest) (*UserDTO, error)

	// GetUser 获取用户详情
	GetUser(ctx context.Context, userID string) (*UserDTO, error)

	// UpdateUser 更新用户
	UpdateUser(ctx context.Context, userID string, req UpdateUserRequest) error

	// ListUsers 获取用户列表
	ListUsers(ctx context.Context, req ListUsersRequest) (*ListUsersResponse, error)

	// DeleteUser 删除用户
	DeleteUser(ctx context.Context, userID string) error

	// SuspendUser 暂停用户
	SuspendUser(ctx context.Context, userID string) error

	// ActivateUser 激活用户
	ActivateUser(ctx context.Context, userID string) error
}

// userServiceImpl 用户应用服务实现
type userServiceImpl struct {
	userRepo  domain.UserRepository
	eventBus  EventBus
	txManager TransactionManager
}

// NewUserService 创建用户应用服务
func NewUserService(
	userRepo domain.UserRepository,
	eventBus EventBus,
	txManager TransactionManager,
) UserService {
	return &userServiceImpl{
		userRepo:  userRepo,
		eventBus:  eventBus,
		txManager: txManager,
	}
}

// CreateUser 创建用户
func (s *userServiceImpl) CreateUser(ctx context.Context, req CreateUserRequest) (*UserDTO, error) {
	// 验证邮箱是否已存在
	email, err := domain.NewEmail(req.Email)
	if err != nil {
		return nil, kernel.NewBusinessRuleError("INVALID_EMAIL", err.Error())
	}

	existingUser, err := s.userRepo.FindByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, kernel.NewBusinessRuleError("EMAIL_EXISTS", "email already exists")
	}

	// 验证用户名是否已存在
	existingUser, err = s.userRepo.FindByUsername(ctx, req.Username)
	if err == nil && existingUser != nil {
		return nil, kernel.NewBusinessRuleError("USERNAME_EXISTS", "username already exists")
	}

	// 创建用户ID
	userID, err := domain.NewUserID(req.UserID)
	if err != nil {
		return nil, kernel.NewBusinessRuleError("INVALID_USER_ID", err.Error())
	}

	// 解析生日
	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		return nil, kernel.NewBusinessRuleError("INVALID_BIRTH_DATE", "invalid birth date format")
	}

	// 创建用户聚合根
	user, err := domain.NewUser(userID, req.Username, req.Email, birthDate)
	if err != nil {
		return nil, kernel.NewBusinessRuleError("CREATE_USER_FAILED", err.Error())
	}

	// 在事务中保存用户
	err = s.txManager.ExecuteInTransaction(ctx, func(ctx context.Context) error {
		if err := s.userRepo.Save(ctx, user); err != nil {
			return err
		}

		// 发布领域事件
		for _, event := range user.GetUncommittedEvents() {
			if err := s.eventBus.Publish(ctx, event); err != nil {
				return err
			}
		}
		user.MarkEventsAsCommitted()

		return nil
	})
	if err != nil {
		return nil, err
	}

	// 转换为DTO返回
	return UserAssembler{}.ToDTO(user), nil
}

// GetUser 获取用户详情
func (s *userServiceImpl) GetUser(ctx context.Context, userID string) (*UserDTO, error) {
	id, err := domain.NewUserID(userID)
	if err != nil {
		return nil, kernel.NewBusinessRuleError("INVALID_USER_ID", err.Error())
	}

	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return UserAssembler{}.ToDTO(user), nil
}

// UpdateUser 更新用户
func (s *userServiceImpl) UpdateUser(ctx context.Context, userID string, req UpdateUserRequest) error {
	id, err := domain.NewUserID(userID)
	if err != nil {
		return kernel.NewBusinessRuleError("INVALID_USER_ID", err.Error())
	}

	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// 在事务中更新用户
	return s.txManager.ExecuteInTransaction(ctx, func(ctx context.Context) error {
		// 更新邮箱（如果提供）
		if req.Email != "" && req.Email != user.Email().Value() {
			if err := user.ChangeEmail(req.Email); err != nil {
				return kernel.NewBusinessRuleError("CHANGE_EMAIL_FAILED", err.Error())
			}
		}

		// 更新用户名（如果提供）
		if req.Username != "" && req.Username != user.Username() {
			if err := user.ChangeUsername(req.Username); err != nil {
				return kernel.NewBusinessRuleError("CHANGE_USERNAME_FAILED", err.Error())
			}
		}

		// 保存用户
		if err := s.userRepo.Save(ctx, user); err != nil {
			return err
		}

		// 发布领域事件
		for _, event := range user.GetUncommittedEvents() {
			if err := s.eventBus.Publish(ctx, event); err != nil {
				return err
			}
		}
		user.MarkEventsAsCommitted()

		return nil
	})
}

// ListUsers 获取用户列表
func (s *userServiceImpl) ListUsers(ctx context.Context, req ListUsersRequest) (*ListUsersResponse, error) {
	users, total, err := s.userRepo.List(ctx, req.Offset, req.Size)
	if err != nil {
		return nil, err
	}

	// 转换为DTO列表
	userDTOs := make([]*UserDTO, len(users))
	for i, user := range users {
		userDTOs[i] = UserAssembler{}.ToDTO(user)
	}

	return &ListUsersResponse{
		Users: userDTOs,
		Total: total,
		Page:  req.Page,
		Size:  req.Size,
	}, nil
}

// DeleteUser 删除用户
func (s *userServiceImpl) DeleteUser(ctx context.Context, userID string) error {
	id, err := domain.NewUserID(userID)
	if err != nil {
		return kernel.NewBusinessRuleError("INVALID_USER_ID", err.Error())
	}

	return s.txManager.ExecuteInTransaction(ctx, func(ctx context.Context) error {
		return s.userRepo.Delete(ctx, id)
	})
}

// SuspendUser 暂停用户
func (s *userServiceImpl) SuspendUser(ctx context.Context, userID string) error {
	id, err := domain.NewUserID(userID)
	if err != nil {
		return kernel.NewBusinessRuleError("INVALID_USER_ID", err.Error())
	}

	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	return s.txManager.ExecuteInTransaction(ctx, func(ctx context.Context) error {
		if err := user.Suspend(); err != nil {
			return kernel.NewBusinessRuleError("SUSPEND_USER_FAILED", err.Error())
		}

		if err := s.userRepo.Save(ctx, user); err != nil {
			return err
		}

		// 发布领域事件
		for _, event := range user.GetUncommittedEvents() {
			if err := s.eventBus.Publish(ctx, event); err != nil {
				return err
			}
		}
		user.MarkEventsAsCommitted()

		return nil
	})
}

// ActivateUser 激活用户
func (s *userServiceImpl) ActivateUser(ctx context.Context, userID string) error {
	id, err := domain.NewUserID(userID)
	if err != nil {
		return kernel.NewBusinessRuleError("INVALID_USER_ID", err.Error())
	}

	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	return s.txManager.ExecuteInTransaction(ctx, func(ctx context.Context) error {
		if err := user.Activate(); err != nil {
			return kernel.NewBusinessRuleError("ACTIVATE_USER_FAILED", err.Error())
		}

		return s.userRepo.Save(ctx, user)
	})
}
