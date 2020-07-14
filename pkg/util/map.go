package util

import (
	"reflect"
	"sort"
)

// GetMapKeys through `reflect`
func GetMapKeys(m interface{}) (keys []string) {
	val := reflect.ValueOf(m)

	if val.Kind() != reflect.Map {
		return
	}

	for _, key := range val.MapKeys() {
		keys = append(keys, key.String())
	}

	return
}

// GetSortedMapKeys through `reflect`
func GetSortedMapKeys(m interface{}) (keys []string) {
	keys = GetMapKeys(m)
	sort.Strings(keys)

	return
}
