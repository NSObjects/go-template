/*
 * OpenAPI3 模板渲染函数
 */

package templates

import (
	"fmt"
	"strings"

	"github.com/NSObjects/go-template/muban/modgen/openapi"
)

// RenderOpenAPIBiz 从OpenAPI3生成业务逻辑模板
func (tr *TemplateRenderer) RenderOpenAPIBiz(module *openapi.APIModule, pascal, packagePath string) (string, error) {
	data := TemplateData{
		Pascal:                pascal,
		PackagePath:           packagePath,
		Operations:            module.Operations,
		HasRequestBodyOrQuery: checkHasRequestBodyOrQuery(module.Operations),
	}

	return tr.Render("biz_openapi", data)
}

// RenderOpenAPIService 从OpenAPI3生成服务层模板
func (tr *TemplateRenderer) RenderOpenAPIService(module *openapi.APIModule, pascal, camel, baseRoute, packagePath string) (string, error) {
	data := TemplateData{
		Pascal:                pascal,
		Camel:                 camel,
		Route:                 baseRoute,
		PackagePath:           packagePath,
		Operations:            module.Operations,
		HasPathParams:         checkHasPathParams(module.Operations),
		HasRequestBodyOrQuery: checkHasRequestBodyOrQuery(module.Operations),
	}

	return tr.Render("service_openapi", data)
}

// RenderOpenAPIParam 从OpenAPI3生成参数模板
func (tr *TemplateRenderer) RenderOpenAPIParam(module *openapi.APIModule, pascal, packagePath string) (string, error) {
	data := TemplateData{
		Pascal:                pascal,
		PackagePath:           packagePath,
		Operations:            module.Operations,
		ResponseDataTypes:     deduplicateResponseData(module.Operations),
		HasTimeFields:         checkHasTimeFieldsAll(module.Operations),
		HasPathParams:         checkHasPathParams(module.Operations),
		HasRequestBodyOrQuery: checkHasRequestBodyOrQuery(module.Operations),
	}

	return tr.Render("param_openapi", data)
}

// deduplicateResponseData 去重响应数据类型
func deduplicateResponseData(operations []openapi.APIOperation) []openapi.ResponseData {
	seen := make(map[string]bool)
	var result []openapi.ResponseData

	for _, op := range operations {
		if op.ResponseData != nil {
			key := op.ResponseData.GoType
			if !seen[key] {
				seen[key] = true
				result = append(result, *op.ResponseData)
			}
		}
	}

	return result
}

// RenderModelFromOpenAPI 从OpenAPI3生成数据模型模板
func RenderModelFromOpenAPI(module *openapi.APIModule, pascal, name, packagePath string) string {
	var structs []string

	// 从schemas生成数据模型
	for schemaName, schema := range module.Schemas {
		if schemaName != "operatorResponse" { // 跳过通用响应结构
			// 将schema名称替换为模块名称
			structs = append(structs, generateModelStruct(pascal, schema))
		}
	}

	return fmt.Sprintf(`/*
 * Generated from OpenAPI3 document
 * Module: %s
 */

package model

import (
	"time"
	"gorm.io/gorm"
)

%s
`, name, strings.Join(structs, "\n"))
}

// generateMethodImplementations 生成方法实现
func generateMethodImplementations(module *openapi.APIModule, pascal string) string {
	var implementations []string

	for _, op := range module.Operations {
		methodName := openapi.GenerateMethodName(op)
		requestType := fmt.Sprintf("%s%sRequest", pascal, methodName)
		returnType, returnValue := generateBizReturnType(methodName, pascal)

		// 根据方法类型决定参数
		switch strings.ToLower(methodName) {
		case "getbyid", "delete":
			implementations = append(implementations, fmt.Sprintf(`
func (h *%sHandler) %s(ctx context.Context, id int64) %s {
	// TODO: 实现业务逻辑
	%s
}`, pascal, methodName, returnType, returnValue))
		case "update":
			implementations = append(implementations, fmt.Sprintf(`
func (h *%sHandler) %s(ctx context.Context, id int64, req param.%s) %s {
	// TODO: 实现业务逻辑
	%s
}`, pascal, methodName, requestType, returnType, returnValue))
		default:
			implementations = append(implementations, fmt.Sprintf(`
func (h *%sHandler) %s(ctx context.Context, req param.%s) %s {
	// TODO: 实现业务逻辑
	%s
}`, pascal, methodName, requestType, returnType, returnValue))
		}
	}

	return strings.Join(implementations, "")
}

