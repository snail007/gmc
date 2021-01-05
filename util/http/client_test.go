package ghttp

import (
	assert2 "github.com/stretchr/testify/assert"
	"testing"
	"time"
)

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

func TestHTTPClient_Get5(t *testing.T) {
	assert := assert2.New(t)
	client := NewHTTPClient()
	_, _, _, err := client.Get(httpsServerURL2+"/hello", time.Second, map[string]string{"token": "200"})
	assert.Contains(err.Error(), "tls: bad certificate")
}

func TestHTTPClient_Get6(t *testing.T) {
	assert := assert2.New(t)
	client := NewHTTPClient()
	client.SetClientCert(cert,key)
	body, _, _, _ := client.Get(httpsServerURL2+"/hello", time.Second, map[string]string{"token": "200"})
	assert.Contains(string(body), "hello")
}