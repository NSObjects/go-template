/*
 * OpenAPI3 文档解析和代码生成
 * 支持从OpenAPI3文档生成API模块
 */

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// OpenAPI3 结构定义
type OpenAPI3 struct {
	OpenAPI    string              `yaml:"openapi" json:"openapi"`
	Info       Info                `yaml:"info" json:"info"`
	Servers    []Server            `yaml:"servers,omitempty" json:"servers,omitempty"`
	Paths      map[string]PathItem `yaml:"paths" json:"paths"`
	Components Components          `yaml:"components,omitempty" json:"components,omitempty"`
	Tags       []Tag               `yaml:"tags,omitempty" json:"tags,omitempty"`
}

type Info struct {
	Title       string `yaml:"title" json:"title"`
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
	Version     string `yaml:"version" json:"version"`
}

type Server struct {
	URL         string `yaml:"url" json:"url"`
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
}

type PathItem struct {
	Get    *Operation `yaml:"get,omitempty" json:"get,omitempty"`
	Post   *Operation `yaml:"post,omitempty" json:"post,omitempty"`
	Put    *Operation `yaml:"put,omitempty" json:"put,omitempty"`
	Delete *Operation `yaml:"delete,omitempty" json:"delete,omitempty"`
	Patch  *Operation `yaml:"patch,omitempty" json:"patch,omitempty"`
}

type Operation struct {
	Tags        []string            `yaml:"tags,omitempty" json:"tags,omitempty"`
	Summary     string              `yaml:"summary,omitempty" json:"summary,omitempty"`
	Description string              `yaml:"description,omitempty" json:"description,omitempty"`
	OperationID string              `yaml:"operationId,omitempty" json:"operationId,omitempty"`
	Parameters  []Parameter         `yaml:"parameters,omitempty" json:"parameters,omitempty"`
	RequestBody *RequestBody        `yaml:"requestBody,omitempty" json:"requestBody,omitempty"`
	Responses   map[string]Response `yaml:"responses" json:"responses"`
}

type Parameter struct {
	Name        string  `yaml:"name" json:"name"`
	In          string  `yaml:"in" json:"in"`
	Description string  `yaml:"description,omitempty" json:"description,omitempty"`
	Required    bool    `yaml:"required,omitempty" json:"required,omitempty"`
	Schema      *Schema `yaml:"schema,omitempty" json:"schema,omitempty"`
	Ref         string  `yaml:"$ref,omitempty" json:"$ref,omitempty"`
}

type RequestBody struct {
	Description string               `yaml:"description,omitempty" json:"description,omitempty"`
	Required    bool                 `yaml:"required,omitempty" json:"required,omitempty"`
	Content     map[string]MediaType `yaml:"content" json:"content"`
}

type Response struct {
	Description string               `yaml:"description" json:"description"`
	Content     map[string]MediaType `yaml:"content,omitempty" json:"content,omitempty"`
	Ref         string               `yaml:"$ref,omitempty" json:"$ref,omitempty"`
}

type MediaType struct {
	Schema *Schema `yaml:"schema,omitempty" json:"schema,omitempty"`
}

type Schema struct {
	Type        string             `yaml:"type,omitempty" json:"type,omitempty"`
	Format      string             `yaml:"format,omitempty" json:"format,omitempty"`
	Description string             `yaml:"description,omitempty" json:"description,omitempty"`
	Properties  map[string]*Schema `yaml:"properties,omitempty" json:"properties,omitempty"`
	Items       *Schema            `yaml:"items,omitempty" json:"items,omitempty"`
	Required    []string           `yaml:"required,omitempty" json:"required,omitempty"`
	Ref         string             `yaml:"$ref,omitempty" json:"$ref,omitempty"`
	Default     interface{}        `yaml:"default,omitempty" json:"default,omitempty"`
	// 验证相关字段
	MinLength *int     `yaml:"minLength,omitempty" json:"minLength,omitempty"`
	MaxLength *int     `yaml:"maxLength,omitempty" json:"maxLength,omitempty"`
	Pattern   string   `yaml:"pattern,omitempty" json:"pattern,omitempty"`
	Minimum   *float64 `yaml:"minimum,omitempty" json:"minimum,omitempty"`
	Maximum   *float64 `yaml:"maximum,omitempty" json:"maximum,omitempty"`
	MinItems  *int     `yaml:"minItems,omitempty" json:"minItems,omitempty"`
	MaxItems  *int     `yaml:"maxItems,omitempty" json:"maxItems,omitempty"`
}

type Components struct {
	Schemas       map[string]*Schema      `yaml:"schemas,omitempty" json:"schemas,omitempty"`
	Parameters    map[string]*Parameter   `yaml:"parameters,omitempty" json:"parameters,omitempty"`
	Responses     map[string]*Response    `yaml:"responses,omitempty" json:"responses,omitempty"`
	RequestBodies map[string]*RequestBody `yaml:"requestBodies,omitempty" json:"requestBodies,omitempty"`
}

type Tag struct {
	Name        string `yaml:"name" json:"name"`
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
}

// API模块信息
type APIModule struct {
	Name          string
	Tag           string
	Operations    []APIOperation
	Schemas       map[string]*Schema
	Parameters    map[string]*Parameter
	RequestBodies map[string]*RequestBody
}

