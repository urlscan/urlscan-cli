package api

import (
	"fmt"
	"io"
	"os"
)

func (c *Client) Download(path, output string) (n int64, err error) {
	resp, err := c.NewRequest().Get(path)
	if err != nil {
		return 0, err
	}

	if resp.IsSuccess() {
		w, err := os.Create(output)
		if err != nil {
			return 0, err
		}
		defer func() {
			closeErr := w.Close()
			if closeErr != nil && err == nil {
				err = closeErr
			}
		}()

		return io.Copy(w, resp.Body)
	}

	return 0, fmt.Errorf("unknown error downloading %q, HTTP response code: %d", resp.Request.RawRequest.URL, resp.StatusCode)
}
