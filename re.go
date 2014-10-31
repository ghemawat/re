/*
Package re combines regular expression matching with fmt.Scan like extraction
of sub-matches into caller-supplied objects.

	hostport := regexp.MustCompile(`(\w+):(\d+)`)

	var host string
	var port int
	if re.Find(hostport, "localhost:10000", &host, &port) {
		...
	}
*/
package re

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
)

// Find returns true iff the regular expression re occurs somewhere in
// input, and for every non-nil result, the corresponding regular expression
// sub-match is succesfully parsed and stored into *result[i].
//
// The following can be passed as result arguments to Find:
//
// nil: the corresponding sub-match is discarded without being saved.
//
// Pointer to a built-in numeric types (*int, *int8, *int16,
// *int32, *int64, *uint, *uintptr, *uint8, *uint16, *uint32,
// *uint64): The digits in the corresponding sub-match will be parsed
// and the result stored into the pointed-to object.  Find will return
// an error if the sub-match cannot be parsed successfully, or the
// parse result is out of range.  Note that byte is equivalent to
// uint8 and rune is equivalent to uint32.  These types are all
// handled via textual parsing of digits (this matches fmt's behavior)
// and therefore Find cannot be used to directly extract a single rune
// ot byte in the input; for that, parse into a string or []byte and
// use the first element.
//
// Pointer to string or []byte: the corresponding sub-match is
// stored in the pointed-to object.  When storing into a []byte, no
// copying is done, and the stored slice is an alias of the input.
//
// func([]byte) error: the function is called with the corresponding
// sub-match.  If the result is a non-nil error, the Find call
// fails with that error.
//
// An error is returned if a result does not have one of the preceding
// types.  Caveat: the set of supported types might be extended in the
// future.
func Find(re *regexp.Regexp, input []byte, result ...interface{}) error {
	matches := re.FindSubmatchIndex(input)
	if matches == nil {
		return fmt.Errorf(`re.Find: could not find "%s" in "%s"`,
			re, input)
	}
	if len(matches) < 2+2*len(result) {
		return fmt.Errorf(`re.Find: only got %d matches from "%s"; need at least %d`,
			len(matches)/2-1, re, len(result))
	}
	for i, r := range result {
		start, limit := matches[2+2*i], matches[2+2*i+1]
		if start < 0 || limit < 0 {
			// Sub-expression is missing; treat as empty.
			start = 0
			limit = 0
		}
		if err := assign(input[start:limit], r); err != nil {
			return err
		}
	}
	return nil
}

func assign(b []byte, r interface{}) error {
	switch v := r.(type) {
	case nil:
		// Discard the match.
	case *string:
		*v = string(b)
	case *[]byte:
		*v = b
	case *bool:
		if err := parseBool(b, v); err != nil {
			return err
		}
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
		// could treat as a number or a raw byte; treat as a number
		// (just like fmt)
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
		// could treat as a number or a rune; treat as a number
		// (just like fmt)
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
	case func([]byte) error:
		if err := v(b); err != nil {
			return err
		}
	default:
		t := reflect.ValueOf(r).Type()
		return parseError(fmt.Sprintf("unsupported type %s", t), b)
	}
	return nil
}

func parseBool(b []byte, v *bool) error {
	switch {
	case len(b) == 1 && b[0] == '0':
		*v = false
	case len(b) == 1 && b[0] == '1':
		*v = true
	case len(b) == 5 && bytes.EqualFold(b, []byte("false")):
		*v = false
	case len(b) == 4 && bytes.EqualFold(b, []byte("true")):
		*v = true
	default:
		return parseError("not a valid bool", b)
	}
	return nil
}

func parseError(explanation string, b []byte) error {
	return fmt.Errorf(`re.Find: parsing "%s": %s`, b, explanation)
}
