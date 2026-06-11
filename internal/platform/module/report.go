package module

const WarningNoEntryPoints = "no_entry_points"

// Warning records non-blocking assembly feedback.
type Warning struct {
	Code    string
	Module  string
	Message string
}

// Report describes the observable result of module assembly.
type Report struct {
	ActiveModules []string
	Capabilities  []CapabilityStatus
	Requirements  []RequirementStatus
	EntryPoints   []EntryPoint
	Warnings      []Warning
}

// HasActiveModule reports whether a module is included in the assembly.
func (r Report) HasActiveModule(name string) bool {
	for _, moduleName := range r.ActiveModules {
		if moduleName == name {
			return true
		}
	}
	return false
}

// Capability returns the reported status for a named capability.
func (r Report) Capability(name string) (CapabilityStatus, bool) {
	for _, capability := range r.Capabilities {
		if capability.Name == name {
			return capability, true
		}
	}
	return CapabilityStatus{}, false
}

// Requirement returns the reported status for one module capability requirement.
func (r Report) Requirement(moduleName, capabilityName string) (RequirementStatus, bool) {
	for _, requirement := range r.Requirements {
		if requirement.Module == moduleName && requirement.Capability == capabilityName {
			return requirement, true
		}
	}
	return RequirementStatus{}, false
}

// RequirementStatus records how a module requirement was satisfied.
type RequirementStatus struct {
	Module     string
	Capability string
	Provider   string
	Satisfied  bool
}
