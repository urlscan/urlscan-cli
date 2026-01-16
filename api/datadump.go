package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
)

type DataDumpList struct {
	Files []DataDumpFile  `json:"files"`
	Raw   json.RawMessage `json:"-"`
}

type DataDumpFile struct {
	Size int64  `json:"size"`
	Path string `json:"path"`
}

func (r *DataDumpList) PrettyJSON() string {
	var jsonBody bytes.Buffer
	err := json.Indent(&jsonBody, r.Raw, "", "  ")
	if err != nil {
		msg := fmt.Sprintf("error formatting JSON response: %s", err)
		panic(msg)
	}
	return jsonBody.String()
}

func (r *DataDumpList) UnmarshalJSON(data []byte) error {
	type result DataDumpList
	var dst result

	err := json.Unmarshal(data, &dst)
	if err != nil {
		return err
	}
	*r = DataDumpList(dst)
	r.Raw = data
	return err
}

func (c *Client) GetDataDumpList(path string) (*DataDumpList, error) {
	path, err := url.JoinPath("/datadump/list/", path)
	if err != nil {
		return nil, err
	}

	req := c.NewRequest().SetPath(PrefixedPath(path)).SetMethod("GET")
	resp, err := req.Do()
	if err != nil {
		return nil, err
	}

	var r DataDumpList
	err = resp.Unmarshal(&r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}
