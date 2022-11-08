// Package config is an interface for dynamic configuration.
package config

import (
	"github.com/boxgo/box/v2/config/reader"
	"github.com/boxgo/box/v2/config/source"
	"github.com/boxgo/box/v2/config/value"
	"github.com/boxgo/box/v2/config/value/json"
)

type (
	// Configurator is an interface abstraction for dynamic configuration
	Configurator interface {
		Load() error                  // Load config sources
		Scan(val interface{}) error   // Scan to val
		Value(key string) value.Value // Value get values through key
		// Watch(key string, cb func(old, new value.Value)) error // TODO Watch field change
		// Close() error                                          // TODO Close stop the config loader/watcher
	}

	Config interface {
		Path() string
	}

	config struct {
		opts   Options
		values value.Values
	}
)

// NewConfig returns new config
func NewConfig(opts ...Option) Configurator {
	options := Options{
		Reader: reader.NewReader(),
	}

	for _, o := range opts {
		o(&options)
	}

	values, _ := json.NewValues(nil)

	c := &config{
		opts:   options,
		values: values,
	}

	return c
}

func (c *config) Load() error {
	changes := make([]*source.ChangeSet, len(c.opts.Source))

	for idx, sour := range c.opts.Source {
		data, err := sour.Read()
		if err != nil {
			return err
		}

		changes[idx] = data
	}

	if mergeData, err := c.opts.Reader.Merge(changes...); err != nil {
		return err
	} else if c.values, err = json.NewValues(mergeData); err != nil {
		return err
	}

	return nil
}

func (c *config) Scan(val interface{}) (err error) {
	switch v := val.(type) {
	case Config:
		err = c.Value(v.Path()).Scan(val)
	default:
		err = c.values.Scan(val)
	}

	if err != nil {
		return
	}

	if c.opts.Validator != nil {
		err = c.opts.Validator.Validate(val)
	}

	return
}

func (c *config) Value(key string) value.Value {
	return c.values.Value(key)
}

func (c *config) Watch(key string, cb func(old value.Value, new value.Value)) error {
	// TODO implement me
	panic("implement me")
}

func (c *config) Close() error {
	// TODO implement me
	panic("implement me")
}
