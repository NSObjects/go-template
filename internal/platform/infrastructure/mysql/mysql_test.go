package mysql

import (
	"context"
	"testing"

	"github.com/NSObjects/go-template/internal/platform/configs"
	"github.com/NSObjects/go-template/internal/platform/infrastructure/resources"
)

func TestOpenDisabledMySQLDoesNotConnect(t *testing.T) {
	resource, err := Open(context.Background(), configs.MySQLConfig{})
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	if resource == nil {
		t.Fatal("Open() resource = nil, want disabled resource")
	}
	if resource.DB() != nil {
		t.Fatal("DB() != nil, want disabled resource without DB")
	}

	status := resource.Check(context.Background())
	if status.State != resources.StateDisabled {
		t.Fatalf("State = %q, want %q", status.State, resources.StateDisabled)
	}
}

func TestCheckEnabledMySQLWithoutDBReportsUnavailable(t *testing.T) {
	resource := &Resource{enabled: true}

	status := resource.Check(context.Background())
	if status.Name != resources.CapabilityMySQL {
		t.Fatalf("Name = %q, want %q", status.Name, resources.CapabilityMySQL)
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

func TestCloseDisabledMySQLIsNoop(t *testing.T) {
	resource := &Resource{}

	if err := resource.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}
