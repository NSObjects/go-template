/*
 *
 * api_test.go
 * apis
 *
 * Created by lintao on 2019-01-01 14:05
 * Copyright © 2017-2019 PYL. All rights reserved.
 *
 */

package apis

import (
	_ "go-template/init"
	"go-template/tools/db"
	"net/http/httptest"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/labstack/echo"
)

var e *echo.Echo

func init() {
	aa, _, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	database, err := db.NewTestDataSource(aa)
	if err != nil {
		panic(err)
	}
	apiServer := NewEchoServer(database)
	apiServer.RegisterRouter()
	apiServer.LoadMiddleware()
	e = apiServer.server
}

func request(method, path string, e *echo.Echo) (int, string) {
	req := httptest.NewRequest(method, path, nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec.Code, rec.Body.String()
}
