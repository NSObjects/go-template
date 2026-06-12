package user

import (
	"strconv"

	"github.com/NSObjects/go-template/internal/code"
	"github.com/NSObjects/go-template/internal/resp"
	"github.com/NSObjects/go-template/internal/utils"
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
	if err := bindAndValidate(ctx, &req); err != nil {
		return err
	}

	list, total, err := c.user.ListUsers(utils.BuildContext(ctx), req)
	if err != nil {
		return err
	}
	return resp.ListDataResponse(ctx, list, total)
}

func (c *controller) Create(ctx echo.Context) error {
	var req CreateRequest
	if err := bindAndValidate(ctx, &req); err != nil {
		return err
	}

	if err := c.user.Create(utils.BuildContext(ctx), req); err != nil {
		return err
	}
	return resp.OperateSuccess(ctx)
}

func (c *controller) GetByID(ctx echo.Context) error {
	id, err := parseID(ctx)
	if err != nil {
		return err
	}

	result, err := c.user.GetByID(utils.BuildContext(ctx), id)
	if err != nil {
		return err
	}
	return resp.OneDataResponse(ctx, result)
}

func (c *controller) Update(ctx echo.Context) error {
	id, err := parseID(ctx)
	if err != nil {
		return err
	}

	var req UpdateRequest
	if err := bindAndValidate(ctx, &req); err != nil {
		return err
	}

	if err := c.user.Update(utils.BuildContext(ctx), id, req); err != nil {
		return err
	}
	return resp.OperateSuccess(ctx)
}

func (c *controller) Delete(ctx echo.Context) error {
	id, err := parseID(ctx)
	if err != nil {
		return err
	}

	if err := c.user.Delete(utils.BuildContext(ctx), id); err != nil {
		return err
	}
	return resp.OperateSuccess(ctx)
}

func bindAndValidate(ctx echo.Context, obj any) error {
	if err := ctx.Bind(obj); err != nil {
		return code.WrapError(err, code.ErrBind, "bind request failed")
	}
	if err := ctx.Validate(obj); err != nil {
		return code.WrapError(err, code.ErrValidation, "validation failed")
	}
	return nil
}

func parseID(ctx echo.Context) (int64, error) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return 0, code.WrapBadRequestError(err, "invalid user id")
	}
	return id, nil
}