// generateInterfaceReturnType 生成接口返回类型
func generateInterfaceReturnType(methodName, pascal string) string {
	switch strings.ToLower(methodName) {
	case "list":
		return fmt.Sprintf("([]param.%sResponse, int64, error)", pascal)
	case "create", "update":
		return "error"
	case "getbyid":
		return fmt.Sprintf("(*param.%sResponse, error)", pascal)
	case "delete":
		return "error"
	default:
		return fmt.Sprintf("(*param.%sResponse, error)", pascal)
	}
}

// generateBizReturnType 生成业务逻辑返回类型和返回值
func generateBizReturnType(methodName, pascal string) (string, string) {
	switch strings.ToLower(methodName) {
	case "list":
		return fmt.Sprintf("([]param.%sResponse, int64, error)", pascal), fmt.Sprintf(generateBizImplementation(methodName, pascal), pascal, pascal)
	case "create":
		return "error", generateBizImplementation(methodName, pascal)
	case "update":
		return "error", fmt.Sprintf(generateBizImplementation(methodName, pascal), pascal)
	case "getbyid":
		return "(*param.%sResponse, error)", fmt.Sprintf(generateBizImplementation(methodName, pascal), pascal)
	case "delete":
		return "error", generateBizImplementation(methodName, pascal)
	default:
		return fmt.Sprintf("(*param.%sResponse, error)", pascal), fmt.Sprintf(generateBizImplementation(methodName, pascal), pascal)
	}
}

// generateBizImplementation 生成业务逻辑实现代码
func generateBizImplementation(methodName, pascal string) string {
	switch strings.ToLower(methodName) {
	case "list":
		return fmt.Sprintf(`// TODO: 实现查询逻辑
	// 使用带context的数据库查询 - context包含链路追踪信息
	// db := h.dataManager.MySQLWithContext(ctx)
	// query := h.dataManager.Query.WithContext(ctx)
	// 
	// 示例实现：
	// var users []model.User
	// var total int64
	// 
	// // 构建查询条件
	// query := h.dataManager.Query.User.WithContext(ctx)
	// if req.Name != "" {
	// 	query = query.Where(h.dataManager.Query.User.Name.Like("%%" + req.Name + "%%"))
	// }
	// if req.Email != "" {
	// 	query = query.Where(h.dataManager.Query.User.Account.Eq(req.Email))
	// }
	// 
	// // 分页查询
	// offset := (req.Page - 1) * req.Count
	// err := query.Count(&total).Offset(offset).Limit(req.Count).Find(&users)
	// if err != nil {
	// 	return nil, 0, err
	// }
	// 
	// // 转换为响应格式
	// var responses []param.%sResponse
	// for _, user := range users {
	// 	responses = append(responses, param.%sResponse{
	// 		ID:   user.ID,
	// 		Name: user.Name,
	// 	})
	// }
	// 
	// return responses, total, nil
	return nil, 0, nil`, pascal, pascal)
	case "create":
		return `// TODO: 实现创建逻辑
	// 使用带context的数据库操作 - context包含链路追踪信息
	// db := h.dataManager.MySQLWithContext(ctx)
	// 
	// 示例实现：
	// user := model.User{
	// 	Name:    req.Name,
	// 	Account: req.Account,
	// 	Phone:   req.Phone,
	// 	Status:  req.Status,
	// }
	// 
	// err := h.dataManager.Query.User.WithContext(ctx).Create(&user)
	// if err != nil {
	// 	return err
	// }
	// 
	// return nil
	return nil`
	case "getbyid":
		return fmt.Sprintf(`// TODO: 实现根据ID查询逻辑
	// 使用带context的数据库查询 - context包含链路追踪信息
	// 
	// 示例实现：
	// var user model.User
	// err := h.dataManager.Query.User.WithContext(ctx).Where(h.dataManager.Query.User.ID.Eq(id)).First(&user)
	// if err != nil {
	// 	return nil, err
	// }
	// 
	// return &param.%sResponse{
	// 	ID:   user.ID,
	// 	Name: user.Name,
	// }, nil
	return nil, nil`, pascal)
	case "update":
		return `// TODO: 实现更新逻辑
	// 使用带context的数据库操作 - context包含链路追踪信息
	// 
	// 示例实现：
	// updates := map[string]interface{}{
	// 	"name":   req.Name,
	// 	"phone":  req.Phone,
	// 	"status": req.Status,
	// }
	// 
	// err := h.dataManager.Query.User.WithContext(ctx).Where(h.dataManager.Query.User.ID.Eq(req.Id)).Updates(updates)
	// if err != nil {
	// 	return err
	// }
	// 
	// return nil
	return nil`
	case "delete":
		return `// TODO: 实现删除逻辑
	// 使用带context的数据库操作 - context包含链路追踪信息
	// 
	// 示例实现：
	// err := h.dataManager.Query.User.WithContext(ctx).Where(h.dataManager.Query.User.ID.Eq(id)).Delete(&model.User{})
	// if err != nil {
	// 	return err
	// }
	// 
	// return nil
	return nil`
	default:
		return `// TODO: 实现业务逻辑
	return nil`
	}
}

