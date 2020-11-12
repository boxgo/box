// +build !configinit

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/boxgo/box/pkg/config/source"
	"github.com/boxgo/box/pkg/config/source/env"
	"github.com/boxgo/box/pkg/config/source/etcd"
	"github.com/boxgo/box/pkg/config/source/file"
	"github.com/boxgo/box/pkg/config/source/redis"
	"github.com/boxgo/box/pkg/util/fputil"
)

var (
	firstInitCfg    = NewConfig()
	firstInitCfgCfg = struct {
		Loader string `config:"loader"`
		Reader string `config:"reader"`
		Source []struct {
			Type string `config:"type"`
		} `config:"source"`
	}{}
)

func init() {
	var (
		wd      string
		path    string
		err     error
		sources []source.Source
	)

	if wd, err = os.Getwd(); err != nil {
		panic(fmt.Errorf("get work dir error: %s", err))
	}

	fps := []string{
		filepath.Join(wd, "box.yaml"),
		filepath.Join(wd, "box.toml"),
		filepath.Join(wd, "box.json"),
	}
	if path = fputil.FirstExistFilePath(fps); path == "" {
		panic(fmt.Errorf("config file\n%s\nnot found", strings.Join(fps, "\n")))
	}

	if err := firstInitCfg.Load(file.NewSource(file.WithPath(path))); err != nil {
		panic(fmt.Errorf("config load error: %s", err))
	}

	if err := firstInitCfg.Sync(); err != nil {
		panic(fmt.Errorf("config sync error: %s", err))
	}

	if err := firstInitCfg.Get().Scan(&firstInitCfgCfg); err != nil {
		panic(fmt.Errorf("config scan error: %s", err))
	}

	for idx, sour := range firstInitCfgCfg.Source {
		cfgData := firstInitCfg.Get("source", strconv.Itoa(idx)).Bytes()

		switch sour.Type {
		case "etcd":
			sources = append(sources, etcd.NewSource(etcd.WithConfig(cfgData)...))
		case "env":
			sources = append(sources, env.NewSource(env.WithConfig(cfgData)...))
		case "file":
			sources = append(sources, file.NewSource(file.WithConfig(cfgData)...))
		case "redis":
			sources = append(sources, redis.NewSource(redis.WithConfig(cfgData)...))
		}
	}

	if err := Default.Load(sources...); err != nil {
		panic(fmt.Errorf("default config load error: %s", err))
	}
}
