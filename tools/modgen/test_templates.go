/*
 * 测试用例模板生成器
 * 基于OpenAPI3文档生成测试用例
 * 支持mock框架和validator验证
 */

package main

import (
	"fmt"
	"strings"
)

// 生成业务逻辑测试用例
func renderBizTestFromOpenAPI(module *APIModule, pascal, packagePath string) string {
	var testCases []string

	// 为每个操作生成测试用例
	for _, op := range module.Operations {
		methodName := generateMethodName(op)
		requestType := fmt.Sprintf("%s%sRequest", pascal, methodName)
		testCases = append(testCases, generateBizTestCases(methodName, requestType, pascal, op))
	}

	// 生成结构体测试
	structTests := generateStructTests(pascal, module)

	// 生成mock接口
	mockInterface := generateMockInterface(pascal, module)

	return fmt.Sprintf(`/*
 * Generated test cases from OpenAPI3 document
 * Module: %s
 */

package biz

import (
	"context"
	"fmt"
	"testing"
	"time"

	"%s/internal/api/service/param"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock%sUseCase 模拟业务逻辑接口
type Mock%sUseCase struct {
	mock.Mock
}

%s

%s

%s
`, module.Name, packagePath, pascal, pascal, mockInterface, strings.Join(testCases, "\n"), structTests)
}

// 生成服务层测试用例
func renderServiceTestFromOpenAPI(module *APIModule, pascal, packagePath string) string {
	var testCases []string

	// 为每个操作生成测试用例
	for _, op := range module.Operations {
		handlerName := generateHandlerName(op)
		methodName := generateMethodName(op)
		requestType := fmt.Sprintf("%s%sRequest", pascal, methodName)
		testCases = append(testCases, generateServiceTestCases(handlerName, requestType, pascal, op))
	}

	// 生成HTTP测试用例
	httpTests := generateHTTPTests(pascal, module)

	// 生成validator测试
	validatorTests := generateValidatorTests(pascal, module)

	return fmt.Sprintf(`/*
 * Generated test cases from OpenAPI3 document
 * Module: %s
 */

package service

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"%s/internal/api/service/param"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// CustomValidator 自定义验证器
type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// Mock%sUseCase 模拟业务逻辑接口
type Mock%sUseCase struct {
	mock.Mock
}

%s

%s

%s

%s
`, module.Name, packagePath, pascal, pascal, generateMockInterface(pascal, module), strings.Join(testCases, "\n"), httpTests, validatorTests)
}

// 生成Mock接口
func generateMockInterface(pascal string, module *APIModule) string {
	var methods []string

	for _, op := range module.Operations {
		methodName := generateMethodName(op)
		requestType := fmt.Sprintf("%s%sRequest", pascal, methodName)
		responseType := fmt.Sprintf("%sResponse", pascal)

		// 根据操作类型确定返回类型
		var returnType string
		switch strings.ToLower(methodName) {
		case "list":
			returnType = fmt.Sprintf("([]param.%s, int64, error)", responseType)
		case "create", "update", "getbyid":
			returnType = fmt.Sprintf("(*param.%s, error)", responseType)
		case "delete":
			returnType = "error"
		default:
			returnType = fmt.Sprintf("(*param.%s, error)", responseType)
		}

		// 根据返回类型生成不同的mock方法
		var mockMethod string
		switch strings.ToLower(methodName) {
		case "list":
			mockMethod = fmt.Sprintf(`
func (m *Mock%sUseCase) %s(ctx context.Context, req param.%s) %s {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]param.%sResponse), args.Get(1).(int64), args.Error(2)
}`, pascal, methodName, requestType, returnType, pascal)
		case "create", "update", "getbyid":
			mockMethod = fmt.Sprintf(`
func (m *Mock%sUseCase) %s(ctx context.Context, req param.%s) %s {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*param.%sResponse), args.Error(1)
}`, pascal, methodName, requestType, returnType, pascal)
		case "delete":
			mockMethod = fmt.Sprintf(`
func (m *Mock%sUseCase) %s(ctx context.Context, req param.%s) %s {
	args := m.Called(ctx, req)
	return args.Error(0)
}`, pascal, methodName, requestType, returnType)
		default:
			mockMethod = fmt.Sprintf(`
func (m *Mock%sUseCase) %s(ctx context.Context, req param.%s) %s {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*param.%sResponse), args.Error(1)
}`, pascal, methodName, requestType, returnType, pascal)
		}

		methods = append(methods, mockMethod)
	}

	return strings.Join(methods, "\n")
}

