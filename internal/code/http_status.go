package code

import "net/http"

// HTTPStatus 根据错误码返回对应的HTTP状态码
// 只支持标准的6个状态码：200, 400, 401, 403, 404, 500
func HTTPStatus(code int) int {
	switch {
	// 成功
	case code == 0:
		return http.StatusOK // 200

	// HTTP状态码错误（直接映射）
	case code >= 100400 && code <= 100404:
		switch code {
		case ErrBadRequest:
			return http.StatusBadRequest // 400
		case ErrUnauthorized:
			return http.StatusUnauthorized // 401
		case ErrForbidden:
			return http.StatusForbidden // 403
		case ErrNotFound:
			return http.StatusNotFound // 404
		default:
			return http.StatusBadRequest // 400
		}

	// 客户端错误 4xx
	case code >= 100001 && code <= 100099: // 基本错误
		return http.StatusBadRequest // 400
	case code >= 100201 && code <= 100299: // 认证授权错误
		if code == ErrTokenInvalid || code == ErrExpired {
			return http.StatusUnauthorized // 401
		}
		return http.StatusForbidden // 403
	case code >= 100301 && code <= 100399: // 编解码错误
		return http.StatusBadRequest // 400
	case code >= 100401 && code <= 100499: // 验证错误
		return http.StatusBadRequest // 400

	// 服务器错误 5xx
	case code >= 100101 && code <= 100199: // 内部错误（数据库、Redis、Kafka等）
		return http.StatusInternalServerError // 500

	// 默认情况
	default:
		return http.StatusInternalServerError // 500
	}
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
