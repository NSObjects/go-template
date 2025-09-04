/*
 * Created by generator on 2025/9/3
 * Enhanced with better UX and detailed templates
 */

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// 颜色定义
const (
	RED    = "\033[0;31m"
	GREEN  = "\033[0;32m"
	YELLOW = "\033[1;33m"
	BLUE   = "\033[0;34m"
	NC     = "\033[0m" // No Color
)

func main() {
	var name string
	var route string
	var force bool
	var openapiFile string
	var generateTests bool

	flag.StringVar(&name, "name", "", "模块名，例如: user, article")
	flag.StringVar(&route, "route", "", "基础路由前缀，例如: /articles (默认使用 name 的复数形式)")
	flag.BoolVar(&force, "force", false, "若目标文件已存在则覆盖")
	flag.StringVar(&openapiFile, "openapi", "", "OpenAPI3文档路径，例如: doc/openapi.yaml")
	flag.BoolVar(&generateTests, "tests", false, "同时生成测试用例")
	flag.Parse()

	if name == "" {
		printError("请使用 --name 指定模块名")
		fmt.Println("用法: go run tools/modgen/main.go --name=user")
		fmt.Println("或者: go run tools/modgen/main.go --name=user --openapi=doc/openapi.yaml")
		fmt.Println("或者: go run tools/modgen/main.go --name=user --openapi=doc/openapi.yaml --tests")
		os.Exit(1)
	}

	printInfo("🚀 开始生成 %s 模块...", name)

	pascal := toPascal(name)
	camel := toCamel(name)
	plural := pluralize(name)
	baseRoute := route
	if baseRoute == "" {
		baseRoute = "/" + plural
	}

	// 工作空间根目录（tools/modgen 相对）
	cwd, _ := os.Getwd()
	repoRoot := findRepoRoot(cwd)
	if repoRoot == "" {
		exitWith("未找到仓库根目录，请在项目内运行")
	}

	// 获取项目包路径
	packagePath, err := getPackagePath(repoRoot)
	if err != nil {
		exitWith(fmt.Sprintf("获取项目包路径失败: %v", err))
	}

	// 生成目标文件路径
	bizFile := filepath.Join(repoRoot, "internal", "api", "biz", fmt.Sprintf("%s.go", name))
	svcFile := filepath.Join(repoRoot, "internal", "api", "service", fmt.Sprintf("%s.go", name))
	paramFile := filepath.Join(repoRoot, "internal", "api", "service", "param", fmt.Sprintf("%s.go", name))
	modelFile := filepath.Join(repoRoot, "internal", "api", "data", "model", fmt.Sprintf("%s.go", name))
	codeFile := filepath.Join(repoRoot, "internal", "code", fmt.Sprintf("%s.go", name))

	// 根据是否提供OpenAPI文档选择生成方式
	if openapiFile != "" {
		printInfo("📄 从OpenAPI3文档生成: %s", openapiFile)
		generateFromOpenAPIDoc(name, pascal, camel, baseRoute, openapiFile, packagePath, repoRoot, force, generateTests)
	} else {
		printInfo("📝 使用默认模板生成")
		// 写入文件
		mustWrite(bizFile, renderBiz(pascal, packagePath), force)
		mustWrite(svcFile, renderService(pascal, camel, baseRoute, packagePath), force)
		mustWrite(paramFile, renderParam(pascal, packagePath), force)
		mustWrite(modelFile, renderModel(pascal, name, packagePath), force)
		mustWrite(codeFile, renderCode(pascal, name, packagePath), force)

		// 生成测试用例（如果启用）
		if generateTests {
			printInfo("🧪 生成测试用例...")
			bizTestFile := filepath.Join(repoRoot, "internal", "api", "biz", fmt.Sprintf("%s_test.go", name))
			svcTestFile := filepath.Join(repoRoot, "internal", "api", "service", fmt.Sprintf("%s_test.go", name))
			mustWrite(bizTestFile, renderBizTest(pascal, packagePath), force)
			mustWrite(svcTestFile, renderServiceTest(pascal, packagePath), force)
		}
	}

	// 注入到 fx.Options
	_ = tryInject(filepath.Join(repoRoot, "internal", "api", "biz", "biz.go"), "New"+pascal+"Handler")
	_ = tryInject(filepath.Join(repoRoot, "internal", "api", "service", "service.go"), "AsRoute(New"+pascal+"Controller)")

	printSuccess("✅ %s 模块生成完成！", name)
	fmt.Printf("\n📁 生成的文件:\n")
	fmt.Printf("  📄 业务逻辑: %s\n", bizFile)
	fmt.Printf("  📄 控制器: %s\n", svcFile)
	fmt.Printf("  📄 参数结构: %s\n", paramFile)
	fmt.Printf("  📄 数据模型: %s\n", modelFile)
	fmt.Printf("  📄 错误码: %s\n", codeFile)

	printUsageInstructions(name, pascal)
}