// generateHandlerImplementations 生成处理器实现
func generateHandlerImplementations(module *openapi.APIModule, pascal string) string {
	var implementations []string

	for _, op := range module.Operations {
		handlerName := openapi.GenerateHandlerName(op)
		methodName := openapi.GenerateMethodName(op)
		requestType := fmt.Sprintf("%s%sRequest", pascal, methodName)
		implementation := generateServiceHandlerImplementation(pascal, handlerName, methodName, requestType)
		implementations = append(implementations, implementation)
	}

	return strings.Join(implementations, "")
}

// generateServiceHandlerImplementation 生成service层处理器实现
func generateServiceHandlerImplementation(pascal, handlerName, methodName, requestType string) string {
	switch strings.ToLower(methodName) {
	case "list":
		return fmt.Sprintf(`
func (c *%sController) %s(ctx echo.Context) error {
	// TODO: 绑定和验证请求参数
	var req param.%s
	if err := BindAndValidate(ctx, &req); err != nil {
		return err
	}
	
	// 调用业务逻辑 - 构造包含链路追踪信息的context
	bizCtx := utils.BuildContext(ctx)
	list, total, err := c.%s.%s(bizCtx, req)
	if err != nil {
		return err
	}
	
	// 返回列表数据
	return resp.ListDataResponse(list, total, ctx)
}`, pascal, handlerName, requestType, strings.ToLower(pascal[:1])+pascal[1:], methodName)
	case "create":
		return fmt.Sprintf(`
func (c *%sController) %s(ctx echo.Context) error {
	// TODO: 绑定和验证请求参数
	var req param.%s
	if err := BindAndValidate(ctx, &req); err != nil {
		return err
	}
	
	// 调用业务逻辑 - 构造包含链路追踪信息的context
	bizCtx := utils.BuildContext(ctx)
	if err := c.%s.%s(bizCtx, req); err != nil {
		return err
	}
	
	// 返回操作成功
	return resp.OperateSuccess(ctx)
}`, pascal, handlerName, requestType, strings.ToLower(pascal[:1])+pascal[1:], methodName)
	case "update":
		return fmt.Sprintf(`
func (c *%sController) %s(ctx echo.Context) error {
	// 获取路径参数
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	
	// 绑定和验证请求体参数
	var req param.%s
	if err := BindAndValidate(ctx, &req); err != nil {
		return err
	}
	
	// 调用业务逻辑 - 构造包含链路追踪信息的context
	bizCtx := utils.BuildContext(ctx)
	if err := c.%s.%s(bizCtx, id, req); err != nil {
		return err
	}
	
	// 返回操作成功
	return resp.OperateSuccess(ctx)
}`, pascal, handlerName, requestType, strings.ToLower(pascal[:1])+pascal[1:], methodName)
	case "getbyid":
		return fmt.Sprintf(`
func (c *%sController) %s(ctx echo.Context) error {
	// 获取路径参数
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	
	// 调用业务逻辑 - 构造包含链路追踪信息的context
	bizCtx := utils.BuildContext(ctx)
	result, err := c.%s.%s(bizCtx, id)
	if err != nil {
		return err
	}
	
	// 返回单个数据
	return resp.OneDataResponse(result, ctx)
}`, pascal, handlerName, strings.ToLower(pascal[:1])+pascal[1:], methodName)
	case "delete":
		return fmt.Sprintf(`
func (c *%sController) %s(ctx echo.Context) error {
	// 获取路径参数
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	
	// 调用业务逻辑 - 构造包含链路追踪信息的context
	bizCtx := utils.BuildContext(ctx)
	err := c.%s.%s(bizCtx, id)
	if err != nil {
		return err
	}
	
	// 返回操作成功
	return resp.OperateSuccess(ctx)
}`, pascal, handlerName, strings.ToLower(pascal[:1])+pascal[1:], methodName)
	default:
		return fmt.Sprintf(`
func (c *%sController) %s(ctx echo.Context) error {
	// TODO: 绑定和验证请求参数
	var req param.%s
	if err := BindAndValidate(ctx, &req); err != nil {
		return err
	}
	
	// 调用业务逻辑 - 构造包含链路追踪信息的context
	bizCtx := utils.BuildContext(ctx)
	result, err := c.%s.%s(bizCtx, req)
	if err != nil {
		return err
	}
	
	// 返回数据
	return resp.OneDataResponse(result, ctx)
}`, pascal, handlerName, requestType, strings.ToLower(pascal[:1])+pascal[1:], methodName)
	}
}

