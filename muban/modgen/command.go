package modgen

import (
	"fmt"
	"os"

	"github.com/NSObjects/go-template/muban/modgen/generator"
	"github.com/NSObjects/go-template/muban/modgen/utils"
	"github.com/spf13/cobra"
)

// Options captures the CLI parameters for module generation.
type Options struct {
	Name          string
	Route         string
	Force         bool
	OpenAPIFile   string
	GenerateTests bool
}

// NewCommand builds the Cobra command for module scaffolding generation.
func NewCommand() *cobra.Command {
	opts := Options{}

	cmd := &cobra.Command{
		Use:   "modgen",
		Short: "Generate module scaffolding and optional OpenAPI-based handlers",
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Name, "name", "n", "", "模块名，例如: user, article")
	cmd.Flags().StringVar(&opts.Route, "route", "", "基础路由前缀，例如: /articles (默认使用 name 的复数形式)")
	cmd.Flags().BoolVar(&opts.Force, "force", false, "若目标文件已存在则覆盖")
	cmd.Flags().StringVar(&opts.OpenAPIFile, "openapi", "", "OpenAPI3文档路径，例如: doc/openapi.yaml")
	cmd.Flags().BoolVar(&opts.GenerateTests, "tests", false, "同时生成测试用例（Table-driven测试）")
	cmd.Example = "  go run ./muban -- new module --name=user\n" +
		"  go run ./muban -- new module --name=article --openapi=doc/openapi.yaml --tests\n" +
		"  go run ./muban -- new module --openapi=doc/openapi.yaml"

	cmd.SilenceUsage = true

	return cmd
}

// Run executes the module generation with the provided options.
func Run(opts Options) error {
	generateAll := false

	if opts.OpenAPIFile != "" {
		if opts.Name == "" {
			generateAll = true
		}
	} else if opts.Name == "" {
		return fmt.Errorf("请使用 --name 指定模块名，或提供 --openapi 文档路径")
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取当前工作目录失败: %w", err)
	}

	repoRoot := utils.FindRepoRoot(cwd)
	if repoRoot == "" {
		return fmt.Errorf("未找到仓库根目录，请在项目内运行")
	}

	packagePath, err := utils.GetPackagePath(repoRoot)
	if err != nil {
		return fmt.Errorf("获取项目包路径失败: %w", err)
	}

	config := &generator.Config{
		Name:          opts.Name,
		Route:         opts.Route,
		Force:         opts.Force,
		OpenAPIFile:   opts.OpenAPIFile,
		GenerateTests: opts.GenerateTests,
		PackagePath:   packagePath,
		RepoRoot:      repoRoot,
		GenerateAll:   generateAll,
	}

	gen := generator.NewGenerator(config)
	if err := gen.Generate(); err != nil {
		return fmt.Errorf("生成失败: %w", err)
	}

	return nil
}
