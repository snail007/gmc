// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package ghttp

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	gsync "github.com/snail007/gmc/util/sync"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/textproto"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	gproxy "github.com/snail007/gmc/util/proxy"
)

var (
	// Client a HTTPClient, all request shared one http.Client object, keep cookies, keepalive etc.
	Client = NewHTTPClient()
)

func Close() {
	Client.Close()
}

func Get(u string, timeout time.Duration, queryData, header map[string]string) (body []byte, code int, resp *http.Response, err error) {
	return Client.Get(u, timeout, queryData, header)
}

func Post(u string, data map[string]string, timeout time.Duration, header map[string]string) (body []byte, code int, resp *http.Response, err error) {
	return Client.Post(u, data, timeout, header)
}

func NewTriableGet(URL string, maxTry int, timeout time.Duration, queryData, header map[string]string) (tr *TriableRequest, err error) {
	return Client.NewTriableGet(URL, maxTry, timeout, queryData, header)
}

func NewTriablePost(URL string, maxTry int, timeout time.Duration, data, header map[string]string) (tr *TriableRequest, err error) {
	return Client.NewTriablePost(URL, maxTry, timeout, data, header)
}

func NewBatchGet(urlArr []string, timeout time.Duration, data, header map[string]string) (br *BatchRequest, err error) {
	return Client.NewBatchGet(urlArr, timeout, data, header)
}

func NewBatchPost(urlArr []string, timeout time.Duration, data, header map[string]string) (br *BatchRequest, err error) {
	return Client.NewBatchPost(urlArr, timeout, data, header)
}

func PostOfReader(u string, r io.Reader, timeout time.Duration, header map[string]string) (body []byte, code int, resp *http.Response, err error) {
	return Client.PostOfReader(u, r, timeout, header)
}

func Upload(u, fieldName, filename string, data map[string]string, timeout time.Duration, header map[string]string) (body string, resp *http.Response, err error) {
	return Client.Upload(u, fieldName, filename, data, timeout, header)
}

func UploadOfReader(u, fieldName string, filename string, reader io.ReadCloser, data map[string]string, timeout time.Duration, header map[string]string) (body string, resp *http.Response, err error) {
	return Client.UploadOfReader(u, fieldName, filename, reader, data, timeout, header)
}

func Download(u string, timeout time.Duration, queryData, header map[string]string) (data []byte, resp *http.Response, err error) {
	return Client.Download(u, timeout, queryData, header)
}

func DownloadToFile(u string, timeout time.Duration, queryData, header map[string]string, file string) (resp *http.Response, err error) {
	return Client.DownloadToFile(u, timeout, queryData, header, file)
}

func DownloadToWriter(u string, timeout time.Duration, queryData, header map[string]string, writer io.Writer) (resp *http.Response, err error) {
	return Client.DownloadToWriter(u, timeout, queryData, header, writer)
}

// HTTPClient do http get, post etc.
type HTTPClient struct {
	pinCert           *x509.Certificate
	clientCert        tls.Certificate
	clientAuth        bool
	caCertPool        *x509.CertPool
	proxyURL          *url.URL
	opts              *x509.VerifyOptions
	setProxyFromEnv   bool
	dns               []string
	connWrap          func(net.Conn) (conn net.Conn, err error)
	jar               http.CookieJar
	dialer            func(network, address string, timeout time.Duration) (net.Conn, error)
	basicAuthUser     string
	basicAuthPass     string
	httpClientFactory func(r *http.Request) *http.Client
	preHandler        func(r *http.Request)
	keepalive         bool
	proxyUsed         *url.URL
	beforeDo          []BeforeDoClientFunc
	afterDo           []AfterDoClientFunc
}

// NewHTTPClient new a HTTPClient, all request shared one http.Client object, keep cookies, keepalive etc.
// Keepalive is enabled in default.
func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		keepalive: true,
		jar:       NewCookieJar(),
	}
}

// SetBeforeDo sets callback call before request sent.
func (s *HTTPClient) SetBeforeDo(beforeDo BeforeDoClientFunc) *HTTPClient {
	return s.setBeforeDo(beforeDo, true)
}

