package code

import (
	"errors"
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
	assert.Equal(t, "Validation failed", info.Message)
	assert.Empty(t, info.Details)
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

func TestNewErrorInfoUnknown(t *testing.T) {
	err := errors.New("unexpected")
	info := NewErrorInfo(err)

	assert.Equal(t, InternalError, info.Type)
	assert.Equal(t, CategorySystem, info.Category)
	assert.Equal(t, ErrUnknown, info.Code)
	assert.Equal(t, "Internal server error", info.Message)
	assert.NotEmpty(t, info.Details)
}
