package api

import (
	"encoding/json"
	"fmt"
)

type SavedSearchOptions struct {
	Search struct {
		ID               string   `json:"_id"`
		Datasource       string   `json:"datasource"`
		Name             string   `json:"name"`
		Query            string   `json:"query"`
		TLP              string   `json:"tlp"`
		Permissions      []string `json:"permissions"`
		Pass             int      `json:"pass"`
		Description      string   `json:"description,omitempty"`
		LongDescription  string   `json:"longDescription,omitempty"`
		OwnerDescription string   `json:"ownerDescription,omitempty"`
		Tags             []string `json:"tags,omitempty"`
		UserTags         []string `json:"userTags,omitempty"`
	} `json:"search"`
}

type SavedSearchOption func(*SavedSearchOptions)

func WithSavedSearchID(id string) SavedSearchOption {
	return func(opts *SavedSearchOptions) {
		opts.Search.ID = id
	}
}

func WithSavedSearchDatasource(datasource string) SavedSearchOption {
	return func(opts *SavedSearchOptions) {
		opts.Search.Datasource = datasource
	}
}

func WithSavedSearchName(name string) SavedSearchOption {
	return func(opts *SavedSearchOptions) {
		opts.Search.Name = name
	}
}

func WithSavedSearchQuery(query string) SavedSearchOption {
	return func(opts *SavedSearchOptions) {
		opts.Search.Query = query
	}
}

func WithSavedSearchTLP(tlp string) SavedSearchOption {
	return func(opts *SavedSearchOptions) {
		opts.Search.TLP = tlp
	}
}

func WithSavedSearchPermissions(permissions []string) SavedSearchOption {
	return func(opts *SavedSearchOptions) {
		opts.Search.Permissions = permissions
	}
}

func WithSavedSearchPass(pass int) SavedSearchOption {
	return func(opts *SavedSearchOptions) {
		opts.Search.Pass = pass
	}
}

func WithSavedSearchDescription(description string) SavedSearchOption {
	return func(opts *SavedSearchOptions) {
		opts.Search.Description = description
	}
}

func WithSavedSearchLongDescription(longDescription string) SavedSearchOption {
	return func(opts *SavedSearchOptions) {
		opts.Search.LongDescription = longDescription
	}
}

func WithSavedSearchOwnerDescription(ownerDescription string) SavedSearchOption {
	return func(opts *SavedSearchOptions) {
		opts.Search.OwnerDescription = ownerDescription
	}
}

func WithSavedSearchTags(tags []string) SavedSearchOption {
	return func(opts *SavedSearchOptions) {
		opts.Search.Tags = tags
	}
}

func WithSavedSearchUserTags(userTags []string) SavedSearchOption {
	return func(opts *SavedSearchOptions) {
		opts.Search.UserTags = userTags
	}
}

func newSavedSearchOptions(opts ...SavedSearchOption) *SavedSearchOptions {
	var o SavedSearchOptions
	for _, fn := range opts {
		fn(&o)
	}
	return &o
}

func (c *Client) CreateSavedSearch(opts ...SavedSearchOption) (*Response, error) {
	savedSearchOptions := newSavedSearchOptions(opts...)
	marshalled, err := json.Marshal(savedSearchOptions)
	if err != nil {
		return nil, err
	}

	return c.NewRequest().SetBodyJSONBytes(marshalled).Post(PrefixedPath("/user/searches/"))
}

func (c *Client) UpdateSavedSearch(opts ...SavedSearchOption) (*Response, error) {
	savedSearchOptions := newSavedSearchOptions(opts...)
	marshalled, err := json.Marshal(savedSearchOptions)
	if err != nil {
		return nil, err
	}
	return c.NewRequest().SetBodyJSONBytes(marshalled).Put(
		PrefixedPath(fmt.Sprintf("/user/searches/%s/", savedSearchOptions.Search.ID)),
	)
}
