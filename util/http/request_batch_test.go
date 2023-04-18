// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package ghttp

import (
	"fmt"
	gcast "github.com/snail007/gmc/util/cast"
	gmap "github.com/snail007/gmc/util/map"
	assert2 "github.com/stretchr/testify/assert"
	"net/http"
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
	resp := r.Execute()
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
	resp := r.WaitFirstSuccess().Execute()
	assert2.True(t, resp.Success())
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
}
