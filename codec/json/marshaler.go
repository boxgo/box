package json

import (
	"encoding/json"

	"github.com/boxgo/box/v2/codec"
	jsoniter "github.com/json-iterator/go"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type marshaler struct{}

var (
	protoMarshler = protojson.MarshalOptions{
		Multiline:       false,
		Indent:          "",
		AllowPartial:    false,
		UseProtoNames:   false,
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
		Resolver:        nil,
	}

	protoUnmarshler = protojson.UnmarshalOptions{
		AllowPartial:   false,
		DiscardUnknown: true,
		Resolver:       nil,
	}

	jsoniterMarshler = jsoniter.ConfigCompatibleWithStandardLibrary
)

func NewMarshaler() codec.Marshaler {
	return &marshaler{}
}

func (j marshaler) String() string {
	return "json"
}

func (marshaler) Marshal(v interface{}) ([]byte, error) {
	switch m := v.(type) {
	case json.Marshaler:
		return m.MarshalJSON()
	case proto.Message:
		return protoMarshler.Marshal(m)
	default:
		return jsoniterMarshler.Marshal(m)
	}
}

func (marshaler) Unmarshal(data []byte, v interface{}) error {
	switch m := v.(type) {
	case json.Unmarshaler:
		return m.UnmarshalJSON(data)
	case proto.Message:
		return protoUnmarshler.Unmarshal(data, m)
	default:
		return jsoniterMarshler.Unmarshal(data, m)
	}
}
