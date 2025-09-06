/*
 * Generated enhanced test cases from OpenAPI3 document
 * Module: User
 * Features: Table-driven tests with comprehensive coverage
 */

package service

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NSObjects/go-template/internal/api/service/param"
	"github.com/NSObjects/go-template/internal/resp"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)


// MockUserUseCase 模拟业务逻辑接口
type MockUserUseCase struct {
	mock.Mock
}



func (m *MockUserUseCase) ListUsers(ctx context.Context, req param.UserListUsersRequest) ([]param.UserListItem, int64, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]param.UserListItem), args.Get(1).(int64), args.Error(2)
}



func (m *MockUserUseCase) Create(ctx context.Context, req param.UserCreateRequest) (*param.UserData, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*param.UserData), args.Error(1)
}




func (m *MockUserUseCase) GetByID(ctx context.Context, id int64) (*param.UserData, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*param.UserData), args.Error(1)
}



func (m *MockUserUseCase) Update(ctx context.Context, id int64, req param.UserUpdateRequest) (*param.UserData, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*param.UserData), args.Error(1)
}




func (m *MockUserUseCase) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}





func TestUserController_ListUsers(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*MockUserUseCase)
		setupRequest   func() (*http.Request, echo.Context)
		expectedStatus int
		expectedError  bool
		validateResponse func(t *testing.T, status int, body string)
	}{
		{
			name: "成功场景",
			setupMock: func(m *MockUserUseCase) {
				
				m.On("ListUsers", mock.Anything, mock.Anything).Return([]param.UserListItem{
					{
						Id:        1,
						Username:  "test",
						Email:     "test@example.com",
						
						Age:       18,
					},
				}, int64(1), nil)
				
			},
			setupRequest: func() (*http.Request, echo.Context) {
				e := echo.New()
				
				reqBody := `{"page":1,"size":10}`
				req := httptest.NewRequest(http.MethodGet, "/user", bytes.NewBufferString(reqBody))
				
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				
				return req, c
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
			validateResponse: func(t *testing.T, status int, body string) {
				assert.Equal(t, http.StatusOK, status)
				var response resp.DataResponse
				err := json.Unmarshal([]byte(body), &response)
				assert.NoError(t, err)
				assert.Equal(t, 200, response.Code)
				assert.Equal(t, "success", response.Msg)
			},
		},
		{
			name: "业务逻辑错误",
			setupMock: func(m *MockUserUseCase) {
				
				m.On("ListUsers", mock.Anything, mock.Anything).Return([]param.UserListItem{}, int64(0), assert.AnError)
				
			},
			setupRequest: func() (*http.Request, echo.Context) {
				e := echo.New()
				
				reqBody := `{"page":1,"size":10}`
				req := httptest.NewRequest(http.MethodGet, "/user", bytes.NewBufferString(reqBody))
				
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				
				return req, c
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
			validateResponse: func(t *testing.T, status int, body string) {
				// 错误响应验证
			},
		},
		
		{
			name: "无效请求体",
			setupMock: func(m *MockUserUseCase) {
				// 无效请求不会调用biz层
			},
			setupRequest: func() (*http.Request, echo.Context) {
				e := echo.New()
				req := httptest.NewRequest(http.MethodGet, "/user", bytes.NewBufferString("invalid json"))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				
				return req, c
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			validateResponse: func(t *testing.T, status int, body string) {
				assert.Equal(t, http.StatusBadRequest, status)
			},
		},
		
		
		
		{
			name: "参数验证失败",
			setupMock: func(m *MockUserUseCase) {
				// 验证失败不会调用biz层
			},
			setupRequest: func() (*http.Request, echo.Context) {
				e := echo.New()
				
				// 测试分页参数无效
				reqBody := `{"page":-1,"size":0}`
				
				req := httptest.NewRequest(http.MethodGet, "/user", bytes.NewBufferString(reqBody))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				
				return req, c
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			validateResponse: func(t *testing.T, status int, body string) {
				assert.Equal(t, http.StatusBadRequest, status)
			},
		},
		
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建mock
			mockUseCase := new(MockUserUseCase)
			tt.setupMock(mockUseCase)

			// 创建控制器
			controller := &UserController{
				user: mockUseCase,
			}

			// 设置请求
			_, c := tt.setupRequest()

			// 执行测试
			err := controller.ListUsers(c)

			// 验证结果
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// 验证响应
			if tt.validateResponse != nil {
				status := c.Response().Status
				body := ""
				// 注意：在Echo中，响应体需要通过其他方式获取
				// 这里简化处理，实际项目中可能需要更复杂的响应体获取逻辑
				tt.validateResponse(t, status, body)
			}

			// 验证mock调用
			mockUseCase.AssertExpectations(t)
		})
	}
}


func TestUserController_ListUsers_Validation(t *testing.T) {
	tests := []struct {
		name        string
		requestBody string
		expectError bool
		errorMsg    string
	}{
		
		{
			name:        "有效分页参数",
			requestBody: `{"page":1,"size":10}`,
			expectError: false,
		},
		{
			name:        "无效页码",
			requestBody: `{"page":0,"size":10}`,
			expectError: true,
			errorMsg:    "页码必须大于0",
		},
		{
			name:        "页面大小超出限制",
			requestBody: `{"page":1,"size":101}`,
			expectError: true,
			errorMsg:    "页面大小超出限制",
		},
		
		{
			name:        "无效JSON格式",
			requestBody: `invalid json`,
			expectError: true,
			errorMsg:    "无效JSON格式",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建mock
			mockUseCase := new(MockUserUseCase)
			
			if !tt.expectError {
				
				mockUseCase.On("ListUsers", mock.Anything, mock.Anything).Return([]param.UserListItem{}, int64(0), nil)
				
			}
			

			// 创建控制器
			controller := &UserController{
				user: mockUseCase,
			}

			// 创建请求
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/user", bytes.NewBufferString(tt.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			

			// 执行测试
			err := controller.ListUsers(c)

			// 验证结果
			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}

			// 验证mock调用
			mockUseCase.AssertExpectations(t)
		})
	}
}





