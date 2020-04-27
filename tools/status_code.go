/*
 *
 * status_code.go
 * api_helper
 *
 * Created by lintao on 2019-01-26 16:57
 * Copyright © 2017-2019 PYL. All rights reserved.
 *
 */

package tools

type ResponetsCode int

const (
	StatusDBErr              ResponetsCode = 407
	StatusParamErr           ResponetsCode = 405
	StatusRelogin            ResponetsCode = 500
	StatusServiceUnavailable ResponetsCode = 401
	StatusOK                 ResponetsCode = 200
	StatusAlerMsg            ResponetsCode = 201
)

var statusCodeMsg = map[ResponetsCode]string{
	// 前三种错误不能在前端显示， 只在调试或者打日志中使用
	StatusDBErr:              "数据库错误",
	StatusParamErr:           "请求参数错误",
	StatusServiceUnavailable: "未知错误",
	// ----------------------------------------
	StatusRelogin: "重新登录",
	StatusOK:      "操作成功",
	StatusAlerMsg: "操作错误", // 用户错误操作流程， 可以直接在前端显示的错误
}
