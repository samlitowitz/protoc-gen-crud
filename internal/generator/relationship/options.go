package relationship

type options struct {
	optionsPkg string
}

type Option interface {
	apply(*options)
}

type optionsPkgOption string

func (o optionsPkgOption) apply(opts *options) {
	opts.optionsPkg = string(o)
}

func WithOptionsPkg(o string) Option {
	return optionsPkgOption(o)
}
