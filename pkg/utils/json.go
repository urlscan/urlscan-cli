package utils

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
)

func MarshalIndent(v any) ([]byte, error) {
	return json.Marshal(v, jsontext.AllowDuplicateNames(true), jsontext.WithIndent("  "))
}
