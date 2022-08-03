package toml

import (
	"bytes"

	"github.com/BurntSushi/toml"
	"github.com/boxgo/box/v2/codec"
)

type marshaler struct{}

func NewMarshaler() codec.Marshaler {
	return &marshaler{}
}

func (t marshaler) String() string {
	return "toml"
}

func (t marshaler) Marshal(v interface{}) ([]byte, error) {
	b := bytes.NewBuffer(nil)
	err := toml.NewEncoder(b).Encode(v)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (t marshaler) Unmarshal(d []byte, v interface{}) error {
	return toml.Unmarshal(d, v)
}
