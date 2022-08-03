package core_test

import (
	"testing"

	"github.com/boxgo/box/v2/logger/core"
)

var (
	m = core.NewMaskers(core.MaskRules{
		{`"password":(\s*)".*?"`, `"password":$1"*"`},
		{`password:(\s*).*?\S*`, `password:$1*`},
		{`password=\w*&`, `password=*&`},
		{`password=\w*\S`, `password=*`},
		{`\\"password\\":(\s*)\\".*?\\"`, `\"password\":$1\"*\"`},
	})
)

func Test_Filter1(t *testing.T) {
	expect(t, `"password":"1234" foo: 1234`, `"password":"*" foo: 1234`)
	expect(t, `"password": "1234" foo: 1234`, `"password": "*" foo: 1234`)
	expect(t, `password:1234 foo: 1234`, `password:* foo: 1234`)
	expect(t, `password: 1234 foo: 1234`, `password: * foo: 1234`)
	expect(t,
		`2019-07-25T19:54:38.160+0800	INFO	{"requestId": "04B6IyNZR", "method": "POST", "path": "/user/login", "ip": "127.0.0.1", "query": "", "body": "{\n\"userid\": \"admin\",\n\"password\": \"123123\",\n\"loginType\": \"trade\",\n\"captcha\": \"7783\"\n}"}`,
		`2019-07-25T19:54:38.160+0800	INFO	{"requestId": "04B6IyNZR", "method": "POST", "path": "/user/login", "ip": "127.0.0.1", "query": "", "body": "{\n\"userid\": \"admin\",\n\"password\": \"*\",\n\"loginType\": \"trade\",\n\"captcha\": \"7783\"\n}"}`)
}

func Test_Filter2(t *testing.T) {
	expect(t, "a=1&password=1234", "a=1&password=*")
	expect(t, "a=1&password=1234 foo", "a=1&password=* foo")
	expect(t, "a=1&password=1234&b=2 foo", "a=1&password=*&b=2 foo")
}

func expect(t *testing.T, origin, expect string) {
	str := string(m.Mask([]byte(origin)))
	if str != expect {
		t.Fatalf("\norigin: %s\nexpect: %s\nactual: %s", origin, expect, str)
	} else {
		t.Logf("\norigin: %s\nexpect: %s\nactual: %s", origin, expect, str)
	}
}
