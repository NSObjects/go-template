/*
 * OpenAPI3 解析器
 */

package openapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/NSObjects/go-template/muban/modgen/utils"
	"gopkg.in/yaml.v2"
)

// ParseOpenAPI3 解析OpenAPI3文档
func ParseOpenAPI3(filePath string) (*OpenAPI3, error) {
	data, format, err := loadOpenAPIData(filePath)
	if err != nil {
		return nil, err
	}

	var openapi OpenAPI3
	if err := unmarshalOpenAPIDoc(data, format, &openapi); err != nil {
		return nil, err
	}

	return &openapi, nil
}

func loadOpenAPIData(source string) ([]byte, string, error) {
	if isRemoteURL(source) {
		client := &http.Client{Timeout: 15 * time.Second}
		req, _ := http.NewRequest("GET", source, nil)
		req.Header.Set("User-Agent", "go-template-modgen/1.0")
		resp, err := client.Do(req)
		if err != nil {
			return nil, "", fmt.Errorf("获取远程OpenAPI文档失败: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, "", fmt.Errorf("获取远程OpenAPI文档失败: HTTP %d", resp.StatusCode)
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, "", fmt.Errorf("读取远程OpenAPI文档失败: %v", err)
		}

		ext := strings.ToLower(filepath.Ext(resp.Request.URL.Path))
		if ext == "" {
			ext = strings.ToLower(filepath.Ext(source))
		}
		if ext == "" {
			contentType := resp.Header.Get("Content-Type")
			switch {
			case strings.Contains(strings.ToLower(contentType), "json"):
				ext = ".json"
			case strings.Contains(strings.ToLower(contentType), "yaml") || strings.Contains(strings.ToLower(contentType), "yml"):
				ext = ".yaml"
			}
		}

		// 仅在识别为标准后缀时保留，否则置空，交由解析函数双路尝试
		if ext != ".json" && ext != ".yaml" && ext != ".yml" {
			ext = ""
		}

		return data, ext, nil
	}

	data, err := os.ReadFile(source)
	if err != nil {
		return nil, "", fmt.Errorf("读取文件失败: %v", err)
	}

	return data, strings.ToLower(filepath.Ext(source)), nil
}

func isRemoteURL(source string) bool {
	parsed, err := url.Parse(source)
	if err != nil {
		return false
	}

	return parsed.Scheme == "http" || parsed.Scheme == "https"
}

func unmarshalOpenAPIDoc(data []byte, format string, openapi *OpenAPI3) error {
	switch format {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, openapi); err != nil {
			return fmt.Errorf("解析OpenAPI文档失败: %v", err)
		}
		return nil
	case ".json":
		if err := json.Unmarshal(data, openapi); err != nil {
			return fmt.Errorf("解析OpenAPI文档失败: %v", err)
		}
		return nil
	default:
		// 无法确定格式时，双路尝试 YAML 和 JSON：哪个不报错就采用
		if err := yaml.Unmarshal(data, openapi); err == nil {
			return nil
		}
		if err := json.Unmarshal(data, openapi); err == nil {
			return nil
		}
		return fmt.Errorf("解析OpenAPI文档失败: 无法识别文档格式，请使用 .yaml/.yml 或 .json 文件")
	}
}

