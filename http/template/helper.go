// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gtemplate

import (
	gcore "github.com/snail007/gmc/core"
	gctx "github.com/snail007/gmc/module/ctx"
	glog "github.com/snail007/gmc/module/log"
)

func Init(ctx gcore.Ctx) (tpl gcore.Template, err error) {
	cfg := ctx.Config()
	tpl, err = NewTemplate(ctx, cfg.GetString("template.dir"))
	if err != nil {
		return nil, err
	}
	tpl.Delims(cfg.GetString("template.delimiterleft"),
		cfg.GetString("template.delimiterright"))
	tpl.Extension(cfg.GetString("template.ext"))
	return tpl, nil
}

var ctx = gctx.NewCtx()

func init() {
	ctx.SetLogger(glog.DiscardLogger)
}

func RenderBytes(tpl []byte, data map[string]interface{}) (result []byte, err error) {
	return RenderBytesWithFunc(tpl, data, nil)
}

func RenderBytesWithFunc(tpl []byte, data map[string]interface{}, funcMap map[string]interface{}) (result []byte, err error) {
	t := New()
	t.SetCtx(ctx)
	t.DdisableLogging()
	t.DisableLoadDefaultBinData()
	t.SetBinBytes(map[string][]byte{"tpl": tpl})
	if len(funcMap) > 0 {
		t.Funcs(funcMap)
	}
	err = t.Parse()
	if err != nil {
		return
	}
	return t.Execute("tpl", data)
}

func RenderString(tpl string, data map[string]interface{}) (result string, err error) {
	v, e := RenderBytes([]byte(tpl), data)
	return string(v), e
}

func RenderStringWithFunc(tpl string, data map[string]interface{}, funcMap map[string]interface{}) (result string, err error) {
	r, e := RenderBytesWithFunc([]byte(tpl), data, funcMap)
	return string(r), e
}
