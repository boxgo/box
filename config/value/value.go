package value

import (
	"time"
)

type (
	// Values is returned by the reader
	Values interface {
		Bytes() []byte
		Value(key string) Value
		Scan(val interface{}) error
		Map() map[string]interface{}
	}

	// Value represents a value of any type
	Value interface {
		Bool() bool
		Uint() uint
		Int() int
		Int32() int32
		Int64() int64
		Float32() float32
		Float64() float64
		String() string
		Time() time.Time
		Duration() time.Duration
		StringSlice() []string
		StringMap() map[string]string
		Scan(val interface{}) error
		Bytes() []byte
	}
)
