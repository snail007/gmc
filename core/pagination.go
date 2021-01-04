// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gcore

import "net/http"

type Paginator interface {
	MaxPages() int
	SetMaxPages(maxPages int)
	PerPageNums() int
	SetPerPageNums(perPageNums int)
	Request() *http.Request
	SetRequest(request *http.Request)
	PageNums() int
	Nums() int64
	SetNums(nums interface{})
	Page() int
	Pages() []int
	PageLink(page int) string
	PageBaseLink() string
	PageLinkPrev() (link string)
	PageLinkNext() (link string)
	PageLinkFirst() (link string)
	PageLinkLast() (link string)
	HasPrev() bool
	HasNext() bool
	IsActive(page int) bool
	Offset() int
	HasPages() bool
}
