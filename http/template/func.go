// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gtemplate

import (
	gcore "github.com/snail007/gmc/core"
	"github.com/snail007/gmc/http/template/sprig"
	"github.com/snail007/gmc/util/cast"
	"html/template"
)

func addFunc(ctx gcore.Ctx) map[string]interface{} {
	i18n, _ := gcore.ProviderI18n()(ctx)
	funcMap := sprig.FuncMap()
	f2 := map[string]interface{}{
		"tr":     i18n.TrV,
		"trs":    i18n.Tr,
		"string": anyToString,
		"tohtml": anyToTplHTML,
		"val":    trimNoValue,
	}
	for k, v := range f2 {
		funcMap[k] = v
	}
	return funcMap
}

func anyToString(v interface{}) string {
	return gcast.ToString(v)
}

func anyToTplHTML(v interface{}) template.HTML {
	return template.HTML(anyToString(v))
}

func trimNoValue(m interface{}, key ...string) interface{} {
	if m == nil {
		if len(key) >= 2 {
			return key[1]
		}
		return ""
	}
	if len(key) == 0 {
		return m
	}
	var v interface{}
	var ok bool
	switch val := m.(type) {
	case map[string]interface{}:
		v, ok = val[key[0]]
	case map[string]string:
		v, ok = val[key[0]]
	}
	if ok {
		if gcast.ToString(v) == "" && len(key) >= 2 {
			return key[1]
		}
		return v
	}
	if len(key) >= 2 {
		return key[1]
	}
	return ""
}
