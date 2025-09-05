/*
 * 代码生成器核心逻辑
 */

package generator

import (
	"fmt"
	"path/filepath"

	"github.com/NSObjects/go-template/tools/modgen/openapi"
	"github.com/NSObjects/go-template/tools/modgen/templates"
	"github.com/NSObjects/go-template/tools/modgen/utils"
)

// Config 生成器配置
type Config struct {
	Name          string
	Route         string
	Force         bool
	OpenAPIFile   string
	GenerateTests bool
	PackagePath   string
	RepoRoot      string
}

// Generator 代码生成器
type Generator struct {
	config *Config
}

// NewGenerator 创建新的代码生成器
func NewGenerator(config *Config) *Generator {
	return &Generator{
		config: config,
	}
}

// Generate 执行代码生成
func (g *Generator) Generate() error {
	// 检查模块名
	if g.config.Name == "" {
		return fmt.Errorf("模块名不能为空")
	}

	pascal := utils.ToPascal(g.config.Name)
	camel := utils.ToCamel(g.config.Name)
	plural := utils.Pluralize(g.config.Name)
	baseRoute := g.config.Route
	if baseRoute == "" {
		baseRoute = "/" + plural
	}

	utils.PrintInfo("🚀 开始生成 %s 模块...", g.config.Name)

	// 根据是否提供OpenAPI文档选择生成方式
	if g.config.OpenAPIFile != "" {
		utils.PrintInfo("📄 从OpenAPI3文档生成: %s", g.config.OpenAPIFile)
		return g.generateFromOpenAPIDoc(pascal, camel, baseRoute)
	} else {
		utils.PrintInfo("📝 使用默认模板生成")
		return g.generateFromDefaultTemplate(pascal, camel, baseRoute)
	}
}

// generateFromDefaultTemplate 使用默认模板生成
func (g *Generator) generateFromDefaultTemplate(pascal, camel, baseRoute string) error {
	// 创建模板渲染器
	renderer, err := templates.NewTemplateRenderer()
	if err != nil {
		return fmt.Errorf("创建模板渲染器失败: %v", err)
	}

	// 生成目标文件路径
	bizFile := filepath.Join(g.config.RepoRoot, "internal", "api", "biz", fmt.Sprintf("%s.go", g.config.Name))
	svcFile := filepath.Join(g.config.RepoRoot, "internal", "api", "service", fmt.Sprintf("%s.go", g.config.Name))
	paramFile := filepath.Join(g.config.RepoRoot, "internal", "api", "service", "param", fmt.Sprintf("%s.go", g.config.Name))
	codeFile := filepath.Join(g.config.RepoRoot, "internal", "code", fmt.Sprintf("%s.go", g.config.Name))

	// 生成模板内容
	bizContent, err := renderer.RenderBiz(pascal, g.config.PackagePath)
	if err != nil {
		return fmt.Errorf("渲染业务逻辑模板失败: %v", err)
	}

	svcContent, err := renderer.RenderService(pascal, camel, baseRoute, g.config.PackagePath)
	if err != nil {
		return fmt.Errorf("渲染服务层模板失败: %v", err)
	}

	paramContent, err := renderer.RenderParam(pascal, g.config.PackagePath)
	if err != nil {
		return fmt.Errorf("渲染参数模板失败: %v", err)
	}

	codeContent, err := renderer.RenderCode(pascal, g.config.Name, g.config.PackagePath)
	if err != nil {
		return fmt.Errorf("渲染错误码模板失败: %v", err)
	}

	// 写入文件
	utils.MustWrite(bizFile, bizContent, g.config.Force)
	utils.MustWrite(svcFile, svcContent, g.config.Force)
	utils.MustWrite(paramFile, paramContent, g.config.Force)
	utils.MustWrite(codeFile, codeContent, g.config.Force)

	// 生成测试用例（如果启用）
	if g.config.GenerateTests {
		utils.PrintInfo("🧪 生成测试用例...")
		bizTestFile := filepath.Join(g.config.RepoRoot, "internal", "api", "biz", fmt.Sprintf("%s_test.go", g.config.Name))
		svcTestFile := filepath.Join(g.config.RepoRoot, "internal", "api", "service", fmt.Sprintf("%s_test.go", g.config.Name))

		bizTestContent, err := renderer.RenderBizTest(pascal, g.config.PackagePath)
		if err != nil {
			return fmt.Errorf("渲染业务逻辑测试模板失败: %v", err)
		}

		svcTestContent, err := renderer.RenderServiceTest(pascal, camel, baseRoute, g.config.PackagePath)
		if err != nil {
			return fmt.Errorf("渲染服务层测试模板失败: %v", err)
		}

		utils.MustWrite(bizTestFile, bizTestContent, g.config.Force)
		utils.MustWrite(svcTestFile, svcTestContent, g.config.Force)
	}

	// 注入到 fx.Options
	_ = utils.TryInject(filepath.Join(g.config.RepoRoot, "internal", "api", "biz", "biz.go"), "New"+pascal+"Handler")
	_ = utils.TryInject(filepath.Join(g.config.RepoRoot, "internal", "api", "service", "service.go"), "AsRoute(New"+pascal+"Controller)")

	utils.PrintSuccess("✅ %s 模块生成完成！", g.config.Name)
	g.printGeneratedFiles(pascal)
	utils.PrintUsageInstructions(g.config.Name, pascal)

	return nil
}

