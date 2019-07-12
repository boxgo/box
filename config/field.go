package config

import (
	"reflect"
	"strings"

	"github.com/boxgo/box/minibox"
)

type (
	// Field config filed
	Field struct {
		Box      minibox.MiniBox // minibox
		Name     string          // field name. tag: config
		Type     string          // field type. reflect.Type.String
		Format   string          // field format. tag: format
		Enum     string          // field enum. tag: enum
		Default  string          // field default value. tag: default
		Required bool            // field required. tag: required
		Help     string          // field help message. tag: help
		Env      string          // field env string.
		Flag     string          // field flag string.
	}
)

const (
	configNameTag     = "config"
	configHelpTag     = "help"
	configDefaultTag  = "default"
	configFormatTag   = "format"
	configEnumTag     = "enum"
	configRequiredTag = "required"
)

func getFieldsInBoxes(cfgs ...minibox.MiniBox) []Field {
	fields := []Field{}
	for _, cfg := range cfgs {
		fields = append(fields, getFieldsInBox(cfg)...)
	}

	return fields
}

func getFieldsInBox(cfg minibox.MiniBox) []Field {
	fields := []Field{}

	value := reflect.ValueOf(cfg).Elem()
	for index := 0; index < value.NumField(); index++ {
		fieldType := value.Type().Field(index)

		if fieldType.Anonymous {
			continue
		}

		name, desc, defaultv, format, enum, required := readTags(fieldType.Tag)
		field := Field{
			Box:      cfg,
			Name:     name,
			Type:     fieldType.Type.String(),
			Format:   format,
			Enum:     enum,
			Env:      envify(cfg.Name(), name),
			Flag:     flagify(cfg.Name(), name),
			Help:     desc,
			Default:  defaultv,
			Required: required,
		}

		fields = append(fields, field)
	}

	return fields
}

// readTags get config description from tag
func readTags(tag reflect.StructTag) (name, help, def, format, enum string, required bool) {
	name = tag.Get(configNameTag)
	help = tag.Get(configHelpTag)
	def = tag.Get(configDefaultTag)
	format = tag.Get(configFormatTag)
	enum = tag.Get(configEnumTag)
	required = tag.Get(configRequiredTag) != ""

	return
}

// envify convert to env string
func envify(path, name string) string {
	words := strings.Split(path, ".")
	if name != "" {
		words = append(words, name)
	}

	str := strings.Join(words, "_")

	return strings.ToUpper(str)
}

// flagify convert to flag string
func flagify(path, name string) string {
	words := strings.Split(path, ".")
	if name != "" {
		words = append(words, name)
	}

	str := strings.Join(words, ".")

	return strings.ToLower(str)
}
