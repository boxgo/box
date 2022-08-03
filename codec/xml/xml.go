package xml

import (
	"encoding/xml"

	"github.com/boxgo/box/v2/codec"
)

type coder struct{}

func NewCoder() codec.Coder {
	return &coder{}
}

func (x coder) String() string {
	return "xml"
}

func (x coder) Marshal(v interface{}) ([]byte, error) {
	return xml.Marshal(v)
}

func (x coder) Unmarshal(d []byte, v interface{}) error {
	return xml.Unmarshal(d, v)
}
