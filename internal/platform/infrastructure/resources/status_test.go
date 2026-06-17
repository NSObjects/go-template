package resources

import (
	"errors"
	"strings"
	"testing"
)

func TestCapabilityStatusReportsDisabledAsUnavailable(t *testing.T) {
	status := Disabled(CapabilityMongoDB)

	if status.Name != CapabilityMongoDB {
		t.Fatalf("Name = %q, want %q", status.Name, CapabilityMongoDB)
	}
	if status.Enabled {
		t.Fatal("Enabled = true, want false")
	}
	if status.Available {
		t.Fatal("Available = true, want false")
	}
	if status.State != StateDisabled {
		t.Fatalf("State = %q, want %q", status.State, StateDisabled)
	}
}

func TestCapabilityStatusReportsAvailableWhenValidated(t *testing.T) {
	status := Available(CapabilityMySQL, "ping ok")

	if !status.Enabled {
		t.Fatal("Enabled = false, want true")
	}
	if !status.Available {
		t.Fatal("Available = false, want true")
	}
	if status.State != StateAvailable {
		t.Fatalf("State = %q, want %q", status.State, StateAvailable)
	}
	if status.Message != "ping ok" {
		t.Fatalf("Message = %q, want ping ok", status.Message)
	}
}

func TestCapabilityStatusReportsUnavailableWithReason(t *testing.T) {
	status := Unavailable(CapabilityRedis, errors.New("dial refused"))

	if !status.Enabled {
		t.Fatal("Enabled = false, want true")
	}
	if status.Available {
		t.Fatal("Available = true, want false")
	}
	if status.State != StateUnavailable {
		t.Fatalf("State = %q, want %q", status.State, StateUnavailable)
	}
	if !strings.Contains(status.Message, "dial refused") {
		t.Fatalf("Message = %q, want failure reason", status.Message)
	}
}

func TestCapabilityErrorIncludesCapabilityName(t *testing.T) {
	err := NewCapabilityError(CapabilityTracing, "shutdown", errors.New("flush failed"))

	if !strings.Contains(err.Error(), CapabilityTracing) {
		t.Fatalf("error = %q, want capability name", err)
	}
	if !strings.Contains(err.Error(), "shutdown") {
		t.Fatalf("error = %q, want operation", err)
	}
}
