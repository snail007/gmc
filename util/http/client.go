package ghttp

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	gproxy "github.com/snail007/gmc/util/proxy"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	// Client a HTTPClient, all request shared one http.Client object, keep cookies, keepalive etc.
	Client = NewHTTPClient()
	// SimpleClient simple HTTPClient, every request using a new http.Client object.
	SimpleClient = NewSimpleHTTPClient()
)

// HTTPClient do http get, post etc.
type HTTPClient struct {
	pinCert         *x509.Certificate
	clientCert      tls.Certificate
	clientAuth      bool
	caCertPool      *x509.CertPool
	proxyURL        *url.URL
	opts            *x509.VerifyOptions
	setProxyFromEnv bool
	dns             string
	connWrap        func(net.Conn) (conn net.Conn, err error)
	simpleClient    bool
	client          *http.Client
	once            *sync.Once
}

// NewSimpleHTTPClient new a simple HTTPClient, every request using a new http.Client object.
func NewSimpleHTTPClient() HTTPClient {
	return HTTPClient{
		once:         &sync.Once{},
		simpleClient: true,
	}
}

// NewHTTPClient new a HTTPClient, all request shared one http.Client object, keep cookies, keepalive etc.
func NewHTTPClient() HTTPClient {
	return HTTPClient{
		once: &sync.Once{},
	}
}

// SetConnWrap called after net.DialTimeout, you can change the conn，return it, if no change just return it.
func (s *HTTPClient) SetConnWrap(fn func(c net.Conn) (conn net.Conn, err error)) {
	s.connWrap = fn
	return
}

// SetProxyFromEnv sets if using http_proxy or all_proxy from system environment to send request.
func (s *HTTPClient) SetProxyFromEnv(set bool) {
	s.setProxyFromEnv = set
	return
}

// SetProxy sets the proxies used for send request.
//
// if SetProxyFromEnv sets true and SetProxy sets a proxies url, SetProxy will be working.
func (s *HTTPClient) SetProxy(proxyURL string) (err error) {
	if proxyURL == "" {
		return
	}
	s.proxyURL, err = url.Parse(proxyURL)
	return
}

// SetDNS sets a dns for resolve domain in url when send request.
func (s *HTTPClient) SetDNS(dns string) (err error) {
	if dns == "" {
		return
	}
	_, _, err = net.SplitHostPort(dns)
	if err != nil {
		return
	}
	s.dns = dns
	return
}

// SetPinCert sets https server's certificate.
func (s *HTTPClient) SetPinCert(pemBytes []byte) (err error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		panic("failed to parse certificate PEM")
	}
	s.pinCert, err = x509.ParseCertificate(block.Bytes)
	return
}

// SetClientCert sets client certificate send to https server, when https server requires client certificate.
func (s *HTTPClient) SetClientCert(certPemBytes, keyBytes []byte) (err error) {
	s.clientCert, err = tls.X509KeyPair(certPemBytes, keyBytes)
	if err != nil {
		return
	}
	s.clientAuth = true
	return
}

// SetRootCaCerts sets root certificates to checking  https server's certificate.
func (s *HTTPClient) SetRootCaCerts(caPemBytes ...[]byte) (err error) {
	if s.caCertPool == nil {
		s.caCertPool = x509.NewCertPool()
	}
	for _, v := range caPemBytes {
		ok := s.caCertPool.AppendCertsFromPEM(v)
		if !ok {
			err = errors.New("failed to parse root certificate")
			return
		}
	}
	s.opts = &x509.VerifyOptions{
		Roots: s.caCertPool,
	}
	return
}

// Get send a HTTP GET request, no header, just passive nil.
func (s *HTTPClient) Get(u string, timeout time.Duration, header map[string]string) (body []byte, code int, resp *http.Response, err error) {
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()
	client, err := s.getClient(u, timeout)
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return
	}
	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}
	resp, err = client.Do(req)
	if err != nil {
		return nil, 0, nil, err
	}
	code = resp.StatusCode
	body, err = ioutil.ReadAll(resp.Body)
	return
}

// Post send a HTTP POST request, no header, just passive nil.
// data is form key value.
func (s *HTTPClient) Post(u string, data map[string]string, timeout time.Duration, header map[string]string) (body []byte, code int, resp *http.Response, err error) {
	postParamsString := ""
	if data != nil {
		postParams := []string{}
		for k, v := range data {
			postParams = append(postParams, url.QueryEscape(k)+"="+url.QueryEscape(v))
		}
		postParamsString = strings.Join(postParams, "&")
	}
	return s.PostOfReader(u, strings.NewReader(postParamsString), timeout, header)
}

