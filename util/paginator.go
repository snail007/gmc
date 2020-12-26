package gutil

import (
	"github.com/snail007/gmc/util/cast"
	"math"
	"net/http"
	"net/url"
	"strconv"
)

type Paginator struct {
	request       *http.Request
	perPageNums   int
	maxPages      int
	nums          int64
	pageRange     []int
	pageNums      int
	page          int
	pageParamName string
}

func (p *Paginator) MaxPages() int {
	return p.maxPages
}

func (p *Paginator) SetMaxPages(maxPages int) {
	p.maxPages = maxPages
}

func (p *Paginator) PerPageNums() int {
	return p.perPageNums
}

func (p *Paginator) SetPerPageNums(perPageNums int) {
	p.perPageNums = perPageNums
}

func (p *Paginator) Request() *http.Request {
	return p.request
}

func (p *Paginator) SetRequest(request *http.Request) {
	p.request = request
}

func (p *Paginator) PageNums() int {
	if p.pageNums != 0 {
		return p.pageNums
	}
	pageNums := math.Ceil(float64(p.nums) / float64((*p).perPageNums))
	if p.maxPages > 0 {
		pageNums = math.Min(pageNums, float64(p.maxPages))
	}
	p.pageNums = int(pageNums)
	return p.pageNums
}

func (p *Paginator) Nums() int64 {
	return p.nums
}

func (p *Paginator) SetNums(nums interface{}) {
	p.nums = gcast.ToInt64(nums)
}

func (p *Paginator) Page() int {
	if p.page != 0 {
		return p.page
	}
	if p.request.Form == nil {
		p.request.ParseForm()
	}
	p.page, _ = strconv.Atoi(p.request.Form.Get(p.pageParamName))
	if p.page > p.PageNums() {
		p.page = p.PageNums()
	}
	if p.page <= 0 {
		p.page = 1
	}
	return p.page
}

func (p *Paginator) Pages() []int {
	if p.pageRange == nil && p.nums > 0 {
		var pages []int
		pageNums := p.PageNums()
		page := p.Page()
		switch {
		case page >= pageNums-4 && pageNums > 9:
			start := pageNums - 9 + 1
			pages = make([]int, 9)
			for i := range pages {
				pages[i] = start + i
			}
		case page >= 5 && pageNums > 9:
			start := page - 5 + 1
			pages = make([]int, int(math.Min(9, float64(page+4+1))))
			for i := range pages {
				pages[i] = start + i
			}
		default:
			pages = make([]int, int(math.Min(9, float64(pageNums))))
			for i := range pages {
				pages[i] = i + 1
			}
		}
		p.pageRange = pages
	}
	return p.pageRange
}

func (p *Paginator) PageLink(page int) string {
	link, _ := url.ParseRequestURI(p.request.RequestURI)
	values := link.Query()
	if page == 1 {
		values.Del(p.pageParamName)
	} else {
		values.Set(p.pageParamName, strconv.Itoa(page))
	}
	link.RawQuery = values.Encode()
	return link.String()
}
func (p *Paginator) PageBaseLink() string {
	link, _ := url.ParseRequestURI(p.request.RequestURI)
	values := link.Query()
	values.Del(p.pageParamName)
	values.Set(p.pageParamName, "")
	link.RawQuery = values.Encode()
	return link.String()
}
func (p *Paginator) PageLinkPrev() (link string) {
	if p.HasPrev() {
		link = p.PageLink(p.Page() - 1)
	}
	return
}

func (p *Paginator) PageLinkNext() (link string) {
	if p.HasNext() {
		link = p.PageLink(p.Page() + 1)
	}
	return
}

func (p *Paginator) PageLinkFirst() (link string) {
	return p.PageLink(1)
}

func (p *Paginator) PageLinkLast() (link string) {
	return p.PageLink(p.PageNums())
}

func (p *Paginator) HasPrev() bool {
	return p.Page() > 1
}

func (p *Paginator) HasNext() bool {
	return p.Page() < p.PageNums()
}

func (p *Paginator) IsActive(page int) bool {
	return p.Page() == page
}

func (p *Paginator) Offset() int {
	return (p.Page() - 1) * p.perPageNums
}

func (p *Paginator) HasPages() bool {
	return p.PageNums() > 1
}

func NewPaginator(req *http.Request, per int, nums int64, pageKey string) *Paginator {
	p := Paginator{}
	p.request = req
	p.pageParamName = pageKey
	if per <= 0 {
		per = 10
	}
	p.perPageNums = per
	p.SetNums(nums)
	return &p
}
