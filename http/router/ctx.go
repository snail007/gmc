package gmcrouter

import (
	"net/http"

	gmchttputil "github.com/snail007/gmc/util/httputil"
)

type Ctx struct {
	Response http.ResponseWriter
	Request  *http.Request
	Param    Params
}

func NewCtx(w http.ResponseWriter, r *http.Request, ps ...Params) *Ctx {
	var ps0 Params
	if len(ps) > 0 {
		ps0 = ps[0]
	} else {
		ps0 = Params{}
	}
	return &Ctx{
		Response: w,
		Request:  r,
		Param:    ps0,
	}
}

func (this *Ctx) Write(data ...interface{}) (n int, err error) {
	return gmchttputil.Write(this.Response, data...)
}

func (this *Ctx) WriteHeader(statusCode int) {
	this.Response.WriteHeader(statusCode)
}

func (this *Ctx) StatusCode() int {
	return gmchttputil.StatusCode(this.Response)
}

//WriteCount acquires outgoing bytes count by writer
func (this *Ctx) WriteCount() int64 {
	return gmchttputil.WriteCount(this.Response)
}
