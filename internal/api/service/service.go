package service

import (
	"net/http/httptest"

	"github.com/NSObjects/echo-admin/internal/code"
	"github.com/labstack/echo/v4"
	"github.com/marmotedu/errors"
	"go.uber.org/fx"
)

var Model = fx.Options(
	fx.Provide(AsRoute(NewUserController)),
)

func AsRoute(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(RegisterRouter)),
		fx.ResultTags(`group:"routes"`),
	)
}

func BindAndValidate(ctx echo.Context, obj any) error {
	if err := ctx.Bind(obj); err != nil {
		return errors.WrapC(err, code.ErrBind, err.Error())
	}

	if err := ctx.Validate(obj); err != nil {
		return errors.WrapC(err, code.ErrValidation, err.Error())
	}

	return nil
}

type RegisterRouter interface {
	RegisterRouter(s *echo.Group, middlewareFunc ...echo.MiddlewareFunc)
}

//func testServer() *echo.Echo {
//	aa, _, err := sqlmock.New()
//	if err != nil {
//		panic(err)
//	}
//	database, err := db.NewTestDataSource(aa)
//	if err != nil {
//		panic(err)
//	}
//	apiServer := server.NewEchoServer(database)
//	return apiServer.Server()
//}

func Request(method, path string) (int, string) {
	//req := httptest.NewRequest(method, path, nil)
	rec := httptest.NewRecorder()
	//testServer().ServeHTTP(rec, req)
	return rec.Code, rec.Body.String()
}
