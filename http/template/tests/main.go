package main

import (
	"fmt"

	template "github.com/snail007/gmc/http/template"
)

func main() {
	t, err := template.NewTemplate("views")
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
