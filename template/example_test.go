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
