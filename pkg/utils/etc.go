package utils

import "fmt"

func GetURLFromMap(json map[string]any) (string, error) {
	url, ok := json["url"].(string)
	if !ok {
		return "", fmt.Errorf("url field is missing or not a string in JSON: %v", json)
	}

	return url, nil
}
