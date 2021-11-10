package file

import (
	"strings"

	"github.com/boxgo/box/pkg/codec"
)

func format(p string, e codec.Marshaler) string {
	parts := strings.Split(p, ".")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return e.String()
}
