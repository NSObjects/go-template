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

// UnavailableCapabilityProviderError blocks assembly when a selected provider
// cannot satisfy a required capability.
type UnavailableCapabilityProviderError struct {
	Module     string
	Capability string
	Provider   string
}

func (e *UnavailableCapabilityProviderError) Error() string {
	return fmt.Sprintf("module %q requires capability %q provider %q, but it is unavailable", e.Module, e.Capability, e.Provider)
}

// MissingCapabilityValueError blocks wiring when the selected capability
// provider has no runtime value of the requested type.
type MissingCapabilityValueError struct {
	Module     string
	Capability string
	Provider   string
}

func (e *MissingCapabilityValueError) Error() string {
	return fmt.Sprintf("module %q requires capability %q provider %q, but it has no matching runtime value", e.Module, e.Capability, e.Provider)
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
