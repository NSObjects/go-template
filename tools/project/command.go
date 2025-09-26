package project

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

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
			return RunWithWriter(opts, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVar(&opts.ModulePath, "module", "", "新项目的 Go Module 路径，例如: github.com/acme/demo")
	cmd.Flags().StringVar(&opts.OutputDir, "output", "", "生成项目的目标目录，默认为当前目录下的模块名")
	cmd.Flags().StringVar(&opts.Name, "name", "", "项目名称，用于替换模板中的标识")
	cmd.Flags().BoolVar(&opts.Force, "force", false, "目标目录存在时覆盖")

	cmd.Example = strings.Join([]string{
		"  # 生成一个新项目到 ../awesome-api",
		"  gt new project --module=github.com/acme/awesome-api --output=../awesome-api",
		"",
		"  # 使用默认输出目录并覆盖已存在内容",
		"  gt new project --module=github.com/acme/awesome-api --force",
		"",
		"  # 从任何目录生成项目（无需在模板目录下）",
		"  gt new project --module=github.com/acme/awesome-api",
	}, "\n")

	cmd.SilenceUsage = true

	return cmd
}

// Run executes the project generation workflow based on options.
func Run(opts Options) error {
	return RunWithWriter(opts, os.Stdout)
}

// RunWithWriter 执行生成逻辑并将日志写入指定 writer。
func RunWithWriter(opts Options, out io.Writer) error {
	modulePath := strings.TrimSpace(opts.ModulePath)
	if modulePath == "" {
		return fmt.Errorf("请使用 --module 指定新项目的 Go Module 路径")
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取当前工作目录失败: %w", err)
	}

	outputDir := strings.TrimSpace(opts.OutputDir)
	if outputDir == "" {
		outputDir = filepath.Base(modulePath)
	}
	if !filepath.IsAbs(outputDir) {
		outputDir = filepath.Join(cwd, outputDir)
	}

	projectName := strings.TrimSpace(opts.Name)
	if projectName == "" {
		projectName = filepath.Base(modulePath)
	}

	// 检查输出目录是否存在
	if _, err := os.Stat(outputDir); err == nil {
		if !opts.Force {
			return fmt.Errorf("目标目录已存在: %s (使用 --force 可覆盖)", outputDir)
		}
		if err := ensureSafeRemoval(outputDir, cwd); err != nil {
			return err
		}
		if err := os.RemoveAll(outputDir); err != nil {
			return fmt.Errorf("清理目标目录失败: %w", err)
		}
	}

	fmt.Fprintf(out, "🚀 正在生成新项目: %s\n", projectName)
	fmt.Fprintf(out, "📦 Module: %s\n", modulePath)
	fmt.Fprintf(out, "📁 目标目录: %s\n", outputDir)

	// 使用嵌入模板生成器
	generator := NewEmbedTemplateGenerator(
		outputDir,
		modulePath,
		projectName,
		WithLogger(func(format string, args ...interface{}) {
			fmt.Fprintf(out, format, args...)
		}),
	)
	if err := generator.Generate(); err != nil {
		return fmt.Errorf("生成项目失败: %w", err)
	}

	fmt.Fprintf(out, "✅ 项目生成完成！接下来可进入目录 %s 并执行 make dev-setup\n", outputDir)
	return nil
}

func ensureSafeRemoval(outputDir, cwd string) error {
	cleaned := filepath.Clean(outputDir)
	if cleaned == "" || cleaned == "." {
		return fmt.Errorf("拒绝清理危险目录: %s", outputDir)
	}

	if cleaned == filepath.Clean(cwd) {
		return fmt.Errorf("拒绝清理当前工作目录: %s", outputDir)
	}

	if isRootPath(cleaned) {
		return fmt.Errorf("拒绝清理系统根目录: %s", outputDir)
	}

	return nil
}

func isRootPath(path string) bool {
	if runtime.GOOS == "windows" {
		if vol := filepath.VolumeName(path); vol != "" {
			rest := strings.TrimPrefix(path, vol)
			rest = filepath.Clean(rest)
			return rest == string(filepath.Separator)
		}
	}

	return path == string(filepath.Separator)
}
