package jsonutil

import (
	"encoding/json"
)

// Copy src to dest interface
func Copy(src, dest interface{}) error {
	data, err := json.Marshal(src)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, dest)
}
