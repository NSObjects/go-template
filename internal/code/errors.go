package code

import (
	"fmt"

	"github.com/marmotedu/errors"
)

// NewError 创建带错误码的错误
func NewError(code int, message string) error {
	return errors.WithCode(code, "%s", message)
}

// NewErrorf 创建带错误码的格式化错误
func NewErrorf(code int, format string, args ...interface{}) error {
	return errors.WithCode(code, format, args...)
}

// WrapError 包装错误并添加错误码
func WrapError(err error, code int, message string) error {
	if err == nil {
		return NewError(code, message)
	}
	return errors.WrapC(err, code, "%s", message)
}

// WrapErrorf 包装错误并添加错误码（格式化）
func WrapErrorf(err error, code int, format string, args ...interface{}) error {
	return errors.WrapC(err, code, format, args...)
}

// ========== 数据源底层错误包装函数 ==========

// WrapDatabaseError 包装数据库错误
func WrapDatabaseError(err error, operation string) error {
	if err == nil {
		return nil
	}
	message := fmt.Sprintf("database %s failed", operation)
	return errors.WrapC(err, ErrDatabase, "%s", message)
}

// WrapRedisError 包装Redis错误
func WrapRedisError(err error, operation string) error {
	if err == nil {
		return nil
	}
	message := fmt.Sprintf("redis %s failed", operation)
	return errors.WrapC(err, ErrRedis, "%s", message)
}

// WrapKafkaError 包装Kafka错误
func WrapKafkaError(err error, operation string) error {
	if err == nil {
		return nil
	}
	message := fmt.Sprintf("kafka %s failed", operation)
	return errors.WrapC(err, ErrKafka, "%s", message)
}

// WrapExternalError 包装第三方服务错误
func WrapExternalError(err error, service, operation string) error {
	if err == nil {
		return nil
	}
	message := fmt.Sprintf("external service %s %s failed", service, operation)
	return errors.WrapC(err, ErrExternalService, "%s", message)
}

// ========== biz层HTTP错误包装函数 ==========

// WrapBadRequestError 包装400错误
func WrapBadRequestError(err error, message string) error {
	if err == nil {
		return NewError(ErrBadRequest, message)
	}
	return errors.WrapC(err, ErrBadRequest, "%s", message)
}

// WrapUnauthorizedError 包装401错误
func WrapUnauthorizedError(err error, message string) error {
	if err == nil {
		return NewError(ErrUnauthorized, message)
	}
	return errors.WrapC(err, ErrUnauthorized, "%s", message)
}

// WrapForbiddenError 包装403错误
func WrapForbiddenError(err error, message string) error {
	if err == nil {
		return NewError(ErrForbidden, message)
	}
	return errors.WrapC(err, ErrForbidden, "%s", message)
}

// WrapNotFoundError 包装404错误
func WrapNotFoundError(err error, message string) error {
	if err == nil {
		return NewError(ErrNotFound, message)
	}
	return errors.WrapC(err, ErrNotFound, "%s", message)
}

// WrapInternalServerError 包装500错误
func WrapInternalServerError(err error, message string) error {
	if err == nil {
		return NewError(ErrInternalServer, message)
	}
	return errors.WrapC(err, ErrInternalServer, "%s", message)
}

// ========== 框架通用错误创建函数 ==========

// NewValidationError 验证错误
func NewValidationError(field, message string) error {
	return NewErrorf(ErrValidation, "validation failed for field %s: %s", field, message)
}

// NewPermissionDeniedError 权限拒绝错误
func NewPermissionDeniedError(resource, action string) error {
	return NewErrorf(ErrPermissionDenied, "permission denied for %s on %s", action, resource)
}

// NewTokenInvalidError Token无效错误
func NewTokenInvalidError() error {
	return NewError(ErrTokenInvalid, "token is invalid")
}

// NewTokenExpiredError Token过期错误
func NewTokenExpiredError() error {
	return NewError(ErrExpired, "token is expired")
}

// NewUnauthorizedError 未授权错误
func NewUnauthorizedError() error {
	return NewError(ErrUnauthorized, "unauthorized")
}

// NewForbiddenError 禁止访问错误
func NewForbiddenError() error {
	return NewError(ErrForbidden, "forbidden")
}

// NewNotFoundError 资源不存在错误
func NewNotFoundError(resource string) error {
	return NewErrorf(ErrNotFound, "%s not found", resource)
}

// NewBadRequestError 请求错误
func NewBadRequestError(message string) error {
	return NewErrorf(ErrBadRequest, "bad request: %s", message)
}
