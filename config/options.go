package config

import (
	"context"

	"github.com/boxgo/box/v2/config/reader"
	"github.com/boxgo/box/v2/config/source"
	"github.com/boxgo/box/v2/config/validator"
)

type (
	Options struct {
		Source    []source.Source
		Reader    reader.Reader
		Validator validator.Validator
		Context   context.Context
	}

	Option func(o *Options)
)

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

// WithValidator sets the config validator
func WithValidator(r validator.Validator) Option {
	return func(o *Options) {
		o.Validator = r
	}
}
