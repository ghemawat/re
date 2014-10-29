package re_test

import (
	"re"
	"regexp"
	"test"
	"testing"
)

var (
	hp  = regexp.MustCompile(`(\w+):(\d+)`)
	all = regexp.MustCompile(`(.*)`)
)

func TestReExtract(t *testing.T) {
	var host string
	var port int
	test.Assert(t, re.Find(hp, "host:1234", &host, &port), "re.Find host:port")
	test.Eq(t, "host", host, "did not extract correct host")
	test.Eq(t, 1234, port, "did not extract correct port")
}

func TestReBoolOk(t *testing.T) {
	var tests = []struct {
		arg    string
		result bool
	}{
		{"0", false},
		{"false", false},
		{"False", false},
		{"1", true},
		{"true", true},
		{"true", true},
	}
	for _, c := range tests {
		var b bool
		test.Assert(t, re.Find(all, c.arg, &b), "could not parse", c.arg)
		test.Eq(t, c.result, b, "unexpected boolean result")
	}
}

func TestReBoolFail(t *testing.T) {
	for _, arg := range []string{
		"0x", "x0", "1x", "x1", "xtrue", "truex",
		"xfalse", "falsex", "tru", "fals",
	} {
		var b bool
		test.Eq(t, false, re.Find(all, arg, &b), "unexpectedly parsed", arg)
	}
}

func TestReFunc(t *testing.T) {
	var arg string
	f := func(a string) bool {
		arg = a
		return true
	}
	test.Assert(t, re.Find(hp, "host:1234", f), "re.Find function")
	test.Eq(t, "host", arg, "wrong argument to function called by Find")
	test.Assert(t, re.Find(hp, "host:1234", nil, f), "re.Find function")
	test.Eq(t, "1234", arg, "wrong argument to function called by Find")
}

func TestReFuncFailure(t *testing.T) {
	var arg string
	f := func(a string) bool {
		arg = a
		return false
	}
	test.Eq(t, false, re.Find(hp, "host:1234", f), "re.Find function")
	test.Eq(t, "host", arg, "wrong argument to function called by Find")
}
