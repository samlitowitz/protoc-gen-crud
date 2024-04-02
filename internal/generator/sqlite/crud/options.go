package crud

type options struct {
	formatOutput bool
}

type Option interface {
	apply(*options)
}

type formatOutputOption bool

func (f formatOutputOption) apply(opts *options) {
	opts.formatOutput = bool(f)
}

func WithFormatOutput(f bool) Option {
	return formatOutputOption(f)
}
