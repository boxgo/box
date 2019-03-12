package config

import (
	"context"
	"encoding/json"

	"github.com/boxgo/box/config/loader"
	"github.com/boxgo/box/config/printer"
	"github.com/boxgo/box/minibox"
	"github.com/boxgo/kit/logger"
)

type (
	// Config 配置管理器
	Config interface {
		Set(cfgs ...minibox.MiniBox) Config
		Scan(ctx context.Context)
	}

	configCenter struct {
		boxes  []minibox.MiniBox
		loader loader.Loader
	}
)

func (boxCfg *configCenter) Set(cfgs ...minibox.MiniBox) Config {
	boxCfg.boxes = append(boxCfg.boxes, cfgs...)

	return boxCfg
}

func (boxCfg *configCenter) Scan(ctx context.Context) {
	printer.PrintPrettyConfigMap(boxCfg.boxes)

	for _, cfg := range boxCfg.boxes {
		cfgHook, ok := cfg.(minibox.ConfigHook)

		if ok {
			cfgHook.ConfigWillLoad(ctx)
		}

		boxCfg.loader.Load(cfg)

		if cfgExt, ok := cfg.(minibox.MiniBoxExt); ok {
			for _, ext := range cfgExt.Exts() {
				boxCfg.loader.Load(ext)
			}
		}

		if ok {
			cfgHook.ConfigDidLoad(ctx)
		}

		cfgStr, err := json.Marshal(cfg)
		if err != nil {
			logger.Default.Errorf("config: %-20s marshal error: %s", cfg.Name(), err.Error())
		} else {
			logger.Default.Infof("config: %-20s %s", cfg.Name(), cfgStr)
		}
	}
}

// NewConfig 新建一个配置管理器
func NewConfig(lder loader.Loader) Config {
	return &configCenter{
		loader: lder,
	}
}
