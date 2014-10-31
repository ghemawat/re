package re

import (
	"bytes"
	"regexp"
	"strconv"
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

func Find(r *regexp.Regexp, data []byte, results ...interface{}) bool {
	return assignResults(data, r.FindSubmatchIndex(data), results)
}

func assignResults(data []byte, matches []int, results []interface{}) bool {
	if matches == nil {
		return false
	}
	if len(matches) < 2+2*len(results) {
		// Not enough matches to fill all results
		return false
	}
	for i, r := range results {
		start, limit := matches[2+2*i], matches[2+2*i+1]
		if start < 0 || limit < 0 {
			// Sub-expression is missing.  Allow match to nil.
			switch r {
			case nil:
				return true
			}
			return false
		}
		if !assign(data[start:limit], r) {
			return false
		}
	}
	return true
}

func assign(b []byte, r interface{}) bool {
	switch v := r.(type) {
	case nil:
		// Discard the match.
	case *string:
		*v = string(b)
	case *[]byte:
		*v = b
	case *bool:
		if !parseBool(b, v) {
			return false
		}
	case *int:
		if i, err := strconv.ParseInt(string(b), 0, 64); err != nil {
			return false
		} else {
			if int64(int(i)) != i {
				return false
			}
			*v = int(i)
		}
	case *int8:
		if i, err := strconv.ParseInt(string(b), 0, 8); err != nil {
			return false
		} else {
			*v = int8(i)
		}
	case *int16:
		if i, err := strconv.ParseInt(string(b), 0, 16); err != nil {
			return false
		} else {
			*v = int16(i)
		}
	case *int32:
		if i, err := strconv.ParseInt(string(b), 0, 32); err != nil {
			return false
		} else {
			*v = int32(i)
		}
	case *uint:
		if u, err := strconv.ParseUint(string(b), 0, 64); err != nil {
			return false
		} else {
			if uint64(uint(u)) != u {
				return false
			}
			*v = uint(u)
		}
	case *uintptr:
		if u, err := strconv.ParseUint(string(b), 0, 64); err != nil {
			return false
		} else {
			if uint64(uintptr(u)) != u {
				return false
			}
			*v = uintptr(u)
		}
	case *uint8:
		// could treat as a number or a raw byte; match fmt and treat like a number
		if u, err := strconv.ParseUint(string(b), 0, 8); err != nil {
			return false
		} else {
			*v = uint8(u)
		}
	case *uint16:
		if u, err := strconv.ParseUint(string(b), 0, 16); err != nil {
			return false
		} else {
			*v = uint16(u)
		}
	case *uint32:
		// could treat as a number or a rune; match fmt and treat like a number
		if u, err := strconv.ParseUint(string(b), 0, 32); err != nil {
			return false
		} else {
			*v = uint32(u)
		}
	case *uint64:
		if u, err := strconv.ParseUint(string(b), 0, 64); err != nil {
			return false
		} else {
			*v = u
		}
	case func([]byte) bool:
		if !v(b) {
			return false
		}
	default:
		return false
	}
	return true
}

func parseBool(b []byte, v *bool) bool {
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
		return false
	}
	return true
}
