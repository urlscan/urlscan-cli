package api

import "encoding/json"

type JSONError struct {
	Status      int             `json:"status"`
	Message     string          `json:"message"`
	Description string          `json:"description,omitempty"`
	Raw         json.RawMessage `json:"-"`
}

func (r *JSONError) UnmarshalJSON(data []byte) error {
	type result JSONError
	var dst result

	err := json.Unmarshal(data, &dst)
	if err != nil {
		return err
	}
	*r = JSONError(dst)
	r.Raw = data
	return err
}

func (e JSONError) Error() string {
	return e.Message
}
