package newcmd

import (
	"github.com/NSObjects/go-template/tools/modgen"
	"github.com/NSObjects/go-template/tools/project"
	"github.com/spf13/cobra"
)

// NewCommand builds the `tool new` command that bundles project and module generators.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new",
		Short: "新建项目与模块的脚手架工具",
	}

	projectCmd := project.NewCommand()
	projectCmd.Use = "project"
	projectCmd.Short = "从模板仓库创建一个全新的项目"
	projectCmd.Aliases = append(projectCmd.Aliases, "proj")

	moduleCmd := modgen.NewCommand()
	moduleCmd.Use = "module"
	moduleCmd.Short = "生成业务模块，支持默认模板与 OpenAPI 文档"
	moduleCmd.Aliases = append(moduleCmd.Aliases, "modgen")

	cmd.AddCommand(projectCmd)
	cmd.AddCommand(moduleCmd)

	return cmd
}
