package project

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/NSObjects/go-template/tools/modgen/utils"
	"github.com/spf13/cobra"
)

// Options defines the CLI parameters for project generation.
type Options struct {
	ModulePath string
	OutputDir  string
	Name       string
	Force      bool
}

// NewCommand constructs the Cobra command for generating a new project from the template.
func NewCommand() *cobra.Command {
	opts := Options{}

	cmd := &cobra.Command{
		Use:   "project",
		Short: "从当前模板仓库生成一个全新的项目骨架",
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(opts)
		},
	}

	cmd.Flags().StringVar(&opts.ModulePath, "module", "", "新项目的 Go Module 路径，例如: github.com/acme/demo")
	cmd.Flags().StringVar(&opts.OutputDir, "output", "", "生成项目的目标目录，默认为当前目录下的模块名")
	cmd.Flags().StringVar(&opts.Name, "name", "", "项目名称，用于替换模板中的 go-template 标识")
	cmd.Flags().BoolVar(&opts.Force, "force", false, "目标目录存在时覆盖")

	cmd.Example = strings.Join([]string{
		"  # 生成一个新项目到 ../awesome-api",
		"  go run ./tools -- new project --module=github.com/acme/awesome-api --output=../awesome-api",
		"",
		"  # 使用默认输出目录并覆盖已存在内容",
		"  go run ./tools -- new project --module=github.com/acme/awesome-api --force",
	}, "\n")

	cmd.SilenceUsage = true

	return cmd
}

// Run executes the project generation workflow based on options.
func Run(opts Options) error {
	if opts.ModulePath == "" {
		return fmt.Errorf("请使用 --module 指定新项目的 Go Module 路径")
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取当前工作目录失败: %w", err)
	}

	repoRoot := utils.FindRepoRoot(cwd)
	if repoRoot == "" {
		return fmt.Errorf("未找到模板仓库的根目录，请在 go-template 项目内运行")
	}

	templateModulePath, err := utils.GetPackagePath(repoRoot)
	if err != nil {
		return fmt.Errorf("获取模板 Module 路径失败: %w", err)
	}

	outputDir := opts.OutputDir
	if outputDir == "" {
		outputDir = filepath.Base(opts.ModulePath)
	}
	if !filepath.IsAbs(outputDir) {
		outputDir = filepath.Join(cwd, outputDir)
	}

	projectName := opts.Name
	if projectName == "" {
		projectName = filepath.Base(opts.ModulePath)
	}

	generator, err := NewGenerator(GeneratorConfig{
		SourceRoot:          repoRoot,
		OutputDir:           outputDir,
		TemplateModulePath:  templateModulePath,
		TargetModulePath:    opts.ModulePath,
		TemplateProjectName: filepath.Base(templateModulePath),
		TargetProjectName:   projectName,
		Force:               opts.Force,
	})
	if err != nil {
		return err
	}

	fmt.Printf("🚀 正在生成新项目: %s\n", projectName)
	fmt.Printf("📦 Module: %s\n", opts.ModulePath)
	fmt.Printf("📁 目标目录: %s\n", outputDir)

	if err := generator.Generate(); err != nil {
		return fmt.Errorf("生成项目失败: %w", err)
	}

	fmt.Printf("✅ 项目生成完成！接下来可进入目录 %s 并执行 make dev-setup\n", outputDir)
	return nil
}
