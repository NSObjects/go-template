/*
 * 字符串处理工具函数
 */

package utils

import "strings"

// 颜色定义
const (
	RED    = "\033[0;31m"
	GREEN  = "\033[0;32m"
	YELLOW = "\033[1;33m"
	BLUE   = "\033[0;34m"
	NC     = "\033[0m" // No Color
)

// ToPascal 转换为Pascal命名
func ToPascal(s string) string {
	parts := splitWords(s)
	for i := range parts {
		if parts[i] == "" {
			continue
		}
		parts[i] = strings.ToUpper(parts[i][:1]) + strings.ToLower(parts[i][1:])
	}
	return strings.Join(parts, "")
}

// ToCamel 转换为Camel命名
func ToCamel(s string) string {
	p := ToPascal(s)
	if p == "" {
		return p
	}
	return strings.ToLower(p[:1]) + p[1:]
}

// SplitWords 分割单词
func splitWords(s string) []string {
	s = strings.ReplaceAll(s, "-", "_")
	parts := strings.Split(s, "_")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

// Pluralize 复数化
func Pluralize(s string) string {
	// 简易复数：以 y 结尾改 ies，其它加 s
	if strings.HasSuffix(s, "y") && len(s) > 1 && !isVowel(s[len(s)-2]) {
		return s[:len(s)-1] + "ies"
	}
	if strings.HasSuffix(s, "s") {
		return s + "es"
	}
	return s + "s"
}

// IsVowel 判断是否为元音字母
func isVowel(b byte) bool {
	switch b {
	case 'a', 'e', 'i', 'o', 'u':
		return true
	default:
		return false
	}
}

// ToSnakeCase 转换为蛇形命名
func ToSnakeCase(str string) string {
	var result []rune
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}
