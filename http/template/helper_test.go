package gtemplate

import (
	"github.com/snail007/gmc/http/template/testdata"
	gctx "github.com/snail007/gmc/module/ctx"
	gmap "github.com/snail007/gmc/util/map"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRender(t *testing.T) {
	tpl := `a{{.user}}c`
	d, err := RenderString(tpl, gmap.M{"user": "b"})
	assert.Nil(t, err)
	assert.Equal(t, "abc", string(d))

	tpl = `a {{ .user}c`
	_, err = RenderString(tpl, gmap.M{"user": "b"})
	assert.NotNil(t, err)
	var testfunc = func(str string) string {
		return str + str
	}
	tpl = `a{{testfunc .user}}c`
	d, err = RenderStringWithFunc(tpl, gmap.M{"user": "b"}, map[string]interface{}{"testfunc": testfunc})
	//fmt.Println(err.Error())
	assert.Equal(t, "abbc", d)
}

func TestNewEmbedTemplateFS(t *testing.T) {
	tpl, _ := NewTemplate(gctx.NewCtx(), "")
	efs := NewEmbedTemplateFS(tpl, testdata.TplFS, "tpl").SetExt(".html")
	efs.Parse()
	assert.Len(t, efs.tpl.binData, 4)

	//err
	tpl.Parse()
	assert.Error(t, NewEmbedTemplateFS(tpl, testdata.TplFS, ".").Parse())
}
