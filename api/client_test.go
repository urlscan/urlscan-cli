package api

import (
	"net/http"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

func newTestClient() *Client {
	c := NewClient("api_key")
	SetHost("http://testserver")
	return c
}

func TestSetHost(t *testing.T) {
	SetHost("https://testserver")
	url := URL("dummy")
	assert.Equal(t, "https://testserver/dummy", url.String())
}

func TestGet(t *testing.T) {
	defer gock.Off()

	gock.New("http://testserver/").
		Get("/bar").
		Reply(200).
		JSON(map[string]string{"foo": "bar"})

	c := newTestClient()
	got, err := c.Get(URL("bar"))
	assert.NoError(t, err)
	assert.Equal(t, "{\"foo\":\"bar\"}\n", string(got.Raw))
	assert.Equal(t, gock.IsDone(), true)
}

type Counter struct {
	count int
}

func (c *Counter) Count() int {
	original := c.count
	c.count++
	return original
}

func TestRetry(t *testing.T) {
	defer gock.Off()

	retryCounter := &Counter{}

	gock.New("http://testserver/").
		Get("/bar").
		AddMatcher(func(req *http.Request, ereq *gock.Request) (bool, error) { return retryCounter.Count() == 0, nil }).
		Reply(http.StatusTooManyRequests).
		SetHeaders(map[string]string{"X-Rate-Limit-Reset-After": "0"}).
		JSON(map[string]string{"foo": "bar"})

	gock.New("http://testserver/").
		Get("/bar").
		AddMatcher(func(req *http.Request, ereq *gock.Request) (bool, error) { return retryCounter.Count() == 1, nil }).
		Reply(http.StatusOK).
		JSON(map[string]string{"foo": "bar"})

	c := newTestClient()
	got, err := c.Get(URL("bar"))
	assert.NoError(t, err)
	assert.Equal(t, "{\"foo\":\"bar\"}\n", string(got.Raw))
	assert.Equal(t, gock.IsDone(), true)
}

func TestWaitAndGetResult(t *testing.T) {
	defer gock.Off()

	gock.New("http://testserver/").
		Post("/api/v1/scan/").
		Reply(http.StatusOK).
		JSON(map[string]string{"uuid": "dummy"})

	gock.New("http://testserver/").
		Get("/api/v1/result/dummy/").
		Reply(http.StatusOK).
		JSON(map[string]string{"foo": "bar"})

	c := newTestClient()
	// do scan to get UUID
	scanRes, err := c.Scan("http://localhost")
	assert.NoError(t, err)

	// wait for the result
	got, err := c.WaitAndGetResult(t.Context(), scanRes.UUID, 1)
	assert.NoError(t, err)
	assert.Equal(t, "{\"foo\":\"bar\"}\n", string(got.Raw))
	assert.Equal(t, gock.IsDone(), true)
}

func TestError(t *testing.T) {
	defer gock.Off()

	gock.New("http://testserver/").
		Get("/foo").
		Reply(http.StatusBadRequest).
		SetHeader("Content-Type", "application/json").
		BodyString(`{"status": 400, "message": "dummy"}`)

	c := newTestClient()
	req, err := http.NewRequest("GET", URL("http://testserver/foo").String(), nil)
	assert.NoError(t, err)

	_, err = c.DoWithJsonParse(req)
	assert.Error(t, err)
	assert.Equal(t, "dummy", err.Error())
}
