package code

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewErrorInfoBusiness(t *testing.T) {
	err := NewValidationError("email", "invalid format")
	info := NewErrorInfo(err)

	assert.Equal(t, BusinessError, info.Type)
	assert.Equal(t, CategoryValidation, info.Category)
	assert.Equal(t, ErrValidation, info.Code)
	assert.Equal(t, "invalid format", info.Message)
	assert.Contains(t, info.Details, "validation failed for field email")
}

func TestNewErrorInfoInternal(t *testing.T) {
	err := WrapDatabaseError(errors.New("connection refused"), "query user")
	info := NewErrorInfo(err)

	assert.Equal(t, InternalError, info.Type)
	assert.Equal(t, CategoryDatabase, info.Category)
	assert.Equal(t, ErrDatabase, info.Code)
	assert.Equal(t, "Database error", info.Message)
	assert.True(t, strings.Contains(info.Details, "database query user failed"))
}

func TestNewErrorInfoBusinessDetailIsNotExposedAsMessage(t *testing.T) {
	err := WrapBadRequestError(errors.New("parse id: strconv.Atoi failed"), "invalid user id")
	info := NewErrorInfo(err)

	assert.Equal(t, BusinessError, info.Type)
	assert.Equal(t, CategoryValidation, info.Category)
	assert.Equal(t, ErrBadRequest, info.Code)
	assert.Equal(t, "invalid user id", info.Message)
	assert.Contains(t, info.Details, "parse id")
}

func TestNewErrorInfoUnknown(t *testing.T) {
	err := errors.New("unexpected")
	info := NewErrorInfo(err)

	assert.Equal(t, InternalError, info.Type)
	assert.Equal(t, CategorySystem, info.Category)
	assert.Equal(t, ErrUnknown, info.Code)
	assert.Equal(t, "Internal server error", info.Message)
	assert.NotEmpty(t, info.Details)
}

func TestNewErrorInfoUnregisteredAppErrorFallsBackToUnknown(t *testing.T) {
	err := NewError(999999, "should not leak")
	info := NewErrorInfo(err)

	assert.Equal(t, InternalError, info.Type)
	assert.Equal(t, CategorySystem, info.Category)
	assert.Equal(t, ErrUnknown, info.Code)
	assert.Equal(t, "Internal server error", info.Message)
	assert.Contains(t, info.Details, "should not leak")
}

func TestNewErrorInfoNestedAppErrorKeepsInnerDetail(t *testing.T) {
	inner := WrapBadRequestError(errors.New("raw parse failure"), "invalid user id")
	outer := WrapError(inner, ErrBadRequest, "bad request wrapper")
	info := NewErrorInfo(outer)

	assert.Equal(t, BusinessError, info.Type)
	assert.Equal(t, CategoryValidation, info.Category)
	assert.Equal(t, ErrBadRequest, info.Code)
	assert.Equal(t, "bad request wrapper", info.Message)
	assert.Contains(t, info.Details, "raw parse failure")
}

func TestParseRegisteredCoder(t *testing.T) {
	t.Run("plain error is not treated as registered code", func(t *testing.T) {
		_, ok := ParseRegisteredCoder(errors.New("unexpected"))
		assert.False(t, ok)
	})

	t.Run("wrapped project error keeps registered code", func(t *testing.T) {
		err := fmt.Errorf("outer context: %w", WrapNotFoundError(errors.New("missing row"), "user not found"))

		coder, ok := ParseRegisteredCoder(err)
		if assert.True(t, ok) {
			assert.Equal(t, ErrNotFound, coder.Code())
		}
	})
}

func TestParseError(t *testing.T) {
	t.Run("business error keeps public message and detail", func(t *testing.T) {
		err := WrapBadRequestError(errors.New("raw parse failure"), "invalid user id")

		appErr, ok := ParseError(err)
		if assert.True(t, ok) {
			assert.Equal(t, ErrBadRequest, appErr.Code())
			assert.Equal(t, "invalid user id", appErr.Message())
			assert.Contains(t, appErr.Detail(), "raw parse failure")
		}
	})

	t.Run("fmt wrapped project error is parsed", func(t *testing.T) {
		err := fmt.Errorf("outer: %w", NewNotFoundError("user"))

		appErr, ok := ParseError(err)
		if assert.True(t, ok) {
			assert.Equal(t, ErrNotFound, appErr.Code())
			assert.Equal(t, "user not found", appErr.Message())
		}
	})
}
