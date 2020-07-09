package json

import (
	"fmt"
	"strconv"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type jsonValue struct {
	value jsoniter.Any
	api   jsoniter.API
}

func (j *jsonValue) Bool(def bool) bool {
	value := j.value.ToBool()
	err := j.value.LastError()
	if err == nil {
		return value
	}

	str, ok := j.value.GetInterface().(string)
	if !ok {
		return def
	}

	value, err = strconv.ParseBool(str)
	if err == nil {
		return value
	}

	return def
}

func (j *jsonValue) Int(def int) int {
	value := j.value.ToInt()
	err := j.value.LastError()
	if err == nil {
		return value
	}

	str, ok := j.value.GetInterface().(string)
	if !ok {
		return def
	}

	value, err = strconv.Atoi(str)
	if err == nil {
		return value
	}

	return def
}

func (j *jsonValue) Uint(def uint) uint {
	value := j.value.ToUint()
	err := j.value.LastError()
	if err == nil {
		return value
	}

	str, ok := j.value.GetInterface().(string)
	if !ok {
		return def
	}

	v, err := strconv.ParseUint(str, 10, 64)
	if err == nil {
		return uint(v)
	}

	return def
}

func (j *jsonValue) String(def string) string {
	value := j.value.ToString()
	err := j.value.LastError()
	if err == nil {
		return value
	}

	str, ok := j.value.GetInterface().(string)
	if !ok {
		return def
	}

	return str
}

func (j *jsonValue) Float64(def float64) float64 {
	value := j.value.ToFloat64()
	err := j.value.LastError()
	if err == nil {
		return value
	}

	str, ok := j.value.GetInterface().(string)
	if !ok {
		return def
	}

	value, err = strconv.ParseFloat(str, 64)
	if err == nil {
		return value
	}

	return def
}

func (j *jsonValue) Duration(def time.Duration) time.Duration {
	if j.value.ValueType() == jsoniter.StringValue {
		v := j.value.ToString()
		err := j.value.LastError()
		if err != nil {
			return def
		}

		value, err := time.ParseDuration(v)
		if err != nil {
			return def
		}

		return value
	} else if j.value.ValueType() == jsoniter.NumberValue {
		v := j.value.ToInt64()
		err := j.value.LastError()
		if err != nil {
			return def
		}

		return time.Duration(v) * time.Second
	}

	return def
}

func (j *jsonValue) StringSlice(def []string) []string {
	arr, ok := j.value.GetInterface().([]interface{})
	if !ok {
		return def
	}

	strs := make([]string, len(arr))
	for idx, it := range arr {
		strs[idx] = fmt.Sprintf("%v", it)
	}

	return strs
}

func (j *jsonValue) StringMap(def map[string]string) map[string]string {
	m, ok := j.value.GetInterface().(map[string]interface{})
	if !ok {
		return def
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