// generateRequestStruct 生成请求结构体
func generateRequestStruct(structName string, requestBody *openapi.RequestBody, module *openapi.APIModule) string {
	var fields []string

	// 从requestBody的content中提取schema
	for _, mediaType := range requestBody.Content {
		if mediaType.Schema != nil {
			// 处理schema引用
			var schema *openapi.Schema
			if mediaType.Schema.Ref != "" {
				schema = openapi.ResolveSchemaRef(mediaType.Schema.Ref, &openapi.OpenAPI3{Components: openapi.Components{Schemas: module.Schemas}})
			} else {
				schema = mediaType.Schema
			}

			if schema != nil {
				fields = generateStructFields(schema, module)
			}
			break
		}
	}

	return fmt.Sprintf(`
// %s 请求结构体
type %s struct {
%s
}`, structName, structName, strings.Join(fields, "\n"))
}

// generateSimpleRequestStruct 生成简单请求结构体（无请求体的操作）
func generateSimpleRequestStruct(structName, methodName string, op openapi.APIOperation, module *openapi.APIModule) string {
	var fields []string

	switch strings.ToLower(methodName) {
	case "list":
		// 列表请求包含分页和查询参数
		fieldSet := make(map[string]bool)
		fields = append(fields, "\tPage  int    `json:\"page\" form:\"page\" query:\"page\" validate:\"min=1\"`")
		fields = append(fields, "\tCount int    `json:\"count\" form:\"count\" query:\"count\" validate:\"min=1,max=100\"`")
		fieldSet["page"] = true
		fieldSet["count"] = true

		// 添加查询参数（去重）
		for _, param := range op.Parameters {
			if param.In == "query" && !fieldSet[param.Name] {
				goType := openapi.GenerateGoTypeFromParam(param)
				validationTag := generateValidationTags(param.Schema, param.Required)
				jsonTag := fmt.Sprintf("`json:\"%s\" form:\"%s\" query:\"%s\"", param.Name, param.Name, param.Name)
				if validationTag != "" {
					jsonTag += " " + validationTag
				}
				jsonTag += "`"
				fields = append(fields, fmt.Sprintf("\t%s %s %s", strings.Title(param.Name), goType, jsonTag))
				fieldSet[param.Name] = true
			}
		}
	case "getbyid":
		// 根据ID获取请求
		fields = append(fields, "\tID int64 `json:\"id\" form:\"id\" param:\"id\" validate:\"required,min=1\"`")
	case "delete":
		// 删除请求
		fields = append(fields, "\tID int64 `json:\"id\" param:\"id\" validate:\"required,min=1\"`")
	default:
		// 其他操作，从参数中提取
		for _, param := range op.Parameters {
			if param.In == "path" {
				goType := openapi.GenerateGoTypeFromParam(param)
				validationTag := generateValidationTags(param.Schema, param.Required)
				jsonTag := fmt.Sprintf("`json:\"%s\" param:\"%s\"", param.Name, param.Name)
				if validationTag != "" {
					jsonTag += " " + validationTag
				}
				jsonTag += "`"
				fields = append(fields, fmt.Sprintf("\t%s %s %s", strings.Title(param.Name), goType, jsonTag))
			}
		}
	}

	return fmt.Sprintf(`
// %s 请求结构体
type %s struct {
%s
}`, structName, structName, strings.Join(fields, "\n"))
}

