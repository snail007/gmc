package cache

import "time"

type Cache interface {
	Has(key string) (bool, error)
	Clear() error
	String() string

	Get(key string) (interface{}, error)
	Set(key string, value interface{}, ttl time.Duration) error
	Del(key string) error

	GetMulti(keys []string) (map[string]interface{}, error)
	SetMulti(values map[string]interface{}, ttl time.Duration) (err error)
	DelMulti(keys []string) (err error)
}
