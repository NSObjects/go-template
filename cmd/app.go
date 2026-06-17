// Package cmd contains CLI entrypoints for the application.
package cmd

import (
	"fmt"
	"os"

	"github.com/NSObjects/go-template/internal/boot"
)

// Run starts the application and exits the process when startup fails.
func Run(cfg string) {
	if err := run(cfg); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "application startup failed:", err)
		os.Exit(1)
	}
}

func run(cfg string) error {
	return boot.Run(cfg, boot.WithModules(boot.BusinessModules()...))
}
