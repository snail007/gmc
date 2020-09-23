package httpserver

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/snail007/gmc/service"

	gmcconfig "github.com/snail007/gmc/config/gmc"

	"github.com/snail007/gmc/http/controller"
	"github.com/snail007/gmc/http/router"
	"github.com/snail007/gmc/http/session/memorystore"
	"github.com/snail007/gmc/http/template"
	"github.com/snail007/gmc/util/logutil"

	"github.com/stretchr/testify/assert"
)

//go test -coverprofile=a.out;go tool cover -html=a.out

func TestNew(t *testing.T) {
	assert := assert.New(t)
	cfg := mockConfig()
	s := New()
	s.Init(cfg)
	err := s.bind(":").Listen()
	assert.Nil(err)
	addr := s.Listener().Addr().String()
	s0 := New()
	s0.Init(cfg)
	err = s0.bind(addr).Listen()
	assert.NotNil(err)
}
func TestRouting(t *testing.T) {
	assert := assert.New(t)
	s := mockHTTPServer()
	s.BeforeRouting(func(w http.ResponseWriter, r *http.Request, tpl *HTTPServer) (isContinue bool) {
		w.Write([]byte("error"))
		return false
	})
	s.router.HandlerFunc("GET", "/routing", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("routing"))
	})
	w, r := mockRequest(s, "/routing")
	s.ServeHTTP(w, r)
	str, _ := result(w)
	assert.Equal("error", str)
}

func TestRouting_1(t *testing.T) {
	assert := assert.New(t)
	s := mockHTTPServer()
	s.RoutingFiliter(func(w http.ResponseWriter, r *http.Request, ps router.Params, server *HTTPServer) (isContinue bool) {
		w.Write([]byte("error"))
		return false
	})
	s.router.HandlerFunc("GET", "/routing", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("routing"))
	})
	w, r := mockRequest(s, "/routing")
	s.ServeHTTP(w, r)
	str, _ := result(w)
	assert.Equal("error", str)
}
func Test_handle50x(t *testing.T) {
	assert := assert.New(t)
	obj := new(controller.Controller)
	objrf := reflect.ValueOf(obj)
	objv := objrf.Interface().(*controller.Controller)
	w := httptest.NewRecorder()
	objv.Response = w
	objv.Request = httptest.NewRequest("GET", "http://example.com/foo", nil)
	s := New()
	s.Init(gmcconfig.New())
	s.handle50x(&objrf, fmt.Errorf("aaa"))

	//response
	resp := w.Result()
	b, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(resp.StatusCode, int(http.StatusInternalServerError))
	assert.Equal("Internal Server Error", string(b))
}
func Test_handle50x_1(t *testing.T) {
	assert := assert.New(t)
	obj := new(controller.Controller)
	objrf := reflect.ValueOf(obj)
	objv := objrf.Interface().(*controller.Controller)
	w := httptest.NewRecorder()
	objv.Response = w
	objv.Request = httptest.NewRequest("GET", "http://example.com/foo", nil)
	s := New()
	s.Init(gmcconfig.New())
	s.SetHandler50x(func(c *controller.Controller, err interface{}) {
		c.Write(fmt.Errorf("%sbbb", err))
	})
	s.handle50x(&objrf, fmt.Errorf("aaa"))

	//response
	resp := w.Result()
	b, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(resp.StatusCode, int(http.StatusOK))
	assert.Equal("aaabbb", string(b))
}
func TestConnCount(t *testing.T) {
	assert := assert.New(t)
	s := mockHTTPServer()
	s.Listen()
	_, port, _ := net.SplitHostPort(s.listener.Addr().String())
	transport := http.Transport{
		DisableKeepAlives: true,
	}
	client := http.Client{
		Transport: &transport,
	}
	resp, err := client.Get("http://127.0.0.1:" + port)
	client.CloseIdleConnections()
	transport.CloseIdleConnections()
	assert.Equal(int64(0), s.ActiveConnCount())
	assert.Nil(err)
	b, _ := ioutil.ReadAll(resp.Body)
	assert.Equal("Page not found", string(b))
}
func TestConnCount_2(t *testing.T) {
	assert := assert.New(t)
	s := mockHTTPServer()
	s.Listen()
	_, port, _ := net.SplitHostPort(s.listener.Addr().String())
	client := http.Client{}
	client.Get("http://127.0.0.1:" + port)
	assert.Equal(int64(1), s.ActiveConnCount())
}