// generateResponseStruct 生成响应结构体
func generateResponseStruct(pascal string, module *openapi.APIModule) string {
	return fmt.Sprintf(`
// %sResponse 响应结构体
type %sResponse struct {
	// TODO: 根据OpenAPI文档定义响应字段
	ID   int64  %s
	Name string %s
}`, pascal, pascal, "`json:\"id\"`", "`json:\"name\"`")
}

// generateQueryStruct 生成查询结构体
func generateQueryStruct(pascal string, module *openapi.APIModule) string {
	var fields []string
	fieldSet := make(map[string]bool) // 用于去重

	// 从operations中提取查询参数
	for _, op := range module.Operations {
		for _, param := range op.Parameters {
			if param.In == "query" && !fieldSet[param.Name] {
				goType := openapi.GenerateGoTypeFromParam(param)
				jsonTag := fmt.Sprintf("`json:\"%s\" form:\"%s\" query:\"%s\"`",
					param.Name, param.Name, param.Name)
				fields = append(fields, fmt.Sprintf("\t%s %s %s",
					strings.Title(param.Name), goType, jsonTag))
				fieldSet[param.Name] = true
			}
		}
	}

	// 添加分页参数（如果不存在）
	if !fieldSet["page"] {
		fields = append(fields, "\tPage  int `json:\"page\" form:\"page\" query:\"page\"`")
	}
	if !fieldSet["count"] {
		fields = append(fields, "\tCount int `json:\"count\" form:\"count\" query:\"count\"`")
	}

	return fmt.Sprintf(`
// %sParam 查询参数结构体
type %sParam struct {
%s
}`, pascal, pascal, strings.Join(fields, "\n"))
}

