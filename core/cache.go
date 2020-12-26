package gcore

import (
	"fmt"
	"time"
)

type Cache interface {
	// Has returns true if cached value exists.
	Has(key string) (bool, error)
	// Clear deletes all cached data.
	Clear() error
	// String returns info about this driver.
	String() string
	// Get gets cached value by given key.
	Get(key string) (string, error)
	// Put puts value into cache with key and expire time.
	Set(key string, value string, ttl time.Duration) error
	// Delete deletes cached value by given key.
	Del(key string) error
	// GetMulti gets multiple keys's values at once.
	GetMulti(keys []string) (map[string]string, error)
	// SetMulti sets multiple keys's values at once.
	SetMulti(values map[string]string, ttl time.Duration) (err error)
	// DelMulti deletes multiple keys's values at once.
	DelMulti(keys []string) (err error)
	// Incr increases cached int-type value by given key as a counter.
	Incr(key string) (int64,error)
	// Decr decreases cached int-type value by given key as a counter.
	Decr(key string) (int64,error)
	// Incr increases N cached int-type value by given key as a counter.
	IncrN(key string,n int64) (int64,error)
	// Decr decreases N cached int-type value by given key as a counter.
	DecrN(key string,n int64) (int64,error)
}

var (
	// ErrKeyNotExists is the error of key not exists
	ErrKeyNotExists = fmt.Errorf("key not exists")
)

func IsNotExits(err error) bool {
	return ErrKeyNotExists == err
}