// 生成默认业务逻辑测试模板
func renderBizTest(pascal, packagePath string) string {
	return fmt.Sprintf(`/*
 * Generated test cases
 * Module: %s
 */

package biz

import (
	"context"
	"testing"

	"%s/internal/api/data"
	"%s/internal/api/service/param"
	"github.com/stretchr/testify/assert"
)

func Test%sHandler_List(t *testing.T) {
	// 创建handler
	handler := &%sHandler{
		dm: &data.DataManager{},
	}
	ctx := context.Background()
	req := param.%sParam{
		Page:  1,
		Count: 10,
	}

	// 测试List方法
	result, total, err := handler.List(ctx, req)

	// 由于biz层实现只是返回默认值，这里只测试方法调用不panic
	// 注意：result为nil，total为0，err为nil，这是biz层的默认实现
	assert.Nil(t, result)
	assert.Equal(t, int64(0), total)
	assert.NoError(t, err)
}

func Test%sHandler_Create(t *testing.T) {
	// 创建handler
	handler := &%sHandler{
		dm: &data.DataManager{},
	}
	ctx := context.Background()
	req := param.%sBody{
		// TODO: 填充测试数据
	}

	err := handler.Create(ctx, req)
	// 由于biz层实现只是返回默认值，这里只测试方法调用不panic
	assert.NoError(t, err)
}

func Test%sHandler_Update(t *testing.T) {
	// 创建handler
	handler := &%sHandler{
		dm: &data.DataManager{},
	}
	ctx := context.Background()
	req := param.%sBody{
		// TODO: 填充测试数据
	}

	err := handler.Update(ctx, 1, req)
	// 由于biz层实现只是返回默认值，这里只测试方法调用不panic
	assert.NoError(t, err)
}

func Test%sHandler_Delete(t *testing.T) {
	// 创建handler
	handler := &%sHandler{
		dm: &data.DataManager{},
	}
	ctx := context.Background()

	err := handler.Delete(ctx, 1)
	// 由于biz层实现只是返回默认值，这里只测试方法调用不panic
	assert.NoError(t, err)
}

func Test%sHandler_Detail(t *testing.T) {
	// 创建handler
	handler := &%sHandler{
		dm: &data.DataManager{},
	}
	ctx := context.Background()

	result, err := handler.Detail(ctx, 1)
	// 由于biz层实现只是返回默认值，这里只测试方法调用不panic
	// 注意：result为nil，err为nil，这是biz层的默认实现
	assert.Nil(t, result)
	assert.NoError(t, err)
}
`, strings.ToLower(pascal), packagePath, packagePath, pascal, pascal, pascal, pascal, pascal, pascal, pascal, pascal, pascal, pascal, pascal, pascal, pascal)
}

// 生成默认服务层测试模板
func renderServiceTest(pascal, packagePath string) string {
	camel := strings.ToLower(pascal[:1]) + pascal[1:]

	header := renderServiceTestHeader(pascal, packagePath)
	mockInterface := renderServiceTestMockInterface(pascal)
	listTest := renderServiceTestList(pascal, camel)
	createTest := renderServiceTestCreate(pascal, camel)
	updateTest := renderServiceTestUpdate(pascal, camel)
	deleteTest := renderServiceTestDelete(pascal, camel)
	detailTest := renderServiceTestDetail(pascal, camel)

	return header + mockInterface + listTest + createTest + updateTest + deleteTest + detailTest
}

// 生成测试文件头部
func renderServiceTestHeader(pascal, packagePath string) string {
	return fmt.Sprintf(`/*
 * Generated test cases
 * Module: %s
 */

package service

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"%s/internal/api/service/param"
	"%s/internal/resp"
	"%s/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

`, strings.ToLower(pascal), packagePath, packagePath, packagePath)
}

// 生成Mock接口
func renderServiceTestMockInterface(pascal string) string {
	return fmt.Sprintf(`// Mock%sUseCase 模拟业务逻辑接口
type Mock%sUseCase struct {
	mock.Mock
}

func (m *Mock%sUseCase) List(ctx context.Context, req param.%sParam) ([]param.%sResponse, int64, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]param.%sResponse), args.Get(1).(int64), args.Error(2)
}

func (m *Mock%sUseCase) Create(ctx context.Context, req param.%sBody) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *Mock%sUseCase) Update(ctx context.Context, id int64, req param.%sBody) error {
	args := m.Called(ctx, id, req)
	return args.Error(0)
}

func (m *Mock%sUseCase) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *Mock%sUseCase) Detail(ctx context.Context, id int64) (*param.%sResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*param.%sResponse), args.Error(1)
}

`, pascal, pascal, pascal, pascal, pascal, pascal, pascal, pascal, pascal, pascal, pascal, pascal, pascal, pascal)
}

// 生成List测试
func renderServiceTestList(pascal, camel string) string {
	return fmt.Sprintf(`func Test%sController_List(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		expectedError  bool
		mockSetup      func(*Mock%sUseCase)
	}{
		{
			name:           "valid request",
			queryParams:    "?page=1&count=10&name=test",
			expectedStatus: http.StatusOK,
			expectedError:  false,
			mockSetup: func(m *Mock%sUseCase) {
				m.On("List", mock.Anything, mock.MatchedBy(func(req param.%sParam) bool {
					return req.Page == 1 && req.Count == 10
				})).Return([]param.%sResponse{}, int64(0), nil)
			},
		},
		{
			name:           "invalid request - invalid page",
			queryParams:    "?page=0&count=10",
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			mockSetup: func(m *Mock%sUseCase) {
				m.On("List", mock.Anything, mock.MatchedBy(func(req param.%sParam) bool {
					return req.Page == 0
				})).Return([]param.%sResponse{}, int64(0), assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			utils.SetupTestValidator(e)
			req := httptest.NewRequest(http.MethodGet, "/test"+tt.queryParams, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// 创建mock
			mockUseCase := new(Mock%sUseCase)
			tt.mockSetup(mockUseCase)

			// 创建控制器并注入mock依赖
			controller := &%sController{
				%s: mockUseCase,
			}
			err := controller.list(c)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
				
				// 验证响应格式是否符合resp包的标准格式
				var response resp.ListResponse
				err = json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotNil(t, response.Data)
			}
			
			// 验证mock调用
			mockUseCase.AssertExpectations(t)
		})
	}
}

`, pascal, pascal, pascal, pascal, pascal, pascal, pascal, pascal, pascal, camel, camel)
}

