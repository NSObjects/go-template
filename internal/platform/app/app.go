// Package app assembles configured platform, capability, and business modules.
package app

import (
	"context"
	"fmt"
	"sync"

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

// CapabilityStopper is implemented by started capability modules that own
// resources which must be released during shutdown.
type CapabilityStopper interface {
	module.Module
	Stop(context.Context) error
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
	cfg                      configs.Config
	store                    *configs.Store
	logger                   applog.Logger
	report                   module.Report
	routes                   []http.Route
	startedCapabilityMu      sync.Mutex
	startedCapabilityModules []module.Module
	server                   *server.EchoServer
}

// Assemble validates modules and prepares runtime adapters.
func Assemble(options Options) (*App, error) {
	return AssembleWithContext(context.Background(), options)
}

// AssembleWithContext validates modules and prepares runtime adapters.
func AssembleWithContext(ctx context.Context, options Options) (*App, error) {
	options.Config = configs.Normalize(options.Config)
	adapters := options.EntryPointAdapters
	if len(adapters) == 0 {
		adapters = []string{http.EntryPointType}
	}
	assembly, err := module.AssembleRuntime(
		options.Modules,
		module.WithEntryPointAdapters(adapters...),
		module.WithCapabilitySelections(capabilitySelections(options)...),
	)
	if err != nil {
		return nil, err
	}
	report := assembly.Report()

	routes, err := http.RoutesFromEntryPoints(report.EntryPoints)
	if err != nil {
		return nil, err
	}

	selectedModules := selectedCapabilityModules(options.Modules, assembly.SelectedCapabilityModuleIndexes())
	for _, mod := range selectedModules {
		validator, ok := mod.(CapabilityValidator)
		if !ok {
			continue
		}
		if err := validator.Validate(); err != nil {
			return nil, err
		}
	}
	startedModules := make([]module.Module, 0, len(selectedModules))
	for _, mod := range selectedModules {
		starter, ok := mod.(CapabilityStarter)
		if !ok {
			continue
		}
		if err := starter.Start(ctx); err != nil {
			_ = stopCapabilityModules(context.Background(), startedModules)
			return nil, err
		}
		startedModules = append(startedModules, mod)
	}

	store := options.Store
	if store == nil {
		store = configs.NewStore(options.Config)
	}
	logger := applog.NewLogger(options.Config)

	return &App{
		cfg:                      options.Config,
		store:                    store,
		logger:                   logger,
		report:                   report,
		routes:                   routes,
		startedCapabilityModules: startedModules,
	}, nil
}

func stopCapabilityModules(ctx context.Context, modules []module.Module) error {
	var stopErr error
	for i := len(modules) - 1; i >= 0; i-- {
		stopper, ok := modules[i].(CapabilityStopper)
		if !ok {
			continue
		}
		if err := stopper.Stop(ctx); err != nil && stopErr == nil {
			stopErr = err
		}
	}
	return stopErr
}

func selectedCapabilityModules(modules []module.Module, indexes []int) []module.Module {
	selected := make([]module.Module, 0, len(indexes))
	for _, index := range indexes {
		if index < 0 || index >= len(modules) {
			continue
		}
		selected = append(selected, modules[index])
	}
	return selected
}

func capabilitySelections(options Options) []module.CapabilitySelection {
	if len(options.CapabilitySelections) > 0 {
		return options.CapabilitySelections
	}

	selections := make([]module.CapabilitySelection, 0, len(options.Config.Capabilities.Providers))
	for capability, provider := range options.Config.Capabilities.Providers {
		selections = append(selections, module.CapabilitySelection{
			Capability: capability,
			Provider:   provider,
		})
	}
	return selections
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

// Stop releases resources owned by started capability modules.
func (a *App) Stop(ctx context.Context) error {
	a.startedCapabilityMu.Lock()
	startedModules := a.startedCapabilityModules
	a.startedCapabilityModules = nil
	a.startedCapabilityMu.Unlock()

	return stopCapabilityModules(ctx, startedModules)
}

// Run starts the assembled server.
func (a *App) Run(ctx context.Context) error {
	defer func() {
		_ = a.Stop(ctx)
	}()
	a.Server().Run(a.cfg.System.Port)
	return nil
}
