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
// Write output data to response
func (this *Ctx) Write(data ...interface{}) (n int, err error) {
	return gmchttputil.Write(this.Response, data...)
}
// WriteHeader sets http code in response
func (this *Ctx) WriteHeader(statusCode int) {
	this.Response.WriteHeader(statusCode)
}
// StatusCode returns http code in response, if not set, default is 200.
func (this *Ctx) StatusCode() int {
	return gmchttputil.StatusCode(this.Response)
}
// WriteCount acquires outgoing bytes count by writer
func (this *Ctx) WriteCount() int64 {
	return gmchttputil.WriteCount(this.Response)
}
// IsPOST returns true if the request is POST request.
func (this *Ctx) IsPOST() bool {
	return http.MethodPost==this.Request.Method
}
// IsGET returns true if the request is GET request.
func (this *Ctx) IsGET() bool {
	return http.MethodGet==this.Request.Method
}
// IsPUT returns true if the request is PUT request.
func (this *Ctx) IsPUT() bool {
	return http.MethodPut==this.Request.Method
}
// IsDELETE returns true if the request is DELETE request.
func (this *Ctx) IsDELETE() bool {
	return http.MethodDelete==this.Request.Method
}
// IsPATCH returns true if the request is PATCH request.
func (this *Ctx) IsPATCH() bool {
	return http.MethodPatch==this.Request.Method
}
// IsHEAD returns true if the request is HEAD request.
func (this *Ctx) IsHEAD() bool {
	return http.MethodHead==this.Request.Method
}
// IsOPTIONS returns true if the request is OPTIONS request.
func (this *Ctx) IsOPTIONS() bool {
	return http.MethodOptions==this.Request.Method
}
//Stop will exit controller method or api handle function at once
func (this *Ctx) Stop(msg ...interface{}) {
	gmchttputil.Stop(this.Response, msg...)
}