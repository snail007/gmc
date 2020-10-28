package gmccachehelper

import (
	"fmt"
	gmccachefile "github.com/snail007/gmc/cache/file"
	gmccachemem "github.com/snail007/gmc/cache/memory"
	gmccore "github.com/snail007/gmc/core"
	gmcerr "github.com/snail007/gmc/error"
	"time"

	gmcconfig "github.com/snail007/gmc/config"

	logutil "github.com/snail007/gmc/util/log"

	gmccacheredis "github.com/snail007/gmc/cache/redis"
	"github.com/snail007/gmc/util/cast"
)

var (
	myCache     = map[string]gmccore.Cache{}
	groupRedis  = map[string]gmccore.Cache{}
	groupMemory = map[string]gmccore.Cache{}
	groupFile   = map[string]gmccore.Cache{}
	logger      = logutil.New("")
	cfg         *gmcconfig.Config
)

func SetLogger(l gmccore.Logger) {
	logger = l
}

//RegistGroup parse app.toml database configuration, `cfg` is Config object of app.toml
func Init(cfg0 *gmcconfig.Config) (err error) {
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
				cfg := &gmccacheredis.RedisCacheConfig{
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
				groupRedis[id] = gmccacheredis.New(cfg)
			} else if k == "memory" {
				cfg := &gmccachemem.MemCacheConfig{
					CleanupInterval: time.Duration(cast.ToInt(vvv["cleanupinterval"])) * time.Second,
				}
				groupMemory[id] = gmccachemem.NewMemCache(cfg)
			} else if k == "file" {
				cfg := &gmccachefile.FileCacheConfig{
					Dir:             cast.ToString(vvv["dir"]),
					CleanupInterval: time.Duration(cast.ToInt(vvv["cleanupinterval"])) * time.Second,
				}

				groupFile[id],err = gmccachefile.NewFileCache(cfg)
				if err!=nil{
					return
				}
			}
		}
	}
	return
}

func Cache(id ...string) gmccore.Cache {
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
func Redis(id ...string) gmccore.Cache {
	return find("redis",id...)
}

func AddCacheU(id string, c gmccore.Cache) {
	myCache[id] = c
}
func CacheU(id ...string) gmccore.Cache {

	return find("user",id...)
}

//Memory acquires a memory cache object associated the id, id default is : `default`
func Memory(id ...string) gmccore.Cache {
	return find("memory",id...)
}
//File acquires a file cache object associated the id, id default is : `default`
func File(id ...string) gmccore.Cache {
	return find("file",id...)
}

func find(typ string,id ...string) gmccore.Cache {
	id0 := "default"
	if len(id) > 0 {
		id0 = id[0]
	}
	var v gmccore.Cache
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
		logger.Infof(gmcerr.New(fmt.Sprintf(f, v...)).String())
	}
}
