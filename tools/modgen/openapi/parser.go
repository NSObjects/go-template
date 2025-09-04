/*
 * OpenAPI3 解析器
 */

package openapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// ParseOpenAPI3 解析OpenAPI3文档
func ParseOpenAPI3(filePath string) (*OpenAPI3, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %v", err)
	}

	var openapi OpenAPI3

	// 根据文件扩展名选择解析方式
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".yaml", ".yml":
		err = yaml.Unmarshal(data, &openapi)
	case ".json":
		err = json.Unmarshal(data, &openapi)
	default:
		return nil, fmt.Errorf("不支持的文件格式: %s", ext)
	}

	if err != nil {
		return nil, fmt.Errorf("解析OpenAPI文档失败: %v", err)
	}

	return &openapi, nil
}

// GenerateFromOpenAPI 从OpenAPI3文档生成API模块
func GenerateFromOpenAPI(openapi *OpenAPI3, moduleName string) (*APIModule, error) {
	// 验证模块名称是否在OpenAPI文档中存在
	if !validateModuleExists(openapi, moduleName) {
		return nil, fmt.Errorf("OpenAPI文档中没有找到模块 '%s' 的相关接口。请检查OpenAPI文档中的tags或确保模块名称正确", moduleName)
	}

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

// parseOperation 解析单个操作
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

// resolveParameterRef 解析参数引用
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

// ResolveSchemaRef 解析Schema引用
func ResolveSchemaRef(ref string, openapi *OpenAPI3) *Schema {
	// 移除 #/components/schemas/ 前缀
	if strings.HasPrefix(ref, "#/components/schemas/") {
		schemaName := strings.TrimPrefix(ref, "#/components/schemas/")
		if schema, exists := openapi.Components.Schemas[schemaName]; exists {
			return schema
		}
	}
	return nil
}

// validateModuleExists 验证模块是否在OpenAPI文档中存在
func validateModuleExists(openapi *OpenAPI3, moduleName string) bool {
	// 检查tags中是否有对应的模块名称
	for _, tag := range openapi.Tags {
		if strings.EqualFold(tag.Name, moduleName) {
			return true
		}
	}

	// 检查paths中的操作是否有对应的tag
	for _, pathItem := range openapi.Paths {
		operations := []*Operation{pathItem.Get, pathItem.Post, pathItem.Put, pathItem.Delete, pathItem.Patch}
		for _, op := range operations {
			if op != nil {
				for _, tag := range op.Tags {
					if strings.EqualFold(tag, moduleName) {
						return true
					}
				}
			}
		}
	}

	return false
}
