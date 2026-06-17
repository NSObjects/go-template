package boot

import (
	"context"
	"errors"

	goredis "github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"gorm.io/gorm"

	"github.com/NSObjects/go-template/internal/platform/configs"
	"github.com/NSObjects/go-template/internal/platform/infrastructure/mongodb"
	inframysql "github.com/NSObjects/go-template/internal/platform/infrastructure/mysql"
	"github.com/NSObjects/go-template/internal/platform/infrastructure/redis"
	infraresources "github.com/NSObjects/go-template/internal/platform/infrastructure/resources"
	"github.com/NSObjects/go-template/internal/platform/infrastructure/tracing"
	"github.com/NSObjects/go-template/internal/platform/server"
)

// Resources exposes framework-owned infrastructure resources to modules.
type Resources struct {
	status  *infraresources.Resources
	mysql   *inframysql.Resource
	redis   *redis.Resource
	mongodb *mongodb.Resource
}

type resourceAssembly struct {
	components []infraresources.Component
	mysql      *inframysql.Resource
	redis      *redis.Resource
	mongodb    *mongodb.Resource
}

func openResources(ctx context.Context, cfg configs.Config) (*Resources, error) {
	assembly := newResourceAssembly()
	if err := assembly.addMySQL(ctx, cfg.MySQL); err != nil {
		return nil, err
	}
	if err := assembly.addRedis(ctx, cfg.Redis); err != nil {
		return nil, errors.Join(err, assembly.closeOpened(context.Background()))
	}
	if err := assembly.addMongoDB(ctx, cfg.MongoDB); err != nil {
		return nil, errors.Join(err, assembly.closeOpened(context.Background()))
	}
	if err := assembly.addTracing(ctx, cfg); err != nil {
		return nil, errors.Join(err, assembly.closeOpened(context.Background()))
	}
	return assembly.resources(), nil
}

func newResourceAssembly() *resourceAssembly {
	return &resourceAssembly{
		components: []infraresources.Component{{
			Name: infraresources.CapabilityLogging,
			Check: func(context.Context) infraresources.CapabilityStatus {
				return infraresources.Available(infraresources.CapabilityLogging, "installed")
			},
		}},
	}
}

func (a *resourceAssembly) addMySQL(ctx context.Context, cfg configs.MySQLConfig) error {
	resource, err := inframysql.Open(ctx, cfg)
	if err != nil {
		return err
	}
	a.mysql = resource
	a.components = append(a.components, infraresources.Component{
		Name:  infraresources.CapabilityMySQL,
		Check: resource.Check,
		Close: func(context.Context) error { return resource.Close() },
	})
	return nil
}

func (a *resourceAssembly) addRedis(ctx context.Context, cfg configs.RedisConfig) error {
	resource, err := redis.Open(ctx, cfg)
	if err != nil {
		return err
	}
	a.redis = resource
	a.components = append(a.components, infraresources.Component{
		Name:  infraresources.CapabilityRedis,
		Check: resource.Check,
		Close: func(context.Context) error { return resource.Close() },
	})
	return nil
}

func (a *resourceAssembly) addMongoDB(ctx context.Context, cfg configs.MongoDBConfig) error {
	resource, err := mongodb.Open(ctx, cfg)
	if err != nil {
		return err
	}
	a.mongodb = resource
	a.components = append(a.components, infraresources.Component{
		Name:  infraresources.CapabilityMongoDB,
		Check: resource.Check,
		Close: resource.Close,
	})
	return nil
}

func (a *resourceAssembly) addTracing(ctx context.Context, cfg configs.Config) error {
	runtime, err := tracing.Start(ctx, cfg.Tracing, cfg.App.Name)
	if err != nil {
		return err
	}
	a.components = append(a.components, infraresources.Component{
		Name:  infraresources.CapabilityTracing,
		Check: runtime.Check,
		Close: runtime.Shutdown,
	})
	return nil
}

func (a *resourceAssembly) closeOpened(ctx context.Context) error {
	return infraresources.New(a.components...).Close(ctx)
}

func (a *resourceAssembly) resources() *Resources {
	return &Resources{
		status:  infraresources.New(a.components...),
		mysql:   a.mysql,
		redis:   a.redis,
		mongodb: a.mongodb,
	}
}

// MySQL returns the configured GORM DB when MySQL is enabled.
func (r *Resources) MySQL() (*gorm.DB, bool) {
	if r == nil || r.mysql == nil || r.mysql.DB() == nil {
		return nil, false
	}
	return r.mysql.DB(), true
}

// Redis returns the configured Redis client when Redis is enabled.
func (r *Resources) Redis() (*goredis.Client, bool) {
	if r == nil || r.redis == nil || r.redis.Client() == nil {
		return nil, false
	}
	return r.redis.Client(), true
}

// MongoDB returns the configured MongoDB client when MongoDB is enabled.
func (r *Resources) MongoDB() (*mongo.Client, bool) {
	if r == nil || r.mongodb == nil || r.mongodb.Client() == nil {
		return nil, false
	}
	return r.mongodb.Client(), true
}

// Close releases every opened infrastructure resource.
func (r *Resources) Close(ctx context.Context) error {
	if r == nil || r.status == nil {
		return nil
	}
	return r.status.Close(ctx)
}

// Shutdown releases every opened infrastructure resource.
func (r *Resources) Shutdown(ctx context.Context) error {
	return r.Close(ctx)
}

// Status returns server-facing capability status records.
func (r *Resources) Status(ctx context.Context) []server.CapabilityStatus {
	if r == nil || r.status == nil {
		return nil
	}
	statuses := r.status.Status(ctx)
	out := make([]server.CapabilityStatus, 0, len(statuses))
	for _, status := range statuses {
		out = append(out, server.CapabilityStatus{
			Name:      status.Name,
			Enabled:   status.Enabled,
			Available: status.Available,
			State:     status.State,
			Message:   status.Message,
		})
	}
	return out
}

// Ready reports whether every enabled resource is available.
func (r *Resources) Ready(ctx context.Context) error {
	if r == nil || r.status == nil {
		return nil
	}
	return r.status.Ready(ctx)
}
