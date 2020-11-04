// Package config is an interface for dynamic configuration.
package config

import (
	"github.com/boxgo/box/pkg/config/reader"
	"github.com/boxgo/box/pkg/config/source"
)

type (
	// Configurator is an interface abstraction for dynamic configuration
	Configurator interface {
		// Load config sources
		Load(source ...source.Source) error
		// Force a source changeset sync
		Sync() error
		// Stop the config loader/watcher
		Close() error
		// Watch field change
		Watch(path ...string) (Watcher, error)
		// Get value through field
		Get(path ...string) reader.Value
		// Scan to val
		Scan(val interface{}) error
		// Bytes get merged config data
		Bytes() []byte
	}
)

var (
	// Default Config Manager
	Default = NewConfig()
)

// NewConfig returns new config
func NewConfig(opts ...Option) Configurator {
	return newConfig(opts...)
}

// Load config sources
func Load(source ...source.Source) error {
	return Default.Load(source...)
}

// Sync force a source changeset sync
func Sync() error {
	return Default.Sync()
}

// Watch a value for changes
func Watch(path ...string) (Watcher, error) {
	return Default.Watch(path...)
}

// Close stop the config loader/watcher
func Close() error {
	return Default.Close()
}

func Scan(val interface{}) error {
	return Default.Scan(val)
}

// Get a value from the config
func Get(path ...string) reader.Value {
	return Default.Get(path...)
}

func Byte() []byte {
	return Default.Bytes()
}
