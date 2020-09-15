package controller

import (
	"net/http"

	"github.com/snail007/gmc/router"
)

type Controller struct {
	w      http.ResponseWriter
	r      *http.Request
	Router *router.HttpRouter
}

func (this *Controller) Pre_(w http.ResponseWriter) {
	//router *router.HttpRouter, w http.ResponseWriter, r *http.Request
	// this.r = r
	// this.w = w
	// this.Router = router
}
func (this *Controller) Post_() {

}
