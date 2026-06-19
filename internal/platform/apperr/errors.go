package apperr

import (
	stderrors "errors"
	"fmt"
)

// Error carries a code, safe message, diagnostic detail, and wrapped cause.
type Error struct {
	code    int
	message string
	detail  string
	cause   error
}

// Error returns the safe client-facing message.
func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	return e.message
}

// Unwrap returns the underlying cause.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.cause
}

// Code returns the application error code.
func (e *Error) Code() int {
	if e == nil {
		return ErrUnknown
	}
	return e.code
}

// Message returns the safe client-facing message.
func (e *Error) Message() string {
	if e == nil {
		return ""
	}
	return e.message
}

// Detail returns internal diagnostic detail for logs.
func (e *Error) Detail() string {
	if e == nil {
		return ""
	}
	return e.detail
}

// New creates an application error.
func New(code int, message string) error {
	return newError(code, message, "", nil)
}

// Newf creates a formatted application error.
func Newf(code int, format string, args ...interface{}) error {
	return New(code, fmt.Sprintf(format, args...))
}

// Wrap adds an application code and safe message to err.
func Wrap(err error, code int, message string) error {
	return wrapOrNew(err, code, message)
}

// Wrapf adds an application code and formatted safe message to err.
func Wrapf(err error, code int, format string, args ...interface{}) error {
	return wrapOrNewf(err, code, format, args...)
}

func wrapIfError(err error, code int, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	def := definitionFor(code)
	detail := fmt.Sprintf(format, args...)
	return newError(code, def.Message, joinDetail(detail, err), err)
}

func wrapOrNew(err error, code int, message string) error {
	if err == nil {
		return nil
	}
	return newError(code, message, joinDetail(message, err), err)
}

func wrapOrNewf(err error, code int, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	message := fmt.Sprintf(format, args...)
	return newError(code, message, joinDetail(message, err), err)
}

func newError(code int, message, detail string, cause error) error {
	def, registered := Lookup(code)
	if message == "" {
		if registered {
			message = def.Message
		} else {
			message = definitionFor(ErrUnknown).Message
		}
	}
	if detail == "" {
		detail = message
	}
	if registered && def.Kind == KindInternal {
		message = def.Message
	}
	return &Error{
		code:    code,
		message: message,
		detail:  detail,
		cause:   cause,
	}
}

func joinDetail(detail string, cause error) string {
	if cause == nil {
		return detail
	}
	causeDetail := fmt.Sprintf("%+v", cause)
	if appErr, ok := Parse(cause); ok && appErr.Detail() != "" {
		causeDetail = appErr.Detail()
	}
	if detail == "" {
		return causeDetail
	}
	return fmt.Sprintf("%s: %s", detail, causeDetail)
}

// ParseRegistered returns the first registered application error definition in
// err's unwrap chain.
func ParseRegistered(err error) (Definition, bool) {
	appErr, ok := Parse(err)
	if !ok {
		return Definition{}, false
	}
	return Lookup(appErr.Code())
}

// Parse returns the first application Error in err's unwrap chain.
func Parse(err error) (*Error, bool) {
	if appErr, ok := stderrors.AsType[*Error](err); ok {
		return appErr, true
	}
	return nil, false
}

// WrapDatabase wraps database errors.
func WrapDatabase(err error, operation string) error {
	return wrapIfError(err, ErrDatabase, "database %s failed", operation)
}

// WrapRedis wraps Redis errors.
func WrapRedis(err error, operation string) error {
	return wrapIfError(err, ErrRedis, "redis %s failed", operation)
}

// WrapKafka wraps Kafka errors.
func WrapKafka(err error, operation string) error {
	return wrapIfError(err, ErrKafka, "kafka %s failed", operation)
}

// WrapExternal wraps external service errors.
func WrapExternal(err error, service, operation string) error {
	return wrapIfError(err, ErrExternalService, "external service %s %s failed", service, operation)
}

// WrapBadRequest wraps a bad request error.
func WrapBadRequest(err error, message string) error {
	return wrapOrNew(err, ErrBadRequest, message)
}

// WrapUnauthorized wraps an unauthorized error.
func WrapUnauthorized(err error, message string) error {
	return wrapOrNew(err, ErrUnauthorized, message)
}

// WrapForbidden wraps a forbidden error.
func WrapForbidden(err error, message string) error {
	return wrapOrNew(err, ErrForbidden, message)
}

// WrapNotFound wraps a not found error.
func WrapNotFound(err error, message string) error {
	return wrapOrNew(err, ErrNotFound, message)
}

// WrapMethodNotAllowed wraps a method-not-allowed error.
func WrapMethodNotAllowed(err error, message string) error {
	return wrapOrNew(err, ErrMethodNotAllowed, message)
}

// WrapConflict wraps a conflict error.
func WrapConflict(err error, message string) error {
	return wrapOrNew(err, ErrConflict, message)
}

// WrapInternal wraps an internal server error.
func WrapInternal(err error, message string) error {
	return wrapOrNew(err, ErrInternalServer, message)
}

// NewValidation creates a validation error.
func NewValidation(field, message string) error {
	return newError(ErrValidation, message, fmt.Sprintf("validation failed for field %s: %s", field, message), nil)
}

// NewTokenInvalid creates an invalid-token error.
func NewTokenInvalid() error {
	return New(ErrTokenInvalid, "token is invalid")
}

// NewTokenExpired creates an expired-token error.
func NewTokenExpired() error {
	return New(ErrExpired, "token is expired")
}

// NewUnauthorized creates an unauthorized error.
func NewUnauthorized() error {
	return New(ErrUnauthorized, "unauthorized")
}

// NewForbidden creates a forbidden error.
func NewForbidden() error {
	return New(ErrForbidden, "forbidden")
}

// NewPermissionDenied creates a permission-denied error.
func NewPermissionDenied(resource, action string) error {
	return Newf(ErrPermissionDenied, "permission denied for %s on %s", action, resource)
}

// NewNotFound creates a not found error for resource.
func NewNotFound(resource string) error {
	return Newf(ErrNotFound, "%s not found", resource)
}

// NewMethodNotAllowed creates a method-not-allowed error.
func NewMethodNotAllowed(message string) error {
	return New(ErrMethodNotAllowed, message)
}

// NewConflict creates a conflict error.
func NewConflict(message string) error {
	return New(ErrConflict, message)
}

// NewBadRequest creates a bad request error.
func NewBadRequest(message string) error {
	return Newf(ErrBadRequest, "bad request: %s", message)
}
