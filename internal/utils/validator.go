/*
 * 验证器工具函数
 * 用于测试环境中的验证器注册和设置
 */

package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// SetupTestValidator 为测试环境设置验证器
func SetupTestValidator(e *echo.Echo) {
	// 创建验证器实例
	v := validator.New()

	// 注册验证器到Echo实例
	e.Validator = &CustomValidator{validator: v}
}

// CustomValidator 自定义验证器结构
type CustomValidator struct {
	validator *validator.Validate
}

// Validate 实现Echo的Validator接口
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// GetValidator 获取验证器实例
func (cv *CustomValidator) GetValidator() *validator.Validate {
	return cv.validator
}
