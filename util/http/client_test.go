// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package ghttp

import (
	"bytes"
	"crypto/md5"
	"fmt"
	assert2 "github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestHTTPClient_Header(t *testing.T) {
	assert := assert2.New(t)
	client := NewHTTPClient()
	body, _, _, _ := client.Post(httpServerURL+"/header", map[string]string{"name": "snail007"}, time.Second, map[string]string{"token": "200"})
	assert.Equal("200", string(body))
}
func TestHTTPClient_Post(t *testing.T) {
	assert := assert2.New(t)
	client := NewHTTPClient()
	body, _, _, _ := client.Post(httpServerURL+"/post", map[string]string{"name": "snail007"}, time.Second, map[string]string{"token": "200"})
	assert.Equal("snail007", string(body))

	body, _, _, _ = Post(httpServerURL+"/post", map[string]string{"name": "snail007"}, time.Second, map[string]string{"token": "200"})
	assert.Equal("snail007", string(body))
}

func TestHTTPClient_PostOfReader(t *testing.T) {
	assert := assert2.New(t)
	client := NewHTTPClient()
	body, _, _, _ := client.PostOfReader(httpServerURL+"/post", strings.NewReader("name=snail007"), time.Second, map[string]string{"token": "200"})
	assert.Equal("snail007", string(body))

	body, _, _, _ = PostOfReader(httpServerURL+"/post", strings.NewReader("name=snail007"), time.Second, map[string]string{"token": "200"})
	assert.Equal("snail007", string(body))
}

func TestHTTPClient_Upload(t *testing.T) {
	assert := assert2.New(t)
	f, _ := os.Create("upload.bin")
	f.WriteString("a")
	f.Close()
	h := md5.New()
	h.Write([]byte("a"))
	s := fmt.Sprintf("%x", h.Sum(nil))
	body, _, err := Upload(httpServerURL+"/upload", "test", "upload.bin", map[string]string{"uid": "007"}, 0, nil)
	assert.Nil(err)
	assert.Equal("007"+s, body)
	assert.FileExists("test.bin")
	os.Remove("test.bin")
	os.Remove("upload.bin")
}

func TestHTTPClient_UploadOfReader(t *testing.T) {
	assert := assert2.New(t)
	f, _ := os.Create("upload.bin")
	f.WriteString("a")
	f.Seek(0, 0)
	defer f.Close()
	h := md5.New()
	h.Write([]byte("a"))
	s := fmt.Sprintf("%x", h.Sum(nil))
	body, _, err := UploadOfReader(httpServerURL+"/upload", "test", "upload.bin", f, map[string]string{"uid": "007"}, 0, map[string]string{"test": "test"})
	assert.Nil(err)
	assert.Equal("007"+s, body)
	assert.FileExists("test.bin")
	os.Remove("test.bin")
	os.Remove("upload.bin")
}

func TestHTTPClient_Get(t *testing.T) {
	assert := assert2.New(t)
	client := NewHTTPClient()
	body, _, _, err := client.Get(httpServerURL+"/hello", 1, nil, map[string]string{"token": "200"})
	assert.Nil(body)
	assert.True(strings.Contains(err.Error(), "deadline") || strings.Contains(err.Error(), "timeout"))
}

func TestHTTPClient_Get1(t *testing.T) {
	assert := assert2.New(t)
	client := NewHTTPClient()
	body, _, _, _ := client.Get(httpServerURL+"/hello", time.Second, nil, map[string]string{"token": "200"})
	assert.Equal("hello", string(body))
}

func TestHTTPClient_Get2(t *testing.T) {
	assert := assert2.New(t)
	client := NewHTTPClient()
	body, _, _, err := client.Get(httpServerURL+"/sleep", time.Second, nil, map[string]string{"token": "200"})
	assert.Nil(body)
	assert.True(strings.Contains(err.Error(), "deadline") || strings.Contains(err.Error(), "timeout"))
}

func TestHTTPClient_Get3(t *testing.T) {
	assert := assert2.New(t)
	client := NewHTTPClient()
	body, _, _, _ := client.Get(httpServerURL+"/sleep", time.Second*3, nil, map[string]string{"token": "200"})
	assert.Equal("hello", string(body))
}

func TestHTTPClient_Get4(t *testing.T) {
	assert := assert2.New(t)
	client := NewHTTPClient()
	body, _, _, _ := client.Get(httpsServerURL+"/hello", time.Second, nil, map[string]string{"token": "200"})
	assert.Contains(string(body), "hello")
}

func TestHTTPClient_Get4_1(t *testing.T) {
	assert := assert2.New(t)
	client := NewHTTPClient()
	client.SetPinCert(cert)
	body, _, _, _ := client.Get(httpsServerURL+"/hello", time.Second, nil, map[string]string{"token": "200"})
	assert.Contains(string(body), "hello")
}

func TestHTTPClient_Get4_2(t *testing.T) {
	assert := assert2.New(t)
	client := NewHTTPClient()
	client.SetRootCaCerts(cert)
	body, _, _, _ := client.Get(httpsServerURL+"/hello", time.Second, nil, map[string]string{"token": "200"})
	assert.Contains(string(body), "hello")
}

