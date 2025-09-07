package api

import (
	"strings"
)

const (
	apiPrefix = "/api/v1/"
)

func PrefixedPath(path string) string {
	// drop leading slashes to avoid double slashes in the URL
	trimmed := strings.TrimPrefix(path, "/")
	joined := strings.Join([]string{apiPrefix, trimmed}, "/")
	return strings.ReplaceAll(joined, "//", "/")
}
