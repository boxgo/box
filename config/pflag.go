package config

import (
	"strconv"
	"strings"
	"time"

	"github.com/spf13/pflag"
)

// bind pflag by field
func bindFieldPflag(field Field) {
	if pflag.CommandLine.Lookup(field.Flag) != nil {
		return
	}

	switch field.Type {
	case "string":
		pflag.String(field.Flag, field.Default, field.Help)
	case "bool":
		b, _ := strconv.ParseBool(field.Default)
		pflag.Bool(field.Flag, b, field.Help)
	case "int":
		i, _ := strconv.ParseInt(field.Default, 10, 64)
		pflag.Int(field.Flag, int(i), field.Help)
	case "uint":
		i, _ := strconv.ParseUint(field.Default, 10, 64)
		pflag.Uint(field.Flag, uint(i), field.Help)
	case "float64":
		f, _ := strconv.ParseFloat(field.Default, 64)
		pflag.Float64(field.Flag, f, field.Help)
	case "float32":
		f, _ := strconv.ParseFloat(field.Default, 32)
		pflag.Float32(field.Flag, float32(f), field.Help)
	case "time.Duration":
		d, _ := time.ParseDuration(field.Default)
		pflag.Duration(field.Flag, d, field.Help)
	case "[]string":
		pflag.StringSlice(field.Flag, strings.Split(field.Default, ","), field.Help)
	case "[]uint":
		pflag.UintSlice(field.Flag, []uint{}, field.Help)
	}
}
