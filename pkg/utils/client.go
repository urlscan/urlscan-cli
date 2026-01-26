package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/briandowns/spinner"
	"github.com/samber/lo"
	"github.com/samber/mo"
	"github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/pkg/version"
)

type APIClient struct {
	*api.Client
}

func GetKey() (string, error) {
	// api key loading precedence:
	// 1. Environment variable (URLSCAN_API_KEY)
	// 2. Keyring
	key := os.Getenv("URLSCAN_API_KEY")
	if key == "" {
		got, err := NewKeyManager().GetKey()
		if err != nil {
			return "", err
		}
		key = got
	}
	return key, nil
}

func NewAPIClient() (*APIClient, error) {
	key, err := GetKey()
	if err != nil {
		return nil, err
	}

	c := api.NewClient(key)
	c.Agent = fmt.Sprintf("urlscan-cli %s", version.Version)

	return &APIClient{c}, nil
}

type DownloadOptions struct {
	client          *APIClient
	path            string
	output          string
	directoryPrefix string
	force           bool
	silent          bool
}

type DownloadOption = func(*DownloadOptions)

func WithDownloadClient(client *APIClient) DownloadOption {
	return func(opts *DownloadOptions) {
		opts.client = client
	}
}

func WithDownloadURL(path string) DownloadOption {
	return func(opts *DownloadOptions) {
		opts.path = path
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
		opts.path = fmt.Sprintf("/dom/%s/", uuid)
	}
}

func WithDownloadScreenshot(uuid string) DownloadOption {
	return func(opts *DownloadOptions) {
		opts.path = fmt.Sprintf("/screenshots/%s.png", uuid)
	}
}

func WithDownloadSilent(silent bool) DownloadOption {
	return func(opts *DownloadOptions) {
		opts.silent = silent
	}
}

func WithDownloadDirectoryPrefix(directoryPrefix string) DownloadOption {
	return func(opts *DownloadOptions) {
		opts.directoryPrefix = directoryPrefix
	}
}

func NewDownloadOptions(opts ...DownloadOption) *DownloadOptions {
	var o DownloadOptions
	for _, fn := range opts {
		fn(&o)
	}
	return &o
}

func Download(opts *DownloadOptions) error {
	output := filepath.Join(opts.directoryPrefix, opts.output)

	if !opts.force {
		err := checkFileExists(output)
		if err != nil {
			return err
		}
	}

	_, err := opts.client.Download(opts.path, output)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("Downloaded: %s from %s\n", output, opts.path)
	if opts.silent {
		// output it to stderr to make the rest of stdout clean (for piping with jq, etc.)
		fmt.Fprint(os.Stderr, msg)
	} else {
		fmt.Print(msg)
	}

	return nil
}

func DownloadWithSpinner(opts *DownloadOptions) error {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	// start the spinner
	s.Start()

	err := Download(opts)

	// stop the spinner
	s.Stop()

	return err
}

type BatchJSONResultPair struct {
	Key    string          `json:"key"`
	Result json.RawMessage `json:"result"`
}

func NewBatchJSONResultPairs(keys []string, results []mo.Result[*api.Response]) []*BatchJSONResultPair {
	return lo.ZipBy2(keys, results, func(url string, result mo.Result[*api.Response]) *BatchJSONResultPair {
		return &BatchJSONResultPair{
			Key:    url,
			Result: *api.BatchResultToRaw(result),
		}
	})
}
