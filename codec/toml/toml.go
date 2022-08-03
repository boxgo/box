package toml

import (
	"bytes"

	"github.com/BurntSushi/toml"
	"github.com/boxgo/box/v2/codec"
)

type coder struct{}

func NewCoder() codec.Coder {
	return &coder{}
}

func (t coder) String() string {
	return "toml"
}

func (t coder) Marshal(v interface{}) ([]byte, error) {
	b := bytes.NewBuffer(nil)

	if err := toml.NewEncoder(b).Encode(v); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (t coder) Unmarshal(d []byte, v interface{}) error {
	return toml.Unmarshal(d, v)
}
