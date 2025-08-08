package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

func newTestClient() *Client {
	c := NewClient("dummy")
	c.SetBaseURL(&url.URL{
		Scheme: "http",
		Host:   "testserver",
	})
	return c
}

func TestGet(t *testing.T) {
	defer gock.Off()

	gock.New("http://testserver/").
		Get("/bar").
		Reply(200).
		JSON(map[string]string{"foo": "bar"})

	c := newTestClient()
	resp, err := c.NewRequest().Get("/bar")
	assert.NoError(t, err)
	assert.Equal(t, "{\"foo\":\"bar\"}\n", string(resp.body))
	assert.Equal(t, gock.IsDone(), true)
}

func TestGetWithQueryParams(t *testing.T) {
	defer gock.Off()

	gock.New("http://testserver/").
		Get("/bar").
		MatchParam("foo", "bar").
		Reply(200).
		JSON(map[string]string{"foo": "bar"})

	c := newTestClient()
	resp, err := c.NewRequest().SetQueryParam("foo", "bar").Get("/bar")
	assert.NoError(t, err)
	assert.Equal(t, "{\"foo\":\"bar\"}\n", string(resp.body))
	assert.Equal(t, gock.IsDone(), true)
}

func TestPost(t *testing.T) {
	defer gock.Off()

	gock.New("http://testserver/").
		Post("/bar").
		MatchType("json").
		JSON(map[string]string{"foo": "bar"}).
		Reply(200).
		JSON(map[string]string{"bar": "baz"})

	c := newTestClient()
	resp, err := c.NewRequest().SetBodyJSONBytes([]byte(`{"foo":"bar"}`)).Post("/bar")
	assert.NoError(t, err)
	assert.Equal(t, "{\"bar\":\"baz\"}\n", string(resp.body))
	assert.Equal(t, gock.IsDone(), true)
}

func TestJSONError(t *testing.T) {
	defer gock.Off()

	jsonErr := JSONError{
		Status:  400,
		Message: "Bad Request",
	}
	marshalled, err := json.Marshal(jsonErr)
	assert.NoError(t, err)

	gock.New("http://testserver/").Get("/bar").Reply(400).SetHeader("Content-Type", "application/json").BodyString(string(marshalled))

	c := newTestClient()
	_, err = c.NewRequest().Get("/bar")
	assert.Error(t, err)
	assert.Equal(t, "Bad Request", err.Error())

	assert.Equal(t, gock.IsDone(), true)
}

func TestNetworkError(t *testing.T) {
	defer gock.Off()

	gock.New("http://testserver/").Get("/bar").ReplyError(fmt.Errorf("network error"))

	c := newTestClient()
	_, err := c.NewRequest().Get("/bar")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "network error")

	assert.Equal(t, gock.IsDone(), true)
}
