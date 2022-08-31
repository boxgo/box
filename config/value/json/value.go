package json

import (
	"fmt"
	"strconv"
	"time"
	"unsafe"

	"github.com/araddon/dateparse"
	jsoniter "github.com/json-iterator/go"
)

type (
	jsonValue struct {
		value jsoniter.Any
		api   jsoniter.API
	}
)

func init() {
	jsoniter.RegisterTypeDecoderFunc("time.Duration", func(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
		switch iter.WhatIsNext() {
		case jsoniter.NumberValue:
			*(*time.Duration)(ptr) = time.Duration(iter.ReadInt64()) * time.Second
		case jsoniter.StringValue:
			str := iter.ReadString()
			duration, err := time.ParseDuration(str)
			if err != nil {
				iter.ReportError("time.ParseDuration", err.Error())
			} else {
				*(*time.Duration)(ptr) = duration
			}
		default:
			*(*interface{})(ptr) = iter.Read()
		}
	})
	jsoniter.RegisterTypeDecoderFunc("time.Time", func(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
		switch iter.WhatIsNext() {
		case jsoniter.NumberValue:
			*(*time.Time)(ptr) = time.Unix(iter.ReadInt64(), 0)
		case jsoniter.StringValue:
			str := iter.ReadString()
			t, err := dateparse.ParseLocal(str)
			if err != nil {
				iter.ReportError("dateparse.ParseLocal", err.Error())
			} else {
				*(*time.Time)(ptr) = t
			}
		default:
			*(*interface{})(ptr) = iter.Read()
		}
	})
}

func (j *jsonValue) Bool() bool {
	value := j.value.ToBool()
	err := j.value.LastError()
	if err == nil {
		return value
	}

	if str, ok := j.value.GetInterface().(string); !ok {
		return false
	} else if v, e := strconv.ParseBool(str); e == nil {
		return v
	}

	return false
}

func (j *jsonValue) Uint() uint {
	value := j.value.ToUint()
	err := j.value.LastError()
	if err == nil {
		return value
	}

	if str, ok := j.value.GetInterface().(string); !ok {
		return 0
	} else if v, e := strconv.ParseUint(str, 10, 64); e == nil {
		return uint(v)
	}

	return 0
}

func (j *jsonValue) Int() int {
	value := j.value.ToInt()
	err := j.value.LastError()
	if err == nil {
		return value
	}

	if str, ok := j.value.GetInterface().(string); !ok {
		return 0
	} else if value, err = strconv.Atoi(str); err == nil {
		return value
	}

	return 0
}

func (j *jsonValue) Int32() int32 {
	value := j.value.ToInt32()
	err := j.value.LastError()
	if err == nil {
		return value
	}

	if str, ok := j.value.GetInterface().(string); !ok {
		return 0
	} else if v, e := strconv.ParseInt(str, 10, 32); e == nil {
		return int32(v)
	}

	return 0
}

func (j *jsonValue) Int64() int64 {
	value := j.value.ToInt64()
	err := j.value.LastError()
	if err == nil {
		return value
	}

	if str, ok := j.value.GetInterface().(string); !ok {
		return 0
	} else if value, err = strconv.ParseInt(str, 10, 64); err == nil {
		return value
	}

	return 0
}

func (j *jsonValue) Float32() float32 {
	value := j.value.ToFloat32()
	err := j.value.LastError()
	if err == nil {
		return value
	}

	if str, ok := j.value.GetInterface().(string); !ok {
		return 0
	} else if v, e := strconv.ParseFloat(str, 32); e == nil {
		return float32(v)
	}

	return 0
}

func (j *jsonValue) Float64() float64 {
	value := j.value.ToFloat64()
	err := j.value.LastError()
	if err == nil {
		return value
	}

	if str, ok := j.value.GetInterface().(string); !ok {
		return 0
	} else if value, err = strconv.ParseFloat(str, 64); err == nil {
		return value
	}

	return 0
}

func (j *jsonValue) String() string {
	value := j.value.ToString()
	err := j.value.LastError()
	if err == nil {
		return value
	}

	if str, ok := j.value.GetInterface().(string); !ok {
		return ""
	} else {
		return str
	}
}

func (j *jsonValue) Time() time.Time {
	if j.value.ValueType() == jsoniter.StringValue {
		val := j.value.ToString()
		err := j.value.LastError()
		if err != nil {
			return time.Time{}
		}

		if dt, e := dateparse.ParseLocal(val); e != nil {
			return time.Time{}
		} else {
			return dt
		}
	} else if j.value.ValueType() == jsoniter.NumberValue {
		val := j.value.ToInt64()
		err := j.value.LastError()
		if err != nil {
			return time.Time{}
		}

		return time.Unix(val, 0)
	}

	return time.Time{}
}

func (j *jsonValue) Duration() time.Duration {
	if j.value.ValueType() == jsoniter.StringValue {
		v := j.value.ToString()
		err := j.value.LastError()
		if err != nil {
			return 0
		}

		value, err := time.ParseDuration(v)
		if err != nil {
			return 0
		}

		return value
	} else if j.value.ValueType() == jsoniter.NumberValue {
		v := j.value.ToInt64()
		err := j.value.LastError()
		if err != nil {
			return 0
		}

		return time.Duration(v) * time.Second
	}

	return 0
}

func (j *jsonValue) StringSlice() []string {
	arr, ok := j.value.GetInterface().([]interface{})
	if !ok {
		return nil
	}

	strs := make([]string, len(arr))
	for idx, it := range arr {
		strs[idx] = fmt.Sprintf("%v", it)
	}

	return strs
}

func (j *jsonValue) StringMap() map[string]string {
	m, ok := j.value.GetInterface().(map[string]interface{})
	if !ok {
		return nil
	}

	res := map[string]string{}

	for k, v := range m {
		res[k] = fmt.Sprintf("%v", v)
	}

	return res
}

func (j *jsonValue) Scan(val interface{}) error {
	b, err := j.api.Marshal(j.value)
	if err != nil {
		return err
	} else if len(b) == 0 {
		return nil
	}

	return j.api.Unmarshal(b, val)
}

func (j *jsonValue) Bytes() []byte {
	b, err := j.api.Marshal(j.value)
	if err == nil {
		return b
	}

	return []byte{}
}
