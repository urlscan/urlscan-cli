package api

type IncidentOptions struct {
	Incident struct {
		// common options
		ExpireAfter                int      `json:"expireAfter"`
		ScanInterval               int      `json:"scanInterval"`
		ScanIntervalMode           string   `json:"scanIntervalMode"`
		WatchedAttributes          []string `json:"watchedAttributes"`
		UserAgents                 []string `json:"userAgents"`
		UserAgentsPerInterval      int      `json:"userAgentsPerInterval"`
		Countries                  []string `json:"countries"`
		CountriesPerInterval       int      `json:"countriesPerInterval"`
		StopDelaySuspended         int      `json:"stopDelaySuspended"`
		StopDelayInactive          int      `json:"stopDelayInactive"`
		StopDelayMalicious         int      `json:"stopDelayMalicious"`
		ScanIntervalAfterSuspended int      `json:"scanIntervalAfterSuspended"`
		ScanIntervalAfterMalicious int      `json:"scanIntervalAfterMalicious"`
		// common fields
		Visibility      string `json:"visibility"`
		IncidentProfile string `json:"incidentProfile,omitempty"`
		ExpireAt        string `json:"expireAt,omitempty"`
		// create & updated fields
		Channels   []string `json:"channels"`
		Observable string   `json:"observable"`
	} `json:"incident"`
}

type IncidentOption func(*IncidentOptions)

func WithIncidentExpireAfter(expireAfter int) IncidentOption {
	return func(opts *IncidentOptions) {
		opts.Incident.ExpireAfter = expireAfter
	}
}

func WithIncidentScanInterval(scanInterval int) IncidentOption {
	return func(opts *IncidentOptions) {
		opts.Incident.ScanInterval = scanInterval
	}
}

func WithIncidentScanIntervalMode(scanIntervalMode string) IncidentOption {
	return func(opts *IncidentOptions) {
		opts.Incident.ScanIntervalMode = scanIntervalMode
	}
}

func WithIncidentWatchedAttributes(watchedAttributes []string) IncidentOption {
	return func(opts *IncidentOptions) {
		opts.Incident.WatchedAttributes = watchedAttributes
	}
}

func WithIncidentUserAgents(userAgents []string) IncidentOption {
	return func(opts *IncidentOptions) {
		opts.Incident.UserAgents = userAgents
	}
}

func WithIncidentUserAgentsPerInterval(userAgentsPerInterval int) IncidentOption {
	return func(opts *IncidentOptions) {
		opts.Incident.UserAgentsPerInterval = userAgentsPerInterval
	}
}

func WithIncidentCountries(countries []string) IncidentOption {
	return func(opts *IncidentOptions) {
		opts.Incident.Countries = countries
	}
}

func WithIncidentCountriesPerInterval(countriesPerInterval int) IncidentOption {
	return func(opts *IncidentOptions) {
		opts.Incident.CountriesPerInterval = countriesPerInterval
	}
}

func WithIncidentStopDelaySuspended(stopDelaySuspended int) IncidentOption {
	return func(opts *IncidentOptions) {
		opts.Incident.StopDelaySuspended = stopDelaySuspended
	}
}

func WithIncidentStopDelayInactive(stopDelayInactive int) IncidentOption {
	return func(opts *IncidentOptions) {
		opts.Incident.StopDelayInactive = stopDelayInactive
	}
}

func WithIncidentStopDelayMalicious(stopDelayMalicious int) IncidentOption {
	return func(opts *IncidentOptions) {
		opts.Incident.StopDelayMalicious = stopDelayMalicious
	}
}

func WithIncidentScanIntervalAfterSuspended(scanIntervalAfterSuspended int) IncidentOption {
	return func(opts *IncidentOptions) {
		opts.Incident.ScanIntervalAfterSuspended = scanIntervalAfterSuspended
	}
}

func WithIncidentScanIntervalAfterMalicious(scanIntervalAfterMalicious int) IncidentOption {
	return func(opts *IncidentOptions) {
		opts.Incident.ScanIntervalAfterMalicious = scanIntervalAfterMalicious
	}
}

func WithIncidentVisibility(visibility string) IncidentOption {
	return func(opts *IncidentOptions) {
		opts.Incident.Visibility = visibility
	}
}

func WithIncidentProfile(incidentProfile string) IncidentOption {
	return func(opts *IncidentOptions) {
		opts.Incident.IncidentProfile = incidentProfile
	}
}

func WithIncidentExpireAt(expireAt string) IncidentOption {
	return func(opts *IncidentOptions) {
		opts.Incident.ExpireAt = expireAt
	}
}

func WithIncidentChannels(channels []string) IncidentOption {
	return func(opts *IncidentOptions) {
		opts.Incident.Channels = channels
	}
}

func WithIncidentObservable(observable string) IncidentOption {
	return func(opts *IncidentOptions) {
		opts.Incident.Observable = observable
	}
}

func newIncidentOptions(opts ...IncidentOption) *IncidentOptions {
	options := &IncidentOptions{}

	for _, opt := range opts {
		opt(options)
	}

	return options
}
