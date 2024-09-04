package config

import (
	"context"

	"github.com/boxgo/box/pkg/config/loader"
	"github.com/boxgo/box/pkg/config/reader"
	"github.com/boxgo/box/pkg/config/source"
	"github.com/boxgo/box/pkg/config/validator"
)

type (
	Options struct {
		Loader    loader.Loader
		Reader    reader.Reader
		Validator validator.Validator
		Source    []source.Source

		// for alternative data
		Context context.Context
	}

	Option func(o *Options)
)

// WithLoader sets the loader for manager config
func WithLoader(l loader.Loader) Option {
	return func(o *Options) {
		o.Loader = l
	}
}

// WithValidator sets the loader for manager config
func WithValidator(v validator.Validator) Option {
	return func(o *Options) {
		o.Validator = v
	}
}

// WithSource appends a source to list of sources
func WithSource(s ...source.Source) Option {
	return func(o *Options) {
		o.Source = append(o.Source, s...)
	}
}

// WithReader sets the config reader
func WithReader(r reader.Reader) Option {
	return func(o *Options) {
		o.Reader = r
	}
}
