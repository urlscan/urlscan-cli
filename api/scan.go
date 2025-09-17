package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/samber/mo"
)

type ScanResult struct {
	UUID string          `json:"uuid"`
	Raw  json.RawMessage `json:"-"`
}

func (r *ScanResult) PrettyJSON() string {
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

func (c *Client) NewScanRequest(url string, opts ...ScanOption) *Request {
	scanOptions := newScanOptions(url, opts...)
	marshalled, err := json.Marshal(scanOptions)
	if err != nil {
		panic(fmt.Sprintf("error marshalling scan options: %s", err))
	}

	return c.NewRequest().
		SetPath(PrefixedPath("/scan/")).
		SetMethod(http.MethodPost).
		SetBodyJSONBytes(marshalled)
}

func (c *Client) Scan(url string, options ...ScanOption) (*ScanResult, error) {
	req := c.NewScanRequest(url, options...)
	resp, err := req.Do()
	if err != nil {
		return nil, err
	}

	var r ScanResult
	err = resp.Unmarshal(&r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (c *Client) NewBatchScanTask(url string, opts ...ScanOption) BatchTask[*Response] {
	return func(c *Client, ctx context.Context) mo.Result[*Response] {
		req := c.NewScanRequest(url, opts...)
		resp, err := req.Do()
		if err != nil {
			return mo.Err[*Response](err)
		}
		return mo.Ok(resp)
	}
}

func (c *Client) NewBatchScanWithWaitTask(url string, maxWait int, opts ...ScanOption) BatchTask[*Response] {
	return func(c *Client, ctx context.Context) mo.Result[*Response] {
		scanReq := c.NewScanRequest(url, opts...)
		scanResp, err := scanReq.Do()
		if err != nil {
			return mo.Err[*Response](err)
		}

		var scanResult ScanResult
		err = scanResp.Unmarshal(&scanResult)
		if err != nil {
			return mo.Err[*Response](err)
		}

		waitResp, err := c.WaitAndGetResult(ctx, scanResult.UUID, maxWait)
		if err != nil {
			return mo.Err[*Response](err)
		}

		return mo.Ok(waitResp)
	}
}
