/*
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
*/
package re

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
)

// Span is a special type designed to be passed via pointer to Scan.  re.Scan
// will store the starting and ending offsets of the corresponding regular
// expression capture group into the Span.
//
// This type can be placed anywhere within the list of arguments to scan, but
// the most typical usage is to find the entire extent of the match, which can
// be achieved by placing it third (immediately after the regular expression
// and input), and wrapping the entire regexp in parentheses so that the Span
// is filled with the extent of the entire match.
type Span struct {
	Start int
	End   int
}

var (
	NotFound = errors.New("not found")
)

// Scan returns nil if regular expression re matches somewhere in
// input, and for every non-nil entry in output, the corresponding
// regular expression sub-match is succesfully parsed and stored into
// *output[i].
//
// The following can be passed as output arguments to Scan:
//
// nil: The corresponding sub-match is discarded without being saved.
//
// Pointer to string or []byte: The corresponding sub-match is
// stored in the pointed-to object.  When storing into a []byte, no
// copying is done, and the stored slice is an alias of the input.
//
// Pointer to some built-in numeric types (int, int8, int16, int32,
// int64, uint, uintptr, uint8, uint16, uint32, uint64, float32,
// float64): The corresponding sub-match will be parsed as a literal
// of the numeric type and the result stored into *output[i].  Scan
// will return an error if the sub-match cannot be parsed
// successfully, or the parse result is out of range for the type.
//
// Pointer to a rune or a byte: rune is an alias of uint32 and byte is
// an alias of uint8, so the preceding rule applies; i.e., Scan treats
// the input as a string of digits to be parsed into the rune or
// byte. Therefore Scan cannot be used to directly extract a single
// rune or byte from the input. For that, parse into a string or
// []byte and use the first element, or pass in a custom parsing
// function (see below).
//
// func([]byte) error: The function is passed the corresponding
// sub-match.  If the result is a non-nil error, the Scan call fails
// with that error. Pass in such a function to provide custom parsing:
// e.g., treating a number as decimal even if it starts with "0"
// (normally Scan would treat such as a number as octal); or parsing
// an otherwise unsupported type like time.Duration.
//
// An error is returned if output[i] does not have one of the preceding
// types.  Caveat: the set of supported types might be extended in the
// future.
//
// Extra sub-matches (ones with no corresponding output) are discarded silently.
func Scan(re *regexp.Regexp, input []byte, output ...interface{}) error {
	matches := re.FindSubmatchIndex(input)
	if matches == nil {
		return fmt.Errorf("regular expression %q: %w", re, NotFound)
	}
	if len(matches) < 2+2*len(output) {
		return fmt.Errorf(`re.Scan: only got %d matches from "%s"; need at least %d`,
			len(matches)/2-1, re, len(output))
	}
	for i, r := range output {
		span := Span{
			Start: matches[2+2*i],
			End:   matches[2+2*i+1],
		}
		var submatch []byte
		if span.Start > -1 && span.End >= span.Start {
			submatch = input[span.Start:span.End]
		}
		if err := assign(r, submatch, span); err != nil {
			return err
		}
	}
	return nil
}

// ScanString behaves the same as Scan, but it matches the regexp against a
// string, rather than a byte array.
func ScanString(re *regexp.Regexp, input string, output ...interface{}) error {
	return Scan(re, []byte(input), output...)
}

func assign(r interface{}, b []byte, s Span) error {
	switch v := r.(type) {
	case nil:
		// Discard the match.
	case func([]byte) error:
		if err := v(b); err != nil {
			return err
		}
	case *Span:
		*v = s
	case *string:
		*v = string(b)
	case *[]byte:
		*v = b
	case *int:
		i, err := strconv.ParseInt(string(b), 0, 64)
		if err != nil {
			return err
		}
		if int64(int(i)) != i {
			return parseError("out of range for int", b)
		}
		*v = int(i)
	case *int8:
		i, err := strconv.ParseInt(string(b), 0, 8)
		if err != nil {
			return err
		}
		*v = int8(i)
	case *int16:
		i, err := strconv.ParseInt(string(b), 0, 16)
		if err != nil {
			return err
		}
		*v = int16(i)
	case *int32:
		i, err := strconv.ParseInt(string(b), 0, 32)
		if err != nil {
			return err
		}
		*v = int32(i)
	case *int64:
		i, err := strconv.ParseInt(string(b), 0, 64)
		if err != nil {
			return err
		}
		*v = i
	case *uint:
		u, err := strconv.ParseUint(string(b), 0, 64)
		if err != nil {
			return err
		}
		if uint64(uint(u)) != u {
			return parseError("out of range for uint", b)
		}
		*v = uint(u)
	case *uintptr:
		u, err := strconv.ParseUint(string(b), 0, 64)
		if err != nil {
			return err
		}
		if uint64(uintptr(u)) != u {
			return parseError("out of range for uintptr", b)
		}
		*v = uintptr(u)
	case *uint8:
		u, err := strconv.ParseUint(string(b), 0, 8)
		if err != nil {
			return err
		}
		*v = uint8(u)
	case *uint16:
		u, err := strconv.ParseUint(string(b), 0, 16)
		if err != nil {
			return err
		}
		*v = uint16(u)
	case *uint32:
		u, err := strconv.ParseUint(string(b), 0, 32)
		if err != nil {
			return err
		}
		*v = uint32(u)
	case *uint64:
		u, err := strconv.ParseUint(string(b), 0, 64)
		if err != nil {
			return err
		}
		*v = u
	case *float32:
		f, err := strconv.ParseFloat(string(b), 32)
		if err != nil {
			return err
		}
		*v = float32(f)
	case *float64:
		f, err := strconv.ParseFloat(string(b), 64)
		if err != nil {
			return err
		}
		*v = f
	default:
		t := reflect.ValueOf(r).Type()
		return parseError(fmt.Sprintf("unsupported type %s", t), b)
	}
	return nil
}

func parseError(explanation string, b []byte) error {
	return fmt.Errorf(`re.Scan: parsing "%s": %s`, b, explanation)
}
