package code

import (
	"fmt"
)

// AppError is the project error value used to carry a business code,
// a safe client-facing message, an internal diagnostic detail, and a cause.
type AppError struct {
	code    int
	message string
	detail  string
	cause   error
}

// Error returns the safe message. Use Detail for internal diagnostics.
func (e *AppError) Error() string {
	if e == nil {
		return ""
	}
	return e.message
}

// Unwrap returns the underlying cause.
func (e *AppError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.cause
}

// Code returns the project error code.
func (e *AppError) Code() int {
	if e == nil {
		return ErrUnknown
	}
	return e.code
}

// Message returns the client-facing message.
func (e *AppError) Message() string {
	if e == nil {
		return ""
	}
	return e.message
}

// Detail returns the internal diagnostic detail for logs.
func (e *AppError) Detail() string {
	if e == nil {
		return ""
	}
	return e.detail
}

// NewError 创建带错误码的错误
func NewError(code int, message string) error {
	return newAppError(code, message, "", nil)
}

// NewErrorf 创建带错误码的格式化错误
func NewErrorf(code int, format string, args ...interface{}) error {
	return NewError(code, fmt.Sprintf(format, args...))
}

// WrapError 包装错误并添加错误码
func WrapError(err error, code int, message string) error {
	return wrapOrNew(err, code, message)
}

// WrapErrorf 包装错误并添加错误码（格式化）
func WrapErrorf(err error, code int, format string, args ...interface{}) error {
	return wrapOrNewf(err, code, format, args...)
}

func wrapIfError(err error, code int, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	message := registeredMessage(code)
	detail := fmt.Sprintf(format, args...)
	return newAppError(code, message, joinDetail(detail, err), err)
}

func wrapOrNew(err error, code int, message string) error {
	return newAppError(code, message, joinDetail(message, err), err)
}

func wrapOrNewf(err error, code int, format string, args ...interface{}) error {
	message := fmt.Sprintf(format, args...)
	return newAppError(code, message, joinDetail(message, err), err)
}

func newAppError(code int, message, detail string, cause error) error {
	if message == "" {
		message = registeredMessage(code)
	}
	if detail == "" {
		detail = message
	}
	if IsServerError(code) {
		message = registeredMessage(code)
	}
	return &AppError{
		code:    code,
		message: message,
		detail:  detail,
		cause:   cause,
	}
}

func registeredMessage(code int) string {
	if coder, ok := Lookup(code); ok {
		return coder.String()
	}
	return "Internal server error"
}

func joinDetail(detail string, cause error) string {
	if cause == nil {
		return detail
	}
	causeDetail := fmt.Sprintf("%+v", cause)
	if appErr, ok := ParseError(cause); ok && appErr.Detail() != "" {
		causeDetail = appErr.Detail()
	}
	if detail == "" {
		return causeDetail
	}
	return fmt.Sprintf("%s: %s", detail, causeDetail)
}

// ========== 数据源底层错误包装函数 ==========

// WrapDatabaseError 包装数据库错误
func WrapDatabaseError(err error, operation string) error {
	return wrapIfError(err, ErrDatabase, "database %s failed", operation)
}

// WrapRedisError 包装Redis错误
func WrapRedisError(err error, operation string) error {
	return wrapIfError(err, ErrRedis, "redis %s failed", operation)
}

// WrapKafkaError 包装Kafka错误
func WrapKafkaError(err error, operation string) error {
	return wrapIfError(err, ErrKafka, "kafka %s failed", operation)
}

// WrapExternalError 包装第三方服务错误
func WrapExternalError(err error, service, operation string) error {
	return wrapIfError(err, ErrExternalService, "external service %s %s failed", service, operation)
}

// ========== biz层HTTP错误包装函数 ==========

// WrapBadRequestError 包装400错误
func WrapBadRequestError(err error, message string) error {
	return wrapOrNew(err, ErrBadRequest, message)
}

// WrapUnauthorizedError 包装401错误
func WrapUnauthorizedError(err error, message string) error {
	return wrapOrNew(err, ErrUnauthorized, message)
}

// WrapForbiddenError 包装403错误
func WrapForbiddenError(err error, message string) error {
	return wrapOrNew(err, ErrForbidden, message)
}

// WrapNotFoundError 包装404错误
func WrapNotFoundError(err error, message string) error {
	return wrapOrNew(err, ErrNotFound, message)
}

// WrapInternalServerError 包装500错误
func WrapInternalServerError(err error, message string) error {
	return wrapOrNew(err, ErrInternalServer, message)
}

// ========== 框架通用错误创建函数 ==========

// NewValidationError 验证错误
func NewValidationError(field, message string) error {
	return newAppError(ErrValidation, message, fmt.Sprintf("validation failed for field %s: %s", field, message), nil)
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
