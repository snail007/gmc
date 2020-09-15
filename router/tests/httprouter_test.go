package router

import (
	"net/http/httptest"
	"testing"

	"github.com/snail007/gmc/router"

	"github.com/snail007/gmc/controller"
)

type Controller struct {
	controller.Controller
}

func (this *Controller) Method1() {

}
func TestRoute(t *testing.T) {
	r := router.NewHttpRouter()
	r.Route("/user/", new(Controller))
	h, args, ok := r.Lookup("GET", "/user/method1")
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	h(w, req, router.Params{})
	t.Log(h, args, ok)
	t.Fail()
}
