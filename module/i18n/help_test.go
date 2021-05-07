// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gi18n

import (
	"bytes"
	gcore "github.com/snail007/gmc/core"
	assert2 "github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestNew(t *testing.T) {
	assert := assert2.New(t)
	cfg := gcore.ProviderConfig()()
	cfg.Set("i18n.enable", true)
	cfg.Set("i18n.dir", "tests")
	cfg.Set("i18n.default", "none")
	e := Init(cfg)
	assert.Nil(e)
	assert.Equal("你好", I18N.Tr("zh-cn", "001"))
	assert.Equal("Hello", I18N.Tr("en-us", "001"))
	assert.Equal("001", I18N.Tr("none", "001"))
	assert.Equal("default", I18N.Tr("none", "001", "default"))
}
func TestNew_2(t *testing.T) {
	assert := assert2.New(t)
	cfg := gcore.ProviderConfig()()
	cfg.Set("i18n.enable", true)
	cfg.Set("i18n.dir", "tests")
	cfg.Set("i18n.default", "zh-cn")
	e := Init(cfg)
	assert.Nil(e)
	assert.Equal("你好", I18N.Tr("zh-cn", "001"))
	assert.Equal("Hello", I18N.Tr("en-us", "001"))
	assert.Equal("你好", I18N.Tr("none", "001"))
	assert.Equal("你好", I18N.Tr("none", "001", "default"))
}

func TestParseAcceptLanguageStr(t *testing.T) {
	assert := assert2.New(t)
	i18n := I18n{}
	v, _ := i18n.ParseAcceptLanguageStr("")
	assert.Equal("", v)
	v, _ = i18n.ParseAcceptLanguageStr("zh-cn")
	assert.Equal("zh-CN", v)
	v, _ = i18n.ParseAcceptLanguageStr("zh")
	assert.Equal("zh", v)
}
func TestParseAcceptLanguage(t *testing.T) {
	assert := assert2.New(t)
	i18n := I18n{}
	r, _ := http.NewRequest("GET", "http://foo/foo", new(bytes.Buffer))
	r.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	v, _ := i18n.ParseAcceptLanguage(r)
	assert.Equal("zh-CN", v[0])
	r.Header.Set("Accept-Language", "")
	v, _ = i18n.ParseAcceptLanguage(r)
	assert.Len(v, 0)
	r.Header.Set("Accept-Language", "zh;q=0.9,en;q=0.8")
	v, _ = i18n.ParseAcceptLanguage(r)
	assert.Equal("zh", v[0])
}

func TestMatch(t *testing.T) {
	assert := assert2.New(t)
	cfg := gcore.ProviderConfig()()
	cfg.Set("i18n.enable", true)
	cfg.Set("i18n.dir", "tests")
	cfg.Set("i18n.default", "zh-cn")
	e := Init(cfg)
	assert.Nil(e)
	r, _ := http.NewRequest("GET", "http://foo/foo", new(bytes.Buffer))
	r.Header.Set("Accept-Language", "zh")
	v, err := I18N.MatchAcceptLanguageT(r)
	assert.Nil(err)
	assert.Equal("zh-CN", v.String())
}
