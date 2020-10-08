package gmci18n_test

import (
	"bytes"
	gmcconfig "github.com/snail007/gmc/config"
	gmci18n "github.com/snail007/gmc/i18n"
	assert2 "github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestNew(t *testing.T) {
	assert:=assert2.New(t)
	cfg:=gmcconfig.New()
	cfg.Set("i18n.enable",true)
	cfg.Set("i18n.dir","tests")
	cfg.Set("i18n.default","none")
	gmci18n.Init(cfg)
	tool:=gmci18n.NewI18nTool(gmci18n.Tr)
	assert.Equal("你好",tool.Lang("zh-cn").Tr("001"))
	assert.Equal("Hello",tool.Lang("en-us").Tr("001"))
	assert.Equal("001",tool.Lang("none").Tr("001"))
	assert.Equal("default",tool.Lang("none").Tr("001","default"))
}
func TestNew_2(t *testing.T) {
	assert:=assert2.New(t)
	cfg:=gmcconfig.New()
	cfg.Set("i18n.enable",true)
	cfg.Set("i18n.dir","tests")
	cfg.Set("i18n.default","zh-cn")
	gmci18n.Init(cfg)
	tool:=gmci18n.NewI18nTool(gmci18n.Tr)
	assert.Equal("你好",tool.Lang("zh-cn").Tr("001"))
	assert.Equal("Hello",tool.Lang("en-us").Tr("001"))
	assert.Equal("你好",tool.Lang("none").Tr("001"))
	assert.Equal("你好",tool.Lang("none").Tr("001","default"))
}

func TestParseAcceptLanguageStr(t *testing.T) {
	assert:=assert2.New(t)
	v,_:=gmci18n.ParseAcceptLanguageStr("")
	assert.Equal("",v)
	v,_=gmci18n.ParseAcceptLanguageStr("zh-cn")
	assert.Equal("zh-CN",v)
	v,_=gmci18n.ParseAcceptLanguageStr("zh")
	assert.Equal("zh",v)
}
func TestParseAcceptLanguage(t *testing.T) {
	assert:=assert2.New(t)
	r,_:=http.NewRequest("GET","http://foo/foo",new(bytes.Buffer))
	r.Header.Set("Accept-Language","zh-CN,zh;q=0.9,en;q=0.8")
	v,_:=gmci18n.ParseAcceptLanguage(r)
	assert.Equal("zh-CN",v[0])
	r.Header.Set("Accept-Language","")
	v,_=gmci18n.ParseAcceptLanguage(r)
	assert.Len(v,0)
	r.Header.Set("Accept-Language","zh;q=0.9,en;q=0.8")
	v,_=gmci18n.ParseAcceptLanguage(r)
	assert.Equal("zh",v[0])
}

func TestMatch(t *testing.T) {
	assert:=assert2.New(t)
	cfg:=gmcconfig.New()
	cfg.Set("i18n.enable",true)
	cfg.Set("i18n.dir","tests")
	cfg.Set("i18n.default","zh-cn")
	gmci18n.Init(cfg)
	r,_:=http.NewRequest("GET","http://foo/foo",new(bytes.Buffer))
	r.Header.Set("Accept-Language","zh")
	v,err:=gmci18n.MatchAcceptLanguageT(r)
	assert.Nil(err)
	assert.Equal("zh-CN",v.String())
}