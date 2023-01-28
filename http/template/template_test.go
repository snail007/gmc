// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gtemplate

import (
	"encoding/base64"
	gctx "github.com/snail007/gmc/module/ctx"
	"github.com/stretchr/testify/assert"
	"testing"
	gotemplate "text/template"
)

func TestParse(t *testing.T) {
	tpl.Execute("user/list", map[string]string{
		"head": "test",
	})
	//funcsM := sprig.GenericFuncMap()
	//fmt.Println(len(funcsM))
	//t.Fail()
}

func TestSetBinBase64(t *testing.T) {
	SetBinBase64(map[string]string{"test": base64.StdEncoding.EncodeToString([]byte("aaa"))})
	assert.Equal(t, defaultTpl.binData["test"], []byte("aaa"))
}

func TestSetBinBytes(t *testing.T) {
	SetBinBytes(map[string][]byte{"test": []byte("aaa")})
	assert.Equal(t, defaultTpl.binData["test"], []byte("aaa"))
}

func TestSetBinString(t *testing.T) {
	SetBinString(map[string]string{"test": "aaa"})
	assert.Equal(t, defaultTpl.BinData()["test"], []byte("aaa"))
}

func TestTemplate_Ctx(t *testing.T) {
	ctx := gctx.NewCtx()
	defaultTpl.SetCtx(ctx)
	assert.Same(t, ctx, defaultTpl.Ctx())
}

func TestTemplate_Ext(t *testing.T) {
	ext := ".txt"
	defaultTpl.SetExt(ext)
	assert.Equal(t, ext, defaultTpl.Ext())
}

func TestTemplate_Tpl(t *testing.T) {
	tpl := gotemplate.New("gmc")
	defaultTpl.SetTpl(tpl)
	assert.Same(t, tpl, defaultTpl.Tpl())
}

func TestTemplate_RootDir(t *testing.T) {
	dir := "tests/views"
	defaultTpl.SetRootDir(dir)
	assert.Contains(t, defaultTpl.RootDir(), dir)
	err := defaultTpl.Parse()
	assert.Nil(t, err)
}
