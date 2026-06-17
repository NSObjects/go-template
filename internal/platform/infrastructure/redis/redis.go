// Package redis owns the process-level Redis resource.
package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"

	"github.com/NSObjects/go-template/internal/platform/configs"
	"github.com/NSObjects/go-template/internal/platform/infrastructure/resources"
)

// Resource wraps an optional Redis client and its lifecycle.
type Resource struct {
	enabled bool
	client  *goredis.Client
}

// Open creates the configured Redis resource. Disabled resources do not connect.
func Open(ctx context.Context, cfg configs.RedisConfig) (*Resource, error) {
	if ctx == nil {
		return nil, resources.NewCapabilityError(resources.CapabilityRedis, "open", errors.New("nil context"))
	}
	cfg = configs.Normalize(configs.Config{Redis: cfg}).Redis
	if !cfg.Enabled {
		return &Resource{}, nil
	}
	if err := configs.Validate(configs.Config{Redis: cfg}); err != nil {
		return nil, resources.NewCapabilityError(resources.CapabilityRedis, "configure", err)
	}

	client := goredis.NewClient(&goredis.Options{
		Addr:         cfg.Address,
		Username:     cfg.Username,
		Password:     cfg.Password,
		DB:           cfg.DB,
		DialTimeout:  seconds(cfg.DialTimeoutSeconds),
		ReadTimeout:  seconds(cfg.ReadTimeoutSeconds),
		WriteTimeout: seconds(cfg.WriteTimeoutSeconds),
	})

	return &Resource{enabled: true, client: client}, nil
}

// Client returns the underlying Redis client, or nil when Redis is disabled.
func (r *Resource) Client() *goredis.Client {
	if r == nil {
		return nil
	}
	return r.client
}

// Check returns the current Redis capability status.
func (r *Resource) Check(ctx context.Context) resources.CapabilityStatus {
	if r == nil || !r.enabled {
		return resources.Disabled(resources.CapabilityRedis)
	}
	if r.client == nil {
		return resources.Unavailable(resources.CapabilityRedis, fmt.Errorf("redis client is nil"))
	}
	if ctx == nil {
		return resources.Unavailable(resources.CapabilityRedis, errors.New("nil context"))
	}
	if err := r.client.Ping(ctx).Err(); err != nil {
		return resources.Unavailable(resources.CapabilityRedis, err)
	}
	return resources.Available(resources.CapabilityRedis, "ping ok")
}

// Close releases the Redis client.
func (r *Resource) Close() error {
	if r == nil || r.client == nil {
		return nil
	}
	if err := r.client.Close(); err != nil {
		return resources.NewCapabilityError(resources.CapabilityRedis, "close", err)
	}
	return nil
}

func seconds(value int) time.Duration {
	if value <= 0 {
		return 0
	}
	return time.Duration(value) * time.Second
}
