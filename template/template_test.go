package template

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
}
func TestMain(m *testing.M) {
	tpl.Delims("{{", "}}")
	tpl.Parse()
	os.Exit(m.Run())
}
