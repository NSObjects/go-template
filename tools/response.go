/*
 * Created by lintao on 2023/7/18 下午3:56
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package tools

import (
	"errors"
	"github.com/NSObjects/go-template/tools/log"
	"github.com/labstack/echo/v4"
	"net/http"
	"reflect"
)

type ListResponse struct {
	Code ResponetsCode `json:"code"`
	Msg  string        `json:"msg"`
	Data ListData      `json:"data"`
}

type ListData struct {
	Total int64       `json:"total"`
	Datas interface{} `json:"datas"`
}

type DataResponse struct {
	Code ResponetsCode `json:"code"`
	Msg  string        `json:"msg"`
	Data interface{}   `json:"data"`
}

func ApiError(err error, c echo.Context) error {

	if err == nil {
		return errors.New("error can't be nil")
	}

	log.ErrorSkip(2, err)
	var rjson struct {
		Code ResponetsCode `json:"code"`
		Msg  string        `json:"msg"`
	}

	if terr, ok := err.(*Error); ok {
		rjson.Code = terr.Code
		rjson.Msg = terr.Err.Error()
	} else {
		rjson.Code = StatusServiceUnavailable
		rjson.Msg = err.Error()
	}

	return c.JSON(http.StatusOK, rjson)
}

func OperateSuccess(c echo.Context) error {
	var rjson struct {
		Code ResponetsCode `json:"code"`
		Msg  string        `json:"msg"`
	}

	rjson.Code = StatusOK
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
			Datas: arr,
			Total: total,
		},
		Code: StatusOK,
	}

	return c.JSON(http.StatusOK, r)
}

func OneDataResponse(data interface{}, c echo.Context) error {
	r := DataResponse{
		Data: data,
		Code: StatusOK,
	}

	return c.JSON(http.StatusOK, r)
}
