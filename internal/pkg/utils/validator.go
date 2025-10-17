/*
 * 验证器工具函数
 * 用于测试环境中的验证器注册和设置
 */

package utils

import (
	"github.com/NSObjects/go-template/internal/pkg/validator"
	"github.com/labstack/echo/v4"
)

// SetupTestValidator 为测试环境设置验证器
func SetupTestValidator(e *echo.Echo) {
	// 使用统一的验证器实现
	customValidator := validator.NewCustomValidator()
	e.Validator = customValidator
}