// 生成业务逻辑测试用例
func generateBizTestCases(methodName, requestType, pascal string, op APIOperation) string {
	// 根据OpenAPI文档生成测试数据
	testData := generateTestDataFromOpenAPI(op, pascal)

	return fmt.Sprintf(`
func Test%sHandler_%s(t *testing.T) {
	tests := []struct {
		name    string
		req     param.%s
		wantErr bool
		mockSetup func(*Mock%sUseCase)
	}{
		{
			name: "valid request",
			req: param.%s{
%s
			},
			wantErr: false,
			mockSetup: func(m *Mock%sUseCase) {
				// TODO: 设置mock期望
				m.On("%s", mock.Anything, mock.MatchedBy(func(req interface{}) bool {
					return true // 接受任何有效的请求参数
				})).Return(nil, nil)
			},
		},
		{
			name: "invalid request",
			req: param.%s{
				// 无效数据
			},
			wantErr: true,
			mockSetup: func(m *Mock%sUseCase) {
				// TODO: 设置mock期望
				m.On("%s", mock.Anything, mock.MatchedBy(func(req interface{}) bool {
					return true // 接受任何有效的请求参数
				})).Return(nil, fmt.Errorf("validation error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建mock
			mockUseCase := new(Mock%sUseCase)
			tt.mockSetup(mockUseCase)
			
			handler := &%sHandler{
				// TODO: 注入mock依赖
			}
			ctx := context.Background()
			
			result, err := handler.%s(ctx, tt.req)
			
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				// TODO: 验证返回结果
			}
			
			// 验证mock调用
			mockUseCase.AssertExpectations(t)
		})
	}
}`, pascal, methodName, requestType, pascal, requestType, testData, pascal, methodName, requestType, pascal, methodName, pascal, pascal, methodName)
}

// 生成服务层测试用例
func generateServiceTestCases(handlerName, requestType, pascal string, op APIOperation) string {
	// 根据OpenAPI文档生成测试数据
	testData := generateTestDataFromOpenAPI(op, pascal)
	methodName := generateMethodName(op)

	// 获取正确的HTTP方法
	httpMethod := getHTTPMethod(op.Method)

	// 根据操作类型生成不同的测试
	switch strings.ToLower(methodName) {
	case "list":
		return generateListTestCases(pascal, handlerName, methodName, httpMethod)
	case "getbyid", "delete":
		return generatePathParamTestCases(pascal, handlerName, methodName, httpMethod)
	case "create", "update":
		return generateRequestBodyTestCases(pascal, handlerName, methodName, httpMethod, testData)
	default:
		return generateRequestBodyTestCases(pascal, handlerName, methodName, httpMethod, testData)
	}
}

// 生成List操作的测试用例
func generateListTestCases(pascal, handlerName, methodName, httpMethod string) string {
	return fmt.Sprintf(`
func Test%sController_%s(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		expectedError  bool
		mockSetup      func(*Mock%sUseCase)
	}{
		{
			name: "valid request",
			queryParams: "?page=1&count=10&name=test&email=test@example.com",
			expectedStatus: http.StatusOK,
			expectedError:  false,
			mockSetup: func(m *Mock%sUseCase) {
				// TODO: 设置mock期望
				m.On("%s", mock.Anything, mock.MatchedBy(func(req interface{}) bool {
					return true // 接受任何有效的请求参数
				})).Return([]param.%sResponse{}, int64(0), nil)
			},
		},
		{
			name: "invalid request",
			queryParams: "?page=0&count=0", // 无效的分页参数
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			mockSetup: func(m *Mock%sUseCase) {
				// TODO: 设置mock期望
				m.On("%s", mock.Anything, mock.MatchedBy(func(req interface{}) bool {
					return true // 接受任何有效的请求参数
				})).Return([]param.%sResponse{}, int64(0), fmt.Errorf("validation error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建Echo实例
			e := echo.New()
			
			// 注册validator
			e.Validator = &CustomValidator{validator: validator.New()}
			
			// 创建请求 - List操作使用查询参数
			req := httptest.NewRequest(http.MethodGet, "/test"+tt.queryParams, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// 创建mock
			mockUseCase := new(Mock%sUseCase)
			tt.mockSetup(mockUseCase)

			// 创建控制器
			controller := &%sController{
				%s: mockUseCase,
			}

			// 执行测试
			err := controller.%s(c)

			// 验证结果
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}
			
			// 验证mock调用
			mockUseCase.AssertExpectations(t)
		})
	}
}`, pascal, handlerName, pascal, pascal, methodName, pascal, pascal, methodName, pascal, pascal, pascal, strings.ToLower(pascal), handlerName)
}

