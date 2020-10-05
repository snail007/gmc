package gmccachehelper

import (
	"fmt"
	gmccachefile "github.com/snail007/gmc/cache/file"
	gmccachemem "github.com/snail007/gmc/cache/memory"
	gmcerr "github.com/snail007/gmc/error"
	"log"
	"time"

	gmcconfig "github.com/snail007/gmc/config"

	"github.com/snail007/gmc/util/logutil"

	gmccache "github.com/snail007/gmc/cache"
	gmccacheredis "github.com/snail007/gmc/cache/redis"
	"github.com/snail007/gmc/util/castutil"
)

var (
	myCache     = map[string]gmccache.Cache{}
	groupRedis  = map[string]gmccache.Cache{}
	groupMemory = map[string]gmccache.Cache{}
	groupFile   = map[string]gmccache.Cache{}
	logger      = logutil.New("")
	cfg         *gmcconfig.Config
)

func SetLogger(l *log.Logger) {
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
			if !castutil.ToBool(vvv["enable"]) {
				continue
			}
			id := castutil.ToString(vvv["id"])
			if k == "redis" {
				if _, ok := groupRedis[id]; ok {
					return
				}
				cfg := &gmccacheredis.RedisCacheConfig{
					Debug:           castutil.ToBool(vvv["debug"]),
					Prefix:          castutil.ToString(vvv["prefix"]),
					Logger:          logger,
					Addr:            castutil.ToString(vvv["address"]),
					Password:        castutil.ToString(vvv["password"]),
					DBNum:           castutil.ToInt(vvv["dbnum"]),
					MaxIdle:         castutil.ToInt(vvv["maxidle"]),
					MaxActive:       castutil.ToInt(vvv["maxactive"]),
					IdleTimeout:     time.Duration(castutil.ToInt(vvv["idletimeout"])) * time.Second,
					Wait:            castutil.ToBool(vvv["wait"]),
					MaxConnLifetime: time.Duration(castutil.ToInt(vvv["maxconnlifetime"])) * time.Second,
					Timeout:         time.Duration(castutil.ToInt(vvv["timeout"])) * time.Second,
				}
				groupRedis[id] = gmccacheredis.New(cfg)
			} else if k == "memory" {
				cfg := &gmccachemem.MemCacheConfig{
					CleanupInterval: time.Duration(castutil.ToInt(vvv["cleanupinterval"])) * time.Second,
				}
				groupMemory[id] = gmccachemem.NewMemCache(cfg)
			} else if k == "file" {
				cfg := &gmccachefile.FileCacheConfig{
					Dir:             castutil.ToString(vvv["dir"]),
					CleanupInterval: time.Duration(castutil.ToInt(vvv["cleanupinterval"])) * time.Second,
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

func Cache(id ...string) gmccache.Cache {
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
func Redis(id ...string) gmccache.Cache {
	return find("redis",id...)
}

func AddCacheU(id string, c gmccache.Cache) {
	myCache[id] = c
}
func CacheU(id ...string) gmccache.Cache {

	return find("user",id...)
}

//Memory acquires a memory cache object associated the id, id default is : `default`
func Memory(id ...string) gmccache.Cache {
	return find("memory",id...)
}
//File acquires a file cache object associated the id, id default is : `default`
func File(id ...string) gmccache.Cache {
	return find("file",id...)
}

func find(typ string,id ...string) gmccache.Cache {
	id0 := "default"
	if len(id) > 0 {
		id0 = id[0]
	}
	var v gmccache.Cache
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
		logger.Println(gmcerr.New(fmt.Sprintf(f, v...)).String())
	}
}
