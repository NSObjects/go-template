package boot

import (
	"context"
	"testing"

	"github.com/NSObjects/go-template/internal/platform/configs"
	"github.com/NSObjects/go-template/internal/platform/infrastructure/resources"
	"github.com/NSObjects/go-template/internal/platform/server"
)

func TestOpenResourcesKeepsDisabledClientsUnavailableForBusinessWiring(t *testing.T) {
	runtime, err := openResources(context.Background(), configs.Config{})
	if err != nil {
		t.Fatalf("openResources() error = %v", err)
	}
	t.Cleanup(func() {
		if err := runtime.Close(context.Background()); err != nil {
			t.Fatalf("Close() error = %v", err)
		}
	})

	if db, ok := runtime.MySQL(); ok || db != nil {
		t.Fatalf("MySQL() = (%v, %v), want disabled resource", db, ok)
	}
	if client, ok := runtime.Redis(); ok || client != nil {
		t.Fatalf("Redis() = (%v, %v), want disabled resource", client, ok)
	}
	if client, ok := runtime.MongoDB(); ok || client != nil {
		t.Fatalf("MongoDB() = (%v, %v), want disabled resource", client, ok)
	}

	statuses := runtime.Status(context.Background())
	if len(statuses) != 5 {
		t.Fatalf("Status() length = %d, want fixed capability statuses", len(statuses))
	}
	if err := runtime.Ready(context.Background()); err != nil {
		t.Fatalf("Ready() error = %v, want disabled resources ignored", err)
	}
	assertCapabilityState(t, statuses, resources.CapabilityLogging, true, true)
	assertCapabilityState(t, statuses, resources.CapabilityMySQL, false, false)
	assertCapabilityState(t, statuses, resources.CapabilityRedis, false, false)
	assertCapabilityState(t, statuses, resources.CapabilityMongoDB, false, false)
	assertCapabilityState(t, statuses, resources.CapabilityTracing, false, false)
}

func assertCapabilityState(t *testing.T, statuses []server.CapabilityStatus, name string, wantEnabled, wantAvailable bool) {
	t.Helper()
	for _, status := range statuses {
		if status.Name != name {
			continue
		}
		if status.Enabled != wantEnabled || status.Available != wantAvailable {
			t.Fatalf("%s status = %#v, want enabled=%v available=%v", name, status, wantEnabled, wantAvailable)
		}
		return
	}
	t.Fatalf("status for %s not found in %#v", name, statuses)
}
