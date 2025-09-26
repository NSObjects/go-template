package project

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// FileTemplateGenerator 基于文件系统的模板生成器
type FileTemplateGenerator struct {
	outputDir   string
	modulePath  string
	projectName string
	templateDir string
}

// NewFileTemplateGenerator 创建新的文件模板生成器
func NewFileTemplateGenerator(outputDir, modulePath, projectName string) *FileTemplateGenerator {
	// 获取当前工作目录
	wd, _ := os.Getwd()
	templateDir := filepath.Join(wd, "tools/project/templates")

	return &FileTemplateGenerator{
		outputDir:   outputDir,
		modulePath:  modulePath,
		projectName: projectName,
		templateDir: templateDir,
	}
}

// Generate 生成项目文件
func (g *FileTemplateGenerator) Generate() error {
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
func (g *FileTemplateGenerator) createDirectories() error {
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
func (g *FileTemplateGenerator) generateFiles() error {
	// 定义要生成的文件映射
	filesToGenerate := map[string]string{
		"main.go":                                               "main.go.tmpl",
		"go.mod":                                                "go.mod.tmpl",
		"README.md":                                             "README.md.tmpl",
		"Makefile":                                              "Makefile.tmpl",
		"env.example":                                           "env.example.tmpl",
		"cmd/fx.go":                                             "cmd/fx.go.tmpl",
		"cmd/root.go":                                           "cmd/root.go.tmpl",
		"cmd/gen.go":                                            "cmd/gen.go.tmpl",
		"configs/config.toml":                                   "configs/config.toml.tmpl",
		"configs/rbac_model.conf":                               "configs/rbac_model.conf.tmpl",
		"internal/configs/config.go":                            "internal/configs/config.go.tmpl",
		"internal/configs/bootstrap.go":                         "internal/configs/bootstrap.go.tmpl",
		"internal/configs/store.go":                             "internal/configs/store.go.tmpl",
		"internal/configs/file_source.go":                       "internal/configs/file_source.go.tmpl",
		"internal/configs/merge.go":                             "internal/configs/merge.go.tmpl",
		"internal/configs/hot_reload.go":                        "internal/configs/hot_reload.go.tmpl",
		"internal/configs/etcd_source.go":                       "internal/configs/etcd_source.go.tmpl",
		"internal/configs/consul_source.go":                     "internal/configs/consul_source.go.tmpl",
		"internal/log/logger.go":                                "internal/log/logger.go.tmpl",
		"internal/log/factory.go":                               "internal/log/factory.go.tmpl",
		"internal/log/console_sink.go":                          "internal/log/console_sink.go.tmpl",
		"internal/log/file_sink.go":                             "internal/log/file_sink.go.tmpl",
		"internal/log/global.go":                                "internal/log/global.go.tmpl",
		"internal/log/slog.go":                                  "internal/log/slog.go.tmpl",
		"internal/resp/response.go":                             "internal/resp/response.go.tmpl",
		"internal/code/code.go":                                 "internal/code/code.go.tmpl",
		"internal/code/base.go":                                 "internal/code/base.go.tmpl",
		"internal/code/errors.go":                               "internal/code/errors.go.tmpl",
		"internal/code/http_status.go":                          "internal/code/http_status.go.tmpl",
		"internal/code/error_types.go":                          "internal/code/error_types.go.tmpl",
		"internal/api/biz/biz.go":                               "internal/api/biz/biz.go.tmpl",
		"internal/api/service/service.go":                       "internal/api/service/service.go.tmpl",
		"internal/api/data/data.go":                             "internal/api/data/data.go.tmpl",
		"internal/api/data/casbin.go":                           "internal/api/data/casbin.go.tmpl",
		"internal/api/data/jwt.go":                              "internal/api/data/jwt.go.tmpl",
		"internal/api/data/db/db.go":                            "internal/api/data/db/db.go.tmpl",
		"internal/api/data/db/mysql.go":                         "internal/api/data/db/mysql.go.tmpl",
		"internal/api/data/db/redis.go":                         "internal/api/data/db/redis.go.tmpl",
		"internal/api/data/db/mongodb.go":                       "internal/api/data/db/mongodb.go.tmpl",
		"internal/api/data/db/kafka.go":                         "internal/api/data/db/kafka.go.tmpl",
		"internal/api/data/db/callback.go":                      "internal/api/data/db/callback.go.tmpl",
		"internal/api/data/model/user.go":                       "internal/api/data/model/user.go.tmpl",
		"internal/api/data/query/gen.go":                        "internal/api/data/query/gen.go.tmpl",
		"internal/server/echo_server.go":                        "internal/server/echo_server.go.tmpl",
		"internal/server/config.go":                             "internal/server/config.go.tmpl",
		"internal/server/middlewares/middleware.go":             "internal/server/middlewares/middleware.go.tmpl",
		"internal/server/middlewares/error.go":                  "internal/server/middlewares/error.go.tmpl",
		"internal/server/middlewares/jwt.go":                    "internal/server/middlewares/jwt.go.tmpl",
		"internal/server/middlewares/casbin.go":                 "internal/server/middlewares/casbin.go.tmpl",
		"internal/server/middlewares/config.go":                 "internal/server/middlewares/config.go.tmpl",
		"internal/utils/context.go":                             "internal/utils/context.go.tmpl",
		"internal/utils/encrypt.go":                             "internal/utils/encrypt.go.tmpl",
		"internal/utils/validator.go":                           "internal/utils/validator.go.tmpl",
		"internal/validator/custom.go":                          "internal/validator/custom.go.tmpl",
		"internal/health/checker.go":                            "internal/health/checker.go.tmpl",
		"internal/metrics/prometheus.go":                        "internal/metrics/prometheus.go.tmpl",
		"internal/middleware/rate_limit.go":                     "internal/middleware/rate_limit.go.tmpl",
		"internal/cache/redis.go":                               "internal/cache/redis.go.tmpl",
		"internal/docs/swagger.go":                              "internal/docs/swagger.go.tmpl",
		"k8s/deployment.yaml":                                   "k8s/deployment.yaml.tmpl",
		"scripts/dev.sh":                                        "scripts/dev.sh.tmpl",
		"sql/create_users_table.sql":                            "sql/create_users_table.sql.tmpl",
		"docker-compose.yaml":                                   "docker-compose.yaml.tmpl",
		"Dockerfile":                                            "Dockerfile.tmpl",
		"doc/example_openapi.yaml":                              "doc/example_openapi.yaml.tmpl",
		"internal/configs/config_simple_test.go":                "internal/configs/config_simple_test.go.tmpl",
		"internal/log/elasticsearch_sink.go":                    "internal/log/elasticsearch_sink.go.tmpl",
		"internal/log/loki_sink.go":                             "internal/log/loki_sink.go.tmpl",
		"internal/log/logger_test.go":                           "internal/log/logger_test.go.tmpl",
		"internal/resp/response_test.go":                        "internal/resp/response_test.go.tmpl",
		"internal/code/user.go":                                 "internal/code/user.go.tmpl",
		"internal/code/code_generated.go":                       "internal/code/code_generated.go.tmpl",
		"internal/code/code_test.go":                            "internal/code/code_test.go.tmpl",
		"internal/code/error_code_generated.md":                 "internal/code/error_code_generated.md.tmpl",
		"internal/api/service/param/.gitkeep":                   "internal/api/service/param/.gitkeep.tmpl",
		"internal/api/data/model/casbin_rule.gen.go":            "internal/api/data/model/casbin_rule.gen.go.tmpl",
		"internal/api/data/query/casbin_rule.gen.go":            "internal/api/data/query/casbin_rule.gen.go.tmpl",
		"internal/server/echo_server_test.go":                   "internal/server/echo_server_test.go.tmpl",
		"internal/server/config_test.go":                        "internal/server/config_test.go.tmpl",
		"internal/server/README.md":                             "internal/server/README.md.tmpl",
		"internal/server/middlewares/README.md":                 "internal/server/middlewares/README.md.tmpl",
		"internal/server/middlewares/middleware_simple_test.go": "internal/server/middlewares/middleware_simple_test.go.tmpl",
		"internal/utils/context_test.go":                        "internal/utils/context_test.go.tmpl",
		"internal/utils/validator_test.go":                      "internal/utils/validator_test.go.tmpl",
		"go.sum":                                                "go.sum.tmpl",
		"LICENSE":                                               "LICENSE.tmpl",
		"configs/README.md":                                     "configs/README.md.tmpl",
		"docs/lint-summary.md":                                  "docs/lint-summary.md.tmpl",
		"docs/lint.md":                                          "docs/lint.md.tmpl",
		"docs/optimization-guide.md":                            "docs/optimization-guide.md.tmpl",
		"internal/code/error_types_test.go":                     "internal/code/error_types_test.go.tmpl",
		"internal/code/errors_test.go":                          "internal/code/errors_test.go.tmpl",
		"internal/code/registration_test.go":                    "internal/code/registration_test.go.tmpl",
		"golangci-report.xml":                                   "golangci-report.xml.tmpl",
		".gitignore":                                            ".gitignore.tmpl",
		".golangci.yml":                                         ".golangci.yml.tmpl",
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
func (g *FileTemplateGenerator) generateFile(outputFile, templateFile string) error {
	// 读取模板文件
	templatePath := filepath.Join(g.templateDir, templateFile)

	// 检查模板文件是否存在
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		// 如果模板文件不存在，创建一个简单的占位符文件
		return g.createPlaceholderFile(outputFile)
	}

	tmpl, err := template.ParseFiles(templatePath)
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
func (g *FileTemplateGenerator) createPlaceholderFile(outputFile string) error {
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
