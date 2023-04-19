// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package ghttp

import (
	"fmt"
	gmap "github.com/snail007/gmc/util/map"
	assert2 "github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestHTTPClient_NewTriableGet(t *testing.T) {
	t.Parallel()
	req, err := NewTriableGet(httpServerURL+"/try1?key=1", 3, time.Second, gmap.Mss{"msg": "he"}, gmap.Mss{"h1": "llo"})
	assert2.Nil(t, err)
	called := false
	req.AppendBeforeDo(func(idx int, req *http.Request) {
		called = true
	})
	resp := req.Execute()
	assert2.True(t, called)
	assert2.Nil(t, resp.Err())
	assert2.Len(t, req.ErrAll(), 0)
	assert2.True(t, string(resp.Body()) == "hello")
}

func TestTriableRequest_CheckErrorFunc(t *testing.T) {
	t.Parallel()
	req, err := NewTriableGet(httpServerURL+"/try1?key=2", 3, time.Second, gmap.Mss{"msg": "he"}, gmap.Mss{"h1": "llo"})
	assert2.Nil(t, err)
	req.CheckErrorFunc(func(idx int, req *http.Request, resp *http.Response) error {
		return fmt.Errorf("fail")
	})
	v := 1
	req.AppendBeforeDo(func(idx int, req *http.Request) {
		v = 2
	})
	req.SetBeforeDo(func(idx int, req *http.Request) {
		v = 3
	})
	resp := req.Execute()
	assert2.Equal(t, 3, v)
	assert2.Equal(t, "fail", resp.Err().Error())
	assert2.Equal(t, 3, resp.Idx())
	assert2.True(t, string(resp.Body()) == "")
}

func TestHTTPClient_NewTriablePost(t *testing.T) {
	t.Parallel()
	req, err := NewTriablePost(httpServerURL+"/try1?key=3", 3,
		time.Second, gmap.Mss{"msg": "he"}, gmap.Mss{"h1": "llo", "host": "example.com"})
	assert2.Nil(t, err)
	v := 1
	req.AppendBeforeDo(func(idx int, req *http.Request) {
		v = 2
	})
	req.SetBeforeDo(nil)
	req.Keepalive(false)
	resp := req.Execute()
	assert2.Equal(t, 1, v)
	assert2.Nil(t, resp.Err())
	assert2.Len(t, req.ErrAll(), 0)
	assert2.NotNil(t, resp.Response)
	assert2.True(t, resp.Request().Close)
	assert2.True(t, string(resp.Body()) == "hello")
	assert2.Equal(t, "example.com", req.req.Host)

	req.Keepalive(false)
	resp = req.Execute()
	assert2.Nil(t, resp.Err())
	assert2.Nil(t, req.Err())
	assert2.Len(t, req.ErrAll(), 0)
	assert2.NotNil(t, resp.Response)
	assert2.True(t, resp.Request().Close)
	assert2.True(t, string(resp.Body()) == "hello")
	assert2.Equal(t, "example.com", req.req.Host)
}

func TestHTTPClient_NewTriableGet2(t *testing.T) {
	t.Parallel()
	req, err := NewTriableGet(httpServerURL+"/try1?key=4", 3, 0, gmap.Mss{"msg": "he"}, gmap.Mss{"h1": "llo"})
	assert2.Nil(t, err)
	resp := req.Execute()
	assert2.Nil(t, resp.Err())
	assert2.True(t, string(resp.Body()) == "hello")
}

func TestHTTPClient_NewTriableGet3(t *testing.T) {
	t.Parallel()
	tr, err := NewTriableRequestByURL(nil, http.MethodGet, httpServerURL+"/try1?key=5", 3, 0, gmap.Mss{"msg": "he"}, gmap.Mss{"h1": "llo"})
	assert2.Nil(t, err)
	resp := tr.Execute()
	assert2.Nil(t, resp.Err())
	assert2.True(t, string(resp.Body()) == "hello")
	resp = tr.Execute()
	assert2.Nil(t, resp.Err())
	assert2.True(t, string(resp.Body()) == "hello")
}

func TestHTTPClient_NewTriableGet4(t *testing.T) {
	t.Parallel()
	req, err := NewTriableGet(httpServerURL+"/try1?key=6", 2, time.Second, gmap.Mss{"msg": "he"}, gmap.Mss{"h1": "llo"})
	assert2.Nil(t, err)
	resp := req.Execute()
	assert2.NotNil(t, resp.Err())
	assert2.NotNil(t, req.Err())
	assert2.Len(t, req.ErrAll(), 3)
}
