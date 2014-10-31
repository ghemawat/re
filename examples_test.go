package re_test

import (
	"fmt"
	"re"
	"regexp"
	"strconv"
	"time"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func ExampleFind() {
	// A regexp that matches a line of simplified ls -l output.
	r := regexp.MustCompile(`^(.{10}) +(\d+) +(\w+) +(\w+) +(\d+) +(\S+) +(\S+)`)
	var s struct {
		mode, user, group, date, name string
		nlinks, size                  int64
	}
	err := re.Find(r, []byte("-rwxr-xr-x 1 root root 110080 2014-03-24  /bin/ls"),
		&s.mode, &s.nlinks, &s.user, &s.group, &s.size, &s.date, &s.name)
	check(err)
	fmt.Printf("%+v\n", s)
	// Output:
	// {mode:-rwxr-xr-x user:root group:root date:2014-03-24 name:/bin/ls nlinks:1 size:110080}
}

func ExampleFind_customParsing() {
	// Define a function that parses a number in binary.
	var number uint64
	parseBinary := func(b []byte) (err error) {
		number, err = strconv.ParseUint(string(b), 2, 64)
		return err
	}

	r := regexp.MustCompile(`([01]+)`)
	err := re.Find(r, []byte("1001"), parseBinary)
	check(err)
	fmt.Println(number)
	// Output:
	// 9
}

func ExampleFind_supportNewType() {
	// A function that returns a custom parser that parses into
	// the specified *time.Duration.
	parseDuration := func(d *time.Duration) func([]byte) error {
		return func(b []byte) (err error) {
			*d, err = time.ParseDuration(string(b))
			return err
		}
	}

	r := regexp.MustCompile(`(.*)`)
	var interval time.Duration
	err := re.Find(r, []byte("3m20s"), parseDuration(&interval))
	check(err)
	fmt.Println(interval)
	// Output:
	// 3m20s
}
