package gcache

import (
	"fmt"
	"github.com/snail007/gmc/core"
	"github.com/snail007/gmc/util/cast"
	"time"
)

type (
	MemCache struct {
		gcore.Cache
		cfg *MemCacheConfig
		c   *MemoryCache
	}
	MemCacheConfig struct {
		CleanupInterval time.Duration
	}
)

func NewMemCacheConfig() *MemCacheConfig {
	return &MemCacheConfig{
		CleanupInterval: time.Second,
	}
}

func NewMemCache(cfg interface{}) *MemCache {
	cfg0 := cfg.(*MemCacheConfig)
	rc := &MemCache{
		cfg: cfg0,
	}

	rc.c = NewMemoryCache(NoExpiration, cfg0.CleanupInterval)
	return rc
}

func (s *MemCache) Has(key string) (bool, error) {
	_, ok := s.c.Get(key)
	return ok, nil
}
func (s *MemCache) Clear() error {
	s.c.Flush()
	return nil
}
func (s *MemCache) String() string {
	return fmt.Sprintf("gmc memory cache: %d", s.cfg.CleanupInterval/time.Second)
}

func (s *MemCache) Get(key string) (string, error) {
	v, b := s.c.Get(key)
	if b {
		return gcast.ToString(v), nil
	}
	return "", gcore.ErrKeyNotExists
}
func (s *MemCache) Set(key string, value string, ttl time.Duration) error {
	s.c.Set(key, value, ttl)
	return nil
}
func (s *MemCache) Del(key string) error {
	s.c.Delete(key)
	return nil
}
func (s *MemCache) Incr(key string) (v int64, err error) {
	err = s.c.Increment(key, 1)
	if nil == err {
		n, _ := s.c.Get(key)
		v = gcast.ToInt64(n)
	}
	return
}
func (s *MemCache) Decr(key string) (v int64, err error) {
	err = s.c.Decrement(key, 1)
	if nil == err {
		n, _ := s.c.Get(key)
		v = gcast.ToInt64(n)
	}
	return
}
func (s *MemCache) IncrN(key string, n int64) (v int64, err error) {
	err = s.c.Increment(key, n)
	if nil == err {
		n, _ := s.c.Get(key)
		v = gcast.ToInt64(n)
	}
	return
}
func (s *MemCache) DecrN(key string, n int64) (v int64, err error) {
	err = s.c.Decrement(key, n)
	if nil == err {
		n, _ := s.c.Get(key)
		v = gcast.ToInt64(n)
	}
	return
}
func (s *MemCache) GetMulti(keys []string) (map[string]string, error) {
	d := map[string]string{}
	for _, key := range keys {
		v, e := s.Get(key)
		if e != nil && gcore.IsNotExits(e) {
			return nil, e
		}
		d[key] = gcast.ToString(v)
	}
	return d, nil
}
func (s *MemCache) SetMulti(values map[string]string, ttl time.Duration) (err error) {
	for k, v := range values {
		s.c.Set(k, v, ttl)
	}
	return nil
}
func (s *MemCache) DelMulti(keys []string) (err error) {
	for _, k := range keys {
		s.c.Delete(k)
	}
	return nil
}
