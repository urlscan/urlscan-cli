package api

import (
	"net/url"
)

const (
	apiPrefix = "/api/v1/"
)

func PrefixedPath(path string) string {
	result, err := url.JoinPath(apiPrefix, path)
	if err != nil {
		panic(err)
	}
	return result
}
