/*
 *
 * user.go
 * ucase
 *
 * Created by lintao on 2019-01-29 16:24
 * Copyright © 2017-2019 PYL. All rights reserved.
 *
 */

package ucase

import (
	"go-template/apis/api_helper"
	"go-template/models"
	"go-template/repository"
	"strconv"

	"github.com/labstack/echo"
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
	var p models.UserParam
	if err := c.Bind(&p); err != nil {
		return api_helper.ApiError(api_helper.NewParamError(err), c)
	}

	users, total, err := this.repository.FindUser(p)
	if err != nil {
		return api_helper.ApiError(api_helper.NewDBError(err), c)
	}

	return api_helper.ListDataResponse(users, total, c)
}

func (this *UserHandler) CreateUser(c echo.Context) (err error) {

	var p models.UserParam
	if err := c.Bind(&p); err != nil {
		return api_helper.ApiError(api_helper.NewParamError(err), c)
	}

	if _, err = this.repository.CreateUser(p); err != nil {
		return api_helper.ApiError(api_helper.NewDBError(err), c)
	}

	return api_helper.OperateSuccess(c)
}

func (this *UserHandler) DeleteUser(c echo.Context) (err error) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return api_helper.ApiError(api_helper.NewParamError(err), c)
	}

	err = this.repository.DeleteUserById(int64(id))
	if err != nil {
		return api_helper.ApiError(api_helper.NewDBError(err), c)
	}

	return api_helper.OperateSuccess(c)
}

func (this *UserHandler) UpdateUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return api_helper.ApiError(api_helper.NewParamError(err), c)
	}

	param := new(models.UserParam)
	if err := c.Bind(param); err != nil {
		return api_helper.ApiError(api_helper.NewParamError(err), c)
	}

	err = this.repository.UpdateUser(*param, int64(id))
	if err != nil {
		return api_helper.ApiError(api_helper.NewDBError(err), c)
	}

	return api_helper.OperateSuccess(c)
}

func (this *UserHandler) GetUserDetail(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return api_helper.ApiError(api_helper.NewParamError(err), c)
	}

	user, err := this.repository.GetUserById(int64(id))
	if err != nil {
		return api_helper.ApiError(api_helper.NewDBError(err), c)
	}

	return api_helper.OneDataResponse(user, c)
}
