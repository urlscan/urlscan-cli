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

	paths := []string{"/bar", "/baz"}
	tasks := make([]BatchTask[*Response], len(paths))
	for i, path := range paths {
		tasks[i] = func(c *Client, ctx context.Context) mo.Result[*Response] {
			resp, err := c.NewRequest().Get(path)
			if err != nil {
				return mo.Err[*Response](err)
			}
			return mo.Ok(resp)
		}
	}

	results, err := Batch(c, tasks, WithBatchMaxConcurrency(2), WithBatchTimeout(5))
	assert.NoError(t, err)

	assert.Equal(t, results[0].MustGet().StatusCode, http.StatusOK)
	assert.Equal(t, results[1].MustGet().StatusCode, http.StatusOK)
}
