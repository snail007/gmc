package router

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/snail007/gmc/router"

	"github.com/snail007/gmc/controller"
)

type Controller struct {
	controller.Controller
}

func (this *Controller) Method1() {
	this.Response.Write([]byte("OKAY"))
}
func TestRoute(t *testing.T) {
	assert := assert.New(t)
	r := router.NewHttpRouter()
	r.Controller("/user/", new(Controller))
	h, args, _ := r.Lookup("GET", "/user/method1")
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	h(w, req, args)
	body, _ := ioutil.ReadAll(w.Body)
	assert.Equal("OKAY", string(body))
}
