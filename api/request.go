package api

import (
	"bytes"
	"context"
	"io"
	"net/http"
)

type GetContentFn func() (io.ReadCloser, error)

type Request struct {
	Headers     http.Header
	QueryParams map[string]string
	Path        string
	Method      string
	Body        []byte
	GetBody     GetContentFn
	client      *Client
	ctx         context.Context
	RawRequest  *http.Request
}

func (r *Request) SetPath(path string) *Request {
	r.Path = path
	return r
}

func (r *Request) SetMethod(method string) *Request {
	r.Method = method
	return r
}

func (r *Request) SetHeaders(headers map[string]string) *Request {
	for k, v := range headers {
		r.SetHeader(k, v)
	}
	return r
}

func (r *Request) SetHeader(key, value string) *Request {
	if r.Headers == nil {
		r.Headers = make(http.Header)
	}
	r.Headers.Set(key, value)
	return r
}

func (r *Request) SetQueryParams(params map[string]string) *Request {
	for k, v := range params {
		r.SetQueryParam(k, v)
	}
	return r
}

func (r *Request) SetQueryParam(key, value string) *Request {
	if r.QueryParams == nil {
		r.QueryParams = make(map[string]string)
	}
	r.QueryParams[key] = value
	return r
}

func (r *Request) SetContext(ctx context.Context) *Request {
	r.ctx = ctx
	return r
}

func (r *Request) Do() (resp *Response, err error) {
	defer func() {
		if resp == nil {
			resp = &Response{Request: r, err: nil, body: nil, Response: nil}
		}
	}()

	if r.Headers == nil {
		r.Headers = make(http.Header)
	}

	resp, _ = r.client.Do(r)
	err = resp.Error()
	return resp, err
}

func (r *Request) Send(method, path string) (*Response, error) {
	r.SetMethod(method)
	r.SetPath(path)
	return r.Do()
}

func (r *Request) SetBodyString(body string) *Request {
	return r.SetBodyBytes([]byte(body))
}

func (r *Request) SetBodyBytes(body []byte) *Request {
	r.Body = body
	r.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(body)), nil
	}
	return r
}

func (r *Request) SetBodyJSONBytes(body []byte) *Request {
	r.SetContentType("application/json")
	r.SetBodyBytes(body)
	return r
}

func (r *Request) SetContentType(contentType string) *Request {
	return r.SetHeader("Content-Type", contentType)
}

func (r *Request) Get(path string) (*Response, error) {
	return r.Send(http.MethodGet, path)
}

func (r *Request) Post(path string) (*Response, error) {
	return r.Send(http.MethodPost, path)
}

func (r *Request) Delete(path string) (*Response, error) {
	return r.Send(http.MethodDelete, path)
}

func (r *Request) Put(path string) (*Response, error) {
	return r.Send(http.MethodPut, path)
}
