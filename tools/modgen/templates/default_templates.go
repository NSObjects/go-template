/*
 * 默认模板渲染函数
 */

package templates

import (
	"fmt"
)

// RenderBiz 生成默认业务逻辑模板
func RenderBiz(pascal, packagePath string) string {
	tr, err := NewTemplateRenderer()
	if err != nil {
		// 如果模板渲染器创建失败，返回错误信息
		return fmt.Sprintf("// 错误: 无法创建模板渲染器: %v", err)
	}

	result, err := tr.RenderBiz(pascal, packagePath)
	if err != nil {
		// 如果模板渲染失败，返回错误信息
		return fmt.Sprintf("// 错误: 渲染biz模板失败: %v", err)
	}

	return result
}

// RenderService 生成默认服务层模板
func RenderService(pascal, camel, baseRoute, packagePath string) string {
	tr, err := NewTemplateRenderer()
	if err != nil {
		// 如果模板渲染器创建失败，返回错误信息
		return fmt.Sprintf("// 错误: 无法创建模板渲染器: %v", err)
	}

	result, err := tr.RenderService(pascal, camel, baseRoute, packagePath)
	if err != nil {
		// 如果模板渲染失败，返回错误信息
		return fmt.Sprintf("// 错误: 渲染service模板失败: %v", err)
	}

	return result
}

// RenderParam 生成默认参数模板
func RenderParam(pascal, packagePath string) string {
	tr, err := NewTemplateRenderer()
	if err != nil {
		// 如果模板渲染器创建失败，返回错误信息
		return fmt.Sprintf("// 错误: 无法创建模板渲染器: %v", err)
	}

	result, err := tr.RenderParam(pascal, packagePath)
	if err != nil {
		// 如果模板渲染失败，返回错误信息
		return fmt.Sprintf("// 错误: 渲染param模板失败: %v", err)
	}

	return result
}

// RenderModel 生成默认数据模型模板
func RenderModel(pascal, table, packagePath string) string {
	tr, err := NewTemplateRenderer()
	if err != nil {
		// 如果模板渲染器创建失败，返回错误信息
		return fmt.Sprintf("// 错误: 无法创建模板渲染器: %v", err)
	}

	result, err := tr.RenderModel(pascal, table, packagePath)
	if err != nil {
		// 如果模板渲染失败，返回错误信息
		return fmt.Sprintf("// 错误: 渲染model模板失败: %v", err)
	}

	return result
}

// RenderCode 生成业务错误码文件
func RenderCode(pascal, table, packagePath string) string {
	tr, err := NewTemplateRenderer()
	if err != nil {
		// 如果模板渲染器创建失败，返回错误信息
		return fmt.Sprintf("// 错误: 无法创建模板渲染器: %v", err)
	}

	result, err := tr.RenderCode(pascal, table, packagePath)
	if err != nil {
		// 如果模板渲染失败，返回错误信息
		return fmt.Sprintf("// 错误: 渲染code模板失败: %v", err)
	}

	return result
}
