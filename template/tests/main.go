package main

import (
	"fmt"

	template "../../template"
)

func main() {
	t, err := template.New("views")
	if err != nil {
		fmt.Println(err)
		return
	}
	t.Funcs(map[string]interface{}{
		"add": add,
	})
	// t.Extension(".tpl")
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
