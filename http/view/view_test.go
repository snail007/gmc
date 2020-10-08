package gmcview_test

import (
	"bytes"
	gmctemplate "github.com/snail007/gmc/http/template"
	gmcview "github.com/snail007/gmc/http/view"
	assert2 "github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	assert := assert2.New(t)
	b := new(bytes.Buffer)
	tpl := mockTpl()
	v := gmcview.New(b, tpl)
	assert.NotNil(v)
}
func TestView_Render(t *testing.T) {
	assert := assert2.New(t)
	b := new(bytes.Buffer)
	tpl := mockTpl()
	v := gmcview.New(b, tpl)
	v.SetMap(map[string]interface{}{
		"name":"b",
	})
	v.Render("user/profile")
	assert.Equal( "b",b.String())
}
func TestView_Layout(t *testing.T) {
	assert := assert2.New(t)
	b := new(bytes.Buffer)
	tpl := mockTpl()
	v := gmcview.New(b, tpl).Layout("layout/page")
	v.Set("name", "b")
	v.Render("user/profile")
	assert.Equal( "abc",b.String())
}
func mockTpl() *gmctemplate.Template {
	t, _ := gmctemplate.New("tests")
	t.Parse()
	return t
}
