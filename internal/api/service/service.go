/*
 * Created by lintao on 2023/7/26 下午2:22
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package service

import (
	"net/http/httptest"

	"github.com/labstack/echo/v4"
)

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
