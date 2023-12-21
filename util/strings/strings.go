package gstrings

import (
	"strings"
	"unsafe"
)

func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func HasPrefixAny(str string, prefix ...string) bool {
	for _, v := range prefix {
		if strings.HasPrefix(str, v) {
			return true
		}
	}
	return false
}

func HasSuffixAny(str string, suffix ...string) bool {
	for _, v := range suffix {
		if strings.HasSuffix(str, v) {
			return true
		}
	}
	return false
}

func HasHTTPPrefix(str string) bool {
	return HasPrefixAny(strings.ToLower(str), "http://", "https://")
}

// Replace str from a list of old, new string
// pairs. Replacements are performed in the order they appear in the
// target string, without overlapping matches. The old string
// comparisons are done in argument order.
//
// Replace panics if given an odd number of oldNew arguments.
func Replace(str string, oldNew ...string) string {
	if len(oldNew) == 0 {
		return str
	}
	r := strings.NewReplacer(oldNew...)
	return r.Replace(str)
}
