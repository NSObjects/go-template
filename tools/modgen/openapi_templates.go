/*
 * OpenAPI3 模板渲染函数
 * 根据OpenAPI3文档生成Go代码模板
 */

package main

import (
	"fmt"
	"strings"
)

// 从OpenAPI3生成业务逻辑模板
func renderBizFromOpenAPI(module *APIModule, pascal, packagePath string) string {
	var methods []string

	// 为每个操作生成方法
	for _, op := range module.Operations {
		methodName := generateMethodName(op)
		requestType := fmt.Sprintf("%s%sRequest", pascal, methodName)
		returnType := generateInterfaceReturnType(methodName, pascal)
		methods = append(methods, fmt.Sprintf(`
	// %s %s
	%s(ctx context.Context, req param.%s) %s`,
			methodName, op.Summary, methodName, requestType, returnType))
	}

	return fmt.Sprintf(`/*
 * Generated from OpenAPI3 document
 * Module: %s
 */

package biz

import (
	"context"
	"%s/internal/api/data"
	"%s/internal/api/service/param"
)

// %sUseCase 业务逻辑接口
type %sUseCase interface {
%s
}

// %sHandler 业务逻辑处理器
type %sHandler struct {
	dataManager *data.DataManager
	// TODO: 注入其他依赖
}

// New%sHandler 创建业务逻辑处理器
func New%sHandler(dataManager *data.DataManager) %sUseCase {
	return &%sHandler{
		dataManager: dataManager,
	}
}

// TODO: 实现业务逻辑方法
%s
`,
		module.Name, packagePath, packagePath, pascal, pascal, strings.Join(methods, ""),
		pascal, pascal, pascal, pascal, pascal, pascal, generateMethodImplementations(module, pascal))
}

// 从OpenAPI3生成服务层模板
func renderServiceFromOpenAPI(module *APIModule, pascal, camel, baseRoute, packagePath string) string {
	var routes []string

	// 为每个操作生成路由
	for _, op := range module.Operations {
		route := generateRoute(op, baseRoute)
		handler := generateHandlerName(op)
		routes = append(routes, fmt.Sprintf(`	g.%s("%s", c.%s).Name = "%s"`,
			strings.ToUpper(op.Method), route, handler, op.Summary))
	}

	return fmt.Sprintf(`/*
 * Generated from OpenAPI3 document
 * Module: %s
 */

package service

import (
	"%s/internal/api/biz"
	"%s/internal/api/service/param"
	"%s/internal/resp"
	"%s/internal/utils"
	"github.com/labstack/echo/v4"
)

type %sController struct {
	%s biz.%sUseCase
}

func New%sController(h biz.%sUseCase) RegisterRouter {
	return &%sController{%s: h}
}

func (c *%sController) RegisterRouter(g *echo.Group, m ...echo.MiddlewareFunc) {
%s
}

// TODO: 实现控制器方法
%s
`,
		module.Name, packagePath, packagePath, packagePath, packagePath,
		pascal, camel, pascal, pascal, pascal, pascal, camel, pascal,
		strings.Join(routes, "\n"), generateHandlerImplementations(module, pascal))
}

// 从OpenAPI3生成参数模板
func renderParamFromOpenAPI(module *APIModule, pascal, packagePath string) string {
	var structs []string
	needsTime := false

	// 生成请求结构体
	for _, op := range module.Operations {
		methodName := generateMethodName(op)
		structName := fmt.Sprintf("%s%sRequest", pascal, methodName)

		if op.RequestBody != nil {
			// 有请求体的操作（如Create, Update）
			structs = append(structs, generateRequestStruct(structName, op.RequestBody, module))
		} else {
			// 没有请求体的操作（如List, GetByID, Delete）
			structs = append(structs, generateSimpleRequestStruct(structName, methodName, op, module))
		}
	}

	// 生成响应结构体
	responseStruct := generateResponseStruct(pascal, module)
	structs = append(structs, responseStruct)

	// 生成查询参数结构体
	queryStruct := generateQueryStruct(pascal, module)
	structs = append(structs, queryStruct)

	// 检查是否需要time包
	for _, structStr := range structs {
		if strings.Contains(structStr, "time.Time") {
			needsTime = true
			break
		}
	}

	imports := ""
	if needsTime {
		imports = `import "time"`
	}

	return fmt.Sprintf(`/*
 * Generated from OpenAPI3 document
 * Module: %s
 */

package param

%s

%s
`, module.Name, imports, strings.Join(structs, "\n"))
}