func TestUserController_Create(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*MockUserUseCase)
		setupRequest   func() (*http.Request, echo.Context)
		expectedStatus int
		expectedError  bool
		validateResponse func(t *testing.T, status int, body string)
	}{
		{
			name: "成功场景",
			setupMock: func(m *MockUserUseCase) {
				
				m.On("Create", mock.Anything, mock.Anything).Return(&param.UserData{Id: 1}, nil)
				
				
			},
			setupRequest: func() (*http.Request, echo.Context) {
				e := echo.New()
				
				reqBody := `{"username":"testuser","email":"test@example.com"}`
				req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBufferString(reqBody))
				
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				
				return req, c
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
			validateResponse: func(t *testing.T, status int, body string) {
				assert.Equal(t, http.StatusOK, status)
				var response resp.DataResponse
				err := json.Unmarshal([]byte(body), &response)
				assert.NoError(t, err)
				assert.Equal(t, 200, response.Code)
				assert.Equal(t, "success", response.Msg)
			},
		},
		{
			name: "业务逻辑错误",
			setupMock: func(m *MockUserUseCase) {
				
				m.On("Create", mock.Anything, mock.Anything).Return(nil, assert.AnError)
				
				
			},
			setupRequest: func() (*http.Request, echo.Context) {
				e := echo.New()
				
				reqBody := `{"username":"testuser","email":"test@example.com"}`
				req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBufferString(reqBody))
				
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				
				return req, c
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
			validateResponse: func(t *testing.T, status int, body string) {
				// 错误响应验证
			},
		},
		
		{
			name: "无效请求体",
			setupMock: func(m *MockUserUseCase) {
				// 无效请求不会调用biz层
			},
			setupRequest: func() (*http.Request, echo.Context) {
				e := echo.New()
				req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBufferString("invalid json"))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				
				return req, c
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			validateResponse: func(t *testing.T, status int, body string) {
				assert.Equal(t, http.StatusBadRequest, status)
			},
		},
		
		
		
		{
			name: "参数验证失败",
			setupMock: func(m *MockUserUseCase) {
				// 验证失败不会调用biz层
			},
			setupRequest: func() (*http.Request, echo.Context) {
				e := echo.New()
				
				// 测试必填字段缺失
				reqBody := `{}`
				
				req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBufferString(reqBody))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				
				return req, c
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			validateResponse: func(t *testing.T, status int, body string) {
				assert.Equal(t, http.StatusBadRequest, status)
			},
		},
		
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建mock
			mockUseCase := new(MockUserUseCase)
			tt.setupMock(mockUseCase)

			// 创建控制器
			controller := &UserController{
				user: mockUseCase,
			}

			// 设置请求
			_, c := tt.setupRequest()

			// 执行测试
			err := controller.Create(c)

			// 验证结果
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// 验证响应
			if tt.validateResponse != nil {
				status := c.Response().Status
				body := ""
				// 注意：在Echo中，响应体需要通过其他方式获取
				// 这里简化处理，实际项目中可能需要更复杂的响应体获取逻辑
				tt.validateResponse(t, status, body)
			}

			// 验证mock调用
			mockUseCase.AssertExpectations(t)
		})
	}
}


