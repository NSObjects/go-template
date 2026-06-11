package mysql

import (
	"context"
	"fmt"
	"strings"

	"github.com/NSObjects/go-template/internal/api/data/db"
	"github.com/NSObjects/go-template/internal/configs"
	"github.com/NSObjects/go-template/internal/platform/module"
)

const Name = "mysql"

// Module exposes MySQL as an observable capability.
type Module struct {
	cfg configs.MysqlConfig
}

// New creates a MySQL capability module from configuration.
func New(cfg configs.MysqlConfig) Module {
	return Module{cfg: cfg}
}

// Descriptor reports whether MySQL is enabled and usable.
func (m Module) Descriptor() module.Descriptor {
	status := module.CapabilityDisabled
	if m.cfg.Enabled {
		status = module.CapabilityEnabled
		if err := validate(m.cfg); err != nil {
			status = module.CapabilityUnavailable
		}
	}
	return module.Descriptor{
		Name: Name,
		Kind: module.CapabilityModule,
		Provides: []module.Capability{
			{Name: Name, Status: status},
		},
	}
}

// Validate returns a startup-blocking error when enabled MySQL config is unusable.
func (m Module) Validate() error {
	if !m.cfg.Enabled {
		return nil
	}
	if err := validate(m.cfg); err != nil {
		return fmt.Errorf("%s capability unavailable: %w", Name, err)
	}
	return nil
}

// Start verifies the enabled MySQL capability can become available.
func (m Module) Start(ctx context.Context) error {
	if !m.cfg.Enabled {
		return nil
	}
	if err := m.Validate(); err != nil {
		return err
	}
	mysqlDB, err := db.NewMysql(m.cfg)
	if err != nil {
		return fmt.Errorf("%s capability unavailable: %w", Name, err)
	}
	sqlDB, err := mysqlDB.DB()
	if err != nil {
		return fmt.Errorf("%s capability unavailable: %w", Name, err)
	}
	defer sqlDB.Close()
	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("%s capability unavailable: %w", Name, err)
	}
	return nil
}

func validate(cfg configs.MysqlConfig) error {
	if blank(cfg.Host) || blank(cfg.Port) || blank(cfg.User) || blank(cfg.Database) {
		return fmt.Errorf("host, port, user, and database are required")
	}
	return nil
}

func blank(value string) bool {
	return strings.TrimSpace(value) == ""
}
