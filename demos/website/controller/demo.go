package controller

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/snail007/gmc"
)

type Demo struct {
	gmc.Controller
}

// Before__ -> Foo-Method -> After__ , if panic occur will stop the function link call.

func (this *Demo) Before__() {
	method := filepath.Base(this.Request.URL.Path)
	if strings.Contains(method, "db") {
		return
	}
	if strings.Contains(method, "hello") {
		this.Write("hello")
	}
}
func (this *Demo) Index__() {
	this.View.Set("title", "i'm index page")
	this.StopE(this.View.Render("welcome").Err(), func() {
		this.Write("execute template fail")
	}, func() {
		this.Stop()
	})
	// never output
	this.Write(">>>>>")
}
func (this *Demo) Hello() {
	this.Write(" ")
}
func (this *Demo) After__() {
	method := filepath.Base(this.Request.URL.Path)
	if strings.Contains(method, "db") {
		return
	}
	if strings.Contains(method, "hello") {
		this.Write("world!")
	}
}
func (this *Demo) Protected() {
	this.Write("Protected")
}
func (this *Demo) DB() {
	db := gmc.DB.MySQL()
	rs, err := db.QuerySQL("select * from test")
	if err != nil {
		this.Stop(err)
	}
	this.Write(rs.Rows())
}
func (this *Demo) SessionSet() {
	err := this.SessionStart()
	if err != nil {
		this.Stop(err)
	}
	this.Session.Set("username", "foo")
	this.Write(this.Session.Get("username"))
}
func (this *Demo) SessionGet() {
	err := this.SessionStart()
	if err != nil {
		this.Stop(err)
	}
	this.Write(this.Session.Get("username"))
}

func (this *Demo) SessionGet1() {
	err := this.SessionStart()
	if err != nil {
		this.Stop(err)
	}
	this.View.Render("sess")
}

func (this *Demo) Error500() {
	a := 0
	a /= a
}

func (this *Demo) Version() {
	this.Write("v1")
}

func (this *Demo) Func() {
	this.View.Set("name", "hello")
	this.View.Render("func").Stop()
	// this will never called
	this.Write("def")
}

func (this *Demo) I18n1() {
	this.View.Render("i18n")
}

func (this *Demo) I18n2() {
	//this.Write(this.Lang," ",gmc.Tr(this.Lang,"001"))
	this.Write(this.Lang, " ", this.Tr("001", "here you should tips yourself, what's 001? 这里的文字提示自己这个001是什么."))
}

func (this *Demo) I18n3() {
	this.View.Render("i18n3")
}

func (this *Demo) Cache() {
	gmc.Cache.Redis().Set("test", "aaa", time.Second)
	v, _ := gmc.Cache.Redis().Get("test")
	this.Write(v)
}

func (this *Demo) List() {
	pager := this.Ctx.NewPager(10, 10000)
	this.View.Set("paginator", pager)
	this.View.Render("list")
}