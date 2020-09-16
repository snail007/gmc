package controller

import (
	"net/http"

	"github.com/snail007/gmc/router"
)

type Controller struct {
	Response http.ResponseWriter
	Request  *http.Request
	Args     router.Params
}

func (this *Controller) PreCall__(w http.ResponseWriter, r *http.Request, ps router.Params) {
	this.Response = w
	this.Request = r
	this.Args = ps
}
func (this *Controller) PostCall__() {

}
