package module

import "fmt"

// MissingCapabilityError blocks assembly when a business module requires a
// capability that no enabled provider can satisfy.
type MissingCapabilityError struct {
	Module     string
	Capability string
}

func (e *MissingCapabilityError) Error() string {
	return fmt.Sprintf("module %q requires missing capability %q", e.Module, e.Capability)
}

// UnsupportedEntryPointError blocks assembly when no platform adapter can
// expose a declared entry point type.
type UnsupportedEntryPointError struct {
	Module         string
	EntryPointType string
}

func (e *UnsupportedEntryPointError) Error() string {
	return fmt.Sprintf("module %q declares unsupported entry point type %q", e.Module, e.EntryPointType)
}
