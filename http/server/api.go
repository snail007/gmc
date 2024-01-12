// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package ghttpserver

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"embed"
	"fmt"
	gcore "github.com/snail007/gmc/core"
	ghttputil "github.com/snail007/gmc/internal/util/http"
	"github.com/snail007/gmc/module/log"
	gfile "github.com/snail007/gmc/util/file"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type APIServer struct {
	listener          net.Listener
	listenerFactory   func(addr string) (net.Listener, error)
	server            *http.Server
	address           string
	router            gcore.HTTPRouter
	logger            gcore.Logger
	handle404         func(ctx gcore.Ctx)
	handle500         func(ctx gcore.Ctx, err interface{})
	isShowErrorStack  bool
	certFile, keyFile string
	middleware0       []gcore.Middleware
	middleware1       []gcore.Middleware
	middleware2       []gcore.Middleware
	middleware3       []gcore.Middleware
	isShutdown        bool
	ext               string
	config            gcore.Config
	remoteAddrDataMap *sync.Map
	connCnt           *int64
	ctx               gcore.Ctx
}

func NewAPIServerForProvider(ctx gcore.Ctx, address string) (*APIServer, error) {
	if address != "" {
		return NewAPIServer(ctx, address), nil
	}
	return NewDefaultAPIServer(ctx, ctx.Config())
}
func NewAPIServer(ctx gcore.Ctx, address string) *APIServer {
	api := &APIServer{
		server: &http.Server{
			TLSConfig: &tls.Config{},
		},
		address:           address,
		router:            gcore.ProviderHTTPRouter()(ctx),
		isShowErrorStack:  true,
		middleware0:       []gcore.Middleware{},
		middleware1:       []gcore.Middleware{},
		middleware2:       []gcore.Middleware{},
		middleware3:       []gcore.Middleware{},
		remoteAddrDataMap: &sync.Map{},
		ctx:               ctx,
	}
	ctx.SetAPIServer(api)
	api.server.Handler = api
	api.server.SetKeepAlivesEnabled(false)
	api.logger = gcore.ProviderLogger()(ctx, "")
	api.server.ErrorLog = func() *log.Logger {
		ns := api.logger.Namespace()
		if ns != "" {
			ns = "[" + ns + "]"
		}
		l := log.New(glog.NewIOWriter(api.logger.Writer()), ns, log.Lmicroseconds|log.LstdFlags)
		return l
	}()
	api.server.ConnState = api.connState
	api.connCnt = new(int64)
	return api
}

func NewDefaultAPIServer(ctx gcore.Ctx, config gcore.Config) (api *APIServer, err error) {
	api = NewAPIServer(ctx, config.GetString("apiserver.listen"))
	if config.GetBool("apiserver.tlsenable") {
		tlsCfg := &tls.Config{}
		if config.GetBool("apiserver.tlsclientauth") {
			tlsCfg.ClientAuth = tls.RequireAndVerifyClientCert
		}
		clientCertPool := x509.NewCertPool()
		caBytes, e := ioutil.ReadFile(config.GetString("apiserver.tlsclientsca"))
		if e != nil {
			return nil, e
		}
		ok := clientCertPool.AppendCertsFromPEM(caBytes)
		if !ok {
			api = nil
			err = gcore.ProviderError()().New("failed to parse tls clients root certificate")
			return
		}
		tlsCfg.ClientCAs = clientCertPool
		api.server.TLSConfig = tlsCfg
	}
	api.config = config
	api.ShowErrorStack(config.GetBool("apiserver.showerrorstack"))
	return
}

func (this *APIServer) ListenerFactory() func(addr string) (net.Listener, error) {
	return this.listenerFactory
}

func (this *APIServer) SetListenerFactory(listenerFactory func(addr string) (net.Listener, error)) {
	this.listenerFactory = listenerFactory
}

func (this *APIServer) Ctx() gcore.Ctx {
	return this.ctx
}

func (this *APIServer) initRequestCtx(w http.ResponseWriter, r *http.Request) gcore.Ctx {
	w = ghttputil.NewResponseWriter(w)
	c0 := this.ctx.CloneWithHTTP(w, r)
	item, _ := this.remoteAddrDataMap.Load(r.RemoteAddr)
	if v, ok := item.(remoteAddrItem); ok {
		c0.SetLocalAddr(v.localAddr)
		c0.SetRemoteAddr(v.remoteAddr)
		c0.SetConn(v.conn)
	}
	gcore.SetCtx(w, c0)
	return c0
}

