package mysql

import (
	"context"
	"errors"
	"time"

	"github.com/NSObjects/go-template/internal/code"
	user "github.com/NSObjects/go-template/internal/modules/user"
	"github.com/NSObjects/go-template/internal/capabilities/userstorage/mysql/gormgen/model"
	"github.com/NSObjects/go-template/internal/capabilities/userstorage/mysql/gormgen/query"
)

type gormRepository struct {
	q *query.Query
}

func newRepository(q *query.Query) user.Repository {
	return gormRepository{q: q}
}

func (r gormRepository) query() (*query.Query, error) {
	if r.q == nil {
		return nil, errors.New("mysql is disabled")
	}
	return r.q, nil
}

func (r gormRepository) ListUsers(ctx context.Context, req user.ListUsersRequest) ([]user.ListItem, int64, error) {
	q, err := r.query()
	if err != nil {
		return nil, 0, code.WrapDatabaseError(err, "查询User列表失败")
	}

	users, err := q.User.WithContext(ctx).Offset(req.Offset()).Limit(req.Limit()).Find()
	if err != nil {
		return nil, 0, code.WrapDatabaseError(err, "查询User列表失败")
	}
	list := make([]user.ListItem, 0, len(users))
	for _, item := range users {
		list = append(list, user.ListItem{
			Id:        item.ID,
			Username:  item.Username,
			Email:     item.Email,
			Age:       int(item.Age),
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		})
	}

	count, err := q.User.Count()
	if err != nil {
		return nil, 0, code.WrapDatabaseError(err, "查询User列表失败")
	}
	return list, count, nil
}

func (r gormRepository) Create(ctx context.Context, req user.CreateRequest) error {
	q, err := r.query()
	if err != nil {
		return code.WrapDatabaseError(err, "创建User失败")
	}

	now := time.Now()
	item := model.User{
		Username:  req.Username,
		Email:     req.Email,
		Age:       int32(req.Age),
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := q.User.WithContext(ctx).Create(&item); err != nil {
		return code.WrapDatabaseError(err, "创建User失败")
	}
	return nil
}

func (r gormRepository) GetByID(ctx context.Context, id int64) (user.Data, error) {
	q, err := r.query()
	if err != nil {
		return user.Data{}, code.WrapDatabaseError(err, "查询User详情失败")
	}

	item, err := q.User.WithContext(ctx).GetByID(uint(id))
	if err != nil {
		return user.Data{}, code.WrapDatabaseError(err, "查询User详情失败")
	}
	return user.Data{
		Username:  item.Username,
		Email:     item.Email,
		Age:       int(item.Age),
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
		Id:        item.ID,
	}, nil
}

func (r gormRepository) Update(ctx context.Context, id int64, req user.UpdateRequest) error {
	q, err := r.query()
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

func (r gormRepository) Delete(ctx context.Context, id int64) error {
	q, err := r.query()
	if err != nil {
		return code.WrapDatabaseError(err, "删除User失败")
	}

	if err := q.User.WithContext(ctx).DeleteByID(uint(id)); err != nil {
		return code.WrapDatabaseError(err, "删除User失败")
	}
	return nil
}