// 从OpenAPI3生成数据模型模板
func renderModelFromOpenAPI(module *APIModule, pascal, name, packagePath string) string {
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

// 生成方法名
func generateMethodName(op APIOperation) string {
	// 优先使用operationId
	if op.OperationID != "" {
		// 处理特殊字符和空格
		operationID := strings.ReplaceAll(op.OperationID, " ", "")
		operationID = strings.ReplaceAll(operationID, "-", "")
		operationID = strings.ReplaceAll(operationID, "_", "")

		// 特殊处理一些常见的操作名
		switch strings.ToLower(operationID) {
		case "createuser":
			return "Create"
		case "findusers":
			return "List"
		case "getuserbyid":
			return "GetByID"
		case "updateuser":
			return "Update"
		case "deleteuser":
			return "Delete"
		default:
			return toPascal(operationID)
		}
	}

	// 根据HTTP方法和路径生成方法名
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

// 生成路由
func generateRoute(op APIOperation, baseRoute string) string {
	// 简化路径处理
	path := op.Path
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}
	return "/" + path
}

// 生成处理器名
func generateHandlerName(op APIOperation) string {
	method := generateMethodName(op)
	return strings.ToLower(method[:1]) + method[1:]
}

// 生成方法实现
func generateMethodImplementations(module *APIModule, pascal string) string {
	var implementations []string

	for _, op := range module.Operations {
		methodName := generateMethodName(op)
		requestType := fmt.Sprintf("%s%sRequest", pascal, methodName)
		returnType, returnValue := generateBizReturnType(methodName, pascal)
		implementations = append(implementations, fmt.Sprintf(`
func (h *%sHandler) %s(ctx context.Context, req param.%s) %s {
	// TODO: 实现业务逻辑
	%s
}`, pascal, methodName, requestType, returnType, returnValue))
	}

	return strings.Join(implementations, "")
}

// 生成接口返回类型
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

// 生成业务逻辑返回类型和返回值
func generateBizReturnType(methodName, pascal string) (string, string) {
	switch strings.ToLower(methodName) {
	case "list":
		return fmt.Sprintf("([]param.%sResponse, int64, error)", pascal), generateBizImplementation(methodName, pascal)
	case "create", "update":
		return "error", generateBizImplementation(methodName, pascal)
	case "getbyid":
		return fmt.Sprintf("(*param.%sResponse, error)", pascal), generateBizImplementation(methodName, pascal)
	case "delete":
		return "error", generateBizImplementation(methodName, pascal)
	default:
		return fmt.Sprintf("(*param.%sResponse, error)", pascal), generateBizImplementation(methodName, pascal)
	}
}

// 生成业务逻辑实现代码
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
	// 	query = query.Where(h.dataManager.Query.User.Name.Like("%" + req.Name + "%"))
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
		return fmt.Sprintf(`// TODO: 实现创建逻辑
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
	return nil`)
	case "getbyid":
		return fmt.Sprintf(`// TODO: 实现根据ID查询逻辑
	// 使用带context的数据库查询 - context包含链路追踪信息
	// 
	// 示例实现：
	// var user model.User
	// err := h.dataManager.Query.User.WithContext(ctx).Where(h.dataManager.Query.User.ID.Eq(req.ID)).First(&user)
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
		return fmt.Sprintf(`// TODO: 实现更新逻辑
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
	return nil`)
	case "delete":
		return `// TODO: 实现删除逻辑
	// 使用带context的数据库操作 - context包含链路追踪信息
	// 
	// 示例实现：
	// err := h.dataManager.Query.User.WithContext(ctx).Where(h.dataManager.Query.User.ID.Eq(req.ID)).Delete(&model.User{})
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

// 生成处理器实现
func generateHandlerImplementations(module *APIModule, pascal string) string {
	var implementations []string

	for _, op := range module.Operations {
		handlerName := generateHandlerName(op)
		methodName := generateMethodName(op)
		requestType := fmt.Sprintf("%s%sRequest", pascal, methodName)
		implementation := generateServiceHandlerImplementation(pascal, handlerName, methodName, requestType)
		implementations = append(implementations, implementation)
	}

	return strings.Join(implementations, "")
}

// 生成service层处理器实现
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
	case "create", "update":
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
	case "getbyid":
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
	
	// 返回单个数据
	return resp.OneDataResponse(result, ctx)
}`, pascal, handlerName, requestType, strings.ToLower(pascal[:1])+pascal[1:], methodName)
	case "delete":
		return fmt.Sprintf(`
func (c *%sController) %s(ctx echo.Context) error {
	// TODO: 绑定和验证请求参数
	var req param.%s
	if err := BindAndValidate(ctx, &req); err != nil {
		return err
	}
	
	// 调用业务逻辑 - 构造包含链路追踪信息的context
	bizCtx := utils.BuildContext(ctx)
	err := c.%s.%s(bizCtx, req)
	if err != nil {
		return err
	}
	
	// 返回操作成功
	return resp.OperateSuccess(ctx)
}`, pascal, handlerName, requestType, strings.ToLower(pascal[:1])+pascal[1:], methodName)
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

// 生成请求结构体
func generateRequestStruct(structName string, requestBody *RequestBody, module *APIModule) string {
	var fields []string

	// 从requestBody的content中提取schema
	for _, mediaType := range requestBody.Content {
		if mediaType.Schema != nil {
			// 处理schema引用
			var schema *Schema
			if mediaType.Schema.Ref != "" {
				schema = resolveSchemaRef(mediaType.Schema.Ref, &OpenAPI3{Components: Components{Schemas: module.Schemas}})
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

// 生成简单请求结构体（无请求体的操作）
func generateSimpleRequestStruct(structName, methodName string, op APIOperation, module *APIModule) string {
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
				goType := generateGoTypeFromParam(param)
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
				goType := generateGoTypeFromParam(param)
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

// 生成响应结构体
func generateResponseStruct(pascal string, module *APIModule) string {
	return fmt.Sprintf(`
// %sResponse 响应结构体
type %sResponse struct {
	// TODO: 根据OpenAPI文档定义响应字段
	ID   int64  %s
	Name string %s
}`, pascal, pascal, "`json:\"id\"`", "`json:\"name\"`")
}

