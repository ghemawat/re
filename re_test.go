package re_test

import (
	"fmt"
	"os"
	"re"
	"reflect"
	"regexp"
	"test"
	"testing"
)

// Table driven testing; data per case:
//	regexp
//	string
//	list of objects that are filled in
//	expected result from function
//	expected value for objects: list of <values>
//      perhaps replace expected value list with a function that checks?

var (
	hp  = regexp.MustCompile(`(\w+):(\d+)`)
	all = regexp.MustCompile(`(.*)`)
)

func TestFind(t *testing.T) {
	type args []interface{}
	b1, b2, b3 := true, true, true
	i := 0
	s := ""
	for _, c := range []struct {
		re       string
		input    string
		result   bool
		args     args
		expected []interface{}
	}{
		// Tests without any argument extraction.
		{`(\w+):(\d+)`, "", false, nil, nil},
		{`(\w+):(\d+)`, "host:1234x", true, nil, nil},
		{`(\w+):(\d+)`, "host:x1234", false, nil, nil},
		{`^(\w+):(\d+)$`, "host:1234", true, nil, nil},
		{`^(\w+):(\d+)$`, "host:1234x", false, nil, nil},
		// not enough matches
		{`^\w+:\d+$`, "host:1234", false, args{&s, &i}, nil},
		// combination of multiple arguments
		{`^(\w+):(\d+)$`, "host:5678", true, args{&s, &i}, args{"host", 5678}},
		// extraction into nil
		{`^(\w+):(\d+)$`, "host:1234", true, args{nil, nil}, args{nil, nil}},
		// boolean
		{`(\w+)`, "0", true, args{&b1}, args{false}},
		{`(\w+)`, "false", true, args{&b2}, args{false}},
		{`(\w+)`, "False", true, args{&b3}, args{false}},
		{`(\w+)`, "1", true, args{&b1}, args{true}},
		{`(\w+)`, "true", true, args{&b2}, args{true}},
		{`(\w+)`, "True", true, args{&b3}, args{true}},
		{`(\w+)`, "x0", false, args{&b1}, args{false}},
		{`(\w+)`, "xfalse", false, args{&b2}, args{false}},
		{`(\w+)`, "falsex", false, args{&b3}, args{false}},
		{`(\w+)`, "x1", false, args{&b1}, args{true}},
		{`(\w+)`, "xtrue", false, args{&b2}, args{true}},
		{`(\w+)`, "truex", false, args{&b3}, args{true}},
		// int
		{`(\d+)`, "1234", true, args{&i}, args{1234}},
		{`(.*)`, "-1234", true, args{&i}, args{-1234}},
		{`(.*)`, "123456789123456789123456789", false, args{&i}, nil},
		{`(.*)`, "-123456789123456789123456789", false, args{&i}, nil},
		// uint
		// uintptr
		// uint8
		// uint16
		// uint32
		// uint64
		// int8
		// int16
		// int32
		// byte
		// rune
		// string
		// []byte
		// func

	} {
		ok := re.Find(regexp.MustCompile(c.re), c.input, c.args...)
		if ok != c.result {
			if c.result {
				t.Errorf("Find(`%s`, `%s`, ...) failed unexpectedly", c.re, c.input)
			} else {
				t.Errorf("Find(`%s`, `%s`, ...) succeeded unexpectedly", c.re, c.input)
			}
			continue
		}
		if !ok {
			continue
		}
		for i, a := range c.args {
			if a == nil && c.expected[i] == nil {
				continue
			}
			av := reflect.Indirect(reflect.ValueOf(a)).Interface()
			fmt.Fprintln(os.Stderr, av, c.expected[i])
			if !reflect.DeepEqual(av, c.expected[i]) {
				t.Errorf("Find(`%s`, `%s`, ...): result[%d] is %#v; expected %#v\n",
					c.re, c.input, i, av, c.expected[i])
			}

		}
	}
}

func TestReStr(t *testing.T) {
	var host string
	test.Assert(t, re.Find(hp, "host:1234", &host), "re.Find host:port")
	test.Eq(t, "host", host, "did not extract correct host")
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
