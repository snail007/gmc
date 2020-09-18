package redis

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
	gmccache "github.com/snail007/gmc/cache"
)

type RedisCacheConfig struct {
	Debug           bool
	Prefix          string
	Logger          *log.Logger
	Addr            string
	Password        string
	DBNum           int
	MaxIdle         int
	MaxActive       int
	IdleTimeout     time.Duration
	Wait            bool
	MaxConnLifetime time.Duration
	Timeout         time.Duration
}

func NewRedisCacheConfig() RedisCacheConfig {
	return RedisCacheConfig{
		Debug:           false,
		Prefix:          "",
		Logger:          log.New(os.Stderr, "", log.LstdFlags),
		Addr:            "127.0.0.1:6379",
		Password:        "",
		DBNum:           0,
		MaxIdle:         3,
		MaxActive:       10,
		IdleTimeout:     time.Second * 300,
		Wait:            true,
		MaxConnLifetime: time.Second * 1800,
	}
}

type RedisCache struct {
	cfg         RedisCacheConfig
	pool        *redis.Pool
	connected   bool
	connectLock *sync.Mutex
}

// New redis cache
func New(cfg interface{}) gmccache.Cache {
	cfg0 := cfg.(RedisCacheConfig)
	rc := &RedisCache{
		cfg:         cfg0,
		connectLock: &sync.Mutex{},
	}

	return rc
}

// Connect to redis server
func (c *RedisCache) connect() {
	if c.connected {
		return
	}
	c.connectLock.Lock()
	defer c.connectLock.Unlock()
	c.newPool()
	c.logf("connect to server %s db is %d", c.cfg.Addr, c.cfg.DBNum)
	c.connected = true
}

// Get value by key
func (c *RedisCache) Get(key string) (val interface{}, err error) {
	c.connect()
	val, err = c.exec("Get", c.key(key))
	return
}

// Set value by key
func (c *RedisCache) Set(key string, val interface{}, ttl time.Duration) (err error) {
	c.connect()
	_, err = c.exec("SetEx", c.key(key), int64(ttl/time.Second), val)
	return
}

// Del value by key
func (c *RedisCache) Del(key string) (err error) {
	c.connect()
	_, err = c.exec("Del", c.key(key))
	return
}

// Has cache key
func (c *RedisCache) Has(key string) (bool, error) {
	c.connect()
	// return 0 OR 1
	one, err := redis.Int(c.exec("Exists", c.key(key)))
	return one == 1, err
}

// GetMulti values by keys
func (c *RedisCache) GetMulti(keys []string) (map[string]interface{}, error) {
	c.connect()
	conn := c.pool.Get()
	defer conn.Close()

	var args []interface{}
	for _, key := range keys {
		args = append(args, c.key(key))
	}

	list, err := redis.Values(conn.Do("MGet", args...))
	if err != nil {
		return nil, err
	}

	values := make(map[string]interface{}, len(keys))
	for i, val := range list {
		values[keys[i]] = val
	}

	return values, nil
}

// SetMulti values
func (c *RedisCache) SetMulti(values map[string]interface{}, ttl time.Duration) (err error) {
	c.connect()
	conn := c.pool.Get()
	defer conn.Close()

	// open multi
	conn.Send("Multi")
	ttlSec := int64(ttl / time.Second)

	for key, val := range values {
		// bs, _ := cache.Marshal(val)
		conn.Send("SetEx", c.key(key), ttlSec, val)
	}

	// do exec
	_, err = conn.Do("Exec")
	return
}

// DelMulti values by keys
func (c *RedisCache) DelMulti(keys []string) (err error) {
	c.connect()
	conn := c.pool.Get()
	defer conn.Close()

	var args []interface{}
	for _, key := range keys {
		args = append(args, c.key(key))
	}

	_, err = conn.Do("Del", args...)
	return
}

// Clear all caches
func (c *RedisCache) Clear() error {
	c.connect()
	conn := c.pool.Get()
	defer conn.Close()
	_, err := conn.Do("FlushDb")
	return err
}

// String get
func (c *RedisCache) String() string {
	pwd := "*"
	if c.cfg.Debug {
		pwd = c.cfg.Password
	}
	return fmt.Sprintf("connection info. url: %s, pwd: %s, dbNum: %d", c.cfg.Addr, pwd, c.cfg.DBNum)
}

// Key build
func (c *RedisCache) key(key string) string {
	if c.cfg.Prefix != "" {
		return fmt.Sprintf("%s:%s", c.cfg.Prefix, key)
	}
	return key
}

// actually do the redis cmds, args[0] must be the key name.
func (c *RedisCache) exec(commandName string, args ...interface{}) (reply interface{}, err error) {
	if len(args) < 1 {
		return nil, errors.New("missing required arguments")
	}

	conn := c.pool.Get()
	defer conn.Close()

	if c.cfg.Debug {
		st := time.Now()
		reply, err = conn.Do(commandName, args...)
		c.logf(
			"operate redis cache. command: %s, key: %v, elapsed time: %.03f\n",
			commandName, args[0], time.Since(st).Seconds()*1000,
		)
		return
	}
	reply, err = conn.Do(commandName, args...)
	if err != nil {
		c.logf("redis error :\n%s\n%s", commandName, args[0])
	}
	return
}

func (c *RedisCache) logf(format string, v ...interface{}) {
	if c.cfg.Debug && c.cfg.Logger != nil {
		c.cfg.Logger.Printf(format, v...)
	}
}

func (c *RedisCache) newPool() {
	c.pool = &redis.Pool{
		MaxIdle:         c.cfg.MaxIdle,
		IdleTimeout:     c.cfg.IdleTimeout,
		MaxActive:       c.cfg.MaxActive,
		MaxConnLifetime: c.cfg.MaxConnLifetime,
		Wait:            c.cfg.Wait,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", c.cfg.Addr,
				redis.DialConnectTimeout(c.cfg.Timeout),
			)
			if err != nil {
				return nil, err
			}
			if c.cfg.Password != "" {
				_, err := conn.Do("AUTH", c.cfg.Password)
				if err != nil {
					_ = conn.Close()
					return nil, err
				}
			}
			_, _ = conn.Do("SELECT", c.cfg.DBNum)
			return conn, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
