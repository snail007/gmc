package controller

import (
	"github.com/snail007/gmc"
)

type Demo struct {
	gmc.Controller
}

func (this *Demo) Index__() {
	this.View.
		Set("title", "It works!").
		Render("welcome")
}

func (this *Demo) Hello() {
	this.Write("fmt.Println(\"Hello GMC!\")")
}

func (this *Demo) Conn() {
	ctx := this.Ctx
	ctx.Write(ctx.Conn().LocalAddr().String(), " ", ctx.Conn().RemoteAddr().String())
}
