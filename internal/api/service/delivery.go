/*
 * Created by lintao on 2023/7/18 下午4:00
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package service

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/NSObjects/go-template/internal/api/data/db"
	"github.com/labstack/echo/v4"
	"net/http/httptest"
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
