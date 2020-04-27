/*
 *
 * user.go
 * ucase
 *
 * Created by lintao on 2019-01-29 16:24
 * Copyright © 2017-2019 PYL. All rights reserved.
 *
 */

package usecase

import (
	"go-template/domain"
	"go-template/repository"
	"go-template/tools"
	"strconv"

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
	repository repository.UserRepository
}

func NewUserHandler(repository repository.UserRepository) *UserHandler {
	return &UserHandler{repository: repository}
}

func (this *UserHandler) GetUser(c echo.Context) (err error) {
	var p domain.UserParam
	if err := c.Bind(&p); err != nil {
		return tools.ApiError(tools.NewParamError(err), c)
	}

	users, total, err := this.repository.FindUser(p)
	if err != nil {
		return tools.ApiError(tools.NewDBError(err), c)
	}

	return tools.ListDataResponse(users, total, c)
}

func (this *UserHandler) CreateUser(c echo.Context) (err error) {

	var p domain.UserParam
	if err := c.Bind(&p); err != nil {
		return tools.ApiError(tools.NewParamError(err), c)
	}

	if _, err = this.repository.CreateUser(p); err != nil {
		return tools.ApiError(tools.NewDBError(err), c)
	}

	return tools.OperateSuccess(c)
}

func (this *UserHandler) DeleteUser(c echo.Context) (err error) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return tools.ApiError(tools.NewParamError(err), c)
	}

	err = this.repository.DeleteUserById(int64(id))
	if err != nil {
		return tools.ApiError(tools.NewDBError(err), c)
	}

	return tools.OperateSuccess(c)
}

func (this *UserHandler) UpdateUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return tools.ApiError(tools.NewParamError(err), c)
	}

	param := new(domain.UserParam)
	if err := c.Bind(param); err != nil {
		return tools.ApiError(tools.NewParamError(err), c)
	}

	err = this.repository.UpdateUser(*param, int64(id))
	if err != nil {
		return tools.ApiError(tools.NewDBError(err), c)
	}

	return tools.OperateSuccess(c)
}

func (this *UserHandler) GetUserDetail(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return tools.ApiError(tools.NewParamError(err), c)
	}

	user, err := this.repository.GetUserById(int64(id))
	if err != nil {
		return tools.ApiError(tools.NewDBError(err), c)
	}

	return tools.OneDataResponse(user, c)
}