// 生成路径参数测试用例 (GetByID, Delete)
func generatePathParamTestCases(pascal, handlerName, methodName, httpMethod string) string {
	var actualMethod string
	switch strings.ToLower(methodName) {
	case "getbyid":
		actualMethod = "http.MethodGet"
	case "delete":
		actualMethod = "http.MethodDelete"
	default:
		actualMethod = "http.MethodGet"
	}
	return fmt.Sprintf(`
func Test%sController_%s(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedError  bool
		mockSetup      func(*Mock%sUseCase)
	}{
		{
			name: "valid request",
			path: "/test/1",
			expectedStatus: http.StatusOK,
			expectedError:  false,
			mockSetup: func(m *Mock%sUseCase) {
				// TODO: 设置mock期望
				m.On("%s", mock.Anything, mock.MatchedBy(func(req interface{}) bool {
					return true // 接受任何有效的请求参数
				})).Return(nil, nil)
			},
		},
		{
			name: "invalid request",
			path: "/test/0", // 无效的ID
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			mockSetup: func(m *Mock%sUseCase) {
				// TODO: 设置mock期望
				m.On("%s", mock.Anything, mock.MatchedBy(func(req interface{}) bool {
					return true // 接受任何有效的请求参数
				})).Return(nil, fmt.Errorf("validation error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建Echo实例
			e := echo.New()
			
			// 注册validator
			e.Validator = &CustomValidator{validator: validator.New()}
			
			// 创建请求 - 使用路径参数
			req := httptest.NewRequest(%s, tt.path, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// 创建mock
			mockUseCase := new(Mock%sUseCase)
			tt.mockSetup(mockUseCase)

			// 创建控制器
			controller := &%sController{
				%s: mockUseCase,
			}

			// 执行测试
			err := controller.%s(c)

			// 验证结果
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}
			
			// 验证mock调用
			mockUseCase.AssertExpectations(t)
		})
	}
}`, pascal, handlerName, pascal, pascal, methodName, pascal, methodName, actualMethod, pascal, pascal, strings.ToLower(pascal), handlerName)
}

// 生成请求体测试用例 (Create, Update)
func generateRequestBodyTestCases(pascal, handlerName, methodName, httpMethod, testData string) string {
	var actualMethod string
	switch strings.ToLower(methodName) {
	case "create":
		actualMethod = "http.MethodPost"
	case "update":
		actualMethod = "http.MethodPut"
	default:
		actualMethod = "http.MethodPost"
	}
	return fmt.Sprintf(`
func Test%sController_%s(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedError  bool
		mockSetup      func(*Mock%sUseCase)
	}{
		{
			name: "valid request",
			requestBody: `+"`"+`%s`+"`"+`,
			expectedStatus: http.StatusOK,
			expectedError:  false,
			mockSetup: func(m *Mock%sUseCase) {
				// TODO: 设置mock期望
				m.On("%s", mock.Anything, mock.MatchedBy(func(req interface{}) bool {
					return true // 接受任何有效的请求参数
				})).Return(nil, nil)
			},
		},
		{
			name: "invalid request",
			requestBody: `+"`"+`{"invalid": "data"}`+"`"+`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			mockSetup: func(m *Mock%sUseCase) {
				// TODO: 设置mock期望
				m.On("%s", mock.Anything, mock.MatchedBy(func(req interface{}) bool {
					return true // 接受任何有效的请求参数
				})).Return(nil, fmt.Errorf("validation error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建Echo实例
			e := echo.New()
			
			// 注册validator
			e.Validator = &CustomValidator{validator: validator.New()}
			
			// 创建请求 - 使用请求体
			req := httptest.NewRequest(%s, "/test", bytes.NewReader([]byte(tt.requestBody)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// 创建mock
			mockUseCase := new(Mock%sUseCase)
			tt.mockSetup(mockUseCase)

			// 创建控制器
			controller := &%sController{
				%s: mockUseCase,
			}

			// 执行测试
			err := controller.%s(c)

			// 验证结果
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}
			
			// 验证mock调用
			mockUseCase.AssertExpectations(t)
		})
	}
}`, pascal, handlerName, pascal, testData, pascal, methodName, pascal, methodName, actualMethod, pascal, pascal, strings.ToLower(pascal), handlerName)
}

