// Package config is an interface for dynamic configuration.
package config

import (
	"time"

	"github.com/boxgo/box/pkg/config/reader"
	"github.com/boxgo/box/pkg/config/source"
	"github.com/boxgo/box/pkg/config/source/file"
)

type (
	// Configurator is an interface abstraction for dynamic configuration
	Configurator interface {
		// Load config sources
		Load(source ...source.Source) error
		// Force a source changeset sync
		Sync() error
		// Watch a value for changes
		Watch(field *Field) (Watcher, error)
		// Stop the config loader/watcher
		Close() error
		// SprintFields registered fields
		SprintFields() string
		// SprintTemplate through encoder
		SprintTemplate(encoder string) string
		// Mounter mount fields to configurator
		Mounter
		// Getter get value through field
		Getter
	}

	// Configurable instance should implements this interface
	Config interface {
		Name() string
		Configure(mg MountGetter)
	}

	// Mounter fields
	Mounter interface {
		Mount(fields ...*Field)
	}

	Getter interface {
		// Get value through field
		Get(field *Field) reader.Value
		// GetString through field
		GetBool(field *Field) bool
		// GetInt through field
		GetInt(field *Field) int
		// GetUint through field
		GetUint(field *Field) uint
		// GetString through field
		GetString(field *Field) string
		// GetFloat64 through field
		GetFloat64(field *Field) float64
		// GetDuration through field
		GetDuration(field *Field) time.Duration
		// GetStringSlice through field
		GetStringSlice(field *Field) []string
		// GetStringMap through field
		GetStringMap(field *Field) map[string]string
	}

	MountGetter interface {
		Mounter
		Getter
	}
)

var (
	// Default Config Manager
	DefaultConfig = NewConfig()
)

// NewConfig returns new config
func NewConfig(opts ...Option) Configurator {
	return newConfig(opts...)
}

// LoadFile is short hand for creating a file source and loading it
func LoadFile(path string) error {
	return Load(file.NewSource(
		file.WithPath(path),
	))
}

// Load config sources
func Load(source ...source.Source) error {
	return DefaultConfig.Load(source...)
}

// Force a source changeset sync
func Sync() error {
	return DefaultConfig.Sync()
}

// Watch a value for changes
func Watch(field *Field) (Watcher, error) {
	return DefaultConfig.Watch(field)
}

// Stop the config loader/watcher
func Close() error {
	return DefaultConfig.Close()
}

// Mount fields
func Mount(fields ...*Field) {
	DefaultConfig.Mount(fields...)
}

// Get a value from the config
func Get(field *Field) reader.Value {
	return DefaultConfig.Get(field)
}

// GetString through field
func GetBool(field *Field) bool {
	return DefaultConfig.GetBool(field)
}

// GetInt through field
func GetInt(field *Field) int {
	return DefaultConfig.GetInt(field)
}

// GetUint through field
func GetUint(field *Field) uint {
	return DefaultConfig.GetUint(field)
}

// GetString through field
func GetString(field *Field) string {
	return DefaultConfig.GetString(field)
}

// GetFloat64 through field
func GetFloat64(field *Field) float64 {
	return DefaultConfig.GetFloat64(field)
}

// GetDuration through field
func GetDuration(field *Field) time.Duration {
	return DefaultConfig.GetDuration(field)
}

// GetStringSlice through field
func GetStringSlice(field *Field) []string {
	return DefaultConfig.GetStringSlice(field)
}

// GetStringMap through field
func GetStringMap(field *Field) map[string]string {
	return DefaultConfig.GetStringMap(field)
}

// SprintFields registered fields
func SprintFields() string {
	return DefaultConfig.SprintFields()
}

// SprintTemplate through encoder
func SprintTemplate(encoder string) string {
	return DefaultConfig.SprintTemplate(encoder)
}
