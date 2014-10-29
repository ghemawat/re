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
		s := data[matches[2+2*i]:matches[2+2*i+1]]
		switch v := r.(type) {
		case nil:
			// Discard the match.
		case *string:
			*v = s
		case *[]byte:
			*v = []byte(s)
		case *byte:
			if len(s) != 1 {
				return false
			}
			*v = s[0]
		case *rune:
			for i, r := range s {
				if i > 0 {
					// More than one rune
					return false
				}
				*v = r
			}
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
		case func(string) bool:
			if !v(s) {
				return false
			}
		default:
			// TODO: Try Scan interface
			return false
		}
		// TODO: other types:
		//   uint8, uint16, uint32, uint64
		//   int8, int16, int32, int64
		//   float32, float64
		//   complex?
		// TODO: support for numeric radices
	}
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
