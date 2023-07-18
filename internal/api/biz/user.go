/*
 * Created by lintao on 2023/7/18 下午3:59
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package biz

import (
	"strconv"

	"github.com/NSObjects/go-template/internal/api/data"
	"github.com/NSObjects/go-template/internal/api/data/model"
	"github.com/NSObjects/go-template/internal/resp"
	"github.com/labstack/echo/v4"
)

type UserUsecase interface {
	GetUser(c echo.Context) error
	GetUserDetail(c echo.Context) error
	CreateUser(c echo.Context) error
	DeleteUser(c echo.Context) error
	UpdateUser(c echo.Context) error
}

type UserHandler struct {
	repository data.UserRepository
}

func NewUserHandler(repository data.UserRepository) *UserHandler {
	return &UserHandler{repository: repository}
}

func (h *UserHandler) GetUser(c echo.Context) (err error) {
	var p model.UserParam
	if err := c.Bind(&p); err != nil {
		return resp.ApiError(resp.NewParamError(err), c)
	}

	users, total, err := h.repository.FindUser(p)
	if err != nil {
		return resp.ApiError(resp.NewDBError(err), c)
	}

	return resp.ListDataResponse(users, total, c)
}

func (h *UserHandler) CreateUser(c echo.Context) (err error) {

	var p model.UserParam
	if err := c.Bind(&p); err != nil {
		return resp.ApiError(resp.NewParamError(err), c)
	}

	if _, err = h.repository.CreateUser(p); err != nil {
		return resp.ApiError(resp.NewDBError(err), c)
	}

	return resp.OperateSuccess(c)
}

func (h *UserHandler) DeleteUser(c echo.Context) (err error) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return resp.ApiError(resp.NewParamError(err), c)
	}

	err = h.repository.DeleteUserById(int64(id))
	if err != nil {
		return resp.ApiError(resp.NewDBError(err), c)
	}

	return resp.OperateSuccess(c)
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return resp.ApiError(resp.NewParamError(err), c)
	}

	param := new(model.UserParam)
	if err := c.Bind(param); err != nil {
		return resp.ApiError(resp.NewParamError(err), c)
	}

	err = h.repository.UpdateUser(*param, int64(id))
	if err != nil {
		return resp.ApiError(resp.NewDBError(err), c)
	}

	return resp.OperateSuccess(c)
}

func (h *UserHandler) GetUserDetail(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return resp.ApiError(resp.NewParamError(err), c)
	}

	user, err := h.repository.GetUserById(int64(id))
	if err != nil {
		return resp.ApiError(resp.NewDBError(err), c)
	}

	return resp.OneDataResponse(user, c)
}