// 生成结构体测试
func generateStructTests(pascal string, module *APIModule) string {
	var tests []string

	// 生成请求结构体测试
	for _, op := range module.Operations {
		if op.RequestBody != nil {
			structName := fmt.Sprintf("%s%sRequest", pascal, generateMethodName(op))
			tests = append(tests, generateStructTest(structName, op.RequestBody, module))
		}
	}

	// 生成响应结构体测试
	responseTest := generateResponseStructTest(pascal, module)
	tests = append(tests, responseTest)

	// 生成查询参数结构体测试
	queryTest := generateQueryStructTest(pascal, module)
	tests = append(tests, queryTest)

	return strings.Join(tests, "\n")
}

// 生成结构体测试
func generateStructTest(structName string, requestBody *RequestBody, module *APIModule) string {
	var fields []string
	var assertions []string

	// 从requestBody的content中提取schema
	for _, mediaType := range requestBody.Content {
		if mediaType.Schema != nil {
			schema := resolveSchemaRef(mediaType.Schema.Ref, &OpenAPI3{Components: Components{Schemas: module.Schemas}})
			if schema.Properties != nil {
				for fieldName, fieldSchema := range schema.Properties {
					goType := generateGoTypeFromSchema(fieldSchema)
					fieldNameTitle := strings.Title(fieldName)

					// 生成测试字段
					fields = append(fields, fmt.Sprintf("\t\t%s: %s,", fieldNameTitle, generateTestValue(fieldSchema, goType)))

					// 生成断言
					assertions = append(assertions, fmt.Sprintf("\tassert.Equal(t, %s, result.%s)", generateTestValue(fieldSchema, goType), fieldNameTitle))
				}
			}
		}
	}

	return fmt.Sprintf(`
func Test%s_Structure(t *testing.T) {
	// 测试结构体创建和字段访问
	result := param.%s{
%s
	}

	// 验证字段值
%s
}`, structName, structName, strings.Join(fields, "\n"), strings.Join(assertions, "\n"))
}

// 生成响应结构体测试
func generateResponseStructTest(pascal string, module *APIModule) string {
	return fmt.Sprintf(`
func Test%sResponse_Structure(t *testing.T) {
	// 测试响应结构体
	response := param.%sResponse{
		ID:   1,
		Name: "test",
		// TODO: 根据OpenAPI文档添加其他字段
	}

	// 验证字段
	assert.Equal(t, int64(1), response.ID)
	assert.Equal(t, "test", response.Name)
	// TODO: 添加其他字段验证
}`, pascal, pascal)
}

// 生成查询参数结构体测试
func generateQueryStructTest(pascal string, module *APIModule) string {
	return fmt.Sprintf(`
func Test%sParam_Structure(t *testing.T) {
	// 测试查询参数结构体
	param := param.%sParam{
		Page:  1,
		Count: 10,
		// TODO: 根据OpenAPI文档添加其他查询参数
	}

	// 验证字段
	assert.Equal(t, 1, param.Page)
	assert.Equal(t, 10, param.Count)
	// TODO: 添加其他字段验证
}`, pascal, pascal)
}

// 生成HTTP测试用例
func generateHTTPTests(pascal string, module *APIModule) string {
	var tests []string

	for _, op := range module.Operations {
		handlerName := generateHandlerName(op)
		tests = append(tests, generateHTTPTest(pascal, handlerName, op))
	}

	return strings.Join(tests, "\n")
}

