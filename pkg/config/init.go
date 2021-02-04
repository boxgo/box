// +build !no_config_init

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/boxgo/box/pkg/config/source"
	"github.com/boxgo/box/pkg/config/source/file"
	"github.com/boxgo/box/pkg/util/fputil"
)

var (
	bootOK = loadBootConfig()
)

func loadBootConfig() bool {
	var (
		path string
		fps  []string
		cfg  = NewConfig()
	)

	if bootCfg := os.Getenv("BOX_BOOT_CONFIG"); fputil.IsFile(bootCfg) {
		fps = []string{bootCfg}
	} else {
		wd, err := os.Getwd()
		if err != nil {
			panic(fmt.Errorf("get work dir error: %s", err))
		}

		if bootCfg == "" {
			bootCfg = "box"
		}

		fps = []string{
			filepath.Join(wd, bootCfg+".yml"),
			filepath.Join(wd, bootCfg+".yaml"),
			filepath.Join(wd, bootCfg+".toml"),
			filepath.Join(wd, bootCfg+".json"),
			filepath.Join(wd, "config", bootCfg+".yml"),
			filepath.Join(wd, "config", bootCfg+".yaml"),
			filepath.Join(wd, "config", bootCfg+".toml"),
			filepath.Join(wd, "config", bootCfg+".json"),
		}
	}

	if path = fputil.FirstExistFilePath(fps); path == "" {
		panic(fmt.Errorf("config file\n%s\nnot found", strings.Join(fps, "\n")))
	}

	if err := cfg.Load(file.NewSource(file.WithPath(path))); err != nil {
		panic(fmt.Errorf("config load error: %s", err))
	}

	if err := cfg.Sync(); err != nil {
		panic(fmt.Errorf("config sync error: %s", err))
	}

	if err := cfg.Get().Scan(&bootCfg); err != nil {
		panic(fmt.Errorf("config scan error: %s", err))
	} else {
		defaultSources = make([]source.Source, len(bootCfg.Source))
	}

	if len(bootCfg.Source) == 0 { // if no source, use input file as config.
		defaultSources = make([]source.Source, 1)
		bootCfg.Source = append([]Source{}, Source{
			Type: "file",
			name: "file",
			data: []byte(fmt.Sprintf("{\"path\":\"%s\",\"type\":\"file\"}", path)),
		})

		return true
	}

	for idx, sour := range bootCfg.Source {
		bootCfg.Source[idx].name = sour.Type
		bootCfg.Source[idx].data = cfg.Get("source", strconv.Itoa(idx)).Bytes()
	}

	return true
}
