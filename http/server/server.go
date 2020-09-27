package gmchttpserver

import (
	"compress/gzip"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net"
	"net/http"
	"path/filepath"
	"reflect"
	"runtime/debug"
	"strings"
	"sync/atomic"
	"time"

	gmcsession "github.com/snail007/gmc/http/session"
	gmcfilestore "github.com/snail007/gmc/http/session/filestore"
	gmcmemorystore "github.com/snail007/gmc/http/session/memorystore"
	gmcredisstore "github.com/snail007/gmc/http/session/redisstore"
	gmctemplate "github.com/snail007/gmc/http/template"

	gmccontroller "github.com/snail007/gmc/http/controller"

	gmcconfig "github.com/snail007/gmc/config/gmc"
	gmcrouter "github.com/snail007/gmc/http/router"
	"github.com/snail007/gmc/http/server/ctxvalue"
	gmchttputil "github.com/snail007/gmc/util/httputil"
	"github.com/snail007/gmc/util/logutil"
)

var (
	bindata = map[string][]byte{}
)

func SetBinData(data map[string]string) {
	bindata = map[string][]byte{}
	for k, v := range data {
		b, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			panic("init static bin data fail, error: " + err.Error())
		}
		bindata[k] = b
	}
}

type HTTPServer struct {
	tpl          *gmctemplate.Template
	sessionStore gmcsession.Store
	router       *gmcrouter.HTTPRouter
	logger       *log.Logger
	addr         string
	listener     net.Listener
	server       *http.Server
	connCnt      *int64
	config       *gmcconfig.GMCConfig
	handler40x   func(ctx *gmcrouter.Ctx, tpl *gmctemplate.Template)
	handler50x   func(c gmccontroller.IController, err interface{})
	//just for testing
	isTestNotClosedError bool
	staticDir            string
	staticUrlpath        string
	middleware0          []func(ctx *gmcrouter.Ctx, server *HTTPServer) (isStop bool)
	middleware1          []func(ctx *gmcrouter.Ctx, server *HTTPServer) (isStop bool)
	middleware2          []func(ctx *gmcrouter.Ctx, server *HTTPServer) (isStop bool)
	middleware3          []func(ctx *gmcrouter.Ctx, server *HTTPServer) (isStop bool)
}

func New() *HTTPServer {
	return &HTTPServer{
		middleware0: []func(ctx *gmcrouter.Ctx, server *HTTPServer) (isStop bool){},
		middleware1: []func(ctx *gmcrouter.Ctx, server *HTTPServer) (isStop bool){},
		middleware2: []func(ctx *gmcrouter.Ctx, server *HTTPServer) (isStop bool){},
		middleware3: []func(ctx *gmcrouter.Ctx, server *HTTPServer) (isStop bool){},
	}
}

//Init implements service.Services Init
func (s *HTTPServer) Init(cfg *gmcconfig.GMCConfig) (err error) {
	connCnt := int64(0)
	s.server = &http.Server{}
	s.logger = logutil.New("")
	s.connCnt = &connCnt
	s.config = cfg
	s.isTestNotClosedError = false
	s.server.ConnState = s.connState
	s.server.Handler = s

	//init base objects
	err = s.initBaseObjets()
	return
}
func (s *HTTPServer) initBaseObjets() (err error) {
	s.tpl, err = gmctemplate.New(s.config.GetString("template.dir"))
	if err != nil {
		return
	}
	//init session store
	err = s.initSessionStore()
	if err != nil {
		return
	}
	//init http server tls configuration
	err = s.initTLSConfig()
	if err != nil {
		return
	}
	//init http server router
	s.router = gmcrouter.NewHTTPRouter()
	s.router.SetHandle50x(s.handle50x)
	s.addr = s.config.GetString("httpserver.listen")

	//init static files handler, must be after router inited
	s.initStatic()
	return
}

