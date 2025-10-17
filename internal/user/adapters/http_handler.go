package adapters

import (
	"github.com/NSObjects/go-template/internal/pkg/code"
	"github.com/NSObjects/go-template/internal/pkg/utils"
	"github.com/NSObjects/go-template/internal/shared/ports/resp"
	"github.com/NSObjects/go-template/internal/user/app"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

// RegisterRouter 路由注册接口
type RegisterRouter interface {
	RegisterRouter(s *echo.Group, middlewareFunc ...echo.MiddlewareFunc)
}

// UserController HTTP控制器（入站适配器）
type UserController struct {
	userService app.UserService
}

// NewUserController 创建用户控制器
func NewUserController(userService app.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// RegisterRouter 注册路由（实现RegisterRouter接口）
func (c *UserController) RegisterRouter(g *echo.Group, middlewareFunc ...echo.MiddlewareFunc) {
	g.GET("/users", c.ListUsers).Name = "获取用户列表"
	g.POST("/users", c.CreateUser).Name = "创建用户"
	g.GET("/users/:id", c.GetUser).Name = "获取用户详情"
	g.PUT("/users/:id", c.UpdateUser).Name = "更新用户"
	g.DELETE("/users/:id", c.DeleteUser).Name = "删除用户"
	g.POST("/users/:id/suspend", c.SuspendUser).Name = "暂停用户"
	g.POST("/users/:id/activate", c.ActivateUser).Name = "激活用户"
}

// ListUsers 获取用户列表
func (c *UserController) ListUsers(ctx echo.Context) error {
	var req app.ListUsersRequest
	if err := ctx.Bind(&req); err != nil {
		return code.WrapValidationError(err, "bind request failed")
	}

	if err := ctx.Validate(&req); err != nil {
		return code.WrapValidationError(err, "validation failed")
	}

	req.CalculateOffset()

	bizCtx := utils.BuildContext(ctx)
	response, err := c.userService.ListUsers(bizCtx, req)
	if err != nil {
		return err
	}

	return resp.ListDataResponse(ctx, response.Users, response.Total)
}

// CreateUser 创建用户
func (c *UserController) CreateUser(ctx echo.Context) error {
	var req app.CreateUserRequest
	if err := ctx.Bind(&req); err != nil {
		return code.WrapValidationError(err, "bind request failed")
	}

	if err := ctx.Validate(&req); err != nil {
		return code.WrapValidationError(err, "validation failed")
	}

	bizCtx := utils.BuildContext(ctx)
	user, err := c.userService.CreateUser(bizCtx, req)
	if err != nil {
		return err
	}

	return resp.OneDataResponse(ctx, user)
}

// GetUser 获取用户详情
func (c *UserController) GetUser(ctx echo.Context) error {
	userID := ctx.Param("id")
	if userID == "" {
		return code.WrapValidationError(nil, "user ID is required")
	}

	bizCtx := utils.BuildContext(ctx)
	user, err := c.userService.GetUser(bizCtx, userID)
	if err != nil {
		return err
	}

	return resp.OneDataResponse(ctx, user)
}

// UpdateUser 更新用户
func (c *UserController) UpdateUser(ctx echo.Context) error {
	userID := ctx.Param("id")
	if userID == "" {
		return code.WrapValidationError(nil, "user ID is required")
	}

	var req app.UpdateUserRequest
	if err := ctx.Bind(&req); err != nil {
		return code.WrapValidationError(err, "bind request failed")
	}

	if err := ctx.Validate(&req); err != nil {
		return code.WrapValidationError(err, "validation failed")
	}

	bizCtx := utils.BuildContext(ctx)
	err := c.userService.UpdateUser(bizCtx, userID, req)
	if err != nil {
		return err
	}

	return resp.OperateSuccess(ctx)
}

// DeleteUser 删除用户
func (c *UserController) DeleteUser(ctx echo.Context) error {
	userID := ctx.Param("id")
	if userID == "" {
		return code.WrapValidationError(nil, "user ID is required")
	}

	bizCtx := utils.BuildContext(ctx)
	err := c.userService.DeleteUser(bizCtx, userID)
	if err != nil {
		return err
	}

	return resp.OperateSuccess(ctx)
}

// SuspendUser 暂停用户
func (c *UserController) SuspendUser(ctx echo.Context) error {
	userID := ctx.Param("id")
	if userID == "" {
		return code.WrapValidationError(nil, "user ID is required")
	}

	bizCtx := utils.BuildContext(ctx)
	err := c.userService.SuspendUser(bizCtx, userID)
	if err != nil {
		return err
	}

	return resp.OperateSuccess(ctx)
}

// ActivateUser 激活用户
func (c *UserController) ActivateUser(ctx echo.Context) error {
	userID := ctx.Param("id")
	if userID == "" {
		return code.WrapValidationError(nil, "user ID is required")
	}

	bizCtx := utils.BuildContext(ctx)
	err := c.userService.ActivateUser(bizCtx, userID)
	if err != nil {
		return err
	}

	return resp.OperateSuccess(ctx)
}

// UserControllerModule 用户控制器模块
var UserControllerModule = fx.Options(
	fx.Provide(NewUserController),
)
