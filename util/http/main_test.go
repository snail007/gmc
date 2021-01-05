package ghttp

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/http"
	"os"
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
	httpsServer = http.Server{
		Handler:   r,
		TLSConfig: tlsConfig,
	}
	l, _ := tls.Listen("tcp", ":",tlsConfig)
	_, p, _ := net.SplitHostPort(l.Addr().String())
	httpsServerURL2 = fmt.Sprintf("https://127.0.0.1:%s", p)
	fmt.Printf("https on %s\n", l.Addr().String())
	go func() { panic(httpsServer2.Serve(l)) }()
}

func initHTTPServer() {
	r := http.NewServeMux()
	r.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})
	r.HandleFunc("/sleep", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
		time.Sleep(time.Second * 2)
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
MIIDSzCCAjOgAwIBAgIBATANBgkqhkiG9w0BAQsFADBGMQswCQYDVQQGEwJBWjER
MA8GA1UEChMIb2RhamcuaW8xETAPBgNVBAsTCG9kYWpnLmlvMREwDwYDVQQDEwhv
ZGFqZy5pbzAgFw0yMTAxMDUxMjAwMTRaGA8xODUxMDMzMDEzNTEwN1owRjELMAkG
A1UEBhMCQVoxETAPBgNVBAoTCG9kYWpnLmlvMREwDwYDVQQLEwhvZGFqZy5pbzER
MA8GA1UEAxMIb2RhamcuaW8wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIB
AQC0mirGO+PEoyDxbAAsCdJujlgUGieHpaPgDrX8tzTt/pu+5H9fR9oH6li+juz7
QzsO+lAOmVVxJyHLe8Vw2Kq8Sw2pIzGBGw1xkXA0JmbAcV7kuRNfuthFQiXTc3FM
OV8GXtFW25AEgqG0NsVlOPFh+DPGUSBVoXoH0eOdILoLdPEAAwS0K/pwg8nYq0IE
nXVlfevnOApglEO4THsjGPGij/EJAy6PS6gdcG1wfMRwMk6MwQM/rDQ1enQA44Hk
Itr4AvY6aWemnsUO0xBzTmlW0KII6o58qcieGP5KhBnZRp7YMmR4MOzVJqvGBESb
64cOEN7UcbV9jqacVUpnr42vAgMBAAGjQjBAMA4GA1UdDwEB/wQEAwICpDAdBgNV
HSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwDwYDVR0TAQH/BAUwAwEB/zANBgkq
hkiG9w0BAQsFAAOCAQEAADvEavt4eBu46DpRxH2tRhHMiYCab16xyOE7VFX9xEcZ
mW5m6ma6vGhCGe1f5AE+YksepDz+iFGgt9x59K7+R131/1naGc7fztjSo8gVoJ7l
odTcBvd3Kn1RcGmLHmsRESVz9afsM45Uldi2OHobuZTCXeY818t0TLT4kHWIWQYv
OE74eTrLUGSoqsgjyLr+AklBUOG2kIXOM28GVZdGK2pLS2Cg7gapmbb6AnWdUVrw
dCPgYQew1E/oDf8xBo9jbZXxknH2k7+6eS3Jh5NQbyIOFxxM6KAZIP81HbYtrwel
0yEnOVEy23J0PL+uRHAkfN4WBi9bRkBiD1gLH2FUDQ==
-----END CERTIFICATE-----`)
	key = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAtJoqxjvjxKMg8WwALAnSbo5YFBonh6Wj4A61/Lc07f6bvuR/
X0faB+pYvo7s+0M7DvpQDplVcSchy3vFcNiqvEsNqSMxgRsNcZFwNCZmwHFe5LkT
X7rYRUIl03NxTDlfBl7RVtuQBIKhtDbFZTjxYfgzxlEgVaF6B9HjnSC6C3TxAAME
tCv6cIPJ2KtCBJ11ZX3r5zgKYJRDuEx7Ixjxoo/xCQMuj0uoHXBtcHzEcDJOjMED
P6w0NXp0AOOB5CLa+AL2Omlnpp7FDtMQc05pVtCiCOqOfKnInhj+SoQZ2Uae2DJk
eDDs1SarxgREm+uHDhDe1HG1fY6mnFVKZ6+NrwIDAQABAoIBAQCljwZTPyenZRuX
9TebJ03ex2J62xcNxIyboyC1kIVW/yZrMjCJeeanhu7fkUoxdAo5ysGFAI9Q8VyL
muT+c8DgZ7UYLgj9n30NdRP02pcrJ0KkTf0yrqf/pYnjc1qTU6yGvHkPNKurVs1B
1UvTZQXyl4Nu84O+vA6QCJqtugiS4pAo+FmKjAn2eAPtfBeX3LMz19uSSY8Htowp
YThDSRaJ8XgnJM7N7nsvdMXdNXwQZlZhWnvKXQoFS2OodoQaeLDsNfyfUIzr5Wk6
o4YNaTGuH8GmDEWcBAg5fPQB+C+mcnrYGnnFXxHPAgedqXpCK/GL65wF7Ki9/E/5
ziHZbt3hAoGBAOIgn+B5q13pQZ6ACDdteLmXYgBw+oBCyhXyS8ycNhvsVY4Wc0ab
MuRNW82y+wHTI8N1081etIBd3WXYGME+pllDAIAitZ/8eZwlAeghcb5mmZLR/baT
iML7Ys2RWQrx1DYnmAH3JkkNr1UnoQQn62hakh75ZY0Zyh0UPBxSH8sZAoGBAMx1
7kPgWOVDfVZv18fwI+B5W/eN6Q+ygUEgt1c7NhKKDqQw920Vwz9yB3XCfJTvSU0K
U7hEAntFJN7Tzq3KfsAngJHwEZvkdiaJEbXr7AdZGdIPR6zP+QsIP9Oukv9MEZCo
IFaOeBPYsugAS5PJBIKVIkiA4z1VrQqif/W6eQAHAoGAF7mMjKS3UhcTB2ovcoFN
1UsIwTsZTTO0uDC/uyv4kV1ubIX2ekX2RPXI2AAbTcm1SuCl5Do3ffBbNkBB+KR2
F49sEgWSQMLgj31igdRgdrWVD05w7CL2il6Nszu4t+k/dp8Y17vyjF+fMbQCtMjr
bftysUVBXliCWCKzW9VR+KECgYBRgdngOTl2+/alVKTC0dqbjAW7pFj6pwCcA/zS
y4n8zgiUL+kTFY/mZQDQUx3zCYlBKxLA7GvI1IGkSu+jnIv28khw5TE/4k2vgwkK
auiG7WA7u1epbqcrXLiFHJ0BJUQDVOK/XsBDuSlpD2URnxsrK2SlXqw4MUVwbeNx
BEtkVQKBgH9YMzzmLNen2iE9jA6UNjsE+1tYQV4g+QPU0bKAhcq1s8aCgTC2QIJH
y98rdCtPB3eNzUMq1YXOdXq35a0ZQduV58MZfgWc/BFQOQKAc9TiViGWoOkMt1rH
KT6Y+LX9uS93FtO1JFHhjP1MbbIUezZiRnp7ppo+tjz/XFbimqD2
-----END RSA PRIVATE KEY-----`)
)