// generateFromOpenAPIDoc 从OpenAPI3文档生成
func (g *Generator) generateFromOpenAPIDoc(pascal, camel, baseRoute string) error {
	// 解析OpenAPI文档
	openapiDoc, err := openapi.ParseOpenAPI3(g.config.OpenAPIFile)
	if err != nil {
		return fmt.Errorf("解析OpenAPI文档失败: %v", err)
	}

	// 生成API模块
	module, err := openapi.GenerateFromOpenAPI(openapiDoc, g.config.Name)
	if err != nil {
		return fmt.Errorf("生成API模块失败: %v", err)
	}

	// 生成文件路径
	bizFile := filepath.Join(g.config.RepoRoot, "internal", "api", "biz", fmt.Sprintf("%s.go", g.config.Name))
	svcFile := filepath.Join(g.config.RepoRoot, "internal", "api", "service", fmt.Sprintf("%s.go", g.config.Name))
	paramFile := filepath.Join(g.config.RepoRoot, "internal", "api", "service", "param", fmt.Sprintf("%s.go", g.config.Name))
	codeFile := filepath.Join(g.config.RepoRoot, "internal", "code", fmt.Sprintf("%s.go", g.config.Name))

	// 生成代码
	utils.MustWrite(bizFile, templates.RenderBizFromOpenAPI(module, pascal, g.config.PackagePath), g.config.Force)
	utils.MustWrite(svcFile, templates.RenderServiceFromOpenAPI(module, pascal, camel, baseRoute, g.config.PackagePath), g.config.Force)
	utils.MustWrite(paramFile, templates.RenderParamFromOpenAPI(module, pascal, g.config.PackagePath), g.config.Force)
	utils.MustWrite(codeFile, templates.RenderCode(pascal, g.config.Name, g.config.PackagePath), g.config.Force)

	// 生成测试用例（如果启用）
	if g.config.GenerateTests {
		utils.PrintInfo("🧪 生成测试用例...")
		bizTestFile := filepath.Join(g.config.RepoRoot, "internal", "api", "biz", fmt.Sprintf("%s_test.go", g.config.Name))
		svcTestFile := filepath.Join(g.config.RepoRoot, "internal", "api", "service", fmt.Sprintf("%s_test.go", g.config.Name))
		utils.MustWrite(bizTestFile, templates.RenderBizTestFromOpenAPI(module, pascal, g.config.PackagePath), g.config.Force)
		utils.MustWrite(svcTestFile, templates.RenderServiceTestFromOpenAPI(module, pascal, g.config.PackagePath), g.config.Force)
	}

	utils.PrintInfo("📊 从OpenAPI文档解析到 %d 个操作", len(module.Operations))

	// 注入到 fx.Options
	_ = utils.TryInject(filepath.Join(g.config.RepoRoot, "internal", "api", "biz", "biz.go"), "New"+pascal+"Handler")
	_ = utils.TryInject(filepath.Join(g.config.RepoRoot, "internal", "api", "service", "service.go"), "AsRoute(New"+pascal+"Controller)")

	utils.PrintSuccess("✅ %s 模块生成完成！", g.config.Name)
	g.printGeneratedFiles(pascal)
	utils.PrintUsageInstructions(g.config.Name, pascal)

	return nil
}

// printGeneratedFiles 打印生成的文件列表
func (g *Generator) printGeneratedFiles(pascal string) {
	fmt.Printf("\n📁 生成的文件:\n")
	fmt.Printf("  📄 业务逻辑: %s\n", filepath.Join(g.config.RepoRoot, "internal", "api", "biz", fmt.Sprintf("%s.go", g.config.Name)))
	fmt.Printf("  📄 控制器: %s\n", filepath.Join(g.config.RepoRoot, "internal", "api", "service", fmt.Sprintf("%s.go", g.config.Name)))
	fmt.Printf("  📄 参数结构: %s\n", filepath.Join(g.config.RepoRoot, "internal", "api", "service", "param", fmt.Sprintf("%s.go", g.config.Name)))
	fmt.Printf("  📄 错误码: %s\n", filepath.Join(g.config.RepoRoot, "internal", "code", fmt.Sprintf("%s.go", g.config.Name)))
}
