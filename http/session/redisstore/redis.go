// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More infomation at https://github.com/snail007/gmc

package redisstore

import (
	"log"
	"os"
	"time"

	gmccache "github.com/snail007/gmc/cache"
	gmcredis "github.com/snail007/gmc/cache/redis"
	"github.com/snail007/gmc/http/session"
)

type RedisStoreConfig struct {
	TTL      int64 //seconds
	RedisCfg gmcredis.RedisCacheConfig
	Logger   *log.Logger
}

func NewRedisStoreConfig() RedisStoreConfig {
	return RedisStoreConfig{
		TTL:      3600,
		Logger:   log.New(os.Stdout, "[redisstore]", log.LstdFlags),
		RedisCfg: gmcredis.NewRedisCacheConfig(),
	}
}

type RedisStore struct {
	session.Store
	cfg   RedisStoreConfig
	cache gmccache.Cache
}

func New(config interface{}) (st session.Store, err error) {
	cfg := config.(RedisStoreConfig)
	s := &RedisStore{
		cfg:   cfg,
		cache: gmcredis.New(cfg.RedisCfg),
	}
	st = s
	return
}

func (s *RedisStore) Load(sessionID string) (sess *session.Session, isExists bool) {
	v, e := s.cache.Get(sessionID)
	if v == nil || e != nil {
		return
	}
	str := string(v.([]byte))
	sess = session.NewSession()
	err := sess.Unserialize(string(str))
	if err != nil {
		sess = nil
		s.cfg.Logger.Printf("redisstore unserialize error: %s", err)
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
func (s *RedisStore) Save(sess *session.Session) (err error) {
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
