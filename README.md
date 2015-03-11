# re package

Package re combines regular expression matching with fmt.Scan like
extraction of sub-matches into caller-supplied objects. Pointers to
variables can be passed as extra arguments to re.Scan.  These
variables are filled in with regular expression sub-matches.  The
sub-matches are parsed appropriately based on the type of the
variable.  E.g., if a *int is passed in, the sub-match is parsed as a
number (and overflow is detected).

For example, the host and port portions of a URL can be extracted as
follows:

	var host string
	var port int
	reg := regexp.MustCompile(`^https?://([^/:]+):(\d+)/`)
	if err := re.Scan(reg, url, &host, &port); err == nil {
		Process(host, port)
	}

A "func([]byte) error" can also be passed in as an extra argument to provide
custom parsing.

## Installation

~~~~
go get github.com/ghemawat/re
~~~~

See godoc for further documentation and examples.

 * [godoc.org/github.com/ghemawat/re](http://godoc.org/github.com/ghemawat/re)
