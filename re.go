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
	"encoding"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"time"
)

type input struct {
	s string
	b []byte
}

func (x input) str() string {
	if x.b != nil {
		return string(x.b)
	}
	return x.s
}

func (x input) bytes() []byte {
	if x.b == nil {
		return []byte(x.s)
	}
	return x.b
}

// Find returns true iff the regular expression re matches data, and
// for every i in [0..len(results)-1] if result[i] is non-nil, the ith
// sub-match (starting to count at zero) is succesfully parsed and
// stored into *result[i].
//
// TODO: Document all supported types.
// TODO: Give examples.
func Find(re *regexp.Regexp, data []byte, results ...interface{}) error {
	matches := re.FindSubmatchIndex(data)
	if matches == nil {
		return fmt.Errorf(`re.Find: could not find "%s" in "%s"`, re, data)
	}
	if len(matches) < 2+2*len(results) {
		return fmt.Errorf(`re.Find: only got %d matches from "%s"; need at least %d`, len(matches)/2-1, re, len(results))
	}
	for i, r := range results {
		start, limit := matches[2+2*i], matches[2+2*i+1]
		if start < 0 || limit < 0 {
			// Sub-expression is missing; treat as empty.
			start = 0
			limit = 0
		}
		if err := assign(data[start:limit], r); err != nil {
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
				return makeError("out of range for int", b)
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
	case *uint:
		if u, err := strconv.ParseUint(string(b), 0, 64); err != nil {
			return err
		} else {
			if uint64(uint(u)) != u {
				return makeError("out of range for uint", b)
			}
			*v = uint(u)
		}
	case *uintptr:
		if u, err := strconv.ParseUint(string(b), 0, 64); err != nil {
			return err
		} else {
			if uint64(uintptr(u)) != u {
				return makeError("out of range for uintptr", b)
			}
			*v = uintptr(u)
		}
	case *uint8:
		// could treat as a number or a raw byte; match fmt and treat like a number
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
		// could treat as a number or a rune; match fmt and treat like a number
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
	case *time.Duration:
		if d, err := time.ParseDuration(string(b)); err != nil {
			return err
		} else {
			*v = d
		}
	case encoding.TextUnmarshaler:
		if err := v.UnmarshalText(b); err != nil {
			return err
		}
	case func([]byte) error:
		if err := v(b); err != nil {
			return err
		}
	default:
		return makeError(fmt.Sprintf("unsupported type %s", reflect.ValueOf(r).Type()), b)
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
		return makeError("not a valid bool", b)
	}
	return nil
}

func makeError(explanation string, b []byte) error {
	return fmt.Errorf(`re.Find: parsing "%s": %s`, b, explanation)
}