func TestUserController_Create_Validation(t *testing.T) {
	tests := []struct {
		name        string
		requestBody string
		expectError bool
		errorMsg    string
	}{
		
		{
			name:        "有效请求体",
			requestBody: `{"username":"testuser","email":"test@example.com","age":25}`,
			expectError: false,
		},
		{
			name:        "缺少必填字段",
			requestBody: `{"username":"testuser"}`,
			expectError: true,
			errorMsg:    "缺少必填字段",
		},
		{
			name:        "无效邮箱格式",
			requestBody: `{"username":"testuser","email":"invalid-email","age":25}`,
			expectError: true,
			errorMsg:    "无效邮箱格式",
		},
		{
			name:        "用户名长度不足",
			requestBody: `{"username":"ab","email":"test@example.com","age":25}`,
			expectError: true,
			errorMsg:    "用户名长度不足",
		},
		{
			name:        "年龄超出范围",
			requestBody: `{"username":"testuser","email":"test@example.com","age":200}`,
			expectError: true,
			errorMsg:    "年龄超出范围",
		},
		
		{
			name:        "无效JSON格式",
			requestBody: `invalid json`,
			expectError: true,
			errorMsg:    "无效JSON格式",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建mock
			mockUseCase := new(MockUserUseCase)
			
			if !tt.expectError {
				
				mockUseCase.On("Create", mock.Anything, mock.Anything).Return(nil)
				
			}
			

			// 创建控制器
			controller := &UserController{
				user: mockUseCase,
			}

			// 创建请求
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBufferString(tt.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			

			// 执行测试
			err := controller.Create(c)

			// 验证结果
			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}

			// 验证mock调用
			mockUseCase.AssertExpectations(t)
		})
	}
}





