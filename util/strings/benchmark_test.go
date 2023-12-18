package gstrings

import "testing"

func Benchmark_NormalBytes2String(b *testing.B) {
	x := []byte("Hello GMC! Hello GMC! Hello GMC!")
	for i := 0; i < b.N; i++ {
		_ = string(x)
	}
}

func Benchmark_StringRef(b *testing.B) {
	x := []byte("Hello GMC! Hello GMC! Hello GMC!")
	for i := 0; i < b.N; i++ {
		_ = StringRef(x)
	}
}

func Benchmark_NormalString2Bytes(b *testing.B) {
	x := "Hello GMC! Hello GMC! Hello GMC!"
	for i := 0; i < b.N; i++ {
		_ = []byte(x)
	}
}

func Benchmark_BytesRef(b *testing.B) {
	x := "Hello GMC! Hello GMC! Hello GMC!"
	for i := 0; i < b.N; i++ {
		_ = BytesRef(x)
	}
}
