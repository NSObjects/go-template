/*
 * Generated from OpenAPI3 document
 * Module: user
 */

package service

import (
	"github.com/NSObjects/echo-admin/internal/api/biz"
	"github.com/NSObjects/echo-admin/internal/api/service/param"
	"github.com/NSObjects/echo-admin/internal/resp"
	"github.com/NSObjects/echo-admin/internal/utils"
	"github.com/labstack/echo/v4"
)

type UserController struct {
	user biz.UserUseCase
}

func NewUserController(h biz.UserUseCase) RegisterRouter {
	return &UserController{user: h}
}

func (c *UserController) RegisterRouter(g *echo.Group, m ...echo.MiddlewareFunc) {
	g.GET("/user/{id}", c.getByID).Name = "根据id查询某个用户"
	g.PUT("/user/{id}", c.update).Name = "更新用户数据"
	g.DELETE("/user/{id}", c.delete).Name = "删除用户"
	g.GET("/users", c.list).Name = "查询用户"
	g.POST("/users", c.create).Name = "创建用户"
}

// TODO: 实现控制器方法

func (c *UserController) getByID(ctx echo.Context) error {
	// TODO: 绑定和验证请求参数
	var req param.UserGetByIDRequest
	if err := BindAndValidate(ctx, &req); err != nil {
		return err
	}

	// 调用业务逻辑 - 构造包含链路追踪信息的context
	bizCtx := utils.BuildContext(ctx)
	result, err := c.user.GetByID(bizCtx, req)
	if err != nil {
		return err
	}

	// 返回单个数据
	return resp.OneDataResponse(result, ctx)
}
func (c *UserController) update(ctx echo.Context) error {
	// TODO: 绑定和验证请求参数
	var req param.UserUpdateRequest
	if err := BindAndValidate(ctx, &req); err != nil {
		return err
	}

	// 调用业务逻辑 - 构造包含链路追踪信息的context
	bizCtx := utils.BuildContext(ctx)
	result, err := c.user.Update(bizCtx, req)
	if err != nil {
		return err
	}

	// 返回单个数据
	return resp.OneDataResponse(result, ctx)
}
func (c *UserController) delete(ctx echo.Context) error {
	// TODO: 绑定和验证请求参数
	var req param.UserDeleteRequest
	if err := BindAndValidate(ctx, &req); err != nil {
		return err
	}

	// 调用业务逻辑 - 构造包含链路追踪信息的context
	bizCtx := utils.BuildContext(ctx)
	err := c.user.Delete(bizCtx, req)
	if err != nil {
		return err
	}

	// 返回操作成功
	return resp.OperateSuccess(ctx)
}
func (c *UserController) list(ctx echo.Context) error {
	// TODO: 绑定和验证请求参数
	var req param.UserListRequest
	if err := BindAndValidate(ctx, &req); err != nil {
		return err
	}

	// 调用业务逻辑 - 构造包含链路追踪信息的context
	bizCtx := utils.BuildContext(ctx)
	list, total, err := c.user.List(bizCtx, req)
	if err != nil {
		return err
	}

	// 返回列表数据
	return resp.ListDataResponse(list, total, ctx)
}
func (c *UserController) create(ctx echo.Context) error {
	// TODO: 绑定和验证请求参数
	var req param.UserCreateRequest
	if err := BindAndValidate(ctx, &req); err != nil {
		return err
	}

	// 调用业务逻辑 - 构造包含链路追踪信息的context
	bizCtx := utils.BuildContext(ctx)
	result, err := c.user.Create(bizCtx, req)
	if err != nil {
		return err
	}

	// 返回单个数据
	return resp.OneDataResponse(result, ctx)
}
