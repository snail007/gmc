package controller

import (
	gmccontroller "github.com/snail007/gmc/http/controller"
)

type Demo struct {
	gmccontroller.Controller
}

func (this *Demo) Before__() {
	this.Write("hello")
}
func (this *Demo) Hello() {
	this.Write(" ")
}
func (this *Demo) After__() {
	this.Write("world!")
}
func (this *Demo) Protected() {
	this.Write("Protected")
}
