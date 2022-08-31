package json

import (
	"strconv"
	"strings"
	"time"

	"github.com/boxgo/box/v2/config/source"
	"github.com/boxgo/box/v2/config/value"
	jsoniter "github.com/json-iterator/go"
)

type (
	jsonValues struct {
		ch  *source.ChangeSet
		api jsoniter.API
	}
)

func NewValues(ch *source.ChangeSet) (value.Values, error) {
	api := jsoniter.Config{
		EscapeHTML:             true,
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
		TagKey:                 "config",
	}.Froze()

	if ch == nil {
		ch = &source.ChangeSet{
			Data:      nil,
			Checksum:  "",
			Format:    "json",
			Source:    "",
			Timestamp: time.Now(),
		}
		ch.Sum()
	}

	return &jsonValues{
		ch:  ch,
		api: api,
	}, nil
}

func (j *jsonValues) Bytes() []byte {
	return j.ch.Data
}

func (j *jsonValues) Value(key string) value.Value {
	var paths []interface{}

	if key != "" {
		items := strings.Split(key, ".")
		paths = make([]interface{}, len(items))

		for idx, it := range items {
			if n, err := strconv.ParseInt(it, 10, 32); err == nil {
				paths[idx] = int(n)
			} else {
				paths[idx] = it
			}
		}
	}

	return &jsonValue{
		api:   j.api,
		value: j.api.Get(j.Bytes(), paths...),
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
