package gstrings

import (
	"reflect"
	"strings"
	"unsafe"
)

func BytesRef(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func StringRef(b []byte) string {
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

func HasHTTPPrefix(str string) bool {
	return HasPrefixAny(strings.ToLower(str), "http://", "https://")
}
