/*
 * Created by lintao on 2023/7/18 下午3:56
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package resp

import (
	"net/http"
	"reflect"
	"time"

	"log/slog"

	"github.com/NSObjects/echo-admin/internal/code"
	"github.com/NSObjects/echo-admin/internal/log"
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

// ErrorResponse 统一错误响应结构
type ErrorResponse struct {
	Code      int    `json:"code"`       // 业务错误码
	Message   string `json:"message"`    // 错误消息
	RequestID string `json:"request_id"` // 请求ID（用于追踪）
	Timestamp int64  `json:"timestamp"`  // 错误发生时间戳
}

// APIError 返回API错误
func APIError(err error, c echo.Context) error {
	if err == nil {
		return errors.New("error can't be nil")
	}

	// 解析错误码
	codeError := errors.ParseCoder(err)
	errorCode := codeError.Code()

	// 获取HTTP状态码（只支持200,400,401,403,404,500）
	httpStatus := code.HTTPStatus(errorCode)

	// 构建错误响应
	rjson := ErrorResponse{
		Code:      errorCode,
		Message:   codeError.String(),
		RequestID: getRequestID(c),
		Timestamp: time.Now().Unix(),
	}

	// 统一记录错误日志（所有错误都打印到日志）
	logError(err, errorCode, codeError.String(), rjson.RequestID, c)

	// 返回对应的HTTP状态码
	return c.JSON(httpStatus, rjson)
}

// logError 统一错误日志记录
func logError(err error, errorCode int, message, requestID string, c echo.Context) {
	// 构建基础日志字段
	fields := []slog.Attr{
		slog.Int("code", errorCode),
		slog.String("message", message),
		slog.String("request_id", requestID),
		slog.String("method", c.Request().Method),
		slog.String("uri", c.Request().RequestURI),
		slog.String("user_agent", c.Request().UserAgent()),
		slog.String("error", err.Error()),
	}

	// 根据错误类型选择日志级别
	if code.IsInternalError(errorCode) {
		// 内部错误：使用Error级别，记录详细信息
		log.Error("Internal Error", fields...)
	} else {
		// 业务错误：使用Warn级别，记录业务信息
		log.Warn("Business Error", fields...)
	}
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

// getRequestID 获取请求ID，用于错误追踪
func getRequestID(c echo.Context) string {
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
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString 生成随机字符串
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
