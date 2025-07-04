package api

import (
	"encoding/json"
	"iter"
	"net/url"
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

func HostnameIteratorNoLimit(noLimit bool) HostnameIteratorOption {
	return func(it *HostnameIterator) error {
		it.noLimit = noLimit
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
	limit     int
	noLimit   bool
	size      int
	count     int
	PageState string
	link      *url.URL
	HasMore   bool
}

func newHostnameIterator(cli *Client, u *url.URL, options ...HostnameIteratorOption) (*HostnameIterator, error) {
	it := &HostnameIterator{
		client:  cli,
		HasMore: true,
		count:   0,
	}

	for _, opt := range options {
		if err := opt(it); err != nil {
			return nil, err
		}
	}

	query := u.Query()

	// size (number of results per batch) is "limit" in this API endpoint
	if it.size > 0 {
		query.Add("limit", strconv.Itoa(it.size))
	}

	if it.PageState != "" {
		query.Add("pageState", it.PageState)
	}

	u.RawQuery = query.Encode()
	it.link = u

	return it, nil
}

func (it *HostnameIterator) getMoreResults() (results []*json.RawMessage, err error) {
	resp, err := it.client.Get(it.link)
	if err != nil {
		return nil, err
	}

	r := &HostnameResults{}
	err = json.Unmarshal(resp.Raw, r)
	if err != nil {
		return nil, err
	}

	for _, result := range r.Results {
		results = append(results, &result)
	}

	// update pageState for the next request
	q := it.link.Query()
	q.Set("pageState", r.PageState)
	it.link.RawQuery = q.Encode()

	// update HasMore based on the number of results
	it.HasMore = len(r.Results) >= it.size

	return results, nil
}

func (it *HostnameIterator) Iterate() iter.Seq2[*json.RawMessage, error] {
	return func(yield func(*json.RawMessage, error) bool) {
		for it.count < it.limit || it.noLimit {
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
				if !it.noLimit && it.count >= it.limit {
					return
				}
			}

			if len(results) == 0 || !it.HasMore {
				return
			}
		}
	}
}
