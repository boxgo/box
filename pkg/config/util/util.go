package util

import (
	"strings"

	"github.com/boxgo/box/pkg/config/field"
	"github.com/olekukonko/tablewriter"
)

// SprintFields registered fields
func SprintFields(sys, user []*field.Field) (str string) {
	builder := &strings.Builder{}
	table := tablewriter.NewWriter(builder)
	table.SetHeader([]string{"Namespace", "Field", "Immutable", "Default", "Description"})
	table.SetAutoMergeCellsByColumnIndex([]int{0})
	table.SetAutoFormatHeaders(false)
	table.SetRowLine(true)

	data := make([][]string, len(sys)+len(user))
	for idx, f := range sys {
		data[idx] = f.Row()
	}

	for idx, f := range user {
		data[len(sys)+idx] = f.Row()
	}

	table.AppendBulk(data)
	table.Render()

	return builder.String()
}
