package config

import (
	"context"
	"encoding/json"
	"flag"
	"strings"

	"github.com/boxgo/box/minibox"
	"github.com/boxgo/config"
	"github.com/boxgo/config/source"
	"github.com/boxgo/logger"
	"github.com/spf13/pflag"
)

type (
	// Config register, dispatcher
	Config interface {
		EnablePflag(bool) Config                 // auto inject pflag
		EnableFlag(bool) Config                  // auto inject flag
		Register(cfgs ...minibox.MiniBox) Config // register box config
		Delivery(ctx context.Context) Config     // delivery config to box
	}

	configCenter struct {
		enablePflag bool
		enableFlag  bool
		config      config.Config
		boxes       []minibox.MiniBox
		sources     []source.Source
	}
)

func (cc *configCenter) EnablePflag(enable bool) Config {
	cc.enablePflag = enable
	cc.enableFlag = !enable

	return cc
}

func (cc *configCenter) EnableFlag(enable bool) Config {
	cc.enablePflag = !enable
	cc.enableFlag = enable

	panic("Flag inject not supported")
}

func (cc *configCenter) Register(cfgs ...minibox.MiniBox) Config {
	cc.boxes = append(cc.boxes, cfgs...)

	return cc
}

func (cc *configCenter) Delivery(ctx context.Context) Config {
	cc.prepare()

	for _, cfg := range cc.boxes {
		cfgHook, ok := cfg.(minibox.ConfigHook)

		if ok {
			cfgHook.ConfigWillLoad(ctx)
		}

		cc.load(cfg)

		if cfgExt, ok := cfg.(minibox.MiniBoxExt); ok {
			for _, ext := range cfgExt.Exts() {
				cc.load(ext)
			}
		}

		if ok {
			cfgHook.ConfigDidLoad(ctx)
		}

		cfgStr, err := json.Marshal(cfg)
		if err != nil {
			logger.Errorf("config %-20s marshal error: %s", cfg.Name(), err.Error())
		} else {
			logger.Infof("config %-20s %s", cfg.Name(), cfgStr)
		}
	}

	return cc
}

func (cc *configCenter) prepare() {
	for _, cfg := range cc.boxes {
		for _, f := range getFieldsInBox(cfg) {
			bindFieldPflag(f)
		}
	}

	if cc.enablePflag && !pflag.Parsed() {
		pflag.Parse()
	} else if cc.enableFlag && !flag.Parsed() {
		flag.Parse()
	}

	if err := cc.config.Sync(); err != nil {
		logger.Panic("Config.Sync.Error", err)
	}
	if err := cc.config.Load(cc.sources...); err != nil {
		logger.Panic("Config.Load.Error", err)
	}
}

func (cc *configCenter) load(cfg minibox.MiniBox) {
	var err error
	if cfg.Name() == "" {
		err = cc.config.Scan(cfg)
	} else {
		paths := strings.Split(cfg.Name(), ".")
		err = cc.config.Get(paths...).Scan(cfg)
	}

	if err != nil {
		logger.Errorf("Load.%s.Error: %s", cfg.Name(), err.Error())
	}
}

// NewConfig 新建一个配置管理器
func NewConfig(sources ...source.Source) Config {
	return &configCenter{
		config:  config.NewConfig(),
		sources: sources,
	}
}
