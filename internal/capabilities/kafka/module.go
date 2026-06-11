package kafka

import (
	"context"
	"fmt"

	"github.com/NSObjects/go-template/internal/api/data/db"
	"github.com/NSObjects/go-template/internal/configs"
	"github.com/NSObjects/go-template/internal/platform/module"
)

const Name = "kafka"

// Module exposes Kafka as an observable capability.
type Module struct {
	cfg configs.KafkaConfig
}

// New creates a Kafka capability module from configuration.
func New(cfg configs.KafkaConfig) Module {
	return Module{cfg: cfg}
}

// Descriptor reports whether Kafka is enabled and usable.
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

// Validate returns a startup-blocking error when enabled Kafka config is unusable.
func (m Module) Validate() error {
	if !m.cfg.Enabled {
		return nil
	}
	if err := validate(m.cfg); err != nil {
		return fmt.Errorf("%s capability unavailable: %w", Name, err)
	}
	return nil
}

// Start verifies the enabled Kafka capability can become available.
func (m Module) Start(context.Context) error {
	if !m.cfg.Enabled {
		return nil
	}
	if err := m.Validate(); err != nil {
		return err
	}
	producer, err := db.NewKafkaProducer(m.cfg)
	if err != nil {
		return fmt.Errorf("%s capability unavailable: %w", Name, err)
	}
	defer producer.Close()
	return nil
}

func validate(cfg configs.KafkaConfig) error {
	if len(cfg.Brokers) == 0 {
		return fmt.Errorf("brokers are required")
	}
	return nil
}
