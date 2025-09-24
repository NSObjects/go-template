package validator

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

// CustomValidator 自定义验证器
type CustomValidator struct {
	validator *validator.Validate
}

// NewCustomValidator 创建自定义验证器
func NewCustomValidator() *CustomValidator {
	v := validator.New()

	// 注册自定义验证规则
	v.RegisterValidation("username", validateUsername)
	v.RegisterValidation("password", validatePassword)
	v.RegisterValidation("phone", validatePhone)
	v.RegisterValidation("no_sql_injection", validateNoSQLInjection)

	return &CustomValidator{validator: v}
}

// Validate 验证结构体
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// validateUsername 验证用户名
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()

	// 用户名长度3-20，只能包含字母、数字、下划线
	if len(username) < 3 || len(username) > 20 {
		return false
	}

	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, username)
	return matched
}

// validatePassword 验证密码
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// 密码长度至少8位，包含大小写字母、数字
	if len(password) < 8 {
		return false
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

	return hasUpper && hasLower && hasNumber
}

// validatePhone 验证手机号
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	matched, _ := regexp.MatchString(`^1[3-9]\d{9}$`, phone)
	return matched
}

// validateNoSQLInjection 防止SQL注入
func validateNoSQLInjection(fl validator.FieldLevel) bool {
	value := strings.ToLower(fl.Field().String())

	// 检查常见的SQL注入关键词
	sqlKeywords := []string{
		"'", "\"", ";", "--", "/*", "*/", "xp_", "sp_",
		"union", "select", "insert", "update", "delete",
		"drop", "create", "alter", "exec", "execute",
	}

	for _, keyword := range sqlKeywords {
		if strings.Contains(value, keyword) {
			return false
		}
	}

	return true
}
