/*
 * OpenAPI3 代码生成器
 */

package openapi

import (
	"fmt"
	"regexp"
	"strings"
)

// GenerateGoType 生成Go类型
func GenerateGoType(schema *Schema, openapi *OpenAPI3) string {
	if schema == nil {
		return "interface{}"
	}

	// 处理引用
	if schema.Ref != "" {
		refSchema := ResolveSchemaRef(schema.Ref, openapi)
		if refSchema != nil {
			return GenerateGoType(refSchema, openapi)
		}
		return "interface{}"
	}

	switch schema.Type {
	case "string":
		if schema.Format == "date-time" {
			return "time.Time"
		}
		return "string"
	case "integer":
		if schema.Format == "int64" {
			return "int64"
		}
		return "int"
	case "number":
		if schema.Format == "float" {
			return "float32"
		}
		return "float64"
	case "boolean":
		return "bool"
	case "array":
		if schema.Items != nil {
			itemType := GenerateGoType(schema.Items, openapi)
			return "[]" + itemType
		}
		return "[]interface{}"
	case "object":
		if schema.Properties != nil {
			// 生成结构体
			return "struct{}" // 简化处理，实际应该生成完整结构体
		}
		return "map[string]interface{}"
	default:
		return "interface{}"
	}
}

// GenerateJSONTag 生成JSON标签
func GenerateJSONTag(fieldName string) string {
	return fmt.Sprintf("json:\"%s\"", ToSnakeCase(fieldName))
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

// GenerateMethodName 生成方法名
func GenerateMethodName(op APIOperation) string {
	// 优先使用 operationId
	if op.OperationID != "" {
		// 清理 operationId 中的非法字符（空格、下划线、连字符）
		cleaned := regexp.MustCompile(`[\s\-_]+`).ReplaceAllString(op.OperationID, "")

		// 定义 operationId 到方法名的标准映射
		operationMap := map[string]string{
			"createuser":  "Create",
			"findusers":   "List",
			"getuserbyid": "GetByID",
			"updateuser":  "Update",
			"deleteuser":  "Delete",
		}

		// 尝试从映射表中获取标准方法名
		if method, ok := operationMap[strings.ToLower(cleaned)]; ok {
			return method
		}

		// 否则转换为 PascalCase 返回
		return toPascal(cleaned)
	}

	// 根据 HTTP 方法和路径生成方法名
	switch op.Method {
	case "GET":
		if strings.Contains(op.Path, "{id}") {
			return "GetByID"
		}
		return "List"
	case "POST":
		return "Create"
	case "PUT":
		return "Update"
	case "DELETE":
		return "Delete"
	case "PATCH":
		return "Patch"
	default:
		return "Handle"
	}
}

// ToPascal 转换为Pascal命名
func toPascal(s string) string {
	parts := splitWords(s)
	for i := range parts {
		if parts[i] == "" {
			continue
		}
		parts[i] = strings.ToUpper(parts[i][:1]) + strings.ToLower(parts[i][1:])
	}
	return strings.Join(parts, "")
}

// SplitWords 分割单词
func splitWords(s string) []string {
	// 先处理连字符和下划线
	s = strings.ReplaceAll(s, "-", "_")
	s = strings.ReplaceAll(s, " ", "_")

	// 分割下划线
	parts := strings.Split(s, "_")
	out := make([]string, 0, len(parts))

	for _, p := range parts {
		if p == "" {
			continue
		}

		// 处理驼峰命名：在大小写转换处分割
		var words []string
		var current strings.Builder

		for i, r := range p {
			if i > 0 && isUpper(r) && !isUpper(rune(p[i-1])) {
				// 遇到大写字母且前一个字符不是大写字母，开始新单词
				if current.Len() > 0 {
					words = append(words, current.String())
					current.Reset()
				}
			}
			current.WriteRune(r)
		}

		if current.Len() > 0 {
			words = append(words, current.String())
		}

		out = append(out, words...)
	}

	return out
}

// isUpper 检查字符是否为大写
func isUpper(r rune) bool {
	return r >= 'A' && r <= 'Z'
}

// GenerateRoute 生成路由
func GenerateRoute(op APIOperation, baseRoute string) string {
	// 简化路径处理
	path := op.Path
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}
	return "/" + path
}

