package data

import (
	"context"
	"time"

	"github.com/NSObjects/go-template/internal/api/data/db"
	appuser "github.com/NSObjects/go-template/internal/application/user"
	"github.com/NSObjects/go-template/internal/code"
	domain "github.com/NSObjects/go-template/internal/domain/user"
)

type userRepository struct {
	d *db.DataManager
}

func NewUserRepository(d *db.DataManager) domain.Repository {
	return userRepository{d: d}
}

func (u userRepository) List(ctx context.Context, query domain.ListUsersQuery) ([]domain.User, int64, error) {
	users, err := u.d.Query.User.WithContext(ctx).Offset(query.Offset()).Limit(query.Limit()).Find()
	if err != nil {
		return nil, 0, code.WrapDatabaseError(err, "查询用户列表失败")
	}

	result := make([]domain.User, 0, len(users))
	for _, item := range users {
		aggregate, convErr := appuser.AssembleDomainUser(*item)
		if convErr != nil {
			return nil, 0, convErr
		}
		result = append(result, aggregate)
	}

	count, err := u.d.Query.User.Count()
	if err != nil {
		return nil, 0, code.WrapDatabaseError(err, "统计用户数量失败")
	}

	return result, count, nil
}

func (u userRepository) Create(ctx context.Context, entity domain.User) error {
	now := time.Now()
	record := appuser.AssembleModelUserForCreate(entity, now)
	if err := u.d.Query.User.WithContext(ctx).Create(&record); err != nil {
		return code.WrapDatabaseError(err, "创建用户失败")
	}
	return nil
}

func (u userRepository) GetByID(ctx context.Context, id domain.ID) (domain.User, error) {
	userModel, err := u.d.Query.User.WithContext(ctx).GetByID(uint(id.Int64()))
	if err != nil {
		return domain.User{}, code.WrapDatabaseError(err, "查询用户详情失败")
	}
	aggregate, convErr := appuser.AssembleDomainUser(*userModel)
	if convErr != nil {
		return domain.User{}, convErr
	}
	return aggregate, nil
}

func (u userRepository) Update(ctx context.Context, entity domain.User) error {
	now := time.Now()
	record := appuser.AssembleModelUserForUpdate(entity, now)
	_, err := u.d.Query.User.WithContext(ctx).Where(u.d.Query.User.ID.Eq(entity.ID.Int64())).Updates(record)
	if err != nil {
		return code.WrapDatabaseError(err, "更新用户失败")
	}
	return nil
}

func (u userRepository) Delete(ctx context.Context, id domain.ID) error {
	if err := u.d.Query.User.WithContext(ctx).DeleteByID(uint(id.Int64())); err != nil {
		return code.WrapDatabaseError(err, "删除用户失败")
	}
	return nil
}
