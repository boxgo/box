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

func (c *JSONEncoder) Encode(reader io.Reader, v interface{}) error {
	dec := json.NewDecoder(reader)

	return dec.Decode(v)
}
