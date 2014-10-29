package re

import (
	"regexp"
	"strconv"
	"strings"
)

func Find(r *regexp.Regexp, data string, results ...interface{}) bool {
	return assignResults(data, r.FindStringSubmatchIndex(data), results)
}

func assignResults(data string, matches []int, results []interface{}) bool {
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

func assign(s string, r interface{}) bool {
	switch v := r.(type) {
	case nil:
		// Discard the match.
	case *string:
		*v = s
	case *[]byte:
		*v = []byte(s)
	case *bool:
		if !parseBool(s, v) {
			return false
		}
	case *int:
		if i, err := strconv.ParseInt(s, 10, 64); err != nil {
			return false
		} else {
			if int64(int(i)) != i {
				return false
			}
			*v = int(i)
		}
	case *int8:
		if i, err := strconv.ParseInt(s, 10, 8); err != nil {
			return false
		} else {
			*v = int8(i)
		}
	case *int16:
		if i, err := strconv.ParseInt(s, 10, 16); err != nil {
			return false
		} else {
			*v = int16(i)
		}
	case *int32:
		if i, err := strconv.ParseInt(s, 10, 32); err != nil {
			return false
		} else {
			*v = int32(i)
		}
	case *uint:
		if u, err := strconv.ParseUint(s, 10, 64); err != nil {
			return false
		} else {
			if uint64(uint(u)) != u {
				return false
			}
			*v = uint(u)
		}
	case *uintptr:
		if u, err := strconv.ParseUint(s, 10, 64); err != nil {
			return false
		} else {
			if uint64(uintptr(u)) != u {
				return false
			}
			*v = uintptr(u)
		}
	case *uint8:
		if u, err := strconv.ParseUint(s, 10, 8); err != nil {
			return false
		} else {
			*v = uint8(u)
		}
	case *uint16:
		if u, err := strconv.ParseUint(s, 10, 16); err != nil {
			return false
		} else {
			*v = uint16(u)
		}
	case *uint32:
		if u, err := strconv.ParseUint(s, 10, 32); err != nil {
			return false
		} else {
			*v = uint32(u)
		}
	case *uint64:
		if u, err := strconv.ParseUint(s, 10, 64); err != nil {
			return false
		} else {
			*v = u
		}
	case func(string) bool:
		if !v(s) {
			return false
		}
	default:
		return false
	}
	// TODO: support for numeric radices
	// Find(..., CRadix(&x), ...)
	// Find(..., Hex(&x), ...)
	// Find(..., Octal(&x), ...)
	// Find(..., Binary(&x), ...)
	return true
}

func parseBool(s string, v *bool) bool {
	if s == "0" || strings.ToLower(s) == "false" {
		*v = false
		return true
	}
	if s == "1" || strings.ToLower(s) == "true" {
		*v = true
		return true
	}
	return false
}
