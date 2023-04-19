package ghttp

import (
	"crypto/md5"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	gctx "github.com/snail007/gmc/module/ctx"
	gcast "github.com/snail007/gmc/util/cast"
	"io"
	"net"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var (
	httpServer      http.Server
	httpsServer     http.Server
	httpsServer2    http.Server
	httpServerURL   string
	httpsServerURL  string
	httpsServerURL2 string
)

func TestMain(m *testing.M) {
	initHTTPServer()
	initHTTPSServer()
	initHTTPSServer2()
	os.Exit(m.Run())
}

func initHTTPSServer() {
	r := http.NewServeMux()
	r.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})
	tlsConfig := &tls.Config{}
	tlsConfig.Certificates = make([]tls.Certificate, 1)
	tlsConfig.Certificates[0], _ = tls.X509KeyPair(cert, key)
	httpsServer = http.Server{
		Handler: r,
	}
	l, _ := tls.Listen("tcp", ":", tlsConfig)
	_, p, _ := net.SplitHostPort(l.Addr().String())
	httpsServerURL = fmt.Sprintf("https://127.0.0.1:%s", p)
	fmt.Printf("https on %s\n", l.Addr().String())
	go func() { panic(httpsServer.Serve(l)) }()
}

func initHTTPSServer2() {
	r := http.NewServeMux()
	r.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})
	clientCertPool := x509.NewCertPool()
	clientCertPool.AppendCertsFromPEM(cert)
	tlsConfig := &tls.Config{
		ClientCAs:  clientCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	tlsConfig.Certificates = make([]tls.Certificate, 1)
	tlsConfig.Certificates[0], _ = tls.X509KeyPair(cert, key)
	httpsServer2 = http.Server{
		Handler:   r,
		TLSConfig: tlsConfig,
	}
	l, _ := tls.Listen("tcp", ":", tlsConfig)
	_, p, _ := net.SplitHostPort(l.Addr().String())
	httpsServerURL2 = fmt.Sprintf("https://127.0.0.1:%s", p)
	fmt.Printf("https on %s\n", l.Addr().String())
	go func() { panic(httpsServer2.Serve(l)) }()
}

func initHTTPServer() {
	r := http.NewServeMux()
	r.HandleFunc("/header", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.Header.Get("token")))
	})
	r.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.FormValue("name")))
	})
	r.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})
	r.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		file, _, err := r.FormFile("test")
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		f, _ := os.Create("test.bin")
		io.Copy(f, file)
		f.Close()
		h := md5.New()
		h.Write([]byte("a"))
		s := fmt.Sprintf("%x", h.Sum(nil))
		w.Write([]byte(r.PostFormValue("uid") + s))
	})
	r.HandleFunc("/sleep", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second * 2)
		w.Write([]byte("hello"))
	})
	cntMap := sync.Map{}
	r.HandleFunc("/try1", func(w http.ResponseWriter, r *http.Request) {
		ctx := gctx.NewCtxWithHTTP(w, r)
		cntKey := ctx.GET("key")
		v, _ := cntMap.LoadOrStore(cntKey, new(int32))
		cnt := v.(*int32)
		if atomic.AddInt32(cnt, 1) <= 3 {
			time.Sleep(time.Second * 3)
		}

		ctx.Write(ctx.GET("msg") + ctx.POST("msg") + ctx.Header("h1"))
	})
	r.HandleFunc("/batch", func(w http.ResponseWriter, r *http.Request) {
		ctx := gctx.NewCtxWithHTTP(w, r)
		idx := ctx.GET("idx")
		if idx == "" {
			idx = ctx.POST("idx")
		}
		sleepStr := ctx.GET("sleep")
		if sleepStr == "" {
			sleepStr = ctx.POST("sleep")
		}
		if ctx.GET("nosleep") == "1" {
			sleepStr = ""
		}
		sleep := gcast.ToInt(sleepStr)
		if sleep > 0 {
			time.Sleep(time.Second * time.Duration(sleep))
		}
		ctx.Write(idx)
	})
	httpServer = http.Server{
		Handler: r,
	}
	l, _ := net.Listen("tcp", ":")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	httpServerURL = fmt.Sprintf("http://127.0.0.1:%s", p)
	fmt.Printf("http on %s\n", l.Addr().String())
	go func() { panic(httpServer.Serve(l)) }()
}

