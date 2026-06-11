package module

// Option customizes assembly validation.
type Option func(*assemblerConfig)

type assemblerConfig struct {
	entryPointAdapters   map[string]struct{}
	capabilitySelections map[string]string
}

type assemblyResult struct {
	report           Report
	capabilityValues map[capabilityProvider]any
}

type capabilityProvider struct {
	capability string
	provider   string
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

// WithCapabilityProviderSelections declares preferred providers from a config map.
func WithCapabilityProviderSelections(selections map[string]string) Option {
	return func(cfg *assemblerConfig) {
		for capability, provider := range selections {
			cfg.capabilitySelections[capability] = provider
		}
	}
}

// Assemble validates module declarations and returns an observable report.
func Assemble(modules []Module, options ...Option) (Report, error) {
	result, err := assemble(modules, options...)
	if err != nil {
		return Report{}, err
	}
	return result.report, nil
}

func assemble(modules []Module, options ...Option) (assemblyResult, error) {
	cfg := assemblerConfig{
		entryPointAdapters:   make(map[string]struct{}),
		capabilitySelections: make(map[string]string),
	}
	for _, option := range options {
		option(&cfg)
	}

	descriptors := make([]Descriptor, 0, len(modules))
	report := Report{}
	capabilityValues := make(map[capabilityProvider]any)
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
			if capability.Value != nil {
				capabilityValues[capabilityProvider{
					capability: capability.Name,
					provider:   provider,
				}] = capability.Value
			}
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
					return assemblyResult{}, &UnavailableCapabilityProviderError{
						Module:     descriptor.Name,
						Capability: requirement.Name,
						Provider:   selectedProvider,
					}
				}
				return assemblyResult{}, &MissingCapabilityError{
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
				return assemblyResult{}, &UnsupportedEntryPointError{
					Module:         descriptor.Name,
					EntryPointType: entryPoint.Type,
				}
			}
		}
	}

	return assemblyResult{
		report:           report,
		capabilityValues: capabilityValues,
	}, nil
}

// ResolveCapabilityValueFromModules returns the runtime value selected for one module requirement.
func ResolveCapabilityValueFromModules[T any](
	providers []Module,
	moduleName string,
	capabilityName string,
	options ...Option,
) (T, error) {
	modules := make([]Module, 0, len(providers)+1)
	modules = append(modules, staticDescriptorModule{descriptor: Descriptor{
		Name: moduleName,
		Kind: BusinessModule,
		Requires: []CapabilityRef{
			{Name: capabilityName},
		},
	}})
	modules = append(modules, providers...)

	result, err := assemble(modules, options...)
	if err != nil {
		var zero T
		return zero, err
	}
	value, ok := resolveCapabilityValue[T](result, moduleName, capabilityName)
	if !ok {
		var zero T
		requirement, _ := result.report.Requirement(moduleName, capabilityName)
		return zero, &MissingCapabilityValueError{
			Module:     moduleName,
			Capability: capabilityName,
			Provider:   requirement.Provider,
		}
	}
	return value, nil
}

func resolveCapabilityValue[T any](result assemblyResult, moduleName, capabilityName string) (T, bool) {
	var zero T

	requirement, ok := result.report.Requirement(moduleName, capabilityName)
	if !ok || !requirement.Satisfied {
		return zero, false
	}
	value, ok := result.capabilityValues[capabilityProvider{
		capability: capabilityName,
		provider:   requirement.Provider,
	}]
	if !ok {
		return zero, false
	}
	typed, ok := value.(T)
	if !ok {
		return zero, false
	}
	return typed, true
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

type staticDescriptorModule struct {
	descriptor Descriptor
}

func (m staticDescriptorModule) Descriptor() Descriptor {
	return m.descriptor
}
