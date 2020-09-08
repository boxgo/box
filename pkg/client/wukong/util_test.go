package wukong

import (
	"testing"
)

func TestUrlJoin(t *testing.T) {
	u, e := urlJoin("http://example.com/", "a", "b", "c")
	assertEqual(t, u, "http://example.com/a/b/c")
	assertEqual(t, e, nil)

	u1, e1 := urlJoin("http://example.com/a/", "b", "c")
	assertEqual(t, u1, "http://example.com/a/b/c")
	assertEqual(t, e1, nil)

	u2, e2 := urlJoin("http://example.com/a/", "/b/", "/c/")
	assertEqual(t, u2, "http://example.com/a/b/c")
	assertEqual(t, e2, nil)
}

func TestUrlFormat(t *testing.T) {
	assertEqual(t, "http://example.com/aa/b/cc/dede", urlFormat("http://example.com/:a/b/:c/:de", map[string]interface{}{"a": "aa", "b": "bb", "c": "cc", "de": "dede"}))
	assertEqual(t, "http://example.com/1/1.1/false", urlFormat("http://example.com/:a/:b/:c", map[string]interface{}{"a": 1, "b": 1.1, "c": false}))
}
