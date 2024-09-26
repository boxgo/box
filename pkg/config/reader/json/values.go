package json

import (
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/boxgo/box/pkg/config/reader"
	"github.com/boxgo/box/pkg/config/source"
	jsoniter "github.com/json-iterator/go"
)

type jsonValues struct {
	ch  *source.ChangeSet
	api jsoniter.API
}

const defaultTagKey = "config"

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
}

func newValues(ch *source.ChangeSet) (reader.Values, error) {
	api := jsoniter.Config{
		EscapeHTML:             true,
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
		CaseSensitive:          true,
		TagKey:                 defaultTagKey,
	}.Froze()

	return &jsonValues{
		ch:  ch,
		api: api,
	}, nil
}

func (j *jsonValues) Bytes() []byte {
	return j.ch.Data
}

func (j *jsonValues) Get(path ...string) reader.Value {
	var p []interface{}
	for _, pit := range path {
		for _, it := range strings.Split(pit, ".") {
			if i, err := strconv.ParseInt(it, 10, 32); err == nil {
				p = append(p, int(i))
			} else {
				p = append(p, it)
			}
		}
	}

	return &jsonValue{
		api:   j.api,
		value: j.api.Get(j.Bytes(), p...),
	}
}

func (j *jsonValues) Map() map[string]interface{} {
	return j.api.Get(j.ch.Data).GetInterface().(map[string]interface{})
}

func (j *jsonValues) Scan(v interface{}) error {
	if len(j.ch.Data) == 0 {
		return nil
	}

	return j.api.Unmarshal(j.ch.Data, v)
}
