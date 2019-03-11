package loader

import (
	"strings"

	"github.com/boxgo/box/minibox"
	"github.com/micro/go-config"
	"github.com/micro/go-config/source/env"
	"github.com/micro/go-config/source/file"
)

type (
	//Micro go-config
	Micro struct {
		config config.Config
	}
)

// Load load config to minibox
func (micro *Micro) Load(cfg minibox.MiniBox) {
	if cfg.Name() == "" {
		micro.config.Scan(cfg)
	} else {
		paths := strings.Split(cfg.Name(), ".")
		micro.config.Get(paths...).Scan(cfg)
	}
}

// NewFileEnvConfig new a default app config from file and env
func NewFileEnvConfig(path string) Loader {
	boxCfg := &Micro{
		config: config.NewConfig(),
	}

	boxCfg.config.Load(
		file.NewSource(file.WithPath(path)),
		env.NewSource(),
	)

	return boxCfg
}

// NewFileConfig new a default app config from file
func NewFileConfig(path string) Loader {
	boxCfg := &Micro{
		config: config.NewConfig(),
	}

	boxCfg.config.Load(
		file.NewSource(file.WithPath(path)),
	)

	return boxCfg
}
