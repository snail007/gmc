package handlers

import gcore "github.com/snail007/gmc/core"

func Init(api gcore.APIServer) {
	initMiddleware(api)
	initError(api)
	initHanlder(api)
}