// GenerateFromOpenAPI 从OpenAPI3文档生成API模块
func GenerateFromOpenAPI(openapi *OpenAPI3, moduleName string) (*APIModule, error) {
	// 清理模块名
	cleanModuleName := utils.CleanModuleName(moduleName)

	// 验证模块名称是否在OpenAPI文档中存在，否则尝试基于路径名推断
	if !validateModuleExists(openapi, moduleName) {
		// 允许调用方传入路径段作为 moduleName（如 users）
		// 若仍找不到，继续流程，最终会因为没有匹配的 operations 返回空集合
	}

	module := &APIModule{
		Name:          cleanModuleName,
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

	// 解析paths，只处理属于当前模块的操作；当 tag 为中文或不匹配时，回退到路径首段匹配
	for path, pathItem := range openapi.Paths {
		fallbackModule := extractPathModule(path)
		// 处理GET操作
		if pathItem.Get != nil && (hasModuleTag(pathItem.Get.Tags, moduleName) || strings.EqualFold(fallbackModule, moduleName)) {
			op := parseOperation("GET", path, pathItem.Get, openapi)
			module.Operations = append(module.Operations, op)
		}

		// 处理POST操作
		if pathItem.Post != nil && (hasModuleTag(pathItem.Post.Tags, moduleName) || strings.EqualFold(fallbackModule, moduleName)) {
			op := parseOperation("POST", path, pathItem.Post, openapi)
			module.Operations = append(module.Operations, op)
		}

		// 处理PUT操作
		if pathItem.Put != nil && (hasModuleTag(pathItem.Put.Tags, moduleName) || strings.EqualFold(fallbackModule, moduleName)) {
			op := parseOperation("PUT", path, pathItem.Put, openapi)
			module.Operations = append(module.Operations, op)
		}

		// 处理DELETE操作
		if pathItem.Delete != nil && (hasModuleTag(pathItem.Delete.Tags, moduleName) || strings.EqualFold(fallbackModule, moduleName)) {
			op := parseOperation("DELETE", path, pathItem.Delete, openapi)
			module.Operations = append(module.Operations, op)
		}

		// 处理PATCH操作
		if pathItem.Patch != nil && (hasModuleTag(pathItem.Patch.Tags, moduleName) || strings.EqualFold(fallbackModule, moduleName)) {
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

	// 判断是否有路径参数
	hasPathParams := false
	for _, param := range operation.Parameters {
		if param.In == "path" {
			hasPathParams = true
			break
		}
	}
	op.HasPathParams = hasPathParams

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

	// 处理错误码
	op.XErrorCodes = operation.XErrorCodes

	return op
}

// parseResponseData 解析响应数据
func parseResponseData(responses map[string]Response, openapi *OpenAPI3) *ResponseData {
	// 查找200响应
	if response200, ok := responses["200"]; ok {
		// 处理 $ref 引用
		if response200.Ref != "" {
			// 解析引用
			refResponse := resolveResponseRef(response200.Ref, openapi)
			if refResponse != nil {
				return parseResponseData(map[string]Response{"200": *refResponse}, openapi)
			}
		}

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

	// 处理 OpenAPI 3.1 组合模式
	if len(schema.AllOf) > 0 {
		// 合并 allOf 中的所有 schema
		mergedSchema := mergeAllOfSchemas(schema.AllOf, openapi)
		if mergedSchema != nil {
			return parseSchemaToResponseData(mergedSchema, openapi)
		}
	}

	// 处理 oneOf 组合模式
	if len(schema.OneOf) > 0 {
		// 对于 oneOf，我们选择第一个有效的 schema
		for _, subSchema := range schema.OneOf {
			if result := parseSchemaToResponseData(subSchema, openapi); result != nil {
				return result
			}
		}
	}

	// 处理 anyOf 组合模式
	if len(schema.AnyOf) > 0 {
		// 对于 anyOf，我们选择第一个有效的 schema
		for _, subSchema := range schema.AnyOf {
			if result := parseSchemaToResponseData(subSchema, openapi); result != nil {
				return result
			}
		}
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

// mergeAllOfSchemas 合并 allOf 中的所有 schema
func mergeAllOfSchemas(allOf []*Schema, openapi *OpenAPI3) *Schema {
	if len(allOf) == 0 {
		return nil
	}

	merged := &Schema{
		Type:       "object",
		Properties: make(map[string]*Schema),
		Required:   []string{},
	}

	for _, schema := range allOf {
		if schema == nil {
			continue
		}

		// 处理引用
		resolvedSchema := schema
		if schema.Ref != "" {
			resolvedSchema = ResolveSchemaRef(schema.Ref, openapi)
			if resolvedSchema == nil {
				continue
			}
		}

		// 合并 properties
		if resolvedSchema.Properties != nil {
			for name, propSchema := range resolvedSchema.Properties {
				merged.Properties[name] = propSchema
			}
		}

		// 合并 required 字段
		if resolvedSchema.Required != nil {
			for _, req := range resolvedSchema.Required {
				// 避免重复添加
				found := false
				for _, existing := range merged.Required {
					if existing == req {
						found = true
						break
					}
				}
				if !found {
					merged.Required = append(merged.Required, req)
				}
			}
		}

		// 合并其他属性
		if resolvedSchema.Type != "" && merged.Type == "" {
			merged.Type = resolvedSchema.Type
		}
		if resolvedSchema.Description != "" && merged.Description == "" {
			merged.Description = resolvedSchema.Description
		}
		// 合并 OpenAPI 3.1 字段
		if len(resolvedSchema.Enum) > 0 && len(merged.Enum) == 0 {
			merged.Enum = resolvedSchema.Enum
		}
		if resolvedSchema.Nullable != nil && merged.Nullable == nil {
			merged.Nullable = resolvedSchema.Nullable
		}
		if resolvedSchema.Const != "" && merged.Const == "" {
			merged.Const = resolvedSchema.Const
		}
		if len(resolvedSchema.Examples) > 0 && len(merged.Examples) == 0 {
			merged.Examples = resolvedSchema.Examples
		}
	}

	return merged
}

// parseDataSchema 解析 data 字段的 Schema
func parseDataSchema(schema *Schema, openapi *OpenAPI3) *ResponseData {
	return parseDataSchemaWithDepth(schema, openapi, 0)
}

func parseDataSchemaWithDepth(schema *Schema, openapi *OpenAPI3, depth int) *ResponseData {
	if schema == nil || depth > 100 { // 防止循环引用导致栈溢出
		return nil
	}

	// 处理引用
	if schema.Ref != "" {
		refSchema := ResolveSchemaRef(schema.Ref, openapi)
		if refSchema != nil {
			return parseDataSchemaWithDepth(refSchema, openapi, depth+1)
		}
		return nil
	}

	// 处理数组类型（list 结构）
	if schema.Type == "array" && schema.Items != nil {
		itemSchema := resolveSchemaRefIfNeeded(schema.Items, openapi)
		if itemSchema == nil {
			return nil
		}

		if itemSchema.Type == "object" {
			return handleObjectType(itemSchema, "ListItem", "列表项")
		} else {
			// 处理其他类型，生成默认的 ListItem
			return &ResponseData{
				GoType:      "ListItem", // 列表项类型
				Description: "列表项",
				Fields:      []Field{}, // 空字段列表
			}
		}
	}

	// 处理对象类型（单个对象）
	if schema.Type == "object" || schema.Type == "" {
		// 检查是否有 list 字段（列表响应）
		if listSchema, ok := schema.Properties["list"]; ok {
			return parseDataSchemaWithDepth(listSchema, openapi, depth+1)
		}

		// 检查是否有 total 字段（列表响应结构）
		if _, hasTotal := schema.Properties["total"]; hasTotal {
			// 这是一个列表响应结构，查找 list 字段
			if listSchema, ok := schema.Properties["list"]; ok {
				return parseDataSchemaWithDepth(listSchema, openapi, depth+1)
			}
		}

		// 普通对象类型
		return handleObjectType(schema, "Data", schema.Description)
	}

	return nil
}

// resolveSchemaRefIfNeeded 解析 schema 的 $ref（如果存在）
func resolveSchemaRefIfNeeded(schema *Schema, openapi *OpenAPI3) *Schema {
	if schema == nil {
		return nil
	}
	if schema.Ref != "" {
		refSchema := ResolveSchemaRef(schema.Ref, openapi)
		if refSchema != nil {
			return refSchema
		}
	}
	return schema
}

// handleObjectType 处理 object 类型的 schema
func handleObjectType(schema *Schema, goType, defaultDesc string) *ResponseData {
	requiredMap := make(map[string]bool)
	for _, req := range schema.Required {
		requiredMap[req] = true
	}

	var fields []Field
	if schema.Properties != nil {
		for name, propSchema := range schema.Properties {
			required := requiredMap[name]
			fields = append(fields, GenerateField(name, propSchema, required))
		}
	}

	// 如果对象为空（如 SingleResponse 的 data: {}），生成一个默认的 Data 类型
	if len(fields) == 0 {
		if goType == "Data" {
			fields = []Field{
				{
					Name:      "ID",
					GoType:    "int64",
					FieldName: "ID",
					Required:  false,
				},
			}
		} else {
			fields = []Field{} // 空字段列表
		}
	}

	description := schema.Description
	if description == "" {
		description = defaultDesc
	}

	return &ResponseData{
		GoType:      goType,
		Description: description,
		Fields:      fields,
	}
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

// resolveResponseRef 解析Response引用
func resolveResponseRef(ref string, openapi *OpenAPI3) *Response {
	// 移除 #/components/responses/ 前缀
	if strings.HasPrefix(ref, "#/components/responses/") {
		responseName := strings.TrimPrefix(ref, "#/components/responses/")
		if response, exists := openapi.Components.Responses[responseName]; exists {
			return response
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
		// 也检查清理后的名称是否匹配
		if strings.EqualFold(utils.CleanModuleName(tag.Name), moduleName) {
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
					// 也检查清理后的名称是否匹配
					if strings.EqualFold(utils.CleanModuleName(tag), moduleName) {
						return true
					}
				}
			}
		}
	}

	return false
}

// ExtractAllModuleNames 从OpenAPI文档中提取所有模块名
func ExtractAllModuleNames(openapi *OpenAPI3) ([]string, error) {
	moduleSet := make(map[string]bool)

	// 1) 基于路径首段提取模块名（优先，解决中文 tag 场景）
	for path, pathItem := range openapi.Paths {
		fallback := extractPathModule(path)
		if fallback != "" {
			moduleSet[fallback] = true
		}
		operations := []*Operation{pathItem.Get, pathItem.Post, pathItem.Put, pathItem.Delete, pathItem.Patch}
		for _, op := range operations {
			if op == nil {
				continue
			}
			// 2) 对于英文 tag，亦加入集合；中文/非 ASCII 的 tag 跳过，避免生成 Module前缀文件名
			for _, tag := range op.Tags {
				if tag == "" || hasNonASCII(tag) {
					continue
				}
				moduleSet[strings.ToLower(tag)] = true
			}
		}
	}

	// 3) 顶层 tags（同样只收英文），补充集合
	for _, tag := range openapi.Tags {
		if tag.Name != "" && !hasNonASCII(tag.Name) {
			moduleSet[strings.ToLower(tag.Name)] = true
		}
	}

	// 转换为切片（不再调用 CleanModuleName，以保留 ascii 路径名如 users）
	var moduleNames []string
	for moduleName := range moduleSet {
		moduleNames = append(moduleNames, moduleName)
	}

	if len(moduleNames) == 0 {
		return nil, fmt.Errorf("OpenAPI文档中没有找到任何模块标签")
	}

	return moduleNames, nil
}

// hasModuleTag 检查操作是否属于指定模块
func hasModuleTag(tags []string, moduleName string) bool {
	for _, tag := range tags {
		if strings.EqualFold(tag, moduleName) {
			return true
		}
	}
	return false
}

// extractPathModule 基于路径首段提取模块名，跳过常见前缀（/api, /v1, /v2, /api/v1 等）
func extractPathModule(p string) string {
	if p == "" {
		return ""
	}
	s := p
	if !strings.HasPrefix(s, "/") {
		s = "/" + s
	}
	parts := strings.Split(s, "/")
	var filtered []string
	for _, seg := range parts {
		if seg == "" {
			continue
		}
		// 跳过常见前缀
		lower := strings.ToLower(seg)
		if lower == "api" || strings.HasPrefix(lower, "v") {
			continue
		}
		filtered = append(filtered, lower)
	}
	if len(filtered) == 0 {
		return ""
	}
	return filtered[0]
}

// hasNonASCII 判断是否包含非 ASCII 字符（用于过滤中文等情况）
func hasNonASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > 127 {
			return true
		}
	}
	return false
}
