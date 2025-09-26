package project

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

//go:embed templates/*
var templateFS embed.FS

// EmbedTemplateGenerator 基于嵌入文件的模板生成器
type EmbedTemplateGenerator struct {
	outputDir   string
	modulePath  string
	projectName string
}

// NewEmbedTemplateGenerator 创建新的嵌入模板生成器
func NewEmbedTemplateGenerator(outputDir, modulePath, projectName string) *EmbedTemplateGenerator {
	return &EmbedTemplateGenerator{
		outputDir:   outputDir,
		modulePath:  modulePath,
		projectName: projectName,
	}
}

// Generate 生成项目文件
func (g *EmbedTemplateGenerator) Generate() error {
	// 创建必要的目录
	if err := g.createDirectories(); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 生成文件
	if err := g.generateFiles(); err != nil {
		return fmt.Errorf("生成文件失败: %w", err)
	}

	return nil
}

// createDirectories 创建项目目录结构
func (g *EmbedTemplateGenerator) createDirectories() error {
	dirs := []string{
		"cmd",
		"configs",
		"doc",
		"docs",
		"internal/api/biz",
		"internal/api/data",
		"internal/api/data/db",
		"internal/api/data/model",
		"internal/api/data/query",
		"internal/api/service",
		"internal/api/service/param",
		"internal/cache",
		"internal/code",
		"internal/configs",
		"internal/health",
		"internal/log",
		"internal/metrics",
		"internal/middleware",
		"internal/resp",
		"internal/server",
		"internal/server/middlewares",
		"internal/utils",
		"internal/validator",
		"internal/docs",
		"k8s",
		"scripts",
		"sql",
	}

	for _, dir := range dirs {
		fullPath := filepath.Join(g.outputDir, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return fmt.Errorf("创建目录 %s 失败: %w", dir, err)
		}
		fmt.Printf("✅ 创建目录: %s\n", fullPath)
	}

	return nil
}

