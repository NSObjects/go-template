/*
 * Created by lintao on 2023/7/18 下午3:59
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package biz

import (
	"github.com/NSObjects/go-template/internal/api/data/repo"

	"github.com/NSObjects/go-template/internal/api/service/param"

	"github.com/NSObjects/go-template/internal/api/data/model"
)

//type UserUsecase interface {
//	ListUser(p param.APIQuery) error
//	GetUserDetail(c echo.Context) error
//	CreateUser(c echo.Context) error
//	DeleteUser(c echo.Context) error
//	UpdateUser(c echo.Context) error
//}

type UserHandler struct {
	repository repo.UserRepository
}

func NewUserHandler(repository repo.UserRepository) *UserHandler {
	return &UserHandler{repository: repository}
}

func (h *UserHandler) ListUser(param model.User, p param.APIQuery) (users []model.User, total int64, err error) {

	users, total, err = h.repository.FindUser(param, p)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (h *UserHandler) CreateUser(param model.User) (err error) {
	if _, err = h.repository.CreateUser(param); err != nil {
		return err
	}

	return nil
}

func (h *UserHandler) DeleteUser(id int64) (err error) {
	if err = h.repository.DeleteUserById(id); err != nil {
		return err
	}

	return err
}

func (h *UserHandler) UpdateUser(user model.User, id int64) error {
	if err := h.repository.UpdateUser(user, id); err != nil {
		return err
	}
	return nil
}

func (h *UserHandler) GetUserDetail(id int64) (model.User, error) {
	user, err := h.repository.GetUserById(id)
	if err != nil {
		return user, err
	}

	return user, nil
}
