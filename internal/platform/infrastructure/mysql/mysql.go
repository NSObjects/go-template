// Package mysql owns the process-level MySQL resource.
package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/NSObjects/go-template/internal/platform/configs"
	"github.com/NSObjects/go-template/internal/platform/infrastructure/resources"
)

// Resource wraps an optional GORM MySQL connection and its lifecycle.
type Resource struct {
	enabled bool
	db      *gorm.DB
}

// Open creates the configured MySQL resource. Disabled resources do not connect.
func Open(ctx context.Context, cfg configs.MySQLConfig) (*Resource, error) {
	if ctx == nil {
		return nil, resources.NewCapabilityError(resources.CapabilityMySQL, "open", errors.New("nil context"))
	}
	cfg = configs.Normalize(configs.Config{MySQL: cfg}).MySQL
	if !cfg.Enabled {
		return &Resource{}, nil
	}
	if err := configs.Validate(configs.Config{MySQL: cfg}); err != nil {
		return nil, resources.NewCapabilityError(resources.CapabilityMySQL, "configure", err)
	}

	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		return nil, resources.NewCapabilityError(resources.CapabilityMySQL, "open", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, resources.NewCapabilityError(resources.CapabilityMySQL, "db", err)
	}
	configurePool(sqlDB, cfg)

	return &Resource{enabled: true, db: db}, nil
}

// DB returns the underlying GORM DB, or nil when MySQL is disabled.
func (r *Resource) DB() *gorm.DB {
	if r == nil {
		return nil
	}
	return r.db
}

// Check returns the current MySQL capability status.
func (r *Resource) Check(ctx context.Context) resources.CapabilityStatus {
	if r == nil || !r.enabled {
		return resources.Disabled(resources.CapabilityMySQL)
	}
	sqlDB, err := r.sqlDB()
	if err != nil {
		return resources.Unavailable(resources.CapabilityMySQL, err)
	}
	if ctx == nil {
		return resources.Unavailable(resources.CapabilityMySQL, errors.New("nil context"))
	}
	if err := sqlDB.PingContext(ctx); err != nil {
		return resources.Unavailable(resources.CapabilityMySQL, err)
	}
	return resources.Available(resources.CapabilityMySQL, "ping ok")
}

// Close releases the MySQL connection pool.
func (r *Resource) Close() error {
	if r == nil || r.db == nil {
		return nil
	}
	sqlDB, err := r.sqlDB()
	if err != nil {
		return resources.NewCapabilityError(resources.CapabilityMySQL, "close", err)
	}
	if err := sqlDB.Close(); err != nil {
		return resources.NewCapabilityError(resources.CapabilityMySQL, "close", err)
	}
	return nil
}

func (r *Resource) sqlDB() (*sql.DB, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("mysql db is nil")
	}
	sqlDB, err := r.db.DB()
	if err != nil {
		return nil, err
	}
	if sqlDB == nil {
		return nil, fmt.Errorf("mysql sql db is nil")
	}
	return sqlDB, nil
}

func configurePool(db *sql.DB, cfg configs.MySQLConfig) {
	if cfg.MaxOpenConns > 0 {
		db.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		db.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetimeSeconds > 0 {
		db.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetimeSeconds) * time.Second)
	}
}
