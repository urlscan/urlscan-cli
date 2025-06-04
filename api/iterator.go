package api

import (
	"encoding/json"
	"fmt"
	"iter"
	"net/url"
	"strconv"
)

const MaxTotal = 10_000

type SearchResult struct {
	Sort []any           `json:"sort"`
	Raw  json.RawMessage `json:"-"`
}

type SearchResults struct {
	Results []SearchResult  `json:"results"`
	HasMore bool            `json:"has_more"`
	Total   int             `json:"total"`
	Raw     json.RawMessage `json:"-"`
}

func (r *SearchResult) UnmarshalJSON(data []byte) error {
	// a hack to prevent infinite UnmarshalJSON recursion loop
	// ref. https://biscuit.ninja/posts/go-avoid-an-infitine-loop-with-custom-json-unmarshallers/
	//      https://github.com/stripe/stripe-go/blob/f1847e5d06c4d13389d10629c6827dc375bfd015/event.go#L79-L91
	type result SearchResult
	var dst result

	err := json.Unmarshal(data, &dst)
	if err != nil {
		return err
	}
	*r = SearchResult(dst)
	r.Raw = data
	return err
}

func (r *SearchResults) UnmarshalJSON(data []byte) error {
	type results SearchResults
	var dst results

	err := json.Unmarshal(data, &dst)
	if err != nil {
		return err
	}
	*r = SearchResults(dst)
	r.Raw = data
	return err
}

type IteratorOption func(*Iterator) error

func IteratorSearchAfter(searchAfter string) IteratorOption {
	return func(it *Iterator) error {
		it.searchAfter = searchAfter
		return nil
	}
}

func IteratorSize(n int) IteratorOption {
	return func(it *Iterator) error {
		it.size = n
		return nil
	}
}

func IteratorLimit(n int) IteratorOption {
	return func(it *Iterator) error {
		it.limit = n
		return nil
	}
}

func IteratorQuery(q string) IteratorOption {
	return func(it *Iterator) error {
		it.q = q
		return nil
	}
}

type Iterator struct {
	client      *Client
	limit       int
	size        int
	q           string
	searchAfter string
	link        *url.URL
	count       int
	HasMore     bool
	Total       int
}

func newIterator(cli *Client, u *url.URL, options ...IteratorOption) (*Iterator, error) {
	it := &Iterator{
		client:  cli,
		HasMore: true,
		count:   0,
		Total:   0,
	}

	for _, opt := range options {
		if err := opt(it); err != nil {
			return nil, err
		}
	}

	query := u.Query()

	if it.q != "" {
		query.Add("q", it.q)
	}

	if it.searchAfter != "" {
		query.Add("search_after", it.searchAfter)
	}

	if it.size > 0 {
		query.Add("size", strconv.Itoa(it.size))
	}

	u.RawQuery = query.Encode()
	it.link = u

	return it, nil
}

func (it *Iterator) getMoreResults() (results []*SearchResult, err error) {
	resp, err := it.client.Get(it.link)
	if err != nil {
		return nil, err
	}

	r := &SearchResults{}
	err = json.Unmarshal(resp.Raw, r)
	if err != nil {
		return nil, err
	}

	// set total only once (= when the first request is made)
	if it.Total == 0 {
		it.Total = r.Total
	}

	for _, result := range r.Results {
		results = append(results, &result)
	}
	// set searchAfter for the next request
	if len(r.Results) > 0 {
		last := r.Results[len(r.Results)-1]

		if len(last.Sort) >= 2 {
			timestamp, ok := last.Sort[0].(float64)
			if !ok {
				return nil, fmt.Errorf("invalid result sort format")
			}

			uuid, ok := last.Sort[1].(string)
			if !ok {
				return nil, fmt.Errorf("invalid result sort format")
			}
			it.searchAfter = fmt.Sprintf("%s,%s", strconv.FormatFloat(timestamp, 'f', -1, 64), uuid)
		}
	}

	q := it.link.Query()
	q.Set("search_after", it.searchAfter)

	it.link.RawQuery = q.Encode()

	// set HasMore
	if r.Total != MaxTotal {
		it.HasMore = r.Total > (it.count + len(r.Results))
	} else {
		it.HasMore = len(r.Results) >= it.size
	}

	return results, nil
}

func (it *Iterator) Iterate() iter.Seq2[*SearchResult, error] {
	return func(yield func(*SearchResult, error) bool) {
		for it.limit == 0 || it.count < it.limit {
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
				if it.count >= it.limit {
					return
				}
			}

			if len(results) == 0 || !it.HasMore {
				return
			}
		}
	}
}
