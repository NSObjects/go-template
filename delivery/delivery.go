/*
 *
 * delivery.go
 * delivery
 *
 * Created by lintao on 2020/4/28 1:59 下午
 * Copyright © 2020-2020 LINTAO. All rights reserved.
 *
 */

package delivery

import (
	"go-template/tools/db"
	"net/http/httptest"

	"github.com/labstack/echo/v4"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

type RegisterRouter interface {
	RegisterRouter(s *echo.Group, middlewareFunc ...echo.MiddlewareFunc)
}

func testServer() *echo.Echo {
	aa, _, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	database, err := db.NewTestDataSource(aa)
	if err != nil {
		panic(err)
	}
	apiServer := NewEchoServer(database)
	return apiServer.Server()
}

func Request(method, path string) (int, string) {
	req := httptest.NewRequest(method, path, nil)
	rec := httptest.NewRecorder()
	testServer().ServeHTTP(rec, req)
	return rec.Code, rec.Body.String()
}
