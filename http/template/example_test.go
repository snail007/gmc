// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gtemplate

import (
	"fmt"
	gcore "github.com/snail007/gmc/core"
	glog "github.com/snail007/gmc/module/log"
	"io/ioutil"
)

func Example() {
	defaultTpl.binData = map[string][]byte{}
	ctx := gcore.ProviderCtx()()
	ctx.Logger().SetOutput(glog.NewLoggerWriter(ioutil.Discard))
	ctx.SetConfig(gcore.ProviderConfig()())
	tpl, err := NewTemplate(ctx, "tests/views")
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
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
	// Output: hello
}

func add(a, b string) string {
	return a + ">>>" + b
}
