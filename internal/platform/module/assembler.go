package module

// Option customizes assembly validation.
type Option func(*assemblerConfig)

type assemblerConfig struct {
	entryPointAdapters   map[string]struct{}
	capabilitySelections map[string]string
}

type assemblyResult struct {
	Assembly
	capabilityValues map[capabilityProvider]any
}

type capabilityProvider struct {
	capability string
	provider   string
}

// Assembly contains the observable report plus runtime assembly metadata.
type Assembly struct {
	report                          Report
	selectedCapabilityModuleIndexes []int
}

// Report returns the observable assembly report.
func (a Assembly) Report() Report {
	return a.report
}

// SelectedCapabilityModuleIndexes returns input module indexes for selected capability providers.
func (a Assembly) SelectedCapabilityModuleIndexes() []int {
	return append([]int(nil), a.selectedCapabilityModuleIndexes...)
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
	return result.Report(), nil
}

// AssembleRuntime validates module declarations and returns runtime assembly metadata.
func AssembleRuntime(modules []Module, options ...Option) (Assembly, error) {
	result, err := assemble(modules, options...)
	if err != nil {
		return Assembly{}, err
	}
	return result.Assembly, nil
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
	capabilityModuleIndexes := make(map[capabilityProvider]int)
	capabilityProviderModules := make(map[capabilityProvider]string)
	enabledCapabilities := make(map[string][]CapabilityStatus)

	for moduleIndex, mod := range modules {
		descriptor := mod.Descriptor()
		descriptors = append(descriptors, descriptor)
		report.ActiveModules = append(report.ActiveModules, descriptor.Name)

		for _, capability := range descriptor.Provides {
			provider := capability.Provider
			if provider == "" {
				provider = descriptor.Name
			}
			providerKey := capabilityProvider{
				capability: capability.Name,
				provider:   provider,
			}
			if firstModule, ok := capabilityProviderModules[providerKey]; ok {
				return assemblyResult{}, &DuplicateCapabilityProviderError{
					Capability:   capability.Name,
					Provider:     provider,
					FirstModule:  firstModule,
					SecondModule: descriptor.Name,
				}
			}
			capabilityProviderModules[providerKey] = descriptor.Name
			status := CapabilityStatus{
				Name:     capability.Name,
				Status:   capability.Status,
				Provider: provider,
				Default:  capability.Default,
			}
			report.Capabilities = append(report.Capabilities, status)
			if descriptor.Kind == CapabilityModule {
				capabilityModuleIndexes[providerKey] = moduleIndex
			}
			if capability.Value != nil {
				capabilityValues[providerKey] = capability.Value
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
			status, resolution := resolveCapabilityProvider(
				enabledCapabilities[requirement.Name],
				selectedProvider,
			)
			if !resolution.ok {
				if resolution.ambiguous {
					return assemblyResult{}, &AmbiguousCapabilityProviderError{
						Module:     descriptor.Name,
						Capability: requirement.Name,
						Providers:  capabilityProviders(enabledCapabilities[requirement.Name]),
					}
				}
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

	selectedCapabilityModuleIndexes := selectedCapabilityModuleIndexes(report, capabilityModuleIndexes)
	return assemblyResult{
		Assembly: Assembly{
			report:                          report,
			selectedCapabilityModuleIndexes: selectedCapabilityModuleIndexes,
		},
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
		requirement, _ := result.Report().Requirement(moduleName, capabilityName)
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

	requirement, ok := result.Report().Requirement(moduleName, capabilityName)
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

func selectedCapabilityModuleIndexes(report Report, capabilityModuleIndexes map[capabilityProvider]int) []int {
	selected := make([]int, 0, len(report.Requirements))
	seen := make(map[int]struct{})
	for _, requirement := range report.Requirements {
		if !requirement.Satisfied {
			continue
		}
		index, ok := capabilityModuleIndexes[capabilityProvider{
			capability: requirement.Capability,
			provider:   requirement.Provider,
		}]
		if !ok {
			continue
		}
		if _, ok := seen[index]; ok {
			continue
		}
		seen[index] = struct{}{}
		selected = append(selected, index)
	}
	return selected
}

type capabilityProviderResolution struct {
	ok        bool
	ambiguous bool
}

func resolveCapabilityProvider(capabilities []CapabilityStatus, selectedProvider string) (CapabilityStatus, capabilityProviderResolution) {
	if selectedProvider != "" {
		for _, capability := range capabilities {
			if capability.Provider == selectedProvider {
				return capability, capabilityProviderResolution{ok: true}
			}
		}
		return CapabilityStatus{}, capabilityProviderResolution{}
	}
	return defaultCapabilityProvider(capabilities)
}

func defaultCapabilityProvider(capabilities []CapabilityStatus) (CapabilityStatus, capabilityProviderResolution) {
	if len(capabilities) == 0 {
		return CapabilityStatus{}, capabilityProviderResolution{}
	}
	var defaultCapability CapabilityStatus
	defaultCount := 0
	for _, capability := range capabilities {
		if capability.Default {
			defaultCapability = capability
			defaultCount++
		}
	}
	if defaultCount == 1 {
		return defaultCapability, capabilityProviderResolution{ok: true}
	}
	if defaultCount > 1 {
		return CapabilityStatus{}, capabilityProviderResolution{ambiguous: true}
	}
	if len(capabilities) == 1 {
		return capabilities[0], capabilityProviderResolution{ok: true}
	}
	return CapabilityStatus{}, capabilityProviderResolution{ambiguous: true}
}

func capabilityProviders(capabilities []CapabilityStatus) []string {
	providers := make([]string, 0, len(capabilities))
	for _, capability := range capabilities {
		providers = append(providers, capability.Provider)
	}
	return providers
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
