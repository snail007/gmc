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
	assert2.False(t, r.Execute().Success())
	assert2.Len(t, r.ErrorAll(), 4)
	resp, err := r.Result()
	assert2.Nil(t, resp)
	assert2.NotNil(t, err)
	assert2.Nil(t, r.ResponseAll())
	resps, errs := r.ResultAll()
	assert2.Len(t, resps, 4)
	assert2.Len(t, errs, 4)
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
	assert2.True(t, r.WaitFirstSuccess().Execute().Success())
	resp, err := r.Result()
	assert2.Nil(t, err)
	assert2.NotNil(t, resp)
	assert2.Equal(t, "4", string(resp.Body()))
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
	assert2.Len(t, r.ErrorAll(), 3)
	resp, err := r.Result()
	assert2.Nil(t, err)
	assert2.NotNil(t, resp)
	assert2.Equal(t, "4", string(resp.Body()))
	assert2.Len(t, r.ResponseAll(), 1)
	resps, errs := r.ResultAll()
	assert2.Len(t, resps, 4)
	assert2.Len(t, errs, 4)
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
	assert2.False(t, r.WaitFirstSuccess().Execute().Success())
	resp, err := r.Result()
	assert2.Nil(t, resp)
	assert2.NotNil(t, err)
	resps, errs := r.ResultAll()
	assert2.Len(t, resps, 4)
	assert2.Len(t, errs, 4)
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
	assert2.Len(t, r.ErrorAll(), 4)
	resp, err := r.Result()
	assert2.Nil(t, resp)
	assert2.NotNil(t, err)
	assert2.Nil(t, r.ResponseAll())
	resps, errs := r.ResultAll()
	assert2.Len(t, resps, 4)
	assert2.Len(t, errs, 4)
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
	resp, err := r.Result()
	assert2.Nil(t, err)
	assert2.NotNil(t, resp)
	assert2.Equal(t, "4", string(resp.Body()))
}
