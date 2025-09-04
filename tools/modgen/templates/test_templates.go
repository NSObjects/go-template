/*
 * 测试模板渲染函数
 */

package templates

import (
	"fmt"
	"strings"
)

// RenderBizTest 生成默认业务逻辑测试模板
func RenderBizTest(pascal, packagePath string) string {
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

// RenderServiceTest 生成默认服务层测试模板
func RenderServiceTest(pascal, packagePath string) string {
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

// renderServiceTestHeader 生成测试文件头部
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

// renderServiceTestMockInterface 生成Mock接口
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
	args := m.Called(ctx, id)
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

`, pascal, pascal, pascal, pascal, pascal, pascal, pascal, pascal, pascal, pascal, pascal, pascal)
}

// renderServiceTestList 生成List测试
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

// renderServiceTestCreate 生成Create测试
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

// renderServiceTestUpdate 生成Update测试
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

// renderServiceTestDelete 生成Delete测试
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

// renderServiceTestDetail 生成Detail测试
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
`, pascal, pascal, pascal, pascal, pascal, camel, camel)
}
