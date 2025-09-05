/*
 * Generated test cases from OpenAPI3 document
 * Module: Product
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

// MockProductUseCase 模拟业务逻辑接口
type MockProductUseCase struct {
	mock.Mock
}

// 确保 MockProductUseCase 实现了 biz.ProductUseCase 接口
var _ biz.ProductUseCase = (*MockProductUseCase)(nil)

func (m *MockProductUseCase) ListProducts(ctx context.Context, req param.ProductListProductsRequest) ([]param.ProductListItem, int64, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]param.ProductListItem), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductUseCase) CreateProduct(ctx context.Context, req param.ProductCreateProductRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func TestProductController_ListProducts(t *testing.T) {
	// 创建模拟对象
	mockUseCase := new(MockProductUseCase)
	controller := &ProductController{product: mockUseCase}

	// 创建测试请求
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// 设置模拟期望

	mockUseCase.On("ListProducts", mock.Anything, mock.Anything).Return([]param.ProductListItem{}, int64(0), nil)

	// 执行测试
	err := controller.ListProducts(c)

	// 验证结果
	assert.NoError(t, err)
	mockUseCase.AssertExpectations(t)
}

func TestProductController_ListProducts_Error(t *testing.T) {
	// 创建模拟对象
	mockUseCase := new(MockProductUseCase)
	controller := &ProductController{product: mockUseCase}

	// 创建测试请求
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// 设置模拟期望返回错误

	mockUseCase.On("ListProducts", mock.Anything, mock.Anything).Return([]param.ProductListItem{}, int64(0), assert.AnError)

	// 执行测试
	err := controller.ListProducts(c)

	// 验证结果
	assert.Error(t, err)
	mockUseCase.AssertExpectations(t)
}

func TestProductController_CreateProduct(t *testing.T) {
	// 创建模拟对象
	mockUseCase := new(MockProductUseCase)
	controller := &ProductController{product: mockUseCase}

	// 创建测试请求
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// 设置模拟期望

	mockUseCase.On("CreateProduct", mock.Anything, mock.Anything).Return(nil)

	// 执行测试
	err := controller.CreateProduct(c)

	// 验证结果
	assert.NoError(t, err)
	mockUseCase.AssertExpectations(t)
}

func TestProductController_CreateProduct_Error(t *testing.T) {
	// 创建模拟对象
	mockUseCase := new(MockProductUseCase)
	controller := &ProductController{product: mockUseCase}

	// 创建测试请求
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// 设置模拟期望返回错误

	mockUseCase.On("CreateProduct", mock.Anything, mock.Anything).Return(assert.AnError)

	// 执行测试
	err := controller.CreateProduct(c)

	// 验证结果
	assert.Error(t, err)
	mockUseCase.AssertExpectations(t)
}
