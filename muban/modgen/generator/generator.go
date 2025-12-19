/*
 * 代码生成器核心逻辑
 */

package generator

import (
	"fmt"
	"path/filepath"

	"github.com/NSObjects/go-template/muban/modgen/openapi"
	"github.com/NSObjects/go-template/muban/modgen/templates"
	"github.com/NSObjects/go-template/muban/modgen/utils"
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
	GenerateAll   bool
}

// Generator 代码生成器
type Generator struct {
	config     *Config
	openAPIDoc *openapi.OpenAPI3
	renderer   *templates.TemplateRenderer
}

// NewGenerator 创建新的代码生成器
func NewGenerator(config *Config) *Generator {
	return &Generator{
		config: config,
	}
}

// Generate 执行代码生成
func (g *Generator) Generate() error {
	if err := g.ensureSupportFiles(); err != nil {
		return err
	}

	// 如果生成所有模块
	if g.config.GenerateAll {
		return g.generateAllModules()
	}

	// 检查模块名
	if g.config.Name == "" {
		return fmt.Errorf("模块名不能为空")
	}

	utils.PrintInfo("🚀 开始生成 %s 模块...", g.config.Name)

	// 根据是否提供OpenAPI文档选择生成方式
	if g.config.OpenAPIFile != "" {
		utils.PrintInfo("📄 从OpenAPI3文档生成: %s", g.config.OpenAPIFile)
		doc, err := g.loadOpenAPIDoc()
		if err != nil {
			return err
		}
		return g.generateOpenAPIModule(doc, g.config.Name, true)
	}

	utils.PrintInfo("📝 使用默认模板生成")
	return g.generateFromDefaultTemplate()
}

// generateFromDefaultTemplate 使用默认模板生成
func (g *Generator) generateFromDefaultTemplate() error {
	pascal, camel, baseRoute := g.naming(g.config.Name)

	renderer, err := g.templateRenderer()
	if err != nil {
		return fmt.Errorf("创建模板渲染器失败: %v", err)
	}

	// 生成目标文件路径
	bizFile, svcFile, paramFile, codeFile := g.moduleFilePaths(g.config.Name)
	dataFile := filepath.Join(g.config.RepoRoot, "internal", "api", "data", fmt.Sprintf("%s.go", g.config.Name))

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

	dataContent, err := renderer.RenderData(pascal, camel, g.config.PackagePath)
	if err != nil {
		return fmt.Errorf("渲染数据层模板失败: %v", err)
	}

	codeContent, err := renderer.RenderCode(pascal, g.config.Name, g.config.PackagePath)
	if err != nil {
		return fmt.Errorf("渲染错误码模板失败: %v", err)
	}

	// 写入文件
	utils.MustWrite(bizFile, bizContent, g.config.Force)
	utils.MustWrite(svcFile, svcContent, g.config.Force)
	utils.MustWrite(paramFile, paramContent, g.config.Force)
	utils.MustWrite(dataFile, dataContent, g.config.Force)
	utils.MustWrite(codeFile, codeContent, g.config.Force)

	// 生成测试用例（如果启用）
	if g.config.GenerateTests {
		utils.PrintInfo("🧪 生成测试用例...")
		bizTestFile, svcTestFile := g.moduleTestFilePaths(g.config.Name)

		// 使用增强测试模板作为默认
		data := templates.TemplateData{
			Pascal:      pascal,
			PackagePath: g.config.PackagePath,
		}
		bizTestContent, err := renderer.RenderBizTestEnhanced(data)
		if err != nil {
			return fmt.Errorf("渲染业务逻辑测试模板失败: %v", err)
		}

		svcTestData := templates.TemplateData{
			Pascal:      pascal,
			Camel:       camel,
			Route:       baseRoute,
			PackagePath: g.config.PackagePath,
		}
		svcTestContent, err := renderer.RenderServiceTestEnhanced(svcTestData)
		if err != nil {
			return fmt.Errorf("渲染服务层测试模板失败: %v", err)
		}

		utils.MustWrite(bizTestFile, bizTestContent, g.config.Force)
		utils.MustWrite(svcTestFile, svcTestContent, g.config.Force)
	}

	// 注入到 fx.Options
	_ = utils.TryInject(filepath.Join(g.config.RepoRoot, "internal", "api", "biz", "biz.go"), "New"+pascal+"Handler")
	_ = utils.TryInject(filepath.Join(g.config.RepoRoot, "internal", "api", "service", "service.go"), "AsRoute(New"+pascal+"Controller)")
	_ = utils.TryInject(filepath.Join(g.config.RepoRoot, "internal", "api", "data", "data.go"), "New"+pascal+"Repository")

	utils.PrintSuccess("✅ %s 模块生成完成！", g.config.Name)
	g.printGeneratedFiles(g.config.Name)
	utils.PrintUsageInstructions(g.config.Name, pascal)

	return nil
}

