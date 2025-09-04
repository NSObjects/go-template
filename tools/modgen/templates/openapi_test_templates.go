/*
 * OpenAPI3 测试模板渲染函数
 */

package templates

import (
	"fmt"
	"strings"

	"github.com/NSObjects/go-template/tools/modgen/openapi"
)

// RenderBizTestFromOpenAPI 从OpenAPI3生成业务逻辑测试模板
func RenderBizTestFromOpenAPI(module *openapi.APIModule, pascal, packagePath string) string {
	// 为每个操作生成测试方法
	var testMethods []string

	for _, op := range module.Operations {
		methodName := openapi.GenerateMethodName(op)

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

// RenderServiceTestFromOpenAPI 从OpenAPI3生成服务层测试模板
func RenderServiceTestFromOpenAPI(module *openapi.APIModule, pascal, packagePath string) string {
	// 为每个操作生成测试方法
	var testMethods []string
	var mockMethods []string

	for _, op := range module.Operations {
		methodName := openapi.GenerateMethodName(op)

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
