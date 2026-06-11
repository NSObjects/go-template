package module

// Option customizes assembly validation.
type Option func(*assemblerConfig)

type assemblerConfig struct {
	entryPointAdapters   map[string]struct{}
	capabilitySelections map[string]string
}

// WithEntryPointAdapters declares the entry point types the platform can expose.
func WithEntryPointAdapters(entryPointTypes ...string) Option {
	return func(cfg *assemblerConfig) {
		for _, entryPointType := range entryPointTypes {
			cfg.entryPointAdapters[entryPointType] = struct{}{}
		}
	}
}

// WithCapabilitySelections declares preferred providers for capabilities.
func WithCapabilitySelections(selections ...CapabilitySelection) Option {
	return func(cfg *assemblerConfig) {
		for _, selection := range selections {
			cfg.capabilitySelections[selection.Capability] = selection.Provider
		}
	}
}

// Assemble validates module declarations and returns an observable report.
func Assemble(modules []Module, options ...Option) (Report, error) {
	cfg := assemblerConfig{
		entryPointAdapters:   make(map[string]struct{}),
		capabilitySelections: make(map[string]string),
	}
	for _, option := range options {
		option(&cfg)
	}

	descriptors := make([]Descriptor, 0, len(modules))
	report := Report{}
	enabledCapabilities := make(map[string][]CapabilityStatus)

	for _, mod := range modules {
		descriptor := mod.Descriptor()
		descriptors = append(descriptors, descriptor)
		report.ActiveModules = append(report.ActiveModules, descriptor.Name)

		for _, capability := range descriptor.Provides {
			provider := capability.Provider
			if provider == "" {
				provider = descriptor.Name
			}
			status := CapabilityStatus{
				Name:     capability.Name,
				Status:   capability.Status,
				Provider: provider,
				Default:  capability.Default,
			}
			report.Capabilities = append(report.Capabilities, status)
			if canSatisfyRequirement(capability.Status) {
				enabledCapabilities[capability.Name] = append(enabledCapabilities[capability.Name], status)
			}
		}

		for _, entryPoint := range descriptor.EntryPoints {
			if entryPoint.Owner == "" {
				entryPoint.Owner = descriptor.Name
			}
			report.EntryPoints = append(report.EntryPoints, entryPoint)
		}

		if descriptor.Kind == BusinessModule && len(descriptor.EntryPoints) == 0 {
			report.Warnings = append(report.Warnings, Warning{
				Code:    WarningNoEntryPoints,
				Module:  descriptor.Name,
				Message: "business module declares no entry points",
			})
		}
	}

	for _, descriptor := range descriptors {
		for _, requirement := range descriptor.Requires {
			selectedProvider := cfg.capabilitySelections[requirement.Name]
			status, ok := resolveCapabilityProvider(
				enabledCapabilities[requirement.Name],
				selectedProvider,
			)
			if !ok {
				if selectedProvider != "" {
					return Report{}, &UnavailableCapabilityProviderError{
						Module:     descriptor.Name,
						Capability: requirement.Name,
						Provider:   selectedProvider,
					}
				}
				return Report{}, &MissingCapabilityError{
					Module:     descriptor.Name,
					Capability: requirement.Name,
				}
			}
			report.Requirements = append(report.Requirements, RequirementStatus{
				Module:     descriptor.Name,
				Capability: requirement.Name,
				Provider:   status.Provider,
				Satisfied:  true,
			})
		}

		for _, entryPoint := range descriptor.EntryPoints {
			if len(cfg.entryPointAdapters) == 0 {
				continue
			}
			if _, ok := cfg.entryPointAdapters[entryPoint.Type]; !ok {
				return Report{}, &UnsupportedEntryPointError{
					Module:         descriptor.Name,
					EntryPointType: entryPoint.Type,
				}
			}
		}
	}

	return report, nil
}

func resolveCapabilityProvider(capabilities []CapabilityStatus, selectedProvider string) (CapabilityStatus, bool) {
	if selectedProvider != "" {
		for _, capability := range capabilities {
			if capability.Provider == selectedProvider {
				return capability, true
			}
		}
		return CapabilityStatus{}, false
	}
	return defaultCapabilityProvider(capabilities)
}

func defaultCapabilityProvider(capabilities []CapabilityStatus) (CapabilityStatus, bool) {
	if len(capabilities) == 0 {
		return CapabilityStatus{}, false
	}
	for _, capability := range capabilities {
		if capability.Default {
			return capability, true
		}
	}
	if len(capabilities) == 1 {
		return capabilities[0], true
	}
	return CapabilityStatus{}, false
}

func canSatisfyRequirement(status CapabilityState) bool {
	return status == CapabilityEnabled || status == CapabilityHealthy
}