// AppendBeforeDo add a callback call before request sent.
func (s *HTTPClient) AppendBeforeDo(beforeDo BeforeDoClientFunc) *HTTPClient {
	return s.setBeforeDo(beforeDo, false)
}

func (s *HTTPClient) setBeforeDo(beforeDo BeforeDoClientFunc, isSet bool) *HTTPClient {
	if isSet {
		if beforeDo != nil {
			s.beforeDo = []BeforeDoClientFunc{beforeDo}
		} else {
			s.beforeDo = []BeforeDoClientFunc{}
		}
	} else if beforeDo != nil {
		s.beforeDo = append(s.beforeDo, beforeDo)
	}
	return s
}

// SetAfterDo sets callback call after request sent.
func (s *HTTPClient) SetAfterDo(afterDo AfterDoClientFunc) *HTTPClient {
	return s.setAfterDo(afterDo, true)
}

// AppendAfterDo add a callback call after request sent.
func (s *HTTPClient) AppendAfterDo(afterDo AfterDoClientFunc) *HTTPClient {
	return s.setAfterDo(afterDo, false)
}

func (s *HTTPClient) setAfterDo(afterDo AfterDoClientFunc, isSet bool) *HTTPClient {
	if isSet {
		if afterDo != nil {
			s.afterDo = []AfterDoClientFunc{afterDo}
		} else {
			s.afterDo = []AfterDoClientFunc{}
		}
	} else if afterDo != nil {
		s.afterDo = append(s.afterDo, afterDo)
	}
	return s
}

func (s *HTTPClient) callBeforeDo(req *http.Request) {
	for _, f := range s.beforeDo {
		f(req)
	}
}

func (s *HTTPClient) callAfterDo(req *http.Request, resp *http.Response, err error) {
	for _, f := range s.afterDo {
		f(req, resp, err)
	}
}

// ProxyUsed returns the final proxy URL used by connect to target
func (s *HTTPClient) ProxyUsed() *url.URL {
	return s.proxyUsed
}

func (s *HTTPClient) SetKeepalive(keepalive bool) {
	s.keepalive = keepalive
}

// SetHttpClientFactory sets a http.Client factory, the http.Client it's returned, will be used to send the processed http.Request.
// the argument is the processed http.Request, you can modify it in your any purpose, and the modified  http.Request
// will be sent by the returned client.
func (s *HTTPClient) SetHttpClientFactory(httpClientFactory func(r *http.Request) *http.Client) {
	s.httpClientFactory = httpClientFactory
}

// SetPreHandler handle the processed http.Request, you can modify it in your any purpose, and the modified  http.Request
// will be sent by the returned client.
func (s *HTTPClient) SetPreHandler(preHandler func(r *http.Request)) {
	s.preHandler = preHandler
}

// SetConnWrap called after net.DialTimeout, you can change the conn，return it, if no change just return it.
func (s *HTTPClient) SetConnWrap(fn func(c net.Conn) (conn net.Conn, err error)) {
	s.connWrap = fn
	return
}

// SetProxyFromEnv sets if using http_proxy or all_proxy from system environment to send request.
func (s *HTTPClient) SetProxyFromEnv(set bool) {
	s.setProxyFromEnv = set
	s.proxyUsed = s.getProxyURL()
	return
}

func (s *HTTPClient) SetDialer(dialer func(network, address string, timeout time.Duration) (net.Conn, error)) {
	s.dialer = dialer
	return
}

func (s *HTTPClient) SetBasicAuth(username, password string) {
	s.basicAuthUser = username
	s.basicAuthPass = password
}

// SetProxy sets the proxies used for send request.
//
// if SetProxyFromEnv sets true and SetProxy sets a proxies url, SetProxy will be working.
func (s *HTTPClient) SetProxy(proxyURL string) (err error) {
	if proxyURL == "" {
		return
	}
	if !strings.Contains(proxyURL, "://") {
		proxyURL = "http://" + proxyURL
	}
	s.proxyURL, err = url.Parse(proxyURL)
	if err != nil {
		return
	}
	s.proxyUsed = s.getProxyURL()
	return
}

