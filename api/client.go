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
	Status      int    `json:"status"`
	Message     string `json:"message"`
	Description string `json:"description"`
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
		retryAfter := res.Header.Get("X-Rate-Limit-Reset-After")
		if retryAfter != "" {
			retryAfterInt, err := strconv.Atoi(retryAfter)
			if err == nil {
				log.Info("Rate limit exceeded", "X-Rate-Limit-Reset-After", retryAfter)
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

func (cli *Client) sendRequest(method string, url *url.URL, body io.Reader, headers map[string]string) (*http.Response, error) {
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

	return (cli.httpClient).Do(req)
}

func (cli *Client) parseResponse(resp *http.Response) (*Response, error) {
	apiResp := &Response{}

	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json") {
		return nil, fmt.Errorf("expecting JSON response from %s %s",
			resp.Request.Method, resp.Request.URL.String())
	}

	var reader = resp.Body
	if resp.StatusCode != http.StatusOK {
		apiErr := &Error{}
		if err := json.NewDecoder(reader).Decode(apiErr); err != nil {
			return nil, err
		}
		return nil, apiErr
	}

	read, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	apiResp.Raw = json.RawMessage(read)

	return apiResp, nil
}

func (cli *Client) Get(url *url.URL, options ...RequestOption) (*Response, error) {
	o := opts(options...)
	httpResp, err := cli.sendRequest("GET", url, nil, o.headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close() //nolint:errcheck

	return cli.parseResponse(httpResp)
}

func (cli *Client) Post(url *url.URL, req *Request, options ...RequestOption) (*Response, error) {
	b := []byte(req.Raw)
	defaultContentTypeOptions := append(
		[]RequestOption{WithHeader("Content-Type", "application/json")},
		options...)
	o := opts(defaultContentTypeOptions...)
	httpResp, err := cli.sendRequest("POST", url, bytes.NewReader(b), o.headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close() //nolint:errcheck

	return cli.parseResponse(httpResp)
}

func (cli *Client) Delete(url *url.URL, req *Request, options ...RequestOption) (*Response, error) {
	b := []byte(req.Raw)
	defaultContentTypeOptions := append(
		[]RequestOption{WithHeader("Content-Type", "application/json")},
		options...)
	o := opts(defaultContentTypeOptions...)
	httpResp, err := cli.sendRequest("DELETE", url, bytes.NewReader(b), o.headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close() //nolint:errcheck

	return cli.parseResponse(httpResp)
}

func (cli *Client) Put(url *url.URL, req *Request, options ...RequestOption) (*Response, error) {
	b := []byte(req.Raw)
	defaultContentTypeOptions := append(
		[]RequestOption{WithHeader("Content-Type", "application/json")},
		options...)
	o := opts(defaultContentTypeOptions...)
	httpResp, err := cli.sendRequest("PUT", url, bytes.NewReader(b), o.headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close() // nolint:errcheck

	return cli.parseResponse(httpResp)
}

func (cli *Client) Download(url *url.URL, w io.Writer) (int64, error) {
	resp, err := cli.sendRequest("GET", url, nil, nil)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close() // nolint:errcheck

	if resp.StatusCode == http.StatusOK {
		return io.Copy(w, resp.Body)
	}

	if _, err := cli.parseResponse(resp); err != nil {
		return 0, err
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

func (cli *Client) Scan(url string, options ...ScanOption) (*ScanResult, error) {
	scanOptions := newScanOptions(url, options...)

	marshalled, err := json.Marshal(scanOptions)
	if err != nil {
		return nil, err
	}

	resp, err := cli.Post(URL("/api/v1/scan/"), &Request{
		Raw: json.RawMessage(marshalled),
	})
	if err != nil {
		return nil, err
	}

	r := &ScanResult{}
	err = json.Unmarshal(resp.Raw, r)
	if err != nil {
		return nil, err
	}

	return r, nil
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
