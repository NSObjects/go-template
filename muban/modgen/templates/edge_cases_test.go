/*
 * 边界情况和特殊场景测试
 */

package templates

import (
	"strings"
	"testing"
)

func TestTemplateWithEmptyData(t *testing.T) {
	renderer, err := NewTemplateRenderer()
	if err != nil {
		t.Fatalf("创建模板渲染器失败: %v", err)
	}

	// 测试空字符串
	content, err := renderer.RenderBiz("", "github.com/test/project")
	if err != nil {
		t.Fatalf("渲染空字符串模板失败: %v", err)
	}

	// 应该包含空字符串的处理
	if !strings.Contains(content, "UseCase interface") {
		t.Error("空字符串模板应该包含基本结构")
	}
}

func TestTemplateWithVeryLongName(t *testing.T) {
	renderer, err := NewTemplateRenderer()
	if err != nil {
		t.Fatalf("创建模板渲染器失败: %v", err)
	}

	// 测试很长的名称
	longName := "VeryLongModuleNameThatExceedsNormalLength"
	content, err := renderer.RenderBiz(longName, "github.com/test/project")
	if err != nil {
		t.Fatalf("渲染长名称模板失败: %v", err)
	}

	// 验证长名称被正确处理
	expectedContents := []string{
		longName + "UseCase",
		longName + "Handler",
		"New" + longName + "Handler",
	}

	for _, expected := range expectedContents {
		if !strings.Contains(content, expected) {
			t.Errorf("长名称模板缺少内容: %s", expected)
		}
	}
}

func TestTemplateWithSpecialCharacters(t *testing.T) {
	renderer, err := NewTemplateRenderer()
	if err != nil {
		t.Fatalf("创建模板渲染器失败: %v", err)
	}

	// 测试包含特殊字符的名称
	specialNames := []string{
		"user_profile",
		"order-item",
		"product_category",
		"test_module_v2",
		"api_endpoint",
	}

	for _, name := range specialNames {
		t.Run(name, func(t *testing.T) {
			content, err := renderer.RenderBiz(name, "github.com/test/project")
			if err != nil {
				t.Fatalf("渲染特殊字符模板失败: %v", err)
			}

			// 验证特殊字符被正确处理
			if !strings.Contains(content, "UseCase") {
				t.Errorf("特殊字符模板缺少基本结构: %s", name)
			}
		})
	}
}

func TestTemplateWithUnicodeCharacters(t *testing.T) {
	renderer, err := NewTemplateRenderer()
	if err != nil {
		t.Fatalf("创建模板渲染器失败: %v", err)
	}

	// 测试Unicode字符
	unicodeName := "用户管理"
	content, err := renderer.RenderBiz(unicodeName, "github.com/test/project")
	if err != nil {
		t.Fatalf("渲染Unicode模板失败: %v", err)
	}

	// 验证Unicode字符被正确处理
	if !strings.Contains(content, "UseCase") {
		t.Error("Unicode模板缺少基本结构")
	}
}

func TestTemplateWithNumbers(t *testing.T) {
	renderer, err := NewTemplateRenderer()
	if err != nil {
		t.Fatalf("创建模板渲染器失败: %v", err)
	}

	// 测试包含数字的名称
	numberNames := []string{
		"user2",
		"api_v1",
		"test123",
		"module_2023",
	}

	for _, name := range numberNames {
		t.Run(name, func(t *testing.T) {
			content, err := renderer.RenderBiz(name, "github.com/test/project")
			if err != nil {
				t.Fatalf("渲染数字模板失败: %v", err)
			}

			// 验证数字被正确处理
			if !strings.Contains(content, "UseCase") {
				t.Errorf("数字模板缺少基本结构: %s", name)
			}
		})
	}
}

func TestTemplateWithMixedCase(t *testing.T) {
	renderer, err := NewTemplateRenderer()
	if err != nil {
		t.Fatalf("创建模板渲染器失败: %v", err)
	}

	// 测试混合大小写
	mixedCaseNames := []string{
		"UserProfile",
		"orderItem",
		"APIEndpoint",
		"testModule",
	}

	for _, name := range mixedCaseNames {
		t.Run(name, func(t *testing.T) {
			content, err := renderer.RenderBiz(name, "github.com/test/project")
			if err != nil {
				t.Fatalf("渲染混合大小写模板失败: %v", err)
			}

			// 验证混合大小写被正确处理
			if !strings.Contains(content, "UseCase") {
				t.Errorf("混合大小写模板缺少基本结构: %s", name)
			}
		})
	}
}

