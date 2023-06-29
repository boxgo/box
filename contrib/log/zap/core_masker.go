package zap

import (
	"regexp"
)

type (
	Maskers []Masker
	Masker  struct {
		reg     *regexp.Regexp
		replace []byte
	}

	MaskRules []MaskRule
	MaskRule  struct {
		Rule    string
		Replace string
	}
)

var (
	DefaultMaskRules = []MaskRule{
		{Rule: `"password":(\s*)".*?"`, Replace: `"password":$1"***"`},
		{Rule: `password:(\s*).*?\S*`, Replace: `password:$1***`},
		{Rule: `password=\w*&`, Replace: `password=***&`},
		{Rule: `password=\w*\S`, Replace: `password=***`},
		{Rule: `\\"password\\":(\s*)\\".*?\\"`, Replace: `\"password\":$1\"***\"`},
	}
)

func NewMaskers(rules MaskRules) *Maskers {
	ms := make(Maskers, len(rules))

	for i, rule := range rules {
		ms[i] = Masker{
			reg:     regexp.MustCompile(rule.Rule),
			replace: []byte(rule.Replace),
		}
	}

	return &ms
}

func (m *Maskers) Mask(data []byte) []byte {
	for _, filter := range *m {
		data = filter.reg.ReplaceAll(data, filter.replace)
	}

	return data
}
