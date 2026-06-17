package boot

// BusinessModules returns the business modules installed by the default runtime.
func BusinessModules() []Module {
	return []Module{
		customerModule(),
		productModule(),
		salesOrderModule(),
	}
}
