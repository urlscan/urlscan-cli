package api

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/samber/mo"

	"golang.org/x/sync/errgroup"
)

type BatchOptions struct {
	MaxConcurrency int
	Timeout        int
}

type BatchOption func(*BatchOptions)

func WithBatchMaxConcurrency(max int) BatchOption {
	return func(opts *BatchOptions) {
		opts.MaxConcurrency = max
	}
}

func WithBatchTimeout(timeout int) BatchOption {
	return func(opts *BatchOptions) {
		opts.Timeout = timeout
	}
}

func newBatchOptions(opts ...BatchOption) *BatchOptions {
	var o BatchOptions
	for _, fn := range opts {
		fn(&o)
	}
	return &o
}

type BatchTask[T any] func(c *Client, ctx context.Context) mo.Result[T]

func Batch[T any](c *Client, tasks []BatchTask[T], opts ...BatchOption) ([]mo.Result[T], error) {
	var timeoutCtx context.Context
	var timeoutCancel context.CancelFunc
	var mu sync.Mutex

	batchOpts := newBatchOptions(opts...)
	if batchOpts.Timeout > 0 {
		timeoutCtx, timeoutCancel = context.WithTimeout(context.Background(), time.Duration(batchOpts.Timeout)*time.Second)
		defer timeoutCancel()
	} else {
		timeoutCtx = context.Background()
	}

	results := make([]mo.Result[T], len(tasks))

	g, ctx := errgroup.WithContext(timeoutCtx)
	g.SetLimit(batchOpts.MaxConcurrency)
	for i, task := range tasks {
		g.Go(func() error {
			result := task(c, ctx)

			mu.Lock()
			results[i] = result
			mu.Unlock()

			return nil
		})
	}

	err := g.Wait()
	if err != nil {
		return results, err
	}

	return results, nil
}

// errorRaw builds a JSON object {"error": "<msg>"} with the message properly
// escaped. Interpolating err.Error() directly produced invalid JSON whenever the
// error string contained quotes (e.g. *url.Error: `Get "http://host": ...`),
// breaking the marshal of the whole batch.
func errorRaw(msg string) *json.RawMessage {
	b, err := json.Marshal(map[string]string{"error": msg})
	if err != nil {
		b = []byte(`{"error": "failed to marshal error message"}`)
	}
	raw := json.RawMessage(b)
	return &raw
}

func BatchResultToRaw(r mo.Result[*Response]) *json.RawMessage {
	err := r.Error()
	if err != nil {
		jsonErr, ok := errors.AsType[*JSONError](err)
		if ok {
			return &jsonErr.Raw
		}
		return errorRaw(err.Error())
	}
	resp := r.MustGet()
	raw, err := resp.ToJSON()
	if err != nil {
		return errorRaw(err.Error())
	}
	return raw
}
