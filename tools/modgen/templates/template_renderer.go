/*
 * 基于 text/template 的模板渲染器
 */

package templates

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/NSObjects/go-template/tools/modgen/openapi"
)

// TemplateData 模板数据
type TemplateData struct {
	Pascal            string                 // 大驼峰命名
	Camel             string                 // 小驼峰命名
	Table             string                 // 表名
	Route             string                 // 基础路由
	PackagePath       string                 // 包路径
	BaseCode          int                    // 错误码基础值
	Operations        []openapi.APIOperation // OpenAPI 操作列表
	ResponseDataTypes []openapi.ResponseData // 去重的响应数据类型
}

// TemplateRenderer 模板渲染器
type TemplateRenderer struct {
	templates map[string]*template.Template
}

// NewTemplateRenderer 创建新的模板渲染器
func NewTemplateRenderer() (*TemplateRenderer, error) {
	tr := &TemplateRenderer{
		templates: make(map[string]*template.Template),
	}

	// 加载所有模板
	templates := []string{"biz", "service", "param", "model", "code", "biz_test", "service_test", "param_openapi", "service_test_openapi", "biz_openapi", "service_openapi", "biz_test_openapi", "service_test_openapi"}
	for _, name := range templates {
		// 尝试多个可能的路径
		paths := []string{
			filepath.Join("tools", "modgen", "templates", "tmpl", name+".tmpl"),
			filepath.Join("tmpl", name+".tmpl"),
			name + ".tmpl",
			filepath.Join("..", "templates", "tmpl", name+".tmpl"),
			filepath.Join("..", "..", "templates", "tmpl", name+".tmpl"),
		}

		// 添加绝对路径
		if wd, err := os.Getwd(); err == nil {
			absPath := filepath.Join(wd, "tools", "modgen", "templates", "tmpl", name+".tmpl")
			paths = append([]string{absPath}, paths...)
		}

		var tmpl *template.Template
		var err error
		for _, path := range paths {
			tmpl, err = template.ParseFiles(path)
			if err == nil {
				break
			}
		}

		if err != nil {
			return nil, fmt.Errorf("加载模板 %s 失败，尝试的路径: %v, 错误: %v", name, paths, err)
		}
		tr.templates[name] = tmpl
	}

	return tr, nil
}

// Render 渲染模板
func (tr *TemplateRenderer) Render(templateName string, data TemplateData) (string, error) {
	tmpl, exists := tr.templates[templateName]
	if !exists {
		return "", fmt.Errorf("模板 %s 不存在", templateName)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("渲染模板 %s 失败: %v", templateName, err)
	}

	return buf.String(), nil
}

// RenderBiz 生成业务逻辑模板
func (tr *TemplateRenderer) RenderBiz(pascal, packagePath string) (string, error) {
	data := TemplateData{
		Pascal:      pascal,
		PackagePath: packagePath,
	}
	return tr.Render("biz", data)
}

// RenderService 生成服务层模板
func (tr *TemplateRenderer) RenderService(pascal, camel, baseRoute, packagePath string) (string, error) {
	data := TemplateData{
		Pascal:      pascal,
		Camel:       camel,
		Route:       baseRoute,
		PackagePath: packagePath,
	}
	return tr.Render("service", data)
}

// RenderParam 生成参数模板
func (tr *TemplateRenderer) RenderParam(pascal, packagePath string) (string, error) {
	data := TemplateData{
		Pascal:      pascal,
		PackagePath: packagePath,
	}
	return tr.Render("param", data)
}

// RenderModel 生成数据模型模板
func (tr *TemplateRenderer) RenderModel(pascal, table, packagePath string) (string, error) {
	data := TemplateData{
		Pascal:      pascal,
		Table:       table,
		PackagePath: packagePath,
	}
	return tr.Render("model", data)
}

// RenderCode 生成错误码模板
func (tr *TemplateRenderer) RenderCode(pascal, table, packagePath string) (string, error) {
	// 计算错误码起始值（基于表名）
	baseCode := 100000
	if len(table) > 0 {
		baseCode += int(table[0])*1000 + int(table[len(table)-1])*10
	}

	data := TemplateData{
		Pascal:      pascal,
		Table:       table,
		PackagePath: packagePath,
		BaseCode:    baseCode,
	}
	return tr.Render("code", data)
}

// RenderBizTest 生成业务逻辑测试模板
func (tr *TemplateRenderer) RenderBizTest(pascal, packagePath string) (string, error) {
	data := TemplateData{
		Pascal:      pascal,
		PackagePath: packagePath,
	}
	return tr.Render("biz_test", data)
}

// RenderServiceTest 生成服务层测试模板
func (tr *TemplateRenderer) RenderServiceTest(pascal, camel, route, packagePath string) (string, error) {
	data := TemplateData{
		Pascal:      pascal,
		Camel:       camel,
		Route:       route,
		PackagePath: packagePath,
	}
	return tr.Render("service_test", data)
}