func TestUserController_GetByID(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*MockUserUseCase)
		setupRequest   func() (*http.Request, echo.Context)
		expectedStatus int
		expectedError  bool
		validateResponse func(t *testing.T, status int, body string)
	}{
		{
			name: "成功场景",
			setupMock: func(m *MockUserUseCase) {
				
				m.On("GetByID", mock.Anything, mock.Anything).Return(&param.UserData{Id: 1}, nil)
				
			},
			setupRequest: func() (*http.Request, echo.Context) {
				e := echo.New()
				
				req := httptest.NewRequest(http.MethodGet, "/user", nil)
				
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				
				c.SetPath("/user/:id")
				c.SetParamNames("id")
				c.SetParamValues("1")
				
				return req, c
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
			validateResponse: func(t *testing.T, status int, body string) {
				assert.Equal(t, http.StatusOK, status)
				var response resp.DataResponse
				err := json.Unmarshal([]byte(body), &response)
				assert.NoError(t, err)
				assert.Equal(t, 200, response.Code)
				assert.Equal(t, "success", response.Msg)
			},
		},
		{
			name: "业务逻辑错误",
			setupMock: func(m *MockUserUseCase) {
				
				m.On("GetByID", mock.Anything, mock.Anything).Return(nil, assert.AnError)
				
			},
			setupRequest: func() (*http.Request, echo.Context) {
				e := echo.New()
				
				req := httptest.NewRequest(http.MethodGet, "/user", nil)
				
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				
				c.SetPath("/user/:id")
				c.SetParamNames("id")
				c.SetParamValues("1")
				
				return req, c
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
			validateResponse: func(t *testing.T, status int, body string) {
				// 错误响应验证
			},
		},
		
		
		{
			name: "无效路径参数",
			setupMock: func(m *MockUserUseCase) {
				// 无效参数不会调用biz层
			},
			setupRequest: func() (*http.Request, echo.Context) {
				e := echo.New()
				
				req := httptest.NewRequest(http.MethodGet, "/user", nil)
				
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				c.SetPath("/user/:id")
				c.SetParamNames("id")
				c.SetParamValues("invalid")
				return req, c
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			validateResponse: func(t *testing.T, status int, body string) {
				assert.Equal(t, http.StatusBadRequest, status)
			},
		},
		
		
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建mock
			mockUseCase := new(MockUserUseCase)
			tt.setupMock(mockUseCase)

			// 创建控制器
			controller := &UserController{
				user: mockUseCase,
			}

			// 设置请求
			_, c := tt.setupRequest()

			// 执行测试
			err := controller.GetByID(c)

			// 验证结果
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// 验证响应
			if tt.validateResponse != nil {
				status := c.Response().Status
				body := ""
				// 注意：在Echo中，响应体需要通过其他方式获取
				// 这里简化处理，实际项目中可能需要更复杂的响应体获取逻辑
				tt.validateResponse(t, status, body)
			}

			// 验证mock调用
			mockUseCase.AssertExpectations(t)
		})
	}
}




func TestUserController_GetByID_PathParams(t *testing.T) {
	tests := []struct {
		name        string
		pathParam   string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "有效ID",
			pathParam:   "123",
			expectError: false,
		},
		{
			name:        "零ID",
			pathParam:   "0",
			expectError: true,
			errorMsg:    "ID必须大于0",
		},
		{
			name:        "负数ID",
			pathParam:   "-1",
			expectError: true,
			errorMsg:    "ID必须大于0",
		},
		{
			name:        "非数字ID",
			pathParam:   "invalid",
			expectError: true,
			errorMsg:    "无效的ID格式",
		},
		{
			name:        "空ID",
			pathParam:   "",
			expectError: true,
			errorMsg:    "ID不能为空",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建mock
			mockUseCase := new(MockUserUseCase)
			if !tt.expectError {
				
				mockUseCase.On("GetByID", mock.Anything, mock.Anything).Return(&param.UserData{Id: 1}, nil)
				
			}

			// 创建控制器
			controller := &UserController{
				user: mockUseCase,
			}

			// 创建请求
			e := echo.New()
			
			req := httptest.NewRequest(http.MethodGet, "/user", nil)
			
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/user/:id")
			c.SetParamNames("id")
			c.SetParamValues(tt.pathParam)

			// 执行测试
			err := controller.GetByID(c)

			// 验证结果
			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}

			// 验证mock调用
			mockUseCase.AssertExpectations(t)
		})
	}
}



