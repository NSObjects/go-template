package data

import (
	"context"
	"errors"
	"time"

	"github.com/NSObjects/go-template/internal/api/biz"
	"github.com/NSObjects/go-template/internal/api/data/model"
	"github.com/NSObjects/go-template/internal/api/data/query"
	"github.com/NSObjects/go-template/internal/api/service/param"
	"github.com/NSObjects/go-template/internal/code"
)

type userRepository struct {
	q *query.Query
}

func NewUserRepository(q *query.Query) biz.UserRepository {
	return userRepository{q: q}
}

func (u userRepository) query() (*query.Query, error) {
	if u.q == nil {
		return nil, errors.New("mysql is disabled")
	}
	return u.q, nil
}

func (u userRepository) ListUsers(ctx context.Context, req param.UserListUsersRequest) ([]param.UserListItem, int64, error) {
	q, err := u.query()
	if err != nil {
		return nil, 0, code.WrapDatabaseError(err, "查询User列表失败")
	}

	users, err := q.User.WithContext(ctx).Offset(req.Offset()).Limit(req.Limit()).Find()
	if err != nil {
		return nil, 0, code.WrapDatabaseError(err, "查询User列表失败")
	}
	var list []param.UserListItem
	for _, user := range users {
		list = append(list, param.UserListItem{
			Id:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Age:       int(user.Age),
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	count, err := q.User.Count()
	if err != nil {
		return nil, 0, code.WrapDatabaseError(err, "查询User列表失败")
	}

	return list, count, nil
}

func (u userRepository) Create(ctx context.Context, req param.UserCreateRequest) error {
	q, err := u.query()
	if err != nil {
		return code.WrapDatabaseError(err, "创建User失败")
	}

	user := model.User{
		Username:  req.Username,
		Email:     req.Email,
		Age:       int32(req.Age),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = q.User.WithContext(ctx).Create(&user)
	if err != nil {
		return code.WrapDatabaseError(err, "创建User失败")
	}

	return nil
}

func (u userRepository) GetByID(ctx context.Context, id int64) (param.UserData, error) {
	q, err := u.query()
	if err != nil {
		return param.UserData{}, code.WrapDatabaseError(err, "查询User详情失败")
	}

	user, err := q.User.WithContext(ctx).GetByID(uint(id))
	if err != nil {
		return param.UserData{}, code.WrapDatabaseError(err, "查询User详情失败")
	}

	return param.UserData{
		Username:  user.Username,
		Email:     user.Email,
		Age:       int(user.Age),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Id:        user.ID,
	}, nil
}

func (u userRepository) Update(ctx context.Context, id int64, req param.UserUpdateRequest) error {
	q, err := u.query()
	if err != nil {
		return code.WrapDatabaseError(err, "更新User失败")
	}

	_, err = q.User.WithContext(ctx).Where(q.User.ID.Eq(id)).Updates(model.User{
		Username: req.Username,
		Email:    req.Email,
		Age:      int32(req.Age),
	})
	if err != nil {
		return code.WrapDatabaseError(err, "更新User失败")
	}

	return nil
}

func (u userRepository) Delete(ctx context.Context, id int64) error {
	q, err := u.query()
	if err != nil {
		return code.WrapDatabaseError(err, "删除User失败")
	}

	err = q.User.WithContext(ctx).DeleteByID(uint(id))
	if err != nil {
		return code.WrapDatabaseError(err, "删除User失败")
	}

	return nil
}
