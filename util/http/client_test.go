package ghttp

import (
	"crypto/md5"
	"fmt"
	assert2 "github.com/stretchr/testify/assert"
	"os"
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
}

func TestHTTPClient_PostOfReader(t *testing.T) {
	assert := assert2.New(t)
	client := NewHTTPClient()
	body, _, _, _ := client.PostOfReader(httpServerURL+"/post", strings.NewReader("name=snail007"), time.Second, map[string]string{"token": "200"})
	assert.Equal("snail007", string(body))
}

func TestHTTPClient_Upload(t *testing.T) {
	assert := assert2.New(t)
	client := NewHTTPClient()
	f, _ := os.Create("upload.bin")
	f.WriteString("a")
	f.Close()
	h := md5.New()
	h.Write([]byte("a"))
	s := fmt.Sprintf("%x", h.Sum(nil))
	body, _, err := client.Upload(httpServerURL+"/upload", "test", "upload.bin", map[string]string{"uid": "007"})
	assert.Nil(err)
	assert.Equal("007"+s, body)
	assert.FileExists("test.bin")
	os.Remove("test.bin")
	os.Remove("upload.bin")
}

func TestHTTPClient_Get(t *testing.T) {
	assert := assert2.New(t)
	client := NewHTTPClient()
	body, _, _, err := client.Get(httpServerURL+"/hello", 1, map[string]string{"token": "200"})
	assert.Nil(body)
	assert.Contains(err.Error(), "timeout")
}

func TestHTTPClient_Get1(t *testing.T) {
	assert := assert2.New(t)
	client := NewHTTPClient()
	body, _, _, _ := client.Get(httpServerURL+"/hello", time.Second, map[string]string{"token": "200"})
	assert.Equal("hello", string(body))
}

func TestHTTPClient_Get2(t *testing.T) {
	assert := assert2.New(t)
	client := NewHTTPClient()
	body, _, _, err := client.Get(httpServerURL+"/sleep", time.Second, map[string]string{"token": "200"})
	assert.Nil(body)
	assert.Contains(err.Error(), "timeout")
}

func TestHTTPClient_Get3(t *testing.T) {
	assert := assert2.New(t)
	client := NewHTTPClient()
	body, _, _, _ := client.Get(httpServerURL+"/sleep", time.Second*3, map[string]string{"token": "200"})
	assert.Equal("hello", string(body))
}

func TestHTTPClient_Get4(t *testing.T) {
	assert := assert2.New(t)
	client := NewHTTPClient()
	body, _, _, _ := client.Get(httpsServerURL+"/hello", time.Second, map[string]string{"token": "200"})
	assert.Contains(string(body), "hello")
}

func TestHTTPClient_Get4_1(t *testing.T) {
	assert := assert2.New(t)
	client := NewHTTPClient()
	client.SetPinCert(cert)
	body, _, _, _ := client.Get(httpsServerURL+"/hello", time.Second, map[string]string{"token": "200"})
	assert.Contains(string(body), "hello")
}

func TestHTTPClient_Get4_2(t *testing.T) {
	assert := assert2.New(t)
	client := NewHTTPClient()
	client.SetRootCaCerts(cert)
	body, _, _, _ := client.Get(httpsServerURL+"/hello", time.Second, map[string]string{"token": "200"})
	assert.Contains(string(body), "hello")
}

func TestHTTPClient_Get5(t *testing.T) {
	assert := assert2.New(t)
	client := NewHTTPClient()
	_, _, _, err := client.Get(httpsServerURL2+"/hello", time.Second, map[string]string{"token": "200"})
	assert.Contains(err.Error(), "tls: bad certificate")
}

func TestHTTPClient_Get6(t *testing.T) {
	assert := assert2.New(t)
	client := NewHTTPClient()
	client.SetClientCert(cert, key)
	body, _, _, _ := client.Get(httpsServerURL2+"/hello", time.Second, map[string]string{"token": "200"})
	assert.Contains(string(body), "hello")
}

func TestHTTPClient_SetProxyFromEnv(t *testing.T) {
	assert := assert2.New(t)
	os.Setenv("HTTP_PROXY", "127.0.0.1:10000")
	client := NewHTTPClient()
	client.SetProxyFromEnv(true)
	_, _, _, err := client.Get(httpServerURL+"/hello", time.Second, nil)
	assert.Contains(err.Error(), "connection refused")
}

func TestHTTPClient_SetProxy(t *testing.T) {
	assert := assert2.New(t)
	client := NewHTTPClient()
	client.SetProxy("127.0.0.1:10000")
	_, _, _, err := client.Get(httpServerURL+"/hello", time.Second, nil)
	assert.Contains(err.Error(), "connection refused")
}

func TestHTTPClient_SetDNS(t *testing.T) {
	assert := assert2.New(t)
	client := NewHTTPClient()
	client.SetDNS("8.8.4.4:53")
	body, _, _, _ := client.Get("http://www.baidu.com/", time.Second*5, map[string]string{"token": "200"})
	assert.Contains(string(body), "STATUS OK")
}
