package api

import (
	"encoding/json"
	"fmt"
	"iter"
	"strconv"
)

type HostnameResults struct {
	Item      string            `json:"item"`
	Results   []json.RawMessage `json:"results"`
	PageState string            `json:"pageState"`
	Raw       json.RawMessage   `json:"-"`
}

func (r *HostnameResults) UnmarshalJSON(data []byte) error {
	type results HostnameResults
	var dst results

	err := json.Unmarshal(data, &dst)
	if err != nil {
		return err
	}
	*r = HostnameResults(dst)
	r.Raw = data
	return err
}

type HostnameIteratorOption func(*HostnameIterator) error

// size is the number of results returned by the iterator in each batch.
func HostnameIteratorSize(size int) HostnameIteratorOption {
	return func(it *HostnameIterator) error {
		it.size = size
		return nil
	}
}

// limit is the maximum number of results that will be returned by the iterator.
// note that this is not the same as the API endpoint's "limit" query parameter.
func HostnameIteratorLimit(limit int) HostnameIteratorOption {
	return func(it *HostnameIterator) error {
		it.limit = limit
		return nil
	}
}

func HostnameIteratorAll(all bool) HostnameIteratorOption {
	return func(it *HostnameIterator) error {
		it.all = all
		return nil
	}
}

func HostnameIteratorPageState(pageState string) HostnameIteratorOption {
	return func(it *HostnameIterator) error {
		it.PageState = pageState
		return nil
	}
}

type HostnameIterator struct {
	client    *Client
	path      string
	request   *Request
	limit     int
	all       bool
	size      int
	count     int
	PageState string
	HasMore   bool
}

func newHostnameIterator(c *Client, path string, options ...HostnameIteratorOption) (*HostnameIterator, error) {
	request := c.NewRequest().SetPath(path)

	it := &HostnameIterator{
		client:  c,
		path:    path,
		request: request,
		// default values
		all:       false,
		count:     0,
		HasMore:   true,
		limit:     0,
		PageState: "",
		size:      0,
	}

	for _, opt := range options {
		if err := opt(it); err != nil {
			return nil, err
		}
	}

	// size (number of results per batch) is "limit" in this API endpoint
	if it.size > 0 {
		it.request.SetQueryParam("limit", strconv.Itoa(it.size))
	}

	if it.PageState != "" {
		it.request.SetQueryParam("pageState", it.PageState)
	}

	return it, nil
}

func (it *HostnameIterator) getMoreResults() (results []*json.RawMessage, err error) {
	resp, err := it.request.Get(it.path)
	if err != nil {
		return nil, err
	}

	var r HostnameResults
	err = resp.Unmarshal(&r)
	if err != nil {
		return nil, err
	}

	for _, result := range r.Results {
		results = append(results, &result)
	}

	// update pageState for the next request
	it.request.SetQueryParam("pageState", r.PageState)

	// update HasMore based on the number of results
	it.HasMore = len(r.Results) >= it.size

	return results, nil
}

func (it *HostnameIterator) Iterate() iter.Seq2[*json.RawMessage, error] {
	return func(yield func(*json.RawMessage, error) bool) {
		for it.count < it.limit || it.all {
			results, err := it.getMoreResults()
			if err != nil {
				yield(nil, err)
				return
			}

			for _, result := range results {
				if !yield(result, nil) {
					return
				}

				it.count++
				if !it.all && it.count >= it.limit {
					return
				}
			}

			if len(results) == 0 || !it.HasMore {
				return
			}
		}
	}
}

func (c *Client) IterateHostname(hostname string, opts ...HostnameIteratorOption) (*HostnameIterator, error) {
	return newHostnameIterator(c, PrefixedPath(fmt.Sprintf("/hostname/%s", hostname)), opts...)
}