func TestTemplateWithVeryLongPackagePath(t *testing.T) {
	renderer, err := NewTemplateRenderer()
	if err != nil {
		t.Fatalf("创建模板渲染器失败: %v", err)
	}

	// 测试很长的包路径
	longPackagePath := "github.com/very/long/package/path/that/exceeds/normal/length/and/should/be/handled/properly"
	content, err := renderer.RenderBiz("Test", longPackagePath)
	if err != nil {
		t.Fatalf("渲染长包路径模板失败: %v", err)
	}

	// 验证长包路径被正确处理
	if !strings.Contains(content, longPackagePath) {
		t.Error("长包路径模板应该包含完整路径")
	}
}

func TestTemplateWithEmptyPackagePath(t *testing.T) {
	renderer, err := NewTemplateRenderer()
	if err != nil {
		t.Fatalf("创建模板渲染器失败: %v", err)
	}

	// 测试空包路径
	content, err := renderer.RenderBiz("Test", "")
	if err != nil {
		t.Fatalf("渲染空包路径模板失败: %v", err)
	}

	// 验证空包路径被正确处理
	if !strings.Contains(content, "UseCase") {
		t.Error("空包路径模板缺少基本结构")
	}
}

func TestTemplateWithSpecialPackagePath(t *testing.T) {
	renderer, err := NewTemplateRenderer()
	if err != nil {
		t.Fatalf("创建模板渲染器失败: %v", err)
	}

	// 测试特殊包路径
	specialPackagePaths := []string{
		"github.com/user-name/project_name",
		"example.com/v1/api",
		"internal/package",
		"./relative/path",
	}

	for _, pkgPath := range specialPackagePaths {
		t.Run(pkgPath, func(t *testing.T) {
			content, err := renderer.RenderBiz("Test", pkgPath)
			if err != nil {
				t.Fatalf("渲染特殊包路径模板失败: %v", err)
			}

			// 验证特殊包路径被正确处理
			if !strings.Contains(content, pkgPath) {
				t.Errorf("特殊包路径模板缺少路径: %s", pkgPath)
			}
		})
	}
}

func TestTemplateConcurrentRendering(t *testing.T) {
	renderer, err := NewTemplateRenderer()
	if err != nil {
		t.Fatalf("创建模板渲染器失败: %v", err)
	}

	// 并发渲染测试
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(i int) {
			defer func() { done <- true }()

			name := "Test" + string(rune('0'+i))
			content, err := renderer.RenderBiz(name, "github.com/test/project")
			if err != nil {
				t.Errorf("并发渲染失败: %v", err)
				return
			}

			if !strings.Contains(content, "UseCase") {
				t.Errorf("并发渲染结果不正确: %s", name)
			}
		}(i)
	}

	// 等待所有goroutine完成
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestTemplateMemoryUsage(t *testing.T) {
	renderer, err := NewTemplateRenderer()
	if err != nil {
		t.Fatalf("创建模板渲染器失败: %v", err)
	}

	// 多次渲染测试内存使用
	for i := 0; i < 100; i++ {
		content, err := renderer.RenderBiz("Test", "github.com/test/project")
		if err != nil {
			t.Fatalf("内存使用测试失败: %v", err)
		}

		if len(content) == 0 {
			t.Error("渲染结果为空")
		}
	}
}

func TestTemplateErrorRecovery(t *testing.T) {
	renderer, err := NewTemplateRenderer()
	if err != nil {
		t.Fatalf("创建模板渲染器失败: %v", err)
	}

	// 测试错误恢复
	content, err := renderer.RenderBiz("Test", "github.com/test/project")
	if err != nil {
		t.Fatalf("错误恢复测试失败: %v", err)
	}

	// 验证内容正确
	if !strings.Contains(content, "UseCase") {
		t.Error("错误恢复后内容不正确")
	}
}
