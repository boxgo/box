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
		err  error
		wd   string
		path string
		cfg  = NewConfig()
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

	for idx, sour := range bootCfg.Source {
		bootCfg.Source[idx].name = sour.Type
		bootCfg.Source[idx].data = cfg.Get("source", strconv.Itoa(idx)).Bytes()
	}

	return true
}
