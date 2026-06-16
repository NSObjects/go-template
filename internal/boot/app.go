// Package boot owns application startup wiring.
package boot

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/NSObjects/go-template/internal/configs"
	"github.com/NSObjects/go-template/internal/server"
)

// Run loads configuration, assembles concrete runtime modules, and blocks until
// the process receives a shutdown signal.
func Run(configPath string) error {
	cfg, store, err := configs.BootstrapE(configPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	return server.New(cfg, store).Run(ctx)
}
