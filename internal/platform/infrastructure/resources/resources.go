package resources

import (
	"context"
	"errors"
)

// Component adapts one concrete capability into the fixed resources bundle.
type Component struct {
	Name  string
	Check func(context.Context) CapabilityStatus
	Close func(context.Context) error
}

// Resources aggregates fixed process-level capability status and lifecycle.
type Resources struct {
	components []Component
}

// New creates a fixed resources bundle.
func New(components ...Component) *Resources {
	copied := make([]Component, 0, len(components))
	for _, component := range components {
		if component.Name == "" {
			continue
		}
		copied = append(copied, component)
	}
	return &Resources{components: copied}
}

// Status returns the observable state for every configured component.
func (r *Resources) Status(ctx context.Context) []CapabilityStatus {
	if r == nil {
		return nil
	}
	statuses := make([]CapabilityStatus, 0, len(r.components))
	for _, component := range r.components {
		if component.Check == nil {
			statuses = append(statuses, Disabled(component.Name))
			continue
		}
		statuses = append(statuses, component.Check(ctx))
	}
	return statuses
}

// Ready returns an error when any enabled component is unavailable.
func (r *Resources) Ready(ctx context.Context) error {
	return ReadyError(r.Status(ctx))
}

// Close closes each component and returns capability-scoped shutdown failures.
func (r *Resources) Close(ctx context.Context) error {
	if r == nil {
		return nil
	}
	var errs []error
	for _, component := range r.components {
		if component.Close == nil {
			continue
		}
		if err := component.Close(ctx); err != nil {
			errs = append(errs, NewCapabilityError(component.Name, "close", err))
		}
	}
	return errors.Join(errs...)
}
