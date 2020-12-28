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
		"string": anyTostring,
		"tohtml": anyToTplHTML,
		"val":    trimNoValue,
	}
	for k, v := range f2 {
		funcMap[k] = v
	}
	return funcMap
}

func anyTostring(v interface{}) string {
	return gcast.ToString(v)
}

func anyToTplHTML(v interface{}) template.HTML {
	return template.HTML(anyTostring(v))
}

func trimNoValue(m interface{}, key string) interface{} {
	switch val := m.(type) {
	case map[string]interface{}:
		if v, ok := val[key]; ok {
			return v
		}
	case map[string]string:
		if v, ok := val[key]; ok {
			return v
		}
	}
	return ""
}
