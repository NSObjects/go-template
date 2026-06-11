package module

// Option customizes assembly validation.
type Option func(*assemblerConfig)

type assemblerConfig struct {
	entryPointAdapters map[string]struct{}
}

// WithEntryPointAdapters declares the entry point types the platform can expose.
func WithEntryPointAdapters(entryPointTypes ...string) Option {
	return func(cfg *assemblerConfig) {
		for _, entryPointType := range entryPointTypes {
			cfg.entryPointAdapters[entryPointType] = struct{}{}
		}
	}
}

// Assemble validates module declarations and returns an observable report.
func Assemble(modules []Module, options ...Option) (Report, error) {
	cfg := assemblerConfig{entryPointAdapters: make(map[string]struct{})}
	for _, option := range options {
		option(&cfg)
	}

	descriptors := make([]Descriptor, 0, len(modules))
	report := Report{}
	enabledCapabilities := make(map[string]CapabilityStatus)

	for _, mod := range modules {
		descriptor := mod.Descriptor()
		descriptors = append(descriptors, descriptor)
		report.ActiveModules = append(report.ActiveModules, descriptor.Name)

		for _, capability := range descriptor.Provides {
			status := CapabilityStatus{
				Name:     capability.Name,
				Status:   capability.Status,
				Provider: descriptor.Name,
			}
			report.Capabilities = append(report.Capabilities, status)
			if canSatisfyRequirement(capability.Status) {
				enabledCapabilities[capability.Name] = status
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
			status, ok := enabledCapabilities[requirement.Name]
			if !ok {
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

func canSatisfyRequirement(status CapabilityState) bool {
	return status == CapabilityEnabled || status == CapabilityHealthy
}
