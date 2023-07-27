/*
 * Created by lintao on 2023/7/18 下午3:56
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package resp

type StatusCode int

const (
	StatusOK                 StatusCode = 0
	StatusDBErr              StatusCode = 100001
	StatusParamErr           StatusCode = 100002
	StatusAuth               StatusCode = 100003
	StatusServiceUnavailable StatusCode = 100004
	StatusAlterMsg           StatusCode = 100005
)

var statusCodeMsg = map[StatusCode]string{
	// 前三种错误不能在前端显示， 只在调试或者打日志中使用
	StatusDBErr:              "数据库错误",
	StatusParamErr:           "请求参数错误",
	StatusServiceUnavailable: "未知错误",
	// ----------------------------------------
	StatusAuth:     "重新登录",
	StatusOK:       "操作成功",
	StatusAlterMsg: "操作错误", // 用户错误操作流程， 可以直接在前端显示的错误
}
