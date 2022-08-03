package proto

import (
	"github.com/boxgo/box/v2/codec"
	"google.golang.org/protobuf/proto"
)

type marshaler struct{}

func NewMarshaler() codec.Marshaler {
	return &marshaler{}
}

func (t marshaler) String() string {
	return "proto"
}

func (t marshaler) Marshal(v interface{}) ([]byte, error) {
	return proto.Marshal(v.(proto.Message))
}

func (t marshaler) Unmarshal(d []byte, v interface{}) error {
	return proto.Unmarshal(d, v.(proto.Message))
}
