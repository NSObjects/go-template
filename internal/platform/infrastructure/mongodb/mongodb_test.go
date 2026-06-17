package mongodb

import (
	"context"
	"testing"

	"github.com/NSObjects/go-template/internal/platform/configs"
	"github.com/NSObjects/go-template/internal/platform/infrastructure/resources"
)

func TestOpenDisabledMongoDBDoesNotConnect(t *testing.T) {
	resource, err := Open(context.Background(), configs.MongoDBConfig{})
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	if resource == nil {
		t.Fatal("Open() resource = nil, want disabled resource")
	}
	if resource.Client() != nil {
		t.Fatal("Client() != nil, want disabled resource without client")
	}

	status := resource.Check(context.Background())
	if status.State != resources.StateDisabled {
		t.Fatalf("State = %q, want %q", status.State, resources.StateDisabled)
	}
}

func TestCheckEnabledMongoDBWithoutClientReportsUnavailable(t *testing.T) {
	resource := &Resource{enabled: true}

	status := resource.Check(context.Background())
	if status.Name != resources.CapabilityMongoDB {
		t.Fatalf("Name = %q, want %q", status.Name, resources.CapabilityMongoDB)
	}
	if !status.Enabled {
		t.Fatal("Enabled = false, want true")
	}
	if status.Available {
		t.Fatal("Available = true, want false")
	}
	if status.State != resources.StateUnavailable {
		t.Fatalf("State = %q, want %q", status.State, resources.StateUnavailable)
	}
}

func TestCloseDisabledMongoDBIsNoop(t *testing.T) {
	resource := &Resource{}

	if err := resource.Close(context.Background()); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}
