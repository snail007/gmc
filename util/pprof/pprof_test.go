package ghttppprof

import (
	"github.com/snail007/gmc"
	ghttp "github.com/snail007/gmc/util/http"
	gtest "github.com/snail007/gmc/util/testing"
	"github.com/stretchr/testify/assert"
	"net"
	"os"
	"testing"
	"time"
)

func TestBindRouter(t *testing.T) {
	//gtest.DebugRunProcess(t)
	if gtest.RunProcess(t, func() {
		api, err := gmc.New.APIServer(gmc.New.Ctx(), "127.0.0.1:"+os.Getenv("API_PORT"))
		assert.NoError(t, err)
		BindRouter(api.Router(), "debug/")
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
	time.Sleep(time.Second * 8)
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
