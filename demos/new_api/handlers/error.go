package handlers

import (
	"github.com/snail007/gmc"
	ghttpserver "github.com/snail007/gmc/http/server"
	"net/http"
)

func initError(api *ghttpserver.APIServer) {
	// sets a function to handle 404 requests.
	api.Handle404(error404)
	// sets a function to handle panic error.
	api.Handle500(error500)
}

func error404(c gmc.C) {
	c.Write("404")
}

func error500(c gmc.C, err interface{}) {
	c.WriteHeader(http.StatusInternalServerError)
	c.Write(gmc.StackE(err))
}
