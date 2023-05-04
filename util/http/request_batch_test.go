// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package ghttp

import (
	"context"
	"fmt"
	gcast "github.com/snail007/gmc/util/cast"
	"github.com/snail007/gmc/util/gpool"
	gmap "github.com/snail007/gmc/util/map"
	assert2 "github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestBatchGet_1(t *testing.T) {
	t.Parallel()
	u := httpServerURL + "/batch?"
	reqUrls := []string{
		u + "idx=1&sleep=3",
		u + "idx=2&sleep=3",
		u + "idx=3&sleep=3",
		u + "idx=4&sleep=3",
	}
	r, err := NewBatchGet(reqUrls, time.Second, nil, nil)
	assert2.Nil(t, err)
	called := false
	r.AppendBeforeDo(func(idx int, req *http.Request) {
		called = true
	})
	resp := r.Pool(gpool.New(4)).Execute()
	assert2.True(t, called)
	assert2.False(t, resp.Success())
	assert2.Len(t, r.RespAll(), 4)
	assert2.Nil(t, resp.Resp())
	assert2.NotNil(t, r.RespAll()[0].Request())
	assert2.Equal(t, 4, resp.ErrorCount())
	assert2.Nil(t, r.RespAll()[0].Body())
}

func TestBatchGet_2(t *testing.T) {
	t.Parallel()
	u := httpServerURL + "/batch?"
	reqUrls := []string{
		u + "idx=1&sleep=3",
		u + "idx=2&sleep=3",
		u + "idx=3&sleep=3",
		u + "idx=4&sleep=0",
	}
	r, err := NewBatchGet(reqUrls, time.Second, nil, nil)
	assert2.Nil(t, err)
	v := 1
	r.AppendBeforeDo(func(idx int, req *http.Request) {
		v = 2
	})
	r.SetBeforeDo(func(idx int, req *http.Request) {
		v = 3
	})
	resp := r.WaitFirstSuccess().Pool(gpool.New(4)).Execute()
	assert2.True(t, resp.Success())
	assert2.Equal(t, 3, v)
	assert2.Nil(t, resp.Resp().Err())
	assert2.NotNil(t, resp.Resp().Response)
	assert2.Equal(t, "4", string(resp.Resp().Body()))
	assert2.Equal(t, 3, resp.Resp().Idx())

}

func TestBatchGet_3(t *testing.T) {
	t.Parallel()
	u := httpServerURL + "/batch?"
	reqUrls := []string{
		u + "idx=1&sleep=3",
		u + "idx=2&sleep=3",
		u + "idx=3&sleep=3",
		u + "idx=4&sleep=0",
	}
	r, err := NewBatchGet(reqUrls, time.Second, nil, nil)
	assert2.Nil(t, err)
	v := 1
	r.AppendBeforeDo(func(idx int, req *http.Request) {
		v = 2
	})
	r.SetBeforeDo(nil)
	assert2.Equal(t, 1, v)
	assert2.False(t, r.Execute().Success())
	assert2.Equal(t, r.ErrorCount(), 3)
	assert2.Nil(t, r.Resp().Err())
	assert2.NotNil(t, r.Resp().Response)
	assert2.Equal(t, "4", string(r.Resp().Body()))
	assert2.Len(t, r.RespAll(), 4)
	assert2.Less(t, r.Resp().UsedTime(), time.Second)
	assert2.Greater(t, r.Resp().UsedTime(), time.Duration(0))
	assert2.False(t, r.Resp().StartTime().IsZero())
	assert2.False(t, r.Resp().EndTime().IsZero())
}

func TestBatchGet_4(t *testing.T) {
	t.Parallel()
	u := httpServerURL + "/batch?"
	reqUrls := []string{
		u + "idx=1&sleep=3",
		u + "idx=2&sleep=3",
		u + "idx=3&sleep=3",
		u + "idx=4&sleep=3",
	}
	r, err := NewBatchGet(reqUrls, time.Second, nil, nil)
	assert2.Nil(t, err)
	resp := r.WaitFirstSuccess().Execute()
	assert2.False(t, resp.Success())
	assert2.Nil(t, resp.Resp())
	assert2.Equal(t, resp.ErrorCount(), 4)
}

func TestBatchGet_5(t *testing.T) {
	t.Parallel()
	u := httpServerURL + "/batch?"
	reqUrls := []string{
		u + "idx=1&sleep=3",
		u + "idx=2&sleep=3",
		u + "idx=3&sleep=3",
		u + "idx=4&sleep=3",
	}
	r, err := NewBatchURL(&http.Client{}, http.MethodGet, reqUrls, time.Second, nil, nil)
	if err != nil {
		return
	}
	assert2.Nil(t, err)
	resp := r.WaitFirstSuccess().Execute()
	assert2.False(t, resp.Success())
	assert2.Nil(t, resp.Resp())
	assert2.Equal(t, resp.ErrorCount(), 4)
}

