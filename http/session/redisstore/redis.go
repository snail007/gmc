// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More infomation at https://github.com/snail007/gmc

package gredisstore

import (
	"github.com/snail007/gmc/core"
	"github.com/snail007/gmc/gmc/log"
	"time"

	gredis "github.com/snail007/gmc/gmc/cache/redis"
	gsession "github.com/snail007/gmc/http/session"
)

type RedisStoreConfig struct {
	TTL      int64 //seconds
	RedisCfg *gredis.RedisCacheConfig
	Logger   gcore.Logger
}

func NewRedisStoreConfig() RedisStoreConfig {
	return RedisStoreConfig{
		TTL:      3600,
		Logger:   glog.NewGMCLog("[redisstore]"),
		RedisCfg: gredis.NewRedisCacheConfig(),
	}
}

type RedisStore struct {
	gcore.SessionStorage
	cfg   RedisStoreConfig
	cache gcore.Cache
}

func New(config interface{}) (st gcore.SessionStorage, err error) {
	cfg := config.(RedisStoreConfig)
	s := &RedisStore{
		cfg:   cfg,
		cache: gredis.New(cfg.RedisCfg),
	}
	st = s
	return
}

func (s *RedisStore) Load(sessionID string) (sess gcore.Session, isExists bool) {
	v, e := s.cache.Get(sessionID)
	if v == "" || e != nil {
		return
	}
	sess = gsession.NewSession()
	err := sess.Unserialize(v)
	if err != nil {
		sess = nil
		s.cfg.Logger.Warnf("redisstore unserialize error: %s", err)
		return
	}
	if time.Now().Unix()-sess.TouchTime() > s.cfg.TTL {
		sess = nil
		s.cache.Del(sessionID)
		return
	}
	isExists = true
	return
}
func (s *RedisStore) Save(sess gcore.Session) (err error) {
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
