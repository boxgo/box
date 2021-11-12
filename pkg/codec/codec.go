// Package codec handles source encoding formats
package codec

type (
	Marshaler interface {
		String() string
		Marshal(interface{}) ([]byte, error)
		Unmarshal([]byte, interface{}) error
	}
)
