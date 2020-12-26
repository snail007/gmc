package gsession

import (
	"fmt"
	gcore "github.com/snail007/gmc/core"
	gconfig "github.com/snail007/gmc/util/config"
	"time"
)

func Init(config *gconfig.Config) (sessionStore gcore.SessionStorage, err error) {
	if !config.GetBool("session.enable") {
		return nil, nil
	}
	typ := config.GetString("session.store")
	if typ == "" {
		typ = "memory"
	}
	ttl := config.GetInt64("session.ttl")
	switch typ {
	case "file":
		cfg := NewFileStoreConfig()
		cfg.TTL = ttl
		cfg.Dir = config.GetString("session.file.dir")
		cfg.GCtime = config.GetInt("session.file.gctime")
		cfg.Prefix = config.GetString("session.file.prefix")
		sessionStore, err = NewFileStore(cfg)
	case "memory":
		cfg := NewMemoryStoreConfig()
		cfg.TTL = ttl
		cfg.GCtime = config.GetInt("session.memory.gctime")
		sessionStore, err = NewMemoryStore(cfg)
	case "redis":
		cfg := NewRedisStoreConfig()
		cfg.RedisCfg.Addr = config.GetString("session.redis.address")
		cfg.RedisCfg.Password = config.GetString("session.redis.password")
		cfg.RedisCfg.Prefix = config.GetString("session.redis.prefix")
		cfg.RedisCfg.Debug = config.GetBool("session.redis.debug")
		cfg.RedisCfg.Timeout = time.Second * config.GetDuration("session.redis.timeout")
		cfg.RedisCfg.DBNum = config.GetInt("session.redis.dbnum")
		cfg.RedisCfg.MaxIdle = config.GetInt("session.redis.maxidle")
		cfg.RedisCfg.MaxActive = config.GetInt("session.redis.maxactive")
		cfg.RedisCfg.MaxConnLifetime = time.Second * config.GetDuration("session.redis.maxconnlifetime")
		cfg.RedisCfg.Wait = config.GetBool("session.redis.wait")
		cfg.TTL = ttl
		sessionStore, err = NewRedisStore(cfg)
	default:
		err = fmt.Errorf("unknown session store type %s", typ)
	}
	return
}
