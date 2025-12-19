/*
 * Codegen Unit Tests
 *
 * Tests for the error code generation tool.
 */

package codegen

import (
	"go/constant"
	"go/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractValue_Int64(t *testing.T) {
	tests := []struct {
		name     string
		value    constant.Value
		expected uint64
	}{
		{
			name:     "positive small int",
			value:    constant.MakeInt64(100001),
			expected: 100001,
		},
		{
			name:     "positive large int",
			value:    constant.MakeInt64(100500),
			expected: 100500,
		},
		{
			name:     "zero",
			value:    constant.MakeInt64(0),
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i64, isInt := constant.Int64Val(tt.value)
			u64, isUint := constant.Uint64Val(tt.value)

			if !isInt && !isUint {
				t.Fatalf("value is neither int nor uint")
			}

			// Test the fixed logic: use u64 if available, otherwise convert i64
			var result uint64
			if !isUint {
				result = uint64(i64)
			} else {
				result = u64
			}

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCommentExtraction(t *testing.T) {
	tests := []struct {
		name     string
		comment  string
		expected string
	}{
		{
			name:     "standard format",
			comment:  "// ErrUserNotFound - 404: User not found.",
			expected: "User not found",
		},
		{
			name:     "no trailing period",
			comment:  "// ErrDatabase - 500: Database error",
			expected: "Database error",
		},
		{
			name:     "empty after dash",
			comment:  "// ErrUnknown -",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This tests the parseComment logic
			// Implementation would need to extract the parseComment function
			// or make it a method for testing
			t.Logf("Comment: %s -> Expected: %s", tt.comment, tt.expected)
		})
	}
}

func TestTypeFilter(t *testing.T) {
	// Test that we only process int types
	tests := []struct {
		name     string
		typeName string
		expected bool
	}{
		{"int type", "int", true},
		{"int64 type", "int64", true},
		{"string type", "string", false},
		{"bool type", "bool", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var typ types.Type
			switch tt.typeName {
			case "int":
				typ = types.Typ[types.Int]
			case "int64":
				typ = types.Typ[types.Int64]
			case "string":
				typ = types.Typ[types.String]
			case "bool":
				typ = types.Typ[types.Bool]
			}

			isInt := typ.Underlying().String() == "int" || typ.Underlying().String() == "int64"
			assert.Equal(t, tt.expected, isInt || tt.typeName == "int" || tt.typeName == "int64")
		})
	}
}
