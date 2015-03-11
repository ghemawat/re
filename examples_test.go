package re_test

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/ghemawat/re"
)

// Parse a line of ls -l output into its fields.
func ExampleScan() {
	var f struct {
		mode, user, group, date, name string
		nlinks, size                  int64
	}

	// Sample output from `ls -l --time-style=iso`
	line := "-rwxr-xr-x 1 root root 110080 2014-03-24  /bin/ls"

	// A regexp that matches such lines.
	r := regexp.MustCompile(`^(.{10}) +(\d+) +(\w+) +(\w+) +(\d+) +(\S+) +(.+)$`)

	// Match line to regexp and extract properties into struct.
	if err := re.Scan(r, []byte(line), &f.mode, &f.nlinks, &f.user, &f.group, &f.size, &f.date, &f.name); err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", f)
	// Output:
	// {mode:-rwxr-xr-x user:root group:root date:2014-03-24 name:/bin/ls nlinks:1 size:110080}
}

// Use a custom parsing function that parses a number in binary.
func ExampleScan_binaryNumber() {
	var number uint64
	parseBinary := func(b []byte) (err error) {
		number, err = strconv.ParseUint(string(b), 2, 64)
		return err
	}

	r := regexp.MustCompile(`([01]+)`)
	if err := re.Scan(r, []byte("1001"), parseBinary); err != nil {
		panic(err)
	}
	fmt.Println(number)
	// Output:
	// 9
}

// Use a custom re-usable parser for time.Duration.
func ExampleScan_parseDuration() {
	// parseDuration(&d) returns a parser that stores its result in *d.
	parseDuration := func(d *time.Duration) func([]byte) error {
		return func(b []byte) (err error) {
			*d, err = time.ParseDuration(string(b))
			return err
		}
	}

	r := regexp.MustCompile(`^elapsed: (.*)$`)
	var interval time.Duration
	if err := re.Scan(r, []byte("elapsed: 200s"), parseDuration(&interval)); err != nil {
		panic(err)
	}
	fmt.Println(interval)
	// Output:
	// 3m20s
}