// 生成HTTP测试
func generateHTTPTest(pascal, handlerName string, op APIOperation) string {
	httpMethod := getHTTPMethod(op.Method)
	methodName := generateMethodName(op)

	return fmt.Sprintf(`
func Test%sController_%s_HTTP(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		requestBody    string
		expectedStatus int
		expectedError  bool
		mockSetup      func(*Mock%sUseCase)
	}{
		{
			name:           "successful request",
			method:         %s,
			path:           "/test",
			requestBody:    `+"`"+`{}`+"`"+`,
			expectedStatus: http.StatusOK,
			expectedError:  false,
			mockSetup: func(m *Mock%sUseCase) {
				// TODO: 设置mock期望
				m.On("%s", mock.Anything, mock.MatchedBy(func(req interface{}) bool {
					return true // 接受任何有效的请求参数
				})).Return(nil, nil)
			},
		},
		{
			name:           "invalid request",
			method:         %s,
			path:           "/test",
			requestBody:    `+"`"+`{"invalid": "data"}`+"`"+`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			mockSetup: func(m *Mock%sUseCase) {
				// TODO: 设置mock期望
				m.On("%s", mock.Anything, mock.MatchedBy(func(req interface{}) bool {
					return true // 接受任何有效的请求参数
				})).Return(nil, fmt.Errorf("validation error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建Echo实例
			e := echo.New()
			
			// 注册validator
			e.Validator = &CustomValidator{validator: validator.New()}
			
			// 创建请求
			req := httptest.NewRequest(tt.method, tt.path, bytes.NewReader([]byte(tt.requestBody)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// 创建mock
			mockUseCase := new(Mock%sUseCase)
			tt.mockSetup(mockUseCase)

			// 创建控制器
			controller := &%sController{
				%s: mockUseCase,
			}

			// 执行测试
			err := controller.%s(c)

			// 验证结果
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}
			
			// 验证mock调用
			mockUseCase.AssertExpectations(t)
		})
	}
}`, pascal, handlerName, pascal, httpMethod, pascal, methodName, httpMethod, pascal, methodName, pascal, pascal, strings.ToLower(pascal), handlerName)
}

// 生成validator测试
func generateValidatorTests(pascal string, module *APIModule) string {
	return fmt.Sprintf(`
func Test%sValidator(t *testing.T) {
	validator := validator.New()
	
	tests := []struct {
		name    string
		request interface{}
		wantErr bool
	}{
		{
			name: "valid request",
			request: param.%sCreateRequest{
				Name:     "test",
				Phone:    "13800138000",
				Account:  "test@example.com",
				Password: "password123",
				Status:   1,
			},
			wantErr: false,
		},
		{
			name: "invalid request - missing required field",
			request: param.%sCreateRequest{
				// 缺少必填字段
			},
			wantErr: true,
		},
		{
			name: "invalid request - invalid email",
			request: param.%sCreateRequest{
				Name:     "test",
				Phone:    "13800138000",
				Account:  "invalid-email",
				Password: "password123",
				Status:   1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Struct(tt.request)
			
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}`, pascal, pascal, pascal, pascal)
}

// 生成测试数据
func generateTestDataFromOpenAPI(op APIOperation, pascal string) string {
	var fields []string
	methodName := generateMethodName(op)

	// 根据操作类型生成不同的测试数据
	switch strings.ToLower(methodName) {
	case "list":
		// List操作使用查询参数，不需要JSON body
		return ""
	case "create", "update":
		// Create和Update操作需要完整的用户数据
		fields = append(fields, "\t\t\tName: \"test\",")
		fields = append(fields, "\t\t\tPhone: \"13800138000\",")
		fields = append(fields, "\t\t\tAccount: \"test@example.com\",")
		fields = append(fields, "\t\t\tPassword: \"password123\",")
		fields = append(fields, "\t\t\tStatus: 1,")
		if strings.ToLower(methodName) == "update" {
			fields = append(fields, "\t\t\tId: 1,")
		}
	case "getbyid", "delete":
		// GetByID和Delete操作只需要ID，使用路径参数
		return ""
	default:
		// 默认情况
		fields = append(fields, "\t\t\tName: \"test\",")
	}

	return strings.Join(fields, "\n")
}

// 获取正确的HTTP方法
func getHTTPMethod(method string) string {
	switch strings.ToUpper(method) {
	case "GET":
		return "http.MethodGet"
	case "POST":
		return "http.MethodPost"
	case "PUT":
		return "http.MethodPut"
	case "DELETE":
		return "http.MethodDelete"
	case "PATCH":
		return "http.MethodPatch"
	default:
		return "http.MethodGet"
	}
}

// 生成测试值
func generateTestValue(schema *Schema, goType string) string {
	switch goType {
	case "string":
		return `"test"`
	case "int", "int32", "int64":
		return "1"
	case "float32", "float64":
		return "1.0"
	case "bool":
		return "true"
	case "time.Time":
		return "time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)"
	default:
		return "nil"
	}
}
