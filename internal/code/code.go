/*
 * Application Error Codes
 * 应用业务错误码定义
 *
 * This package wraps go-kit/code to provide application-specific error codes.
 * Generated code from muban/codegen will register additional errors here.
 */

package code

import (
	"github.com/NSObjects/go-kit/code"
	"github.com/NSObjects/go-kit/errors"
)

// Re-export common error codes from go-kit/code
const (
	ErrSuccess        = code.ErrSuccess
	ErrUnknown        = code.ErrUnknown
	ErrBind           = code.ErrBind
	ErrValidation     = code.ErrValidation
	ErrTokenInvalid   = code.ErrTokenInvalid
	ErrDatabase       = code.ErrDatabase
	ErrRedis          = code.ErrRedis
	ErrKafka          = code.ErrKafka
	ErrBadRequest     = code.ErrBadRequest
	ErrUnauthorized   = code.ErrUnauthorized
	ErrForbidden      = code.ErrForbidden
	ErrNotFound       = code.ErrNotFound
	ErrInternalServer = code.ErrInternalServer
)

// register wraps errors.Register for use in generated code.
// This function is called by codegen-generated init() functions.
func register(c int, httpStatus int, message string) {
	errors.Register(c, httpStatus, message)
}

// WrapDatabaseError wraps an error with database error code context.
// message is used directly as the error message (not formatted into "database %s failed").
func WrapDatabaseError(err error, message string) error {
	if err == nil {
		return nil
	}
	return errors.WrapCode(err, code.ErrDatabase, "%s", message)
}

// WrapValidationError wraps an error with validation error code context.
func WrapValidationError(err error, message string) error {
	return code.WrapValidationError(err, message)
}

// WrapNotFoundError wraps an error with not found error code context.
func WrapNotFoundError(err error, message string) error {
	return code.WrapNotFoundError(err, message)
}

// WrapUnauthorizedError wraps an error with unauthorized error code context.
func WrapUnauthorizedError(err error, message string) error {
	return code.WrapUnauthorizedError(err, message)
}

// WrapForbiddenError wraps an error with forbidden error code context.
func WrapForbiddenError(err error, message string) error {
	return code.WrapForbiddenError(err, message)
}

// WrapInternalServerError wraps an error with internal server error code context.
func WrapInternalServerError(err error, message string) error {
	return code.WrapInternalServerError(err, message)
}

// WrapBindError wraps an error with bind error code context.
func WrapBindError(err error, message string) error {
	return code.WrapBindError(err, message)
}

// WrapRedisError wraps an error with redis error code context.
// message is used directly as the error message (not formatted into "redis %s failed").
func WrapRedisError(err error, message string) error {
	if err == nil {
		return nil
	}
	return errors.WrapCode(err, code.ErrRedis, "%s", message)
}

// WrapKafkaError wraps an error with kafka error code context.
// message is used directly as the error message (not formatted into "kafka %s failed").
func WrapKafkaError(err error, message string) error {
	if err == nil {
		return nil
	}
	return errors.WrapCode(err, code.ErrKafka, "%s", message)
}

// NewValidationError creates a new validation error.
func NewValidationError(field, message string) error {
	return code.NewValidationError(field, message)
}

// NewNotFoundError creates a new not found error.
func NewNotFoundError(resource string) error {
	return code.NewNotFoundError(resource)
}
