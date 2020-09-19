package httpserver

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"reflect"
	"strings"
	"sync/atomic"
	"time"

	"github.com/snail007/gmc/http/controller"
	"github.com/snail007/gmc/http/router"
	"github.com/snail007/gmc/http/server/ctxvalue"
	"github.com/snail007/gmc/http/session"
	"github.com/snail007/gmc/http/session/filestore"
	"github.com/snail007/gmc/http/session/memorystore"
	"github.com/snail007/gmc/http/session/redisstore"
	"github.com/snail007/gmc/http/template"
	"github.com/snail007/gmc/util/logutil"
	"github.com/spf13/viper"
)

type HTTPServer struct {
	tpl          *template.Template
	sessionStore session.Store
	router       *router.HTTPRouter
	logger       *log.Logger
	addr         string
	listener     net.Listener
	server       *http.Server
	connCnt      *int64
	config       *viper.Viper
	handler40x   func(w http.ResponseWriter, r *http.Request, tpl *template.Template)
	handler50x   func(c *controller.Controller, err interface{})
	//just for testing
	isTestNotClosedError bool
}

func New(appconfig *viper.Viper) (s *HTTPServer, err error) {
	connCnt := int64(0)
	s = &HTTPServer{
		server:               &http.Server{},
		logger:               logutil.New(""),
		connCnt:              &connCnt,
		config:               appconfig,
		isTestNotClosedError: false,
	}
	s.server.ConnState = s.connState
	s.server.Handler = s
	//init base objects
	err = s.initBaseObjets()
	return
}
func (s *HTTPServer) initBaseObjets() (err error) {
	s.tpl, err = template.New(s.config.GetString("template.dir"))
	if err != nil {
		return
	}
	err = s.initSessionStore()
	if err != nil {
		return
	}
	err = s.initTLSConfig()
	if err != nil {
		return
	}
	s.router = router.NewHTTPRouter()
	s.addr = s.config.GetString("httpserver.listen")
	return
}

func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx = context.WithValue(ctx, ctxvalue.CtxValueKey, ctxvalue.CtxValue{})
	r = r.WithContext(ctx)
	h, args, _ := s.router.Lookup(r.Method, r.URL.Path)
	if h != nil {
		h(w, r, args)
	} else {
		//404
		s.handle40x(w, r, args)
	}
}
func (s *HTTPServer) SetHandler40x(fn func(w http.ResponseWriter, r *http.Request, tpl *template.Template)) *HTTPServer {
	s.handler40x = fn
	return s
}
func (s *HTTPServer) SetHandler50x(fn func(c *controller.Controller, err interface{})) *HTTPServer {
	s.handler50x = fn
	return s
}

func (s *HTTPServer) handle40x(w http.ResponseWriter, r *http.Request, ps router.Params) {
	if s.handler40x == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Page not found"))
	} else {
		s.handler40x(w, r, s.tpl)
	}
	return
}

func (s *HTTPServer) handle50x(objv *reflect.Value, err interface{}) {
	if s.handler50x == nil {
		c := objv.Interface().(*controller.Controller)
		c.Response.WriteHeader(http.StatusInternalServerError)
		c.Write("Internal Server Error")
	} else {
		c := objv.Interface().(*controller.Controller)
		s.handler50x(c, err)
	}
}

func (s *HTTPServer) SetConfig(c *viper.Viper) *HTTPServer {
	s.config = c
	return s
}
func (s *HTTPServer) Config() *viper.Viper {
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
func (s *HTTPServer) SetRouter(r *router.HTTPRouter) *HTTPServer {
	s.router = r
	s.router.SetHandle50x(s.handle50x)
	return s
}
func (s *HTTPServer) Router() *router.HTTPRouter {
	return s.router
}
func (s *HTTPServer) SetTpl(t *template.Template) *HTTPServer {
	s.tpl = t
	return s
}
func (s *HTTPServer) Tpl() *template.Template {
	return s.tpl
}
func (s *HTTPServer) SetSessionStore(st session.Store) *HTTPServer {
	s.sessionStore = st
	return s
}
func (s *HTTPServer) SessionStore() session.Store {
	return s.sessionStore
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
		cfg := filestore.NewConfig()
		cfg.TTL = ttl
		dir := s.config.GetString("session.file.dir")
		if dir == "" {
			dir = os.TempDir()
		}
		cfg.Dir = dir
		cfg.GCtime = s.config.GetInt("session.file.gctime")
		cfg.Prefix = s.config.GetString("session.file.prefix")
		s.sessionStore, err = filestore.New(cfg)
	case "memory":
		cfg := memorystore.NewConfig()
		cfg.TTL = ttl
		cfg.GCtime = s.config.GetInt("session.memory.gctime")
		s.sessionStore, err = memorystore.New(cfg)
	case "redis":
		cfg := redisstore.NewRedisStoreConfig()
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
		s.sessionStore, err = redisstore.New(cfg)
	default:
		err = fmt.Errorf("unknown session store type %s", typ)
	}
	return
}
