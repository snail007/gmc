package gmctemplate

import (
	gmci18n "github.com/snail007/gmc/i18n"
	"github.com/snail007/gmc/util/castutil"
)

func addFunc() map[string]interface{}{
	funcMap:=map[string]interface{}{
		"tr":gmci18n.TrV,
		"string":anyTostring,
	}
	return funcMap
}

func anyTostring(v interface{})string  {
	return castutil.ToString(v)
}
