/*
Package re combines regular expression matching with fmt.Scan like extraction
of sub-matches into caller-supplied objects.  For example, the host and port
portions of a URL can be extracted as follows:

	var host string
	var port int
	reg := regexp.MustCompile(`^https?://([^/:]+):(\d+)/`)
	if err := re.Scan(reg, url, &host, &port); err == nil {
		Process(host, port)
	}
*/
package re

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
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
		return fmt.Errorf(`re.Scan: could not find "%s" in "%s"`,
			re, input)
	}
	if len(matches) < 2+2*len(output) {
		return fmt.Errorf(`re.Scan: only got %d matches from "%s"; need at least %d`,
			len(matches)/2-1, re, len(output))
	}
	for i, r := range output {
		start, limit := matches[2+2*i], matches[2+2*i+1]
		if start < 0 || limit < 0 {
			// Sub-expression is missing; treat as empty.
			start = 0
			limit = 0
		}
		if err := assign(r, input[start:limit]); err != nil {
			return err
		}
	}
	return nil
}

func assign(r interface{}, b []byte) error {
	switch v := r.(type) {
	case nil:
		// Discard the match.
	case func([]byte) error:
		if err := v(b); err != nil {
			return err
		}
	case *string:
		*v = string(b)
	case *[]byte:
		*v = b
	case *int:
		if i, err := strconv.ParseInt(string(b), 0, 64); err != nil {
			return err
		} else {
			if int64(int(i)) != i {
				return parseError("out of range for int", b)
			}
			*v = int(i)
		}
	case *int8:
		if i, err := strconv.ParseInt(string(b), 0, 8); err != nil {
			return err
		} else {
			*v = int8(i)
		}
	case *int16:
		if i, err := strconv.ParseInt(string(b), 0, 16); err != nil {
			return err
		} else {
			*v = int16(i)
		}
	case *int32:
		if i, err := strconv.ParseInt(string(b), 0, 32); err != nil {
			return err
		} else {
			*v = int32(i)
		}
	case *int64:
		if i, err := strconv.ParseInt(string(b), 0, 64); err != nil {
			return err
		} else {
			*v = i
		}
	case *uint:
		if u, err := strconv.ParseUint(string(b), 0, 64); err != nil {
			return err
		} else {
			if uint64(uint(u)) != u {
				return parseError("out of range for uint", b)
			}
			*v = uint(u)
		}
	case *uintptr:
		if u, err := strconv.ParseUint(string(b), 0, 64); err != nil {
			return err
		} else {
			if uint64(uintptr(u)) != u {
				return parseError("out of range for uintptr", b)
			}
			*v = uintptr(u)
		}
	case *uint8:
		if u, err := strconv.ParseUint(string(b), 0, 8); err != nil {
			return err
		} else {
			*v = uint8(u)
		}
	case *uint16:
		if u, err := strconv.ParseUint(string(b), 0, 16); err != nil {
			return err
		} else {
			*v = uint16(u)
		}
	case *uint32:
		if u, err := strconv.ParseUint(string(b), 0, 32); err != nil {
			return err
		} else {
			*v = uint32(u)
		}
	case *uint64:
		if u, err := strconv.ParseUint(string(b), 0, 64); err != nil {
			return err
		} else {
			*v = u
		}
	case *float32:
		if f, err := strconv.ParseFloat(string(b), 32); err != nil {
			return err
		} else {
			*v = float32(f)
		}
	case *float64:
		if f, err := strconv.ParseFloat(string(b), 64); err != nil {
			return err
		} else {
			*v = f
		}
	default:
		t := reflect.ValueOf(r).Type()
		return parseError(fmt.Sprintf("unsupported type %s", t), b)
	}
	return nil
}

func parseError(explanation string, b []byte) error {
	return fmt.Errorf(`re.Scan: parsing "%s": %s`, b, explanation)
}
