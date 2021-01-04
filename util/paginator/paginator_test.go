// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package paginator

import (
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestNewPaginator(t *testing.T) {
	assert := assert.New(t)
	r := httptest.NewRequest("GET", "http://foo.com/?a=1&page=10", nil)
	p := NewPaginator(r, 10, 1000, "page")
	p.SetMaxPages(100)
	assert.Equal(10, p.Page())
	assert.Equal(100, p.PageNums())
	assert.Equal(10, p.PerPageNums())
	assert.Equal("http://foo.com/?a=1&page=", p.PageBaseLink())
	assert.Equal(true, p.HasNext())
	assert.Equal(true, p.HasPrev())
	assert.Equal(true, p.IsActive(10))
	assert.Equal(100, p.MaxPages())
	assert.Equal(int64(1000), p.Nums())
	assert.Equal(true, p.HasPages())
	assert.Equal("http://foo.com/?a=1", p.PageLinkFirst())
	assert.Equal("http://foo.com/?a=1&page=100", p.PageLinkLast())
	assert.Equal("http://foo.com/?a=1&page=9", p.PageLinkPrev())
	assert.Equal("http://foo.com/?a=1&page=11", p.PageLinkNext())
	assert.Equal([]int{6, 7, 8, 9, 10, 11, 12, 13, 14}, p.Pages())
}
