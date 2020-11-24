package debug

type (
	Config struct{}
)

func StdConfig() *Config {
	cfg := DefaultConfig()

	return cfg
}

func DefaultConfig() *Config {
	return &Config{}
}

func (c *Config) Path() string {
	return "debug"
}

func (c *Config) Build() *Debug {
	return newDebug(c)
}
