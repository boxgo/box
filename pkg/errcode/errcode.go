package errcode

import (
	"github.com/boxgo/box/pkg/util/protoutil"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"google.golang.org/genproto/googleapis/rpc/status"
	gst "google.golang.org/grpc/status"
)

type (
	ErrorCodeBuilder struct {
		modNo  uint
		subNo  uint
		bizNo  uint
		msg    string
		status *status.Status
	}
)

func Parse(err error) *gst.Status {
	return gst.Convert(err)
}

func (builder *ErrorCodeBuilder) Build(inputs ...interface{}) error {
	if details, err := builder.convert2Details(inputs...); err != nil {
		return gst.Convert(err).Err()
	} else {
		s := &status.Status{
			Code:    builder.Code(),
			Message: builder.msg,
			Details: details,
		}

		return gst.ErrorProto(s)
	}
}

func (builder *ErrorCodeBuilder) Code() int32 {
	return int32(builder.modNo*1000000 + builder.subNo*1000 + builder.bizNo*1)
}

func (builder *ErrorCodeBuilder) Message() string {
	return builder.msg
}

func (builder *ErrorCodeBuilder) convert2Details(inputs ...interface{}) ([]*any.Any, error) {
	details := make([]*any.Any, len(inputs))

	for idx, input := range inputs {
		var (
			err    error
			detail *any.Any
		)

		switch val := input.(type) {
		case proto.Message:
			detail, err = protoutil.MarshalAnyProtoMessage(val)
		case *proto.Message:
			detail, err = protoutil.MarshalAnyProtoMessage(*val)
		default:
			detail, err = protoutil.MarshalAny(val)
		}

		if err != nil {
			return nil, err
		}

		details[idx] = detail
	}

	return details, nil
}