// generateFiles 生成所有文件
func (g *EmbedTemplateGenerator) generateFiles() error {
	// 定义要生成的文件映射
	filesToGenerate := map[string]string{
		"main.go":                                               "templates/main.go.tmpl",
		"go.mod":                                                "templates/go.mod.tmpl",
		"README.md":                                             "templates/README.md.tmpl",
		"Makefile":                                              "templates/Makefile.tmpl",
		"env.example":                                           "templates/env.example.tmpl",
		"LICENSE":                                               "templates/LICENSE.tmpl",
		".gitignore":                                            "templates/.gitignore.tmpl",
		".golangci.yml":                                         "templates/.golangci.yml.tmpl",
		"docker-compose.yaml":                                   "templates/docker-compose.yaml.tmpl",
		"Dockerfile":                                            "templates/Dockerfile.tmpl",
		"golangci-report.xml":                                   "templates/golangci-report.xml.tmpl",
		"cmd/fx.go":                                             "templates/cmd/fx.go.tmpl",
		"cmd/root.go":                                           "templates/cmd/root.go.tmpl",
		"cmd/gen.go":                                            "templates/cmd/gen.go.tmpl",
		"configs/config.toml":                                   "templates/configs/config.toml.tmpl",
		"configs/rbac_model.conf":                               "templates/configs/rbac_model.conf.tmpl",
		"configs/README.md":                                     "templates/configs/README.md.tmpl",
		"doc/example_openapi.yaml":                              "templates/doc/example_openapi.yaml.tmpl",
		"docs/lint-summary.md":                                  "templates/docs/lint-summary.md.tmpl",
		"docs/lint.md":                                          "templates/docs/lint.md.tmpl",
		"docs/optimization-guide.md":                            "templates/docs/optimization-guide.md.tmpl",
		"internal/configs/config.go":                            "templates/internal/configs/config.go.tmpl",
		"internal/configs/bootstrap.go":                         "templates/internal/configs/bootstrap.go.tmpl",
		"internal/configs/store.go":                             "templates/internal/configs/store.go.tmpl",
		"internal/configs/file_source.go":                       "templates/internal/configs/file_source.go.tmpl",
		"internal/configs/merge.go":                             "templates/internal/configs/merge.go.tmpl",
		"internal/configs/hot_reload.go":                        "templates/internal/configs/hot_reload.go.tmpl",
		"internal/configs/etcd_source.go":                       "templates/internal/configs/etcd_source.go.tmpl",
		"internal/configs/consul_source.go":                     "templates/internal/configs/consul_source.go.tmpl",
		"internal/configs/config_simple_test.go":                "templates/internal/configs/config_simple_test.go.tmpl",
		"internal/log/logger.go":                                "templates/internal/log/logger.go.tmpl",
		"internal/log/factory.go":                               "templates/internal/log/factory.go.tmpl",
		"internal/log/console_sink.go":                          "templates/internal/log/console_sink.go.tmpl",
		"internal/log/file_sink.go":                             "templates/internal/log/file_sink.go.tmpl",
		"internal/log/global.go":                                "templates/internal/log/global.go.tmpl",
		"internal/log/slog.go":                                  "templates/internal/log/slog.go.tmpl",
		"internal/log/elasticsearch_sink.go":                    "templates/internal/log/elasticsearch_sink.go.tmpl",
		"internal/log/loki_sink.go":                             "templates/internal/log/loki_sink.go.tmpl",
		"internal/log/logger_test.go":                           "templates/internal/log/logger_test.go.tmpl",
		"internal/resp/response.go":                             "templates/internal/resp/response.go.tmpl",
		"internal/resp/response_test.go":                        "templates/internal/resp/response_test.go.tmpl",
		"internal/code/code.go":                                 "templates/internal/code/code.go.tmpl",
		"internal/code/base.go":                                 "templates/internal/code/base.go.tmpl",
		"internal/code/errors.go":                               "templates/internal/code/errors.go.tmpl",
		"internal/code/http_status.go":                          "templates/internal/code/http_status.go.tmpl",
		"internal/code/error_types.go":                          "templates/internal/code/error_types.go.tmpl",
		"internal/code/user.go":                                 "templates/internal/code/user.go.tmpl",
		"internal/code/code_generated.go":                       "templates/internal/code/code_generated.go.tmpl",
		"internal/code/code_test.go":                            "templates/internal/code/code_test.go.tmpl",
		"internal/code/error_code_generated.md":                 "templates/internal/code/error_code_generated.md.tmpl",
		"internal/code/error_types_test.go":                     "templates/internal/code/error_types_test.go.tmpl",
		"internal/code/errors_test.go":                          "templates/internal/code/errors_test.go.tmpl",
		"internal/code/registration_test.go":                    "templates/internal/code/registration_test.go.tmpl",
		"internal/api/biz/biz.go":                               "templates/internal/api/biz/biz.go.tmpl",
		"internal/api/service/service.go":                       "templates/internal/api/service/service.go.tmpl",
		"internal/api/service/example_router.go":                "templates/internal/api/service/example_router.go.tmpl",
		"internal/api/service/param/.gitkeep":                   "templates/internal/api/service/param/.gitkeep.tmpl",
		"internal/api/data/data.go":                             "templates/internal/api/data/data.go.tmpl",
		"internal/api/data/casbin.go":                           "templates/internal/api/data/casbin.go.tmpl",
		"internal/api/data/jwt.go":                              "templates/internal/api/data/jwt.go.tmpl",
		"internal/api/data/db/db.go":                            "templates/internal/api/data/db/db.go.tmpl",
		"internal/api/data/db/mysql.go":                         "templates/internal/api/data/db/mysql.go.tmpl",
		"internal/api/data/db/redis.go":                         "templates/internal/api/data/db/redis.go.tmpl",
		"internal/api/data/db/mongodb.go":                       "templates/internal/api/data/db/mongodb.go.tmpl",
		"internal/api/data/db/kafka.go":                         "templates/internal/api/data/db/kafka.go.tmpl",
		"internal/api/data/db/callback.go":                      "templates/internal/api/data/db/callback.go.tmpl",
		"internal/api/data/model/user.go":                       "templates/internal/api/data/model/user.go.tmpl",
		"internal/api/data/model/casbin_rule.gen.go":            "templates/internal/api/data/model/casbin_rule.gen.go.tmpl",
		"internal/api/data/query/gen.go":                        "templates/internal/api/data/query/gen.go.tmpl",
		"internal/api/data/query/casbin_rule.gen.go":            "templates/internal/api/data/query/casbin_rule.gen.go.tmpl",
		"internal/server/echo_server.go":                        "templates/internal/server/echo_server.go.tmpl",
		"internal/server/config.go":                             "templates/internal/server/config.go.tmpl",
		"internal/server/echo_server_test.go":                   "templates/internal/server/echo_server_test.go.tmpl",
		"internal/server/config_test.go":                        "templates/internal/server/config_test.go.tmpl",
		"internal/server/README.md":                             "templates/internal/server/README.md.tmpl",
		"internal/server/middlewares/middleware.go":             "templates/internal/server/middlewares/middleware.go.tmpl",
		"internal/server/middlewares/error.go":                  "templates/internal/server/middlewares/error.go.tmpl",
		"internal/server/middlewares/jwt.go":                    "templates/internal/server/middlewares/jwt.go.tmpl",
		"internal/server/middlewares/casbin.go":                 "templates/internal/server/middlewares/casbin.go.tmpl",
		"internal/server/middlewares/config.go":                 "templates/internal/server/middlewares/config.go.tmpl",
		"internal/server/middlewares/README.md":                 "templates/internal/server/middlewares/README.md.tmpl",
		"internal/server/middlewares/middleware_simple_test.go": "templates/internal/server/middlewares/middleware_simple_test.go.tmpl",
		"internal/utils/context.go":                             "templates/internal/utils/context.go.tmpl",
		"internal/utils/encrypt.go":                             "templates/internal/utils/encrypt.go.tmpl",
		"internal/utils/validator.go":                           "templates/internal/utils/validator.go.tmpl",
		"internal/utils/context_test.go":                        "templates/internal/utils/context_test.go.tmpl",
		"internal/utils/validator_test.go":                      "templates/internal/utils/validator_test.go.tmpl",
		"internal/validator/custom.go":                          "templates/internal/validator/custom.go.tmpl",
		"internal/health/checker.go":                            "templates/internal/health/checker.go.tmpl",
		"internal/metrics/prometheus.go":                        "templates/internal/metrics/prometheus.go.tmpl",
		"internal/middleware/rate_limit.go":                     "templates/internal/middleware/rate_limit.go.tmpl",
		"internal/cache/redis.go":                               "templates/internal/cache/redis.go.tmpl",
		"internal/docs/swagger.go":                              "templates/internal/docs/swagger.go.tmpl",
		"k8s/deployment.yaml":                                   "templates/k8s/deployment.yaml.tmpl",
		"scripts/dev.sh":                                        "templates/scripts/dev.sh.tmpl",
		"sql/create_users_table.sql":                            "templates/sql/create_users_table.sql.tmpl",
	}

	// 生成每个文件
	for outputFile, templateFile := range filesToGenerate {
		if err := g.generateFile(outputFile, templateFile); err != nil {
			return fmt.Errorf("生成文件 %s 失败: %w", outputFile, err)
		}
		fmt.Printf("✅ 生成文件: %s\n", filepath.Join(g.outputDir, outputFile))
	}

	return nil
}

