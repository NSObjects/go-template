/*
 * Created by lintao on 2023/7/18 下午3:59
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package biz

import (
	"github.com/NSObjects/go-template/internal/api/data"

	"github.com/NSObjects/go-template/internal/api/service/param"

	"github.com/NSObjects/go-template/internal/api/data/model"
)

type UserHandler struct {
	repository data.UserRepository
}

func NewUserHandler(repository data.UserRepository) *UserHandler {
	return &UserHandler{repository: repository}
}

func (h *UserHandler) ListUser(u model.User, p param.APIQuery) ([]param.UserResponse, int64, error) {

	users, total, err := h.repository.FindUser(u, p)
	if err != nil {
		return nil, 0, err
	}

	resp := make([]param.UserResponse, len(users))
	for i, user := range users {
		resp[i] = param.UserResponse{
			Name:     user.Name,
			Phone:    user.Phone,
			Status:   user.Status,
			Password: user.Password,
		}
	}

	return resp, total, nil
}

func (h *UserHandler) CreateUser(param model.User) (err error) {
	if _, err = h.repository.CreateUser(param); err != nil {
		return err
	}
	return nil
}

func (h *UserHandler) DeleteUser(id int64) (err error) {
	if err = h.repository.DeleteUserByID(id); err != nil {
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

func (h *UserHandler) GetUserDetail(id int64) (param.UserResponse, error) {
	user, err := h.repository.GetUserByID(id)
	if err != nil {
		return param.UserResponse{}, err
	}

	return param.UserResponse{
		Name:     user.Name,
		Phone:    user.Phone,
		Status:   user.Status,
		Password: user.Password,
	}, nil
}
