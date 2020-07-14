package config

import (
	"fmt"
	"strings"

	"github.com/boxgo/box/pkg/util"
	"github.com/olekukonko/tablewriter"
)

// sprintFields registered fields
func sprintFields(fields map[string]*Field) (str string) {
	builder := &strings.Builder{}
	table := tablewriter.NewWriter(builder)
	table.SetHeader([]string{"Namespace", "Field", "Default", "Description"})
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)

	data := make([][]string, len(fields))
	for idx, key := range util.GetSortedMapKeys(fields) {
		field := fields[key]
		data[idx] = []string{field.name, field.path, fmt.Sprintf("%v", field.def), field.desc}
	}

	table.AppendBulk(data)
	table.Render()

	return builder.String()
}

// sprintTemplate through encoder
func sprintTemplate(fields map[string]*Field, encoder string) (str string) {
	lastParent := ""
	for _, key := range util.GetSortedMapKeys(fields) {
		isNewParent := true
		field := fields[key]
		words := strings.Split(field.String(), ".")

		if words[0] == lastParent {
			isNewParent = false
		}

		lastParent = words[0]
		for i, word := range words {
			if i == len(words)-1 { // final field
				str += fmt.Sprintf("%s// %s \n", strings.Repeat("\t", i), field.desc)
				str += fmt.Sprintf("%s%s: %v \n", strings.Repeat("\t", i), word, field.def)
			} else if isNewParent { // parent
				str += fmt.Sprintf("%s%s:\n", strings.Repeat("\t", i), word)
			}
		}
	}

	return
}

// field2path convert field to config path
func field2path(field *Field) []string {
	return strings.Split(field.String(), ".")
}