// 生成Create测试
func renderServiceTestCreate(pascal, camel string) string {
	return fmt.Sprintf(`func Test%sController_Create(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedError  bool
		mockSetup      func(*Mock%sUseCase)
	}{
		{
			name: "valid request",
			requestBody: "{\"name\": \"test\", \"description\": \"test description\"}",
			expectedStatus: http.StatusOK,
			expectedError:  false,
			mockSetup: func(m *Mock%sUseCase) {
				m.On("Create", mock.Anything, mock.MatchedBy(func(req param.%sBody) bool {
					return req.Name == "test"
				})).Return(nil)
			},
		},
		{
			name:           "invalid request - missing required field",
			requestBody:    "{\"invalid\": \"data\"}",
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			mockSetup: func(m *Mock%sUseCase) {
				// 无效请求不会调用biz层，所以不需要设置mock期望
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			utils.SetupTestValidator(e)
			req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader([]byte(tt.requestBody)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// 创建mock
			mockUseCase := new(Mock%sUseCase)
			tt.mockSetup(mockUseCase)

			// 创建控制器并注入mock依赖
			controller := &%sController{
				%s: mockUseCase,
			}
			err := controller.create(c)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
				
				// 验证响应格式是否符合resp包的标准格式
				var response struct {
					Code int    `+"`json:\"code\"`"+`
					Msg  string `+"`json:\"msg\"`"+`
				}
				err = json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "success", response.Msg)
			}
			
			// 验证mock调用
			mockUseCase.AssertExpectations(t)
		})
	}
}

`, pascal, pascal, pascal, pascal, pascal, pascal, camel, camel)
}

// 生成Update测试
func renderServiceTestUpdate(pascal, camel string) string {
	return fmt.Sprintf(`func Test%sController_Update(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedError  bool
		mockSetup      func(*Mock%sUseCase)
	}{
		{
			name: "valid request",
			requestBody: "{\"name\": \"updated test\", \"description\": \"updated description\"}",
			expectedStatus: http.StatusOK,
			expectedError:  false,
			mockSetup: func(m *Mock%sUseCase) {
				m.On("Update", mock.Anything, int64(1), mock.MatchedBy(func(req param.%sBody) bool {
					return req.Name == "updated test"
				})).Return(nil)
			},
		},
		{
			name:           "invalid request - missing required field",
			requestBody:    "{\"invalid\": \"data\"}",
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			mockSetup: func(m *Mock%sUseCase) {
				// 无效请求不会调用biz层，所以不需要设置mock期望
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			utils.SetupTestValidator(e)
			req := httptest.NewRequest(http.MethodPut, "/test/1", bytes.NewReader([]byte(tt.requestBody)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues("1")

			// 创建mock
			mockUseCase := new(Mock%sUseCase)
			tt.mockSetup(mockUseCase)

			// 创建控制器并注入mock依赖
			controller := &%sController{
				%s: mockUseCase,
			}
			err := controller.update(c)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
				
				// 验证响应格式是否符合resp包的标准格式
				var response struct {
					Code int    `+"`json:\"code\"`"+`
					Msg  string `+"`json:\"msg\"`"+`
				}
				err = json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "success", response.Msg)
			}
			
			// 验证mock调用
			mockUseCase.AssertExpectations(t)
		})
	}
}

`, pascal, pascal, pascal, pascal, pascal, pascal, camel, camel)
}

// 生成Delete测试
func renderServiceTestDelete(pascal, camel string) string {
	return fmt.Sprintf(`func Test%sController_Delete(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedError  bool
		mockSetup      func(*Mock%sUseCase)
	}{
		{
			name:           "valid request",
			path:           "/test/1",
			expectedStatus: http.StatusOK,
			expectedError:  false,
			mockSetup: func(m *Mock%sUseCase) {
				m.On("Delete", mock.Anything, int64(1)).Return(nil)
			},
		},
		{
			name:           "invalid request - invalid id",
			path:           "/test/0",
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			mockSetup: func(m *Mock%sUseCase) {
				m.On("Delete", mock.Anything, int64(0)).Return(assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			utils.SetupTestValidator(e)
			req := httptest.NewRequest(http.MethodDelete, tt.path, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			// 从路径中提取ID
			pathParts := strings.Split(tt.path, "/")
			id := pathParts[len(pathParts)-1]
			c.SetParamValues(id)

			// 创建mock
			mockUseCase := new(Mock%sUseCase)
			tt.mockSetup(mockUseCase)

			// 创建控制器并注入mock依赖
			controller := &%sController{
				%s: mockUseCase,
			}
			err := controller.remove(c)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
				
				// 验证响应格式是否符合resp包的标准格式
				var response struct {
					Code int    `+"`json:\"code\"`"+`
					Msg  string `+"`json:\"msg\"`"+`
				}
				err = json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "success", response.Msg)
			}
			
			// 验证mock调用
			mockUseCase.AssertExpectations(t)
		})
	}
}

`, pascal, pascal, pascal, pascal, pascal, camel, camel)
}

