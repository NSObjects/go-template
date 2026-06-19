package boot

import (
	"errors"
	"fmt"

	"github.com/labstack/echo/v5"
	"github.com/samber/do/v2"

	"github.com/NSObjects/go-template/internal/platform/configs"
	"github.com/NSObjects/go-template/internal/platform/server"
)

// Module is one installable business capability in the application runtime.
type Module struct {
	name      string
	providers []func(do.Injector)
	mounts    []func(Context) error
	err       error
}

// Context is the framework surface given to a module during route mounting.
type Context struct {
	Config    configs.Config
	Server    *server.Server
	Router    *echo.Group
	Container do.Injector
	Resources *Resources
}

// Option customizes application assembly.
type Option func(*appOptions)

// ModuleOption customizes one boot module.
type ModuleOption func(*Module)

type appOptions struct {
	modules []Module
}

// NewModule creates a boot module with explicit providers and route mounts.
func NewModule(name string, opts ...ModuleOption) Module {
	module := Module{name: name}
	for _, opt := range opts {
		if opt != nil {
			opt(&module)
		}
	}
	return module
}

// Name returns the module name used in startup diagnostics.
func (m Module) Name() string {
	return m.name
}

// Provide registers one do provider when the module is installed.
func Provide[T any](provider do.Provider[T]) ModuleOption {
	return func(module *Module) {
		if provider == nil {
			module.addError(errors.New("provider is nil"))
			return
		}
		module.providers = append(module.providers, func(i do.Injector) {
			do.Provide(i, provider)
		})
	}
}

// Use installs a raw do package or binding as part of the module.
func Use(register func(do.Injector)) ModuleOption {
	return func(module *Module) {
		if register == nil {
			module.addError(errors.New("do registration is nil"))
			return
		}
		module.providers = append(module.providers, register)
	}
}

// Route mounts an HTTP route set after its handler has been resolved from do.
func Route[H any](register func(*echo.Group, H)) ModuleOption {
	return func(module *Module) {
		if register == nil {
			module.addError(errors.New("route registrar is nil"))
			return
		}
		module.mounts = append(module.mounts, func(ctx Context) error {
			handler, err := do.Invoke[H](ctx.Container)
			if err != nil {
				return err
			}
			register(ctx.Router, handler)
			return nil
		})
	}
}

// WithModules installs business modules into the boot-owned dependency graph.
func WithModules(modules ...Module) Option {
	return func(options *appOptions) {
		options.modules = append(options.modules, modules...)
	}
}

func newAppOptions(opts ...Option) (appOptions, error) {
	var options appOptions
	for _, opt := range opts {
		if opt != nil {
			opt(&options)
		}
	}
	if err := validateModules(options.modules); err != nil {
		return appOptions{}, err
	}
	return options, nil
}

func validateModules(modules []Module) error {
	seen := map[string]struct{}{}
	for _, module := range modules {
		name := module.Name()
		if name == "" {
			return fmt.Errorf("module name is required")
		}
		if module.err != nil {
			return fmt.Errorf("module %s: %w", name, module.err)
		}
		if len(module.mounts) == 0 {
			return fmt.Errorf("module %s has no route mounts", name)
		}
		if _, ok := seen[name]; ok {
			return fmt.Errorf("module %q is registered more than once", name)
		}
		seen[name] = struct{}{}
	}
	return nil
}

func provideModules(i do.Injector, modules []Module) error {
	for _, module := range modules {
		if err := module.provide(i); err != nil {
			return fmt.Errorf("register module %s: %w", module.Name(), err)
		}
	}
	return nil
}

func mountModules(ctx Context, modules []Module) error {
	for _, module := range modules {
		if err := module.mount(ctx); err != nil {
			return fmt.Errorf("mount module %s: %w", module.Name(), err)
		}
	}
	return nil
}

func (m Module) provide(i do.Injector) (err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			if recoveredErr, ok := recovered.(error); ok {
				err = fmt.Errorf("provider registration failed: %w", recoveredErr)
				return
			}
			err = fmt.Errorf("provider registration failed: %v", recovered)
		}
	}()
	for _, provider := range m.providers {
		provider(i)
	}
	return nil
}

func (m Module) mount(ctx Context) error {
	for _, mount := range m.mounts {
		if err := mount(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (m *Module) addError(err error) {
	m.err = errors.Join(m.err, err)
}
