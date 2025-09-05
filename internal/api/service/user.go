/*
 * Generated from OpenAPI3 document
 * Module: User
 */

package service

import (
	"strconv"
	"github.com/NSObjects/go-template/internal/api/biz"
	"github.com/NSObjects/go-template/internal/api/service/param"
	"github.com/NSObjects/go-template/internal/resp"
	"github.com/NSObjects/go-template/internal/utils"
	"github.com/labstack/echo/v4"
)

type UserController struct {
	user biz.UserUseCase
}

func NewUserController(h biz.UserUseCase) RegisterRouter {
	return &UserController{user: h}
}

func (c *UserController) RegisterRouter(g *echo.Group, m ...echo.MiddlewareFunc) {

	g.GET("/users", c.ListUsers).Name = "获取用户列表"

	g.POST("/users", c.Create).Name = "创建用户"

	g.GET("/users/{id}", c.GetByID).Name = "获取用户详情"

	g.PUT("/users/{id}", c.Update).Name = "更新用户"

	g.DELETE("/users/{id}", c.Delete).Name = "删除用户"

}

// TODO: 实现控制器方法



func (c *UserController) ListUsers(ctx echo.Context) error {
	// TODO: 绑定和验证请求参数
	var req param.UserListUsersRequest
	if err := BindAndValidate(ctx, &req); err != nil {
		return err
	}
	
	// 调用业务逻辑 - 构造包含链路追踪信息的context
	bizCtx := utils.BuildContext(ctx)
	list, total, err := c.user.ListUsers(bizCtx, req)
	if err != nil {
		return err
	}
	
	// 返回列表数据
	return resp.ListDataResponse(list, total, ctx)
}



func (c *UserController) Create(ctx echo.Context) error {
	// TODO: 绑定和验证请求参数
	var req param.UserCreateRequest
	if err := BindAndValidate(ctx, &req); err != nil {
		return err
	}
	
	// 调用业务逻辑 - 构造包含链路追踪信息的context
	bizCtx := utils.BuildContext(ctx)
	if err := c.user.Create(bizCtx, req); err != nil {
		return err
	}
	
	// 返回操作成功
	return resp.OperateSuccess(ctx)
}



func (c *UserController) GetByID(ctx echo.Context) error {
	// 获取路径参数
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	
	// 调用业务逻辑 - 构造包含链路追踪信息的context
	bizCtx := utils.BuildContext(ctx)
	result, err := c.user.GetByID(bizCtx, id)
	if err != nil {
		return err
	}
	
	// 返回数据
	return resp.OneDataResponse(result, ctx)
}



func (c *UserController) Update(ctx echo.Context) error {
	// 获取路径参数
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	
	// 绑定和验证请求体参数
	var req param.UserUpdateRequest
	if err := BindAndValidate(ctx, &req); err != nil {
		return err
	}
	
	// 调用业务逻辑 - 构造包含链路追踪信息的context
	bizCtx := utils.BuildContext(ctx)
	if err := c.user.Update(bizCtx, id, req); err != nil {
		return err
	}
	
	// 返回操作成功
	return resp.OperateSuccess(ctx)
}



func (c *UserController) Delete(ctx echo.Context) error {
	// 获取路径参数
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	
	// 调用业务逻辑 - 构造包含链路追踪信息的context
	bizCtx := utils.BuildContext(ctx)
	if err := c.user.Delete(bizCtx, id); err != nil {
		return err
	}
	
	// 返回操作成功
	return resp.OperateSuccess(ctx)
}


