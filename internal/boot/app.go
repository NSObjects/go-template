// Package boot owns application startup wiring.
package boot

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/NSObjects/go-template/internal/configs"
	"github.com/NSObjects/go-template/internal/server"
)

// App is the assembled application runtime.
type App struct {
	server *server.Server
}

// NewApp assembles concrete runtime pieces from config.
func NewApp(cfg configs.Config) (*App, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("create server: %w", err)
	}

	return &App{server: srv}, nil
}

// Server returns the HTTP runtime so tests and boot wiring can register routes
// through the same interface used in production.
func (a *App) Server() *server.Server {
	if a == nil {
		return nil
	}
	return a.server
}

// Run starts the assembled application runtime.
func (a *App) Run(ctx context.Context) error {
	if a == nil || a.server == nil {
		return errors.New("run app: nil server")
	}
	return a.server.Run(ctx)
}

// Run loads configuration, assembles concrete runtime pieces, and blocks until
// the process receives a shutdown signal.
func Run(configPath string) error {
	cfg, err := configs.Load(configPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	app, err := NewApp(cfg)
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	return app.Run(ctx)
}
