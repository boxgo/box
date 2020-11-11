package env

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/boxgo/box/pkg/config/source"
)

type strippedPrefixKey struct{}
type prefixKey struct{}

func WithConfig(data []byte) []source.Option {
	type (
		opt struct {
			Prefix      string `config:"prefix"`
			StripPrefix string `config:"stripPrefix"`
		}
	)

	v := &opt{}
	if err := json.Unmarshal(data, v); err != nil {
		log.Fatal(err)
	}

	return []source.Option{
		WithPrefix(v.Prefix),
		WithStrippedPrefix(v.StripPrefix),
	}
}

// WithStrippedPrefix sets the environment variable prefixes to scope to.
// These prefixes will be removed from the actual config entries.
func WithStrippedPrefix(p ...string) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}

		o.Context = context.WithValue(o.Context, strippedPrefixKey{}, appendUnderscore(p))
	}
}

// WithPrefix sets the environment variable prefixes to scope to.
// These prefixes will not be removed. Each prefix will be considered a top level config entry.
func WithPrefix(p ...string) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, prefixKey{}, appendUnderscore(p))
	}
}

func appendUnderscore(prefixes []string) []string {
	var result []string
	for _, p := range prefixes {
		if !strings.HasSuffix(p, "_") {
			result = append(result, p+"_")
			continue
		}

		result = append(result, p)
	}

	return result
}