// PostOfReader send a HTTP POST request, no header, just passive nil.
// data is an io.Reader.
func (s *HTTPClient) PostOfReader(u string, r io.Reader, timeout time.Duration, header map[string]string) (body []byte, code int, resp *http.Response, err error) {
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "POST", u, r)
	if err != nil {
		return
	}
	client, err := s.getClient(u, timeout)
	if err != nil {
		return
	}
	foundCT := false
	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v)
			if strings.TrimSpace(strings.ToLower(k)) == "content-type" {
				foundCT = true
			}
		}
	}
	if !foundCT {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	resp, err = client.Do(req)
	if err != nil {
		return
	}
	body, err = ioutil.ReadAll(resp.Body)
	code = resp.StatusCode
	if err != nil {
		return
	}
	return
}

// Close close the idle connections when keepalive enabled.
//
// Close do nothing in simple HTTPClient
func (s *HTTPClient) Close() {
	if s.client != nil {
		s.client.Transport.(*http.Transport).CloseIdleConnections()
	}
}

// get a client for requests
func (s *HTTPClient) getClient(u string, timeout time.Duration) (client *http.Client, err error) {
	if s.client != nil {
		return s.client, nil
	}
	defer func() {
		if !s.simpleClient && err == nil {
			s.client = client
		}
	}()
	if s.simpleClient {
		client, err = s.newClient(u, timeout)
	} else {
		s.once.Do(func() {
			s.client, err = s.newClient(u, timeout)
		})
		client = s.client
	}
	return
}

// new a http.Client
func (s *HTTPClient) newClient(u string, timeout time.Duration) (client *http.Client, err error) {
	tr := &http.Transport{}
	if strings.Contains(u, "https://") {
		conf := &tls.Config{}
		conf.InsecureSkipVerify = true
		conf.VerifyPeerCertificate = func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			if s.pinCert == nil && s.caCertPool != nil {
				found := false
				for _, rawCert := range rawCerts {
					cert, _ := x509.ParseCertificate(rawCert)
					_, e := cert.Verify(*s.opts)
					if e == nil {
						found = true
						break
					}
				}
				if !found {
					return fmt.Errorf("unknown server certificate recieved")
				}
			} else {
				if s.pinCert != nil {
					found := false
					for _, rawCert := range rawCerts {
						cert, _ := x509.ParseCertificate(rawCert)
						if s.pinCert.Equal(cert) {
							found = true
						}
					}
					if !found {
						return fmt.Errorf("unknown server certificate recieved")
					}
				}
			}
			return nil
		}
		if s.caCertPool != nil {
			conf.RootCAs = s.caCertPool
		}
		if s.clientAuth {
			conf.Certificates = []tls.Certificate{s.clientCert}
		}
		tr.TLSClientConfig = conf
	}
	resolver := net.DefaultResolver
	if s.dns != "" {
		resolver = &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{
					Timeout:   timeout,
					KeepAlive: timeout,
				}
				return d.DialContext(ctx, "udp", s.dns)
			},
		}
	}
	//proxy
	tr.Proxy = func(req *http.Request) (u *url.URL, err error) {
		// do nothing, disable default, process in DialContext
		return
	}
	if s.proxyURL == nil && s.setProxyFromEnv {
		proxyENV := ""
		for _, k := range []string{"http_proxy", "all_proxy"} {
			proxyENV = os.Getenv(k)
			if proxyENV == "" {
				proxyENV = os.Getenv(strings.ToUpper(k))
			}
			if proxyENV != "" {
				break
			}
		}
		if proxyENV != "" {
			if !strings.HasPrefix(proxyENV, "http://") {
				proxyENV = "http://" + proxyENV
			}
			s.proxyURL, _ = url.Parse(proxyENV)
		}
	}
	tr.DialContext = func(c context.Context, network, addr string) (conn net.Conn, err error) {
		ctx, cancel := context.WithTimeout(c, timeout)
		defer cancel()
		if s.proxyURL != nil {
			var j gproxy.Jumper
			j, err = gproxy.NewJumper(s.proxyURL.String(), timeout)
			if err != nil {
				return
			}
			conn, err = j.Dial(addr)
		} else {
			host, port, _ := net.SplitHostPort(addr)
			iparr, e := resolver.LookupIPAddr(ctx, host)
			if e != nil {
				return nil, e
			}
			if len(iparr) == 0 {
				return nil, fmt.Errorf("can not reslove domain %s", host)
			}
			ip := iparr[rand.Intn(len(iparr))]
			addr = net.JoinHostPort(ip.String(), port)
			conn, err = net.DialTimeout(network, addr, timeout)
		}
		if err == nil && s.connWrap != nil {
			conn, err = s.connWrap(conn)
		}
		return
	}
	client = &http.Client{
		Transport: tr,
	}
	if s.simpleClient {
		tr.DisableKeepAlives = true
	} else {
		client.Jar, _ = cookiejar.New(&cookiejar.Options{})
	}
	return
}
