/*
 * Created by lintao on 2023/7/18 下午3:56
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package resp

import (
	"net/http"
	"reflect"

	"github.com/NSObjects/go-template/internal/log"
	"github.com/labstack/echo/v4"
	"github.com/marmotedu/errors"
)

type ListResponse struct {
	Code int      `json:"code"`
	Msg  string   `json:"msg"`
	Data ListData `json:"data"`
}

type ListData struct {
	Total int64       `json:"total"`
	List  interface{} `json:"list" `
}

type DataResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// APIError 返回API错误
func APIError(err error, c echo.Context) error {
	if err == nil {
		return errors.New("error can't be nil")
	}
	var rjson struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	codeError := errors.ParseCoder(err)
	rjson.Code = codeError.Code()
	rjson.Msg = codeError.String()
	log.Error(err)
	return c.JSON(codeError.HTTPStatus(), rjson)
}

func OperateSuccess(c echo.Context) error {
	var rjson struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}

	rjson.Msg = "success"

	return c.JSON(http.StatusOK, rjson)
}

func ListDataResponse(arr interface{}, total int64, c echo.Context) error {
	if arr == nil {
		arr = make([]interface{}, 0)
	} else if reflect.ValueOf(arr).IsNil() {
		arr = make([]interface{}, 0)
	}

	r := ListResponse{
		Data: ListData{
			List:  arr,
			Total: total,
		},
	}

	return c.JSONPretty(http.StatusOK, r, "  ")
}

func OneDataResponse(data interface{}, c echo.Context) error {
	r := DataResponse{
		Data: data,
	}

	return c.JSON(http.StatusOK, r)
}
