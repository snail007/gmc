package gmchttpserver

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"runtime/debug"

	"github.com/snail007/gmc/util/logutil"

	gmcrouter "github.com/snail007/gmc/http/router"
)

type APIServer struct {
	server            *http.Server
	address           string
	router            *gmcrouter.HTTPRouter
	logger            *log.Logger
	before            func(w http.ResponseWriter, r *http.Request) bool
	after             func(w http.ResponseWriter, r *http.Request, ps gmcrouter.Params, isPanic bool)
	handle404         func(w http.ResponseWriter, r *http.Request)
	handle500         func(w http.ResponseWriter, r *http.Request, ps gmcrouter.Params, err interface{})
	isShowErrorStack  bool
	certFile, keyFile string
}

func NewAPIServer(address string) *APIServer {
	api := &APIServer{
		server: &http.Server{
			TLSConfig: &tls.Config{},
		},
		address:          address,
		logger:           logutil.New(""),
		router:           gmcrouter.NewHTTPRouter(),
		isShowErrorStack: true,
	}
	api.server.Handler = api
	api.server.SetKeepAlivesEnabled(false)
	api.server.ErrorLog = api.logger
	return api
}

func (this *APIServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h, args, _ := this.router.Lookup(r.Method, r.URL.Path)
	if h != nil {
		if this.before != nil && !this.before(w, r) {
			return
		}
		status := ""
		if this.after != nil {
			defer func() {
				this.after(w, r, args, status != "")
			}()
		}
		err := this.call(func() { h(w, r, args) })
		if err != nil {
			status = fmt.Sprintf("%s", err)
		}
		switch status {
		case "__STOP__", "":
		default:
			//exception
			this.handler500(w, r, args, err)
		}
	} else {
		this.handler404(w, r)
	}
}
func (this *APIServer) Server() *http.Server {
	return this.server
}
func (this *APIServer) SetTLSFile(certFile, keyFile string) *APIServer {
	this.certFile, this.keyFile = certFile, keyFile
	return this
}
func (this *APIServer) SetLogger(l *log.Logger) *APIServer {
	this.logger = l
	return this
}
func (this *APIServer) Logger() *log.Logger {
	return this.logger
}
func (this *APIServer) Before(handle func(w http.ResponseWriter, r *http.Request) bool) *APIServer {
	this.before = handle
	return this
}
func (this *APIServer) After(handle func(w http.ResponseWriter, r *http.Request, ps gmcrouter.Params, isPanic bool)) *APIServer {
	this.after = handle
	return this
}
func (this *APIServer) Handle404(handle func(w http.ResponseWriter, r *http.Request)) *APIServer {
	this.handle404 = handle
	return this
}
func (this *APIServer) Handle500(handle func(w http.ResponseWriter, r *http.Request, ps gmcrouter.Params, err interface{})) *APIServer {
	this.handle500 = handle
	return this
}
func (this *APIServer) ShowErrorStack(isShow bool) *APIServer {
	this.isShowErrorStack = isShow
	return this
}
func (this *APIServer) API(path string, handle func(w http.ResponseWriter, r *http.Request, ps gmcrouter.Params)) *APIServer {
	this.router.HandleAny(path, handle)
	return this
}
func (this *APIServer) Run() (err error) {
	l, err := net.Listen("tcp", this.address)
	if err != nil {
		return
	}
	this.address = l.Addr().String()
	go func() {
		var err error
		if this.certFile != "" && this.keyFile != "" {
			this.logger.Printf("api server on https://%s", this.address)
			err = this.server.ServeTLS(l, this.certFile, this.keyFile)
		} else {
			this.logger.Printf("api server on http://%s", this.address)
			err = this.server.Serve(l)
		}
		if err != nil {
			this.logger.Printf("api server exited unexpectedly on %s, error : %s", this.address, err)
		}
	}()
	return
}
func (this *APIServer) handler404(w http.ResponseWriter, r *http.Request) *APIServer {
	if this.handle404 == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Page not found"))
	} else {
		this.handle404(w, r)
	}
	return this
}
func (this *APIServer) handler500(w http.ResponseWriter, r *http.Request, ps gmcrouter.Params, err interface{}) *APIServer {
	if this.handle500 == nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "text/plain")
		msg := fmt.Sprintf("Internal Server Error")
		if err != nil && this.isShowErrorStack {
			msg += fmt.Sprintf("\n%s\n", err) + string(debug.Stack())
		}
		w.Write([]byte(msg))
	} else {
		this.handle500(w, r, ps, err)
	}
	return this
}
func (s *APIServer) call(fn func()) (err interface{}) {
	func() {
		defer func() {
			err = recover()
		}()
		fn()
	}()
	return
}
