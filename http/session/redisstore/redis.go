// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More infomation at https://github.com/snail007/gmc

package gmcredisstore

import (
	"github.com/snail007/gmc/core"
	logutil "github.com/snail007/gmc/util/log"
	"time"

	gmcredis "github.com/snail007/gmc/cache/redis"
	gmcsession "github.com/snail007/gmc/http/session"
)

type RedisStoreConfig struct {
	TTL      int64 //seconds
	RedisCfg *gmcredis.RedisCacheConfig
	Logger   gmccore.Logger
}

func NewRedisStoreConfig() RedisStoreConfig {
	return RedisStoreConfig{
		TTL:      3600,
		Logger:   logutil.New("[redisstore]"),
		RedisCfg: gmcredis.NewRedisCacheConfig(),
	}
}

type RedisStore struct {
	gmcsession.Store
	cfg   RedisStoreConfig
	cache gmccore.Cache
}

func New(config interface{}) (st gmcsession.Store, err error) {
	cfg := config.(RedisStoreConfig)
	s := &RedisStore{
		cfg:   cfg,
		cache: gmcredis.New(cfg.RedisCfg),
	}
	st = s
	return
}

func (s *RedisStore) Load(sessionID string) (sess *gmcsession.Session, isExists bool) {
	v, e := s.cache.Get(sessionID)
	if v == "" || e != nil {
		return
	}
	sess = gmcsession.NewSession()
	err := sess.Unserialize(v)
	if err != nil {
		sess = nil
		s.cfg.Logger.Warnf("redisstore unserialize error: %s", err)
		return
	}
	if time.Now().Unix()-sess.Touchtime() > s.cfg.TTL {
		sess = nil
		s.cache.Del(sessionID)
		return
	}
	isExists = true
	return
}
func (s *RedisStore) Save(sess *gmcsession.Session) (err error) {
	str, err := sess.Serialize()
	if err != nil {
		return
	}
	err = s.cache.Set(sess.SessionID(), str, time.Second*time.Duration(s.cfg.TTL))
	return
}

func (s *RedisStore) Delete(sessionID string) (err error) {
	err = s.cache.Del(sessionID)
	return
}
