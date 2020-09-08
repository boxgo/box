package wukong

import (
	"fmt"
	"net/url"
	"path"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		return
	}

	t.Errorf("Received %v (type %v), expected %v (type %v)", a, reflect.TypeOf(a), b, reflect.TypeOf(b))
}

func urlJoin(baseUrl string, segments ...string) (string, error) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return "", err
	}

	u.Path = path.Join(append([]string{u.Path}, segments...)...)

	return u.String(), nil
}

func urlFormat(rawUrl string, pathParam map[string]interface{}) string {
	re := regexp.MustCompile(`/:(\w+)`)

	return re.ReplaceAllStringFunc(rawUrl, func(s string) string {
		key := strings.Replace(s, "/:", "", -1)
		if val, ok := pathParam[key]; ok {
			return fmt.Sprintf("/%v", val)
		}

		return s
	})
}
