package printer

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"

	"github.com/boxgo/box/minibox"
)

type (
	showConfig struct {
		name   string
		scale  float64
		width  int
		header string
	}
)

var (
	line           = ""
	extLine        = ""
	showTotalWidth = 150.0
	showConfigs    = []*showConfig{
		{name: "Prefix", scale: 0.2},
		{name: "Field", scale: 0.15},
		{name: "Env", scale: 0.25},
		{name: "Desc", scale: 0.4},
	}
)

var (
	configFieldValueTags = []string{"json", "yaml", "xml", "cfg", "config"}
	configFieldDescTags  = "desc"
)

func init() {
	realTotalWidth := 0

	for _, it := range showConfigs {
		nameWidth := len(it.name)
		fieldWidth := int(math.Floor(showTotalWidth*it.scale)) - nameWidth
		if fieldWidth%2 != 0 {
			fieldWidth++
		}
		blank := strings.Repeat(" ", fieldWidth/2)

		it.header = blank + it.name + blank
		it.width = len(it.header)
		realTotalWidth += it.width
	}

	line = strings.Repeat("-", realTotalWidth)
	extLine = strings.Repeat("-    -", realTotalWidth/6)
}

// PrintPrettyConfigMap Print a pretty config detail map
func PrintPrettyConfigMap(configs []minibox.MiniBox) {
	fmt.Println(line)
	fmt.Println(strings.Join([]string{
		showConfigs[0].header,
		showConfigs[1].header,
		showConfigs[2].header,
		showConfigs[3].header,
	}, "|"))
	fmt.Println(line)

	for _, cfg := range configs {
		cfgType := reflect.TypeOf(cfg).Elem()
		read(cfg.Name(), 0, cfgType)

		if boxExt, ok := cfg.(minibox.MiniBoxExt); ok {
			for _, ext := range boxExt.Exts() {
				fmt.Println(extLine)
				extType := reflect.TypeOf(ext).Elem()
				read(ext.Name(), 0, extType)
			}
		}

		fmt.Println(line)
	}
}

func read(prefix string, showedIndex int, cfgType reflect.Type) {
	for index := 0; index < cfgType.NumField(); index++ {
		structField := cfgType.Field(index)

		if structField.Anonymous {
			switch structField.Type.Kind() {
			case reflect.Struct:
				read(prefix, showedIndex, structField.Type)
			case reflect.Ptr:
				read(prefix, showedIndex, structField.Type.Elem())
			}

			return
		}

		if name, desc := readField(structField.Tag); name != "" {
			env := ""
			if prefix == "" {
				env = name
			} else {
				env = fmt.Sprintf("%s_%s", strings.Join(strings.Split(prefix, "."), "_"), name)
			}
			env = strings.ToUpper(env)

			format := strings.Join([]string{
				"",
				"%-" + strconv.Itoa(showConfigs[0].width-2) + "s |",
				"%-" + strconv.Itoa(showConfigs[1].width-2) + "s |",
				"%-" + strconv.Itoa(showConfigs[2].width-2) + "s |",
				"%-" + strconv.Itoa(showConfigs[3].width) + "s",
			}, " ") + "\n"

			if showedIndex == 0 {
				fmt.Printf(format, prefix, name, env, desc)
			} else {
				fmt.Printf(format, "", name, env, desc)
			}

			showedIndex++
		}
	}
}

func readField(tag reflect.StructTag) (name, desc string) {
	for _, tagFlag := range configFieldValueTags {
		if cfgFieldName, ok := tag.Lookup(tagFlag); ok && cfgFieldName != "-" {
			name = cfgFieldName
			desc = tag.Get(configFieldDescTags)
		}
	}

	return
}
