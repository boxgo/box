package file

import (
	"context"
	"encoding/json"
	"log"

	"github.com/boxgo/box/pkg/config/source"
)

type filePathKey struct{}

func WithConfig(data []byte) []source.Option {
	type (
		opt struct {
			Path string `config:"path"`
		}
	)

	v := &opt{}
	if err := json.Unmarshal(data, v); err != nil {
		log.Fatal(err)
	}

	return []source.Option{
		WithPath(v.Path),
	}
}

// WithPath sets the path to file
func WithPath(p string) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, filePathKey{}, p)
	}
}
