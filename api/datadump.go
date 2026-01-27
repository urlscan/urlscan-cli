package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
)

var (
	datePattern       = regexp.MustCompile(`^/?[a-z]+/[a-z]+/\d{8}/?$`)
	daysWindowPattern = regexp.MustCompile(`^/?days/[a-z]+/?`)
	fileTypePattern   = regexp.MustCompile(`^/?[a-z]+/[a-z]+/?`)
)

func hasDateInPath(path string) bool {
	return datePattern.MatchString(path)
}

func hasDaysWindowInPath(path string) bool {
	return daysWindowPattern.MatchString(path)
}

func hasFileTypeInPath(path string) bool {
	return fileTypePattern.MatchString(path)
}

func expandPath(path string) ([]string, error) {
	if !hasFileTypeInPath(path) {
		return nil, fmt.Errorf("path must include file type (api, search, screenshot, dom)")
	}

	// if date is provided or time window is "days", no expansion is needed
	if hasDateInPath(path) || hasDaysWindowInPath(path) {
		return []string{path}, nil
	}

	days := GetLast7Days(GetToday())
	paths := make([]string, 0, len(days))
	for _, day := range days {
		paths = append(paths, filepath.Join(path, day)+"/")
	}
	return paths, nil
}

type DataDumpList struct {
	Files []DataDumpFile  `json:"files"`
	Raw   json.RawMessage `json:"-"`
}

type DataDumpFile struct {
	Size      int64  `json:"size"`
	Path      string `json:"path"`
	Timestamp string `json:"timestamp"`
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

func (c *Client) BulkGetDataDumpList(path string) (*DataDumpList, error) {
	files := []DataDumpFile{}

	paths, err := expandPath(path)
	if err != nil {
		return nil, err
	}
	for _, p := range paths {
		list, err := c.GetDataDumpList(p)
		if err != nil {
			return nil, err
		}
		files = append(files, list.Files...)
	}

	// build aggregated raw JSON
	rawData, err := json.Marshal(map[string]any{"files": files})
	if err != nil {
		return nil, err
	}

	return &DataDumpList{Files: files, Raw: rawData}, nil
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
