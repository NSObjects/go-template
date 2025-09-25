package codegen

import "github.com/spf13/cobra"

// NewCommand constructs the Cobra command for the error code generator.
func NewCommand() *cobra.Command {
	opts := Options{}

	cmd := &cobra.Command{
		Use:   "codegen",
		Short: "Generate error code helpers or documentation for constant definitions",
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(opts, args)
		},
	}

	cmd.Flags().StringVarP(&opts.TypeNames, "type", "t", "", "comma-separated list of type names; must be set")
	cmd.Flags().StringVarP(&opts.Output, "output", "o", "", "output file name; default srcdir/<type>_string.go")
	cmd.Flags().StringVar(&opts.TrimPrefix, "trimprefix", "", "trim the `prefix` from the generated constant names")
	cmd.Flags().StringVar(&opts.BuildTags, "tags", "", "comma-separated list of build tags to apply")
	cmd.Flags().BoolVar(&opts.Doc, "doc", false, "generate error code documentation in markdown format")

	_ = cmd.MarkFlagRequired("type")
	cmd.SilenceUsage = true

	return cmd
}
