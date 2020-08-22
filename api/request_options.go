package api

// RequestOptions represent the available options on the service requests
type RequestOptions struct {
	cache  bool
	term   string
	status string
}

func (options *RequestOptions) apply(opts []RequestOption) {
	for _, o := range opts {
		o.apply(options)
	}
}

// RequestOption overrides behavior of the service requests
type RequestOption interface {
	apply(*RequestOptions)
}

type optionFunc func(*RequestOptions)

func (f optionFunc) apply(o *RequestOptions) {
	f(o)
}

// WithStatus returns requests matching given status
func WithStatus(status string) RequestOption {
	return optionFunc(func(o *RequestOptions) {
		o.status = status
	})
}

// WithCache returns requests using cache
func WithCache(cache bool) RequestOption {
	return optionFunc(func(o *RequestOptions) {
		o.cache = cache
	})
}

// WithTerm returns requests matching the search term
func WithTerm(term string) RequestOption {
	return optionFunc(func(o *RequestOptions) {
		o.term = term
	})
}

// NewRequestOptions returns RequestOptions with defaults
func NewRequestOptions(opts []RequestOption) *RequestOptions {
	options := &RequestOptions{
		cache: true,
	}
	for _, o := range opts {
		o.apply(options)
	}
	return options
}
