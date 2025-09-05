/*
 * 模板渲染器测试用例
 */

package templates

import (
	"strings"
	"testing"
)

func TestNewTemplateRenderer(t *testing.T) {
	renderer, err := NewTemplateRenderer()
	if err != nil {
		t.Fatalf("创建模板渲染器失败: %v", err)
	}

	if renderer == nil {
		t.Fatal("模板渲染器不应为 nil")
	}

	if len(renderer.templates) == 0 {
		t.Fatal("模板不应为空")
	}

	// 验证所有模板都已加载
	expectedTemplates := []string{"biz", "service", "param", "model", "code"}
	for _, name := range expectedTemplates {
		if _, exists := renderer.templates[name]; !exists {
			t.Errorf("模板 %s 未加载", name)
		}
	}
}

func TestRenderBiz(t *testing.T) {
	renderer, err := NewTemplateRenderer()
	if err != nil {
		t.Fatalf("创建模板渲染器失败: %v", err)
	}

	content, err := renderer.RenderBiz("User", "github.com/test/project")
	if err != nil {
		t.Fatalf("渲染业务逻辑模板失败: %v", err)
	}

	// 验证关键内容
	expectedContents := []string{
		"package biz",
		"type UserUseCase interface",
		"type UserHandler struct",
		"func NewUserHandler",
		"func (h *UserHandler) List",
		"func (h *UserHandler) Create",
		"func (h *UserHandler) Update",
		"func (h *UserHandler) Delete",
		"func (h *UserHandler) Detail",
		"github.com/test/project/internal/api/data",
		"github.com/test/project/internal/api/data/model",
		"github.com/test/project/internal/api/service/param",
	}

	for _, expected := range expectedContents {
		if !strings.Contains(content, expected) {
			t.Errorf("业务逻辑模板缺少内容: %s", expected)
		}
	}
}

func TestRenderService(t *testing.T) {
	renderer, err := NewTemplateRenderer()
	if err != nil {
		t.Fatalf("创建模板渲染器失败: %v", err)
	}

	content, err := renderer.RenderService("User", "user", "/users", "github.com/test/project")
	if err != nil {
		t.Fatalf("渲染服务层模板失败: %v", err)
	}

	// 验证关键内容
	expectedContents := []string{
		"package service",
		"type userController struct",
		"func NewUserController",
		"func (c *userController) RegisterRouter",
		"g.GET(\"/users\", c.list)",
		"g.POST(\"/users\", c.create)",
		"g.GET(\"/users/:id\", c.detail)",
		"g.PUT(\"/users/:id\", c.update)",
		"g.DELETE(\"/users/:id\", c.remove)",
		"github.com/test/project/internal/api/biz",
		"github.com/test/project/internal/api/service/param",
		"github.com/test/project/internal/resp",
		"github.com/test/project/internal/utils",
		"github.com/labstack/echo/v4",
	}

	for _, expected := range expectedContents {
		if !strings.Contains(content, expected) {
			t.Errorf("服务层模板缺少内容: %s", expected)
		}
	}
}

func TestRenderParam(t *testing.T) {
	renderer, err := NewTemplateRenderer()
	if err != nil {
		t.Fatalf("创建模板渲染器失败: %v", err)
	}

	content, err := renderer.RenderParam("User", "github.com/test/project")
	if err != nil {
		t.Fatalf("渲染参数模板失败: %v", err)
	}

	// 验证关键内容
	expectedContents := []string{
		"package param",
		"type UserParam struct",
		"type UserBody struct",
		"type UserResponse struct",
		"func (p UserParam) Limit() int",
		"func (p UserParam) Offset() int",
		"json:\"page\" form:\"page\" query:\"page\"",
		"json:\"name\" validate:\"required\"",
		"json:\"id\"",
		"json:\"created_at\"",
		"json:\"updated_at\"",
		"time.Time",
	}

	for _, expected := range expectedContents {
		if !strings.Contains(content, expected) {
			t.Errorf("参数模板缺少内容: %s", expected)
		}
	}
}

func TestRenderModel(t *testing.T) {
	renderer, err := NewTemplateRenderer()
	if err != nil {
		t.Fatalf("创建模板渲染器失败: %v", err)
	}

	content, err := renderer.RenderModel("User", "users", "github.com/test/project")
	if err != nil {
		t.Fatalf("渲染模型模板失败: %v", err)
	}

	// 验证关键内容
	expectedContents := []string{
		"package model",
		"type User struct",
		"func (User) TableName() string",
		"return \"users\"",
		"gorm:\"primaryKey;autoIncrement\"",
		"gorm:\"column:name;type:varchar(100);not null\"",
		"gorm:\"column:description;type:text\"",
		"gorm:\"column:status;type:int;default:1\"",
		"gorm:\"column:created_at\"",
		"gorm:\"column:updated_at\"",
		"gorm:\"column:deleted_at;index\"",
		"gorm.DeletedAt",
		"time.Time",
	}

	for _, expected := range expectedContents {
		if !strings.Contains(content, expected) {
			t.Errorf("模型模板缺少内容: %s", expected)
		}
	}
}

