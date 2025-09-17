package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Response struct {
	*http.Response
	err     error
	Request *Request
	body    []byte
}

func (r *Response) ToBytes() (body []byte, err error) {
	if r.err != nil {
		return nil, r.err
	}

	if r.body != nil {
		return r.body, nil
	}

	if r == nil || r.Body == nil {
		return []byte{}, nil
	}

	defer func() {
		closeErr := r.Body.Close()
		if closeErr != nil && err == nil {
			err = closeErr
		}

		if err != nil {
			r.err = err
		}
		r.body = body
	}()

	body, err = io.ReadAll(r.Body)

	return
}

func (r *Response) IsSuccess() bool {
	return r.StatusCode >= 200 && r.StatusCode <= 299
}

func (r *Response) IsError() bool {
	return !r.IsSuccess()
}

func (r *Response) Error() error {
	if r.err != nil {
		return r.err
	}

	if r.IsSuccess() {
		return nil
	}

	var jsonErr JSONError
	err := json.Unmarshal(r.body, &jsonErr)
	if err != nil {
		return err
	}
	return jsonErr
}

func (r *Response) Unmarshal(v any) error {
	if r.IsError() {
		return r.Error()
	}

	body, err := r.ToBytes()
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, v)
	if err != nil {
		return err
	}

	return nil
}

func (r *Response) ToString() (string, error) {
	b, err := r.ToBytes()
	return string(b), err
}

func (r *Response) GetContentType() string {
	if r.Response == nil {
		return ""
	}
	return r.Header.Get("Content-Type")
}

func (r *Response) ToJSON() (*json.RawMessage, error) {
	body, err := r.ToBytes()
	if err != nil {
		return nil, err
	}

	contentType := r.GetContentType()
	if !strings.Contains(contentType, "application/json") {
		return nil, fmt.Errorf("response is not JSON, content-type: %s", contentType)
	}
	raw := json.RawMessage(body)
	return &raw, nil
}

func (r *Response) PrettyJSON() string {
	var jsonBody bytes.Buffer
	err := json.Indent(&jsonBody, r.body, "", "  ")
	if err != nil {
		log.Info("error formatting JSON response, fallback to the original", "error", err)
		return string(r.body)
	}
	return jsonBody.String()
}
