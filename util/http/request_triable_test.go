// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package ghttp

import (
	gmap "github.com/snail007/gmc/util/map"
	assert2 "github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestHTTPClient_NewTriableGet(t *testing.T) {
	t.Parallel()
	req, err := NewTriableGet(httpServerURL+"/try1", 3, time.Second, gmap.Mss{"msg": "he"}, gmap.Mss{"h1": "llo"})
	assert2.Nil(t, err)
	resp := req.Execute()
	assert2.Nil(t, resp.Err())
	assert2.True(t, string(resp.Body()) == "hello")
	assert2.True(t, string(resp.Body()) == "hello")
}

func TestHTTPClient_NewTriablePost(t *testing.T) {
	t.Parallel()
	req, err := NewTriablePost(httpServerURL+"/try2", 3, time.Second, gmap.Mss{"msg": "he"}, gmap.Mss{"h1": "llo"})
	assert2.Nil(t, err)
	resp := req.Execute()
	assert2.Nil(t, resp.Err())
	assert2.Len(t, req.ErrAll(), 0)
	assert2.NotNil(t, resp.Response)
	assert2.True(t, string(resp.Body()) == "hello")
}
