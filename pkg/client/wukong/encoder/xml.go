package encoder

import (
	"bytes"
	"encoding/xml"
	"io"
)

type (
	XMLEncoder struct{}
)

func (c *XMLEncoder) MimeType() string {
	return MimeTypeXML
}

func (c *XMLEncoder) Decode(v interface{}) (io.Reader, error) {
	data, err := xml.Marshal(v)

	return bytes.NewReader(data), err
}

func (c *XMLEncoder) Encode(reader io.Reader, v interface{}) error {
	dec := xml.NewDecoder(reader)

	return dec.Decode(v)
}
