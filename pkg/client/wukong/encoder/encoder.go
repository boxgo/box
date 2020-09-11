package encoder

import (
	"io"
)

type (
	Encoder interface {
		MimeType() string                             // mime type
		Encode(reader io.Reader, v interface{}) error // encode
		Decode(v interface{}) (io.Reader, error)
	}
)

const (
	MimeTypeJSON       = "application/json"
	MimeTypeXML        = "application/xml"
	MimeTypeForm       = "application/x-www-form-urlencoded"
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

func Decode(contentType string, v interface{}) (io.Reader, error) {
	var coder Encoder
	switch contentType {
	case MimeTypeJSON:
		coder = jsonEncoder
	case MimeTypeXML:
		coder = xmlEncoder
	default:
		coder = jsonEncoder
	}

	return coder.Decode(v)
}

func Encode(contentType string, reader io.Reader, v interface{}) error {
	var coder Encoder
	switch contentType {
	case MimeTypeJSON:
		coder = jsonEncoder
	case MimeTypeXML:
		coder = xmlEncoder
	default:
		coder = jsonEncoder
	}

	return coder.Encode(reader, v)
}
