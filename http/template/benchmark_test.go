// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More infomation at https://github.com/snail007/gmc

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