// generateModelStruct 生成数据模型结构体
func generateModelStruct(schemaName string, schema *openapi.Schema) string {
	var fields []string
	fieldSet := make(map[string]bool) // 用于去重

	if schema.Properties != nil {
		for fieldName, fieldSchema := range schema.Properties {
			// 跳过id字段，因为我们会添加标准的ID字段
			if strings.ToLower(fieldName) == "id" {
				continue
			}

			goType := openapi.GenerateGoTypeFromSchema(fieldSchema)
			jsonTag := fmt.Sprintf("`json:\"%s\" gorm:\"column:%s\"`",
				openapi.ToSnakeCase(fieldName), openapi.ToSnakeCase(fieldName))
			fields = append(fields, fmt.Sprintf("\t%s %s %s",
				strings.Title(fieldName), goType, jsonTag))
			fieldSet[strings.ToLower(fieldName)] = true
		}
	}

	// 添加通用字段
	fields = append(fields, "\tID        int64          `json:\"id\" gorm:\"primaryKey\"`")
	fields = append(fields, "\tCreatedAt time.Time      `json:\"created_at\" gorm:\"autoCreateTime\"`")
	fields = append(fields, "\tUpdatedAt time.Time      `json:\"updated_at\" gorm:\"autoUpdateTime\"`")
	fields = append(fields, "\tDeletedAt gorm.DeletedAt `json:\"deleted_at\" gorm:\"index\"`")

	return fmt.Sprintf(`
// %s 数据模型
type %s struct {
%s
}

// TableName 指定表名
func (%s) TableName() string {
	return "%s"
}`, schemaName, schemaName, strings.Join(fields, "\n"),
		schemaName, openapi.ToSnakeCase(schemaName))
}

// generateStructFields 生成结构体字段
func generateStructFields(schema *openapi.Schema, module *openapi.APIModule) []string {
	var fields []string

	if schema.Properties != nil {
		for fieldName, fieldSchema := range schema.Properties {
			goType := openapi.GenerateGoTypeFromSchema(fieldSchema)
			jsonTag := fmt.Sprintf("`json:\"%s\"`", openapi.ToSnakeCase(fieldName))

			// 检查字段是否必填
			isRequired := false
			if schema.Required != nil {
				for _, requiredField := range schema.Required {
					if requiredField == fieldName {
						isRequired = true
						break
					}
				}
			}

			validationTag := generateValidationTags(fieldSchema, isRequired)

			tag := jsonTag
			if validationTag != "" {
				tag = fmt.Sprintf("`%s %s`", strings.Trim(jsonTag, "`"), validationTag)
			}

			fields = append(fields, fmt.Sprintf("\t%s %s %s",
				strings.Title(fieldName), goType, tag))
		}
	}

	return fields
}

// generateValidationTags 生成验证标签
func generateValidationTags(schema *openapi.Schema, isRequired bool) string {
	var tags []string

	// 必填字段
	if isRequired {
		tags = append(tags, "required")
	}

	// 字符串类型验证
	if schema.Type == "string" {
		// 最小长度
		if schema.MinLength != nil && *schema.MinLength > 0 {
			tags = append(tags, fmt.Sprintf("min=%d", *schema.MinLength))
		}
		// 最大长度
		if schema.MaxLength != nil && *schema.MaxLength > 0 {
			tags = append(tags, fmt.Sprintf("max=%d", *schema.MaxLength))
		}
		// 邮箱格式
		if schema.Format == "email" {
			tags = append(tags, "email")
		}
		// URL格式
		if schema.Format == "uri" {
			tags = append(tags, "url")
		}
		// 正则表达式 - 使用len验证代替regexp，因为validator需要注册regexp验证器
		if schema.Pattern != "" {
			// 对于手机号等特殊格式，使用len验证
			if strings.Contains(schema.Pattern, "\\d{9}") {
				tags = append(tags, "len=11")
			} else {
				// 其他正则表达式暂时跳过，避免验证器注册问题
				// tags = append(tags, fmt.Sprintf("regexp=%s", schema.Pattern))
			}
		}
	}

	// 数字类型验证
	if schema.Type == "integer" || schema.Type == "number" {
		// 最小值
		if schema.Minimum != nil {
			tags = append(tags, fmt.Sprintf("min=%v", *schema.Minimum))
		}
		// 最大值
		if schema.Maximum != nil {
			tags = append(tags, fmt.Sprintf("max=%v", *schema.Maximum))
		}
	}

	// 数组类型验证
	if schema.Type == "array" {
		// 最小长度
		if schema.MinItems != nil && *schema.MinItems > 0 {
			tags = append(tags, fmt.Sprintf("min=%d", *schema.MinItems))
		}
		// 最大长度
		if schema.MaxItems != nil && *schema.MaxItems > 0 {
			tags = append(tags, fmt.Sprintf("max=%d", *schema.MaxItems))
		}
	}

	if len(tags) == 0 {
		return ""
	}

	return fmt.Sprintf("validate:\"%s\"", strings.Join(tags, ","))
}

