// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gview_test

import (
	"bytes"
	gcore "github.com/snail007/gmc/core"
	gtemplate "github.com/snail007/gmc/http/template"
	assert2 "github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	assert := assert2.New(t)
	b := new(bytes.Buffer)
	tpl, _ := mockTpl()
	v := gcore.Providers.View("")(b, tpl)
	assert.NotNil(v)
}
func TestView_Render(t *testing.T) {
	assert := assert2.New(t)
	b := new(bytes.Buffer)
	tpl, _ := mockTpl()
	v := gcore.Providers.View("")(b, tpl)
	v.SetMap(map[string]interface{}{
		"name": "b",
	})
	v.Render("user/profile")
	assert.Equal("b", b.String())
}
func TestView_Layout(t *testing.T) {
	assert := assert2.New(t)
	b := new(bytes.Buffer)
	tpl, _ := mockTpl()
	v := gcore.Providers.View("")(b, tpl).Layout("layout/page")
	v.Set("name", "b")
	v.Render("user/profile")
	assert.Equal("abc", b.String())
}
func mockTpl() (t *gtemplate.Template, ctx gcore.Ctx) {
	ctx = gcore.Providers.Ctx("")()
	ctx.SetConfig(gcore.Providers.Config("")())
	t, _ = gtemplate.NewTemplate(ctx, "tests")
	t.Parse()
	return
}
