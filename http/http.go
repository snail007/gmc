package gmchttp

import (
	gcore "github.com/snail007/gmc/core"
	ghttputil "github.com/snail007/gmc/internal/util/http"
	"net/http"
)

func GetCtx(w http.ResponseWriter) gcore.Ctx {
	if v, ok := w.(*ghttputil.ResponseWriter); ok {
		ctx := v.Data("ctx")
		if ctx == nil {
			return nil
		}
		return ctx.(gcore.Ctx)
	}
	return nil
}
