package utils

import (
	"encoding/json/jsontext"
)

type SearchResults struct {
	Results []jsontext.Value `json:"results"`
	HasMore bool             `json:"has_more"`
	Total   int              `json:"total"`
}

func NewSearchResults() SearchResults {
	return SearchResults{
		Total:   0,
		Results: make([]jsontext.Value, 0),
		HasMore: false,
	}
}
