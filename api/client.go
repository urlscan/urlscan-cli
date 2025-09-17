package api

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strconv"
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

type RetryTransport struct {
	Transport http.RoundTripper
}

func (t *RetryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	res, err := t.Transport.RoundTrip(req)
	if err == nil && res.StatusCode == http.StatusTooManyRequests {
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

type Client struct {
	APIKey     string
	Agent      string
	Err        error
	BaseURL    *url.URL
	httpClient *http.Client
}

func SetHost(host string) {
	baseURL.Host = host
}

func (c *Client) SetBaseURL(url *url.URL) *Client {
	c.BaseURL = url
	return c
}

func (c *Client) SetTransport(transport http.RoundTripper) *Client {
	if c.httpClient == nil || c.httpClient.Transport == nil {
		c.httpClient = &http.Client{}
	}
	c.httpClient.Transport = transport
	return c
}

func (c *Client) SetRetryTransport() *Client {
	c.SetTransport(&RetryTransport{
		Transport: http.DefaultTransport,
	})
	return c
}

func (c *Client) SetAPIKey(key string) *Client {
	c.APIKey = key
	return c
}

func (c *Client) SetAgent(agent string) *Client {
	c.Agent = agent
	return c
}

func NewClient(APIKey string) *Client {
	c := &Client{httpClient: &http.Client{}, BaseURL: &baseURL, APIKey: "", Agent: "", Err: nil}
	c.SetAPIKey(APIKey)
	c.SetAgent(fmt.Sprintf("urlscan-go/%s", version))
	c.SetRetryTransport()
	return c
}

func (c *Client) NewRequest() *Request {
	headers := make(http.Header)
	if c.APIKey != "" {
		headers.Set("API-Key", c.APIKey)
	}
	if c.Agent != "" {
		headers.Set("User-Agent", c.Agent)
	}
	return &Request{
		client:  c,
		Headers: headers,
		// default values
		Body:        nil,
		ctx:         nil,
		GetBody:     nil,
		Method:      "",
		Path:        "",
		QueryParams: make(map[string]string),
		RawRequest:  nil,
	}
}

func (c *Client) URL(path string) (*url.URL, error) {
	url, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	return c.BaseURL.ResolveReference(url), nil
}

func (c *Client) Do(r *Request) (resp *Response, err error) {
	resp = &Response{Request: r, err: nil, body: nil, Response: nil}
	defer func() {
		if err != nil {
			resp.err = err
		} else {
			err = resp.err
		}
	}()

	url, err := c.URL(r.Path)
	if err != nil {
		return nil, fmt.Errorf("error formatting URL: %w", err)
	}

	var reqBody io.ReadCloser
	if r.GetBody != nil {
		reqBody, resp.err = r.GetBody()
		if resp.err != nil {
			return
		}
	}

	req, err := http.NewRequest(r.Method, url.String(), reqBody)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}

	// set headers
	headers := r.Headers
	if headers == nil {
		headers = make(http.Header)
	}
	req.Header = headers
	req.Header.Set("User-Agent", c.Agent)
	req.Header.Set("API-Key", c.APIKey)

	// set query parameters
	if r.QueryParams != nil {
		query := req.URL.Query()
		for k, v := range r.QueryParams {
			query.Set(k, v)
		}
		req.URL.RawQuery = query.Encode()
	}

	ctx := r.ctx
	if ctx != nil {
		req = req.WithContext(ctx)
	}
	r.RawRequest = req

	resp.Response, resp.err = c.httpClient.Do(req)
	if resp.err == nil && resp.StatusCode >= 200 {
		// set resp.body
		_, err = resp.ToBytes()
		if err != nil {
			return nil, err
		}
		// restore body for re-reading
		resp.Body = io.NopCloser(bytes.NewReader(resp.body))
	}

	return
}