// SetDNS sets a dns for resolve domain in url when send request.
func (s *HTTPClient) SetDNS(dns ...string) (err error) {
	if len(dns) == 0 {
		return
	}
	for _, v := range dns {
		_, _, err = net.SplitHostPort(v)
		if err != nil {
			return
		}
	}
	s.dns = dns
	return
}

// SetPinCert sets https server's certificate.
func (s *HTTPClient) SetPinCert(pemBytes []byte) (err error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		err = fmt.Errorf("failed to parse certificate PEM")
		return
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

func (s *HTTPClient) setBasicAuth(req *http.Request) {
	if s.basicAuthUser != "" {
		req.SetBasicAuth(s.basicAuthUser, s.basicAuthPass)
	}
}

// NewTriableGet new a triable request with max retry count
func (s *HTTPClient) NewTriableGet(URL string, maxTry int, timeout time.Duration, queryData, header map[string]string) (tr *TriableRequest, err error) {
	return s.newTriableGetPost(http.MethodGet, URL, maxTry, timeout, queryData, header)

}

// NewTriablePost new a triable request with max retry count
func (s *HTTPClient) NewTriablePost(URL string, maxTry int, timeout time.Duration, data, header map[string]string) (tr *TriableRequest, err error) {
	return s.newTriableGetPost(http.MethodPost, URL, maxTry, timeout, data, header)
}

func (s *HTTPClient) newTriableGetPost(method string, URL string, maxTry int, timeout time.Duration, data, header map[string]string) (tr *TriableRequest, err error) {
	tr, err = NewTriableRequestByURL(nil, method, URL, maxTry, timeout, data, header)
	if err != nil {
		return
	}
	tr.DoFunc(func(req *http.Request) (*http.Response, error) {
		return s.Do(req, timeout)
	})
	return tr, nil
}

// NewBatchGet  new a batch get requests
func (s *HTTPClient) NewBatchGet(urlArr []string, timeout time.Duration, data, header map[string]string) (br *BatchRequest, err error) {
	return s.batchGetPost(http.MethodGet, urlArr, timeout, data, header)
}

// NewBatchPost new a batch post requests
func (s *HTTPClient) NewBatchPost(urlArr []string, timeout time.Duration, data, header map[string]string) (br *BatchRequest, err error) {
	return s.batchGetPost(http.MethodPost, urlArr, timeout, data, header)
}

func (s *HTTPClient) batchGetPost(method string, urlArr []string, timeout time.Duration, data, header map[string]string) (br *BatchRequest, err error) {
	br, err = NewBatchURL(nil, method, urlArr, timeout, data, header)
	if err != nil {
		return
	}
	br.DoFunc(func(idx int, req *http.Request) (*http.Response, error) {
		return s.Do(req, timeout)
	})
	return br, nil
}

// Get send a HTTP GET request, no header, just passive nil.
func (s *HTTPClient) Get(u string, timeout time.Duration, queryData, header map[string]string) (body []byte, code int, resp *http.Response, err error) {
	req, cancel, err := NewGet(u, timeout, queryData, header)
	defer cancel()
	if err != nil {
		return
	}
	return s.call(req, timeout)
}

// Post send an HTTP POST request, no header, just passive nil.
// data is form key value.
func (s *HTTPClient) Post(u string, data map[string]string, timeout time.Duration, header map[string]string) (body []byte, code int, resp *http.Response, err error) {
	req, cancel, err := NewPost(u, timeout, data, header)
	defer cancel()
	if err != nil {
		return
	}
	return s.call(req, timeout)
}

// PostOfReader send a HTTP POST request, no header, just passive nil.
// data is an io.Reader.
func (s *HTTPClient) PostOfReader(u string, r io.Reader, timeout time.Duration, header map[string]string) (body []byte, code int, resp *http.Response, err error) {
	ctx, cancel := getTimeoutContext(timeout)
	defer cancel()
	req, err := NewPostReaderWithContext(ctx, u, r, header)
	if err != nil {
		return
	}
	return s.call(req, timeout)
}

func (s *HTTPClient) Do(req *http.Request, timeout time.Duration) (resp *http.Response, err error) {
	s.setBasicAuth(req)
	client, err := s.newClient(req, timeout)
	if err != nil {
		return
	}
	if s.preHandler != nil {
		s.preHandler(req)
	}
	s.callBeforeDo(req)
	resp, err = client.Do(req)
	s.callAfterDo(req, resp, err)
	return
}

func (s *HTTPClient) Close() {
	if v, ok := defaultClient().Transport.(*http.Transport); ok {
		v.CloseIdleConnections()
	}
	return
}

func (s *HTTPClient) call(req *http.Request, timeout time.Duration) (body []byte, code int, resp *http.Response, err error) {
	defer CloseResponse(resp)
	resp, err = s.Do(req, timeout)
	if err != nil {
		return
	}
	body, err = GetResponseBodyE(resp)
	if err != nil {
		return
	}
	code = resp.StatusCode
	return
}

// Upload uploads a file from a file `filename`,
// fieldName is the form filed name in form, filename is the value of filed `fieldName` and also the file to upload.
// data is the additional form data.
func (s *HTTPClient) Upload(u, fieldName, filename string, data map[string]string, timeout time.Duration, header map[string]string) (body string, resp *http.Response, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	return s.UploadOfReader(u, fieldName, filename, file, data, timeout, header)
}

// UploadOfReader upload a file from a io.Reader `reader`,
// fieldName is the form filed name in form, filename is the value of filed `fieldName`.
// data is the additional form data.
func (s *HTTPClient) UploadOfReader(u, fieldName string, filename string, reader io.ReadCloser, data map[string]string, timeout time.Duration, header map[string]string) (body string, resp *http.Response, err error) {
	r, w := io.Pipe()
	m := multipart.NewWriter(w)
	go func() {
		defer w.Close()
		defer m.Close()
		defer reader.Close()
		if len(data) > 0 {
			for k, v := range data {
				m.WriteField(k, v)
			}
		}
		var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition",
			fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
				quoteEscaper.Replace(fieldName), quoteEscaper.Replace(filename)))
		h.Set("Content-Type", "application/octet-stream")
		part, err := m.CreatePart(h)
		if err != nil {
			return
		}
		if _, err = io.Copy(part, reader); err != nil {
			return
		}
	}()
	ctx, cancel := getTimeoutContext(timeout)
	defer cancel()
	if header == nil {
		header = map[string]string{}
	}
	header["Content-Type"] = m.FormDataContentType()
	req, err := NewPostReaderWithContext(ctx, u, r, header)
	if err != nil {
		return
	}
	b, _, resp, err := s.call(req, timeout)
	if err == nil {
		body = string(b)
	}
	return
}

