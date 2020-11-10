package errcode

import (
	"errors"
	"fmt"
	"sync"
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
	registered      sync.Map
	regDup          = func(key string) error { return fmt.Errorf("errcode %s is registered", key) }
)

func Build(modNo, subNo uint) *ModuleErrorBuilder {
	if modNo < 100 || modNo > 999 {
		panic(ErrModNoInvalid)
	}
	if subNo > 999 {
		panic(ErrSubNoInvalid)
	}

	key := fmt.Sprintf("%d%d", modNo, subNo)
	if _, exist := registered.Load(key); exist {
		panic(regDup(key))
	}
	registered.Store(key, 0)

	return &ModuleErrorBuilder{
		modNo: modNo,
		subNo: subNo,
	}
}

func (builder *ModuleErrorBuilder) Build(bizNo uint, msg string) *ErrorCodeBuilder {
	if bizNo > 999 {
		panic(ErrBizNoInvalid)
	}

	key := fmt.Sprintf("%d%d%d", builder.modNo, builder.subNo, bizNo)
	if _, exist := registered.Load(key); exist {
		panic(regDup(key))
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
