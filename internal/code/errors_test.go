package code

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewError(t *testing.T) {
	tests := []struct {
		name     string
		code     int
		message  string
		expected string
	}{
		{
			name:     "success",
			code:     ErrSuccess,
			message:  "operation successful",
			expected: "OK",
		},
		{
			name:     "bad request",
			code:     ErrBadRequest,
			message:  "invalid request",
			expected: "Bad request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewError(tt.code, tt.message)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), tt.expected)
		})
	}
}

func TestNewErrorf(t *testing.T) {
	tests := []struct {
		name     string
		code     int
		format   string
		args     []interface{}
		expected string
	}{
		{
			name:     "formatted error",
			code:     ErrBadRequest,
			format:   "invalid %s: %s",
			args:     []interface{}{"field", "value"},
			expected: "Bad request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewErrorf(tt.code, tt.format, tt.args...)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), tt.expected)
		})
	}
}

func TestWrapError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		code     int
		message  string
		expected string
	}{
		{
			name:     "wrap existing error",
			err:      errors.New("original error"),
			code:     ErrInternalServer,
			message:  "wrapped error",
			expected: "Internal server error",
		},
		{
			name:     "wrap nil error",
			err:      nil,
			code:     ErrInternalServer,
			message:  "wrapped error",
			expected: "Internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := WrapError(tt.err, tt.code, tt.message)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), tt.expected)
		})
	}
}

func TestWrapDatabaseError(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		operation string
		expected  string
	}{
		{
			name:      "database error",
			err:       errors.New("connection failed"),
			operation: "query",
			expected:  "Database error",
		},
		{
			name:      "nil error",
			err:       nil,
			operation: "query",
			expected:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := WrapDatabaseError(tt.err, tt.operation)
			if tt.err == nil {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tt.expected)
			}
		})
	}
}

func TestWrapRedisError(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		operation string
		expected  string
	}{
		{
			name:      "redis error",
			err:       errors.New("connection failed"),
			operation: "get",
			expected:  "Redis error",
		},
		{
			name:      "nil error",
			err:       nil,
			operation: "get",
			expected:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := WrapRedisError(tt.err, tt.operation)
			if tt.err == nil {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tt.expected)
			}
		})
	}
}

func TestNewValidationError(t *testing.T) {
	tests := []struct {
		name     string
		field    string
		message  string
		expected string
	}{
		{
			name:     "validation error",
			field:    "email",
			message:  "invalid format",
			expected: "Validation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewValidationError(tt.field, tt.message)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), tt.expected)
		})
	}
}

func TestNewPermissionDeniedError(t *testing.T) {
	tests := []struct {
		name     string
		resource string
		action   string
		expected string
	}{
		{
			name:     "permission denied",
			resource: "user",
			action:   "delete",
			expected: "Permission denied",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewPermissionDeniedError(tt.resource, tt.action)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), tt.expected)
		})
	}
}

func TestNewNotFoundError(t *testing.T) {
	tests := []struct {
		name     string
		resource string
		expected string
	}{
		{
			name:     "not found",
			resource: "user",
			expected: "Not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewNotFoundError(tt.resource)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), tt.expected)
		})
	}
}
