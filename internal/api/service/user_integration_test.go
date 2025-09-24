package service

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/NSObjects/go-template/internal/api/service/param"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserUseCase 模拟用户用例
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

func TestUserController_Integration(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name           string
		method         string
		path           string
		body           interface{}
		setupMock      func(*MockUserUseCase)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:   "GET /users - 成功获取用户列表",
			method: "GET",
			path:   "/users?page=1&size=10",
			setupMock: func(m *MockUserUseCase) {
				m.On("ListUsers", mock.Anything, param.UserListUsersRequest{
					Page: 1,
					Size: 10,
				}).Return([]param.UserListItem{
					{
						Id:       1,
						Username: "testuser",
						Email:    "test@example.com",
						Age:      25,
					},
				}, int64(1), nil)
			},
			expectedStatus: 200,
		},
		{
			name:   "POST /users - 成功创建用户",
			method: "POST",
			path:   "/users",
			body: param.UserCreateRequest{
				Username: "newuser",
				Email:    "new@example.com",
				Age:      30,
			},
			setupMock: func(m *MockUserUseCase) {
				m.On("Create", mock.Anything, param.UserCreateRequest{
					Username: "newuser",
					Email:    "new@example.com",
					Age:      30,
				}).Return(&param.UserData{
					Id:       2,
					Username: "newuser",
					Email:    "new@example.com",
					Age:      30,
				}, nil)
			},
			expectedStatus: 200,
		},
		{
			name:   "GET /users/:id - 成功获取用户",
			method: "GET",
			path:   "/users/1",
			setupMock: func(m *MockUserUseCase) {
				m.On("GetByID", mock.Anything, int64(1)).Return(&param.UserData{
					Id:       1,
					Username: "testuser",
					Email:    "test@example.com",
					Age:      25,
				}, nil)
			},
			expectedStatus: 200,
		},
		{
			name:   "PUT /users/:id - 成功更新用户",
			method: "PUT",
			path:   "/users/1",
			body: param.UserUpdateRequest{
				Username: "updateduser",
				Age:      26,
			},
			setupMock: func(m *MockUserUseCase) {
				m.On("Update", mock.Anything, int64(1), param.UserUpdateRequest{
					Username: "updateduser",
					Age:      26,
				}).Return(&param.UserData{
					Id:       1,
					Username: "updateduser",
					Email:    "test@example.com",
					Age:      26,
				}, nil)
			},
			expectedStatus: 200,
		},
		{
			name:   "DELETE /users/:id - 成功删除用户",
			method: "DELETE",
			path:   "/users/1",
			setupMock: func(m *MockUserUseCase) {
				m.On("Delete", mock.Anything, int64(1)).Return(nil)
			},
			expectedStatus: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建模拟对象
			mockUseCase := new(MockUserUseCase)
			tt.setupMock(mockUseCase)

			// 创建控制器
			controller := &UserController{user: mockUseCase}

			// 注册路由
			controller.RegisterRouter(e.Group("/api"))

			// 准备请求
			var reqBody []byte
			if tt.body != nil {
				reqBody, _ = json.Marshal(tt.body)
			}

			req := httptest.NewRequest(tt.method, tt.path, bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			// 执行请求
			c := e.NewContext(req, rec)
			c.SetPath(tt.path)

			// 根据方法调用相应的处理函数
			var err error
			switch tt.method {
			case "GET":
				if tt.path == "/users" {
					err = controller.ListUsers(c)
				} else {
					err = controller.GetByID(c)
				}
			case "POST":
				err = controller.Create(c)
			case "PUT":
				err = controller.Update(c)
			case "DELETE":
				err = controller.Delete(c)
			}

			// 验证结果
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			// 验证模拟对象调用
			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestUserController_ValidationErrors(t *testing.T) {
	e := echo.New()
	mockUseCase := new(MockUserUseCase)
	controller := &UserController{user: mockUseCase}
	controller.RegisterRouter(e.Group("/api"))

	tests := []struct {
		name           string
		method         string
		path           string
		body           interface{}
		expectedStatus int
	}{
		{
			name:           "POST /users - 无效的邮箱格式",
			method:         "POST",
			path:           "/users",
			body:           param.UserCreateRequest{Username: "test", Email: "invalid-email", Age: 25},
			expectedStatus: 400,
		},
		{
			name:           "POST /users - 用户名太短",
			method:         "POST",
			path:           "/users",
			body:           param.UserCreateRequest{Username: "ab", Email: "test@example.com", Age: 25},
			expectedStatus: 400,
		},
		{
			name:           "POST /users - 年龄超出范围",
			method:         "POST",
			path:           "/users",
			body:           param.UserCreateRequest{Username: "testuser", Email: "test@example.com", Age: 200},
			expectedStatus: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(tt.method, tt.path, bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)
			c.SetPath(tt.path)

			var err error
			switch tt.method {
			case "POST":
				err = controller.Create(c)
			}

			// 验证返回400错误
			assert.Error(t, err)
		})
	}
}
