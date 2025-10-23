package ghttppprof

import (
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/snail007/gmc"
	gcore "github.com/snail007/gmc/core"
	gctx "github.com/snail007/gmc/module/ctx"
	ghttp "github.com/snail007/gmc/util/http"
	gtest "github.com/snail007/gmc/util/testing"
	"github.com/stretchr/testify/assert"
)

func TestBindRouter(t *testing.T) {
	//gtest.DebugRunProcess(t)
	if gtest.RunProcess(t, func() {
		api, err := gmc.New.APIServer(gmc.New.Ctx(), "127.0.0.1:"+os.Getenv("API_PORT"))
		assert.NoError(t, err)
		a := []string{}
		BindRouter(api.Router(), "debug/", func(ctx gcore.Ctx) bool {
			a = append(a, ctx.Request().URL.Path)
			return true
		})
		api.Router().HandleAny("/lena", func(w http.ResponseWriter, r *http.Request, ps gcore.Params) {
			ctx := gctx.NewCtxWithHTTP(w, r)
			ctx.Write(len(a))
		})
		assert.NoError(t, api.Run())
		select {}
	}) {
		return
	}
	port := getFreeTCPPort()
	os.Setenv("API_PORT", port)
	p := gtest.NewProcess(t).Verbose(false)
	err := p.Start()
	defer p.Kill()
	assert.NoError(t, err)
	time.Sleep(time.Second * 15)
	_, _, err = ghttp.Download("http://127.0.0.1:"+port+"/debug/", time.Second, nil, nil)
	assert.NoError(t, err)
	_, _, err = ghttp.Download("http://127.0.0.1:"+port+"/debug/allocs?debug=1", time.Second, nil, nil)
	assert.NoError(t, err)
	_, _, err = ghttp.Download("http://127.0.0.1:"+port+"/debug/block?debug=1", time.Second, nil, nil)
	assert.NoError(t, err)
	_, _, err = ghttp.Download("http://127.0.0.1:"+port+"/debug/cmdline", time.Second, nil, nil)
	assert.NoError(t, err)
	_, _, err = ghttp.Download("http://127.0.0.1:"+port+"/debug/goroutine?debug=1", time.Second, nil, nil)
	assert.NoError(t, err)
	_, _, err = ghttp.Download("http://127.0.0.1:"+port+"/debug/heap?debug=1", time.Second, nil, nil)
	assert.NoError(t, err)
	_, _, err = ghttp.Download("http://127.0.0.1:"+port+"/debug/mutex?debug=1", time.Second, nil, nil)
	assert.NoError(t, err)
	_, _, err = ghttp.Download("http://127.0.0.1:"+port+"/debug/profile?seconds=1", time.Second*2, nil, nil)
	assert.NoError(t, err)
	_, _, err = ghttp.Download("http://127.0.0.1:"+port+"/debug/threadcreate?debug=1", time.Second, nil, nil)
	assert.NoError(t, err)
	_, _, err = ghttp.Download("http://127.0.0.1:"+port+"/debug/trace?seconds=1", time.Second*2, nil, nil)
	assert.NoError(t, err)
	b, _, err := ghttp.Download("http://127.0.0.1:"+port+"/lena", time.Second*2, nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, string(b), "10")
}

func getFreeTCPPort() (p string) {
	for {
		l, err := net.Listen("tcp", ":0")
		if err == nil {
			l.Close()
			_, p, _ = net.SplitHostPort(l.Addr().String())
			return
		}
	}
}