func TestUserController_Update(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*MockUserUseCase)
		setupRequest   func() (*http.Request, echo.Context)
		expectedStatus int
		expectedError  bool
		validateResponse func(t *testing.T, status int, body string)
	}{
		{
			name: "成功场景",
			setupMock: func(m *MockUserUseCase) {
				
				m.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(&param.UserData{Id: 1}, nil)
				
			},
			setupRequest: func() (*http.Request, echo.Context) {
				e := echo.New()
				
				reqBody := `{"username":"testuser","email":"test@example.com"}`
				req := httptest.NewRequest(http.MethodPut, "/user", bytes.NewBufferString(reqBody))
				
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				
				c.SetPath("/user/:id")
				c.SetParamNames("id")
				c.SetParamValues("1")
				
				return req, c
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
			validateResponse: func(t *testing.T, status int, body string) {
				assert.Equal(t, http.StatusOK, status)
				var response resp.DataResponse
				err := json.Unmarshal([]byte(body), &response)
				assert.NoError(t, err)
				assert.Equal(t, 200, response.Code)
				assert.Equal(t, "success", response.Msg)
			},
		},
		{
			name: "业务逻辑错误",
			setupMock: func(m *MockUserUseCase) {
				
				m.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError)
				
			},
			setupRequest: func() (*http.Request, echo.Context) {
				e := echo.New()
				
				reqBody := `{"username":"testuser","email":"test@example.com"}`
				req := httptest.NewRequest(http.MethodPut, "/user", bytes.NewBufferString(reqBody))
				
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				
				c.SetPath("/user/:id")
				c.SetParamNames("id")
				c.SetParamValues("1")
				
				return req, c
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
			validateResponse: func(t *testing.T, status int, body string) {
				// 错误响应验证
			},
		},
		
		{
			name: "无效请求体",
			setupMock: func(m *MockUserUseCase) {
				// 无效请求不会调用biz层
			},
			setupRequest: func() (*http.Request, echo.Context) {
				e := echo.New()
				req := httptest.NewRequest(http.MethodPut, "/user", bytes.NewBufferString("invalid json"))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				
				c.SetPath("/user/:id")
				c.SetParamNames("id")
				c.SetParamValues("1")
				
				return req, c
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			validateResponse: func(t *testing.T, status int, body string) {
				assert.Equal(t, http.StatusBadRequest, status)
			},
		},
		
		
		{
			name: "无效路径参数",
			setupMock: func(m *MockUserUseCase) {
				// 无效参数不会调用biz层
			},
			setupRequest: func() (*http.Request, echo.Context) {
				e := echo.New()
				
				reqBody := `{"username":"testuser","email":"test@example.com"}`
				req := httptest.NewRequest(http.MethodPut, "/user", bytes.NewBufferString(reqBody))
				
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				c.SetPath("/user/:id")
				c.SetParamNames("id")
				c.SetParamValues("invalid")
				return req, c
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			validateResponse: func(t *testing.T, status int, body string) {
				assert.Equal(t, http.StatusBadRequest, status)
			},
		},
		
		
		{
			name: "参数验证失败",
			setupMock: func(m *MockUserUseCase) {
				// 验证失败不会调用biz层
			},
			setupRequest: func() (*http.Request, echo.Context) {
				e := echo.New()
				
				// 测试必填字段缺失
				reqBody := `{}`
				
				req := httptest.NewRequest(http.MethodPut, "/user", bytes.NewBufferString(reqBody))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				
				c.SetPath("/user/:id")
				c.SetParamNames("id")
				c.SetParamValues("1")
				
				return req, c
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			validateResponse: func(t *testing.T, status int, body string) {
				assert.Equal(t, http.StatusBadRequest, status)
			},
		},
		
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建mock
			mockUseCase := new(MockUserUseCase)
			tt.setupMock(mockUseCase)

			// 创建控制器
			controller := &UserController{
				user: mockUseCase,
			}

			// 设置请求
			_, c := tt.setupRequest()

			// 执行测试
			err := controller.Update(c)

			// 验证结果
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// 验证响应
			if tt.validateResponse != nil {
				status := c.Response().Status
				body := ""
				// 注意：在Echo中，响应体需要通过其他方式获取
				// 这里简化处理，实际项目中可能需要更复杂的响应体获取逻辑
				tt.validateResponse(t, status, body)
			}

			// 验证mock调用
			mockUseCase.AssertExpectations(t)
		})
	}
}


func TestUserController_Update_Validation(t *testing.T) {
	tests := []struct {
		name        string
		requestBody string
		expectError bool
		errorMsg    string
	}{
		
		{
			name:        "有效请求体",
			requestBody: `{"username":"testuser","email":"test@example.com","age":25}`,
			expectError: false,
		},
		{
			name:        "缺少必填字段",
			requestBody: `{"username":"testuser"}`,
			expectError: true,
			errorMsg:    "缺少必填字段",
		},
		{
			name:        "无效邮箱格式",
			requestBody: `{"username":"testuser","email":"invalid-email","age":25}`,
			expectError: true,
			errorMsg:    "无效邮箱格式",
		},
		{
			name:        "用户名长度不足",
			requestBody: `{"username":"ab","email":"test@example.com","age":25}`,
			expectError: true,
			errorMsg:    "用户名长度不足",
		},
		{
			name:        "年龄超出范围",
			requestBody: `{"username":"testuser","email":"test@example.com","age":200}`,
			expectError: true,
			errorMsg:    "年龄超出范围",
		},
		
		{
			name:        "无效JSON格式",
			requestBody: `invalid json`,
			expectError: true,
			errorMsg:    "无效JSON格式",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建mock
			mockUseCase := new(MockUserUseCase)
			

			// 创建控制器
			controller := &UserController{
				user: mockUseCase,
			}

			// 创建请求
			e := echo.New()
			req := httptest.NewRequest(http.MethodPut, "/user", bytes.NewBufferString(tt.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			
			c.SetPath("/user/:id")
			c.SetParamNames("id")
			c.SetParamValues("1")
			

			// 执行测试
			err := controller.Update(c)

			// 验证结果
			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}

			// 验证mock调用
			mockUseCase.AssertExpectations(t)
		})
	}
}



func TestUserController_Update_PathParams(t *testing.T) {
	tests := []struct {
		name        string
		pathParam   string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "有效ID",
			pathParam:   "123",
			expectError: false,
		},
		{
			name:        "零ID",
			pathParam:   "0",
			expectError: true,
			errorMsg:    "ID必须大于0",
		},
		{
			name:        "负数ID",
			pathParam:   "-1",
			expectError: true,
			errorMsg:    "ID必须大于0",
		},
		{
			name:        "非数字ID",
			pathParam:   "invalid",
			expectError: true,
			errorMsg:    "无效的ID格式",
		},
		{
			name:        "空ID",
			pathParam:   "",
			expectError: true,
			errorMsg:    "ID不能为空",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建mock
			mockUseCase := new(MockUserUseCase)
			if !tt.expectError {
				
				mockUseCase.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(&param.UserData{Id: 1}, nil)
				
			}

			// 创建控制器
			controller := &UserController{
				user: mockUseCase,
			}

			// 创建请求
			e := echo.New()
			
			reqBody := `{"username":"testuser","email":"test@example.com"}`
			req := httptest.NewRequest(http.MethodPut, "/user", bytes.NewBufferString(reqBody))
			
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/user/:id")
			c.SetParamNames("id")
			c.SetParamValues(tt.pathParam)

			// 执行测试
			err := controller.Update(c)

			// 验证结果
			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}

			// 验证mock调用
			mockUseCase.AssertExpectations(t)
		})
	}
}



