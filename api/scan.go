package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/samber/mo"
)

type ScanResult struct {
	UUID string          `json:"uuid"`
	Raw  json.RawMessage `json:"-"`
}

func (r *ScanResult) PrettyJson() string {
	var jsonBody bytes.Buffer
	err := json.Indent(&jsonBody, r.Raw, "", "  ")
	if err != nil {
		msg := fmt.Sprintf("error formatting JSON response: %s", err)
		panic(msg)
	}
	return jsonBody.String()
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

func (cli *Client) Scan(url string, options ...ScanOption) (*ScanResult, error) {
	scanOptions := newScanOptions(url, options...)

	marshalled, err := json.Marshal(scanOptions)
	if err != nil {
		return nil, err
	}

	resp, err := cli.Post(URL("/api/v1/scan/"), &Request{
		Raw: json.RawMessage(marshalled),
	})
	if err != nil {
		return nil, err
	}

	r := &ScanResult{}
	err = json.Unmarshal(resp.Raw, r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (cli *Client) NewBatchScanTask(url string, opts ...ScanOption) BatchTask[*json.RawMessage] {
	return func(cli *Client, ctx context.Context) mo.Result[*json.RawMessage] {
		r, err := cli.Scan(url, opts...)
		if err != nil {
			return mo.Err[*json.RawMessage](err)
		}
		return mo.Ok(&r.Raw)
	}
}

func (cli *Client) NewBatchScanWitWaitTask(url string, maxWait int, opts ...ScanOption) BatchTask[*json.RawMessage] {
	return func(cli *Client, ctx context.Context) mo.Result[*json.RawMessage] {
		scanResult, err := cli.Scan(url, opts...)
		if err != nil {
			return mo.Err[*json.RawMessage](err)
		}

		waitResult, err := cli.WaitAndGetResult(ctx, scanResult.UUID, maxWait)
		if err != nil {
			return mo.Err[*json.RawMessage](err)
		}
		return mo.Ok(&waitResult.Raw)
	}
}
