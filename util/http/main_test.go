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
	r.HandleFunc("/sleep", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second * 2)
		w.Write([]byte("hello"))
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
MIIDTzCCAjegAwIBAgIBATANBgkqhkiG9w0BAQsFADBJMQswCQYDVQQGEwJUTjES
MBAGA1UEChMJMTI3LjAuMC4xMRIwEAYDVQQLEwkxMjcuMC4wLjExEjAQBgNVBAMT
CTEyNy4wLjAuMTAeFw0yMTAxMDYwMjE1NTRaFw0yMjAxMDYwMzE1NTRaMEkxCzAJ
BgNVBAYTAlROMRIwEAYDVQQKEwkxMjcuMC4wLjExEjAQBgNVBAsTCTEyNy4wLjAu
MTESMBAGA1UEAxMJMTI3LjAuMC4xMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIB
CgKCAQEAp5JEvj67R+SUa35umhvayUU6kvS/AXghu7cmO4WltzvRsYTVdMhXX+Q/
+rbibvjw43kNZZXAbjJRvcVX97d4WwuRMTz7iL0PPcjSuyu2xciHTj99qhkRD5Tk
JKFxluS6VeWa5D8DKSiUCPi3yKXPvZh0Y3H5cw3+ausFPzyNewep8f0MO+N1goKj
htg+Dn7uiHGc0Wj/c3oeED0R/NVriClDvjLc/dw5tZEp0j3APQCrQzQt+MDsrwT2
7lBVDEMEr3YVi2aYs9MbjN7sAuLziYXCyfTH0mfHowb1swUz8/3PpMmMM6ovSwTD
S45fv7JkRw0pFpXFtNqVMsHIrZo8NwIDAQABo0IwQDAOBgNVHQ8BAf8EBAMCAqQw
HQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMA8GA1UdEwEB/wQFMAMBAf8w
DQYJKoZIhvcNAQELBQADggEBAIc4bjtgHwc2fWPNfWIq3xtkWjV8kxUeaf6YsFDk
9X6FjeUiBrFVIO8+D4J7K7hHn5hPq6vwM9oTtNFCX/G39A+xHeTJUJs2kVC6D5Nt
QZ7xmtCIdpXLVuoQOXixDv7vaVzLgNXVWNeQvge95Ajqmsz1K2mmS0/dP4UaQTNm
qo+zxknUzoodUiOf6qVKzahMRvEyfM2sBNz0A/LY3YAwIGFCe0l4uVhzVL9tbczL
gF0LG5WmQEvXY17BLFlRrPYNoBEmupw7/nKrhESyYZl1zRBr+/gMb1t0WQN1hqU0
Cjymrf8tCV73TAn+zS40fDuyC437srgJ4snJ1c3WpQGtqkA=
-----END CERTIFICATE-----`)
	key = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAp5JEvj67R+SUa35umhvayUU6kvS/AXghu7cmO4WltzvRsYTV
dMhXX+Q/+rbibvjw43kNZZXAbjJRvcVX97d4WwuRMTz7iL0PPcjSuyu2xciHTj99
qhkRD5TkJKFxluS6VeWa5D8DKSiUCPi3yKXPvZh0Y3H5cw3+ausFPzyNewep8f0M
O+N1goKjhtg+Dn7uiHGc0Wj/c3oeED0R/NVriClDvjLc/dw5tZEp0j3APQCrQzQt
+MDsrwT27lBVDEMEr3YVi2aYs9MbjN7sAuLziYXCyfTH0mfHowb1swUz8/3PpMmM
M6ovSwTDS45fv7JkRw0pFpXFtNqVMsHIrZo8NwIDAQABAoIBAHim/S4JlujNsPuw
vcviMGZonSMAa6KQL6Gr3jBPKyFCRdpOLS73rMmTW2mWUnTacv8lwrqY10PAoVBF
DfCDPno3WuQb53PtxFKVDP3NHL1Ng/aYCk/12m3go1+oilO9/lgoiJy8rfbti6Dm
C9XBZVE6utp7TsNDmSK8czuzyp6AnHOuAr9AC4XA4kiqU++hSFssGiStPshPlGez
4iOAZk2YEDTzaNGZWse4E5w+xkcgDRcHO57lyQxCK2irU8Mfv9CqQN386xCrKbt6
GiInC9wfAIyqaHa4loIVd5sgrsLs1yeMQAe1667cDHm08jY/G2AVYt851VwAuBly
XdWczaECgYEA2mnNtkt5r1K194HaJanHJWBMQBnhsC8PDtATO6qn7EeIaPGdec8l
Evt2MSPwCUDA121Oa0F7PYhRddvQvEJHTqPi31kMXkSYKhIu47gNDEALxae5GIAh
fBBig7CgLtU+WOPRDvFtev/bRbtFBNvnSqwlm4oFJ2RukArxPyPROVECgYEAxGif
tiHgRetUux1rROsgaqyRVrOX2UE56JmNopFV0uYJNb35umojD/T4RRforC8sknud
V1mb5CqNSmzDxJFu8XRfq1CqEiu06TUuOaZNyNhd5hbUXb89axWUi7MyFUt5wTfi
SUC1HW4hDAlnARtUYyjDuDWv3rT+i0QhuWTsOwcCgYBNDSZVOSskfrlTJ6wdvVdU
CDTeKENGNFPLlfwzAHFdGZ815ob3gexCVhPMIjF8Eiv108nmbKNdgcm7GmD5CSi+
xXIz+OY0G17S+Lcx/qwbtjxw7kqOKiWl7uHSM21PGEt2cGhALUvCKKDiaL5giHOA
FFrwFDDdRMD8b9/LtocJAQKBgQCcvRGXe3lK0v6PRG7yVOFNv+FMW432poLcCI5r
Cah/4WvAI5dDGKhad5gZK3dW0V60l0l9B9nMP9j5Z8ri91yd+8zNHlZaod6BrRry
jrDMcz6b++QF3DPbXSFqStrQ+6Zyd3JyGt1uWxCsVmSJEZJKf6GRQ+bRx4bLBNgU
52FNyQKBgQCakmdRsg+X+pOtlu4dwJgrvpm62fmsWZudv3x6cZ0AiFkdJvyFmYKO
j6iq+Lmi0t1MzXdYWrV6QZUKPQJKvMZLln2m7tbxtoPbqeK1dWBAF7mCSPiD08OY
gnSgL4xe75xshCEvTBlSES993KdqJnqN0+Cie9g/LVxJzc6ywIYRVw==
-----END RSA PRIVATE KEY-----`)
)
