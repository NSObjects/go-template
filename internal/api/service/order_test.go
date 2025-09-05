/*
 * Generated test cases from OpenAPI3 document
 * Module: Order
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

// MockOrderUseCase 模拟业务逻辑接口
type MockOrderUseCase struct {
	mock.Mock
}

// 确保 MockOrderUseCase 实现了 biz.OrderUseCase 接口
var _ biz.OrderUseCase = (*MockOrderUseCase)(nil)

func (m *MockOrderUseCase) ListOrders(ctx context.Context, req param.OrderListOrdersRequest) ([]param.OrderListItem, int64, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]param.OrderListItem), args.Get(1).(int64), args.Error(2)
}

func (m *MockOrderUseCase) CreateOrder(ctx context.Context, req param.OrderCreateOrderRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func TestOrderController_ListOrders(t *testing.T) {
	// 创建模拟对象
	mockUseCase := new(MockOrderUseCase)
	controller := &OrderController{order: mockUseCase}

	// 创建测试请求
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// 设置模拟期望

	mockUseCase.On("ListOrders", mock.Anything, mock.Anything).Return([]param.OrderListItem{}, int64(0), nil)

	// 执行测试
	err := controller.ListOrders(c)

	// 验证结果
	assert.NoError(t, err)
	mockUseCase.AssertExpectations(t)
}

func TestOrderController_ListOrders_Error(t *testing.T) {
	// 创建模拟对象
	mockUseCase := new(MockOrderUseCase)
	controller := &OrderController{order: mockUseCase}

	// 创建测试请求
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// 设置模拟期望返回错误

	mockUseCase.On("ListOrders", mock.Anything, mock.Anything).Return([]param.OrderListItem{}, int64(0), assert.AnError)

	// 执行测试
	err := controller.ListOrders(c)

	// 验证结果
	assert.Error(t, err)
	mockUseCase.AssertExpectations(t)
}

func TestOrderController_CreateOrder(t *testing.T) {
	// 创建模拟对象
	mockUseCase := new(MockOrderUseCase)
	controller := &OrderController{order: mockUseCase}

	// 创建测试请求
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// 设置模拟期望

	mockUseCase.On("CreateOrder", mock.Anything, mock.Anything).Return(nil)

	// 执行测试
	err := controller.CreateOrder(c)

	// 验证结果
	assert.NoError(t, err)
	mockUseCase.AssertExpectations(t)
}

func TestOrderController_CreateOrder_Error(t *testing.T) {
	// 创建模拟对象
	mockUseCase := new(MockOrderUseCase)
	controller := &OrderController{order: mockUseCase}

	// 创建测试请求
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// 设置模拟期望返回错误

	mockUseCase.On("CreateOrder", mock.Anything, mock.Anything).Return(assert.AnError)

	// 执行测试
	err := controller.CreateOrder(c)

	// 验证结果
	assert.Error(t, err)
	mockUseCase.AssertExpectations(t)
}
