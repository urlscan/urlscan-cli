package utils

import (
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/briandowns/spinner"
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
	Client *APIClient
	URL    *url.URL
	Output string
	Force  bool
}

func NewDownloadOptions(client *APIClient, url *url.URL, output string, force bool) DownloadOptions {
	return DownloadOptions{
		Client: client,
		URL:    url,
		Output: output,
		Force:  force,
	}
}

func Download(opt DownloadOptions) error {
	// init and start the spinner
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Start()

	if !opt.Force {
		// check if the file already exists
		if _, err := os.Stat(opt.Output); err == nil {
			s.Stop()
			return fmt.Errorf("%s already exists", opt.Output)
		}
	}

	w, err := os.Create(opt.Output)
	if err != nil {
		return err
	}
	defer func() {
		closeErr := w.Close()
		if closeErr != nil {
			err = closeErr
		}
	}()

	_, err = opt.Client.Download(opt.URL, w)
	if err != nil {
		return err
	}

	// stop the spinner
	s.Stop()

	fmt.Printf("Downloaded: %s from %s\n", opt.Output, opt.URL.String())

	return nil
}
