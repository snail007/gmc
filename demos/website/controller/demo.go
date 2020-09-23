package controller

import (
	gcontroller "github.com/snail007/gmc/http/controller"
)

type Demo struct {
	gcontroller.Controller
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
