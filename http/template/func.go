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
	i18n, _ := gcore.Providers.I18n("")(ctx)
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
		return ""
	}
	if len(key) == 0 {
		return m
	}
	switch val := m.(type) {
	case map[string]interface{}:
		if v, ok := val[key[0]]; ok {
			return v
		}
	case map[string]string:
		if v, ok := val[key[0]]; ok {
			return v
		}
	}
	return ""
}
