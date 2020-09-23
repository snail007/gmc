package controller

import "github.com/snail007/gmc"

type Demo struct {
	gmc.Controller
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
