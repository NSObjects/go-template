/*
 * Generated test cases from OpenAPI3 document
 * Module: user
 */

package service

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NSObjects/echo-admin/internal/api/service/param"
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

// MockUserUseCase 模拟业务逻辑接口
type MockUserUseCase struct {
	mock.Mock
}


func (m *MockUserUseCase) GetByID(ctx context.Context, req param.UserGetByIDRequest) (*param.UserResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*param.UserResponse), args.Error(1)
}

func (m *MockUserUseCase) Update(ctx context.Context, req param.UserUpdateRequest) (*param.UserResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*param.UserResponse), args.Error(1)
}

func (m *MockUserUseCase) Delete(ctx context.Context, req param.UserDeleteRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockUserUseCase) List(ctx context.Context, req param.UserListRequest) ([]param.UserResponse, int64, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]param.UserResponse), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserUseCase) Create(ctx context.Context, req param.UserCreateRequest) (*param.UserResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*param.UserResponse), args.Error(1)
}


func TestUserController_getByID(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedError  bool
		mockSetup      func(*MockUserUseCase)
	}{
		{
			name: "valid request",
			path: "/test/1",
			expectedStatus: http.StatusOK,
			expectedError:  false,
			mockSetup: func(m *MockUserUseCase) {
				// TODO: 设置mock期望
				m.On("GetByID", mock.Anything, mock.MatchedBy(func(req interface{}) bool {
					return true // 接受任何有效的请求参数
				})).Return(nil, nil)
			},
		},
		{
			name: "invalid request",
			path: "/test/0", // 无效的ID
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			mockSetup: func(m *MockUserUseCase) {
				// TODO: 设置mock期望
				m.On("GetByID", mock.Anything, mock.MatchedBy(func(req interface{}) bool {
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
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// 创建mock
			mockUseCase := new(MockUserUseCase)
			tt.mockSetup(mockUseCase)

			// 创建控制器
			controller := &UserController{
				user: mockUseCase,
			}

			// 执行测试
			err := controller.getByID(c)

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
}

func TestUserController_update(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedError  bool
		mockSetup      func(*MockUserUseCase)
	}{
		{
			name: "valid request",
			requestBody: `			Name: "test",
			Phone: "13800138000",
			Account: "test@example.com",
			Password: "password123",
			Status: 1,
			Id: 1,`,
			expectedStatus: http.StatusOK,
			expectedError:  false,
			mockSetup: func(m *MockUserUseCase) {
				// TODO: 设置mock期望
				m.On("Update", mock.Anything, mock.MatchedBy(func(req interface{}) bool {
					return true // 接受任何有效的请求参数
				})).Return(nil, nil)
			},
		},
		{
			name: "invalid request",
			requestBody: `{"invalid": "data"}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			mockSetup: func(m *MockUserUseCase) {
				// TODO: 设置mock期望
				m.On("Update", mock.Anything, mock.MatchedBy(func(req interface{}) bool {
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
			req := httptest.NewRequest(http.MethodPut, "/test", bytes.NewReader([]byte(tt.requestBody)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// 创建mock
			mockUseCase := new(MockUserUseCase)
			tt.mockSetup(mockUseCase)

			// 创建控制器
			controller := &UserController{
				user: mockUseCase,
			}

			// 执行测试
			err := controller.update(c)

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
}

func TestUserController_delete(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedError  bool
		mockSetup      func(*MockUserUseCase)
	}{
		{
			name: "valid request",
			path: "/test/1",
			expectedStatus: http.StatusOK,
			expectedError:  false,
			mockSetup: func(m *MockUserUseCase) {
				// TODO: 设置mock期望
				m.On("Delete", mock.Anything, mock.MatchedBy(func(req interface{}) bool {
					return true // 接受任何有效的请求参数
				})).Return(nil, nil)
			},
		},
		{
			name: "invalid request",
			path: "/test/0", // 无效的ID
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			mockSetup: func(m *MockUserUseCase) {
				// TODO: 设置mock期望
				m.On("Delete", mock.Anything, mock.MatchedBy(func(req interface{}) bool {
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
			req := httptest.NewRequest(http.MethodDelete, tt.path, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// 创建mock
			mockUseCase := new(MockUserUseCase)
			tt.mockSetup(mockUseCase)

			// 创建控制器
			controller := &UserController{
				user: mockUseCase,
			}

			// 执行测试
			err := controller.delete(c)

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
}

func TestUserController_list(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		expectedError  bool
		mockSetup      func(*MockUserUseCase)
	}{
		{
			name: "valid request",
			queryParams: "?page=1&count=10&name=test&email=test@example.com",
			expectedStatus: http.StatusOK,
			expectedError:  false,
			mockSetup: func(m *MockUserUseCase) {
				// TODO: 设置mock期望
				m.On("List", mock.Anything, mock.MatchedBy(func(req interface{}) bool {
					return true // 接受任何有效的请求参数
				})).Return([]param.UserResponse{}, int64(0), nil)
			},
		},
		{
			name: "invalid request",
			queryParams: "?page=0&count=0", // 无效的分页参数
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			mockSetup: func(m *MockUserUseCase) {
				// TODO: 设置mock期望
				m.On("List", mock.Anything, mock.MatchedBy(func(req interface{}) bool {
					return true // 接受任何有效的请求参数
				})).Return([]param.UserResponse{}, int64(0), fmt.Errorf("validation error"))
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
			mockUseCase := new(MockUserUseCase)
			tt.mockSetup(mockUseCase)

			// 创建控制器
			controller := &UserController{
				user: mockUseCase,
			}

			// 执行测试
			err := controller.list(c)

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
}

func TestUserController_create(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedError  bool
		mockSetup      func(*MockUserUseCase)
	}{
		{
			name: "valid request",
			requestBody: `			Name: "test",
			Phone: "13800138000",
			Account: "test@example.com",
			Password: "password123",
			Status: 1,`,
			expectedStatus: http.StatusOK,
			expectedError:  false,
			mockSetup: func(m *MockUserUseCase) {
				// TODO: 设置mock期望
				m.On("Create", mock.Anything, mock.MatchedBy(func(req interface{}) bool {
					return true // 接受任何有效的请求参数
				})).Return(nil, nil)
			},
		},
		{
			name: "invalid request",
			requestBody: `{"invalid": "data"}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			mockSetup: func(m *MockUserUseCase) {
				// TODO: 设置mock期望
				m.On("Create", mock.Anything, mock.MatchedBy(func(req interface{}) bool {
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
			req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader([]byte(tt.requestBody)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// 创建mock
			mockUseCase := new(MockUserUseCase)
			tt.mockSetup(mockUseCase)

			// 创建控制器
			controller := &UserController{
				user: mockUseCase,
			}

			// 执行测试
			err := controller.create(c)

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
}


func TestUserController_getByID_HTTP(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		requestBody    string
		expectedStatus int
		expectedError  bool
		mockSetup      func(*MockUserUseCase)
	}{
		{
			name:           "successful request",
			method:         http.MethodGet,
			path:           "/test",
			requestBody:    `{}`,
			expectedStatus: http.StatusOK,
			expectedError:  false,
			mockSetup: func(m *MockUserUseCase) {
				// TODO: 设置mock期望
				m.On("GetByID", mock.Anything, mock.MatchedBy(func(req interface{}) bool {
					return true // 接受任何有效的请求参数
				})).Return(nil, nil)
			},
		},
		{
			name:           "invalid request",
			method:         http.MethodGet,
			path:           "/test",
			requestBody:    `{"invalid": "data"}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			mockSetup: func(m *MockUserUseCase) {
				// TODO: 设置mock期望
				m.On("GetByID", mock.Anything, mock.MatchedBy(func(req interface{}) bool {
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
			mockUseCase := new(MockUserUseCase)
			tt.mockSetup(mockUseCase)

			// 创建控制器
			controller := &UserController{
				user: mockUseCase,
			}

			// 执行测试
			err := controller.getByID(c)

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
}

func TestUserController_update_HTTP(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		requestBody    string
		expectedStatus int
		expectedError  bool
		mockSetup      func(*MockUserUseCase)
	}{
		{
			name:           "successful request",
			method:         http.MethodPut,
			path:           "/test",
			requestBody:    `{}`,
			expectedStatus: http.StatusOK,
			expectedError:  false,
			mockSetup: func(m *MockUserUseCase) {
				// TODO: 设置mock期望
				m.On("Update", mock.Anything, mock.MatchedBy(func(req interface{}) bool {
					return true // 接受任何有效的请求参数
				})).Return(nil, nil)
			},
		},
		{
			name:           "invalid request",
			method:         http.MethodPut,
			path:           "/test",
			requestBody:    `{"invalid": "data"}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			mockSetup: func(m *MockUserUseCase) {
				// TODO: 设置mock期望
				m.On("Update", mock.Anything, mock.MatchedBy(func(req interface{}) bool {
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
			mockUseCase := new(MockUserUseCase)
			tt.mockSetup(mockUseCase)

			// 创建控制器
			controller := &UserController{
				user: mockUseCase,
			}

			// 执行测试
			err := controller.update(c)

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
}

func TestUserController_delete_HTTP(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		requestBody    string
		expectedStatus int
		expectedError  bool
		mockSetup      func(*MockUserUseCase)
	}{
		{
			name:           "successful request",
			method:         http.MethodDelete,
			path:           "/test",
			requestBody:    `{}`,
			expectedStatus: http.StatusOK,
			expectedError:  false,
			mockSetup: func(m *MockUserUseCase) {
				// TODO: 设置mock期望
				m.On("Delete", mock.Anything, mock.MatchedBy(func(req interface{}) bool {
					return true // 接受任何有效的请求参数
				})).Return(nil, nil)
			},
		},
		{
			name:           "invalid request",
			method:         http.MethodDelete,
			path:           "/test",
			requestBody:    `{"invalid": "data"}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			mockSetup: func(m *MockUserUseCase) {
				// TODO: 设置mock期望
				m.On("Delete", mock.Anything, mock.MatchedBy(func(req interface{}) bool {
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
			mockUseCase := new(MockUserUseCase)
			tt.mockSetup(mockUseCase)

			// 创建控制器
			controller := &UserController{
				user: mockUseCase,
			}

			// 执行测试
			err := controller.delete(c)

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
}

func TestUserController_list_HTTP(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		requestBody    string
		expectedStatus int
		expectedError  bool
		mockSetup      func(*MockUserUseCase)
	}{
		{
			name:           "successful request",
			method:         http.MethodGet,
			path:           "/test",
			requestBody:    `{}`,
			expectedStatus: http.StatusOK,
			expectedError:  false,
			mockSetup: func(m *MockUserUseCase) {
				// TODO: 设置mock期望
				m.On("List", mock.Anything, mock.MatchedBy(func(req interface{}) bool {
					return true // 接受任何有效的请求参数
				})).Return(nil, nil)
			},
		},
		{
			name:           "invalid request",
			method:         http.MethodGet,
			path:           "/test",
			requestBody:    `{"invalid": "data"}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			mockSetup: func(m *MockUserUseCase) {
				// TODO: 设置mock期望
				m.On("List", mock.Anything, mock.MatchedBy(func(req interface{}) bool {
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
			mockUseCase := new(MockUserUseCase)
			tt.mockSetup(mockUseCase)

			// 创建控制器
			controller := &UserController{
				user: mockUseCase,
			}

			// 执行测试
			err := controller.list(c)

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
}

func TestUserController_create_HTTP(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		requestBody    string
		expectedStatus int
		expectedError  bool
		mockSetup      func(*MockUserUseCase)
	}{
		{
			name:           "successful request",
			method:         http.MethodPost,
			path:           "/test",
			requestBody:    `{}`,
			expectedStatus: http.StatusOK,
			expectedError:  false,
			mockSetup: func(m *MockUserUseCase) {
				// TODO: 设置mock期望
				m.On("Create", mock.Anything, mock.MatchedBy(func(req interface{}) bool {
					return true // 接受任何有效的请求参数
				})).Return(nil, nil)
			},
		},
		{
			name:           "invalid request",
			method:         http.MethodPost,
			path:           "/test",
			requestBody:    `{"invalid": "data"}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			mockSetup: func(m *MockUserUseCase) {
				// TODO: 设置mock期望
				m.On("Create", mock.Anything, mock.MatchedBy(func(req interface{}) bool {
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
			mockUseCase := new(MockUserUseCase)
			tt.mockSetup(mockUseCase)

			// 创建控制器
			controller := &UserController{
				user: mockUseCase,
			}

			// 执行测试
			err := controller.create(c)

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
}


func TestUserValidator(t *testing.T) {
	validator := validator.New()
	
	tests := []struct {
		name    string
		request interface{}
		wantErr bool
	}{
		{
			name: "valid request",
			request: param.UserCreateRequest{
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
			request: param.UserCreateRequest{
				// 缺少必填字段
			},
			wantErr: true,
		},
		{
			name: "invalid request - invalid email",
			request: param.UserCreateRequest{
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
}
