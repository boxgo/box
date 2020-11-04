package logger

import (
	"github.com/boxgo/box/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	// Config of logger
	Config struct {
		path              string
		Level             zap.AtomicLevel        `config:"level" json:"level" desc:"debug,info,warn,error,dpanic,panic,fatal"`
		Development       bool                   `config:"development" json:"development"`
		DisableCaller     bool                   `config:"disableCaller" json:"disableCaller"`
		DisableStacktrace bool                   `config:"disableStacktrace" json:"disableStacktrace"`
		Sampling          *SamplingConfig        `config:"sampling" json:"sampling"`
		Encoding          string                 `config:"encoding" json:"encoding"`
		EncoderConfig     *EncoderConfig         `config:"encoderConfig" json:"encoderConfig"`
		OutputPaths       []string               `config:"outputPaths" json:"outputPaths"`
		ErrorOutputPaths  []string               `config:"errorOutputPaths" json:"errorOutputPaths"`
		InitialFields     map[string]interface{} `config:"initialFields" json:"initialFields"`
	}

	// SamplingConfig of logger
	SamplingConfig struct {
		Initial    int `config:"initial" json:"initial"`
		Thereafter int `config:"thereafter" json:"thereafter"`
	}

	// EncoderConfig of logger
	EncoderConfig struct {
		MessageKey     string `config:"messageKey" json:"messageKey"`
		LevelKey       string `config:"levelKey" json:"levelKey"`
		TimeKey        string `config:"timeKey" json:"timeKey"`
		NameKey        string `config:"nameKey" json:"nameKey"`
		CallerKey      string `config:"callerKey" json:"callerKey"`
		StacktraceKey  string `config:"stacktraceKey" json:"stacktraceKey"`
		LineEnding     string `config:"lineEnding" json:"lineEnding"`
		EncodeLevel    string `config:"levelEncoder" json:"levelEncoder"`
		EncodeTime     string `config:"timeEncoder" json:"timeEncoder"`
		EncodeDuration string `config:"durationEncoder" json:"durationEncoder"`
		EncodeCaller   string `config:"callerEncoder" json:"callerEncoder"`
		EncodeName     string `config:"nameEncoder" json:"nameEncoder"`
	}
)

// StdConfig new a logger from config path "logger.{{name}}"
func StdConfig(name string) *Config {
	cfg := DefaultConfig("logger." + name)

	config.Scan(cfg)

	return cfg
}

// DefaultConfig of logger
func DefaultConfig(path string) *Config {
	return &Config{
		path:              path,
		Level:             zap.NewAtomicLevel(),
		Development:       false,
		Encoding:          "console",
		DisableCaller:     true,
		DisableStacktrace: false,
		Sampling: &SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: &EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			MessageKey:     "msg",
			CallerKey:      "caller",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    "capitalColor",
			EncodeTime:     "iso8601",
			EncodeDuration: "ms",
			EncodeCaller:   "short",
		},
	}
}

// Build a logger
func (cfg *Config) Build() (*Logger, error) {
	return newLogger(cfg)
}

func (cfg *Config) Path() string {
	return cfg.path
}
