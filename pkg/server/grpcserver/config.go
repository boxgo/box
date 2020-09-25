package grpcserver

import (
	"github.com/boxgo/box/pkg/config"
	"github.com/boxgo/box/pkg/config/field"
)

type (
	Config struct {
		config.SubConfigurator
		network *field.Field
		address *field.Field
	}
)

func newConfig(name string, cfg config.SubConfigurator) *Config {
	c := &Config{
		SubConfigurator: cfg,
		network:         field.New(true, "grpc.server", "network", "The network must be \"tcp\", \"tcp4\", \"tcp6\", \"unix\" or \"unixpacket\".", "tcp4"),
		address:         field.New(true, "grpc.server", "address", "format: host:port", ":9092"),
	}

	c.Mount(c.Fields()...)

	return c
}

func (cfg *Config) Address() string {
	return cfg.GetString(cfg.address)
}

func (cfg *Config) Network() string {
	return cfg.GetString(cfg.network)
}

func (cfg *Config) Fields() []*field.Field {
	return []*field.Field{
		cfg.network,
		cfg.address,
	}
}
