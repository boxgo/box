package protoutil

import (
	"fmt"
	"reflect"

	rpcStatus "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func MarshalAny(obj interface{}) (*anypb.Any, error) {
	typ := reflect.TypeOf(obj)
	val := fmt.Sprintf("%+v", obj)

	return &anypb.Any{TypeUrl: typ.Name(), Value: []byte(val)}, nil
}

func MarshalAnyProtoMessage(msg proto.Message) (*anypb.Any, error) {
	val, err := proto.Marshal(msg)

	return &anypb.Any{
		TypeUrl: string(msg.ProtoReflect().Descriptor().FullName()),
		Value:   val,
	}, err
}

func ConvertToStatus(err error) *status.Status {
	return status.Convert(err)
}

func ConvertToStatusError(err error) error {
	return ConvertToStatus(err).Err()
}

func ErrorProto(s *rpcStatus.Status) error {
	return status.ErrorProto(s)
}
