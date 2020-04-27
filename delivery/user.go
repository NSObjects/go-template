/*
 *
 * user.go
 * apis
 *
 * Created by lintao on 2019-01-29 16:17
 * Copyright © 2017-2019 PYL. All rights reserved.
 *
 */

package delivery

import (
	"go-template/usecase"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	ucase usecase.UserUsecase
}

func (this *UserController) RegisterRouter(s *echo.Group, middlewareFunc ...echo.MiddlewareFunc) {
	s.GET("/users", this.getUser).Name = "用户查询"
	s.POST("/users", this.createUser).Name = "创建用户"
	s.DELETE("/users/:id", this.deleteUser).Name = "删除用户"
	s.PUT("/users/:id", this.updateUser).Name = "更新用户"
	s.GET("/users/:id", this.getUserDetail).Name = "获取某个用户信息"
}

func NewUserController(ucase usecase.UserUsecase) *UserController {
	return &UserController{
		ucase: ucase,
	}
}

func (this *UserController) getUser(c echo.Context) (err error) {
	return this.ucase.GetUser(c)
}

func (this *UserController) getUserDetail(c echo.Context) (err error) {
	return this.ucase.GetUserDetail(c)
}

func (this *UserController) createUser(c echo.Context) (err error) {
	return this.ucase.CreateUser(c)
}

func (this *UserController) updateUser(c echo.Context) (err error) {
	return this.ucase.UpdateUser(c)
}

func (this *UserController) deleteUser(c echo.Context) (err error) {
	return this.ucase.DeleteUser(c)
}
