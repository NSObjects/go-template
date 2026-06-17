package boot

import (
	"context"
	"strings"
	"testing"

	goredis "github.com/redis/go-redis/v9"
	"github.com/samber/do/v2"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"gorm.io/gorm"

	"github.com/NSObjects/go-template/internal/platform/configs"
	"github.com/NSObjects/go-template/internal/platform/server"
)

func TestInjectorProvidesFrameworkRuntimeAndKeepsResourceClientsLazy(t *testing.T) {
	injector := newInjector(context.Background(), configs.Config{})
	t.Cleanup(func() {
		report := injector.ShutdownWithContext(context.Background())
		if report != nil && !report.Succeed {
			t.Fatalf("ShutdownWithContext() report = %v", report)
		}
	})

	if _, err := do.Invoke[*server.Server](injector); err != nil {
		t.Fatalf("Invoke[*server.Server]() error = %v", err)
	}
	if _, err := do.Invoke[*Resources](injector); err != nil {
		t.Fatalf("Invoke[*Resources]() error = %v", err)
	}
	assertDisabledResourceClient[*gorm.DB](t, injector, "mysql")
	assertDisabledResourceClient[*goredis.Client](t, injector, "redis")
	assertDisabledResourceClient[*mongo.Client](t, injector, "mongodb")
}

func assertDisabledResourceClient[T any](t *testing.T, injector do.Injector, capability string) {
	t.Helper()
	_, err := do.Invoke[T](injector)
	if err == nil {
		t.Fatalf("Invoke[%T]() error = nil, want disabled %s resource", *new(T), capability)
	}
	if !strings.Contains(err.Error(), capability) {
		t.Fatalf("Invoke[%T]() error = %q, want %s capability", *new(T), err, capability)
	}
}
