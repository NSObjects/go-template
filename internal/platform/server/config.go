// Package server owns the Echo HTTP runtime.
package server

import (
	"time"

	"github.com/NSObjects/go-template/internal/platform/configs"
)

const (
	defaultServerPort     = configs.DefaultPort
	defaultReadTimeout    = 30 * time.Second
	defaultWriteTimeout   = 30 * time.Second
	defaultIdleTimeout    = 120 * time.Second
	defaultShutdownPeriod = 10 * time.Second
)

// Config contains HTTP server runtime settings.
type Config struct {
	Port string

	ReadTimeout time.Duration

	WriteTimeout time.Duration

	IdleTimeout time.Duration

	ShutdownTimeout time.Duration

	HideBanner bool
}

// DefaultConfig returns the default HTTP server settings.
func DefaultConfig() *Config {
	return &Config{
		Port:            defaultServerPort,
		ReadTimeout:     defaultReadTimeout,
		WriteTimeout:    defaultWriteTimeout,
		IdleTimeout:     defaultIdleTimeout,
		ShutdownTimeout: defaultShutdownPeriod,
		HideBanner:      true,
	}
}

// FromAppConfig derives HTTP server settings from application config.
func FromAppConfig(cfg configs.Config) *Config {
	cfg = configs.Normalize(cfg)

	return &Config{
		Port:            cfg.System.Port,
		ReadTimeout:     defaultReadTimeout,
		WriteTimeout:    defaultWriteTimeout,
		IdleTimeout:     defaultIdleTimeout,
		ShutdownTimeout: defaultShutdownPeriod,
		HideBanner:      true,
	}
}
