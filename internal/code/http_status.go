package code

import "net/http"

// HTTPStatus 根据错误码返回对应的HTTP状态码
// 只支持标准的6个状态码：200, 400, 401, 403, 404, 500
func HTTPStatus(code int) int {
	if code == 0 {
		return http.StatusOK
	}

	if coder, ok := Lookup(code); ok {
		return coder.HTTPStatus()
	}

	return http.StatusInternalServerError
}

// IsClientError 判断是否为客户端错误（4xx）
func IsClientError(code int) bool {
	status := HTTPStatus(code)
	return status >= 400 && status < 500
}

// IsServerError 判断是否为服务器错误（5xx）
func IsServerError(code int) bool {
	status := HTTPStatus(code)
	return status >= 500
}

// IsInternalError 判断是否为内部错误（需要记录详细日志）
func IsInternalError(code int) bool {
	// 内部错误码范围：100101-100199
	return code >= 100101 && code <= 100199
}
