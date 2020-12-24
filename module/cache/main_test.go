package gcache

import (
	"github.com/alicebob/miniredis"
	gcore "github.com/snail007/gmc/core"
	"os"
	"testing"
)

var (
	cFile       gcore.Cache
	cMem        gcore.Cache
	cRedis      gcore.Cache
	redisServer *miniredis.Miniredis
)

func TestMain(m *testing.M) {
	var err error
	redisServer, err = miniredis.Run()
	if err != nil {
		panic(err)
	}
	cFile, err = NewFileCache(NewFileCacheConfig())
	if err != nil {
		panic(err)
	}
	cMem = NewMemCache(NewMemCacheConfig())
	code := m.Run()
	redisServer.Close()
	os.Exit(code)
}
