package core

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
		Rule    string `config:"rule" json:"rule"`
		Replace string `config:"replace" json:"replace"`
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
	if m == nil || len(*m) == 0 {
		return data
	}

	for _, filter := range *m {
		data = filter.reg.ReplaceAll(data, filter.replace)
	}

	return data
}
