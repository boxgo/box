package encoder

import (
	"io"
)

type (
	Encoder interface {
		MimeType() string // mime type
		// Encode(v interface{}) (io.Writer, error) // encode
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
)

func Decode(contentType string, v interface{}) (io.Reader, error) {
	var coder Encoder
	switch contentType {
	case MimeTypeJSON:
		coder = jsonEncoder
	default:
		coder = jsonEncoder
	}

	return coder.Decode(v)
}
