package logger

import (
	"fmt"
	"time"

	"github.com/boxgo/box/pkg/config"
	"github.com/boxgo/box/pkg/logger/core"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	// Config of logger
	Config struct {
		path              string
		SplitLen          int                    `config:"splitLen" json:"splitLen" desc:"split to multi line when reach the len."`
		WatchInterval     time.Duration          `config:"watchInterval" json:"watchInterval" desc:"config change watch interval, default is 5s"`
		Level             zap.AtomicLevel        `config:"level" json:"level" desc:"debug,info,warn,error,panic,fatal"`
		DisableCaller     bool                   `config:"disableCaller" json:"disableCaller"`
		DisableStacktrace bool                   `config:"disableStacktrace" json:"disableStacktrace"`
		Sampling          *SamplingConfig        `config:"sampling" json:"sampling"`
		Encoding          string                 `config:"encoding" json:"encoding"`
		EncoderConfig     *EncoderConfig         `config:"encoderConfig" json:"encoderConfig"`
		OutputPaths       []string               `config:"outputPaths" json:"outputPaths"`
		ErrorOutputPaths  []string               `config:"errorOutputPaths" json:"errorOutputPaths"`
		InitialFields     map[string]interface{} `config:"initialFields" json:"initialFields"`
		Mask              bool                   `config:"mask" json:"mask"`
		MaskRules         []core.MaskRule        `config:"maskRules" json:"maskRules"`
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
func StdConfig(key string) *Config {
	cfg := DefaultConfig(key)

	if err := config.Scan(cfg); err != nil {
		panic(fmt.Errorf("logger build error: %s", err))
	}

	return cfg
}

// DefaultConfig of logger
func DefaultConfig(key string) *Config {
	return &Config{
		path:              "logger." + key,
		SplitLen:          5 * 1024,
		WatchInterval:     time.Second * 10,
		Level:             zap.NewAtomicLevel(),
		Encoding:          "console",
		DisableCaller:     true,
		DisableStacktrace: true,
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
			EncodeLevel:    "capital",
			EncodeTime:     "iso8601",
			EncodeDuration: "ms",
			EncodeCaller:   "short",
		},
		Mask: true,
		MaskRules: []core.MaskRule{
			{Rule: `"password":(\s*)".*?"`, Replace: `"password":$1"*"`},
			{Rule: `password:(\s*).*?\S*`, Replace: `password:$1*`},
			{Rule: `password=\w*&`, Replace: `password=*&`},
			{Rule: `password=\w*\S`, Replace: `password=*`},
			{Rule: `\\"password\\":(\s*)\\".*?\\"`, Replace: `\"password\":$1\"*\"`},
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
