/*
 * Created by lintao on 2023/7/18 下午4:00
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package service

import (
	"github.com/NSObjects/go-template/internal/api/biz"
	"github.com/labstack/echo/v4"
)

type UserController struct {
	user biz.UserUsecase
}

func (u *UserController) RegisterRouter(s *echo.Group, middlewareFunc ...echo.MiddlewareFunc) {
	s.GET("/users", u.getUser).Name = "用户查询"
	s.POST("/users", u.createUser).Name = "创建用户"
	s.DELETE("/users/:id", u.deleteUser).Name = "删除用户"
	s.PUT("/users/:id", u.updateUser).Name = "更新用户"
	s.GET("/users/:id", u.getUserDetail).Name = "获取某个用户信息"
}

func NewUserController(ucase biz.UserUsecase) *UserController {
	return &UserController{
		user: ucase,
	}
}

func (u *UserController) getUser(c echo.Context) (err error) {
	return u.user.GetUser(c)
}

func (u *UserController) getUserDetail(c echo.Context) (err error) {
	return u.user.GetUserDetail(c)
}

func (u *UserController) createUser(c echo.Context) (err error) {
	return u.user.CreateUser(c)
}

func (u *UserController) updateUser(c echo.Context) (err error) {
	return u.user.UpdateUser(c)
}

func (u *UserController) deleteUser(c echo.Context) (err error) {
	return u.user.DeleteUser(c)
}