func TestRenderCode(t *testing.T) {
	renderer, err := NewTemplateRenderer()
	if err != nil {
		t.Fatalf("创建模板渲染器失败: %v", err)
	}

	content, err := renderer.RenderCode("User", "users", "github.com/test/project")
	if err != nil {
		t.Fatalf("渲染错误码模板失败: %v", err)
	}

	// 验证关键内容
	expectedContents := []string{
		"package code",
		"//go:generate codegen -type=int",
		"//go:generate codegen -type=int -doc -output ./error_code_generated.md",
		"// User相关错误码",
		"const (",
		"ErrUserNotFound int = iota +",
		"ErrUserAlreadyExists",
		"ErrUserInvalidData",
		"ErrUserPermissionDenied",
		"ErrUserInUse",
		"ErrUserCreateFailed",
		"ErrUserUpdateFailed",
		"ErrUserDeleteFailed",
	}

	for _, expected := range expectedContents {
		if !strings.Contains(content, expected) {
			t.Errorf("错误码模板缺少内容: %s", expected)
		}
	}
}

func TestRenderWithSpecialCharacters(t *testing.T) {
	renderer, err := NewTemplateRenderer()
	if err != nil {
		t.Fatalf("创建模板渲染器失败: %v", err)
	}

	// 测试下划线命名
	content, err := renderer.RenderBiz("UserProfile", "github.com/test/project")
	if err != nil {
		t.Fatalf("渲染下划线命名模板失败: %v", err)
	}

	expectedContents := []string{
		"type UserProfileUseCase interface",
		"type UserProfileHandler struct",
		"func NewUserProfileHandler",
		"param.UserProfileParam",
		"param.UserProfileResponse",
		"param.UserProfileBody",
	}

	for _, expected := range expectedContents {
		if !strings.Contains(content, expected) {
			t.Errorf("下划线命名模板缺少内容: %s", expected)
		}
	}
}

func TestRenderErrorHandling(t *testing.T) {
	renderer, err := NewTemplateRenderer()
	if err != nil {
		t.Fatalf("创建模板渲染器失败: %v", err)
	}

	// 测试不存在的模板
	_, err = renderer.Render("nonexistent", TemplateData{})
	if err == nil {
		t.Error("应该返回模板不存在的错误")
	}

	expectedError := "模板 nonexistent 不存在"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("错误信息不正确，期望包含: %s，实际: %s", expectedError, err.Error())
	}
}

func TestTemplateData(t *testing.T) {
	renderer, err := NewTemplateRenderer()
	if err != nil {
		t.Fatalf("创建模板渲染器失败: %v", err)
	}

	// 测试各种命名转换
	testCases := []struct {
		name   string
		pascal string
		camel  string
		table  string
		route  string
	}{
		{"user", "User", "user", "user", "/users"},
		{"user_profile", "UserProfile", "userProfile", "user_profile", "/user-profiles"},
		{"order_item", "OrderItem", "orderItem", "order_item", "/order-items"},
		{"product", "Product", "product", "product", "/products"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 测试业务逻辑模板
			content, err := renderer.RenderBiz(tc.pascal, "github.com/test/project")
			if err != nil {
				t.Fatalf("渲染失败: %v", err)
			}

			// 验证大驼峰命名
			if !strings.Contains(content, tc.pascal+"UseCase") {
				t.Errorf("缺少 %sUseCase", tc.pascal)
			}
			if !strings.Contains(content, tc.pascal+"Handler") {
				t.Errorf("缺少 %sHandler", tc.pascal)
			}
		})
	}
}

func TestTemplateConsistency(t *testing.T) {
	renderer, err := NewTemplateRenderer()
	if err != nil {
		t.Fatalf("创建模板渲染器失败: %v", err)
	}

	// 测试多次渲染的一致性
	content1, err1 := renderer.RenderBiz("User", "github.com/test/project")
	content2, err2 := renderer.RenderBiz("User", "github.com/test/project")

	if err1 != nil {
		t.Fatalf("第一次渲染失败: %v", err1)
	}
	if err2 != nil {
		t.Fatalf("第二次渲染失败: %v", err2)
	}

	if content1 != content2 {
		t.Error("多次渲染结果不一致")
	}
}
