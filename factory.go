package belt

type Factory struct {
	items   []I
	options FactoryOptions
}

type FactoryOptions struct {
	sync bool
}

func NewFactory(items []I, options FactoryOptions) *Factory {
	f := Factory{items, options}
	return &f
}

func (f *Factory) ConfigWithSync(isOn bool) *Factory {
	f.options = FactoryOptions{sync: isOn}
	return f
}
