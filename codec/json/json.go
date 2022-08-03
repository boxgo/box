package json

import (
	"encoding/json"

	"github.com/boxgo/box/v2/codec"
	jsonIter "github.com/json-iterator/go"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type coder struct{}

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

	jsonIterMarshler = jsonIter.ConfigCompatibleWithStandardLibrary
)

func NewCoder() codec.Coder {
	return &coder{}
}

func (j coder) String() string {
	return "json"
}

func (coder) Marshal(v interface{}) ([]byte, error) {
	switch m := v.(type) {
	case json.Marshaler:
		return m.MarshalJSON()
	case proto.Message:
		return protoMarshler.Marshal(m)
	default:
		return jsonIterMarshler.Marshal(m)
	}
}

func (coder) Unmarshal(data []byte, v interface{}) error {
	switch m := v.(type) {
	case json.Unmarshaler:
		return m.UnmarshalJSON(data)
	case proto.Message:
		return protoUnmarshler.Unmarshal(data, m)
	default:
		return jsonIterMarshler.Unmarshal(data, m)
	}
}
