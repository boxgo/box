package jsonutil

import (
	"github.com/boxgo/box/pkg/codec/json"
)

var (
	j = json.NewMarshaler()
)

// Copy src to dest interface
func Copy(src, dest interface{}) error {
	data, err := j.Marshal(src)
	if err != nil {
		return err
	}

	return j.Unmarshal(data, dest)
}
