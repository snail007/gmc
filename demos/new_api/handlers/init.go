package handlers

import ghttpserver "github.com/snail007/gmc/http/server"

func Init(api *ghttpserver.APIServer)   {
	initMiddleware(api)
	initError(api)
	initHanlder(api)
}
