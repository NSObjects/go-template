// Package mongodb owns the process-level MongoDB resource.
package mongodb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"

	"github.com/NSObjects/go-template/internal/platform/configs"
	"github.com/NSObjects/go-template/internal/platform/infrastructure/resources"
)

// Resource wraps an optional MongoDB client and its lifecycle.
type Resource struct {
	enabled bool
	client  *mongo.Client
}

// Open creates the configured MongoDB resource. Disabled resources do not connect.
func Open(ctx context.Context, cfg configs.MongoDBConfig) (*Resource, error) {
	if ctx == nil {
		return nil, resources.NewCapabilityError(resources.CapabilityMongoDB, "open", errors.New("nil context"))
	}
	cfg = configs.Normalize(configs.Config{MongoDB: cfg}).MongoDB
	if !cfg.Enabled {
		return &Resource{}, nil
	}
	if err := configs.Validate(configs.Config{MongoDB: cfg}); err != nil {
		return nil, resources.NewCapabilityError(resources.CapabilityMongoDB, "configure", err)
	}

	client, err := mongo.Connect(clientOptions(cfg))
	if err != nil {
		return nil, resources.NewCapabilityError(resources.CapabilityMongoDB, "connect", err)
	}

	return &Resource{enabled: true, client: client}, nil
}

// Client returns the underlying MongoDB client, or nil when MongoDB is disabled.
func (r *Resource) Client() *mongo.Client {
	if r == nil {
		return nil
	}
	return r.client
}

// Check returns the current MongoDB capability status.
func (r *Resource) Check(ctx context.Context) resources.CapabilityStatus {
	if r == nil || !r.enabled {
		return resources.Disabled(resources.CapabilityMongoDB)
	}
	if r.client == nil {
		return resources.Unavailable(resources.CapabilityMongoDB, fmt.Errorf("mongodb client is nil"))
	}
	if ctx == nil {
		return resources.Unavailable(resources.CapabilityMongoDB, errors.New("nil context"))
	}
	if err := r.client.Ping(ctx, readpref.Primary()); err != nil {
		return resources.Unavailable(resources.CapabilityMongoDB, err)
	}
	return resources.Available(resources.CapabilityMongoDB, "ping ok")
}

// Close disconnects the MongoDB client.
func (r *Resource) Close(ctx context.Context) error {
	if r == nil || r.client == nil {
		return nil
	}
	if ctx == nil {
		return resources.NewCapabilityError(resources.CapabilityMongoDB, "close", errors.New("nil context"))
	}
	if err := r.client.Disconnect(ctx); err != nil {
		return resources.NewCapabilityError(resources.CapabilityMongoDB, "close", err)
	}
	return nil
}

func clientOptions(cfg configs.MongoDBConfig) *options.ClientOptions {
	timeout := time.Duration(cfg.ConnectTimeoutSeconds) * time.Second
	return options.Client().
		ApplyURI(cfg.URI).
		SetConnectTimeout(timeout).
		SetServerSelectionTimeout(timeout)
}
