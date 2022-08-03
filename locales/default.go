package locales

import (
	"github.com/boxgo/box/v2/server/ginserver"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	Default = StdConfig("default").Build()
)

func NewLocalizer(languages ...string) *i18n.Localizer {
	return Default.NewLocalizer(languages...)
}

func MustLocalize(ctx *ginserver.Context, localizeConfig *i18n.LocalizeConfig) string {
	return Default.MustLocalize(ctx, localizeConfig)
}
