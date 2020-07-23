package config

import (
	"fmt"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// sprintFields registered fields
func sprintFields(sys, user []*Field) (str string) {
	builder := &strings.Builder{}
	table := tablewriter.NewWriter(builder)
	table.SetHeader([]string{"Namespace", "Field", "Default", "Description"})
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)

	data := make([][]string, len(sys)+len(user))
	for idx, field := range sys {
		data[idx] = []string{field.name, field.path, fmt.Sprintf("%v", field.def), field.desc}
	}

	for idx, field := range user {
		data[len(sys)+idx] = []string{field.name, field.path, fmt.Sprintf("%v", field.def), field.desc}
	}

	table.AppendBulk(data)
	table.Render()

	return builder.String()
}

// sprintTemplate through encoder
func sprintTemplate(sys, user []*Field, encoder string) (tmpl string) {
	lastParent := ""

	sprint := func(f *Field) (str string) {
		isNewParent := true
		words := strings.Split(f.String(), ".")

		if words[0] == lastParent {
			isNewParent = false
		}

		lastParent = words[0]
		for i, word := range words {
			if i == len(words)-1 { // final field
				str += fmt.Sprintf("%s// %s \n", strings.Repeat("\t", i), f.desc)
				str += fmt.Sprintf("%s%s: %v \n", strings.Repeat("\t", i), word, f.def)
			} else if isNewParent { // parent
				str += fmt.Sprintf("%s%s:\n", strings.Repeat("\t", i), word)
			}
		}

		return
	}

	for _, field := range sys {
		tmpl += sprint(field)
	}
	for _, field := range user {
		tmpl += sprint(field)
	}

	return
}

// field2path convert field to config path
func field2path(field *Field) []string {
	return strings.Split(field.String(), ".")
}
