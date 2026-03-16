package api

import (
	"encoding/json"
	"maps"
	"net/url"
)

const (
	apiPrefix = "/api/v1/"
)

// marshalJSONWithExtra marshals v (which should be a type alias to avoid
// infinite recursion) and merges extra into the top-level JSON object.
func marshalJSONWithExtra(v any, extra map[string]any) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	if len(extra) == 0 {
		return b, nil
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	maps.Copy(m, extra)
	return json.Marshal(m)
}

func PrefixedPath(path string) string {
	result, err := url.JoinPath(apiPrefix, path)
	if err != nil {
		panic(err)
	}
	return result
}