func TestInit_0(t *testing.T) {
	assert := assert.New(t)
	cfg := gmcconfig.New()
	cfg.SetConfigFile("../../app/app.toml")
	err := cfg.ReadInConfig()
	assert.Nil(err)
	s := New()
	s.Init(cfg)
	err = s.bind(":").Listen()
	assert.Nil(err)
}
func TestInit_1(t *testing.T) {
	assert := assert.New(t)
	cfg := gmcconfig.New()
	cfg.SetConfigFile("../../app/app.toml")
	err := cfg.ReadInConfig()
	assert.Nil(err)
	cfg.Set("session.store", "file")
	s := New()
	s.Init(cfg)
	err = s.bind(":").Listen()
	assert.Nil(err)
}
func TestInit_2(t *testing.T) {
	assert := assert.New(t)
	cfg := gmcconfig.New()
	cfg.SetConfigFile("../../app/app.toml")
	err := cfg.ReadInConfig()
	assert.Nil(err)
	cfg.Set("session.store", "redis")
	s := New()
	s.Init(cfg)
	s.bind(":")
	err = s.Listen()
	assert.Nil(err)
}
func TestInit_3(t *testing.T) {
	assert := assert.New(t)
	cfg := mockConfig()
	cfg.Set("session.store", "none")
	s := New()
	err := s.Init(cfg)
	assert.Equal("unknown session store type none", err.Error())
}
func TestInit_4(t *testing.T) {
	assert := assert.New(t)
	cfg := gmcconfig.New()
	cfg.SetConfigFile("../../app/app.toml")
	err := cfg.ReadInConfig()
	assert.Nil(err)
	cfg.Set("session.store", "")
	s := New()
	s.Init(cfg)
	s.bind(":")
	err = s.Listen()
	assert.Nil(err)
}
func TestHelper(t *testing.T) {
	assert := assert.New(t)
	s := New()
	s.Init(gmcconfig.New())
	err := s.bind(":").Listen()
	assert.Nil(err)
	assert.NotNil(s.Server().Addr)
	s.SetLogger(logutil.New(""))
	assert.NotNil(s.logger)
	r := router.NewHTTPRouter()
	s.SetRouter(r)
	assert.NotNil(s.router)
	tpl, err := template.New("views")
	assert.Nil(err)
	s.SetTpl(tpl)
	assert.NotNil(s.tpl)
	st, err := memorystore.New(memorystore.NewConfig())
	assert.Nil(err)
	s.SetSessionStore(st)
	assert.NotNil(s.sessionStore)
	s.Close()
}

