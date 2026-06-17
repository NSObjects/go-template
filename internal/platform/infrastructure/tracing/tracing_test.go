package tracing

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/NSObjects/go-template/internal/platform/configs"
	"github.com/NSObjects/go-template/internal/platform/infrastructure/resources"
)

func TestStartDisabledTracingDoesNotRequireExporter(t *testing.T) {
	runtime, err := Start(context.Background(), configs.TracingConfig{}, "go-template")
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	if runtime == nil {
		t.Fatal("Start() runtime = nil, want disabled runtime")
	}
	if runtime.Enabled() {
		t.Fatal("Enabled() = true, want false")
	}

	status := runtime.Check(context.Background())
	if status.State != resources.StateDisabled {
		t.Fatalf("State = %q, want %q", status.State, resources.StateDisabled)
	}
}

func TestCheckEnabledTracingWithoutProviderReportsUnavailable(t *testing.T) {
	runtime := &Runtime{enabled: true}

	status := runtime.Check(context.Background())
	if status.Name != resources.CapabilityTracing {
		t.Fatalf("Name = %q, want %q", status.Name, resources.CapabilityTracing)
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

func TestShutdownTracingReportsCapabilityFailure(t *testing.T) {
	runtime := &Runtime{
		enabled: true,
		shutdown: func(context.Context) error {
			return errors.New("flush failed")
		},
	}

	err := runtime.Shutdown(context.Background())
	if err == nil {
		t.Fatal("Shutdown() error = nil, want flush failure")
	}
	if !strings.Contains(err.Error(), resources.CapabilityTracing) {
		t.Fatalf("Shutdown() error = %q, want tracing capability", err)
	}
}
