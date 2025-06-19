package utils

import (
	"encoding/json"
)

type SearchResults struct {
	Results []json.RawMessage `json:"results"`
	HasMore bool              `json:"has_more"`
	Total   int               `json:"total"`
}

func NewSearchResults() SearchResults {
	return SearchResults{
		Total:   0,
		Results: make([]json.RawMessage, 0),
		HasMore: false,
	}
}
