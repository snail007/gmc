package redis

import (
	"os"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"

	"github.com/alicebob/miniredis"

	"github.com/stretchr/testify/assert"
)

var (
	s *miniredis.Miniredis
)

func TestNew(t *testing.T) {
	assert := assert.New(t)
	cfg := NewRedisCacheConfig()
	cfg.Addr = s.Addr()
	rd := New(cfg)
	//Set
	err := rd.Set("k3", "aaa", time.Minute)
	assert.Nil(err)
	//Get
	data, err := redis.String(rd.Get("k3"))
	assert.Nil(err)
	assert.Equal("aaa", data)
	//Del
	err = rd.Del("k3")
	assert.Nil(err)
	data, err = redis.String(rd.Get("k3"))
	assert.NotNil(err)
	assert.Equal("", data)

	//SetMulti
	err = rd.SetMulti(map[string]interface{}{
		"k1": 111,
		"k2": 222,
	}, time.Minute)
	assert.Nil(err)

	//GetMulti
	data0, err := rd.GetMulti([]string{"k1", "k2"})
	assert.Equal("111", string(data0["k1"].([]byte)))
	assert.Equal("222", string(data0["k2"].([]byte)))

	//DelMulti
	err = rd.DelMulti([]string{"k1", "k2"})
	assert.Nil(err)
	_data, err := rd.GetMulti([]string{"k1", "k2"})
	assert.Nil(_data["k1"])
	assert.Nil(_data["k2"])

	//Clear
	rd.Set("k3", "aaa", time.Minute)
	rd.Set("k4", "bbb", time.Minute)
	b, _ := rd.Has("k4")
	assert.True(b)
	rd.Clear()
	data, err = redis.String(rd.Get("k4"))
	assert.NotNil(err)
	assert.Equal("", data)
	t.Log(rd.String())
}
func TestNew_2(t *testing.T) {
	assert := assert.New(t)
	s.RequireAuth("123")
	cfg := NewRedisCacheConfig()
	cfg.Addr = s.Addr()
	cfg.Debug = true
	cfg.Password = "123"
	cfg.Prefix = "__pre__"
	rd := New(cfg)
	//Set
	err := rd.Set("k3", "aaa", time.Minute)
	assert.Nil(err)
	//Get
	data, err := redis.String(rd.Get("k3"))
	assert.Nil(err)
	assert.Equal("aaa", data)
}
func TestMain(m *testing.M) {
	var err error
	s, err = miniredis.Run()
	if err != nil {
		panic(err)
	}
	code := m.Run()
	s.Close()
	os.Exit(code)
}
