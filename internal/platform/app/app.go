// Package app assembles configured platform, capability, and business modules.
package app

import (
	"context"
	"fmt"

	"github.com/NSObjects/go-template/internal/configs"
	applog "github.com/NSObjects/go-template/internal/log"
	"github.com/NSObjects/go-template/internal/platform/http"
	"github.com/NSObjects/go-template/internal/platform/module"
	"github.com/NSObjects/go-template/internal/server"
)

// CapabilityValidator is implemented by modules that can detect unavailable
// enabled capabilities before runtime starts accepting triggers.
type CapabilityValidator interface {
	module.Module
	Validate() error
}

// CapabilityStarter is implemented by enabled capability modules that need a
// runtime availability check before serving entry point triggers.
type CapabilityStarter interface {
	module.Module
	Start(context.Context) error
}

// Options contains the inputs needed to assemble an application.
type Options struct {
	Config               configs.Config
	Store                *configs.Store
	Modules              []module.Module
	EntryPointAdapters   []string
	CapabilitySelections []module.CapabilitySelection
}

// App is the assembled runtime surface.
type App struct {
	cfg    configs.Config
	store  *configs.Store
	logger applog.Logger
	report module.Report
	routes []http.Route
	server *server.EchoServer
}

// Assemble validates modules and prepares runtime adapters.
func Assemble(options Options) (*App, error) {
	return AssembleWithContext(context.Background(), options)
}

// AssembleWithContext validates modules and prepares runtime adapters.
func AssembleWithContext(ctx context.Context, options Options) (*App, error) {
	for _, mod := range options.Modules {
		validator, ok := mod.(CapabilityValidator)
		if !ok {
			continue
		}
		if err := validator.Validate(); err != nil {
			return nil, err
		}
	}
	for _, mod := range options.Modules {
		starter, ok := mod.(CapabilityStarter)
		if !ok {
			continue
		}
		if err := starter.Start(ctx); err != nil {
			return nil, err
		}
	}

	adapters := options.EntryPointAdapters
	if len(adapters) == 0 {
		adapters = []string{http.EntryPointType}
	}
	report, err := module.Assemble(
		options.Modules,
		module.WithEntryPointAdapters(adapters...),
		module.WithCapabilitySelections(options.CapabilitySelections...),
	)
	if err != nil {
		return nil, err
	}

	routes, err := http.RoutesFromEntryPoints(report.EntryPoints)
	if err != nil {
		return nil, err
	}

	store := options.Store
	if store == nil {
		store = configs.NewStore(options.Config)
	}
	logger := applog.NewLogger(options.Config)

	return &App{
		cfg:    options.Config,
		store:  store,
		logger: logger,
		report: report,
		routes: routes,
	}, nil
}

// NewFromConfigFile loads config and assembles an app from the provided modules.
func NewFromConfigFile(path string, modules []module.Module) (*App, error) {
	cfg, store, err := configs.BootstrapE(path)
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}
	return Assemble(Options{
		Config:  cfg,
		Store:   store,
		Modules: modules,
	})
}

// Report returns the observable assembly report.
func (a *App) Report() module.Report {
	return a.report
}

// Routes returns HTTP routes exposed by assembled modules.
func (a *App) Routes() []http.Route {
	return append([]http.Route(nil), a.routes...)
}

// Server returns the Echo server, constructing it on first use.
func (a *App) Server() *server.EchoServer {
	if a.server == nil {
		a.server = server.NewEchoServer(a.routes, a.cfg, a.store)
	}
	return a.server
}

// Run starts the assembled server.
func (a *App) Run(_ context.Context) error {
	a.Server().Run(a.cfg.System.Port)
	return nil
}