func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rctx := r.Context()
	rctx = context.WithValue(rctx, ctxvalue.CtxValueKey, ctxvalue.CtxValue{
		Tpl:          s.tpl,
		SessionStore: s.sessionStore,
		Router:       s.router,
		Config:       s.config,
	})
	w = gmchttputil.NewResponseWriter(w)
	r = r.WithContext(rctx)
	c0 := gmcrouter.NewCtx(w, r)

	defer func() {
		// middleware3
		s.callMiddleware(c0, s.middleware3)
	}()

	//middleware0
	if s.callMiddleware(c0, s.middleware0) {
		return
	}

	h, args, _ := s.router.Lookup(r.Method, r.URL.Path)
	if h != nil {
		c := gmcrouter.NewCtx(w, r, args)
		c0 = c
		// middleware1
		if s.callMiddleware(c, s.middleware1) {
			return
		}

		h(w, r, args)

		// middleware2
		if s.callMiddleware(c, s.middleware2) {
			return
		}
	} else {
		//404
		s.handle40x(gmcrouter.NewCtx(w, r))
	}

}
func (s *HTTPServer) SetHandler40x(fn func(ctx *gmcrouter.Ctx, tpl *gmctemplate.Template)) *HTTPServer {
	s.handler40x = fn
	return s
}
func (s *HTTPServer) SetHandler50x(fn func(c gmccontroller.IController, err interface{})) *HTTPServer {
	s.handler50x = fn
	return s
}

// called in httpserver
func (s *HTTPServer) handle40x(ctx *gmcrouter.Ctx) {
	if s.handler40x == nil {
		ctx.Response.WriteHeader(http.StatusNotFound)
		ctx.Response.Write([]byte("Page not found"))
	} else {
		s.handler40x(ctx, s.tpl)
	}
	return
}

// called in httprouter
func (s *HTTPServer) handle50x(objv *reflect.Value, err interface{}) {
	c := objv.Interface().(gmccontroller.IController)
	if s.handler50x == nil {
		c.Response__().WriteHeader(http.StatusInternalServerError)
		c.Response__().Header().Set("Content-Type", "text/plain")
		c.Write("Internal Server Error")
		if err != nil && s.config.GetBool("httpserver.showerrorstack") {
			c.Write("\n", err, "\n", string(debug.Stack()))
		}
	} else {
		s.handler50x(c, err)
	}
}

func (s *HTTPServer) SetConfig(c *gmcconfig.GMCConfig) *HTTPServer {
	s.config = c
	return s
}
func (s *HTTPServer) Config() *gmcconfig.GMCConfig {
	return s.config
}

func (s *HTTPServer) ActiveConnCount() int64 {
	return atomic.LoadInt64(s.connCnt)
}
func (s *HTTPServer) Close() *HTTPServer {
	s.server.Close()
	return s
}
func (s *HTTPServer) Listener() net.Listener {
	return s.listener
}
func (s *HTTPServer) Server() *http.Server {
	return s.server
}
func (s *HTTPServer) SetLogger(l *log.Logger) *HTTPServer {
	s.logger = l
	s.server.ErrorLog = s.logger
	return s
}
func (s *HTTPServer) Logger() *log.Logger {
	return s.logger
}
func (s *HTTPServer) SetRouter(r *gmcrouter.HTTPRouter) *HTTPServer {
	s.router = r
	s.router.SetHandle50x(s.handle50x)
	return s
}
func (s *HTTPServer) Router() *gmcrouter.HTTPRouter {
	return s.router
}
func (s *HTTPServer) SetTpl(t *gmctemplate.Template) *HTTPServer {
	s.tpl = t
	return s
}
func (s *HTTPServer) Tpl() *gmctemplate.Template {
	return s.tpl
}
func (s *HTTPServer) SetSessionStore(st gmcsession.Store) *HTTPServer {
	s.sessionStore = st
	return s
}
func (s *HTTPServer) SessionStore() gmcsession.Store {
	return s.sessionStore
}
func (s *HTTPServer) AddMiddleware0(m func(ctx *gmcrouter.Ctx, server *HTTPServer) (isStop bool)) *HTTPServer {
	s.middleware0 = append(s.middleware0, m)
	return s
}
func (s *HTTPServer) AddMiddleware1(m func(ctx *gmcrouter.Ctx, server *HTTPServer) (isStop bool)) *HTTPServer {
	s.middleware1 = append(s.middleware1, m)
	return s
}
func (s *HTTPServer) AddMiddleware2(m func(ctx *gmcrouter.Ctx, server *HTTPServer) (isStop bool)) *HTTPServer {
	s.middleware2 = append(s.middleware2, m)
	return s
}
func (s *HTTPServer) AddMiddleware3(m func(ctx *gmcrouter.Ctx, server *HTTPServer) (isStop bool)) *HTTPServer {
	s.middleware3 = append(s.middleware3, m)
	return s
}

