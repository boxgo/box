package zap_test

import (
	"testing"
)

func Benchmark_Filter1(b *testing.B) {
	expectB(b, `"password":"1234" foo: 1234`, `"password":"*" foo: 1234`)
	expectB(b, `"password": "1234" foo: 1234`, `"password": "*" foo: 1234`)
	expectB(b, `password:1234 foo: 1234`, `password:* foo: 1234`)
	expectB(b, `password: 1234 foo: 1234`, `password: * foo: 1234`)
	expectB(b,
		`2019-07-25T19:54:38.160+0800	INFO	{"requestId": "04B6IyNZR", "method": "POST", "path": "/user/login", "ip": "127.0.0.1", "query": "", "body": "{\n\"userid\": \"admin\",\n\"password\": \"123123\",\n\"loginType\": \"trade\",\n\"captcha\": \"7783\"\n}"}`,
		`2019-07-25T19:54:38.160+0800	INFO	{"requestId": "04B6IyNZR", "method": "POST", "path": "/user/login", "ip": "127.0.0.1", "query": "", "body": "{\n\"userid\": \"admin\",\n\"password\": \"*\",\n\"loginType\": \"trade\",\n\"captcha\": \"7783\"\n}"}`)
}

func Benchmark_Filter2(b *testing.B) {
	expectB(b, "a=1&password=1234", "a=1&password=*")
	expectB(b, "a=1&password=1234 foo", "a=1&password=* foo")
	expectB(b, "a=1&password=1234&b=2 foo", "a=1&password=*&b=2 foo")
}

func expectB(b *testing.B, origin, expect string) {
	str := string(masker.Mask([]byte(origin)))
	if str != expect {
		b.Fatalf("\norigin: %s\nexpect: %s\nactual: %s", origin, expect, str)
	}
}
