package generator

func (g *Generator) ensureSupportFiles() error {
	// Support files are now provided by go-kit, so we don't need to generate them locally.
	return nil
}