//just for testing
func (s *HTTPServer) bind(addr string) *HTTPServer {
	s.addr = addr
	return s
}
func (s *HTTPServer) createListener() (err error) {
	s.listener, err = net.Listen("tcp", s.addr)
	if err == nil {
		s.addr = s.listener.Addr().String()
	}
	return
}
func (s *HTTPServer) Listen() (err error) {
	err = s.createListener()
	if err != nil {
		return
	}
	go func() {
		for {
			err := s.server.Serve(s.listener)
			if err != nil {
				if !s.isTestNotClosedError && strings.Contains(err.Error(), "closed") {
					s.logger.Printf("http server closed on %s", s.addr)
					s.server.Close()
					break
				} else {
					s.logger.Printf("http server Serve fail on %s , error : %s", s.addr, err)
					time.Sleep(time.Second * 3)
					continue
				}
			}
		}
	}()
	s.logger.Printf("http server listen on >>> %s", s.listener.Addr())
	return
}
func (s *HTTPServer) ListenTLS() (err error) {
	err = s.createListener()
	if err != nil {
		return
	}
	go func() {
		for {
			err := s.server.ServeTLS(s.listener, s.config.GetString("httpserver.tlscert"),
				s.config.GetString("httpserver.tlskey"))
			if err != nil {
				if !s.isTestNotClosedError && strings.Contains(err.Error(), "closed") {
					s.logger.Printf("https server closed.")
					s.server.Close()
					break
				} else {
					s.logger.Printf("http server ServeTLS fail , error : %s", err)
					time.Sleep(time.Second * 3)
					continue
				}
			}
		}
	}()
	s.logger.Printf("https server listen on >>> %s", s.listener.Addr())
	return
}

//ConnState count the active conntions
func (s *HTTPServer) connState(c net.Conn, st http.ConnState) {
	switch st {
	case http.StateNew:
		atomic.AddInt64(s.connCnt, 1)
	case http.StateClosed:
		atomic.AddInt64(s.connCnt, -1)
	}
}

// must be called after router inited
func (s *HTTPServer) initStatic() {
	s.staticDir = s.config.GetString("static.dir")
	s.staticUrlpath = s.config.GetString("static.urlpath")
	if s.staticDir != "" && s.staticUrlpath != "" {
		if strings.HasSuffix(s.staticUrlpath, "/") {
			s.staticUrlpath = strings.TrimRight(s.staticUrlpath, "/")
		}
		s.router.HandlerFunc("GET", s.staticUrlpath+"/*filepath", s.serveStatic)
		s.staticUrlpath += "/"
	}
}
func (s *HTTPServer) initTLSConfig() (err error) {
	if s.config.GetBool("httpserver.tlsenable") {
		tlsCfg := &tls.Config{}
		if s.config.GetBool("httpserver.tlsclientauth") {
			tlsCfg.ClientAuth = tls.RequireAndVerifyClientCert
		}
		clientCertPool := x509.NewCertPool()
		caBytes, e := ioutil.ReadFile(s.config.GetString("httpserver.tlsclientsca"))
		if e != nil {
			return e
		}
		ok := clientCertPool.AppendCertsFromPEM(caBytes)
		if !ok {
			err = errors.New("failed to parse tls clients root certificate")
			return
		}
		tlsCfg.ClientCAs = clientCertPool
		s.server.TLSConfig = tlsCfg
	}
	return
}
func (s *HTTPServer) initSessionStore() (err error) {
	if !s.config.GetBool("session.enable") {
		return
	}
	typ := s.config.GetString("session.store")
	if typ == "" {
		typ = "memory"
	}
	ttl := s.config.GetInt64("session.ttl")
	switch typ {
	case "file":
		cfg := gmcfilestore.NewConfig()
		cfg.TTL = ttl
		cfg.Dir = s.config.GetString("session.file.dir")
		cfg.GCtime = s.config.GetInt("session.file.gctime")
		cfg.Prefix = s.config.GetString("session.file.prefix")
		s.sessionStore, err = gmcfilestore.New(cfg)
	case "memory":
		cfg := gmcmemorystore.NewConfig()
		cfg.TTL = ttl
		cfg.GCtime = s.config.GetInt("session.memory.gctime")
		s.sessionStore, err = gmcmemorystore.New(cfg)
	case "redis":
		cfg := gmcredisstore.NewRedisStoreConfig()
		cfg.RedisCfg.Addr = s.config.GetString("session.redis.address")
		cfg.RedisCfg.Password = s.config.GetString("session.redis.password")
		cfg.RedisCfg.Prefix = s.config.GetString("session.redis.prefix")
		cfg.RedisCfg.Debug = s.config.GetBool("session.redis.debug")
		cfg.RedisCfg.Timeout = time.Second * s.config.GetDuration("session.redis.timeout")
		cfg.RedisCfg.DBNum = s.config.GetInt("session.redis.dbnum")
		cfg.RedisCfg.MaxIdle = s.config.GetInt("session.redis.maxidle")
		cfg.RedisCfg.MaxActive = s.config.GetInt("session.redis.maxactive")
		cfg.RedisCfg.MaxConnLifetime = time.Second * s.config.GetDuration("session.redis.maxconnlifetime")
		cfg.RedisCfg.Wait = s.config.GetBool("session.redis.wait")
		cfg.TTL = ttl
		s.sessionStore, err = gmcredisstore.New(cfg)
	default:
		err = fmt.Errorf("unknown session store type %s", typ)
	}
	return
}

