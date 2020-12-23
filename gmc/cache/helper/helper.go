package gcachehelper

import (
	"fmt"
	gcore "github.com/snail007/gmc/core"
	gcachefile "github.com/snail007/gmc/gmc/cache/file"
	gcachemem "github.com/snail007/gmc/gmc/cache/memory"
	gerr "github.com/snail007/gmc/gmc/error"
	"github.com/snail007/gmc/gmc/log"
	"time"

	gconfig "github.com/snail007/gmc/config"

	gcacheredis "github.com/snail007/gmc/gmc/cache/redis"
	"github.com/snail007/gmc/util/cast"
)

var (
	myCache     = map[string]gcore.Cache{}
	groupRedis  = map[string]gcore.Cache{}
	groupMemory = map[string]gcore.Cache{}
	groupFile   = map[string]gcore.Cache{}
	logger      = log.NewGMCLog()
	cfg         *gconfig.Config
)

func SetLogger(l gcore.Logger) {
	logger = l
}

//RegistGroup parse app.toml database configuration, `cfg` is Config object of app.toml
func Init(cfg0 *gconfig.Config) (err error) {
	cfg = cfg0
	for k, v := range cfg.Sub("cache").AllSettings() {
		if _, ok := v.([]interface{}); !ok {
			continue
		}
		for _, vv := range v.([]interface{}) {
			vvv := vv.(map[string]interface{})
			if !cast.ToBool(vvv["enable"]) {
				continue
			}
			id := cast.ToString(vvv["id"])
			if k == "redis" {
				if _, ok := groupRedis[id]; ok {
					return
				}
				cfg := &gcacheredis.RedisCacheConfig{
					Debug:           cast.ToBool(vvv["debug"]),
					Prefix:          cast.ToString(vvv["prefix"]),
					Logger:          logger,
					Addr:            cast.ToString(vvv["address"]),
					Password:        cast.ToString(vvv["password"]),
					DBNum:           cast.ToInt(vvv["dbnum"]),
					MaxIdle:         cast.ToInt(vvv["maxidle"]),
					MaxActive:       cast.ToInt(vvv["maxactive"]),
					IdleTimeout:     time.Duration(cast.ToInt(vvv["idletimeout"])) * time.Second,
					Wait:            cast.ToBool(vvv["wait"]),
					MaxConnLifetime: time.Duration(cast.ToInt(vvv["maxconnlifetime"])) * time.Second,
					Timeout:         time.Duration(cast.ToInt(vvv["timeout"])) * time.Second,
				}
				groupRedis[id] = gcacheredis.New(cfg)
			} else if k == "memory" {
				cfg := &gcachemem.MemCacheConfig{
					CleanupInterval: time.Duration(cast.ToInt(vvv["cleanupinterval"])) * time.Second,
				}
				groupMemory[id] = gcachemem.NewMemCache(cfg)
			} else if k == "file" {
				cfg := &gcachefile.FileCacheConfig{
					Dir:             cast.ToString(vvv["dir"]),
					CleanupInterval: time.Duration(cast.ToInt(vvv["cleanupinterval"])) * time.Second,
				}

				groupFile[id],err = gcachefile.NewFileCache(cfg)
				if err!=nil{
					return
				}
			}
		}
	}
	return
}

func Cache(id ...string) gcore.Cache {
	switch cfg.GetString("cache.default") {
	case "redis":
		return Redis(id...)
	case "memory":
		return Memory(id...)
	case "file":
		return File(id...)
	default:
		return CacheU(id...)
	}
	return nil
}

//Redis acquires a redis cache object associated the id, id default is : `default`
func Redis(id ...string) gcore.Cache {
	return find("redis",id...)
}

func AddCacheU(id string, c gcore.Cache) {
	myCache[id] = c
}
func CacheU(id ...string) gcore.Cache {

	return find("user",id...)
}

//Memory acquires a memory cache object associated the id, id default is : `default`
func Memory(id ...string) gcore.Cache {
	return find("memory",id...)
}
//File acquires a file cache object associated the id, id default is : `default`
func File(id ...string) gcore.Cache {
	return find("file",id...)
}

func find(typ string,id ...string) gcore.Cache {
	id0 := "default"
	if len(id) > 0 {
		id0 = id[0]
	}
	var v gcore.Cache
	var ok bool
	switch typ {
	case "file":
		v, ok = groupFile[id0]
	case "memory":
		v, ok = groupMemory[id0]
	case "redis":
		v, ok = groupRedis[id0]
	case "user":
		v, ok = myCache[id0]
	default:
		logf("[warn] %s cache not found", typ)
	}
	if !ok {
		logf("[warn] %s cache `id`:%s not found",typ, id0)
	}
	return v
}
func logf(f string, v ...interface{}) {
	if logger != nil {
		logger.Infof(gerr.New(fmt.Sprintf(f, v...)).String())
	}
}
