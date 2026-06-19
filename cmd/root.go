package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/NSObjects/go-template/internal/platform/configs"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   configs.DefaultAppName,
	Short: "Start the configured application",
	Long:  "Start the configured application through the boot composition root.",
	Run: func(_ *cobra.Command, _ []string) {
		Run(cfgFile)
	},
}

// Execute runs the root command.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "configs/config.toml",
		"config file")
}
