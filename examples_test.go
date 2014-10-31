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
	var host string
	var port int
	r := regexp.MustCompile(`//([^/]+):(\d+)`)
	check(re.Find(r, []byte("https://localhost:80/index.html"), &host, &port))
	fmt.Println(host, port)
	// Output:
	// localhost 80
}

func ExampleFind_skipNilMatch() {
	var port int
	r := regexp.MustCompile(`(\w+):(\d+)`)
	check(re.Find(r, []byte("localhost:80"), nil, &port))
	// Passing nil caused the first sub-match to be discarded silently.
	fmt.Println(port)
	// Output:
	// 80
}

func ExampleFind_skipTrailingMatches() {
	r := regexp.MustCompile(`(\w+):(\d+)`)
	var host string
	check(re.Find(r, []byte("localhost:80"), &host))
	// Passing fewer arguments than sub-matches caused the extra
	// sub-matches to be discarded silently.
	fmt.Println(host)
	// Output:
	// localhost
}

func ExampleFind_customType() {
	// Define a function to parse a duration.
	var interval time.Duration
	parser := func(b []byte) (err error) {
		interval, err = time.ParseDuration(string(b))
		return err
	}

	r := regexp.MustCompile(`(.*)`)
	check(re.Find(r, []byte("3m20s"), parser))
	fmt.Println(interval)
	// Output:
	// 3m20s
}

func ExampleFind_customParsing() {
	// Define a function that parses a number in binary.
	var number uint64
	parser := func(b []byte) (err error) {
		number, err = strconv.ParseUint(string(b), 2, 64)
		return err
	}

	r := regexp.MustCompile(`([01]+)`)
	check(re.Find(r, []byte("1001"), parser))
	fmt.Println(number)
	// Output:
	// 9
}
