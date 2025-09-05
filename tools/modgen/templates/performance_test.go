/*
 * 性能测试
 */

package templates

import (
	"testing"
)

func BenchmarkTemplateRendererCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		renderer, err := NewTemplateRenderer()
		if err != nil {
			b.Fatalf("创建模板渲染器失败: %v", err)
		}
		_ = renderer
	}
}

func BenchmarkRenderBiz(b *testing.B) {
	renderer, err := NewTemplateRenderer()
	if err != nil {
		b.Fatalf("创建模板渲染器失败: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := renderer.RenderBiz("User", "github.com/test/project")
		if err != nil {
			b.Fatalf("渲染失败: %v", err)
		}
	}
}

func BenchmarkRenderService(b *testing.B) {
	renderer, err := NewTemplateRenderer()
	if err != nil {
		b.Fatalf("创建模板渲染器失败: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := renderer.RenderService("User", "user", "/users", "github.com/test/project")
		if err != nil {
			b.Fatalf("渲染失败: %v", err)
		}
	}
}

func BenchmarkRenderParam(b *testing.B) {
	renderer, err := NewTemplateRenderer()
	if err != nil {
		b.Fatalf("创建模板渲染器失败: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := renderer.RenderParam("User", "github.com/test/project")
		if err != nil {
			b.Fatalf("渲染失败: %v", err)
		}
	}
}

func BenchmarkRenderModel(b *testing.B) {
	renderer, err := NewTemplateRenderer()
	if err != nil {
		b.Fatalf("创建模板渲染器失败: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := renderer.RenderModel("User", "users", "github.com/test/project")
		if err != nil {
			b.Fatalf("渲染失败: %v", err)
		}
	}
}

func BenchmarkRenderCode(b *testing.B) {
	renderer, err := NewTemplateRenderer()
	if err != nil {
		b.Fatalf("创建模板渲染器失败: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := renderer.RenderCode("User", "users", "github.com/test/project")
		if err != nil {
			b.Fatalf("渲染失败: %v", err)
		}
	}
}

func BenchmarkRenderAllTemplates(b *testing.B) {
	renderer, err := NewTemplateRenderer()
	if err != nil {
		b.Fatalf("创建模板渲染器失败: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 渲染所有模板
		_, err := renderer.RenderBiz("User", "github.com/test/project")
		if err != nil {
			b.Fatalf("渲染Biz失败: %v", err)
		}

		_, err = renderer.RenderService("User", "user", "/users", "github.com/test/project")
		if err != nil {
			b.Fatalf("渲染Service失败: %v", err)
		}

		_, err = renderer.RenderParam("User", "github.com/test/project")
		if err != nil {
			b.Fatalf("渲染Param失败: %v", err)
		}

		_, err = renderer.RenderModel("User", "users", "github.com/test/project")
		if err != nil {
			b.Fatalf("渲染Model失败: %v", err)
		}

		_, err = renderer.RenderCode("User", "users", "github.com/test/project")
		if err != nil {
			b.Fatalf("渲染Code失败: %v", err)
		}
	}
}

func BenchmarkRenderWithDifferentNames(b *testing.B) {
	renderer, err := NewTemplateRenderer()
	if err != nil {
		b.Fatalf("创建模板渲染器失败: %v", err)
	}

	names := []string{
		"User", "Product", "Order", "Category", "Item",
		"UserProfile", "OrderItem", "ProductCategory", "TestModule", "ApiEndpoint",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		name := names[i%len(names)]
		_, err := renderer.RenderBiz(name, "github.com/test/project")
		if err != nil {
			b.Fatalf("渲染失败: %v", err)
		}
	}
}

func BenchmarkRenderWithLongNames(b *testing.B) {
	renderer, err := NewTemplateRenderer()
	if err != nil {
		b.Fatalf("创建模板渲染器失败: %v", err)
	}

	longName := "VeryLongModuleNameThatExceedsNormalLengthAndShouldBeHandledProperly"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := renderer.RenderBiz(longName, "github.com/test/project")
		if err != nil {
			b.Fatalf("渲染失败: %v", err)
		}
	}
}

func BenchmarkRenderWithLongPackagePath(b *testing.B) {
	renderer, err := NewTemplateRenderer()
	if err != nil {
		b.Fatalf("创建模板渲染器失败: %v", err)
	}

	longPackagePath := "github.com/very/long/package/path/that/exceeds/normal/length/and/should/be/handled/properly"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := renderer.RenderBiz("User", longPackagePath)
		if err != nil {
			b.Fatalf("渲染失败: %v", err)
		}
	}
}

func BenchmarkConcurrentRendering(b *testing.B) {
	renderer, err := NewTemplateRenderer()
	if err != nil {
		b.Fatalf("创建模板渲染器失败: %v", err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := renderer.RenderBiz("User", "github.com/test/project")
			if err != nil {
				b.Fatalf("并发渲染失败: %v", err)
			}
		}
	})
}

func BenchmarkMemoryUsage(b *testing.B) {
	renderer, err := NewTemplateRenderer()
	if err != nil {
		b.Fatalf("创建模板渲染器失败: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		content, err := renderer.RenderBiz("User", "github.com/test/project")
		if err != nil {
			b.Fatalf("渲染失败: %v", err)
		}
		_ = content // 防止编译器优化
	}
}
