package api

import (
	"context"
	"github.com/samber/mo"
	"net/http"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

type BatchOptions struct {
	MaxConcurrency int
	TotalTimeout   int
}

type BatchOption func(*BatchOptions)

type BatchResult = mo.Result[*http.Response]
type JSONBatchResult = mo.Result[*Response]
type ScanBatchResult = mo.Result[*ScanResult]

func WithBatchMaxConcurrency(max int) BatchOption {
	return func(opts *BatchOptions) {
		opts.MaxConcurrency = max
	}
}

func WithBatchTotalTimeout(timeout int) BatchOption {
	return func(opts *BatchOptions) {
		opts.TotalTimeout = timeout
	}
}

func newBatchOptions(opts ...BatchOption) *BatchOptions {
	batchOpts := &BatchOptions{}

	for _, opt := range opts {
		opt(batchOpts)
	}

	return batchOpts
}

type BatchTask[T any] func(cli *Client, ctx context.Context) mo.Result[T]

func Batch[T any](cli *Client, tasks []BatchTask[T], opts ...BatchOption) ([]mo.Result[T], error) {
	var totalTimeoutCtx context.Context
	var totalTimeoutCancel context.CancelFunc
	var mu sync.Mutex

	batchOpts := newBatchOptions(opts...)
	if batchOpts.TotalTimeout > 0 {
		totalTimeoutCtx, totalTimeoutCancel = context.WithTimeout(context.Background(), time.Duration(batchOpts.TotalTimeout)*time.Second)
		defer totalTimeoutCancel()
	} else {
		totalTimeoutCtx = context.Background()
	}

	results := make([]mo.Result[T], len(tasks))

	g, ctx := errgroup.WithContext(totalTimeoutCtx)
	g.SetLimit(batchOpts.MaxConcurrency)
	for i, task := range tasks {
		i, task := i, task // capture loop variables
		g.Go(func() error {
			result := task(cli, ctx)

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
