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
	var queryParams []Parameter
	for _, param := range operation.Parameters {
		// 处理引用参数
		if param.Ref != "" {
			refParam := resolveParameterRef(param.Ref, openapi)
			if refParam != nil {
				param = *refParam
			}
		}

		// 添加生成相关字段
		param.GoType = GenerateGoTypeFromParam(param)
		param.FieldName = GenerateFieldName(param.Name)
		param.ValidationRules = GenerateValidationRules(param.Schema)

		// 分离路径参数和查询参数
		if param.In == "query" {
			queryParams = append(queryParams, param)
		}

		op.Parameters = append(op.Parameters, param)
	}

	// 设置查询参数
	op.QueryParameters = queryParams

	// 判断是否有请求体或查询参数
	op.HasRequestBodyOrQuery = len(queryParams) > 0 || operation.RequestBody != nil

	// 确定标签
	if len(operation.Tags) > 0 {
		op.Tag = operation.Tags[0]
	}

	// 生成方法名
	op.MethodName = GenerateMethodName(op)

	// 处理响应数据
	op.ResponseData = parseResponseData(operation.Responses, openapi)

	// 处理请求体字段
	if operation.RequestBody != nil {
		op.RequestBody = parseRequestBodyFields(operation.RequestBody, openapi)
	}

	return op
}

// parseResponseData 解析响应数据
func parseResponseData(responses map[string]Response, openapi *OpenAPI3) *ResponseData {
	// 查找200响应
	if response200, ok := responses["200"]; ok {
		if response200.Content != nil {
			if jsonContent, ok := response200.Content["application/json"]; ok && jsonContent.Schema != nil {
				return parseSchemaToResponseData(jsonContent.Schema, openapi)
			}
		}
	}
	return nil
}

// parseSchemaToResponseData 将Schema转换为ResponseData
func parseSchemaToResponseData(schema *Schema, openapi *OpenAPI3) *ResponseData {
	if schema == nil {
		return nil
	}

	// 处理引用
	if schema.Ref != "" {
		refSchema := ResolveSchemaRef(schema.Ref, openapi)
		if refSchema != nil {
			return parseSchemaToResponseData(refSchema, openapi)
		}
		return nil
	}

	// 处理对象类型，查找 data 字段
	if schema.Type == "object" && schema.Properties != nil {
		// 查找 data 字段
		if dataSchema, ok := schema.Properties["data"]; ok {
			return parseDataSchema(dataSchema, openapi)
		}
	}

	return nil
}

// parseDataSchema 解析 data 字段的 Schema
func parseDataSchema(schema *Schema, openapi *OpenAPI3) *ResponseData {
	if schema == nil {
		return nil
	}

	// 处理引用
	if schema.Ref != "" {
		refSchema := ResolveSchemaRef(schema.Ref, openapi)
		if refSchema != nil {
			return parseDataSchema(refSchema, openapi)
		}
		return nil
	}

	// 处理数组类型（list 结构）
	if schema.Type == "array" && schema.Items != nil {
		itemSchema := schema.Items
		if itemSchema.Ref != "" {
			refSchema := ResolveSchemaRef(itemSchema.Ref, openapi)
			if refSchema != nil {
				itemSchema = refSchema
			}
		}

		if itemSchema != nil && itemSchema.Type == "object" && itemSchema.Properties != nil {
			var fields []Field
			for name, propSchema := range itemSchema.Properties {
				required := false
				for _, req := range itemSchema.Required {
					if req == name {
						required = true
						break
					}
				}
				fields = append(fields, GenerateField(name, propSchema, required))
			}

			return &ResponseData{
				GoType:      "ListItem", // 列表项类型
				Description: itemSchema.Description,
				Fields:      fields,
			}
		}
	}

	// 处理对象类型（单个对象）
	if schema.Type == "object" && schema.Properties != nil {
		var fields []Field
		for name, propSchema := range schema.Properties {
			required := false
			for _, req := range schema.Required {
				if req == name {
					required = true
					break
				}
			}
			fields = append(fields, GenerateField(name, propSchema, required))
		}

		return &ResponseData{
			GoType:      "Data", // 单个数据对象
			Description: schema.Description,
			Fields:      fields,
		}
	}

	return nil
}

// parseSchemaToFields 将Schema转换为字段列表
func parseSchemaToFields(schema *Schema, openapi *OpenAPI3) []Field {
	if schema == nil {
		return nil
	}

	// 处理引用
	if schema.Ref != "" {
		refSchema := ResolveSchemaRef(schema.Ref, openapi)
		if refSchema != nil {
			return parseSchemaToFields(refSchema, openapi)
		}
		return nil
	}

	// 处理对象类型
	if schema.Type == "object" && schema.Properties != nil {
		var fields []Field
		for name, propSchema := range schema.Properties {
			required := false
			for _, req := range schema.Required {
				if req == name {
					required = true
					break
				}
			}
			fields = append(fields, GenerateField(name, propSchema, required))
		}
		return fields
	}

	return nil
}

// parseRequestBodyFields 解析请求体字段
func parseRequestBodyFields(requestBody *RequestBody, openapi *OpenAPI3) *RequestBody {
	if requestBody == nil {
		return nil
	}

	// 处理引用
	if requestBody.Ref != "" {
		refRequestBody := resolveRequestBodyRef(requestBody.Ref, openapi)
		if refRequestBody != nil {
			return parseRequestBodyFields(refRequestBody, openapi)
		}
		return requestBody
	}

	// 处理JSON内容
	if jsonContent, ok := requestBody.Content["application/json"]; ok && jsonContent.Schema != nil {
		// 解析请求体字段
		requestBody.Fields = parseSchemaToFields(jsonContent.Schema, openapi)
		return requestBody
	}

	return requestBody
}

// resolveRequestBodyRef 解析请求体引用
func resolveRequestBodyRef(ref string, openapi *OpenAPI3) *RequestBody {
	// 移除 #/components/requestBodies/ 前缀
	if strings.HasPrefix(ref, "#/components/requestBodies/") {
		bodyName := strings.TrimPrefix(ref, "#/components/requestBodies/")
		if body, exists := openapi.Components.RequestBodies[bodyName]; exists {
			return body
		}
	}
	return nil
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
