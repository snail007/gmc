package gmcrouter

import (
	gmcutil "github.com/snail007/gmc/util"
	"net"
	"net/http"
	"strings"
	"time"

	gmchttputil "github.com/snail007/gmc/util/http"
)

type Ctx struct {
	Response  http.ResponseWriter
	Request   *http.Request
	Param     Params
	timeUsed  time.Duration
	LocalAddr string
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

func (this *Ctx) SetParam(param Params) *Ctx {
	this.Param = param
	return this
}

// acquires the method cost time, only for middleware2 and middleware3.
func (this *Ctx) TimeUsed() time.Duration {
	return this.timeUsed
}

// sets the method cost time, only for middleware2 and middleware3, do not call this.
func (this *Ctx) SetTimeUsed(t time.Duration) *Ctx {
	this.timeUsed = t
	return this
}

// Write output data to response
func (this *Ctx) Write(data ...interface{}) (n int, err error) {
	return gmchttputil.Write(this.Response, data...)
}

// WriteE outputs data to response and sets http status code 500
func (this *Ctx) WriteE(data ...interface{}) (n int, err error) {
	this.Response.WriteHeader(http.StatusInternalServerError)
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
	return http.MethodPost == this.Request.Method
}

// IsGET returns true if the request is GET request.
func (this *Ctx) IsGET() bool {
	return http.MethodGet == this.Request.Method
}

// IsPUT returns true if the request is PUT request.
func (this *Ctx) IsPUT() bool {
	return http.MethodPut == this.Request.Method
}

// IsDELETE returns true if the request is DELETE request.
func (this *Ctx) IsDELETE() bool {
	return http.MethodDelete == this.Request.Method
}

// IsPATCH returns true if the request is PATCH request.
func (this *Ctx) IsPATCH() bool {
	return http.MethodPatch == this.Request.Method
}

// IsHEAD returns true if the request is HEAD request.
func (this *Ctx) IsHEAD() bool {
	return http.MethodHead == this.Request.Method
}

// IsOPTIONS returns true if the request is OPTIONS request.
func (this *Ctx) IsOPTIONS() bool {
	return http.MethodOptions == this.Request.Method
}

// Stop will exit controller method or api handle function at once
func (this *Ctx) Stop(msg ...interface{}) {
	gmchttputil.Stop(this.Response, msg...)
}

// ClientIP acquires the real client ip, search in X-Forwarded-For, X-Real-IP, request.RemoteAddr.
func (this *Ctx) ClientIP() (ip string) {
	// X-Forwarded-For
	xForwardedFor := this.Request.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		proxyIps := strings.Split(strings.Replace(xForwardedFor, " ", "", -1), ",")
		if len(proxyIps) > 0 {
			ip = proxyIps[0]
			return
		}
	}
	// X-Real-IP
	ip = this.Request.Header.Get("X-Real-IP")
	if ip != "" {
		return
	}
	// RemoteAddr
	ip, _, _ = net.SplitHostPort(this.Request.RemoteAddr)
	return
}

// NewPager create a new paginator used for template
func (this *Ctx) NewPager(perPage int, total int64) *gmcutil.Paginator {
	return gmcutil.NewPaginator(this.Request, perPage, total, "page")
}

// GET gets the first value associated with the given key in url query.
func (this *Ctx) GET(key string, Default ...string) (val string) {
	val = this.Request.URL.Query().Get(key)
	if val == "" && len(Default) > 0 {
		return Default[0]
	}
	return
}

// POST gets the first value associated with the given key in post body.
func (this *Ctx) POST(key string, Default ...string) (val string) {
	val = this.Request.PostFormValue(key)
	if val == "" && len(Default) > 0 {
		return Default[0]
	}
	return
}

// Redirect redirects to the url, using http location header, and sets http code 302
func (this *Ctx) Redirect(url string) (val string) {
	http.Redirect(this.Response, this.Request, url, http.StatusFound)
	this.Stop()
	return
}
