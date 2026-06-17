package resources

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestResourcesReadyReportsUnavailableCapability(t *testing.T) {
	bundle := New(
		Component{
			Name: CapabilityRedis,
			Check: func(context.Context) CapabilityStatus {
				return Unavailable(CapabilityRedis, errors.New("dial refused"))
			},
		},
		Component{
			Name: CapabilityMongoDB,
			Check: func(context.Context) CapabilityStatus {
				return Disabled(CapabilityMongoDB)
			},
		},
	)

	err := bundle.Ready(context.Background())
	if err == nil {
		t.Fatal("Ready() error = nil, want unavailable Redis")
	}
	if !strings.Contains(err.Error(), CapabilityRedis) {
		t.Fatalf("Ready() error = %q, want redis capability", err)
	}
}

func TestResourcesReadyReportsAvailableCapabilities(t *testing.T) {
	bundle := New(
		Component{
			Name: CapabilityMySQL,
			Check: func(context.Context) CapabilityStatus {
				return Available(CapabilityMySQL, "ping ok")
			},
		},
		Component{
			Name: CapabilityRedis,
			Check: func(context.Context) CapabilityStatus {
				return Available(CapabilityRedis, "ping ok")
			},
		},
	)

	if err := bundle.Ready(context.Background()); err != nil {
		t.Fatalf("Ready() error = %v", err)
	}
}

func TestResourcesCloseReportsCapabilityFailure(t *testing.T) {
	bundle := New(Component{
		Name: CapabilityTracing,
		Close: func(context.Context) error {
			return errors.New("flush failed")
		},
	})

	err := bundle.Close(context.Background())
	if err == nil {
		t.Fatal("Close() error = nil, want tracing shutdown failure")
	}
	if !strings.Contains(err.Error(), CapabilityTracing) {
		t.Fatalf("Close() error = %q, want tracing capability", err)
	}
}

func TestResourcesCloseCleanly(t *testing.T) {
	closed := false
	bundle := New(Component{
		Name: CapabilityMySQL,
		Close: func(context.Context) error {
			closed = true
			return nil
		},
	})

	if err := bundle.Close(context.Background()); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	if !closed {
		t.Fatal("Close() did not invoke component close")
	}
}
