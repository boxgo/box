// Package codec handles data encoding
package codec

type (
	Coder interface {
		String() string
		Marshal(v interface{}) ([]byte, error)   // Marshal returns the encoded data of v.
		Unmarshal(d []byte, v interface{}) error // Unmarshal parses the encoded data of d and stores the result in the value pointed to by v.
	}
)
