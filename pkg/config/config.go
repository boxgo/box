// Package config is an interface for dynamic configuration.
package config

import (
	"fmt"
	"sync"

	"github.com/boxgo/box/pkg/config/field"
	"github.com/boxgo/box/pkg/config/reader"
	"github.com/boxgo/box/pkg/config/source"
)

type (
	// Config is set of config fields. "." is a level splitter.
	// For example:
	//	father.child
	// It means that config data's struct same as this.
	// 	{
	// 	  "father": {
	// 	    "child": {xxx}
	// 	  }
	// 	}
	Config interface {
		Path() string
	}
	// Configurator is an interface abstraction for dynamic configuration
	Configurator interface {
		// Load config sources
		Load(source ...source.Source) error
		// Force a source change set sync
		Sync() error
		// Stop the config loader/watcher
		Close() error
		// Bytes get merged config data
		Bytes() []byte
		// Scan to val
		Scan(val Config) error
		// Watch field change
		Watch(path ...string) (Watcher, error)
		// Get value through field
		Get(path ...string) reader.Value
		// Scanned fields
		Fields() *field.Fields
	}

	bootConfig struct {
		Name    string `config:"name"`
		Version string `config:"version"`
		Loader  string `config:"loader"`
		Reader  string `config:"reader"`
		Source  []struct {
			Type string `config:"type"`
			name string
			data []byte
		} `config:"source"`
	}
)

var (
	// Default Config Manager
	Default        = NewConfig()
	bootCfg        = bootConfig{Name: "box", Version: "unknown"}
	defaultOnce    sync.Once
	defaultSources []source.Source
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
	lazyLoad()
	return Default.Sync()
}

// Close stop the config loader/watcher
func Close() error {
	return Default.Close()
}

// Byte return config raw data
func Byte() []byte {
	lazyLoad()
	return Default.Bytes()
}

// Scan config to val
func Scan(val Config) error {
	lazyLoad()
	return Default.Scan(val)
}

// Watch a value for changes
func Watch(path ...string) (Watcher, error) {
	lazyLoad()
	return Default.Watch(path...)
}

// Get a value from the config
func Get(path ...string) reader.Value {
	lazyLoad()
	return Default.Get(path...)
}

func Fields() *field.Fields {
	lazyLoad()
	return Default.Fields()
}

func ServiceName() string {
	return bootCfg.Name
}

func ServiceVersion() string {
	return bootCfg.Version
}

func lazyLoad() {
	defaultOnce.Do(func() {
		var validSources []source.Source
		for _, s := range defaultSources {
			if s != nil && s.String() != "" {
				validSources = append(validSources, s)
			}
		}

		if err := Default.Load(validSources...); err != nil {
			panic(fmt.Errorf("default load %s error: %s", validSources, err))
		}
	})
}
