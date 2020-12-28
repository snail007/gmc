package main

import (
	"fmt"
	_ "github.com/snail007/gmc"
	gcore "github.com/snail007/gmc/core"
	gtemplate "github.com/snail007/gmc/http/template"
)

func main() {
	ctx := gcore.Providers.Ctx("")()
	t, err := gtemplate.NewTemplate(ctx, "views")
	if err != nil {
		fmt.Println(err)
		return
	}
	t.Funcs(map[string]interface{}{
		"add": add,
	})
	t.Extension(".html")
	err = t.Parse()
	if err != nil {
		fmt.Println(err)
		return
	}
	html, err := t.Execute("layout/list", map[string]string{
		"head": "test",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(html))
}
func add(a, b string) string {
	return a + ">>>" + b
}
