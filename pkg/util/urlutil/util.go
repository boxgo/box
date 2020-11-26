package urlutil

import (
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strings"
)

func UrlJoin(baseUrl string, segments ...string) (string, error) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return "", err
	}

	u.Path = path.Join(append([]string{u.Path}, segments...)...)

	return u.String(), nil
}

func UrlFormat(rawUrl string, pathParam map[string]interface{}) string {
	re := regexp.MustCompile(`/:(\w+)`)

	return re.ReplaceAllStringFunc(rawUrl, func(s string) string {
		key := strings.Replace(s, "/:", "", -1)
		if val, ok := pathParam[key]; ok {
			return fmt.Sprintf("/%v", val)
		}

		return s
	})
}