func (this *APIServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqCtx := this.initRequestCtx(w, r)
	defer func() {
		// middleware3
		this.callMiddleware(reqCtx, this.middleware3)
	}()
	//middleware0
	if this.callMiddleware(reqCtx, this.middleware0) {
		return
	}

	h, params, _ := this.router.Lookup(r.Method, r.URL.Path)
	if h != nil {
		reqCtx.SetParam(params)
		// middleware1
		if this.callMiddleware(reqCtx, this.middleware1) {
			return
		}

		start := time.Now()
		status := ""
		err := this.call(func() { h(reqCtx.Response(), reqCtx.Request(), reqCtx.Param()) })
		reqCtx.SetTimeUsed(time.Now().Sub(start))
		if err != nil {
			status = fmt.Sprintf("%s", err)
			switch status {
			case "__STOP__", "__DIE__":
			default:
				//exception
				this.handler500(reqCtx, err)
			}
		}

		// middleware2
		if this.callMiddleware(reqCtx, this.middleware2) {
			return
		}

	} else {
		this.handler404(reqCtx)
	}
}
func (this *APIServer) Address() string {
	return this.address
}
func (this *APIServer) Server() *http.Server {
	return this.server
}
func (this *APIServer) Router() gcore.HTTPRouter {
	return this.router
}
func (this *APIServer) SetTLSFile(certFile, keyFile string) {
	this.certFile, this.keyFile = certFile, keyFile
}
func (this *APIServer) SetLogger(l gcore.Logger) {
	this.logger = l
}
func (this *APIServer) Logger() gcore.Logger {
	return this.logger
}
func (this *APIServer) AddMiddleware0(m gcore.Middleware) {
	this.middleware0 = append(this.middleware0, m)
}
func (this *APIServer) AddMiddleware1(m gcore.Middleware) {
	this.middleware1 = append(this.middleware1, m)
}
func (this *APIServer) AddMiddleware2(m gcore.Middleware) {
	this.middleware2 = append(this.middleware2, m)
}
func (this *APIServer) AddMiddleware3(m gcore.Middleware) {
	this.middleware3 = append(this.middleware3, m)
}
func (this *APIServer) SetNotFoundHandler(handle func(ctx gcore.Ctx)) {
	this.handle404 = handle
}
func (this *APIServer) SetErrorHandler(handle func(ctx gcore.Ctx, err interface{})) {
	this.handle500 = handle
}
func (this *APIServer) ShowErrorStack(isShow bool) {
	this.isShowErrorStack = isShow
}

func (this *APIServer) Ext(ext string) {
	this.ext = ext
}

func (this *APIServer) API(path string, handle func(ctx gcore.Ctx), ext ...string) {
	// default
	ext1 := this.ext
	if len(ext) > 0 {
		// cover
		ext1 = ext[0]
	}
	this.router.HandleAny(path+ext1, func(w http.ResponseWriter, _ *http.Request, _ gcore.Params) {
		handle(gcore.GetCtx(w))
	})
}

func (this *APIServer) Group(path string) gcore.APIServer {
	newAPI := *this
	newAPI.router = this.router.Group(path)
	return &newAPI
}

// PrintRouteTable dump all routes into `w`, if `w` is nil, os.Stdout will be used.
func (this *APIServer) PrintRouteTable(w io.Writer) {
	this.router.PrintRouteTable(w)
}

func (this *APIServer) createListener() (err error) {
	defer func() {
		if err == nil {
			this.address = this.listener.Addr().String()
		}
	}()
	if this.listener != nil {
		return
	}
	if this.listenerFactory != nil {
		this.listener, err = this.listenerFactory(this.address)
		return
	}
	this.listener, err = net.Listen("tcp", this.address)
	return
}

