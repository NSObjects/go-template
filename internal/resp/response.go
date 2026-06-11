/*
 * Created by lintao on 2023/7/18 下午3:56
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package resp

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

// ErrorResponse 统一错误响应结构
type ErrorResponse struct {
	Code      int    `json:"code"`       // 业务错误码
	Message   string `json:"message"`    // 错误消息
	RequestID string `json:"request_id"` // 请求ID（用于追踪）
	Timestamp int64  `json:"timestamp"`  // 错误发生时间戳
}

// APIError 返回API错误
func APIError(c echo.Context, err error) error {
	if err == nil {
		return errors.New("error can't be nil")
	}

	info := code.NewErrorInfo(err)

	// 构建错误响应
	rjson := ErrorResponse{
		Code:      info.Code,
		Message:   info.Message,
		RequestID: RequestID(c),
		Timestamp: time.Now().Unix(),
	}

	// 返回对应的HTTP状态码
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

// RequestID 获取请求ID，用于错误追踪。
func RequestID(c echo.Context) string {
	// 优先从请求头获取
	if requestID := c.Request().Header.Get("X-Request-ID"); requestID != "" {
		return requestID
	}

	// 从Echo上下文获取
	if requestID := c.Response().Header().Get("X-Request-ID"); requestID != "" {
		return requestID
	}

	// 生成新的请求ID
	requestID := generateRequestID()
	c.Response().Header().Set("X-Request-ID", requestID)
	return requestID
}

// generateRequestID 生成请求ID
func generateRequestID() string {
	return uuid.NewString()
}