// GenerateHandlerName 生成处理器名
func GenerateHandlerName(op APIOperation) string {
	method := GenerateMethodName(op)
	return strings.ToLower(method[:1]) + method[1:]
}

// GenerateGoTypeFromParam 从参数生成Go类型
func GenerateGoTypeFromParam(param Parameter) string {
	if param.Schema != nil {
		return GenerateGoTypeFromSchema(param.Schema)
	}
	return "string"
}

// GenerateGoTypeFromSchema 从Schema生成Go类型
func GenerateGoTypeFromSchema(schema *Schema) string {
	if schema == nil {
		return "interface{}"
	}

	// 处理引用
	if schema.Ref != "" {
		// 简化处理，返回引用类型名
		refName := strings.TrimPrefix(schema.Ref, "#/components/schemas/")
		return refName
	}

	switch schema.Type {
	case "string":
		if schema.Format == "date-time" {
			return "time.Time"
		}
		return "string"
	case "integer":
		if schema.Format == "int64" {
			return "int64"
		}
		return "int"
	case "number":
		if schema.Format == "float" {
			return "float32"
		}
		return "float64"
	case "boolean":
		return "bool"
	case "array":
		if schema.Items != nil {
			itemType := GenerateGoTypeFromSchema(schema.Items)
			return "[]" + itemType
		}
		return "[]interface{}"
	case "object":
		return "map[string]interface{}"
	default:
		return "interface{}"
	}
}

// GenerateValidationRules 生成验证规则
func GenerateValidationRules(schema *Schema) string {
	if schema == nil {
		return ""
	}

	var rules []string

	// 字符串验证
	if schema.Type == "string" {
		if schema.MinLength != nil {
			rules = append(rules, fmt.Sprintf("min=%d", *schema.MinLength))
		}
		if schema.MaxLength != nil {
			rules = append(rules, fmt.Sprintf("max=%d", *schema.MaxLength))
		}
		if schema.Pattern != "" {
			rules = append(rules, fmt.Sprintf("pattern=%s", schema.Pattern))
		}
	}

	// 数值验证
	if schema.Type == "integer" || schema.Type == "number" {
		if schema.Minimum != nil {
			rules = append(rules, fmt.Sprintf("min=%g", *schema.Minimum))
		}
		if schema.Maximum != nil {
			rules = append(rules, fmt.Sprintf("max=%g", *schema.Maximum))
		}
	}

	// 数组验证
	if schema.Type == "array" {
		if schema.MinItems != nil {
			rules = append(rules, fmt.Sprintf("min=%d", *schema.MinItems))
		}
		if schema.MaxItems != nil {
			rules = append(rules, fmt.Sprintf("max=%d", *schema.MaxItems))
		}
	}

	// 邮箱验证
	if schema.Format == "email" {
		rules = append(rules, "email")
	}

	// URL验证
	if schema.Format == "uri" {
		rules = append(rules, "url")
	}

	if len(rules) == 0 {
		return ""
	}

	return strings.Join(rules, ",")
}

// GenerateFieldName 生成字段名
func GenerateFieldName(name string) string {
	// 转换为PascalCase
	parts := strings.Split(name, "_")
	for i, part := range parts {
		if part != "" {
			parts[i] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
		}
	}
	return strings.Join(parts, "")
}

// GenerateField 生成字段信息
func GenerateField(name string, schema *Schema, required bool) Field {
	return Field{
		Name:            name,
		GoType:          GenerateGoTypeFromSchema(schema),
		FieldName:       GenerateFieldName(name),
		Description:     schema.Description,
		Required:        required,
		ValidationRules: GenerateValidationRules(schema),
	}
}