// 生成Detail测试
func renderServiceTestDetail(pascal, camel string) string {
	return fmt.Sprintf(`func Test%sController_Detail(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedError  bool
		mockSetup      func(*Mock%sUseCase)
	}{
		{
			name:           "valid request",
			path:           "/test/1",
			expectedStatus: http.StatusOK,
			expectedError:  false,
			mockSetup: func(m *Mock%sUseCase) {
				m.On("Detail", mock.Anything, int64(1)).Return(&param.%sResponse{ID: 1, Name: "test"}, nil)
			},
		},
		{
			name:           "invalid request - invalid id",
			path:           "/test/0",
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			mockSetup: func(m *Mock%sUseCase) {
				// 无效ID可能不会调用biz层，或者会调用但返回错误
				m.On("Detail", mock.Anything, int64(0)).Return(nil, assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			utils.SetupTestValidator(e)
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			// 从路径中提取ID
			pathParts := strings.Split(tt.path, "/")
			id := pathParts[len(pathParts)-1]
			c.SetParamValues(id)

			// 创建mock
			mockUseCase := new(Mock%sUseCase)
			tt.mockSetup(mockUseCase)

			// 创建控制器并注入mock依赖
			controller := &%sController{
				%s: mockUseCase,
			}
			err := controller.detail(c)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
				
				// 验证响应格式是否符合resp包的标准格式
				var response resp.DataResponse
				err = json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotNil(t, response.Data)
			}
			
			// 验证mock调用
			mockUseCase.AssertExpectations(t)
		})
	}
}
`, pascal, pascal, pascal, pascal, pascal, pascal, camel, camel)
}

// 从OpenAPI3文档生成代码
func generateFromOpenAPIDoc(name, pascal, camel, baseRoute, openapiFile, packagePath, repoRoot string, force, generateTests bool) {
	// 解析OpenAPI文档
	openapi, err := parseOpenAPI3(openapiFile)
	if err != nil {
		exitWith(fmt.Sprintf("解析OpenAPI文档失败: %v", err))
	}

	// 生成API模块
	module, err := generateFromOpenAPI(openapi, name)
	if err != nil {
		exitWith(fmt.Sprintf("生成API模块失败: %v", err))
	}

	// 生成文件路径
	bizFile := filepath.Join(repoRoot, "internal", "api", "biz", fmt.Sprintf("%s.go", name))
	svcFile := filepath.Join(repoRoot, "internal", "api", "service", fmt.Sprintf("%s.go", name))
	paramFile := filepath.Join(repoRoot, "internal", "api", "service", "param", fmt.Sprintf("%s.go", name))
	modelFile := filepath.Join(repoRoot, "internal", "api", "data", "model", fmt.Sprintf("%s.go", name))
	codeFile := filepath.Join(repoRoot, "internal", "code", fmt.Sprintf("%s.go", name))

	// 生成代码
	mustWrite(bizFile, renderBizFromOpenAPI(module, pascal, packagePath), force)
	mustWrite(svcFile, renderServiceFromOpenAPI(module, pascal, camel, baseRoute, packagePath), force)
	mustWrite(paramFile, renderParamFromOpenAPI(module, pascal, packagePath), force)
	mustWrite(modelFile, renderModelFromOpenAPI(module, pascal, name, packagePath), force)
	mustWrite(codeFile, renderCode(pascal, name, packagePath), force)

	// 生成测试用例（如果启用）
	if generateTests {
		printInfo("🧪 生成测试用例...")
		bizTestFile := filepath.Join(repoRoot, "internal", "api", "biz", fmt.Sprintf("%s_test.go", name))
		svcTestFile := filepath.Join(repoRoot, "internal", "api", "service", fmt.Sprintf("%s_test.go", name))
		mustWrite(bizTestFile, renderBizTestFromOpenAPI(module, pascal, packagePath), force)
		mustWrite(svcTestFile, renderServiceTestFromOpenAPI(module, pascal, packagePath), force)
	}

	printInfo("📊 从OpenAPI文档解析到 %d 个操作", len(module.Operations))
}

// 彩色输出函数
func printInfo(format string, args ...interface{}) {
	fmt.Printf(BLUE+"[INFO]"+NC+" "+format+"\n", args...)
}

func printSuccess(format string, args ...interface{}) {
	fmt.Printf(GREEN+"[SUCCESS]"+NC+" "+format+"\n", args...)
}

func printWarning(format string, args ...interface{}) {
	fmt.Printf(YELLOW+"[WARNING]"+NC+" "+format+"\n", args...)
}

func printError(format string, args ...interface{}) {
	fmt.Printf(RED+"[ERROR]"+NC+" "+format+"\n", args...)
}

func exitWith(msg string) {
	printError(msg)
	os.Exit(1)
}

