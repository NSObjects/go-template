package redis

import (
	"context"
	"fmt"
	"strings"

	"github.com/NSObjects/go-template/internal/api/data/db"
	"github.com/NSObjects/go-template/internal/configs"
	"github.com/NSObjects/go-template/internal/platform/module"
)

const Name = "redis"

// Module exposes Redis as an observable capability.
type Module struct {
	cfg configs.RedisConfig
}

// New creates a Redis capability module from configuration.
func New(cfg configs.RedisConfig) Module {
	return Module{cfg: cfg}
}

// Descriptor reports whether Redis is enabled and usable.
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

// Validate returns a startup-blocking error when enabled Redis config is unusable.
func (m Module) Validate() error {
	if !m.cfg.Enabled {
		return nil
	}
	if err := validate(m.cfg); err != nil {
		return fmt.Errorf("%s capability unavailable: %w", Name, err)
	}
	return nil
}

// Start verifies the enabled Redis capability can become available.
func (m Module) Start(ctx context.Context) error {
	if !m.cfg.Enabled {
		return nil
	}
	if err := m.Validate(); err != nil {
		return err
	}
	client := db.NewRedis(m.cfg)
	defer client.Close()
	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("%s capability unavailable: %w", Name, err)
	}
	return nil
}

func validate(cfg configs.RedisConfig) error {
	if blank(cfg.Host) || blank(cfg.Port) {
		return fmt.Errorf("host and port are required")
	}
	return nil
}

func blank(value string) bool {
	return strings.TrimSpace(value) == ""
}
