package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	version = "0.1.0"
)

var log = slog.New(slog.NewTextHandler(os.Stderr, nil))

var baseURL = url.URL{
	Scheme: "https",
	Host:   "urlscan.io",
}

type Request struct {
	Raw json.RawMessage `json:"-"`
}

type Response struct {
	Raw json.RawMessage `json:"-"`
}

type Error struct {
	Status      int             `json:"status"`
	Message     string          `json:"message"`
	Description string          `json:"description,omitempty"`
	Raw         json.RawMessage `json:"-"`
}

func (r *Error) UnmarshalJSON(data []byte) error {
	type result Error
	var dst result

	err := json.Unmarshal(data, &dst)
	if err != nil {
		return err
	}
	*r = Error(dst)
	r.Raw = data
	return err
}

func (e Error) Error() string {
	return e.Message
}

func URL(pathFmt string, a ...any) *url.URL {
	path := fmt.Sprintf(pathFmt, a...)
	url, err := url.Parse(path)
	if err != nil {
		msg := fmt.Sprintf("error formatting URL \"%s\": %s", pathFmt, err)
		panic(msg)
	}
	return baseURL.ResolveReference(url)
}

func (r *Response) PrettyJson() string {
	var jsonBody bytes.Buffer
	err := json.Indent(&jsonBody, r.Raw, "", "  ")
	if err != nil {
		log.Info("error formatting JSON response, fallback to the original", "error", err)
		return string(r.Raw)
	}
	return jsonBody.String()
}

func SetHost(host string) {
	if strings.HasPrefix(host, "https://") {
		baseURL.Scheme = "https"
		baseURL.Host = strings.TrimPrefix(host, "https://")
		return
	}

	if strings.HasPrefix(host, "http://") {
		baseURL.Scheme = "http"
		baseURL.Host = strings.TrimPrefix(host, "http://")
		return
	}

	baseURL.Host = host
}

type requestOptions struct {
	headers map[string]string
}

type RequestOption func(*requestOptions)

type RetryTransport struct {
	Transport http.RoundTripper
}

func (t *RetryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	res, err := t.Transport.RoundTrip(req)

	if res.StatusCode == http.StatusTooManyRequests {
		// rate limit headers: https://urlscan.io/docs/api/#ratelimit
		limitAction := res.Header.Get("X-Rate-Limit-Action")
		limitLimit := res.Header.Get("X-Rate-Limit-Limit")
		limitReset := res.Header.Get("X-Rate-Limit-Reset")
		limitScope := res.Header.Get("X-Rate-Limit-Scope")
		limitWindow := res.Header.Get("X-Rate-Limit-Window")
		retryAfter := res.Header.Get("X-Rate-Limit-Reset-After")

		if retryAfter != "" {
			retryAfterInt, err := strconv.Atoi(retryAfter)
			if err == nil {
				log.Info(fmt.Sprintf("Sleeping for %s seconds", retryAfter),
					"X-Rate-Limit-Action", limitAction,
					"X-Rate-Limit-Limit", limitLimit,
					"X-Rate-Limit-Reset-After", retryAfter,
					"X-Rate-Limit-Reset", limitReset,
					"X-Rate-Limit-Scope", limitScope,
					"X-Rate-Limit-Window", limitWindow,
				)
				time.Sleep(time.Duration(retryAfterInt) * time.Second)
			}
		}
		res, err = t.Transport.RoundTrip(req)
	}

	return res, err
}

type APIClient interface {
	Get(url *url.URL, options ...RequestOption) (any, error)
	Post(url *url.URL, req *Request, options ...RequestOption) (any, error)
}

type Client struct {
	APIKey     string
	Agent      string
	httpClient *http.Client
	headers    map[string]string
}

func WithHeader(header, value string) RequestOption {
	return func(opts *requestOptions) {
		if opts.headers == nil {
			opts.headers = make(map[string]string)
		}
		opts.headers[header] = value
	}
}

