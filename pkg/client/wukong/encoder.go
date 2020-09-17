package wukong

import (
	"encoding/json"
	"encoding/xml"
)

type (
	Encoder interface {
		MimeType() string                          // mime type
		Encode(data []byte, val interface{}) error // encode data to val struct
		Decode(val interface{}) ([]byte, error)    // decode value to data
	}

	JSONEncoder struct{}
	XMLEncoder  struct{}
)

const (
	MimeTypeJSON       = "application/json"
	MimeTypeXML        = "application/xml"
	MimeTypeFormData   = "application/x-www-form-urlencoded"
	MimeTypeUrlencoded = "application/x-www-form-urlencoded"
	MimeTypeHTML       = "text/html"
	MimeTypeText       = "text/plain"
	MimeTypeMultipart  = "multipart/form-data"
)

var (
	jsonEncoder = &JSONEncoder{}
	xmlEncoder  = &XMLEncoder{}
)

func Decode(contentType string, val interface{}) ([]byte, error) {
	var coder Encoder
	switch contentType {
	case MimeTypeJSON:
		coder = jsonEncoder
	case MimeTypeXML:
		coder = xmlEncoder
	default:
		coder = jsonEncoder
	}

	return coder.Decode(val)
}

func Encode(contentType string, data []byte, val interface{}) error {
	var coder Encoder
	switch contentType {
	case MimeTypeJSON:
		coder = jsonEncoder
	case MimeTypeXML:
		coder = xmlEncoder
	default:
		coder = jsonEncoder
	}

	return coder.Encode(data, val)
}

func (c *JSONEncoder) MimeType() string {
	return MimeTypeJSON
}

func (c *JSONEncoder) Decode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (c *JSONEncoder) Encode(data []byte, val interface{}) error {
	return json.Unmarshal(data, val)
}

func (c *XMLEncoder) MimeType() string {
	return MimeTypeXML
}

func (c *XMLEncoder) Decode(v interface{}) ([]byte, error) {
	return xml.Marshal(v)
}

func (c *XMLEncoder) Encode(data []byte, val interface{}) error {
	return xml.Unmarshal(data, val)
}
