package encoder

import (
	"bytes"
	"encoding/json"
	"io"
)

type (
	JSONEncoder struct{}
)

func (c *JSONEncoder) MimeType() string {
	return MimeTypeJSON
}

func (c *JSONEncoder) Decode(v interface{}) (io.Reader, error) {
	data, err := json.Marshal(v)

	return bytes.NewReader(data), err
}
