/*
 * Created by lintao on 2023/7/26 下午2:22
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package service

import (
	"github.com/NSObjects/go-template/internal/code"
	"github.com/marmotedu/errors"
	"net/http/httptest"

	"go.uber.org/fx"

	"github.com/labstack/echo/v4"
)

var Model = fx.Options(
	fx.Provide(
		AsRoute(NewUserController),
	),
)

func AsRoute(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(RegisterRouter)),
		fx.ResultTags(`group:"routes"`),
	)
}

func BindAndValidate(obj any, ctx echo.Context) error {
	if err := ctx.Bind(&obj); err != nil {
		return errors.WrapC(err, code.ErrBind, err.Error())
	}

	if err := ctx.Validate(&obj); err != nil {
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