var (
	cert = []byte(`-----BEGIN CERTIFICATE-----
MIIDUjCCAjqgAwIBAgIBATANBgkqhkiG9w0BAQsFADA6MQswCQYDVQQGEwJMSzEN
MAsGA1UEChMEYi5tdDENMAsGA1UECxMEYi5tdDENMAsGA1UEAxMEYi5tdDAgFw0y
MjAxMTMwMzIzNDJaGA8yMTIxMTIyMDA0MjM0MlowOjELMAkGA1UEBhMCTEsxDTAL
BgNVBAoTBGIubXQxDTALBgNVBAsTBGIubXQxDTALBgNVBAMTBGIubXQwggEiMA0G
CSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDB6BN0Ffimon5bm8IFMmavXs4sURYI
Twp7JMWzKoagfpTqY4ckCCVK0dEc8NsxuY/goCvJQB4NmAgnPFE3yDQh6jAKDNRF
98P0Kt7hIJVpQTJe67PHrp/lWB1h1lJ1eybA5EVi8w/rDTdGSrHMPinEZGGHyLFP
x8aMb5qfSQepbV4j/8Z+4j05YFtDb0tA2Ib9JYL/ldoFrsNETvuy1reymcVbt2na
KoViWGV+EWNSE3uMHiyaby5BaLEDiOCgZfnc5dGm+/SipPdVFhRImtBZ8aBphb3i
cb75GoffsTUWqFxgM5b+svNbwQe5gHpvLz+cxBpSMB1vjbBUfqsrppcfAgMBAAGj
YTBfMA4GA1UdDwEB/wQEAwICpDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUH
AwIwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQUUXK5xVPLkexk5zLEvLHE1s4X
7+wwDQYJKoZIhvcNAQELBQADggEBACFggktQigP81ZZNfDA6X3XjB5ftue/RTq+y
qvzNHGtSwZThTvMiq8/e2XiHsQ4yw37sWqnyTAzvsVQnMQSUdkwKrqM6tRXR2iyo
M8XZxJeC2DYVk5RdBH2VHD5ZkBoMmkwC745iwLLm2+guFcGRfN7Udeq59CXbyZBa
z+gdC7ihpGd5yc5K8eC99MUKJa7Cf9IfOBi0T6TTjERHh+5FG57UvvYPgWjiCdrI
/WI4siAdUZF3cU5Fz1dAdjJbfy3+BgnJ/pnGpmRTmfvJYuusQ2G1wBwkSfuFiUYz
unMKFUDdiHbKlJc7esLCKuHvv+9S8GZ+DUROQcGmlHAk90LJxns=
-----END CERTIFICATE-----`)
	key = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAwegTdBX4pqJ+W5vCBTJmr17OLFEWCE8KeyTFsyqGoH6U6mOH
JAglStHRHPDbMbmP4KAryUAeDZgIJzxRN8g0IeowCgzURffD9Cre4SCVaUEyXuuz
x66f5VgdYdZSdXsmwORFYvMP6w03RkqxzD4pxGRhh8ixT8fGjG+an0kHqW1eI//G
fuI9OWBbQ29LQNiG/SWC/5XaBa7DRE77sta3spnFW7dp2iqFYlhlfhFjUhN7jB4s
mm8uQWixA4jgoGX53OXRpvv0oqT3VRYUSJrQWfGgaYW94nG++RqH37E1FqhcYDOW
/rLzW8EHuYB6by8/nMQaUjAdb42wVH6rK6aXHwIDAQABAoIBAFvVS5FI7pAMmQdN
xx+q9RLNNZurc2HP/UjA65ik5UmRaXlwVYptCSxcHks0jrsIBOn/MX2IjjMl84aR
fG2MtZPyU3oPBWF1rCoxO4knY1uL8w0dV/GT9Eor5w508GyPPJVSBsKMFtfdVHZH
3di7ABDw8XfbRo5gMNpF6NbTQXutaVrDrwdc6R3HRJ8QUTCOe+OryH73nwa2K+iz
EiVd4bFKnYuwXxwBtjWKJyMxal5OsDm+QGg7cZJmA4O+lM4Kw8F8XiwJonPEQEYI
7vSX0zNVnZ9T0CtSRWPhrQTwvv5MLw5/OrdDK/XDgsEVrqC2zTJO8iHihpu6VEJt
pf8jAgECgYEAwr7+laRjYwv40o1lqLu79AYlAIPl4G82ivQGlgbbvw4gNve1A6EA
Ihr3jD9TOV88wh/FAjSawTpl+qxQC65ne1VakO/5Qb9JNnKgR3Q6H7QecrlMGd/v
D+FXdO1fkcCRyLedF57q9N7AzvRxLHK1R0KZsfqslUlOlHq8WsFWuIECgYEA/uV7
k1/7HKqr0uP0WmBLohtYXZNbqhdb14ZhnMfw0UpBlCUtWRseqE3K0A/weuE6mkaA
tH9d3xTWa+WMQ3quyaZXPnfVMvNnM4dgVK1kNJKNC+Wgh6Ci+iGmmDp5U+ZeQcAw
8FFMeqlphwFEKO7jdlKCqDUIK96Du1wcR4BCf58CgYAMh++nv1kp0WZkXfbRoarZ
a9/LpbEP/Pf8fvFBjBVtuMH35359SknQ5/1Px+9Z/LfTIeoyVyIyFsjjFV1dMw6z
j+1w8BAQ2/chCsUnc+IdkiB3b1bnP1KJqg1Pl8qTfVmkGbSBBZfGw+KSLoZtvr/N
YwqyuheKz5m/0hn2mQQ0gQKBgF0Ofp0BL3X5wR0O58iO203lWc9f2tkwCfGXN8+7
FunxiBuDrxiW1AxxyhdHmm3iCDkGgDplPWoR+24MsbZ49ZLczYEa0pT1U7n2NG71
ll2zGxc6z+5z8MwMuPtebaj5s3OhrLwvkhI+Ay6sgavH+vbZjKXIJqGNbN5b9F8O
LjjVAoGBAK0yUx5tnKtvPnUKOHW7u0rU3xYv/6aCIsMmtqldVbUXMAhJnJOWCkWX
wg5jVZUp70wOhKa7NuRbN0RnIaGbTsFSNR/CGk8bpeWubm2D7k12OVAu6RJTLgr7
TI6ovh/ClCYQDP5sqZfbiNb6Kks01TSO2jVDhuQ3VJBT6ps7mtq1
-----END RSA PRIVATE KEY-----`)
)
