package urlutil

import (
	"testing"

	"github.com/boxgo/box/v2/util/testutil"
)

func TestUrlJoin(t *testing.T) {
	u, e := UrlJoin("http://example.com/", "a", "b", "c")
	testutil.ExpectEqual(t, u, "http://example.com/a/b/c")
	testutil.ExpectEqual(t, e, nil)

	u1, e1 := UrlJoin("http://example.com/a/", "b", "c")
	testutil.ExpectEqual(t, u1, "http://example.com/a/b/c")
	testutil.ExpectEqual(t, e1, nil)

	u2, e2 := UrlJoin("http://example.com/a/", "/b/", "/c/")
	testutil.ExpectEqual(t, u2, "http://example.com/a/b/c")
	testutil.ExpectEqual(t, e2, nil)
}

func TestUrlFormat(t *testing.T) {
	testutil.ExpectEqual(t, "http://example.com/aa/b/cc/dede", UrlFormat("http://example.com/:a/b/:c/:de", map[string]interface{}{"a": "aa", "b": "bb", "c": "cc", "de": "dede"}))
	testutil.ExpectEqual(t, "http://example.com/1/1.1/false", UrlFormat("http://example.com/:a/:b/:c", map[string]interface{}{"a": 1, "b": 1.1, "c": false}))
}
