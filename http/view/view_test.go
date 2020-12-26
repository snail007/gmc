package gview_test

import (
	"bytes"
	gtemplate "github.com/snail007/gmc/http/template"
	gview "github.com/snail007/gmc/http/view"
	assert2 "github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	assert := assert2.New(t)
	b := new(bytes.Buffer)
	tpl := mockTpl()
	v := gview.New(b, tpl)
	assert.NotNil(v)
}
func TestView_Render(t *testing.T) {
	assert := assert2.New(t)
	b := new(bytes.Buffer)
	tpl := mockTpl()
	v := gview.New(b, tpl)
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
	v := gview.New(b, tpl).Layout("layout/page")
	v.Set("name", "b")
	v.Render("user/profile")
	assert.Equal( "abc",b.String())
}
func mockTpl() *gtemplate.Template {
	t, _ := gtemplate.NewTemplate("tests")
	t.Parse()
	return t
}
