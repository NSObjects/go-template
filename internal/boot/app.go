// Package boot owns application startup wiring.
package boot

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/samber/do/v2"

	"github.com/NSObjects/go-template/internal/platform/configs"
	"github.com/NSObjects/go-template/internal/platform/server"
)

// App is the assembled application runtime.
type App struct {
	config    configs.Config
	injector  do.Injector
	modules   []Module
	resources *Resources
	server    *server.Server
}

// NewApp assembles concrete runtime pieces from config.
func NewApp(cfg configs.Config, opts ...Option) (*App, error) {
	options, err := newAppOptions(opts...)
	if err != nil {
		return nil, err
	}

	injector := newInjector(context.Background(), cfg)
	app := &App{injector: injector, modules: options.modules}
	if err := provideModules(injector, options.modules); err != nil {
		if closeErr := app.close(context.Background()); closeErr != nil {
			return nil, errors.Join(err, closeErr)
		}
		return nil, err
	}

	if err := app.resolveRuntime(); err != nil {
		if closeErr := app.close(context.Background()); closeErr != nil {
			return nil, errors.Join(err, closeErr)
		}
		return nil, err
	}
	if err := app.mountModules(); err != nil {
		if closeErr := app.close(context.Background()); closeErr != nil {
			return nil, errors.Join(err, closeErr)
		}
		return nil, err
	}
	return app, nil
}

func (a *App) resolveRuntime() error {
	cfg, err := do.Invoke[configs.Config](a.injector)
	if err != nil {
		return err
	}
	resourceSet, err := do.Invoke[*Resources](a.injector)
	if err != nil {
		return err
	}
	srv, err := do.Invoke[*server.Server](a.injector)
	if err != nil {
		return err
	}
	a.config = cfg
	a.resources = resourceSet
	a.server = srv
	return nil
}

func (a *App) mountModules() error {
	if a == nil {
		return errors.New("mount modules: nil app")
	}
	return mountModules(Context{
		Config:    a.config,
		Server:    a.server,
		Router:    a.server.API(),
		Container: a.injector,
		Resources: a.resources,
	}, a.modules)
}

// Server returns the HTTP runtime so tests and boot wiring can register routes
// through the same interface used in production.
func (a *App) Server() *server.Server {
	if a == nil {
		return nil
	}
	return a.server
}

// Resources returns configured infrastructure resources owned by the app.
func (a *App) Resources() *Resources {
	if a == nil {
		return nil
	}
	return a.resources
}

// Container returns the boot-owned dependency graph.
func (a *App) Container() do.Injector {
	if a == nil {
		return nil
	}
	return a.injector
}

// Run starts the assembled application runtime.
func (a *App) Run(ctx context.Context) error {
	if a == nil || a.server == nil {
		return errors.New("run app: nil server")
	}
	runErr := a.server.Run(ctx)
	closeErr := a.close(context.Background())
	if runErr != nil && closeErr != nil {
		return errors.Join(runErr, closeErr)
	}
	if closeErr != nil {
		return closeErr
	}
	return runErr
}

func (a *App) close(ctx context.Context) error {
	if a.injector != nil {
		report := a.injector.ShutdownWithContext(ctx)
		if report == nil || report.Succeed {
			return nil
		}
		return report
	}
	if a.resources != nil {
		return a.resources.Close(ctx)
	}
	return nil
}

// Run loads configuration, assembles concrete runtime pieces, and blocks until
// the process receives a shutdown signal.
func Run(configPath string, opts ...Option) error {
	cfg, err := configs.Load(configPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	app, err := NewApp(cfg, opts...)
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	return app.Run(ctx)
}
