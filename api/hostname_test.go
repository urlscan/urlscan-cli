package api

import (
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

func TestHostname(t *testing.T) {
	defer gock.Off()

	gock.New("http://testserver/").
		Get("/api/v1/hostname/example.com").
		Reply(200).
		SetHeader("Content-Type", "application/json").
		BodyString(`{"results":["dummy"], "pageState": "dummy"}`)

	// empty results & nulled pageState should stop the iteration
	gock.New("http://testserver/").
		Get("/api/v1/hostname/example.com").
		MatchParam("pageState", "dummy").
		Reply(200).
		SetHeader("Content-Type", "application/json").
		BodyString(`{"results":[], "pageState": null}`)

	c := newTestClient()
	it, err := c.IterateHostname("example.com", HostnameIteratorAll(true))
	assert.NoError(t, err)

	count := 0
	for result, err := range it.Iterate() {
		assert.NoError(t, err)
		assert.NotNil(t, result)
		count++
	}
	assert.Equal(t, 1, count)
}