func TestHTTPClient_Get5(t *testing.T) {
	assert := assert2.New(t)
	client := NewHTTPClient()
	_, _, _, err := client.Get(httpsServerURL2+"/hello", time.Second, nil, map[string]string{"token": "200"})
	assert.Contains(err.Error(), "certificate")
}

func TestHTTPClient_Get6(t *testing.T) {
	assert := assert2.New(t)
	client := NewHTTPClient()
	client.SetClientCert(cert, key)
	body, _, resp, _ := client.Get(httpsServerURL2+"/hello", time.Second, nil, map[string]string{"token": "200"})
	assert.Contains(string(body), "hello")
	assert.Contains(string(GetResponseBody(resp)), "hello")
}

func TestHTTPClient_Before(t *testing.T) {
	assert := assert2.New(t)
	client := NewHTTPClient()
	client.SetClientCert(cert, key)
	var r *http.Request
	client.AppendBeforeDo(func(req *http.Request) {
		r = req
		req.Header.Set("test", "abc")
	})
	client.AppendAfterDo(func(req *http.Request, resp *http.Response, err error) {
		resp.Header.Set("test", "def")
	})
	body, _, resp, _ := client.Get(httpsServerURL2+"/hello", time.Second, nil, map[string]string{"token": "200"})
	assert.Contains(string(body), "hello")
	assert.Contains(string(GetResponseBody(resp)), "hello")
	assert.Equal(r.Header.Get("test"), "abc")
	assert.Equal(resp.Header.Get("test"), "def")
}

func TestHTTPClient_BeforeSet(t *testing.T) {
	assert := assert2.New(t)
	client := NewHTTPClient()
	client.SetClientCert(cert, key)
	var r *http.Request
	client.SetBeforeDo(func(req *http.Request) {
		r = req
		req.Header.Set("test", "abc")
	})
	client.SetAfterDo(func(req *http.Request, resp *http.Response, err error) {
		resp.Header.Set("test", "def")
	})
	body, _, resp, _ := client.Get(httpsServerURL2+"/hello", time.Second, nil, map[string]string{"token": "200"})
	assert.Contains(string(body), "hello")
	assert.Contains(string(GetResponseBody(resp)), "hello")
	assert.Equal(r.Header.Get("test"), "abc")
	assert.Equal(resp.Header.Get("test"), "def")
}

func TestHTTPClient_BeforeSetNil(t *testing.T) {
	assert := assert2.New(t)
	client := NewHTTPClient()
	client.SetClientCert(cert, key)
	var r *http.Request
	client.SetBeforeDo(func(req *http.Request) {
		r = req
		req.Header.Set("test", "abc")
	})
	client.SetBeforeDo(nil)
	client.SetAfterDo(func(req *http.Request, resp *http.Response, err error) {
		resp.Header.Set("test", "def")
	})
	client.SetAfterDo(nil)
	body, _, resp, _ := client.Get(httpsServerURL2+"/hello", time.Second, nil, map[string]string{"token": "200"})
	assert.Contains(string(body), "hello")
	assert.Contains(string(GetResponseBody(resp)), "hello")
	assert.Nil(r)
	assert.NotEqual(resp.Header.Get("test"), "def")
}

func TestHTTPClient_SetProxyFromEnv(t *testing.T) {
	assert := assert2.New(t)
	os.Setenv("HTTP_PROXY", "ftp://127.0.0.1:1000")
	client := NewHTTPClient()
	client.SetProxyFromEnv(true)
	_, _, _, err := client.Get(httpServerURL+"/hello", time.Second, nil, nil)
	assert.NotNil(err)
	Close()
	client.SetProxy("127.0.0.1:20000")
	assert.Equal("http://127.0.0.1:20000", client.ProxyUsed().String())
}

func TestHTTPClient_SetProxy(t *testing.T) {
	assert := assert2.New(t)
	client := NewHTTPClient()
	client.SetProxy("127.0.0.1:10000")
	_, _, _, err := client.Get(httpServerURL+"/hello", time.Second, nil, nil)
	assert.Contains(err.Error(), "connection refused")
}

