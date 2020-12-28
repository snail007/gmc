package gcache

import (
	"fmt"
	gcore "github.com/snail007/gmc/core"
	"time"

	"github.com/snail007/gmc/util/cast"
)

var (
	myCache     = map[string]gcore.Cache{}
	groupRedis  = map[string]gcore.Cache{}
	groupMemory = map[string]gcore.Cache{}
	groupFile   = map[string]gcore.Cache{}
	logger      gcore.Logger
	defaultCache       string
)

func SetLogger(l gcore.Logger) {
	logger = l
}

//RegistGroup parse app.toml database configuration, `cfg` is Config object of app.toml
func Init(cfg0 gcore.Config) (err error) {
	defaultCache = cfg0.GetString("cache.default")
	for k, v := range cfg0.Sub("cache").AllSettings() {
		if _, ok := v.([]interface{}); !ok {
			continue
		}
		for _, vv := range v.([]interface{}) {
			vvv := vv.(map[string]interface{})
			if !gcast.ToBool(vvv["enable"]) {
				continue
			}
			id := gcast.ToString(vvv["id"])
			if k == "redis" {
				if _, ok := groupRedis[id]; ok {
					return
				}
				cfg := &RedisCacheConfig{
					Debug:           gcast.ToBool(vvv["debug"]),
					Prefix:          gcast.ToString(vvv["prefix"]),
					Logger:          logger,
					Addr:            gcast.ToString(vvv["address"]),
					Password:        gcast.ToString(vvv["password"]),
					DBNum:           gcast.ToInt(vvv["dbnum"]),
					MaxIdle:         gcast.ToInt(vvv["maxidle"]),
					MaxActive:       gcast.ToInt(vvv["maxactive"]),
					IdleTimeout:     time.Duration(gcast.ToInt(vvv["idletimeout"])) * time.Second,
					Wait:            gcast.ToBool(vvv["wait"]),
					MaxConnLifetime: time.Duration(gcast.ToInt(vvv["maxconnlifetime"])) * time.Second,
					Timeout:         time.Duration(gcast.ToInt(vvv["timeout"])) * time.Second,
				}
				groupRedis[id] = NewRedisCache(cfg)
			} else if k == "memory" {
				cfg := &MemCacheConfig{
					CleanupInterval: time.Duration(gcast.ToInt(vvv["cleanupinterval"])) * time.Second,
				}
				groupMemory[id] = NewMemCache(cfg)
			} else if k == "file" {
				cfg := &FileCacheConfig{
					Dir:             gcast.ToString(vvv["dir"]),
					CleanupInterval: time.Duration(gcast.ToInt(vvv["cleanupinterval"])) * time.Second,
				}

				groupFile[id], err = NewFileCache(cfg)
				if err != nil {
					return
				}
			}
		}
	}
	return
}

func Cache(id ...string) gcore.Cache {
	switch defaultCache {
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
func Redis(id ...string) *RedisCache {
	return find("redis", id...).(*RedisCache)
}

func AddCacheU(id string, c gcore.Cache) {
	myCache[id] = c
}
func CacheU(id ...string) gcore.Cache {
	return find("user", id...)
}

//Memory acquires a memory cache object associated the id, id default is : `default`
func Memory(id ...string) *MemCache {
	return find("memory", id...).(*MemCache)
}

//File acquires a file cache object associated the id, id default is : `default`
func File(id ...string) *FileCache {
	return find("file", id...).(*FileCache)
}

func find(typ string, id ...string) gcore.Cache {
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
		logf("[warn] %s cache `id`:%s not found", typ, id0)
	}
	return v
}
func logf(f string, v ...interface{}) {
	if logger != nil {
		logger.Infof(gcore.Providers.Error("")().New(fmt.Sprintf(f, v...)).Error())
	}
}
