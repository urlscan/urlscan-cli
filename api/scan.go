package api

import "encoding/json"

type ScanResult struct {
	UUID string          `json:"uuid"`
	Raw  json.RawMessage `json:"-"`
}

func (r *ScanResult) UnmarshalJSON(data []byte) error {
	type result ScanResult
	var dst result

	err := json.Unmarshal(data, &dst)
	if err != nil {
		return err
	}
	*r = ScanResult(dst)
	r.Raw = data
	return err
}

type ScanOptions struct {
	URL            string    `json:"url"`
	CustomAgent    *string   `json:"customagent,omitempty"`
	Referer        *string   `json:"referer,omitempty"`
	Visibility     *string   `json:"visibility,omitempty"`
	Tags           *[]string `json:"tags,omitempty"`
	OverrideSafety *string   `json:"overrideSafety,omitempty"`
	Country        *string   `json:"country,omitempty"`
}

type ScanOption func(*ScanOptions)

func WithScanCustomAgent(customAgent string) ScanOption {
	return func(opts *ScanOptions) {
		if customAgent != "" {
			opts.CustomAgent = &customAgent
		}
	}
}

func WithScanReferer(referer string) ScanOption {
	return func(opts *ScanOptions) {
		if referer != "" {
			opts.Referer = &referer
		}
	}
}

func WithScanVisibility(visibility string) ScanOption {
	return func(opts *ScanOptions) {
		if visibility != "" {
			opts.Visibility = &visibility
		}
	}
}

func WithScanTags(tags []string) ScanOption {
	return func(opts *ScanOptions) {
		if len(tags) > 0 {
			opts.Tags = &tags
		}
	}
}

func WithScanOverrideSafety(overrideSafety string) ScanOption {
	return func(opts *ScanOptions) {
		if overrideSafety != "" {
			opts.OverrideSafety = &overrideSafety
		}
	}
}

func WithScanCountry(country string) ScanOption {
	return func(opts *ScanOptions) {
		if country != "" {
			opts.Country = &country
		}
	}
}

func newScanOptions(url string, opts ...ScanOption) *ScanOptions {
	opt := &ScanOptions{
		URL:            url,
		CustomAgent:    nil,
		Referer:        nil,
		Visibility:     nil,
		Tags:           nil,
		OverrideSafety: nil,
		Country:        nil,
	}
	for _, o := range opts {
		o(opt)
	}
	return opt
}
