package util

import (
	"testing"
)

func TestUrlJoin(t *testing.T) {
	u, e := UrlJoin("http://example.com/", "a", "b", "c")
	AssertEqual(t, u, "http://example.com/a/b/c")
	AssertEqual(t, e, nil)

	u1, e1 := UrlJoin("http://example.com/a/", "b", "c")
	AssertEqual(t, u1, "http://example.com/a/b/c")
	AssertEqual(t, e1, nil)

	u2, e2 := UrlJoin("http://example.com/a/", "/b/", "/c/")
	AssertEqual(t, u2, "http://example.com/a/b/c")
	AssertEqual(t, e2, nil)
}

func TestUrlFormat(t *testing.T) {
	AssertEqual(t, "http://example.com/aa/b/cc/dede", UrlFormat("http://example.com/:a/b/:c/:de", map[string]interface{}{"a": "aa", "b": "bb", "c": "cc", "de": "dede"}))
	AssertEqual(t, "http://example.com/1/1.1/false", UrlFormat("http://example.com/:a/:b/:c", map[string]interface{}{"a": 1, "b": 1.1, "c": false}))
}
