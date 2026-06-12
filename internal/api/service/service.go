package service

import (
	"net/http/httptest"

	"github.com/NSObjects/go-template/internal/code"
	"github.com/labstack/echo/v4"
)

func BindAndValidate(ctx echo.Context, obj any) error {
	if err := ctx.Bind(obj); err != nil {
		return code.WrapError(err, code.ErrBind, "bind request failed")
	}

	if err := ctx.Validate(obj); err != nil {
		return code.WrapError(err, code.ErrValidation, "validation failed")
	}

	return nil
}

type RegisterRouter interface {
	RegisterRouter(s *echo.Group, middlewareFunc ...echo.MiddlewareFunc)
}

func Request(method, path string) (int, string) {
	//req := httptest.NewRequest(method, path, nil)
	rec := httptest.NewRecorder()
	//testServer().ServeHTTP(rec, req)
	return rec.Code, rec.Body.String()
}
