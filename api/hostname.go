package api

type HostnameOptions struct {
	Limit     int
	PageState string
}

type HostnameOption func(*HostnameOptions)

func WithHostnameLimit(limit int) HostnameOption {
	return func(opts *HostnameOptions) {
		opts.Limit = limit
	}
}

func WithHostnamePageState(pageState string) HostnameOption {
	return func(opts *HostnameOptions) {
		opts.PageState = pageState
	}
}

func newHostnameOptions(opts ...HostnameOption) *HostnameOptions {
	options := &HostnameOptions{}

	for _, opt := range opts {
		opt(options)
	}

	return options
}
