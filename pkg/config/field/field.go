package field

import (
	"fmt"
	"strings"
)

type (
	Field struct {
		Path    string      `json:"path"`
		Name    string      `json:"name"`
		Desc    string      `json:"desc"`
		Type    string      `json:"type"`
		Default interface{} `json:"default"`
	}
)

func (f *Field) Row() []string {
	return []string{f.Path, f.Name, f.Type, fmt.Sprintf("%v", f.Default), f.Desc}
}

func (f *Field) String() string {
	if f.Name != "" {
		return f.Name + "." + f.Path
	}

	return f.Path
}

func (f *Field) Paths() []string {
	return strings.Split(f.String(), ".")
}
