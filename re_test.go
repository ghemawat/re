package re_test

import (
	"re"
	"reflect"
	"regexp"
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

type testcase struct {
	re       string
	input    string
	result   bool
	args     []interface{}
	expected []interface{}
}

func c(re, input string, result bool, argexpect ...interface{}) testcase {
	t := testcase{re, input, result, nil, nil}
	for i := 0; i < len(argexpect); i += 2 {
		t.args = append(t.args, argexpect[i])
		t.expected = append(t.expected, argexpect[i+1])
	}
	return t
}

func newtrue() *bool {
	r := new(bool)
	*r = true
	return r
}

func TestFind(t *testing.T) {
	for _, c := range []testcase{
		// Tests without any argument extraction.
		c(`(\w+):(\d+)`, "", false),
		c(`(\w+):(\d+)`, "host:1234x", true),
		c(`(\w+):(\d+)`, "host:x1234", false),
		c(`^(\w+):(\d+)$`, "host:1234", true, nil, nil),
		c(`^(\w+):(\d+)$`, "host:1234x", false, nil, nil),

		// not enough matches
		c(`^\w+:\d+$`, "host:1234", false, new(string), nil, new(int), nil),

		// extraction into nil
		c(`^(\w+):(\d+)$`, "host:1234", true, nil, nil, nil, nil),

		// missing sub-expression
		c(`^(\w+):((\d+))?`, "host:", true, nil, nil, nil, nil, nil, nil),
		c(`^(\w+):((\d+))?`, "host:", false, nil, nil, new(int), nil, nil, nil),

		// string
		c(`(.*):\d+`, "host:1234", true, new(string), "host"),

		// []byte
		c(`(.*):\d+`, "host:1234", true, new([]byte), []byte("host")),
		c(`(.*):\d+`, ":1234", true, new([]byte), []byte("")),

		// boolean
		c(`(\w+)`, "0", true, newtrue(), false),
		c(`(\w+)`, "false", true, newtrue(), false),
		c(`(\w+)`, "False", true, newtrue(), false),
		c(`(\w+)`, "1", true, new(bool), true),
		c(`(\w+)`, "true", true, new(bool), true),
		c(`(\w+)`, "True", true, new(bool), true),
		c(`(\w+)`, "x0", false, new(bool), false),
		c(`(\w+)`, "xfalse", false, new(bool), false),
		c(`(\w+)`, "falsex", false, new(bool), false),
		c(`(\w+)`, "x1", false, new(bool), true),
		c(`(\w+)`, "xtrue", false, new(bool), true),
		c(`(\w+)`, "truex", false, new(bool), true),

		// int
		c(`(\d+)`, "1234", true, new(int), 1234),
		c(`(.*)`, "-1234", true, new(int), -1234),
		c(`(.*)`, "123456789123456789123456789", false, new(int), nil),
		c(`(.*)`, "-123456789123456789123456789", false, new(int), nil),
		c(`(.*)`, "0x10", true, new(int), 0x10),
		c(`(.*)`, "010", true, new(int), 010),

		// uint
		c(`(\d+)`, "1234", true, new(uint), uint(1234)),
		c(`(\d+)`, "123456789123456789123456789", false, new(uint), nil),

		// uintptr
		c(`(\d+)`, "1234", true, new(uintptr), uintptr(1234)),
		c(`(\d+)`, "123456789123456789123456789", false, new(uintptr), nil),

		// uint8
		c(`(.*)`, "0", true, new(uint8), uint8(0)),
		c(`(.*)`, "17", true, new(uint8), uint8(17)),
		c(`(.*)`, "255", true, new(uint8), uint8(255)),
		c(`(.*)`, "256", false, new(uint8), nil),
		c(`(.*)`, "x", false, new(uint8), nil),

		// uint16
		c(`(.*)`, "0", true, new(uint16), uint16(0)),
		c(`(.*)`, "17", true, new(uint16), uint16(17)),
		c(`(.*)`, "65535", true, new(uint16), uint16(65535)),
		c(`(.*)`, "65536", false, new(uint16), nil),
		c(`(.*)`, "x", false, new(uint16), nil),

		// uint32
		c(`(.*)`, "0", true, new(uint32), uint32(0)),
		c(`(.*)`, "17", true, new(uint32), uint32(17)),
		c(`(.*)`, "4294967295", true, new(uint32), uint32(4294967295)),
		c(`(.*)`, "4294967296", false, new(uint32), nil),
		c(`(.*)`, "x", false, new(uint32), nil),

		// uint64
		c(`(.*)`, "0", true, new(uint64), uint64(0)),
		c(`(.*)`, "17", true, new(uint64), uint64(17)),
		c(`(.*)`, "18446744073709551615", true, new(uint64), uint64(18446744073709551615)),
		c(`(.*)`, "18446744073709551616", false, new(uint64), nil),
		c(`(.*)`, "x", false, new(uint64), nil),

		// int8
		c(`(.*)`, "0", true, new(int8), int8(0)),
		c(`(.*)`, "17", true, new(int8), int8(17)),
		c(`(.*)`, "127", true, new(int8), int8(127)),
		c(`(.*)`, "128", false, new(int8), nil),
		c(`(.*)`, "x", false, new(int8), nil),

		// int16
		c(`(.*)`, "0", true, new(int16), int16(0)),
		c(`(.*)`, "17", true, new(int16), int16(17)),
		c(`(.*)`, "32767", true, new(int16), int16(32767)),
		c(`(.*)`, "32768", false, new(int16), nil),
		c(`(.*)`, "x", false, new(int16), nil),

		// int32
		c(`(.*)`, "0", true, new(int32), int32(0)),
		c(`(.*)`, "17", true, new(int32), int32(17)),
		c(`(.*)`, "2147483647", true, new(int32), int32(2147483647)),
		c(`(.*)`, "2147483648", false, new(int32), nil),
		c(`(.*)`, "x", false, new(int32), nil),

		// combination of multiple arguments
		c(`^(\w+):(\d+)$`, "host:5678", true, new(string), "host", new(int), 5678),
	} {
		ok := re.Find(regexp.MustCompile(c.re), []byte(c.input), c.args...)
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
			if !reflect.DeepEqual(av, c.expected[i]) {
				t.Errorf("Find(`%s`, `%s`, ...): result[%d] is %#v; expected %#v\n",
					c.re, c.input, i, av, c.expected[i])
			}

		}
	}
}

func TestReFunc(t *testing.T) {
	var arg string
	savearg := func(a []byte) bool {
		arg = string(a)
		return true
	}
	hp := `^(\w+):(\d+)$`
	str := "host:1234"
	if !re.Find(regexp.MustCompile(hp), []byte(str), savearg) {
		t.Fatalf("Find(`%s`, `%s`, savearg): failed unexpectedly", hp, str)
	}
	if arg != "host" {
		t.Fatalf("Find(`%s`, `%s`, savearg): did not call function", hp, str)
	}

	fail := func(a []byte) bool {
		arg = string(a)
		return false
	}
	if re.Find(regexp.MustCompile(hp), []byte(str), fail) {
		t.Fatalf("Find(`%s`, `%s`, fail): succeeded unexpectedly", hp, str)
	}
}
