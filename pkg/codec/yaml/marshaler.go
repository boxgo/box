package yaml

import (
	"github.com/boxgo/box/pkg/codec"
	"github.com/ghodss/yaml"
)

type marshaler struct{}

func NewMarshaler() codec.Marshaler {
	return &marshaler{}
}

func (y marshaler) String() string {
	return "yaml"
}

func (y marshaler) Marshal(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

func (y marshaler) Unmarshal(d []byte, v interface{}) error {
	return yaml.Unmarshal(d, v)
}
