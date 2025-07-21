package utils

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/samber/lo"
	"github.com/samber/mo"
	api "github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/pkg/version"
)

type APIClient struct {
	*api.Client
}

func NewAPIClient() (*APIClient, error) {
	// api key loading precedence:
	// 1. Environment variable (URLSCAN_API_KEY)
	// 2. Keyring
	key := os.Getenv("URLSCAN_API_KEY")
	if key == "" {
		got, err := NewKeyManager().GetKey()
		if err != nil {
			return nil, err
		}
		key = got
	}

	c := api.NewClient(key)
	c.Agent = fmt.Sprintf("urlscan-cli %s", version.Version)

	return &APIClient{c}, nil
}

type DownloadOptions struct {
	client *APIClient
	url    *url.URL
	output string
	force  bool
	silent bool
}

type DownloadOption = func(*DownloadOptions)

func WithDownloadClient(client *APIClient) DownloadOption {
	return func(opts *DownloadOptions) {
		opts.client = client
	}
}

func WithDownloadURL(url *url.URL) DownloadOption {
	return func(opts *DownloadOptions) {
		opts.url = url
	}
}

func WithDownloadOutput(output string) DownloadOption {
	return func(opts *DownloadOptions) {
		opts.output = output
	}
}

func WithDownloadForce(force bool) DownloadOption {
	return func(opts *DownloadOptions) {
		opts.force = force
	}
}

func WithDownloadDOM(uuid string) DownloadOption {
	return func(opts *DownloadOptions) {
		opts.url = api.URL("%s", fmt.Sprintf("/dom/%s/", uuid))
	}
}

func WithDownloadScreenshot(uuid string) DownloadOption {
	return func(opts *DownloadOptions) {
		opts.url = api.URL("%s", fmt.Sprintf("/screenshots/%s.png", uuid))
	}
}

func WithDownloadSilent(silent bool) DownloadOption {
	return func(opts *DownloadOptions) {
		opts.silent = silent
	}
}

func NewDownloadOptions(opts ...DownloadOption) *DownloadOptions {
	downloadOpts := &DownloadOptions{}
	for _, o := range opts {
		o(downloadOpts)
	}
	return downloadOpts
}

func DownloadWithSpinner(opts *DownloadOptions) error {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Start()

	if !opts.force {
		if fileExists(opts.output) {
			s.Stop()
			return fmt.Errorf("%s already exists", opts.output)
		}
	}

	_, err := opts.client.Download(opts.url, opts.output)
	if err != nil {
		return err
	}

	// stop the spinner
	s.Stop()

	msg := fmt.Sprintf("Downloaded: %s from %s\n", opts.output, opts.url.String())
	if opts.silent {
		// output it to stderr to make the rest of stdout clean (for piping with jq, etc.)
		fmt.Fprint(os.Stderr, msg)
	} else {
		fmt.Print(msg)
	}

	return nil
}

type BatchJSONResultPair struct {
	Key    string          `json:"key"`
	Result json.RawMessage `json:"result"`
}

func NewBatchJSONResultPairs(keys []string, results []mo.Result[*json.RawMessage]) []*BatchJSONResultPair {
	return lo.ZipBy2(keys, results, func(url string, result mo.Result[*json.RawMessage]) *BatchJSONResultPair {
		return &BatchJSONResultPair{
			Key:    url,
			Result: *api.BatchJSONResultToRaw(&result),
		}
	})
}
