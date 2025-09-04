/*
 * 验证器工具函数测试
 * 测试验证器相关功能
 */

package utils

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// TestStruct 测试结构体
type TestStruct struct {
	Name  string `validate:"required,min=3,max=10"`
	Email string `validate:"required,email"`
	Age   int    `validate:"min=18,max=100"`
}

func TestSetupTestValidator(t *testing.T) {
	// 创建Echo实例
	e := echo.New()

	// 设置测试验证器
	SetupTestValidator(e)

	// 验证验证器已设置
	assert.NotNil(t, e.Validator)
	assert.IsType(t, &CustomValidator{}, e.Validator)
}

func TestCustomValidator_Validate(t *testing.T) {
	// 创建自定义验证器
	cv := &CustomValidator{validator: validator.New()}

	tests := []struct {
		name      string
		input     interface{}
		wantError bool
	}{
		{
			name: "valid struct",
			input: TestStruct{
				Name:  "John",
				Email: "john@example.com",
				Age:   25,
			},
			wantError: false,
		},
		{
			name: "invalid struct - missing required field",
			input: TestStruct{
				Name:  "",
				Email: "john@example.com",
				Age:   25,
			},
			wantError: true,
		},
		{
			name: "invalid struct - invalid email",
			input: TestStruct{
				Name:  "John",
				Email: "invalid-email",
				Age:   25,
			},
			wantError: true,
		},
		{
			name: "invalid struct - age too young",
			input: TestStruct{
				Name:  "John",
				Email: "john@example.com",
				Age:   15,
			},
			wantError: true,
		},
		{
			name: "invalid struct - name too short",
			input: TestStruct{
				Name:  "Jo",
				Email: "john@example.com",
				Age:   25,
			},
			wantError: true,
		},
		{
			name: "invalid struct - name too long",
			input: TestStruct{
				Name:  "VeryLongName",
				Email: "john@example.com",
				Age:   25,
			},
			wantError: true,
		},
		{
			name: "invalid struct - age too old",
			input: TestStruct{
				Name:  "John",
				Email: "john@example.com",
				Age:   150,
			},
			wantError: true,
		},
		{
			name:      "nil input",
			input:     nil,
			wantError: true,
		},
		{
			name:      "non-struct input",
			input:     "not a struct",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cv.Validate(tt.input)

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCustomValidator_GetValidator(t *testing.T) {
	// 创建自定义验证器
	validator := validator.New()
	cv := &CustomValidator{validator: validator}

	// 获取验证器实例
	result := cv.GetValidator()

	// 验证返回的验证器实例
	assert.Equal(t, validator, result)
	assert.NotNil(t, result)
}

func TestCustomValidator_Integration(t *testing.T) {
	// 创建Echo实例
	e := echo.New()

	// 设置测试验证器
	SetupTestValidator(e)

	// 创建测试处理器
	handler := func(c echo.Context) error {
		var testData TestStruct
		if err := c.Bind(&testData); err != nil {
			return err
		}
		return c.JSON(200, testData)
	}

	// 注册路由
	e.POST("/test", handler)

	// 测试有效数据
	validData := TestStruct{
		Name:  "John",
		Email: "john@example.com",
		Age:   25,
	}

	// 这里我们只测试验证器是否正确设置，不测试HTTP请求
	// 因为实际的HTTP测试需要更复杂的设置
	assert.NotNil(t, e.Validator)

	// 直接测试验证器
	err := e.Validator.Validate(validData)
	assert.NoError(t, err)

	// 测试无效数据
	invalidData := TestStruct{
		Name:  "Jo", // 太短
		Email: "invalid-email",
		Age:   15, // 太年轻
	}

	err = e.Validator.Validate(invalidData)
	assert.Error(t, err)
}

func TestCustomValidator_Struct(t *testing.T) {
	// 测试CustomValidator结构体
	cv := &CustomValidator{validator: validator.New()}

	// 验证字段
	assert.NotNil(t, cv.validator)
	assert.IsType(t, &validator.Validate{}, cv.validator)
}

func TestValidator_FieldValidation(t *testing.T) {
	cv := &CustomValidator{validator: validator.New()}

	// 测试各种字段验证
	tests := []struct {
		name  string
		field string
		value interface{}
		tag   string
		valid bool
	}{
		{
			name:  "required field - valid",
			field: "Name",
			value: "John",
			tag:   "required",
			valid: true,
		},
		{
			name:  "required field - invalid",
			field: "Name",
			value: "",
			tag:   "required",
			valid: false,
		},
		{
			name:  "email field - valid",
			field: "Email",
			value: "john@example.com",
			tag:   "email",
			valid: true,
		},
		{
			name:  "email field - invalid",
			field: "Email",
			value: "invalid-email",
			tag:   "email",
			valid: false,
		},
		{
			name:  "min length - valid",
			field: "Name",
			value: "John",
			tag:   "min=3",
			valid: true,
		},
		{
			name:  "min length - invalid",
			field: "Name",
			value: "Jo",
			tag:   "min=3",
			valid: false,
		},
		{
			name:  "max length - valid",
			field: "Name",
			value: "John",
			tag:   "max=10",
			valid: true,
		},
		{
			name:  "max length - invalid",
			field: "Name",
			value: "VeryLongName",
			tag:   "max=10",
			valid: false,
		},
		{
			name:  "min value - valid",
			field: "Age",
			value: 25,
			tag:   "min=18",
			valid: true,
		},
		{
			name:  "min value - invalid",
			field: "Age",
			value: 15,
			tag:   "min=18",
			valid: false,
		},
		{
			name:  "max value - valid",
			field: "Age",
			value: 25,
			tag:   "max=100",
			valid: true,
		},
		{
			name:  "max value - invalid",
			field: "Age",
			value: 150,
			tag:   "max=100",
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建测试结构体
			testStruct := TestStruct{
				Name:  "John",
				Email: "john@example.com",
				Age:   25,
			}

			// 设置测试字段值
			switch tt.field {
			case "Name":
				testStruct.Name = tt.value.(string)
			case "Email":
				testStruct.Email = tt.value.(string)
			case "Age":
				testStruct.Age = tt.value.(int)
			}

			// 验证
			err := cv.Validate(testStruct)

			if tt.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}


