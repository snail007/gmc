// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More infomation at https://github.com/snail007/gmc

package template

import "fmt"

func Example() {
	tpl, err := New("tests/views")
	if err != nil {
		fmt.Println(err)
		return
	}
	tpl.Funcs(map[string]interface{}{
		"add": add,
	})
	tpl.Extension(".html")
	err = tpl.Parse()
	if err != nil {
		fmt.Println(err)
		return
	}
	b, err := tpl.Execute("user/list", map[string]string{
		"head": "hello",
	})
	fmt.Println(string(b))
	// Output: hello
}

func add(a, b string) string {
	return a + ">>>" + b
}
