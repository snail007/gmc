package gmccachehelper

import (
	"log"
	"time"

	"github.com/snail007/gmc/util/logutil"

	gmccache "github.com/snail007/gmc/cache"
	gmccacheredis "github.com/snail007/gmc/cache/redis"
	gmcconfig "github.com/snail007/gmc/config/gmc"
	"github.com/snail007/gmc/util/castutil"
)

var (
	groupRedis = map[string]gmccache.Cache{}
	logger     = logutil.New("")
)

func SetLogger(l *log.Logger) {
	logger = l
}

//RegistGroup parse app.toml database configuration, `cfg` is GMCConfig object of app.toml
func RegistGroup(cfg *gmcconfig.GMCConfig) (err error) {
	for k, v := range cfg.Sub("cache").AllSettings() {
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
				cfg := gmccacheredis.RedisCacheConfig{
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
			}
		}
	}
	return
}

//Redis acquires a redis cache object associated the id, id default is : `default`
func Redis(id ...string) gmccache.Cache {
	id0 := "default"
	if len(id) > 0 {
		id0 = id[0]
	}
	v, ok := groupRedis[id0]
	if !ok {
		logf("[warn] redis cache `id`:%s not found", id0)
	}
	return v
}

func logf(f string, v ...interface{}) {
	if logger != nil {
		logger.Printf(f, v...)
	}
}
