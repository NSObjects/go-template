package data

import (
	"context"
	"time"

	"github.com/NSObjects/go-template/internal/api/biz"
	"github.com/NSObjects/go-template/internal/api/data/db"
	"github.com/NSObjects/go-template/internal/api/data/model"
	"github.com/NSObjects/go-template/internal/api/service/param"
	"github.com/NSObjects/go-template/internal/code"
)

type userRepository struct {
	d *db.DataManager
}

func NewUserRepository(d *db.DataManager) biz.UserRepository {
	return userRepository{d: d}
}

func (u userRepository) ListUsers(ctx context.Context, req param.UserListUsersRequest) ([]param.UserListItem, int64, error) {

	users, err := u.d.Query.User.WithContext(ctx).Offset(req.Offset()).Limit(req.Limit()).Find()
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

	count, err := u.d.Query.User.Count()
	if err != nil {
		return nil, 0, err
	}

	return list, count, nil
}

func (u userRepository) Create(ctx context.Context, req param.UserCreateRequest) error {
	user := model.User{
		Username:  req.Username,
		Email:     req.Email,
		Age:       int32(req.Age),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := u.d.Query.User.WithContext(ctx).Create(&user)
	if err != nil {
		return err
	}

	return nil
}

func (u userRepository) GetByID(ctx context.Context, id int64) (param.UserData, error) {
	user, err := u.d.Query.User.WithContext(ctx).GetByID(uint(id))
	if err != nil {
		return param.UserData{}, err
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
	_, err := u.d.Query.User.WithContext(ctx).Where(u.d.Query.User.ID.Eq(id)).Updates(model.User{
		Username: req.Username,
		Email:    req.Email,
		Age:      int32(req.Age),
	})
	if err != nil {
		return err
	}

	return nil
}

func (u userRepository) Delete(ctx context.Context, id int64) error {
	err := u.d.Query.User.WithContext(ctx).DeleteByID(uint(id))
	if err != nil {
		return err
	}

	return nil
}
