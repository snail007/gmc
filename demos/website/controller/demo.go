package controller

import (
	"path/filepath"

	"github.com/snail007/gmc"
)

type Demo struct {
	gmc.Controller
}

func (this *Demo) Before__() {
	method := filepath.Base(this.Request.URL.Path)
	if method == "db" {
		return
	}
	if method == "hello" {
		this.Write("hello")
	}
}
func (this *Demo) Hello() {
	this.Write(" ")
}
func (this *Demo) After__() {
	method := filepath.Base(this.Request.URL.Path)
	if method == "db" {
		return
	}
	if method == "hello" {
		this.Write("world!")
	}
}
func (this *Demo) Protected() {
	this.Write("Protected")
}
func (this *Demo) DB() {
	db := gmc.DBMySQL()
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
func (this *Demo) Error500() {
	a := 0
	a /= a
}