func TestBatchGet_6(t *testing.T) {
	t.Parallel()
	u := httpServerURL + "/batch?"
	reqUrls := []string{
		u + "idx=1&nosleep=1",
		u + "idx=2&nosleep=1",
		u + "idx=3&nosleep=1",
		u + "idx=4&nosleep=1",
	}
	r, err := NewBatchURL(&http.Client{}, http.MethodGet, reqUrls, time.Second, nil, nil)
	if err != nil {
		return
	}
	assert2.Nil(t, err)
	resp := r.Execute()
	assert2.True(t, resp.Success())
}

func TestBatchGet_7(t *testing.T) {
	t.Parallel()
	u := httpServerURL + "/batch?"
	reqUrls := []string{
		u + "idx=1&nosleep=1",
	}
	var reqs []*http.Request
	var cancels []context.CancelFunc
	for _, v := range reqUrls {
		r, cancel, _ := NewRequest(http.MethodGet, v, time.Second, nil, nil)
		cancels = append(cancels, cancel)
		reqs = append(reqs, r)
	}
	r := NewBatchRequest(reqs, nil).SetAfterDo(func(resp *Response) {
		cancels[resp.Idx()]()
	})
	resp := r.Execute()
	assert2.True(t, resp.Success())
	r.client = &http.Client{}
	resp = r.SetAfterDo(nil).Execute()
	assert2.Empty(t, resp.RespAll()[0].Err())
	assert2.True(t, resp.Success())
}

func TestBatchPost_1(t *testing.T) {
	t.Parallel()
	u := httpServerURL + "/batch?"
	reqUrls := []string{
		u + "idx=1",
		u + "idx=2",
		u + "idx=3",
		u + "idx=4",
	}
	r, err := NewBatchPost(reqUrls, time.Second, gmap.Mss{"sleep": "3"}, nil)
	assert2.Nil(t, err)
	assert2.False(t, r.Execute().Success())
	assert2.Equal(t, 4, r.ErrorCount())
	assert2.Nil(t, r.Resp())
	assert2.Len(t, r.RespAll(), 4)
}

func TestBatchPost_2(t *testing.T) {
	t.Parallel()
	u := httpServerURL + "/batch?"
	reqUrls := []string{
		u + "idx=1",
		u + "idx=2",
		u + "idx=3",
		u + "idx=4&nosleep=1",
	}
	r, err := NewBatchPost(reqUrls, time.Second, gmap.Mss{"sleep": "3"}, nil)
	assert2.Nil(t, err)
	assert2.True(t, r.WaitFirstSuccess().Execute().Success())
	resp := r.Resp()
	assert2.Nil(t, resp.Err())
	assert2.NotNil(t, resp.Response)
	assert2.Equal(t, "4", string(resp.Body()))
}

func TestBatchRequest_CheckErrorFunc(t *testing.T) {
	t.Parallel()
	u := httpServerURL + "/batch?"
	reqUrls := []string{
		u + "idx=1",
		u + "idx=2",
		u + "idx=3",
		u + "idx=4&nosleep=1",
	}
	r, err := NewBatchPost(reqUrls, time.Second, gmap.Mss{"sleep": "3"}, nil)
	assert2.Nil(t, err)
	assert2.False(t, r.WaitFirstSuccess().MaxTry(2).
		CheckErrorFunc(func(idx int, req *http.Request, resp *http.Response) error {
			return fmt.Errorf("fail")
		}).
		Execute().Success())
	assert2.Equal(t, 4, r.ErrorCount())
	assert2.Nil(t, r.Resp())
	assert2.Equal(t, "fail", r.RespAll()[3].Err().Error())
	assert2.True(t, strings.Contains(r.Err().Error(), "context") || strings.Contains(r.Err().Error(), "timeout"))
	assert2.Contains(t, r.ErrAll()[3].Error(), "fail")

}

func TestBatchRequest_CheckErrorFunc_1(t *testing.T) {
	t.Parallel()
	u := httpServerURL + "/batch?"
	reqUrls := []string{
		u + "idx=1&nosleep=1",
		u + "idx=2&nosleep=1",
		u + "idx=3&nosleep=1",
		u + "idx=4&nosleep=1",
	}
	r, err := NewBatchPost(reqUrls, time.Second, gmap.Mss{"sleep": "3"}, nil)
	assert2.Nil(t, err)
	i := -1
	assert2.True(t, r.WaitFirstSuccess().MaxTry(2).CheckErrorFunc(func(idx int, req *http.Request, resp *http.Response) error {
		if idx == 0 {
			i++
			if i >= 1 {
				return nil
			}
			return fmt.Errorf("fail")
		}
		return nil
	}).Execute().Success())
	assert2.Equal(t, 0, r.ErrorCount())
	assert2.NotNil(t, r.Resp())
	assert2.Greater(t, gcast.ToInt(string(r.Resp().Body())), 0)
	r.Close()
}