// generateFile 生成单个文件
func (g *EmbedTemplateGenerator) generateFile(outputFile, templateFile string) error {
	// 从嵌入的文件系统读取模板文件
	templateContent, err := templateFS.ReadFile(templateFile)
	if err != nil {
		// 如果模板文件不存在，创建一个简单的占位符文件
		return g.createPlaceholderFile(outputFile)
	}

	// 解析模板
	tmpl, err := template.New(templateFile).Parse(string(templateContent))
	if err != nil {
		return fmt.Errorf("解析模板文件失败: %w", err)
	}

	// 创建输出文件
	outputPath := filepath.Join(g.outputDir, outputFile)
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("创建输出文件失败: %w", err)
	}
	defer file.Close()

	// 准备模板数据
	data := map[string]interface{}{
		"ModulePath":  g.modulePath,
		"ProjectName": g.projectName,
	}

	// 执行模板
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("执行模板失败: %w", err)
	}

	return nil
}

// createPlaceholderFile 创建占位符文件
func (g *EmbedTemplateGenerator) createPlaceholderFile(outputFile string) error {
	outputPath := filepath.Join(g.outputDir, outputFile)

	// 根据文件类型创建不同的占位符内容
	var content string
	switch filepath.Ext(outputFile) {
	case ".go":
		content = fmt.Sprintf(`package %s

// TODO: 实现 %s
`, filepath.Base(filepath.Dir(outputFile)), filepath.Base(outputFile))
	case ".toml", ".conf":
		content = "# TODO: 配置 " + filepath.Base(outputFile) + "\n"
	case ".md":
		content = "# TODO: 文档 " + filepath.Base(outputFile) + "\n"
	default:
		content = "# TODO: " + filepath.Base(outputFile) + "\n"
	}

	return os.WriteFile(outputPath, []byte(content), 0644)
}