func TestGetter(t *testing.T) {
	assert := assert.New(t)
	s := mockHTTPServer()
	assert.NotNil(s.Config())
	assert.NotNil(s.Logger())
	assert.NotNil(s.Router())
	assert.NotNil(s.Tpl())
	assert.NotNil(s.SessionStore())
}
func TestListenTLS(t *testing.T) {
	assert := assert.New(t)
	s := mockHTTPServer()
	s.SetConfig(mockConfig())
	s.config.Set("httpserver.listen", ":")
	s.config.Set("httpserver.tlscert", "server.crt")
	s.config.Set("httpserver.tlskey", "server.key")
	s.config.Set("httpserver.tlsenable", true)
	s.config.Set("httpserver.tlsclientauth", true)
	s.config.Set("httpserver.tlsclientsca", "server.crt")
	s.initTLSConfig()
	err := s.ListenTLS()
	assert.Nil(err)
	s.Close()
}
func TestListenTLS_0(t *testing.T) {
	assert := assert.New(t)
	s := mockHTTPServer()
	s.SetConfig(mockConfig())
	s.config.Set("httpserver.listen", ":")
	s.config.Set("httpserver.tlscert", "server.crt")
	s.config.Set("httpserver.tlskey", "server.key")
	s.config.Set("httpserver.tlsenable", true)
	s.config.Set("httpserver.tlsclientauth", true)
	s.config.Set("httpserver.tlsclientsca", "server.crt")
	s.initTLSConfig()
	err := s.ListenTLS()
	assert.Nil(err)
	s.isTestNotClosedError = true
	s.Close()
}
func TestListenTLS_1(t *testing.T) {
	assert := assert.New(t)
	s := mockHTTPServer()
	s.SetConfig(mockConfig())
	s.config.Set("httpserver.listen", ":")
	s.config.Set("httpserver.tlscert", "server.crt")
	s.config.Set("httpserver.tlskey", "server.key")
	s.config.Set("httpserver.tlsenable", true)
	s.config.Set("httpserver.tlsclientauth", true)
	s.config.Set("httpserver.tlsclientsca", "none.crt")
	err := s.initTLSConfig()
	assert.NotNil(err)
}
func TestListenTLS_2(t *testing.T) {
	assert := assert.New(t)
	s := mockHTTPServer()
	s.SetConfig(mockConfig())
	s.config.Set("httpserver.listen", ":")
	s.config.Set("httpserver.tlscert", "server.crt")
	s.config.Set("httpserver.tlskey", "server.key")
	s.config.Set("httpserver.tlsenable", true)
	s.config.Set("httpserver.tlsclientauth", true)
	s.config.Set("httpserver.tlsclientsca", "server.key")
	err := s.initTLSConfig()
	assert.NotNil(err)
	s.initBaseObjets()
}
func TestListen(t *testing.T) {
	assert := assert.New(t)
	s := mockHTTPServer()
	err := s.Listen()
	assert.Nil(err)
	s0 := mockHTTPServer()
	s0.addr = s.listener.Addr().String()
	err = s0.Listen()
	assert.NotNil(err)
}
func TestListen_1(t *testing.T) {
	assert := assert.New(t)
	s := mockHTTPServer()
	err := s.Listen()
	assert.Nil(err)
	s.isTestNotClosedError = true
	s.Close()
}
func TestListen_2(t *testing.T) {
	assert := assert.New(t)
	s := mockHTTPServer()
	s.SetConfig(mockConfig())
	s.config.Set("httpserver.listen", ":")
	s.config.Set("httpserver.tlscert", "server.crt")
	s.config.Set("httpserver.tlskey", "server.key")
	s.config.Set("httpserver.tlsenable", true)
	err := s.ListenTLS()
	assert.Nil(err)

	s0 := mockHTTPServer()
	s.SetConfig(mockConfig())
	s.config.Set("httpserver.listen", ":")
	s.config.Set("httpserver.tlscert", "server.crt")
	s.config.Set("httpserver.tlskey", "server.key")
	s.config.Set("httpserver.tlsenable", true)

	s0.addr = s.listener.Addr().String()
	err = s0.ListenTLS()
	assert.NotNil(err)
}

func Test_handler40x(t *testing.T) {
	assert := assert.New(t)
	s := mockHTTPServer()
	s.SetHandler40x(func(w http.ResponseWriter, r *http.Request, tpl *template.Template) {
		w.Write([]byte("404"))
	})
	w, r := mockRequest(s, "/foo")
	s.handle40x(w, r, router.Params{})
	str, _ := result(w)
	assert.Equal("404", str)
}

type User struct {
	controller.Controller
}

func (this *User) URL() {
	a := "a"
	if this.Args != nil {
		a = ""
	}
	this.Write(this.Request.URL.Path + a)
}
func (this *User) Ps() {
	this.Write(this.Args.ByName("args") + this.Args.MatchedRoutePath())
}
func Test_Controller(t *testing.T) {
	assert := assert.New(t)
	s := mockHTTPServer()
	s.router.Controller("/user/", new(User))
	h, _, _ := s.router.Lookup("GET", "/user/url")
	assert.NotNil(h)
	w, r := mockRequest(s, "/user/url")
	s.ServeHTTP(w, r)
	str, _ := result(w)
	assert.Equal("/user/url", str)
}
func Test_Controller_Args(t *testing.T) {
	assert := assert.New(t)
	s := mockHTTPServer()
	s.router.ControllerMethod("/user/:args", new(User), "Ps")
	h, _, _ := s.router.Lookup("GET", "/user/hello")
	assert.NotNil(h)
	w, r := mockRequest(s, "/user/hello")
	s.ServeHTTP(w, r)
	str, _ := result(w)
	assert.Equal("hello/user/:args", str)
}
func Test_SetBindata(t *testing.T) {
	assert := assert.New(t)
	str := base64.StdEncoding.EncodeToString([]byte("{}"))
	SetBinData(map[string]string{
		"js/juqery.js": str,
	})
	assert.NotNil(bindata["js/juqery.js"])
	s := mockHTTPServer()
	w, r := mockRequest(s, "/static/js/juqery.js")
	r.Header.Set("Accept-Encoding", "gzip")
	s.router.HandlerFunc("GET", "/static/", s.serveStatic)
	s.serveStatic(w, r)
	resp := w.Result()
	b, _ := ioutil.ReadAll(resp.Body)
	gz, _ := gzip.NewReader(bytes.NewReader(b))
	buf := make([]byte, 128)
	n, _ := gz.Read(buf)
	assert.Equal("{}", string(buf[:n]))
	bindata = nil
}
func Test_SetBindata_1(t *testing.T) {
	assert := assert.New(t)
	cfg := mockConfig()
	cfg.Set("static.urlpath", "/static")
	s := mockHTTPServer(cfg)
	w, r := mockRequest(s, "/static/js/juqery.js")
	r.Header.Set("Accept-Encoding", "gzip")
	s.serveStatic(w, r)
	resp := w.Result()
	assert.Equal(404, resp.StatusCode)
}

