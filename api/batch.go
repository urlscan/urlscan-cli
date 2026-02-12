package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

func BatchResultToRaw(r mo.Result[*Response]) *json.RawMessage {
	err := r.Error()
	if err != nil {
		jsonErr, ok := errors.AsType[*JSONError](err)
		if ok {
			return &jsonErr.Raw
		}
		errRaw := json.RawMessage(fmt.Sprintf(`{"error": "%s"}`, err.Error()))
		return &errRaw
	}
	resp := r.MustGet()
	raw, err := resp.ToJSON()
	if err != nil {
		errRaw := json.RawMessage(fmt.Sprintf(`{"error": "%s"}`, err.Error()))
		return &errRaw
	}
	return raw
}
