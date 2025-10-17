package adapters

import (
	"context"

	"github.com/NSObjects/go-template/internal/infra/persistence"
	"github.com/NSObjects/go-template/internal/pkg/code"
	"github.com/NSObjects/go-template/internal/user/domain"
	"gorm.io/gorm"
)

// userRepositoryImpl 用户仓储实现
type userRepositoryImpl struct {
	db     *persistence.DataManager
	mapper UserMapper
}

// NewUserRepository 创建用户仓储实现
func NewUserRepository(db *persistence.DataManager) domain.UserRepository {
	return &userRepositoryImpl{
		db:     db,
		mapper: UserMapper{},
	}
}

// Save 保存用户
func (r *userRepositoryImpl) Save(ctx context.Context, user *domain.User) error {
	po := r.mapper.ToPO(user)

	// 检查用户是否已存在
	userID, err := domain.NewUserID(user.ID().Value())
	if err != nil {
		return code.WrapDatabaseError(err, "invalid user ID")
	}

	exists, err := r.Exists(ctx, userID)
	if err != nil {
		return code.WrapDatabaseError(err, "check user exists failed")
	}

	if exists {
		// 更新用户
		err = r.db.MySQLWithContext(ctx).Model(&UserPO{}).Where("id = ?", po.ID).Updates(po).Error
	} else {
		// 创建用户
		err = r.db.MySQLWithContext(ctx).Create(po).Error
	}

	if err != nil {
		return code.WrapDatabaseError(err, "save user failed")
	}

	return nil
}

// FindByID 根据ID查找用户
func (r *userRepositoryImpl) FindByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
	var po UserPO
	err := r.db.MySQLWithContext(ctx).Where("id = ?", id.Value()).First(&po).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &domain.UserNotFoundError{UserID: id.Value()}
		}
		return nil, code.WrapDatabaseError(err, "find user by id failed")
	}

	user, err := r.mapper.ToEntity(&po)
	if err != nil {
		return nil, code.WrapDatabaseError(err, "convert user entity failed")
	}

	return user, nil
}

// FindByEmail 根据邮箱查找用户
func (r *userRepositoryImpl) FindByEmail(ctx context.Context, email domain.Email) (*domain.User, error) {
	var po UserPO
	err := r.db.MySQLWithContext(ctx).Where("email = ?", email.Value()).First(&po).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &domain.UserNotFoundError{UserID: email.Value()}
		}
		return nil, code.WrapDatabaseError(err, "find user by email failed")
	}

	user, err := r.mapper.ToEntity(&po)
	if err != nil {
		return nil, code.WrapDatabaseError(err, "convert user entity failed")
	}

	return user, nil
}

// FindByUsername 根据用户名查找用户
func (r *userRepositoryImpl) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	var po UserPO
	err := r.db.MySQLWithContext(ctx).Where("username = ?", username).First(&po).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &domain.UserNotFoundError{UserID: username}
		}
		return nil, code.WrapDatabaseError(err, "find user by username failed")
	}

	user, err := r.mapper.ToEntity(&po)
	if err != nil {
		return nil, code.WrapDatabaseError(err, "convert user entity failed")
	}

	return user, nil
}

// List 分页查询用户列表
func (r *userRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*domain.User, int64, error) {
	var pos []UserPO
	var total int64

	// 查询总数
	err := r.db.MySQLWithContext(ctx).Model(&UserPO{}).Count(&total).Error
	if err != nil {
		return nil, 0, code.WrapDatabaseError(err, "count users failed")
	}

	// 查询数据
	err = r.db.MySQLWithContext(ctx).Offset(offset).Limit(limit).Find(&pos).Error
	if err != nil {
		return nil, 0, code.WrapDatabaseError(err, "list users failed")
	}

	// 转换为领域实体
	users := make([]*domain.User, len(pos))
	for i, po := range pos {
		user, err := r.mapper.ToEntity(&po)
		if err != nil {
			return nil, 0, code.WrapDatabaseError(err, "convert user entity failed")
		}
		users[i] = user
	}

	return users, total, nil
}

// Delete 删除用户
func (r *userRepositoryImpl) Delete(ctx context.Context, id domain.UserID) error {
	err := r.db.MySQLWithContext(ctx).Where("id = ?", id.Value()).Delete(&UserPO{}).Error
	if err != nil {
		return code.WrapDatabaseError(err, "delete user failed")
	}

	return nil
}

// Exists 检查用户是否存在
func (r *userRepositoryImpl) Exists(ctx context.Context, id domain.UserID) (bool, error) {
	var count int64
	err := r.db.MySQLWithContext(ctx).Model(&UserPO{}).Where("id = ?", id.Value()).Count(&count).Error
	if err != nil {
		return false, code.WrapDatabaseError(err, "check user exists failed")
	}

	return count > 0, nil
}
