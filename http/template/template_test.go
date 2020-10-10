// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More infomation at https://github.com/snail007/gmc

package gmctemplate

import (
	"os"
	"testing"
)

var (
	tpl, _ = New("tests/views")
)

func TestParse(t *testing.T) {
	tpl.Execute("user/list", map[string]string{
		"head": "test",
	})
	//funcsM := sprig.GenericFuncMap()
	//fmt.Println(len(funcsM))
	//t.Fail()
}
func TestMain(m *testing.M) {
	tpl.Delims("{{", "}}")
	tpl.Parse()
	os.Exit(m.Run())
}
