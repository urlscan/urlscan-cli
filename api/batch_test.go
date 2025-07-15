package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/h2non/gock"
	"github.com/samber/mo"
	"github.com/stretchr/testify/assert"
)

func TestBatch(t *testing.T) {
	defer gock.Off()

	gock.New("http://testserver/").
		Get("/bar").
		Reply(http.StatusOK)
	gock.New("http://testserver/").
		Get("/baz").
		Reply(http.StatusOK)

	c := newTestClient()

	requests := make([]*http.Request, 2)
	req1, err := http.NewRequest("GET", URL("bar").String(), nil)
	assert.NoError(t, err)
	requests[0] = req1

	req2, err := http.NewRequest("GET", URL("baz").String(), nil)
	assert.NoError(t, err)
	requests[1] = req2

	tasks := make([]BatchTask[*http.Response], len(requests))
	for i, req := range requests {
		tasks[i] = func(cli *Client, ctx context.Context) mo.Result[*http.Response] {
			resp, err := cli.Do(req)
			if err != nil {
				return mo.Err[*http.Response](err)
			}
			defer resp.Body.Close() // nolint:errcheck
			return mo.Ok(resp)
		}
	}

	results, err := Batch(c, tasks, WithBatchMaxConcurrency(2), WithBatchTimeout(5))
	assert.NoError(t, err)

	assert.Equal(t, results[0].MustGet().StatusCode, http.StatusOK)
	assert.Equal(t, results[1].MustGet().StatusCode, http.StatusOK)
}
