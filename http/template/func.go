package gtemplate

import (
	gi18n "github.com/snail007/gmc/gmc/i18n"
	"github.com/snail007/gmc/http/template/sprig"
	"github.com/snail007/gmc/util/cast"
	"html/template"
)

func addFunc() map[string]interface{} {
	funcMap := sprig.FuncMap()
	f2 := map[string]interface{}{
		"tr":     gi18n.TrV,
		"trs":    gi18n.Tr,
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
	return cast.ToString(v)
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
