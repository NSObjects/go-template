package user

import (
	platformhttp "github.com/NSObjects/go-template/internal/platform/http"
	"github.com/NSObjects/go-template/internal/resp"
	"github.com/labstack/echo/v4"
)

type controller struct {
	user UseCase
}

func newController(useCase UseCase) *controller {
	return &controller{user: useCase}
}

func (c *controller) ListUsers(ctx echo.Context) error {
	var req ListUsersRequest
	if err := platformhttp.BindAndValidate(ctx, &req); err != nil {
		return err
	}

	list, total, err := c.user.ListUsers(platformhttp.RequestContext(ctx), req)
	if err != nil {
		return err
	}
	return resp.ListDataResponse(ctx, list, total)
}

func (c *controller) Create(ctx echo.Context) error {
	var req CreateRequest
	if err := platformhttp.BindAndValidate(ctx, &req); err != nil {
		return err
	}

	if err := c.user.Create(platformhttp.RequestContext(ctx), req); err != nil {
		return err
	}
	return resp.OperateSuccess(ctx)
}

func (c *controller) GetByID(ctx echo.Context) error {
	id, err := platformhttp.PathInt64(ctx, "id")
	if err != nil {
		return err
	}

	result, err := c.user.GetByID(platformhttp.RequestContext(ctx), id)
	if err != nil {
		return err
	}
	return resp.OneDataResponse(ctx, result)
}

func (c *controller) Update(ctx echo.Context) error {
	id, err := platformhttp.PathInt64(ctx, "id")
	if err != nil {
		return err
	}

	var req UpdateRequest
	if err := platformhttp.BindAndValidate(ctx, &req); err != nil {
		return err
	}

	if err := c.user.Update(platformhttp.RequestContext(ctx), id, req); err != nil {
		return err
	}
	return resp.OperateSuccess(ctx)
}

func (c *controller) Delete(ctx echo.Context) error {
	id, err := platformhttp.PathInt64(ctx, "id")
	if err != nil {
		return err
	}

	if err := c.user.Delete(platformhttp.RequestContext(ctx), id); err != nil {
		return err
	}
	return resp.OperateSuccess(ctx)
}
