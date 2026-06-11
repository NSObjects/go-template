package mongodb

import (
	"context"
	"fmt"
	"strings"

	"github.com/NSObjects/go-template/internal/api/data/db"
	"github.com/NSObjects/go-template/internal/configs"
	"github.com/NSObjects/go-template/internal/platform/module"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

const Name = "mongodb"

// Module exposes MongoDB as an observable capability.
type Module struct {
	cfg configs.Mongodb
}

// New creates a MongoDB capability module from configuration.
func New(cfg configs.Mongodb) Module {
	return Module{cfg: cfg}
}

// Descriptor reports whether MongoDB is enabled and usable.
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

// Validate returns a startup-blocking error when enabled MongoDB config is unusable.
func (m Module) Validate() error {
	if !m.cfg.Enabled {
		return nil
	}
	if err := validate(m.cfg); err != nil {
		return fmt.Errorf("%s capability unavailable: %w", Name, err)
	}
	return nil
}

// Start verifies the enabled MongoDB capability can become available.
func (m Module) Start(ctx context.Context) error {
	if !m.cfg.Enabled {
		return nil
	}
	if err := m.Validate(); err != nil {
		return err
	}
	database, err := db.MongoClient(m.cfg)
	if err != nil {
		return fmt.Errorf("%s capability unavailable: %w", Name, err)
	}
	defer database.Client().Disconnect(ctx)
	if err := database.Client().Ping(ctx, readpref.Primary()); err != nil {
		return fmt.Errorf("%s capability unavailable: %w", Name, err)
	}
	return nil
}

func validate(cfg configs.Mongodb) error {
	if blank(cfg.Host) || blank(cfg.Port) || blank(cfg.DataBase) {
		return fmt.Errorf("host, port, and database are required")
	}
	return nil
}

func blank(value string) bool {
	return strings.TrimSpace(value) == ""
}
