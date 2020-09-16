package template

import (
	"testing"
)

func BenchmarkParse_one(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tpl.Execute("user/list", map[string]string{
			"head": "test",
		})
	}
}
func BenchmarkParse_two(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tpl.Execute("common/head", nil)
	}
}
