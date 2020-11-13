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

type (
	sourceConfig struct {
		idx  int
		name string
		data []byte
	}
)

var (
	bootOK        = bootConfig()
	sourceConfigs []sourceConfig
)

func bootConfig() bool {
	var (
		err        error
		wd         string
		path       string
		bootCfg    = NewConfig()
		bootCfgCfg = struct {
			Loader string `config:"loader"`
			Reader string `config:"reader"`
			Source []struct {
				Type string `config:"type"`
			} `config:"source"`
		}{}
	)

	if wd, err = os.Getwd(); err != nil {
		panic(fmt.Errorf("get work dir error: %s", err))
	}

	fps := []string{
		filepath.Join(wd, "box.yml"),
		filepath.Join(wd, "box.yaml"),
		filepath.Join(wd, "box.toml"),
		filepath.Join(wd, "box.json"),
	}
	if path = fputil.FirstExistFilePath(fps); path == "" {
		panic(fmt.Errorf("config file\n%s\nnot found", strings.Join(fps, "\n")))
	}

	if err := bootCfg.Load(file.NewSource(file.WithPath(path))); err != nil {
		panic(fmt.Errorf("config load error: %s", err))
	}

	if err := bootCfg.Sync(); err != nil {
		panic(fmt.Errorf("config sync error: %s", err))
	}

	if err := bootCfg.Get().Scan(&bootCfgCfg); err != nil {
		panic(fmt.Errorf("config scan error: %s", err))
	}

	defaultSources = make([]source.Source, len(bootCfgCfg.Source))
	sourceConfigs = make([]sourceConfig, len(bootCfgCfg.Source))
	for idx, sour := range bootCfgCfg.Source {
		sourceConfigs[idx] = sourceConfig{
			idx:  idx,
			name: sour.Type,
			data: bootCfg.Get("source", strconv.Itoa(idx)).Bytes(),
		}
	}

	return true
}
