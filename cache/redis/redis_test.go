package gmccacheredis

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
	assert.Equal("redigo: nil returned",err.Error())
	assert.Equal("", data)

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
func TestIncr(t *testing.T) {
	assert := assert.New(t)
	cfg := NewRedisCacheConfig()
	cfg.Addr = "127.0.0.1:6379"
	cfg.Debug = true
	rd := New(cfg)
	// Set
	err := rd.Set("k3", "1", time.Minute)
	assert.Nil(err)
	// incr
	data, err := rd.Incr("k3")
	assert.EqualValues(2,data)
	assert.Nil(err)
	// decr
	data, err = rd.Decr("k3")
	assert.EqualValues(1,data)
	assert.Nil(err)
	// incr N
	data, err = rd.IncrN("k3",3)
	assert.EqualValues(4,data)
	assert.Nil(err)
	// decr N
	data, err = rd.DecrN("k3",3)
	assert.EqualValues(1,data)
	assert.Nil(err)
	//Get
	d, err := rd.Get("k3")
	assert.Nil(err)
	assert.Equal("1", d)
}
func Test_Multi(t *testing.T) {
	assert := assert.New(t)
	cfg := NewRedisCacheConfig()
	cfg.Addr = "127.0.0.1:6379"
	cfg.Debug = true
	rd := New(cfg)
	//SetMulti
	err := rd.SetMulti(map[string]string{
		"k1": "111",
		"k2": "222",
	}, time.Minute)
	assert.Nil(err)

	//GetMulti
	data0, err := rd.GetMulti([]string{"k1", "k2"})
	assert.Equal("111", data0["k1"])
	assert.Equal("222", data0["k2"])

	//DelMulti
	err = rd.DelMulti([]string{"k1", "k2"})
	assert.Nil(err)

	_data, err := rd.GetMulti([]string{"k1", "k2"})

	_,ok:=_data["k1"]
	assert.False(ok)
	_,ok=_data["k2"]
	assert.False(ok)

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
