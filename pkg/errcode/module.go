package errcode

import (
	"errors"
	"fmt"

	"github.com/boxgo/box/pkg/logger"
)

type (
	ModuleErrorBuilder struct {
		modNo uint
		subNo uint
	}
)

var (
	ErrModNoInvalid = errors.New("mod no should in range 100-999")
	ErrSubNoInvalid = errors.New("sub no should in range 0-999")
	ErrBizNoInvalid = errors.New("biz no should in range 0-999")
)

func Build(modNo, subNo uint) *ModuleErrorBuilder {
	if modNo < 100 || modNo > 999 {
		logger.Panic(ErrModNoInvalid)
	}
	if subNo > 999 {
		logger.Panic(ErrSubNoInvalid)
	}

	key := fmt.Sprintf("%d%d", modNo, subNo)
	if _, exist := registered.Load(key); exist {
		logger.Panic(registerError(key))
	}
	registered.Store(key, 0)

	return &ModuleErrorBuilder{
		modNo: modNo,
		subNo: subNo,
	}
}

func (builder *ModuleErrorBuilder) Build(bizNo uint, msg string) *ErrorCodeBuilder {
	if bizNo > 999 {
		logger.Panic(ErrBizNoInvalid)
	}

	key := fmt.Sprintf("%d%d%d", builder.modNo, builder.subNo, bizNo)
	if _, exist := registered.Load(key); exist {
		logger.Panic(registerError(key))
	}

	ecb := &ErrorCodeBuilder{
		modNo:  builder.modNo,
		subNo:  builder.subNo,
		bizNo:  bizNo,
		msg:    msg,
		status: nil,
	}
	registered.Store(key, ecb)

	return ecb
}
