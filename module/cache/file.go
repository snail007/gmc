package gcache

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"github.com/snail007/gmc/core"
	"github.com/snail007/gmc/util"
	"github.com/snail007/gmc/util/cast"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	folder = ".gmcfilecache"
)

type FileCacheConfig struct {
	Dir             string
	CleanupInterval time.Duration
}

func NewFileCacheConfig() *FileCacheConfig {
	return &FileCacheConfig{
		CleanupInterval: time.Second,
		Dir:           os.TempDir() ,
	}
}

// Item represents a cache item.
type Item struct {
	Val     string
	Created int64
	TTL     int64
}

func (item *Item) hasExpired() bool {
	return item.TTL > 0 &&
		(time.Now().Unix()-item.Created) >= item.TTL
}

// FileCache represents a file cache adapter implementation.
type FileCache struct {
	gcore.Cache
	cfg *FileCacheConfig
}

// NewFileCache creates and returns a new file cache.
func NewFileCache(cfg interface{}) (cache *FileCache, err error) {
	cfg0 := cfg.(*FileCacheConfig)
	c := &FileCache{
		cfg: cfg0,
	}
	c.cfg.Dir = strings.Replace(c.cfg.Dir, "{tmp}", os.TempDir(), 1)
	if c.cfg.Dir==""{
		c.cfg.Dir="."
	}
	c.cfg.Dir, err = filepath.Abs(c.cfg.Dir)
	if err != nil {
		return
	}
	c.cfg.Dir = filepath.Join(c.cfg.Dir, folder)
	err = os.MkdirAll(c.cfg.Dir, 0700);
	if err != nil {
		return
	}
	go c.startGC()
	cache=c
	return
}

func (c *FileCache) filepath(key string) string {
	m := md5.Sum([]byte(key))
	hash := hex.EncodeToString(m[:])
	return filepath.Join(c.cfg.Dir, string(hash[0]), string(hash[1]), hash)
}

// Put puts value into cache with key and expire time.
// If expired is 0, it will be deleted by next GC operation.
func (c *FileCache) Set(key string, val string, ttl time.Duration) error {
	filename := c.filepath(key)
	item := &Item{val, time.Now().Unix(), int64(ttl / time.Second)}
	data, err := encodeGob(item)
	if err != nil {
		return err
	}

	os.MkdirAll(filepath.Dir(filename), 0700)
	return ioutil.WriteFile(filename, data, 0700)
}

func (c *FileCache) read(key string) (*Item, error) {
	filename := c.filepath(key)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	item := new(Item)
	return item, decodeGob(data, item)
}

// Get gets cached value by given key.
func (c *FileCache) Get(key string) (val string, err error) {
	item, err := c.read(key)
	if err != nil {
		if os.IsNotExist(err){
			err= gcore.KEY_NOT_EXISTS
		}
		return
	}

	if item.hasExpired() {
		os.Remove(c.filepath(key))
		err= gcore.KEY_NOT_EXISTS
		return
	}
	return gcast.ToString(item.Val), nil
}

// Delete deletes cached value by given key.
func (c *FileCache) Del(key string) error {
	return os.Remove(c.filepath(key))
}

func (c *FileCache) String() string {
	return fmt.Sprintf("gmc file cache, gc: %ds, dir: %s", c.cfg.CleanupInterval/time.Second, c.cfg.Dir)
}

// Incr increases cached int-type value by given key as a counter.
func (c *FileCache) Incr(key string) (int64, error) {
	item, err := c.read(key)
	if err != nil {
		return 0, err
	}
	item.Val = gcast.ToString(gcast.ToInt64(item.Val) + int64(1))
	err = c.Set(key, item.Val, time.Second*time.Duration(item.TTL))
	if err != nil {
		return 0, err
	}
	return gcast.ToInt64(item.Val), nil
}

// Decrease cached int value.
func (c *FileCache) Decr(key string) (int64, error) {
	item, err := c.read(key)
	if err != nil {
		return 0, err
	}
	item.Val = gcast.ToString(gcast.ToInt64(item.Val) - int64(1))
	err = c.Set(key, item.Val, time.Second*time.Duration(item.TTL))
	if err != nil {
		return 0, err
	}
	return gcast.ToInt64(item.Val), nil
}

// Incr value N by key
func (c *FileCache) IncrN(key string, n int64) (val int64, err error) {
	item, err := c.read(key)
	if err != nil {
		return 0, err
	}
	item.Val = gcast.ToString(gcast.ToInt64(item.Val) + int64(n))
	err = c.Set(key, item.Val, time.Second*time.Duration(item.TTL))
	if err != nil {
		return 0, err
	}
	return gcast.ToInt64(item.Val) , nil
}

// Decr value N by key
func (c *FileCache) DecrN(key string, n int64) (val int64, err error) {
	item, err := c.read(key)
	if err != nil {
		return 0, err
	}
	item.Val = gcast.ToString(gcast.ToInt64(item.Val) - int64(n))
	err = c.Set(key, item.Val, time.Second*time.Duration(item.TTL))
	if err != nil {
		return 0, err
	}
	return gcast.ToInt64(item.Val), nil
}

// IsExist returns true if cached value exists.
func (c *FileCache) Has(key string) (bool, error) {
	return gutil.Exists(c.filepath(key)), nil
}

// Flush deletes all cached data.
func (c *FileCache) Clear() error {
	return os.RemoveAll(c.cfg.Dir)
}

func (c *FileCache) GetMulti(keys []string) (map[string]string, error) {
	d := map[string]string{}
	for _, key := range keys {
		v, e := c.Get(key)
		if e != nil && !gcore.IsNotExits(e) {
			return nil, e
		}
		if !gcore.IsNotExits(e){
			d[key] = v
		}
	}
	return d, nil
}
func (c *FileCache) SetMulti(values map[string]string, ttl time.Duration) (err error) {
	for k, v := range values {
		err = c.Set(k, v, ttl)
		if nil != err {
			return
		}
	}
	return nil
}
func (c *FileCache) DelMulti(keys []string) (err error) {
	for _, k := range keys {
		err = c.Del(k)
		if nil != err {
			return
		}
	}
	return nil
}

func (c *FileCache) startGC() {
	if c.cfg.CleanupInterval < 1 {
		return
	}

	if err := filepath.Walk(c.cfg.Dir, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("Walk: %v", err)
		}

		if fi.IsDir() {
			return nil
		}

		data, err := ioutil.ReadFile(path)
		if err != nil && !os.IsNotExist(err) {
			fmt.Errorf("ReadFile: %v", err)
		}

		item := new(Item)
		if err = decodeGob(data, item); err != nil {
			return err
		}
		if item.hasExpired() {
			if err = os.Remove(path); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("Remove: %v", err)
			}
		}
		return nil
	}); err != nil {
		log.Printf("error gc cache files: %v", err)
	}

	time.AfterFunc(time.Duration(c.cfg.CleanupInterval)*time.Second, func() { c.startGC() })
}

func encodeGob(item *Item) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := gob.NewEncoder(buf).Encode(item)
	return buf.Bytes(), err
}

func decodeGob(data []byte, out *Item) error {
	buf := bytes.NewBuffer(data)
	return gob.NewDecoder(buf).Decode(&out)
}
