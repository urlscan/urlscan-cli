package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

func (c *Client) GetResult(uuid string) (*Response, error) {
	return c.NewRequest().Get(
		PrefixedPath(fmt.Sprintf("/result/%s/", uuid)),
	)
}

func (c *Client) WaitAndGetResult(ctx context.Context, uuid string, maxWait int) (*Response, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(maxWait)*time.Second)
	defer cancel()

	log.Info("Waiting for scan to finish", "uuid", uuid)

	delay := 1 * time.Second

	for {
		result, err := c.GetResult(uuid)
		if err == nil {
			return result, nil
		}

		// raise an error if it's not 404 error
		jsonErr, ok := errors.AsType[*JSONError](err)
		if ok {
			if jsonErr.Status != http.StatusNotFound {
				return nil, err
			}
		}

		select {
		case <-time.After(delay):
			delay += 1 * time.Second
			log.Info("Got 404 error, waiting for a scan result...", "delay", delay, "error", err.Error(), "uuid", uuid)
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

type ResultVisibilityOptions struct {
	Visibility string `json:"visibility"`
}

type ResultVisibilityOption func(*ResultVisibilityOptions)

func WithResultVisibility(visibility string) ResultVisibilityOption {
	return func(opts *ResultVisibilityOptions) {
		if visibility != "" {
			opts.Visibility = visibility
		}
	}
}

func (c *Client) UpdateResultVisibility(uuid string, opts ...ResultVisibilityOption) (*Response, error) {
	var options ResultVisibilityOptions
	for _, opt := range opts {
		opt(&options)
	}

	marshalled, err := json.Marshal(options)
	if err != nil {
		return nil, err
	}

	return c.NewRequest().SetBodyJSONBytes(marshalled).Put(
		PrefixedPath(fmt.Sprintf("/result/%s/visibility/", uuid)),
	)
}
