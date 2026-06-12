package mysql

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/NSObjects/go-template/internal/configs"
	"github.com/NSObjects/go-template/internal/platform/module"
	"gorm.io/gorm"
)

const (
	ModuleName   = "database-gorm-mysql"
	Capability   = "database.gorm"
	ProviderName = "mysql"
)

var (
	ErrInvalidConfig = errors.New("invalid mysql database config")
	ErrNotStarted    = errors.New("mysql database is not started")
)

// DBProvider returns the selected GORM database runtime.
type DBProvider interface {
	DB(context.Context) (*gorm.DB, error)
}

// Module exposes MySQL as a provider for the database.gorm capability.
type Module struct {
	cfg             configs.MysqlConfig
	provider        *provider
	defaultProvider bool
}

// Option customizes the MySQL database provider.
type Option func(*Module)

// WithDefault controls whether MySQL is the default database.gorm provider.
func WithDefault(defaultProvider bool) Option {
	return func(m *Module) {
		m.defaultProvider = defaultProvider
	}
}

// New creates a MySQL database capability module.
func New(cfg configs.MysqlConfig, options ...Option) *Module {
	m := &Module{
		cfg: cfg,
		provider: &provider{
			cfg:            cfg,
			runtimeFactory: newRuntime,
		},
	}
	for _, option := range options {
		option(m)
	}
	return m
}

// Descriptor reports whether MySQL can satisfy database.gorm.
func (m Module) Descriptor() module.Descriptor {
	status := module.CapabilityDisabled
	if m.cfg.Enabled {
		status = module.CapabilityEnabled
		if !hasRequiredConfig(m.cfg) {
			status = module.CapabilityUnavailable
		}
	}
	return module.Descriptor{
		Name: ModuleName,
		Kind: module.CapabilityModule,
		Provides: []module.Capability{
			{
				Name:     Capability,
				Provider: ProviderName,
				Status:   status,
				Default:  m.defaultProvider,
				Value:    m.provider,
			},
		},
	}
}

// Provider returns the lazy GORM DB provider owned by this module.
func (m Module) Provider() DBProvider {
	return m.provider
}

// Validate returns a startup-blocking error when selected MySQL config is unusable.
func (m Module) Validate() error {
	if !m.cfg.Enabled {
		return nil
	}
	if !hasRequiredConfig(m.cfg) {
		return fmt.Errorf("%w: host, port, user, and database are required", ErrInvalidConfig)
	}
	return nil
}

// Start verifies the selected MySQL database can initialize and accept connections.
func (m Module) Start(ctx context.Context) error {
	if !m.cfg.Enabled {
		return nil
	}
	if err := m.Validate(); err != nil {
		return err
	}
	return m.provider.Start(ctx)
}

// Stop releases the selected MySQL database runtime.
func (m Module) Stop(ctx context.Context) error {
	if m.provider == nil {
		return nil
	}
	return m.provider.Stop(ctx)
}

type provider struct {
	cfg            configs.MysqlConfig
	runtimeFactory runtimeFactory
	mu             sync.Mutex
	runtime        *runtime
}

type runtimeFactory func(context.Context, configs.MysqlConfig) (*runtime, error)

type runtime struct {
	db       *gorm.DB
	shutdown func(context.Context) error
}

func (p *provider) DB(ctx context.Context) (*gorm.DB, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.runtime == nil {
		return nil, ErrNotStarted
	}
	return p.runtime.db, nil
}

func (p *provider) Start(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.runtime != nil {
		return nil
	}
	runtime, err := p.runtimeFactory(ctx, p.cfg)
	if err != nil {
		return err
	}
	p.runtime = runtime
	return nil
}

func (p *provider) Stop(ctx context.Context) error {
	p.mu.Lock()
	runtime := p.runtime
	if runtime == nil {
		p.mu.Unlock()
		return nil
	}
	p.runtime = nil
	p.mu.Unlock()

	if runtime.shutdown == nil {
		return nil
	}
	return runtime.shutdown(ctx)
}

func newRuntime(ctx context.Context, cfg configs.MysqlConfig) (*runtime, error) {
	mysqlDB, err := newMysqlConnection(cfg)
	if err != nil {
		return nil, err
	}
	sqlDB, err := mysqlDB.DB()
	if err != nil {
		return nil, err
	}
	if err := sqlDB.PingContext(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, err
	}
	return &runtime{
		db: mysqlDB,
		shutdown: func(context.Context) error {
			return sqlDB.Close()
		},
	}, nil
}

func hasRequiredConfig(cfg configs.MysqlConfig) bool {
	return !blank(cfg.Host) && !blank(cfg.Port) && !blank(cfg.User) && !blank(cfg.Database)
}

func blank(value string) bool {
	return strings.TrimSpace(value) == ""
}

var _ DBProvider = (*provider)(nil)