// Download gets url bytes contents.
func (s *HTTPClient) Download(u string, timeout time.Duration, queryData, header map[string]string) (data []byte, resp *http.Response, err error) {
	data, _, resp, err = s.Get(u, timeout, queryData, header)
	return
}

// DownloadToFile gets url bytes contents and save to the file.
func (s *HTTPClient) DownloadToFile(u string, timeout time.Duration, queryData, header map[string]string, file string) (resp *http.Response, err error) {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return
	}
	defer f.Close()
	return s.DownloadToWriter(u, timeout, queryData, header, f)
}

// DownloadToWriter gets url bytes contents and copy to the writer at realtime.
func (s *HTTPClient) DownloadToWriter(u string, timeout time.Duration, queryData, header map[string]string, writer io.Writer) (resp *http.Response, err error) {
	req, cancel, err := NewGet(u, timeout, queryData, header)
	defer cancel()
	if err != nil {
		return
	}
	resp, err = s.Do(req, timeout)
	defer CloseResponse(resp)
	if err != nil {
		return
	}
	_, err = io.Copy(writer, resp.Body)
	return
}

func (s *HTTPClient) newTransport(timeout time.Duration) (tr *http.Transport, err error) {
	resolver := s.newResolver(timeout)
	tr = &http.Transport{
		ResponseHeaderTimeout: timeout,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   timeout,
		ExpectContinueTimeout: 1 * time.Second,
		DisableKeepAlives:     !s.keepalive,
		Proxy: func(req *http.Request) (u *url.URL, err error) {
			// do nothing, disable default, process in DialContext
			return
		},
		DialContext: func(c context.Context, network, addr string) (conn net.Conn, err error) {
			ctx, cancel := context.WithTimeout(c, timeout)
			defer cancel()
			if s.proxyUsed != nil {
				var j *gproxy.Jumper
				j, err = gproxy.NewJumper(s.proxyUsed.String(), timeout)
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
					return nil, fmt.Errorf("can not resolve domain %s", host)
				}
				ip := iparr[0]
				addr = net.JoinHostPort(ip.String(), port)
				if s.dialer != nil {
					conn, err = s.dialer(network, addr, timeout)
				} else {
					conn, err = net.DialTimeout(network, addr, timeout)
				}
			}
			if err == nil && s.connWrap != nil {
				conn, err = s.connWrap(conn)
			}
			return
		},
	}
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
				return fmt.Errorf("unknown server certificate received")
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
					return fmt.Errorf("unknown server certificate received")
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
	return
}