// checkHasTimeFields 检查操作中是否包含时间字段（仅检查请求体和查询参数，不检查响应数据）
func checkHasTimeFields(operations []openapi.APIOperation) bool {
	for _, op := range operations {
		// 检查查询参数
		for _, param := range op.QueryParameters {
			if param.GoType == "time.Time" {
				return true
			}
		}

		// 检查请求体字段
		if op.RequestBody != nil {
			for _, field := range op.RequestBody.Fields {
				if field.GoType == "time.Time" {
					return true
				}
			}
		}
	}
	return false
}

// checkHasTimeFieldsAll 检查操作中是否包含时间字段（包括所有字段：请求体、查询参数、响应数据）
func checkHasTimeFieldsAll(operations []openapi.APIOperation) bool {
	for _, op := range operations {
		// 检查查询参数
		for _, param := range op.QueryParameters {
			if param.GoType == "time.Time" {
				return true
			}
		}

		// 检查请求体字段
		if op.RequestBody != nil {
			for _, field := range op.RequestBody.Fields {
				if field.GoType == "time.Time" {
					return true
				}
			}
		}

		// 检查响应数据字段
		if op.ResponseData != nil {
			for _, field := range op.ResponseData.Fields {
				if field.GoType == "time.Time" {
					return true
				}
			}
		}
	}
	return false
}

// checkHasPathParams 检查操作中是否包含路径参数
func checkHasPathParams(operations []openapi.APIOperation) bool {
	for _, op := range operations {
		if op.HasPathParams {
			return true
		}
	}
	return false
}

// checkHasRequestBodyOrQuery 检查操作中是否包含请求体或查询参数
func checkHasRequestBodyOrQuery(operations []openapi.APIOperation) bool {
	for _, op := range operations {
		if op.HasRequestBodyOrQuery {
			return true
		}
	}
	return false
}

// RenderOpenAPIBizTests 从OpenAPI3生成增强的业务逻辑测试模板
func (tr *TemplateRenderer) RenderOpenAPIBizTests(module *openapi.APIModule, pascal, camel, packagePath string) (string, error) {
	data := TemplateData{
		Pascal:                pascal,
		Camel:                 camel,
		PackagePath:           packagePath,
		Operations:            module.Operations,
		HasTimeFields:         checkHasTimeFields(module.Operations),
		HasPathParams:         checkHasPathParams(module.Operations),
		HasRequestBodyOrQuery: checkHasRequestBodyOrQuery(module.Operations),
	}

	return tr.RenderBizTestEnhanced(data)
}

// RenderOpenAPIServiceTests 从OpenAPI3生成增强的服务层测试模板
func (tr *TemplateRenderer) RenderOpenAPIServiceTests(module *openapi.APIModule, pascal, camel, packagePath string) (string, error) {
	data := TemplateData{
		Pascal:                pascal,
		Camel:                 camel,
		PackagePath:           packagePath,
		Operations:            module.Operations,
		HasTimeFields:         checkHasTimeFields(module.Operations),
		HasPathParams:         checkHasPathParams(module.Operations),
		HasRequestBodyOrQuery: checkHasRequestBodyOrQuery(module.Operations),
	}

	return tr.RenderServiceTestEnhanced(data)
}

// RenderOpenAPICode 从OpenAPI3生成错误码模板
func (tr *TemplateRenderer) RenderOpenAPICode(module *openapi.APIModule, pascal, packagePath string) (string, error) {
	var allErrorCodes []openapi.ErrorCode
	for _, op := range module.Operations {
		allErrorCodes = append(allErrorCodes, op.XErrorCodes...)
	}

	data := TemplateData{
		Pascal:      pascal,
		PackagePath: packagePath,
		ErrorCodes:  allErrorCodes,
	}

	return tr.Render("code_openapi", data)
}
