package locales

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/boxgo/box/v2/codec/toml"
	"github.com/boxgo/box/v2/codec/yaml"
	"github.com/boxgo/box/v2/logger"
	"github.com/boxgo/box/v2/server/ginserver"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type (
	Locales struct {
		cfg       *Config
		bundle    *i18n.Bundle
		languages sync.Map
	}
)

func newLocales(c *Config) *Locales {
	bundle := i18n.NewBundle(language.Chinese)

	bundle.RegisterUnmarshalFunc("toml", toml.NewCoder().Unmarshal)
	bundle.RegisterUnmarshalFunc("yaml", yaml.NewCoder().Unmarshal)
	bundle.RegisterUnmarshalFunc("yml", yaml.NewCoder().Unmarshal)

	files, err := os.ReadDir(c.LanguageDir)
	if err != nil {
		logger.Fatalw("Locales init error", "err", err)
	}

	for _, file := range files {
		if !file.Type().IsRegular() {
			continue
		}
		if !strings.HasPrefix(file.Name(), "active.") {
			continue
		}
		if !strings.HasSuffix(file.Name(), "toml") && !strings.HasSuffix(file.Name(), "yaml") && !strings.HasSuffix(file.Name(), "yml") {
			continue
		}

		if _, err = bundle.LoadMessageFile(filepath.Join(c.LanguageDir, file.Name())); err != nil {
			logger.Fatalw("Locales load message file error", "err", err)
		}
	}

	return &Locales{
		cfg:       c,
		bundle:    bundle,
		languages: sync.Map{},
	}
}

func (locales *Locales) NewLocalizer(languages ...string) *i18n.Localizer {
	key := strings.Join(languages, ";")
	if val, ok := locales.languages.Load(key); ok {
		return val.(*i18n.Localizer)
	}

	loc := i18n.NewLocalizer(locales.bundle, languages...)

	locales.languages.Store(key, loc)

	return loc
}

func (locales *Locales) MustLocalize(ctx *ginserver.Context, localizeConfig *i18n.LocalizeConfig) string {
	lang := ctx.DefaultQuery("language", "zh")
	accept := ctx.GetHeader("Accept-Language")

	return locales.NewLocalizer(lang, accept).MustLocalize(localizeConfig)
}
