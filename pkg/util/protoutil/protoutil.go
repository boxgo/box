package protoutil

import (
	"fmt"
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
)

func MarshalAny(obj interface{}) (*any.Any, error) {
	typ := reflect.TypeOf(obj)
	val := fmt.Sprintf("%+v", obj)

	return &any.Any{TypeUrl: typ.Name(), Value: []byte(val)}, nil
}

func MarshalAnyProtoMessage(msg proto.Message) (*any.Any, error) {
	val, err := proto.Marshal(msg)

	return &any.Any{
		TypeUrl: string(proto.MessageReflect(msg).Descriptor().FullName()),
		Value:   val,
	}, err
}