func opts(opts ...RequestOption) *requestOptions {
	o := &requestOptions{}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

type ClientOption func(*Client)

func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

func NewClient(APIKey string, opts ...ClientOption) *Client {
	c := &Client{APIKey: APIKey, httpClient: &http.Client{
		Transport: &RetryTransport{
			Transport: http.DefaultTransport,
		},
	}}
	for _, o := range opts {
		o(c)
	}
	return c
}

func (cli *Client) NewRequest(method string, url *url.URL, body io.Reader, headers map[string]string) (*http.Request, error) {
	req, err := http.NewRequest(method, url.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", fmt.Sprintf("urlscan-go/%s", version))
	req.Header.Set("API-Key", cli.APIKey)

	for k, v := range cli.headers {
		req.Header.Set(k, v)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return req, nil
}

func (cli *Client) Do(req *http.Request) (*http.Response, error) {
	resp, err := cli.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() //nolint:errcheck

	return resp, nil
}

func (cli *Client) parseResponse(resp *http.Response) (*Response, error) {
	jsonResp := &Response{}

	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json") {
		return nil, fmt.Errorf("expecting JSON response from %s %s",
			resp.Request.Method, resp.Request.URL.String())
	}

	var reader = resp.Body
	read, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	// consider 2xx response as successful
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		jsonResp.Raw = json.RawMessage(read)
		return jsonResp, nil
	}

	jsonErr := &Error{}
	err = json.Unmarshal(read, jsonErr)
	if err != nil {
		return nil, err
	}
	return nil, jsonErr
}

func (cli *Client) DoWithJsonParse(req *http.Request) (*Response, error) {
	resp, err := cli.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() //nolint:errcheck

	return cli.parseResponse(resp)
}

func (cli *Client) Get(url *url.URL, options ...RequestOption) (*Response, error) {
	o := opts(options...)
	req, err := cli.NewRequest("GET", url, nil, o.headers)
	if err != nil {
		return nil, err
	}
	return cli.DoWithJsonParse(req)
}

func (cli *Client) Post(url *url.URL, req *Request, options ...RequestOption) (*Response, error) {
	b := []byte(req.Raw)
	defaultContentTypeOptions := append(
		[]RequestOption{WithHeader("Content-Type", "application/json")},
		options...)
	o := opts(defaultContentTypeOptions...)

	httpReq, err := cli.NewRequest("POST", url, bytes.NewReader(b), o.headers)
	if err != nil {
		return nil, err
	}
	return cli.DoWithJsonParse(httpReq)
}

func (cli *Client) Delete(url *url.URL, options ...RequestOption) (*Response, error) {
	o := opts(options...)
	req, err := cli.NewRequest("DELETE", url, nil, o.headers)
	if err != nil {
		return nil, err
	}
	return cli.DoWithJsonParse(req)
}

func (cli *Client) Put(url *url.URL, req *Request, options ...RequestOption) (*Response, error) {
	b := []byte(req.Raw)
	defaultContentTypeOptions := append(
		[]RequestOption{WithHeader("Content-Type", "application/json")},
		options...)
	o := opts(defaultContentTypeOptions...)
	httpReq, err := cli.NewRequest("PUT", url, bytes.NewReader(b), o.headers)
	if err != nil {
		return nil, err
	}
	return cli.DoWithJsonParse(httpReq)
}

func (cli *Client) Download(url *url.URL, output string) (int64, error) {
	req, err := cli.NewRequest("GET", url, nil, nil)
	if err != nil {
		return 0, err
	}

	resp, err := cli.Do(req)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode == http.StatusOK {
		w, err := os.Create(output)
		if err != nil {
			return 0, err
		}
		defer func() {
			closeErr := w.Close()

			if closeErr != nil {
				err = closeErr
			}
		}()

		return io.Copy(w, resp.Body)
	}

	return 0, fmt.Errorf("unknown error downloading %q, HTTP response code: %d", url, resp.StatusCode)
}

func (cli *Client) Search(q string, options ...IteratorOption) (*Iterator, error) {
	u := URL("/api/v1/search/")
	query := u.Query()
	query.Add("q", q)
	u.RawQuery = query.Encode()

	return newIterator(cli, u, options...)
}

func (cli *Client) StructureSearch(uuid string, options ...IteratorOption) (*Iterator, error) {
	u := URL("/api/v1/pro/result/%s/similar/", uuid)
	return newIterator(cli, u, options...)
}

func (c *Client) IterateHostname(hostname string, opts ...HostnameIteratorOption) (*HostnameIterator, error) {
	u := URL("/api/v1/hostname/%s", hostname)
	return newHostnameIterator(c, u, opts...)
}

func (cli *Client) GetResult(uuid string) (*Response, error) {
	url := URL("%s", fmt.Sprintf("/api/v1/result/%s/", uuid))
	result, err := cli.Get(url)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (cli *Client) WaitAndGetResult(ctx context.Context, uuid string, maxWait int) (*Response, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(maxWait)*time.Second)
	defer cancel()

	log.Info("Waiting for scan to finish", "uuid", uuid)

	delay := 1 * time.Second

	for {
		result, err := cli.GetResult(uuid)
		if err == nil {
			return result, nil
		}

		// raise an error if it's not 404 error
		var apiErr *Error
		if errors.As(err, &apiErr) {
			if apiErr.Status != http.StatusNotFound {
				return nil, err
			}
		}

		select {
		case <-time.After(delay):
			delay += 1 * time.Second
			log.Info("Got 404 error, waiting for a scan result...", "delay", delay, "error", err.Error())
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func (cli *Client) CreateSubscription(opts ...SubscriptionOption) (*Response, error) {
	subscriptionOptions := newSubscriptionOptions(opts...)
	marshalled, err := json.Marshal(subscriptionOptions)
	if err != nil {
		return nil, err
	}

	return cli.Post(URL("/api/v1/user/subscriptions/"), &Request{
		Raw: json.RawMessage(marshalled),
	})
}

func (cli *Client) UpdateSubscription(opts ...SubscriptionOption) (*Response, error) {
	subscriptionOptions := newSubscriptionOptions(opts...)
	marshalled, err := json.Marshal(subscriptionOptions)
	if err != nil {
		return nil, err
	}

	url := URL("/api/v1/user/subscriptions/%s/", subscriptionOptions.Subscription.ID)
	return cli.Put(url, &Request{
		Raw: json.RawMessage(marshalled),
	})
}

func (cli *Client) CreateSavedSearch(opts ...SavedSearchOption) (*Response, error) {
	savedSearchOptions := newSavedSearchOptions(opts...)
	marshalled, err := json.Marshal(savedSearchOptions)
	if err != nil {
		return nil, err
	}

	return cli.Post(URL("/api/v1/user/searches/"), &Request{
		Raw: json.RawMessage(marshalled),
	})
}

func (cli *Client) UpdateSavedSearch(opts ...SavedSearchOption) (*Response, error) {
	savedSearchOptions := newSavedSearchOptions(opts...)
	marshalled, err := json.Marshal(savedSearchOptions)
	if err != nil {
		return nil, err
	}

	url := URL("/api/v1/user/searches/%s/", savedSearchOptions.Search.ID)
	return cli.Put(url, &Request{
		Raw: json.RawMessage(marshalled),
	})
}

func (cli *Client) CreateIncident(opts ...IncidentOption) (*Response, error) {
	incidentOpts := newIncidentOptions(opts...)
	marshalled, err := json.Marshal(incidentOpts)
	if err != nil {
		return nil, err
	}

	return cli.Post(URL("/api/v1/user/incidents/"), &Request{
		Raw: json.RawMessage(marshalled),
	})
}

func (cli *Client) UpdateIncident(id string, opts ...IncidentOption) (*Response, error) {
	incidentOpts := newIncidentOptions(opts...)
	marshalled, err := json.Marshal(incidentOpts)
	if err != nil {
		return nil, err
	}

	url := URL("/api/v1/user/incidents/%s/", id)
	return cli.Put(url, &Request{
		Raw: json.RawMessage(marshalled),
	})
}

func (cli *Client) TriggerNonBlockingLiveScan(id string, opts ...LiveScanOption) (*Response, error) {
	liveScanOpts := newLiveScanOptions(opts...)
	marshalled, err := json.Marshal(liveScanOpts)
	if err != nil {
		return nil, err
	}

	url := URL("/api/v1/livescan/%s/task/", id)
	return cli.Post(url, &Request{
		Raw: json.RawMessage(marshalled),
	})
}

func (cli *Client) TriggerLiveScan(id string, opts ...LiveScanOption) (*Response, error) {
	liveScanOpts := newLiveScanOptions(opts...)
	marshalled, err := json.Marshal(liveScanOpts)
	if err != nil {
		return nil, err
	}

	url := URL("/api/v1/livescan/%s/scan/", id)
	return cli.Post(url, &Request{
		Raw: json.RawMessage(marshalled),
	})
}

func (cli *Client) StoreLiveScanResult(scannerId string, scanId string, opts ...LiveScanStoreOption) (*Response, error) {
	liveScanStoreOpts := newLiveScanStoreOptions(opts...)
	marshalled, err := json.Marshal(liveScanStoreOpts)
	if err != nil {
		return nil, err
	}

	url := URL("/api/v1/livescan/%s/%s/", scannerId, scanId)
	return cli.Put(url, &Request{
		Raw: json.RawMessage(marshalled),
	})
}
