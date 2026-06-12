package http

import (
	"context"
	"strconv"

	"github.com/NSObjects/go-template/internal/code"
	"github.com/NSObjects/go-template/internal/utils"
	"github.com/labstack/echo/v4"
)

// BindAndValidate binds request data and runs the configured Echo validator.
func BindAndValidate(ctx echo.Context, obj any) error {
	if err := ctx.Bind(obj); err != nil {
		return code.WrapError(err, code.ErrBind, "bind request failed")
	}
	if err := ctx.Validate(obj); err != nil {
		return code.WrapError(err, code.ErrValidation, "validation failed")
	}
	return nil
}

// RequestContext returns the request-scoped context enriched by server middleware.
func RequestContext(ctx echo.Context) context.Context {
	return utils.BuildContext(ctx)
}

// PathInt64 parses an int64 path parameter.
func PathInt64(ctx echo.Context, name string) (int64, error) {
	id, err := strconv.ParseInt(ctx.Param(name), 10, 64)
	if err != nil {
		return 0, code.WrapBadRequestError(err, "invalid "+name)
	}
	return id, nil
}