// 生成查询结构体
func generateQueryStruct(pascal string, module *APIModule) string {
	var fields []string
	fieldSet := make(map[string]bool) // 用于去重

	// 从operations中提取查询参数
	for _, op := range module.Operations {
		for _, param := range op.Parameters {
			if param.In == "query" && !fieldSet[param.Name] {
				goType := generateGoTypeFromParam(param)
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

// 生成数据模型结构体
func generateModelStruct(schemaName string, schema *Schema) string {
	var fields []string
	fieldSet := make(map[string]bool) // 用于去重

	if schema.Properties != nil {
		for fieldName, fieldSchema := range schema.Properties {
			// 跳过id字段，因为我们会添加标准的ID字段
			if strings.ToLower(fieldName) == "id" {
				continue
			}

			goType := generateGoTypeFromSchema(fieldSchema)
			jsonTag := fmt.Sprintf("`json:\"%s\" gorm:\"column:%s\"`",
				toSnakeCase(fieldName), toSnakeCase(fieldName))
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
		schemaName, toSnakeCase(schemaName))
}

// 生成结构体字段
func generateStructFields(schema *Schema, module *APIModule) []string {
	var fields []string

	if schema.Properties != nil {
		for fieldName, fieldSchema := range schema.Properties {
			goType := generateGoTypeFromSchema(fieldSchema)
			jsonTag := fmt.Sprintf("`json:\"%s\"`", toSnakeCase(fieldName))

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

// 生成验证标签
func generateValidationTags(schema *Schema, isRequired bool) string {
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

// 从OpenAPI3生成业务逻辑测试模板
func renderBizTestFromOpenAPI(module *APIModule, pascal, packagePath string) string {
	// 为每个操作生成测试方法
	var testMethods []string

	for _, op := range module.Operations {
		methodName := generateMethodName(op)

		// 生成测试方法
		switch strings.ToLower(methodName) {
		case "list":
			testMethods = append(testMethods, fmt.Sprintf(`
func Test%sHandler_%s(t *testing.T) {
	// 创建handler
	handler := &%sHandler{
		dataManager: &data.DataManager{},
	}
	ctx := context.Background()
	req := param.%s%sRequest{
		Page:  1,
		Count: 10,
	}

	// 测试%s方法
	result, total, err := handler.%s(ctx, req)

	// 由于biz层实现只是返回默认值，这里只测试方法调用不panic
	assert.Nil(t, result)
	assert.Equal(t, int64(0), total)
	assert.NoError(t, err)
}

func Test%sHandler_%s_Validation(t *testing.T) {
	// 测试参数验证
	validator := validator.New()
	
	// 测试有效数据
	validReq := param.%s%sRequest{
		Page:  1,
		Count: 10,
	}
	err := validator.Struct(validReq)
	assert.NoError(t, err, "有效数据应该通过验证")
	
	// 测试无效数据
	invalidReq := param.%s%sRequest{
		Page:  0, // 无效页码
		Count: 10,
	}
	err = validator.Struct(invalidReq)
	assert.Error(t, err, "无效数据应该验证失败")
}`, pascal, methodName, pascal, pascal, methodName, methodName, methodName, pascal, methodName, pascal, methodName, pascal, methodName))
		case "create", "update":
			testMethods = append(testMethods, fmt.Sprintf(`
func Test%sHandler_%s(t *testing.T) {
	// 创建handler
	handler := &%sHandler{
		dataManager: &data.DataManager{},
	}
	ctx := context.Background()
	req := param.%s%sRequest{}

	// 测试%s方法
	err := handler.%s(ctx, req)

	// 由于biz层实现只是返回默认值，这里只测试方法调用不panic
	assert.NoError(t, err)
}

func Test%sHandler_%s_Validation(t *testing.T) {
	// 测试参数验证
	validator := validator.New()
	
	// 测试有效数据
	validReq := param.%s%sRequest{
		Name:     "测试用户",
		Phone:    "13812345678",
		Account:  "test@example.com",
		Password: "123456",
		Status:   1,
		Id:       1,
	}
	err := validator.Struct(validReq)
	assert.NoError(t, err, "有效数据应该通过验证")
	
	// 测试无效数据
	invalidReq := param.%s%sRequest{
		// 空结构体，所有必填字段都缺失
	}
	err = validator.Struct(invalidReq)
	assert.Error(t, err, "无效数据应该验证失败")
}`, pascal, methodName, pascal, pascal, methodName, methodName, methodName, pascal, methodName, pascal, methodName, pascal, methodName))
		case "getbyid":
			testMethods = append(testMethods, fmt.Sprintf(`
func Test%sHandler_%s(t *testing.T) {
	// 创建handler
	handler := &%sHandler{
		dataManager: &data.DataManager{},
	}
	ctx := context.Background()
	req := param.%s%sRequest{}

	// 测试%s方法
	result, err := handler.%s(ctx, req)

	// 由于biz层实现只是返回默认值，这里只测试方法调用不panic
	assert.Nil(t, result)
	assert.NoError(t, err)
}

func Test%sHandler_%s_Validation(t *testing.T) {
	// 测试参数验证
	validator := validator.New()
	
	// 测试有效数据
	validReq := param.%s%sRequest{
		ID: 123, // 有效ID
	}
	err := validator.Struct(validReq)
	assert.NoError(t, err, "有效数据应该通过验证")
	
	// 测试无效数据
	invalidReq := param.%s%sRequest{
		ID: 0, // 无效ID
	}
	err = validator.Struct(invalidReq)
	assert.Error(t, err, "无效数据应该验证失败")
}`, pascal, methodName, pascal, pascal, methodName, methodName, methodName, pascal, methodName, pascal, methodName, pascal, methodName))
		case "delete":
			testMethods = append(testMethods, fmt.Sprintf(`
func Test%sHandler_%s(t *testing.T) {
	// 创建handler
	handler := &%sHandler{
		dataManager: &data.DataManager{},
	}
	ctx := context.Background()
	req := param.%s%sRequest{}

	// 测试%s方法
	err := handler.%s(ctx, req)

	// 由于biz层实现只是返回默认值，这里只测试方法调用不panic
	assert.NoError(t, err)
}

func Test%sHandler_%s_Validation(t *testing.T) {
	// 测试参数验证
	validator := validator.New()
	
	// 测试有效数据
	validReq := param.%s%sRequest{
		ID: 123, // 有效ID
	}
	err := validator.Struct(validReq)
	assert.NoError(t, err, "有效数据应该通过验证")
	
	// 测试无效数据
	invalidReq := param.%s%sRequest{
		ID: 0, // 无效ID
	}
	err = validator.Struct(invalidReq)
	assert.Error(t, err, "无效数据应该验证失败")
}`, pascal, methodName, pascal, pascal, methodName, methodName, methodName, pascal, methodName, pascal, methodName, pascal, methodName))
		}
	}

	return fmt.Sprintf(`/*
 * Generated test cases from OpenAPI3 document
 * Module: %s
 */

package biz

import (
	"context"
	"testing"

	"%s/internal/api/data"
	"%s/internal/api/service/param"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)
%s
`, module.Name, packagePath, packagePath, strings.Join(testMethods, ""))
}

// 从OpenAPI3生成服务层测试模板
func renderServiceTestFromOpenAPI(module *APIModule, pascal, packagePath string) string {
	// 为每个操作生成测试方法
	var testMethods []string
	var mockMethods []string

	for _, op := range module.Operations {
		methodName := generateMethodName(op)

		// 生成Mock方法
		switch strings.ToLower(methodName) {
		case "list":
			mockMethods = append(mockMethods, fmt.Sprintf(`
func (m *Mock%sUseCase) %s(ctx context.Context, req param.%s%sRequest) ([]param.%sResponse, int64, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]param.%sResponse), args.Get(1).(int64), args.Error(2)
}`, pascal, methodName, pascal, methodName, pascal, pascal))
		case "create", "update":
			mockMethods = append(mockMethods, fmt.Sprintf(`
func (m *Mock%sUseCase) %s(ctx context.Context, req param.%s%sRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}`, pascal, methodName, pascal, methodName))
		case "getbyid":
			mockMethods = append(mockMethods, fmt.Sprintf(`
func (m *Mock%sUseCase) %s(ctx context.Context, req param.%s%sRequest) (*param.%sResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*param.%sResponse), args.Error(1)
}`, pascal, methodName, pascal, methodName, pascal, pascal))
		case "delete":
			mockMethods = append(mockMethods, fmt.Sprintf(`
func (m *Mock%sUseCase) %s(ctx context.Context, req param.%s%sRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}`, pascal, methodName, pascal, methodName))
		}

		// 生成测试方法
		switch strings.ToLower(methodName) {
		case "list":
			testMethods = append(testMethods, fmt.Sprintf(`
func Test%sController_%s(t *testing.T) {
	// 创建controller并注入mock依赖
	controller := &%sController{
		%s: nil, // 简化测试，不依赖mock
	}

	// 验证controller创建成功
	assert.NotNil(t, controller)
}`, pascal, methodName, pascal, strings.ToLower(pascal)))
		case "create", "update":
			testMethods = append(testMethods, fmt.Sprintf(`
func Test%sController_%s(t *testing.T) {
	// 创建controller并注入mock依赖
	controller := &%sController{
		%s: nil, // 简化测试，不依赖mock
	}

	// 验证controller创建成功
	assert.NotNil(t, controller)
}

func Test%sController_%s_Validation(t *testing.T) {
	// 测试参数验证
	validator := validator.New()
	
	// 测试有效数据
	validReq := param.%s%sRequest{
		Name:     "测试用户",
		Phone:    "13812345678",
		Account:  "test@example.com",
		Password: "123456",
		Status:   1,
		Id:       1,
	}
	err := validator.Struct(validReq)
	assert.NoError(t, err, "有效数据应该通过验证")
	
	// 测试无效数据
	invalidReq := param.%s%sRequest{
		// 空结构体，所有必填字段都缺失
	}
	err = validator.Struct(invalidReq)
	assert.Error(t, err, "无效数据应该验证失败")
}`, pascal, methodName, pascal, strings.ToLower(pascal), pascal, methodName, pascal, methodName, pascal, methodName))
		case "getbyid", "delete":
			testMethods = append(testMethods, fmt.Sprintf(`
func Test%sController_%s(t *testing.T) {
	// 创建controller并注入mock依赖
	controller := &%sController{
		%s: nil, // 简化测试，不依赖mock
	}

	// 验证controller创建成功
	assert.NotNil(t, controller)
}

func Test%sController_%s_Validation(t *testing.T) {
	// 测试参数验证
	validator := validator.New()
	
	// 测试有效数据
	validReq := param.%s%sRequest{
		ID: 123, // 有效ID
	}
	err := validator.Struct(validReq)
	assert.NoError(t, err, "有效数据应该通过验证")
	
	// 测试无效数据
	invalidReq := param.%s%sRequest{
		ID: 0, // 无效ID
	}
	err = validator.Struct(invalidReq)
	assert.Error(t, err, "无效数据应该验证失败")
}`, pascal, methodName, pascal, strings.ToLower(pascal), pascal, methodName, pascal, methodName, pascal, methodName))
		}
	}

	return fmt.Sprintf(`/*
 * Generated test cases from OpenAPI3 document
 * Module: %s
 */

package service

import (
	"context"
	"testing"

	"%s/internal/api/service/param"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock%sUseCase 模拟业务逻辑接口
type Mock%sUseCase struct {
	mock.Mock
}
%s
%s
`, module.Name, packagePath, pascal, pascal, strings.Join(mockMethods, ""), strings.Join(testMethods, ""))
}

// 从参数生成Go类型
func generateGoTypeFromParam(param Parameter) string {
	if param.Schema != nil {
		return generateGoTypeFromSchema(param.Schema)
	}
	return "string"
}

// 从Schema生成Go类型
func generateGoTypeFromSchema(schema *Schema) string {
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
			itemType := generateGoTypeFromSchema(schema.Items)
			return "[]" + itemType
		}
		return "[]interface{}"
	case "object":
		return "map[string]interface{}"
	default:
		return "interface{}"
	}
}
