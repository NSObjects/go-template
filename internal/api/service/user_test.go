/*
 * Generated test cases from OpenAPI3 document
 * Module: User
 */

package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NSObjects/go-template/internal/api/biz"
	"github.com/NSObjects/go-template/internal/api/service/param"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserUseCase 模拟业务逻辑接口
type MockUserUseCase struct {
	mock.Mock
}

// 确保 MockUserUseCase 实现了 biz.UserUseCase 接口
var _ biz.UserUseCase = (*MockUserUseCase)(nil)

func (m *MockUserUseCase) ListUsers(ctx context.Context, req param.UserListUsersRequest) ([]param.UserListItem, int64, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]param.UserListItem), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserUseCase) Create(ctx context.Context, req param.UserCreateRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func TestUserController_ListUsers(t *testing.T) {
	// 创建模拟对象
	mockUseCase := new(MockUserUseCase)
	controller := &UserController{user: mockUseCase}

	// 创建测试请求
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// 设置模拟期望

	mockUseCase.On("ListUsers", mock.Anything, mock.Anything).Return([]param.UserListItem{}, int64(0), nil)

	// 执行测试
	err := controller.ListUsers(c)

	// 验证结果
	assert.NoError(t, err)
	mockUseCase.AssertExpectations(t)
}

func TestUserController_ListUsers_Error(t *testing.T) {
	// 创建模拟对象
	mockUseCase := new(MockUserUseCase)
	controller := &UserController{user: mockUseCase}

	// 创建测试请求
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// 设置模拟期望返回错误

	mockUseCase.On("ListUsers", mock.Anything, mock.Anything).Return([]param.UserListItem{}, int64(0), assert.AnError)

	// 执行测试
	err := controller.ListUsers(c)

	// 验证结果
	assert.Error(t, err)
	mockUseCase.AssertExpectations(t)
}

func TestUserController_Create(t *testing.T) {
	// 创建模拟对象
	mockUseCase := new(MockUserUseCase)
	controller := &UserController{user: mockUseCase}

	// 创建测试请求
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// 设置模拟期望

	mockUseCase.On("Create", mock.Anything, mock.Anything).Return(nil)

	// 执行测试
	err := controller.Create(c)

	// 验证结果
	assert.NoError(t, err)
	mockUseCase.AssertExpectations(t)
}

func TestUserController_Create_Error(t *testing.T) {
	// 创建模拟对象
	mockUseCase := new(MockUserUseCase)
	controller := &UserController{user: mockUseCase}

	// 创建测试请求
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// 设置模拟期望返回错误

	mockUseCase.On("Create", mock.Anything, mock.Anything).Return(assert.AnError)

	// 执行测试
	err := controller.Create(c)

	// 验证结果
	assert.Error(t, err)
	mockUseCase.AssertExpectations(t)
}