func TestUserController_Delete(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*MockUserUseCase)
		setupRequest   func() (*http.Request, echo.Context)
		expectedStatus int
		expectedError  bool
		validateResponse func(t *testing.T, status int, body string)
	}{
		{
			name: "成功场景",
			setupMock: func(m *MockUserUseCase) {
				
				m.On("Delete", mock.Anything, mock.Anything).Return(nil)
				
			},
			setupRequest: func() (*http.Request, echo.Context) {
				e := echo.New()
				
				req := httptest.NewRequest(http.MethodDelete, "/user", nil)
				
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				
				c.SetPath("/user/:id")
				c.SetParamNames("id")
				c.SetParamValues("1")
				
				return req, c
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
			validateResponse: func(t *testing.T, status int, body string) {
				assert.Equal(t, http.StatusOK, status)
				var response resp.DataResponse
				err := json.Unmarshal([]byte(body), &response)
				assert.NoError(t, err)
				assert.Equal(t, 200, response.Code)
				assert.Equal(t, "success", response.Msg)
			},
		},
		{
			name: "业务逻辑错误",
			setupMock: func(m *MockUserUseCase) {
				
				m.On("Delete", mock.Anything, mock.Anything).Return(assert.AnError)
				
			},
			setupRequest: func() (*http.Request, echo.Context) {
				e := echo.New()
				
				req := httptest.NewRequest(http.MethodDelete, "/user", nil)
				
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				
				c.SetPath("/user/:id")
				c.SetParamNames("id")
				c.SetParamValues("1")
				
				return req, c
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
			validateResponse: func(t *testing.T, status int, body string) {
				// 错误响应验证
			},
		},
		
		
		{
			name: "无效路径参数",
			setupMock: func(m *MockUserUseCase) {
				// 无效参数不会调用biz层
			},
			setupRequest: func() (*http.Request, echo.Context) {
				e := echo.New()
				
				req := httptest.NewRequest(http.MethodDelete, "/user", nil)
				
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				c.SetPath("/user/:id")
				c.SetParamNames("id")
				c.SetParamValues("invalid")
				return req, c
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			validateResponse: func(t *testing.T, status int, body string) {
				assert.Equal(t, http.StatusBadRequest, status)
			},
		},
		
		
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建mock
			mockUseCase := new(MockUserUseCase)
			tt.setupMock(mockUseCase)

			// 创建控制器
			controller := &UserController{
				user: mockUseCase,
			}

			// 设置请求
			_, c := tt.setupRequest()

			// 执行测试
			err := controller.Delete(c)

			// 验证结果
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// 验证响应
			if tt.validateResponse != nil {
				status := c.Response().Status
				body := ""
				// 注意：在Echo中，响应体需要通过其他方式获取
				// 这里简化处理，实际项目中可能需要更复杂的响应体获取逻辑
				tt.validateResponse(t, status, body)
			}

			// 验证mock调用
			mockUseCase.AssertExpectations(t)
		})
	}
}




