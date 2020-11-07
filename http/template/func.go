package gmctemplate

import (
	"github.com/snail007/gmc/http/template/sprig"
	gmci18n "github.com/snail007/gmc/i18n"
	"github.com/snail007/gmc/util/cast"
	"html/template"
)

func addFunc() map[string]interface{}{
	funcMap:=sprig.FuncMap()
	f2:=map[string]interface{}{
		"tr":gmci18n.TrV,
		"trs":gmci18n.Tr,
		"string":anyTostring,
		"tohtml":anyToTplHTML,
		"val":trimNoValue,
	}
	for k,v:=range f2{
		funcMap[k]=v
	}
	return funcMap
}

func anyTostring(v interface{})string  {
	return cast.ToString(v)
}

func anyToTplHTML(v interface{})template.HTML  {
	return template.HTML(anyTostring(v))
}
func trimNoValue(m map[string]interface{},key string)interface{}  {
	if v,ok:=m[key];ok{
		return v
	}
	return ""
}