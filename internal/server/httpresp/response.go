/*
 * Created by lintao on 2023/7/18 下午3:56
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package httpresp

import (
	"errors"
	"net/http"
	"reflect"
	"time"

	"github.com/NSObjects/go-template/internal/code"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const successCode = 0

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

// ErrorResponse is the standard JSON error response for HTTP adapters.
type ErrorResponse struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
	Timestamp int64  `json:"timestamp"`
}

// APIError renders a project error as a JSON HTTP response.
func APIError(c echo.Context, err error) error {
	if err == nil {
		return errors.New("error can't be nil")
	}

	info := code.NewErrorInfo(err)
	rjson := ErrorResponse{
		Code:      info.Code,
		Message:   info.Message,
		RequestID: RequestID(c),
		Timestamp: time.Now().Unix(),
	}

	return c.JSON(code.HTTPStatus(info.Code), rjson)
}

func OperateSuccess(c echo.Context) error {
	rjson := DataResponse{
		Code: successCode,
		Msg:  "success",
		Data: map[string]interface{}{},
	}

	return c.JSON(http.StatusOK, rjson)
}

func ListDataResponse(c echo.Context, arr interface{}, total int64) error {
	if arr == nil {
		arr = make([]interface{}, 0)
	} else if reflect.ValueOf(arr).IsNil() {
		arr = make([]interface{}, 0)
	}

	r := ListResponse{
		Code: successCode,
		Msg:  "success",
		Data: ListData{
			List:  arr,
			Total: total,
		},
	}

	return c.JSONPretty(http.StatusOK, r, "  ")
}

func OneDataResponse(c echo.Context, data interface{}) error {
	r := DataResponse{
		Code: successCode,
		Msg:  "success",
		Data: data,
	}

	return c.JSON(http.StatusOK, r)
}

// RequestID returns the request ID used in HTTP error responses.
func RequestID(c echo.Context) string {
	if requestID := c.Request().Header.Get("X-Request-ID"); requestID != "" {
		return requestID
	}

	if requestID := c.Response().Header().Get("X-Request-ID"); requestID != "" {
		return requestID
	}

	requestID := generateRequestID()
	c.Response().Header().Set("X-Request-ID", requestID)
	return requestID
}

func generateRequestID() string {
	return uuid.NewString()
}