func TestUserController_Delete_PathParams(t *testing.T) {
	tests := []struct {
		name        string
		pathParam   string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "有效ID",
			pathParam:   "123",
			expectError: false,
		},
		{
			name:        "零ID",
			pathParam:   "0",
			expectError: true,
			errorMsg:    "ID必须大于0",
		},
		{
			name:        "负数ID",
			pathParam:   "-1",
			expectError: true,
			errorMsg:    "ID必须大于0",
		},
		{
			name:        "非数字ID",
			pathParam:   "invalid",
			expectError: true,
			errorMsg:    "无效的ID格式",
		},
		{
			name:        "空ID",
			pathParam:   "",
			expectError: true,
			errorMsg:    "ID不能为空",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建mock
			mockUseCase := new(MockUserUseCase)
			if !tt.expectError {
				
				mockUseCase.On("Delete", mock.Anything, mock.Anything).Return(nil)
				
			}

			// 创建控制器
			controller := &UserController{
				user: mockUseCase,
			}

			// 创建请求
			e := echo.New()
			
			req := httptest.NewRequest(http.MethodDelete, "/user", nil)
			
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/user/:id")
			c.SetParamNames("id")
			c.SetParamValues(tt.pathParam)

			// 执行测试
			err := controller.Delete(c)

			// 验证结果
			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}

			// 验证mock调用
			mockUseCase.AssertExpectations(t)
		})
	}
}




