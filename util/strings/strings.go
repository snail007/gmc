package strings

import (
	"reflect"
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
