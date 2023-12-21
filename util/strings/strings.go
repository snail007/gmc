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