func (this *APIServer) Run() (err error) {
	err = this.createListener()
	if err != nil {
		return
	}
	if this.config != nil && this.config.GetBool("apiserver.printroute") {
		this.router.PrintRouteTable(os.Stdout)
	}
	go func() {
		var err error
		if this.certFile != "" && this.keyFile != "" {
			this.logger.Infof("api server on https://%s", this.address)
			err = this.server.ServeTLS(this.listener, this.certFile, this.keyFile)
		} else {
			this.logger.Infof("api server on http://%s", this.address)
			err = this.server.Serve(this.listener)
		}
		if err != nil {
			if strings.Contains(err.Error(), "closed") {
				if this.isShutdown {
					this.logger.Infof("api server graceful shutdown on %s", this.address)
				} else {
					this.logger.Infof("api server closed on %s", this.address)
					this.server.Close()
				}
			} else {
				this.logger.Warnf("api server exited unexpectedly on %s, error : %s", this.address, err)
			}
		}
	}()
	return
}
func (s *APIServer) ActiveConnCount() int64 {
	return atomic.LoadInt64(s.connCnt)
}

// ConnState count the active conntions
func (s *APIServer) connState(c net.Conn, st http.ConnState) {
	switch st {
	case http.StateNew:
		atomic.AddInt64(s.connCnt, 1)
		s.remoteAddrDataMap.Store(c.RemoteAddr().String(), remoteAddrItem{
			remoteAddr: c.RemoteAddr().String(),
			localAddr:  c.LocalAddr().String(),
			conn:       c,
		})
	case http.StateClosed:
		atomic.AddInt64(s.connCnt, -1)
		s.remoteAddrDataMap.Delete(c.RemoteAddr().String())
	}
}
func (this *APIServer) handler404(ctx gcore.Ctx) {
	if this.handle404 == nil {
		ctx.Response().WriteHeader(http.StatusNotFound)
		ctx.Response().Write([]byte("Page not found"))
	} else {
		this.handle404(ctx)
	}
}
func (this *APIServer) handler500(ctx gcore.Ctx, err interface{}) {
	msg := gcore.ProviderError()().StackError(err)
	if this.handle500 == nil {
		ctx.WriteHeader(http.StatusInternalServerError)
		ctx.Response().Header().Set("Content-Type", "text/plain")
		info := fmt.Sprintf("Internal Server Error")
		if err != nil && this.isShowErrorStack {
			info += "\n" + msg
		}
		ctx.Write(info)
	} else {
		this.handle500(ctx, err)
	}
	this.logger.Warn(ctx.Request().URL.String() + "\n" + msg)
}
func (s *APIServer) call(fn func()) (err interface{}) {
	func() {
		defer gcore.ProviderError()().Recover(func(e interface{}) {
			err = gcore.ProviderError()().Wrap(e)
		})
		fn()
	}()
	return
}
func (s *APIServer) callMiddleware(ctx gcore.Ctx, middleware []gcore.Middleware) (isStop bool) {
	for _, fn := range middleware {
		func() {
			defer gcore.ProviderError()().Recover(func(e interface{}) {
				s.logger.Warn("middleware panic error : %s", gcore.ProviderError()().StackError(e))
				isStop = false
			})
			isStop = fn(ctx)
		}()
		if isStop {
			return
		}
	}
	return
}

// Init implements gcore.Service Init
func (s *APIServer) Init(cfg gcore.Config) (err error) {

	return
}

// Start implements gcore.Service Start
func (this *APIServer) Start() (err error) {
	this.Run()
	return
}

// Stop implements gcore.Service Stop
func (this *APIServer) Stop() {
	this.server.Close()
}

// GracefulStop implements gcore.Service GracefulStop
func (this *APIServer) GracefulStop() {
	if this.isShutdown {
		return
	}
	this.isShutdown = true
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	this.server.Shutdown(ctx)
	return
}

// SetLog implements gcore.Service SetLog
func (this *APIServer) SetLog(l gcore.Logger) {
	this.logger = l
}

// InjectListeners implements gcore.Service InjectListeners
func (this *APIServer) InjectListeners(l []net.Listener) {
	this.listener = l[0]
}

// Listeners implements gcore.Service Listeners
func (this *APIServer) Listeners() []net.Listener {
	return []net.Listener{this.listener}
}

func (this *APIServer) Listener() net.Listener {
	return this.listener
}

func (s *APIServer) ServeEmbedFS(fs embed.FS, urlPath string) {
	serveEmbedFS(s.router, fs, urlPath)
}

func (s *APIServer) ServeFiles(rootPath, urlPath string) {
	serveFiles(s.router, gfile.Abs(rootPath), urlPath)
}
