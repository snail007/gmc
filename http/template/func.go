package gmctemplate

import (
	"github.com/snail007/gmc/http/template/sprig"
	gmci18n "github.com/snail007/gmc/i18n"
	"github.com/snail007/gmc/util/cast"
)

func addFunc() map[string]interface{}{
	funcMap:=sprig.FuncMap()
	f2:=map[string]interface{}{
		"tr":gmci18n.TrV,
		"string":anyTostring,
	}
	for k,v:=range f2{
		funcMap[k]=v
	}
	return funcMap
}

func anyTostring(v interface{})string  {
	return cast.ToString(v)
}