// 打印使用说明
func printUsageInstructions(name, pascal string) {
	fmt.Printf("\n📖 %s 模块使用说明:\n", name)
	fmt.Println("1. 参数结构: internal/api/service/param/" + name + ".go")
	fmt.Println("2. 业务逻辑: internal/api/biz/" + name + ".go")
	fmt.Println("3. 控制器: internal/api/service/" + name + ".go")
	fmt.Println("4. 数据模型: internal/api/data/model/" + name + ".go")
	fmt.Println("5. 错误码: internal/code/" + name + ".go")
	fmt.Println("\n🔧 下一步操作:")
	fmt.Println("1. 根据业务需求修改参数结构和数据模型")
	fmt.Println("2. 实现具体的业务逻辑")
	fmt.Println("3. 配置路由和中间件")
	fmt.Println("4. 运行 'make gen-code' 生成错误码文档")
	fmt.Println("5. 运行 'make run' 启动服务")
	fmt.Println("\n💡 提示:")
	fmt.Printf("- 如未自动注册，请手动将 New%[1]sHandler 和 AsRoute(New%[1]sController) 加入 fx.Options\n", pascal)
	fmt.Println("- 使用 'make db-gen' 生成数据库模型")
	fmt.Println("- 使用 'make gen-code' 生成错误码文档")
}

