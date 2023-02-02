package config

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/boxgo/box/v2/config/reader"
	"github.com/boxgo/box/v2/config/source"
	"github.com/boxgo/box/v2/config/validator"
)

type (
	Options struct {
		Debug     io.Writer
		Interval  time.Duration
		Source    []source.Source
		Reader    reader.Reader
		Validator validator.Validator
		Context   context.Context
	}

	Option func(o *Options)
)

// WithDebug print debug log
func WithDebug(debug bool) Option {
	return func(o *Options) {
		if debug {
			o.Debug = os.Stdout
		} else {
			o.Debug = io.Discard
		}
	}
}

// WithInterval sets watch interval
func WithInterval(interval time.Duration) Option {
	return func(o *Options) {
		o.Interval = interval
	}
}

// WithSource appends a source to list of sources, the latter has higher priority.
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
