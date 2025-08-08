package api

import (
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

func TestSearch(t *testing.T) {
	defer gock.Off()

	gock.New("http://testserver/").
		Get("/api/v1/search").
		MatchParam("q", "test").
		Reply(200).
		SetHeader("Content-Type", "application/json").
		BodyString(`{"results":[{"sort":[1,"dummy"]}], "total": 2, "has_more": true}`)

	gock.New("http://testserver/").
		Get("/api/v1/search").
		MatchParam("q", "test").
		MatchParam("search_after", "1,dummy").
		Reply(200).
		SetHeader("Content-Type", "application/json").
		BodyString(`{"results":[{"sort":[2,"dummy"]}], "total": 2, "has_more": false}`)

	gock.New("http://testserver/").
		Get("/api/v1/search").
		MatchParam("q", "test").
		MatchParam("search_after", "2,dummy").
		Reply(200).
		SetHeader("Content-Type", "application/json").
		BodyString(`{"results":[], "total": 2, "has_more": false}`)

	c := newTestClient()
	it, err := c.Search("test", IteratorAll(true))
	assert.NoError(t, err)

	count := 0
	for result, err := range it.Iterate() {
		assert.NoError(t, err)
		assert.NotNil(t, result)
		count++
	}
	assert.Equal(t, 2, count)
}
