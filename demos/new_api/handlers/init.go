// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package handlers

import gcore "github.com/snail007/gmc/core"

func Init(api gcore.APIServer) {
	initMiddleware(api)
	initError(api)
	initHanlder(api)
}