func findRepoRoot(start string) string {
	dir := start
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

// getPackagePath 从go.mod文件获取项目包路径
func getPackagePath(repoRoot string) (string, error) {
	goModPath := filepath.Join(repoRoot, "go.mod")
	file, err := os.Open(goModPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		line := scanner.Text()
		// 解析 module github.com/NSObjects/go-template
		parts := strings.Fields(line)
		if len(parts) >= 2 && parts[0] == "module" {
			return parts[1], nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", fmt.Errorf("无法从go.mod解析模块路径")
}

func mustWrite(path, content string, force bool) {
	if !force {
		if _, err := os.Stat(path); err == nil {
			fmt.Printf("跳过已存在文件: %s (使用 --force 可覆盖)\n", path)
			return
		}
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		exitWith(err.Error())
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		exitWith(err.Error())
	}
	fmt.Printf("写入: %s\n", path)
}

func tryInject(filePath, item string) error {
	b, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	s := string(b)
	if strings.Contains(s, item) {
		return nil
	}
	// 优先匹配 fx.Provide(...)，回退匹配 var Model = fx.Options(...)
	patterns := []string{
		`fx\.Provide\(((?s:.*?))\)`,
		`var\s+Model\s*=\s*fx\.Options\(((?s:.*?))\)`,
	}
	var (
		before string
		inside string
		after  string
		found  bool
	)
	for _, p := range patterns {
		re := regexp.MustCompile(p)
		loc := re.FindStringSubmatchIndex(s)
		if loc == nil {
			continue
		}
		// loc: [start0,end0,start1,end1]
		before = s[:loc[2]]
		inside = s[loc[2]:loc[3]]
		after = s[loc[3]:]
		found = true
		break
	}
	if !found {
		return nil
	}

	// 尝试获取上一行缩进
	indent := "\t"
	if li := strings.LastIndex(inside, "\n"); li >= 0 {
		line := inside[li+1:]
		indent = leadingWhitespace(line)
		if indent == "" {
			indent = "\t"
		}
	}
	// 若内部非空且末尾没有逗号，补一个逗号
	trimmed := strings.TrimSpace(inside)
	if trimmed != "" && !strings.HasSuffix(strings.TrimSpace(trimmed), ",") {
		inside = inside + ",\n"
	}
	inside = inside + indent + item + ",\n"
	out := before + inside + after
	return os.WriteFile(filePath, []byte(out), 0o644)
}

func leadingWhitespace(s string) string {
	for i, r := range s {
		if r != ' ' && r != '\t' {
			return s[:i]
		}
	}
	return s
}

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

func toCamel(s string) string { // nolint: revive
	p := toPascal(s)
	if p == "" {
		return p
	}
	return strings.ToLower(p[:1]) + p[1:]
}

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

func pluralize(s string) string {
	// 简易复数：以 y 结尾改 ies，其它加 s
	if strings.HasSuffix(s, "y") && len(s) > 1 && !isVowel(s[len(s)-2]) {
		return s[:len(s)-1] + "ies"
	}
	if strings.HasSuffix(s, "s") {
		return s + "es"
	}
	return s + "s"
}

func isVowel(b byte) bool {
	switch b {
	case 'a', 'e', 'i', 'o', 'u':
		return true
	default:
		return false
	}
}

func renderBiz(pascal, packagePath string) string {
	return "package biz\n\n" +
		"import (\n" +
		"\t\"context\"\n" +
		fmt.Sprintf("\t\"%s/internal/api/data\"\n", packagePath) +
		fmt.Sprintf("\t\"%s/internal/api/data/model\"\n", packagePath) +
		fmt.Sprintf("\t\"%s/internal/api/service/param\"\n", packagePath) +
		")\n\n" +
		"// " + pascal + "UseCase " + pascal + "业务用例接口\n" +
		fmt.Sprintf("type %sUseCase interface {\n", pascal) +
		fmt.Sprintf("\tList(ctx context.Context, p param.%[1]sParam) ([]param.%[1]sResponse, int64, error)\n", pascal) +
		fmt.Sprintf("\tCreate(ctx context.Context, b param.%[1]sBody) error\n", pascal) +
		fmt.Sprintf("\tUpdate(ctx context.Context, id int64, b param.%[1]sBody) error\n", pascal) +
		fmt.Sprintf("\tDelete(ctx context.Context, id int64) error\n") +
		fmt.Sprintf("\tDetail(ctx context.Context, id int64) (*param.%[1]sResponse, error)\n", pascal) +
		"}\n\n" +
		"// " + pascal + "Handler " + pascal + "业务处理器\n" +
		fmt.Sprintf("type %sHandler struct {\n", pascal) +
		"\tdm *data.DataManager\n" +
		"}\n\n" +
		fmt.Sprintf("// New%[1]sHandler 创建%[1]s业务处理器\n", pascal) +
		fmt.Sprintf("func New%[1]sHandler(dm *data.DataManager) *%[1]sHandler {\n", pascal) +
		fmt.Sprintf("\treturn &%[1]sHandler{dm: dm}\n", pascal) +
		"}\n\n" +
		"// List 获取" + pascal + "列表\n" +
		fmt.Sprintf("func (h *%[1]sHandler) List(ctx context.Context, p param.%[1]sParam) ([]param.%[1]sResponse, int64, error) {\n", pascal) +
		"\t// TODO: 实现列表查询逻辑\n" +
		"\t// 示例：\n" +
		"\t// var models []model." + pascal + "\n" +
		"\t// if err := h.dm.MySQLWithContext(ctx).Offset(p.Offset()).Limit(p.Limit()).Find(&models).Error; err != nil {\n" +
		"\t//     return nil, 0, code.WrapDatabaseError(err, \"query " + pascal + " list\")\n" +
		"\t// }\n" +
		"\t// var total int64\n" +
		"\t// h.dm.MySQLWithContext(ctx).Model(&model." + pascal + "{}).Count(&total)\n" +
		"\t// return convert%sToResponses(models), total, nil\n" +
		"\treturn nil, 0, nil\n" +
		"}\n\n" +
		"// Create 创建" + pascal + "\n" +
		fmt.Sprintf("func (h *%[1]sHandler) Create(ctx context.Context, b param.%[1]sBody) error {\n", pascal) +
		"\t// TODO: 实现创建逻辑\n" +
		"\t// 示例：\n" +
		"\t// model := &model." + pascal + "{\n" +
		"\t//     // 设置字段\n" +
		"\t//     CreatedAt: time.Now(),\n" +
		"\t// }\n" +
		"\t// if err := h.dm.MySQLWithContext(ctx).Create(model).Error; err != nil {\n" +
		"\t//     return nil, code.WrapDatabaseError(err, \"create " + pascal + "\")\n" +
		"\t// }\n" +
		"\t// 创建成功\n" +
		"\treturn nil\n" +
		"}\n\n" +
		"// Update 更新" + pascal + "\n" +
		fmt.Sprintf("func (h *%[1]sHandler) Update(ctx context.Context, id int64, b param.%[1]sBody) error {\n", pascal) +
		"\t// TODO: 实现更新逻辑\n" +
		"\t// 示例：\n" +
		"\t// var model model." + pascal + "\n" +
		"\t// if err := h.dm.MySQLWithContext(ctx).First(&model, id).Error; err != nil {\n" +
		"\t//     if errors.Is(err, gorm.ErrRecordNotFound) {\n" +
		"\t//         return nil, code.WrapNotFoundError(nil, \"" + pascal + " not found\")\n" +
		"\t//     }\n" +
		"\t//     return nil, code.WrapDatabaseError(err, \"query " + pascal + "\")\n" +
		"\t// }\n" +
		"\t// // 更新字段\n" +
		"\t// model.UpdatedAt = time.Now()\n" +
		"\t// if err := h.dm.MySQLWithContext(ctx).Save(&model).Error; err != nil {\n" +
		"\t//     return code.WrapDatabaseError(err, \"update " + pascal + "\")\n" +
		"\t// }\n" +
		"\t// 更新成功\n" +
		"\treturn nil\n" +
		"}\n\n" +
		"// Delete 删除" + pascal + "\n" +
		fmt.Sprintf("func (h *%[1]sHandler) Delete(ctx context.Context, id int64) error {\n", pascal) +
		"\t// TODO: 实现删除逻辑\n" +
		"\t// 示例：\n" +
		"\t// if err := h.dm.MySQLWithContext(ctx).Delete(&model." + pascal + "{}, id).Error; err != nil {\n" +
		"\t//     return code.WrapDatabaseError(err, \"delete " + pascal + "\")\n" +
		"\t// }\n" +
		"\t// return nil\n" +
		"\treturn nil\n" +
		"}\n\n" +
		"// Detail 获取" + pascal + "详情\n" +
		fmt.Sprintf("func (h *%[1]sHandler) Detail(ctx context.Context, id int64) (*param.%[1]sResponse, error) {\n", pascal) +
		"\t// TODO: 实现详情查询逻辑\n" +
		"\t// 示例：\n" +
		"\t// var model model." + pascal + "\n" +
		"\t// if err := h.dm.MySQLWithContext(ctx).First(&model, id).Error; err != nil {\n" +
		"\t//     if errors.Is(err, gorm.ErrRecordNotFound) {\n" +
		"\t//         return nil, code.WrapNotFoundError(nil, \"" + pascal + " not found\")\n" +
		"\t//     }\n" +
		"\t//     return nil, code.WrapDatabaseError(err, \"query " + pascal + "\")\n" +
		"\t// }\n" +
		"\t// return convert%sToResponse(&model), nil\n" +
		"\treturn nil, nil\n" +
		"}\n\n" +
		"// convert%sToResponse 转换为响应结构\n" +
		fmt.Sprintf("func convert%sToResponse(model *model.%[1]s) *param.%[1]sResponse {\n", pascal, pascal) +
		"\t// TODO: 实现转换逻辑\n" +
		"\treturn &param." + pascal + "Response{\n" +
		"\t\t// ID: model.ID,\n" +
		"\t\t// CreatedAt: model.CreatedAt,\n" +
		"\t\t// UpdatedAt: model.UpdatedAt,\n" +
		"\t}\n" +
		"}\n\n" +
		"// convert%sToResponses 转换为响应结构列表\n" +
		fmt.Sprintf("func convert%sToResponses(models []model.%[1]s) []param.%[1]sResponse {\n", pascal, pascal) +
		"\tresponses := make([]param." + pascal + "Response, len(models))\n" +
		"\tfor i, model := range models {\n" +
		"\t\tresponses[i] = *convert" + pascal + "ToResponse(&model)\n" +
		"\t}\n" +
		"\treturn responses\n" +
		"}\n"
}

func renderService(pascal, camel, baseRoute, packagePath string) string {
	return "package service\n\n" +
		"import (\n" +
		"\t\"strconv\"\n" +
		fmt.Sprintf("\t\"%s/internal/api/biz\"\n", packagePath) +
		fmt.Sprintf("\t\"%s/internal/api/service/param\"\n", packagePath) +
		fmt.Sprintf("\t\"%s/internal/resp\"\n", packagePath) +
		fmt.Sprintf("\t\"%s/internal/utils\"\n", packagePath) +
		"\t\"github.com/labstack/echo/v4\"\n" +
		")\n\n" +
		fmt.Sprintf("type %sController struct {\n\t%s biz.%sUseCase\n}\n\n", toCamel(pascal), camel, pascal) +
		fmt.Sprintf("func New%[1]sController(h *biz.%[1]sHandler) RegisterRouter {\n\treturn &%[2]sController{%[2]s: h}\n}\n\n", pascal, toCamel(pascal)) +
		fmt.Sprintf("func (c *%[1]sController) RegisterRouter(g *echo.Group, m ...echo.MiddlewareFunc) {\n\tg.GET(\"%[2]s\", c.list).Name = \"列表示例\"\n\tg.POST(\"%[2]s\", c.create).Name = \"创建示例\"\n\tg.GET(\"%[2]s/:id\", c.detail).Name = \"详情示例\"\n\tg.PUT(\"%[2]s/:id\", c.update).Name = \"更新示例\"\n\tg.DELETE(\"%[2]s/:id\", c.remove).Name = \"删除示例\"\n}\n\n", toCamel(pascal), baseRoute) +
		fmt.Sprintf("func (c *%[1]sController) list(ctx echo.Context) error {\n\tvar p param.%[2]sParam\n\tif err := BindAndValidate(ctx, &p); err != nil { return err }\n\tbizCtx := utils.BuildContext(ctx)\n\titems, total, err := c.%[3]s.List(bizCtx, p)\n\tif err != nil { return err }\n\treturn resp.ListDataResponse(items, total, ctx)\n}\n\n", toCamel(pascal), pascal, camel) +
		fmt.Sprintf("func (c *%[1]sController) detail(ctx echo.Context) error {\n\tid, _ := strconv.ParseInt(ctx.Param(\"id\"), 10, 64)\n\tbizCtx := utils.BuildContext(ctx)\n\titem, err := c.%[2]s.Detail(bizCtx, id)\n\tif err != nil { return err }\n\treturn resp.OneDataResponse(item, ctx)\n}\n\n", toCamel(pascal), camel) +
		fmt.Sprintf("func (c *%[1]sController) create(ctx echo.Context) error {\n\tvar b param.%[2]sBody\n\tif err := BindAndValidate(ctx, &b); err != nil { return err }\n\tbizCtx := utils.BuildContext(ctx)\n\tif err := c.%[3]s.Create(bizCtx, b); err != nil { return err }\n\treturn resp.OperateSuccess(ctx)\n}\n\n", toCamel(pascal), pascal, camel) +
		fmt.Sprintf("func (c *%[1]sController) update(ctx echo.Context) error {\n\tid, _ := strconv.ParseInt(ctx.Param(\"id\"), 10, 64)\n\tvar b param.%[2]sBody\n\tif err := BindAndValidate(ctx, &b); err != nil { return err }\n\tbizCtx := utils.BuildContext(ctx)\n\tif err := c.%[3]s.Update(bizCtx, id, b); err != nil { return err }\n\treturn resp.OperateSuccess(ctx)\n}\n\n", toCamel(pascal), pascal, camel) +
		fmt.Sprintf("func (c *%[1]sController) remove(ctx echo.Context) error {\n\tid, _ := strconv.ParseInt(ctx.Param(\"id\"), 10, 64)\n\tbizCtx := utils.BuildContext(ctx)\n\tif err := c.%[2]s.Delete(bizCtx, id); err != nil { return err }\n\treturn resp.OperateSuccess(ctx)\n}\n", toCamel(pascal), camel)
}

func renderParam(pascal, packagePath string) string {
	return "package param\n\n" +
		"import \"time\"\n\n" +
		"// " + pascal + "Param 查询参数\n" +
		fmt.Sprintf("type %sParam struct {\n", pascal) +
		"\tPage  int    `json:\"page\" form:\"page\" query:\"page\"`\n" +
		"\tCount int    `json:\"count\" form:\"count\" query:\"count\"`\n" +
		"\tName  string `json:\"name\" form:\"name\" query:\"name\"`\n" +
		"\t// TODO: 添加更多查询字段\n" +
		"}\n\n" +
		"// Limit 获取限制数量\n" +
		fmt.Sprintf("func (p %sParam) Limit() int {\n", pascal) +
		"\tif p.Count <= 0 {\n" +
		"\t\treturn 10\n" +
		"\t}\n" +
		"\treturn p.Count\n" +
		"}\n\n" +
		"// Offset 获取偏移量\n" +
		fmt.Sprintf("func (p %sParam) Offset() int {\n", pascal) +
		"\tif p.Page <= 1 {\n" +
		"\t\treturn 0\n" +
		"\t}\n" +
		"\treturn (p.Page - 1) * p.Limit()\n" +
		"}\n\n" +
		"// " + pascal + "Body 创建/更新请求体\n" +
		fmt.Sprintf("type %sBody struct {\n", pascal) +
		"\tName        string `json:\"name\" validate:\"required\"`\n" +
		"\tDescription string `json:\"description\"`\n" +
		"\t// TODO: 添加更多字段\n" +
		"}\n\n" +
		"// " + pascal + "Response 响应结构\n" +
		fmt.Sprintf("type %sResponse struct {\n", pascal) +
		"\tID          uint      `json:\"id\"`\n" +
		"\tName        string    `json:\"name\"`\n" +
		"\tDescription string    `json:\"description\"`\n" +
		"\tCreatedAt   time.Time `json:\"created_at\"`\n" +
		"\tUpdatedAt   time.Time `json:\"updated_at\"`\n" +
		"\t// TODO: 添加更多返回字段\n" +
		"}\n"
}

func renderModel(pascal, table, packagePath string) string {
	return "package model\n\n" +
		"import (\n" +
		"\t\"time\"\n" +
		"\t\"gorm.io/gorm\"\n" +
		")\n\n" +
		"// " + pascal + " 数据模型\n" +
		fmt.Sprintf("type %s struct {\n", pascal) +
		"\tID          uint           `gorm:\"primaryKey;autoIncrement\" json:\"id\"`\n" +
		"\tName        string         `gorm:\"column:name;type:varchar(100);not null\" json:\"name\"`\n" +
		"\tDescription string         `gorm:\"column:description;type:text\" json:\"description\"`\n" +
		"\tStatus      int            `gorm:\"column:status;type:int;default:1\" json:\"status\"`\n" +
		"\tCreatedAt   time.Time      `gorm:\"column:created_at\" json:\"created_at\"`\n" +
		"\tUpdatedAt   time.Time      `gorm:\"column:updated_at\" json:\"updated_at\"`\n" +
		"\tDeletedAt   gorm.DeletedAt `gorm:\"column:deleted_at;index\" json:\"-\"`\n" +
		"\t// TODO: 添加更多字段\n" +
		"}\n\n" +
		"// TableName 指定表名\n" +
		fmt.Sprintf("func (%s) TableName() string {\n", pascal) +
		fmt.Sprintf("\treturn \"%s\"\n", table) +
		"}\n"
}

// renderCode 生成业务错误码文件
func renderCode(pascal, table, packagePath string) string {
	// 计算错误码起始值（基于表名）
	baseCode := 100000 + int(table[0])*1000 + int(table[len(table)-1])*10

	return "package code\n\n" +
		"//go:generate codegen -type=int\n" +
		"//go:generate codegen -type=int -doc -output ./error_code_generated.md\n\n" +
		fmt.Sprintf("// %s相关错误码\n", pascal) +
		fmt.Sprintf("const (\n") +
		fmt.Sprintf("\t// Err%sNotFound - 404: %s not found.\n", pascal, pascal) +
		fmt.Sprintf("\tErr%sNotFound int = iota + %d\n", pascal, baseCode) +
		fmt.Sprintf("\t// Err%sAlreadyExists - 400: %s already exists.\n", pascal, pascal) +
		fmt.Sprintf("\tErr%sAlreadyExists\n", pascal) +
		fmt.Sprintf("\t// Err%sInvalidData - 400: %s invalid data.\n", pascal, pascal) +
		fmt.Sprintf("\tErr%sInvalidData\n", pascal) +
		fmt.Sprintf("\t// Err%sPermissionDenied - 403: %s permission denied.\n", pascal, pascal) +
		fmt.Sprintf("\tErr%sPermissionDenied\n", pascal) +
		fmt.Sprintf("\t// Err%sInUse - 400: %s is in use.\n", pascal, pascal) +
		fmt.Sprintf("\tErr%sInUse\n", pascal) +
		fmt.Sprintf("\t// Err%sCreateFailed - 500: %s create failed.\n", pascal, pascal) +
		fmt.Sprintf("\tErr%sCreateFailed\n", pascal) +
		fmt.Sprintf("\t// Err%sUpdateFailed - 500: %s update failed.\n", pascal, pascal) +
		fmt.Sprintf("\tErr%sUpdateFailed\n", pascal) +
		fmt.Sprintf("\t// Err%sDeleteFailed - 500: %s delete failed.\n", pascal, pascal) +
		fmt.Sprintf("\tErr%sDeleteFailed\n", pascal) +
		fmt.Sprintf(")\n")
}
