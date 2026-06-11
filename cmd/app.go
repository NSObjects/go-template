package cmd

import (
	"context"
	"fmt"
	"os"

	included "github.com/NSObjects/go-template/internal/app"
	"github.com/NSObjects/go-template/internal/configs"
	platformapp "github.com/NSObjects/go-template/internal/platform/app"
)

func Run(cfg string) {
	if err := run(cfg); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "application startup failed:", err)
		os.Exit(1)
	}
}

func run(cfg string) error {
	loaded, store, err := configs.BootstrapE(cfg)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	app, err := platformapp.Assemble(platformapp.Options{
		Config:  loaded,
		Store:   store,
		Modules: included.Modules(loaded),
	})
	if err != nil {
		return err
	}
	return app.Run(context.Background())
}
