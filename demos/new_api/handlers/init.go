package handlers

import gmchttpserver "github.com/snail007/gmc/http/server"

func Init(api *gmchttpserver.APIServer)   {
	initMiddleware(api)
	initError(api)
	initHanlder(api)
}
