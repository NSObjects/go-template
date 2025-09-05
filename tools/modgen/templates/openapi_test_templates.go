/*
 * OpenAPI3 测试模板渲染函数
 */

package templates

import (
	"fmt"

	"github.com/NSObjects/go-template/tools/modgen/openapi"
)

// RenderBizTestFromOpenAPI 从OpenAPI3生成业务逻辑测试模板
func RenderBizTestFromOpenAPI(module *openapi.APIModule, pascal, packagePath string) string {
	// 使用新的模板渲染器
	renderer, err := NewTemplateRenderer()
	if err != nil {
		return fmt.Sprintf("// 错误: %v", err)
	}

	// 准备模板数据
	data := TemplateData{
		Pascal:      pascal,
		PackagePath: packagePath,
		Operations:  module.Operations,
	}

	// 渲染模板
	content, err := renderer.Render("biz_test_openapi", data)
	if err != nil {
		return fmt.Sprintf("// 错误: %v", err)
	}

	return content
}

// RenderServiceTestFromOpenAPI 从OpenAPI3生成服务层测试模板
func RenderServiceTestFromOpenAPI(module *openapi.APIModule, pascal, packagePath string) string {
	// 使用新的模板渲染器
	renderer, err := NewTemplateRenderer()
	if err != nil {
		return fmt.Sprintf("// 错误: %v", err)
	}

	// 准备模板数据
	data := TemplateData{
		Pascal:      pascal,
		PackagePath: packagePath,
		Operations:  module.Operations,
	}

	// 渲染模板
	content, err := renderer.Render("service_test_openapi", data)
	if err != nil {
		return fmt.Sprintf("// 错误: %v", err)
	}

	return content
}
