/*
 * Generated from OpenAPI3 document
 * Module: User
 */

package service

import (
	"strconv"

	"github.com/NSObjects/go-template/internal/api/service/param"
	appuser "github.com/NSObjects/go-template/internal/application/user"
	"github.com/NSObjects/go-template/internal/code"
	"github.com/NSObjects/go-template/internal/resp"
	"github.com/NSObjects/go-template/internal/utils"
	"github.com/labstack/echo/v4"
)

type UserController struct {
	user appuser.Service
}

func NewUserController(svc appuser.Service) RegisterRouter {
	return &UserController{user: svc}
}

func (c *UserController) RegisterRouter(g *echo.Group, m ...echo.MiddlewareFunc) {
	g.GET("/users", c.ListUsers).Name = "获取用户列表"
	g.POST("/users", c.Create).Name = "创建用户"
	g.GET("/users/{id}", c.GetByID).Name = "获取用户详情"
	g.PUT("/users/{id}", c.Update).Name = "更新用户"
	g.DELETE("/users/{id}", c.Delete).Name = "删除用户"
}

func (c *UserController) ListUsers(ctx echo.Context) error {
	var req param.UserListUsersRequest
	if err := BindAndValidate(ctx, &req); err != nil {
		return err
	}

	bizCtx := utils.BuildContext(ctx)
	query := appuser.AssembleListUsersQuery(req)
	list, total, err := c.user.ListUsers(bizCtx, query)
	if err != nil {
		return err
	}

	response := appuser.AssembleUserListResponse(list)
	return resp.ListDataResponse(ctx, response, total)
}

func (c *UserController) Create(ctx echo.Context) error {
	var req param.UserCreateRequest
	if err := BindAndValidate(ctx, &req); err != nil {
		return err
	}

	aggregate, err := appuser.AssembleCreateUser(req)
	if err != nil {
		return code.WrapValidationError(err, "用户数据不合法")
	}

	bizCtx := utils.BuildContext(ctx)
	if err := c.user.Create(bizCtx, aggregate); err != nil {
		return err
	}

	return resp.OperateSuccess(ctx)
}

func (c *UserController) GetByID(ctx echo.Context) error {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)

	bizCtx := utils.BuildContext(ctx)
	domainID := appuser.AssembleUserID(id)
	result, err := c.user.GetByID(bizCtx, domainID)
	if err != nil {
		return err
	}

	response := appuser.AssembleUserDataResponse(result)
	return resp.OneDataResponse(ctx, response)
}

func (c *UserController) Update(ctx echo.Context) error {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)

	var req param.UserUpdateRequest
	if err := BindAndValidate(ctx, &req); err != nil {
		return err
	}

	aggregate, err := appuser.AssembleUpdateUser(id, req)
	if err != nil {
		return code.WrapValidationError(err, "用户数据不合法")
	}

	bizCtx := utils.BuildContext(ctx)
	if err := c.user.Update(bizCtx, aggregate); err != nil {
		return err
	}

	return resp.OperateSuccess(ctx)
}

func (c *UserController) Delete(ctx echo.Context) error {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)

	bizCtx := utils.BuildContext(ctx)
	domainID := appuser.AssembleUserID(id)
	if err := c.user.Delete(bizCtx, domainID); err != nil {
		return err
	}

	return resp.OperateSuccess(ctx)
}
