package xml

import (
	"encoding/xml"

	"github.com/boxgo/box/pkg/codec"
)

type marshaler struct{}

func NewMarshaler() codec.Marshaler {
	return &marshaler{}
}

func (x marshaler) String() string {
	return "xml"
}

func (x marshaler) Marshal(v interface{}) ([]byte, error) {
	return xml.Marshal(v)
}

func (x marshaler) Unmarshal(d []byte, v interface{}) error {
	return xml.Unmarshal(d, v)
}