type APIOperation struct {
	Method      string
	Path        string
	Summary     string
	Description string
	OperationID string
	Parameters  []Parameter
	RequestBody *RequestBody
	Responses   map[string]Response
	Tag         string
}

// 解析OpenAPI3文档
func parseOpenAPI3(filePath string) (*OpenAPI3, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %v", err)
	}

	var openapi OpenAPI3

	// 根据文件扩展名选择解析方式
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext == ".yaml" || ext == ".yml" {
		err = yaml.Unmarshal(data, &openapi)
	} else if ext == ".json" {
		err = json.Unmarshal(data, &openapi)
	} else {
		return nil, fmt.Errorf("不支持的文件格式: %s", ext)
	}

	if err != nil {
		return nil, fmt.Errorf("解析OpenAPI文档失败: %v", err)
	}

	return &openapi, nil
}

// 从OpenAPI3文档生成API模块
func generateFromOpenAPI(openapi *OpenAPI3, moduleName string) (*APIModule, error) {
	module := &APIModule{
		Name:          moduleName,
		Operations:    []APIOperation{},
		Schemas:       make(map[string]*Schema),
		Parameters:    make(map[string]*Parameter),
		RequestBodies: make(map[string]*RequestBody),
	}

	// 复制schemas
	if openapi.Components.Schemas != nil {
		module.Schemas = openapi.Components.Schemas
	}

	// 复制parameters
	if openapi.Components.Parameters != nil {
		module.Parameters = openapi.Components.Parameters
	}

	// 复制requestBodies
	if openapi.Components.RequestBodies != nil {
		module.RequestBodies = openapi.Components.RequestBodies
	}

	// 解析paths
	for path, pathItem := range openapi.Paths {
		// 处理GET操作
		if pathItem.Get != nil {
			op := parseOperation("GET", path, pathItem.Get, openapi)
			module.Operations = append(module.Operations, op)
		}

		// 处理POST操作
		if pathItem.Post != nil {
			op := parseOperation("POST", path, pathItem.Post, openapi)
			module.Operations = append(module.Operations, op)
		}

		// 处理PUT操作
		if pathItem.Put != nil {
			op := parseOperation("PUT", path, pathItem.Put, openapi)
			module.Operations = append(module.Operations, op)
		}

		// 处理DELETE操作
		if pathItem.Delete != nil {
			op := parseOperation("DELETE", path, pathItem.Delete, openapi)
			module.Operations = append(module.Operations, op)
		}

		// 处理PATCH操作
		if pathItem.Patch != nil {
			op := parseOperation("PATCH", path, pathItem.Patch, openapi)
			module.Operations = append(module.Operations, op)
		}
	}

	// 确定模块标签
	if len(module.Operations) > 0 {
		module.Tag = module.Operations[0].Tag
	}

	return module, nil
}

// 解析单个操作
func parseOperation(method, path string, operation *Operation, openapi *OpenAPI3) APIOperation {
	op := APIOperation{
		Method:      method,
		Path:        path,
		Summary:     operation.Summary,
		Description: operation.Description,
		OperationID: operation.OperationID,
		Parameters:  []Parameter{},
		RequestBody: operation.RequestBody,
		Responses:   operation.Responses,
	}

	// 处理参数
	for _, param := range operation.Parameters {
		// 处理引用参数
		if param.Ref != "" {
			refParam := resolveParameterRef(param.Ref, openapi)
			if refParam != nil {
				op.Parameters = append(op.Parameters, *refParam)
			}
		} else {
			op.Parameters = append(op.Parameters, param)
		}
	}

	// 确定标签
	if len(operation.Tags) > 0 {
		op.Tag = operation.Tags[0]
	}

	return op
}

// 解析参数引用
func resolveParameterRef(ref string, openapi *OpenAPI3) *Parameter {
	// 移除 #/components/parameters/ 前缀
	if strings.HasPrefix(ref, "#/components/parameters/") {
		paramName := strings.TrimPrefix(ref, "#/components/parameters/")
		if param, exists := openapi.Components.Parameters[paramName]; exists {
			return param
		}
	}
	return nil
}

// 解析Schema引用
func resolveSchemaRef(ref string, openapi *OpenAPI3) *Schema {
	// 移除 #/components/schemas/ 前缀
	if strings.HasPrefix(ref, "#/components/schemas/") {
		schemaName := strings.TrimPrefix(ref, "#/components/schemas/")
		if schema, exists := openapi.Components.Schemas[schemaName]; exists {
			return schema
		}
	}
	return nil
}

// 生成Go类型
func generateGoType(schema *Schema, openapi *OpenAPI3) string {
	if schema == nil {
		return "interface{}"
	}

	// 处理引用
	if schema.Ref != "" {
		refSchema := resolveSchemaRef(schema.Ref, openapi)
		if refSchema != nil {
			return generateGoType(refSchema, openapi)
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
			itemType := generateGoType(schema.Items, openapi)
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

// 生成JSON标签
func generateJSONTag(fieldName string) string {
	return fmt.Sprintf("json:\"%s\"", toSnakeCase(fieldName))
}

// 转换为蛇形命名
func toSnakeCase(str string) string {
	var result []rune
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}