// generateAllModules 生成所有API模块
func (g *Generator) generateAllModules() error {
	utils.PrintInfo("🚀 开始生成所有API模块...")
	utils.PrintInfo("📄 从OpenAPI3文档生成: %s", g.config.OpenAPIFile)

	// 解析OpenAPI文档
	openapiDoc, err := g.loadOpenAPIDoc()
	if err != nil {
		return err
	}

	// 提取所有模块名
	moduleNames, err := openapi.ExtractAllModuleNames(openapiDoc)
	if err != nil {
		return fmt.Errorf("提取模块名失败: %v", err)
	}

	if len(moduleNames) == 0 {
		utils.PrintInfo("⚠️  OpenAPI文档中没有找到任何模块")
		return nil
	}

	utils.PrintInfo("📊 发现 %d 个模块: %v", len(moduleNames), moduleNames)

	// 为每个模块生成代码
	successCount := 0
	for _, moduleName := range moduleNames {
		utils.PrintInfo("\n🔄 正在生成模块: %s", moduleName)

		if err := g.generateOpenAPIModule(openapiDoc, moduleName, false); err != nil {
			utils.PrintError("❌ 生成模块 %s 失败: %v", moduleName, err)
			continue
		}

		successCount++
		utils.PrintSuccess("✅ 模块 %s 生成完成", moduleName)
	}

	utils.PrintSuccess("\n🎉 所有模块生成完成！成功生成 %d/%d 个模块", successCount, len(moduleNames))
	return nil
}

// generateOpenAPIModule 使用 OpenAPI 文档生成模块
func (g *Generator) generateOpenAPIModule(doc *openapi.OpenAPI3, moduleName string, showSummary bool) error {
	renderer, err := g.templateRenderer()
	if err != nil {
		return fmt.Errorf("创建模板渲染器失败: %w", err)
	}

	// 生成API模块
	module, err := openapi.GenerateFromOpenAPI(doc, moduleName)
	if err != nil {
		return fmt.Errorf("生成API模块失败: %w", err)
	}

	pascal, camel, baseRoute := g.naming(moduleName)
	bizFile, svcFile, paramFile, codeFile := g.moduleFilePaths(moduleName)
	dataFile := filepath.Join(g.config.RepoRoot, "internal", "api", "data", fmt.Sprintf("%s.go", moduleName))

	// 生成代码
	bizContent, err := renderer.RenderOpenAPIBiz(module, pascal, g.config.PackagePath)
	if err != nil {
		return fmt.Errorf("渲染业务逻辑模板失败: %w", err)
	}
	utils.MustWrite(bizFile, bizContent, g.config.Force)

	svcContent, err := renderer.RenderOpenAPIService(module, pascal, camel, baseRoute, g.config.PackagePath)
	if err != nil {
		return fmt.Errorf("渲染服务层模板失败: %w", err)
	}
	utils.MustWrite(svcFile, svcContent, g.config.Force)

	paramContent, err := renderer.RenderOpenAPIParam(module, pascal, g.config.PackagePath)
	if err != nil {
		return fmt.Errorf("渲染参数模板失败: %w", err)
	}
	utils.MustWrite(paramFile, paramContent, g.config.Force)

	codeContent, err := renderer.RenderOpenAPICode(module, pascal, g.config.PackagePath)
	if err != nil {
		return fmt.Errorf("渲染错误码模板失败: %w", err)
	}
	utils.MustWrite(codeFile, codeContent, g.config.Force)

	// data 层 Repository 实现骨架
	dataContent, err := renderer.RenderOpenAPIData(module, pascal, camel, g.config.PackagePath)
	if err != nil {
		return fmt.Errorf("渲染数据层模板失败: %w", err)
	}
	utils.MustWrite(dataFile, dataContent, g.config.Force)

	// 生成测试用例（如果启用）
	if g.config.GenerateTests {
		bizTestFile, svcTestFile := g.moduleTestFilePaths(moduleName)
		bizTestContent, err := renderer.RenderOpenAPIBizTests(module, pascal, camel, g.config.PackagePath)
		if err != nil {
			return fmt.Errorf("渲染业务逻辑测试模板失败: %w", err)
		}
		utils.MustWrite(bizTestFile, bizTestContent, g.config.Force)

		svcTestContent, err := renderer.RenderOpenAPIServiceTests(module, pascal, camel, g.config.PackagePath)
		if err != nil {
			return fmt.Errorf("渲染服务层测试模板失败: %w", err)
		}
		utils.MustWrite(svcTestFile, svcTestContent, g.config.Force)
	}

	// 注入到 fx.Options
	_ = utils.TryInject(filepath.Join(g.config.RepoRoot, "internal", "api", "biz", "biz.go"), "New"+pascal+"Handler")
	_ = utils.TryInject(filepath.Join(g.config.RepoRoot, "internal", "api", "service", "service.go"), "AsRoute(New"+pascal+"Controller)")
	_ = utils.TryInject(filepath.Join(g.config.RepoRoot, "internal", "api", "data", "data.go"), "New"+pascal+"Repository")

	if showSummary {
		utils.PrintInfo("📊 从OpenAPI文档解析到 %d 个操作", len(module.Operations))
		utils.PrintSuccess("✅ %s 模块生成完成！", moduleName)
		g.printGeneratedFiles(moduleName)
		utils.PrintUsageInstructions(moduleName, pascal)
	}

	return nil
}