func (s *HTTPClient) newResolver(timeout time.Duration) (resolver *net.Resolver) {
	resolver = net.DefaultResolver
	if len(s.dns) > 0 {
		resolver = &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{
					Timeout:   timeout,
					KeepAlive: timeout,
				}
				dnsLen := len(s.dns)
				errsChn := make(chan error, dnsLen)
				connChn := make(chan net.Conn)
				g := sync.WaitGroup{}
				g.Add(dnsLen)
				for _, v := range s.dns {
					go func(addr string) {
						defer g.Done()
						c, e := d.Dial("udp", addr)
						if e != nil {
							errsChn <- e
							return
						}
						select {
						case connChn <- c:
						default:
							c.Close()
						}
					}(v)
				}
				select {
				case c := <-connChn:
					return c, nil
				case <-gsync.WaitGroupToChan(&g):
					//case selected randomly, so here need check double if it is success
					if len(connChn) > 0 {
						c := <-connChn
						return c, nil
					}
					return nil, <-errsChn
				}
			},
		}
	}
	return
}
func (s *HTTPClient) getProxyURL() (proxyURL *url.URL) {
	if s.proxyURL != nil {
		return s.proxyURL
	}
	if s.setProxyFromEnv {
		proxyENV := ""
		for _, k := range []string{"http_proxy", "https_proxy", "all_proxy"} {
			proxyENV = os.Getenv(k)
			if proxyENV == "" {
				proxyENV = os.Getenv(strings.ToUpper(k))
			}
			if proxyENV != "" {
				break
			}
		}
		if proxyENV != "" {
			if !strings.Contains(proxyENV, "://") {
				proxyENV = "http://" + proxyENV
			}
			proxyURL, _ = url.Parse(proxyENV)
		}
	}
	return
}

// new a http.Client
func (s *HTTPClient) newClient(req *http.Request, timeout time.Duration) (client *http.Client, err error) {
	if s.httpClientFactory != nil {
		client = s.httpClientFactory(req)
	} else {
		client = defaultClient()
	}
	client.Jar = s.jar
	client.Transport, err = s.newTransport(timeout)
	return
}

type CookieJar struct {
	jar *cookiejar.Jar
}

func NewCookieJar() *CookieJar {
	jar, _ := cookiejar.New(&cookiejar.Options{})
	return &CookieJar{
		jar: jar,
	}
}

func (s *CookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	for k, v := range cookies {
		if net.ParseIP(v.Domain) != nil {
			v.Domain = ""
		}
		cookies[k] = v
	}
	s.jar.SetCookies(u, cookies)
}

func (s *CookieJar) Cookies(u *url.URL) []*http.Cookie {
	return s.jar.Cookies(u)
}

var defaultHTTPClient = &http.Client{
	Transport: http.DefaultTransport,
}

func defaultClient() *http.Client {
	return defaultHTTPClient
}
