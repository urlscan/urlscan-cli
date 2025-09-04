package api

import (
	"encoding/json"
	"fmt"
)

type LiveScanOptions struct {
	Task struct {
		URL        string `json:"url"`
		Visibility string `json:"visibility"`
	} `json:"task"`
	Scanner struct {
		PageTimeout     int               `json:"pageTimeout"`
		CaptureDelay    int               `json:"captureDelay"`
		ExtraHeaders    map[string]string `json:"extraHeaders,omitempty"`
		EnableFeatures  []string          `json:"enableFeatures,omitempty"`
		DisableFeatures []string          `json:"disableFeatures,omitempty"`
	} `json:"scanner"`
}

type LiveScanOption func(*LiveScanOptions)

func WithLiveScanTaskURL(url string) LiveScanOption {
	return func(opts *LiveScanOptions) {
		opts.Task.URL = url
	}
}

func WithLiveScanTaskVisibility(visibility string) LiveScanOption {
	return func(opts *LiveScanOptions) {
		opts.Task.Visibility = visibility
	}
}

func WithLiveScanScannerPageTimeout(timeout int) LiveScanOption {
	return func(opts *LiveScanOptions) {
		opts.Scanner.PageTimeout = timeout
	}
}

func WithLiveScanScannerCaptureDelay(delay int) LiveScanOption {
	return func(opts *LiveScanOptions) {
		opts.Scanner.CaptureDelay = delay
	}
}

func WithLiveScanScannerExtraHeaders(headers map[string]string) LiveScanOption {
	return func(opts *LiveScanOptions) {
		opts.Scanner.ExtraHeaders = headers
	}
}

func WithLiveScanScannerEnableFeatures(features []string) LiveScanOption {
	return func(opts *LiveScanOptions) {
		opts.Scanner.EnableFeatures = features
	}
}

func WithLiveScanScannerDisableFeatures(features []string) LiveScanOption {
	return func(opts *LiveScanOptions) {
		opts.Scanner.DisableFeatures = features
	}
}

func newLiveScanOptions(opts ...LiveScanOption) *LiveScanOptions {
	var o LiveScanOptions
	for _, fn := range opts {
		fn(&o)
	}
	return &o
}

type LiveScanStoreOptions struct {
	Task struct {
		Visibility string `json:"visibility"`
	} `json:"task"`
}

type LiveScanStoreOption func(*LiveScanStoreOptions)

func WithLiveScanStoreTaskVisibility(visibility string) LiveScanStoreOption {
	return func(opts *LiveScanStoreOptions) {
		opts.Task.Visibility = visibility
	}
}

func newLiveScanStoreOptions(opts ...LiveScanStoreOption) *LiveScanStoreOptions {
	var options LiveScanStoreOptions
	for _, opt := range opts {
		opt(&options)
	}
	return &options
}

func (c *Client) TriggerNonBlockingLiveScan(id string, opts ...LiveScanOption) (*Response, error) {
	liveScanOpts := newLiveScanOptions(opts...)
	marshalled, err := json.Marshal(liveScanOpts)
	if err != nil {
		return nil, err
	}
	return c.NewRequest().SetBodyJSONBytes(marshalled).Post(PrefixedPath(fmt.Sprintf("/livescan/%s/task/", id)))
}

func (c *Client) TriggerLiveScan(id string, opts ...LiveScanOption) (*Response, error) {
	liveScanOpts := newLiveScanOptions(opts...)
	marshalled, err := json.Marshal(liveScanOpts)
	if err != nil {
		return nil, err
	}

	return c.NewRequest().SetBodyJSONBytes(marshalled).Post(PrefixedPath(fmt.Sprintf("/livescan/%s/scan/", id)))
}

func (c *Client) StoreLiveScanResult(scannerId string, scanId string, opts ...LiveScanStoreOption) (*Response, error) {
	liveScanStoreOpts := newLiveScanStoreOptions(opts...)
	marshalled, err := json.Marshal(liveScanStoreOpts)
	if err != nil {
		return nil, err
	}
	return c.NewRequest().SetBodyJSONBytes(marshalled).Put(PrefixedPath(fmt.Sprintf("/livescan/%s/%s/", scannerId, scanId)))
}