// loadOpenAPIDoc 解析并缓存 OpenAPI 文档
func (g *Generator) loadOpenAPIDoc() (*openapi.OpenAPI3, error) {
	if g.config.OpenAPIFile == "" {
		return nil, fmt.Errorf("未指定OpenAPI文档路径")
	}

	if g.openAPIDoc != nil {
		return g.openAPIDoc, nil
	}

	doc, err := openapi.ParseOpenAPI3(g.config.OpenAPIFile)
	if err != nil {
		return nil, fmt.Errorf("解析OpenAPI文档失败: %w", err)
	}

	g.openAPIDoc = doc
	return g.openAPIDoc, nil
}

// naming 生成命名相关信息
func (g *Generator) naming(moduleName string) (pascal, camel, baseRoute string) {
	pascal = utils.ToPascal(moduleName)
	camel = utils.ToCamel(moduleName)
	baseRoute = g.config.Route
	if baseRoute == "" {
		baseRoute = "/" + utils.Pluralize(moduleName)
	}
	return
}

// moduleFilePaths 返回模块相关文件路径
func (g *Generator) moduleFilePaths(moduleName string) (biz, svc, param, code string) {
	biz = filepath.Join(g.config.RepoRoot, "internal", "api", "biz", fmt.Sprintf("%s.go", moduleName))
	svc = filepath.Join(g.config.RepoRoot, "internal", "api", "service", fmt.Sprintf("%s.go", moduleName))
	param = filepath.Join(g.config.RepoRoot, "internal", "api", "service", "param", fmt.Sprintf("%s.go", moduleName))
	code = filepath.Join(g.config.RepoRoot, "internal", "code", fmt.Sprintf("%s.go", moduleName))
	return
}

// moduleTestFilePaths 返回模块测试文件路径
func (g *Generator) moduleTestFilePaths(moduleName string) (bizTest, svcTest string) {
	bizTest = filepath.Join(g.config.RepoRoot, "internal", "api", "biz", fmt.Sprintf("%s_test.go", moduleName))
	svcTest = filepath.Join(g.config.RepoRoot, "internal", "api", "service", fmt.Sprintf("%s_test.go", moduleName))
	return
}

// printGeneratedFiles 打印生成的文件列表
func (g *Generator) printGeneratedFiles(moduleName string) {
	bizFile, svcFile, paramFile, codeFile := g.moduleFilePaths(moduleName)
	dataFile := filepath.Join(g.config.RepoRoot, "internal", "api", "data", fmt.Sprintf("%s.go", moduleName))

	fmt.Printf("\n📁 生成的文件:\n")
	fmt.Printf("  📄 业务逻辑: %s\n", bizFile)
	fmt.Printf("  📄 控制器: %s\n", svcFile)
	fmt.Printf("  📄 参数结构: %s\n", paramFile)
	fmt.Printf("  📄 数据层: %s\n", dataFile)
	fmt.Printf("  📄 错误码: %s\n", codeFile)
}

func (g *Generator) templateRenderer() (*templates.TemplateRenderer, error) {
	if g.renderer != nil {
		return g.renderer, nil
	}

	r, err := templates.NewTemplateRenderer()
	if err != nil {
		return nil, err
	}

	g.renderer = r
	return g.renderer, nil
}
