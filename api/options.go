package api

type options struct {
	cache  bool
	term   string
	status string
}

func (options *options) apply(opts []Option) {
	for _, o := range opts {
		o.apply(options)
	}
}

// Option overrides behavior of service requests
type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(o *options) {
	f(o)
}

// WithStatus returns requests matching given status
func WithStatus(status string) Option {
	return optionFunc(func(o *options) {
		o.status = status
	})
}

// WithCache returns requests using cache
func WithCache(cache bool) Option {
	return optionFunc(func(o *options) {
		o.cache = cache
	})
}

// WithTerm returns requests matching the search term
func WithTerm(term string) Option {
	return optionFunc(func(o *options) {
		o.term = term
	})
}

// Options returns options with defaults
func Options(opts []Option) *options {
	options := &options{
		cache: true,
	}
	for _, o := range opts {
		o.apply(options)
	}
	return options
}
