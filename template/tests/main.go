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
	err = t.Parse()
	if err != nil {
		fmt.Println(err)
		return
	}
	html, err := t.Execute("layout/list.html", map[string]string{
		"head": "test",
	})
	fmt.Println(string(html), err, t)
}