func (s *HTTPServer) serveStatic(w http.ResponseWriter, r *http.Request) {
	pathA := strings.Split(r.URL.Path, "?")
	path := gmcrouter.CleanPath(pathA[0])
	path = strings.TrimPrefix(path, s.staticUrlpath)
	var b []byte
	var ok bool
	//1. find in bindata
	if len(bindata) > 0 {
		b, ok = bindata[path]
	}
	//2. find in system path
	if !ok {
		var e error
		b, e = ioutil.ReadFile(filepath.Join(s.staticDir, path))
		ok = e == nil
	}
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
		return
	}
	cacheSince := time.Now().Format(http.TimeFormat)
	cacheUntil := time.Now().AddDate(66, 0, 0).Format(http.TimeFormat)
	ext := filepath.Ext(path)
	typ := mime.TypeByExtension(ext)
	w.Header().Set("Cache-Control", "max-age:290304000, public")
	w.Header().Set("Last-Modified", cacheSince)
	w.Header().Set("Expires", cacheUntil)
	w.Header().Set("Content-Type", typ)
	gizpCheck := map[string]bool{".js": true, ".css": true}
	if _, ok := gizpCheck[ext]; ok {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			w.Header().Set("Content-Encoding", "gzip")
			gz := gzip.NewWriter(w)
			defer gz.Close()
			gz.Write(b)
			return
		}
	}
	w.Write(b)
}

// PrintRouteTable dump all routes into `w`, if `w` is nil, os.Stdout will be used.
func (s *HTTPServer) PrintRouteTable(w io.Writer) {
	s.router.PrintRouteTable(w)
}

//Start implements service.Services Start
func (s *HTTPServer) Start() (err error) {
	defer func() {
		if err == nil && s.config.GetBool("httpserver.printroute") {
			s.PrintRouteTable(nil)
		}
	}()
	if s.config.GetBool("httpserver.tlsenable") {
		return s.ListenTLS()
	}
	return s.Listen()
}

//Stop implements service.Services Stop
func (s *HTTPServer) Stop() {
	s.Close()
	return
}

//GracefulStop implements service.Services GracefulStop
func (s *HTTPServer) GracefulStop() {
	s.Close()
	return
}

//SetLog implements service.Services SetLog
func (s *HTTPServer) SetLog(l *log.Logger) {
	s.logger = l
	return
}
func (s *HTTPServer) callMiddleware(ctx *gmcrouter.Ctx, middleware []func(ctx *gmcrouter.Ctx, server *HTTPServer) (isStop bool)) (isStop bool) {
	for _, fn := range middleware {
		var isNext bool
		func() {
			defer func() {
				if e := recover(); e != nil {
					s.logger.Printf("middleware pani error : %s", e)
					isNext = true
					isStop = false
				}
			}()
			isStop = fn(ctx, s)
		}()
		if isStop {
			return
		}
		if !isNext {
			break
		}
	}
	return
}
