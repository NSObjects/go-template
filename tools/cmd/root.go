package cmd

import (
        "github.com/NSObjects/go-template/tools/codegen"
        dynamicsql "github.com/NSObjects/go-template/tools/dynamic-sql-gen"
        "github.com/NSObjects/go-template/tools/newcmd"
        "github.com/spf13/cobra"
)

// NewRootCommand assembles the shared CLI entry point for project tools.
func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "tool",
		Short: "辅助 go-template 项目的开发工具集合",
	}

        rootCmd.AddCommand(codegen.NewCommand())
        rootCmd.AddCommand(dynamicsql.NewCommand())
        rootCmd.AddCommand(newcmd.NewCommand())

	return rootCmd
}

// Execute runs the CLI.
func Execute() error {
	return NewRootCommand().Execute()
}