func TestHTTPClient_SetDNS(t *testing.T) {
	u1, _ := url.Parse("http://foo.com/")
	u2, _ := url.Parse("http://127.0.0.1/")
	jar := NewCookieJar()
	jar.SetCookies(u1, []*http.Cookie{{Name: "foo", Value: "value1"}})
	jar.SetCookies(u2, []*http.Cookie{{Name: "foo", Value: "value2", Domain: "127.0.0.1"}})
	assert2.Equal(t, "value1", jar.Cookies(u1)[0].Value)
	assert := assert2.New(t)
	client := NewHTTPClient()
	client.SetDisableDNS(false)
	client.SetDNS("114.114.114.114:53")
	body1, _, _, _ := client.Get("http://www.baidu.com/", time.Second*5, nil, map[string]string{"token": "200"})
	client.SetDNS("8.8.8.8:53")
	body2, _, _, _ := client.Get("http://www.baidu.com/", time.Second*5, nil, map[string]string{"token": "200"})
	assert.True(strings.Contains(string(body1), "STATUS OK") || strings.Contains(string(body2), "STATUS OK"))
	client.SetDNS("8.8.8.8:5555", "114.114.114.114:53")
	body2, _, _, _ = client.Get("http://www.baidu.com/", time.Second*5, nil, map[string]string{"token": "200"})
	assert.True(strings.Contains(string(body1), "STATUS OK") || strings.Contains(string(body2), "STATUS OK"))
	assert2.Error(t, client.SetDNS("8.8.8.8"))
	assert2.Nil(t, client.SetDNS())
	client.SetDNS("8.8.8.8:5555", "114.114.114.114:5555")
	_, _, _, e := client.Get("http://www.baidu.com/", time.Second*5, nil, map[string]string{"token": "200"})
	assert2.NotNil(t, e)
}

func TestHTTPClient_SetDNS_AllFail(t *testing.T) {
	client := NewHTTPClient()
	client.SetDNS("0.0.0.256:555", "0.0.0.257:555")
	body1, _, _, err := client.Get("http://www.baidu.com/", time.Second*5, nil, map[string]string{"token": "200"})
	assert2.NotNil(t, err)
	assert2.Nil(t, body1)
	t.Log(string(body1))
}

func TestDownload(t *testing.T) {
	type args struct {
		u       string
		timeout time.Duration
		header  map[string]string
	}
	tests := []struct {
		name     string
		args     args
		wantData []byte
		wantErr  bool
	}{
		{"", args{
			u:       httpServerURL + "/hello",
			timeout: time.Second,
			header:  nil,
		}, []byte("hello"), true},
		{"", args{
			u:       httpServerURL + "/none",
			timeout: time.Second,
			header:  nil,
		}, []byte("404 page not found\n"), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotData, _, err := Download(tt.args.u, tt.args.timeout, nil, tt.args.header)
			if (err == nil) != tt.wantErr {
				t.Errorf("Download() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("Download() gotData = %v, want %v", gotData, tt.wantData)
			}
		})
	}
}

func TestDownloadToFile(t *testing.T) {
	type args struct {
		u       string
		timeout time.Duration
		header  map[string]string
		file    string
	}
	os.Mkdir("abc.txt", 0755)
	defer func() {
		os.Remove("abc.txt")
		os.Remove("a.txt")
		os.Remove("b.txt")
	}()
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"normal", args{
			u:       httpServerURL + "/hello",
			timeout: time.Second,
			header:  nil,
			file:    "a.txt",
		}, false},
		{"wrong_file", args{
			u:       httpServerURL + "/hello",
			timeout: time.Second,
			header:  nil,
			file:    "abc.txt",
		}, true},
		{"wrong_url", args{
			u:       httpServerURL + "/none",
			timeout: time.Second,
			header:  nil,
			file:    "b.txt",
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := DownloadToFile(tt.args.u, tt.args.timeout, nil, tt.args.header, tt.args.file); (err != nil) != tt.wantErr {
				t.Errorf("DownloadToFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDownloadToWriter(t *testing.T) {
	type args struct {
		u       string
		timeout time.Duration
		header  map[string]string
	}
	tests := []struct {
		name       string
		args       args
		wantWriter string
		wantErr    bool
	}{
		{"okay_url", args{
			u:       httpServerURL + "/hello",
			timeout: 0,
			header:  nil,
		}, "hello", false},
		{"wrong_url", args{
			u:       httpServerURL + "/none",
			timeout: 0,
			header:  nil,
		}, "404 page not found\n", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			resp, err := DownloadToWriter(tt.args.u, tt.args.timeout, nil, tt.args.header, writer)
			if (err != nil || resp.StatusCode != 200) != tt.wantErr {
				t.Errorf("DownloadToWriter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("DownloadToWriter() gotWriter = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}

func TestGet(t *testing.T) {
	type args struct {
		u       string
		timeout time.Duration
		header  map[string]string
	}
	tests := []struct {
		name     string
		args     args
		wantBody []byte
		wantCode int
		wantErr  bool
	}{
		{"okay_url", args{
			u:       httpServerURL + "/hello",
			timeout: time.Second,
			header:  nil,
		}, []byte("hello"), 200, false},
		{"wrong_url", args{
			u:       httpServerURL + "/none",
			timeout: time.Second,
			header:  nil,
		}, []byte("404 page not found\n"), 404, false},
		{"wrong_host", args{
			u:       "http://none/none",
			timeout: time.Second,
			header:  nil,
		}, nil, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBody, gotCode, _, err := Get(tt.args.u, tt.args.timeout, nil, tt.args.header)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBody, tt.wantBody) {
				t.Errorf("Get() gotBody = %v, want %v", string(gotBody), string(tt.wantBody))
			}
			if gotCode != tt.wantCode {
				t.Errorf("Get() gotCode = %v, want %v", gotCode, tt.wantCode)
			}
		})
	}
}
