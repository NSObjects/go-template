/*
 * 基于 text/template 的模板渲染器
 */

package templates

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/NSObjects/go-template/muban/modgen/openapi"
)

//go:embed tmpl/*.tmpl
var embeddedTemplates embed.FS

// TemplateData 模板数据
type TemplateData struct {
	Pascal                string                 // 大驼峰命名
	Camel                 string                 // 小驼峰命名
	Table                 string                 // 表名
	Route                 string                 // 基础路由
	PackagePath           string                 // 包路径
	BaseCode              int                    // 错误码基础值
	Operations            []openapi.APIOperation // OpenAPI 操作列表
	ResponseDataTypes     []openapi.ResponseData // 去重的响应数据类型
	HasTimeFields         bool                   // 是否包含时间字段
	HasPathParams         bool                   // 是否包含路径参数
	HasRequestBodyOrQuery bool                   // 是否包含请求体或查询参数
	ErrorCodes            []openapi.ErrorCode    // 错误码列表
}

// TemplateRenderer 模板渲染器
type TemplateRenderer struct {
	templates map[string]*template.Template
}

var (
	rendererOnce sync.Once
	renderer     *TemplateRenderer
	rendererErr  error
)

// NewTemplateRenderer 创建新的模板渲染器
func NewTemplateRenderer() (*TemplateRenderer, error) {
	rendererOnce.Do(func() {
		renderer, rendererErr = loadRenderer()
	})
	return renderer, rendererErr
}

func loadRenderer() (*TemplateRenderer, error) {
	tr := &TemplateRenderer{templates: make(map[string]*template.Template)}

	funcMap := template.FuncMap{
		"hasPrefix": strings.HasPrefix,
		"contains":  strings.Contains,
	}

	templateNames := []string{
		"biz",
		"service",
		"param",
		"model",
		"code",
		"param_openapi",
		"biz_openapi",
		"service_openapi",
		"service_test_enhanced",
		"biz_test_enhanced",
		"code_openapi",
		"context_support",
		"resp_support",
	}

	if templateDir, err := locateTemplateDir(); err == nil {
		if err := tr.loadFromDirectory(templateDir, templateNames, funcMap); err == nil {
			return tr, nil
		}
	}

	if err := tr.loadFromEmbeddedFS(templateNames, funcMap); err != nil {
		return nil, err
	}

	return tr, nil
}

func (tr *TemplateRenderer) loadFromDirectory(templateDir string, templateNames []string, funcMap template.FuncMap) error {
	for _, name := range templateNames {
		templatePath := filepath.Join(templateDir, name+".tmpl")

		tmpl, err := template.New(filepath.Base(templatePath)).Funcs(funcMap).ParseFiles(templatePath)
		if err != nil {
			return fmt.Errorf("加载模板 %s 失败，路径: %s, 错误: %w", name, templatePath, err)
		}
		tr.templates[name] = tmpl
	}

	return nil
}

func (tr *TemplateRenderer) loadFromEmbeddedFS(templateNames []string, funcMap template.FuncMap) error {
	for _, name := range templateNames {
		content, err := fs.ReadFile(embeddedTemplates, fmt.Sprintf("tmpl/%s.tmpl", name))
		if err != nil {
			return fmt.Errorf("加载嵌入模板 %s 失败: %w", name, err)
		}

		tmpl, err := template.New(name + ".tmpl").Funcs(funcMap).Parse(string(content))
		if err != nil {
			return fmt.Errorf("解析嵌入模板 %s 失败: %w", name, err)
		}

		tr.templates[name] = tmpl
	}

	return nil
}

func locateTemplateDir() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("无法获取工作目录: %w", err)
	}

	candidates := collectTemplateCandidates(wd)
	for _, path := range candidates {
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			return path, nil
		}
	}

	return "", fmt.Errorf("无法找到模板目录，尝试的路径: %v", candidates)
}

func collectTemplateCandidates(wd string) []string {
	seen := make(map[string]struct{})
	var paths []string

	add := func(path string) {
		if _, ok := seen[path]; ok {
			return
		}
		seen[path] = struct{}{}
		paths = append(paths, path)
	}

	relPaths := []string{
		"tmpl",
		filepath.Join("templates", "tmpl"),
		filepath.Join("muban", "modgen", "templates", "tmpl"),
	}

	for dir := wd; ; dir = filepath.Dir(dir) {
		for _, rel := range relPaths {
			add(filepath.Join(dir, rel))
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
	}

	// 兼容相对路径引用
	add("tmpl")

	return paths
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

// RenderContextSupport 生成 utils 上下文支持文件
func (tr *TemplateRenderer) RenderContextSupport() (string, error) {
	return tr.Render("context_support", TemplateData{})
}

// RenderRespSupport 生成 resp 辅助函数文件
func (tr *TemplateRenderer) RenderRespSupport() (string, error) {
	return tr.Render("resp_support", TemplateData{})
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

func (tr *TemplateRenderer) RenderServiceTestEnhanced(data TemplateData) (string, error) {
	return tr.Render("service_test_enhanced", data)
}

func (tr *TemplateRenderer) RenderBizTestEnhanced(data TemplateData) (string, error) {
	return tr.Render("biz_test_enhanced", data)
}