func Test_SetBindata_2(t *testing.T) {
	assert := assert.New(t)
	str := base64.StdEncoding.EncodeToString([]byte("{}"))
	SetBinData(map[string]string{
		"js/juqery.js": str,
	})
	assert.NotNil(bindata["js/juqery.js"])
	s := mockHTTPServer()
	w, r := mockRequest(s, "/static/js/juqery.js")
	s.serveStatic(w, r)
	resp := w.Result()
	b, _ := ioutil.ReadAll(resp.Body)
	assert.Equal("{}", string(b))
	bindata = nil
}
func Test_SetBindata_3(t *testing.T) {
	assert := assert.New(t)
	assert.Panics(func() {
		SetBinData(map[string]string{
			"/js/juqery.js": ".",
		})
	})
}
func Test_SetBindata_4(t *testing.T) {
	assert := assert.New(t)
	cfg := mockConfig()
	cfg.Set("static.dir", "tests")
	cfg.Set("static.urlpath", "/static")
	s := mockHTTPServer(cfg)
	w, r := mockRequest(s, "/static/d.txt")
	r.Header.Set("Accept-Encoding", "gzip")
	s.serveStatic(w, r)
	resp := w.Result()
	b, _ := ioutil.ReadAll(resp.Body)
	assert.Equal("d", string(b))
}
func Test_Service(t *testing.T) {
	assert := assert.New(t)
	s := mockHTTPServer()
	s0 := (service.Service)(s)
	err := s0.Start()
	assert.Nil(err)
	s0.Stop()
	s0.GracefulStop()
}
func Test_Service_2(t *testing.T) {
	assert := assert.New(t)
	cfg := mockConfig()
	cfg.Set("httpserver.tlsenable", true)
	s := mockHTTPServer(cfg)
	s0 := (service.Service)(s)
	err := s0.Start()
	assert.Nil(err)
	s.SetLog(s.logger)
}
func result(w *httptest.ResponseRecorder) (str string, resp *http.Response) {
	resp = w.Result()
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	str = string(body)
	return
}
func mockConfig() *gmcconfig.GMCConfig {
	cfg := gmcconfig.New()
	cfg.SetConfigFile("../../app/app.toml")
	cfg.ReadInConfig()
	return cfg
}
func mockController(obj interface{}, s *HTTPServer, w http.ResponseWriter, r *http.Request) interface{} {
	objv := reflect.ValueOf(obj).Interface().(*controller.Controller)
	objv.Response = w
	objv.Request = r
	objv.Tpl = s.tpl
	objv.SessionStore = s.sessionStore
	objv.Router = s.router
	objv.Config = s.config
	return objv
}
func mockRequest(s *HTTPServer, uri string) (w *httptest.ResponseRecorder, r *http.Request) {
	w = httptest.NewRecorder()
	r = httptest.NewRequest("GET", "http://example.com"+uri, nil)
	return
}

func mockHTTPServer(cfg ...*gmcconfig.GMCConfig) *HTTPServer {
	c := mockConfig()
	if len(cfg) > 0 {
		c = cfg[0]
	}
	s := New()
	s.Init(c)
	s.bind(":")
	s.SetRouter(router.NewHTTPRouter())
	s.SetLogger(logutil.New(""))
	st, _ := memorystore.New(memorystore.NewConfig())
	s.SetSessionStore(st)
	tpl, _ := template.New("../template/tests/views")
	s.SetTpl(tpl)
	return s
}
