package newcmd

import (
	"fmt"
	"strings"

	"github.com/NSObjects/go-template/muban/project"
	"github.com/spf13/cobra"
)

// NewCommand builds the `tool new` command that bundles project and module generators.
func NewCommand() *cobra.Command {
	projectOpts := project.Options{}

	cmd := &cobra.Command{
		Use:   "new",
		Short: "新建项目与模块的脚手架工具",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return fmt.Errorf("未知参数: %s", strings.Join(args, " "))
			}

			if strings.TrimSpace(projectOpts.ModulePath) == "" {
				return fmt.Errorf("请使用 -m 或 --module 指定新项目的 Go Module 路径")
			}

			return project.RunWithWriter(projectOpts, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&projectOpts.ModulePath, "module", "m", "", "新项目的 Go Module 路径，例如: github.com/acme/demo")
	cmd.Flags().StringVarP(&projectOpts.OutputDir, "output", "o", "", "生成项目的目标目录，默认为当前目录下的模块名")
	cmd.Flags().StringVarP(&projectOpts.Name, "name", "n", "", "项目名称，用于替换模板中的标识")
	cmd.Flags().BoolVarP(&projectOpts.Force, "force", "f", false, "目标目录存在时覆盖")

	cmd.Example = strings.Join([]string{
		"  # 生成一个新项目到 ../awesome-api",
		"  muban new -m github.com/acme/awesome-api -o ../awesome-api",
		"",
		"  # 使用默认输出目录并覆盖已存在内容",
		"  muban new -m github.com/acme/awesome-api -f",
	}, "\n")

	cmd.SilenceUsage = true

	return cmd
}
