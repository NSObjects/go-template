// Package resources coordinates fixed process-level infrastructure resources.
package resources

import (
	"errors"
	"fmt"
)

const (
	// CapabilityLogging identifies the process logger capability.
	CapabilityLogging = "logging"
	// CapabilityMySQL identifies the MySQL resource capability.
	CapabilityMySQL = "mysql"
	// CapabilityRedis identifies the Redis resource capability.
	CapabilityRedis = "redis"
	// CapabilityMongoDB identifies the MongoDB resource capability.
	CapabilityMongoDB = "mongodb"
	// CapabilityTracing identifies the OpenTelemetry/Jaeger tracing capability.
	CapabilityTracing = "tracing"
)

const (
	// StateDisabled means the capability was not enabled by configuration.
	StateDisabled = "disabled"
	// StateAvailable means the enabled capability has passed validation.
	StateAvailable = "available"
	// StateUnavailable means the enabled capability failed validation.
	StateUnavailable = "unavailable"
)

// CapabilityStatus is the observable runtime state for one fixed capability.
type CapabilityStatus struct {
	Name      string `json:"name"`
	Enabled   bool   `json:"enabled"`
	Available bool   `json:"available"`
	State     string `json:"state"`
	Message   string `json:"message,omitempty"`
}

// Disabled reports a configured-off capability.
func Disabled(name string) CapabilityStatus {
	return CapabilityStatus{
		Name:    name,
		Enabled: false,
		State:   StateDisabled,
	}
}

// Available reports an enabled and validated capability.
func Available(name, message string) CapabilityStatus {
	return CapabilityStatus{
		Name:      name,
		Enabled:   true,
		Available: true,
		State:     StateAvailable,
		Message:   message,
	}
}

// Unavailable reports an enabled capability that failed validation.
func Unavailable(name string, err error) CapabilityStatus {
	message := "unavailable"
	if err != nil {
		message = err.Error()
	}
	return CapabilityStatus{
		Name:    name,
		Enabled: true,
		State:   StateUnavailable,
		Message: message,
	}
}

// CapabilityError is an error scoped to a named infrastructure capability.
type CapabilityError struct {
	Name string
	Op   string
	Err  error
}

// NewCapabilityError wraps err with a capability name and operation.
func NewCapabilityError(name, op string, err error) error {
	if err == nil {
		return nil
	}
	return CapabilityError{Name: name, Op: op, Err: err}
}

// Error returns a capability-scoped error message.
func (e CapabilityError) Error() string {
	if e.Op == "" {
		return fmt.Sprintf("%s: %v", e.Name, e.Err)
	}
	return fmt.Sprintf("%s %s: %v", e.Name, e.Op, e.Err)
}

// Unwrap returns the underlying error.
func (e CapabilityError) Unwrap() error {
	return e.Err
}

// UnavailableError returns an aggregate error for unavailable capabilities.
type UnavailableError struct {
	Statuses []CapabilityStatus
}

// Error reports every unavailable capability by name.
func (e UnavailableError) Error() string {
	if len(e.Statuses) == 0 {
		return "no unavailable capabilities"
	}
	message := "unavailable capabilities:"
	for _, status := range e.Statuses {
		message += " " + status.Name
		if status.Message != "" {
			message += "(" + status.Message + ")"
		}
	}
	return message
}

// ReadyError returns nil only when every enabled status is available.
func ReadyError(statuses []CapabilityStatus) error {
	var unavailable []CapabilityStatus
	for _, status := range statuses {
		if status.Enabled && !status.Available {
			unavailable = append(unavailable, status)
		}
	}
	if len(unavailable) == 0 {
		return nil
	}
	return UnavailableError{Statuses: unavailable}
}

// JoinCapabilityErrors joins non-nil errors into one error.
func JoinCapabilityErrors(errs ...error) error {
	var nonNil []error
	for _, err := range errs {
		if err != nil {
			nonNil = append(nonNil, err)
		}
	}
	return errors.Join(nonNil...)
}
