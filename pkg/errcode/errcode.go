package errcode

import (
	"github.com/boxgo/box/pkg/util/protoutil"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
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

func ParseStatus(err error) *status.Status {
	return protoutil.ConvertToStatus(err).Proto()
}

func (builder *ErrorCodeBuilder) Build(inputs ...interface{}) error {
	if details, err := builder.convert2Details(inputs...); err != nil {
		return protoutil.ConvertToStatusError(err)
	} else {
		s := &status.Status{
			Code:    builder.Code(),
			Message: builder.Message(),
			Details: details,
		}

		return protoutil.ErrorProto(s)
	}
}

func (builder *ErrorCodeBuilder) Code() int32 {
	return int32(builder.modNo*1000000 + builder.subNo*1000 + builder.bizNo*1)
}

func (builder *ErrorCodeBuilder) Message() string {
	return builder.msg
}

func (builder *ErrorCodeBuilder) convert2Details(inputs ...interface{}) ([]*anypb.Any, error) {
	details := make([]*anypb.Any, len(inputs))

	for idx, input := range inputs {
		var (
			err    error
			detail *anypb.Any
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
