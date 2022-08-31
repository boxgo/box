// Package reader parses change sets and provides config values
package reader

import (
	"time"

	"github.com/boxgo/box/v2/codec"
	"github.com/boxgo/box/v2/codec/json"
	"github.com/boxgo/box/v2/codec/toml"
	"github.com/boxgo/box/v2/codec/xml"
	"github.com/boxgo/box/v2/codec/yaml"
	"github.com/boxgo/box/v2/config/source"
	"github.com/imdario/mergo"
)

type (
	// Reader is an interface for merging change sets
	Reader interface {
		Merge(...*source.ChangeSet) (*source.ChangeSet, error) // merge multi to one
	}

	reader struct {
		targetCoder  codec.Coder
		sourceCoders map[string]codec.Coder
	}
)

// NewReader creates a coder reader
func NewReader(opts ...Option) Reader {
	options := Options{
		TargetEncoder: json.NewCoder(),
		SourceEncoders: map[string]codec.Coder{
			"json": json.NewCoder(),
			"yml":  yaml.NewCoder(),
			"yaml": yaml.NewCoder(),
			"toml": toml.NewCoder(),
			"xml":  xml.NewCoder(),
		},
	}

	for _, o := range opts {
		o(&options)
	}

	return &reader{
		targetCoder:  options.TargetEncoder,
		sourceCoders: options.SourceEncoders,
	}
}

func (r *reader) Merge(changes ...*source.ChangeSet) (*source.ChangeSet, error) {
	var merged map[string]interface{}

	for _, change := range changes {
		if change == nil || len(change.Data) == 0 {
			continue
		}

		coder, ok := r.sourceCoders[change.Format]
		if !ok {
			// fallback
			coder = r.targetCoder
		}

		var data map[string]interface{}
		if err := coder.Unmarshal(change.Data, &data); err != nil {
			return nil, err
		}

		if err := mergo.Map(&merged, data, mergo.WithOverride); err != nil {
			return nil, err
		}
	}

	data, err := r.targetCoder.Marshal(merged)
	if err != nil {
		return nil, err
	}

	cs := &source.ChangeSet{
		Timestamp: time.Now(),
		Data:      data,
		Source:    "reader",
		Format:    r.targetCoder.String(),
	}
	cs.Checksum = cs.Sum()

	return cs, nil
}
