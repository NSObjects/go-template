package cmd

import (
	"github.com/NSObjects/go-template/muban/codegen"
	dynamicsql "github.com/NSObjects/go-template/muban/dynamic-sql-gen"
	"github.com/NSObjects/go-template/muban/modgen"
	"github.com/NSObjects/go-template/muban/newcmd"
	"github.com/spf13/cobra"
)

// NewRootCommand assembles the shared CLI entry point for project tools.
func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "gt",
		Short: "🚀 Go Template CLI - 快速构建 Go 项目的脚手架工具",
		Long: `🚀 Go Template CLI

一个强大的 Go 项目脚手架工具，帮助你快速创建和生成：
• 完整的 Go 项目结构
• 业务模块代码
• 错误码定义
• 动态 SQL 查询

让 Go 开发更高效！`,
	}

	rootCmd.AddCommand(codegen.NewCommand())
	rootCmd.AddCommand(dynamicsql.NewCommand())
	rootCmd.AddCommand(modgen.NewCommand())
	rootCmd.AddCommand(newcmd.NewCommand())

	return rootCmd
}

// Execute runs the CLI.
func Execute() error {
	return NewRootCommand().Execute()
}
