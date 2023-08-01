/*
 * Created by lintao on 2023/7/27 下午1:44
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package service

import (
	"strconv"

	"github.com/NSObjects/go-template/internal/api/biz"
	"github.com/NSObjects/go-template/internal/api/data/model"
	"github.com/NSObjects/go-template/internal/api/service/param"
	"github.com/NSObjects/go-template/internal/resp"
	"github.com/labstack/echo/v4"
)

type Controller struct {
	user *biz.UserHandler
}

func (u *Controller) RegisterRouter(s *echo.Group, middlewareFunc ...echo.MiddlewareFunc) {
	s.GET("/users", u.getUser).Name = "用户查询"
	s.POST("/users", u.createUser).Name = "创建用户"
	s.DELETE("/users/:id", u.deleteUser).Name = "删除用户"
	s.PUT("/users/:id", u.updateUser).Name = "更新用户"
	s.GET("/users/:id", u.getUserDetail).Name = "获取某个用户信息"
}

func NewUserController(u *biz.UserHandler) RegisterRouter {
	return &Controller{
		user: u,
	}
}

func (u *Controller) getUser(c echo.Context) (err error) {
	var user param.UserParam
	if err = BindAndValidate(&user, c); err != nil {
		return err
	}

	listUser, total, err := u.user.ListUser(user.User, user.APIQuery)
	if err != nil {
		return err
	}
	return resp.ListDataResponse(listUser, total, c)
}

func (u *Controller) getUserDetail(c echo.Context) (err error) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	detail, err := u.user.GetUserDetail(id)
	if err != nil {
		return err
	}

	return resp.OneDataResponse(detail, c)
}

func (u *Controller) createUser(c echo.Context) (err error) {
	var user model.User
	if err = BindAndValidate(&user, c); err != nil {
		return err
	}

	if err = u.user.CreateUser(user); err != nil {
		return err
	}

	return resp.OperateSuccess(c)
}

func (u *Controller) updateUser(c echo.Context) (err error) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var user model.User
	if err = BindAndValidate(&user, c); err != nil {
		return err
	}

	if err = u.user.UpdateUser(user, id); err != nil {
		return err
	}

	return resp.OperateSuccess(c)

}

func (u *Controller) deleteUser(c echo.Context) (err error) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err = u.user.DeleteUser(id); err != nil {
		return err
	}

	return resp.OperateSuccess(c)
}
