// Package module contains the small assembly model used to compose business,
// platform, and capability modules before the application starts.
package module

const EntryPointHTTP = "http"

// ModuleKind identifies the role a module plays during application assembly.
type ModuleKind string

const (
	BusinessModule   ModuleKind = "business"
	CapabilityModule ModuleKind = "capability"
	PlatformModule   ModuleKind = "platform"
)

// CapabilityState describes whether a capability can satisfy business needs.
type CapabilityState string

const (
	CapabilityEnabled     CapabilityState = "enabled"
	CapabilityDisabled    CapabilityState = "disabled"
	CapabilityHealthy     CapabilityState = "healthy"
	CapabilityUnavailable CapabilityState = "unavailable"
)

// Module describes one unit that can participate in application assembly.
type Module interface {
	Descriptor() Descriptor
}

// Descriptor is the assembly-time declaration made by a module.
type Descriptor struct {
	Name        string
	Kind        ModuleKind
	Requires    []CapabilityRef
	Provides    []Capability
	EntryPoints []EntryPoint
}

// CapabilityRef names a capability required by a business module.
type CapabilityRef struct {
	Name string
}

// Capability names a capability provided by a capability module.
type Capability struct {
	Name   string
	Status CapabilityState
}

// CapabilityStatus is the reportable state of one capability.
type CapabilityStatus struct {
	Name     string
	Status   CapabilityState
	Provider string
}

// EntryPoint declares an externally reachable trigger owned by a module.
type EntryPoint struct {
	Owner string
	Type  string
	Name  string
	Value any
}
