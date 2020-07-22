package grpcserver

import (
	"github.com/boxgo/box/pkg/config"
)

type (
	Config struct {
		config.SubConfigurator
		network *config.Field
		address *config.Field
	}
)

func newConfig(name string, cfg config.SubConfigurator) *Config {
	return &Config{
		SubConfigurator: cfg,
		network:         config.NewField(name, "network", "The network must be \"tcp\", \"tcp4\", \"tcp6\", \"unix\" or \"unixpacket\".", "tcp"),
		address:         config.NewField(name, "address", "format: host:port", ":9092"),
	}

}

func (cfg *Config) Address() string {
	return cfg.GetString(cfg.address)
}

func (cfg *Config) Network() string {
	return cfg.GetString(cfg.network)
}

func (cfg *Config) Fields() []*config.Field {
	return []*config.Field{
		cfg.network,
		cfg.address,
	}
}
