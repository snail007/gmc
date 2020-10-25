package gmccachefile_test

import (
	"fmt"
	gmccachefile "github.com/snail007/gmc/cache/file"
	"github.com/snail007/gmc/core"
	assert "github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

var (
	c gmccore.Cache

)

func TestMemCache_Get(t *testing.T) {
	assert := assert.New(t)
	_, err := c.Get("test")
	fmt.Println(err)
	assert.True(gmccore.IsNotExits(err))
}
func TestMemCache_Set(t *testing.T) {
	assert := assert.New(t)
	err := c.Set("test", "aaa", time.Second)
	assert.Nil(err)
	v, err := c.Get("test")
	assert.Nil(err)
	assert.Equal("aaa", v)
}
func TestMemCache_Expire(t *testing.T) {
	assert := assert.New(t)
	err := c.Set("test", "aaa", time.Second)
	assert.Nil(err)
	time.Sleep(time.Second * 2)
	_, err = c.Get("test")
	assert.True(gmccore.IsNotExits(err))
}
func TestMemCache_Delete(t *testing.T) {
	assert := assert.New(t)
	err := c.Set("test", "aaa", time.Second)
	assert.Nil(err)
	c.Del("test")
	_, err = c.Get("test")
	assert.True(gmccore.IsNotExits(err))
}
func TestMemCache_Has(t *testing.T) {
	assert := assert.New(t)
	err := c.Set("test", "aaa", time.Second)
	assert.Nil(err)
	assert.True(c.Has("test"))
}

func TestMemCache_Clean(t *testing.T) {
	assert := assert.New(t)
	err := c.Set("test", "aaa", time.Second)
	assert.Nil(err)
	c.Clear()
	_, err = c.Get("test")
	assert.True(gmccore.IsNotExits(err))
}
func TestMemCache_String(t *testing.T) {
	assert := assert.New(t)
	assert.Contains(c.String(), "gmc file cache, gc: 1s, dir: ")
}
func TestIncr(t *testing.T) {
	assert := assert.New(t)
	// Set
	err := c.Set("k3", "1", time.Minute)
	assert.Nil(err)

	// incr
	data, err := c.Incr("k3")
	assert.Nil(err)
	assert.EqualValues(2, data)

	// decr
	data, err = c.Decr("k3")
	assert.Nil(err)
	assert.EqualValues(1, data)

	// incr N
	data, err = c.IncrN("k3", 3)
	assert.Nil(err)
	assert.EqualValues(4, data)

	// decr N
	data, err = c.DecrN("k3", 3)
	assert.Nil(err)
	assert.EqualValues(1, data)

	//Get
	d, err := c.Get("k3")
	assert.Nil(err)
	assert.Equal("1", d)
}
func Test_Multi(t *testing.T) {
	assert := assert.New(t)
	//SetMulti
	err := c.SetMulti(map[string]string{
		"k1": "111",
		"k2": "222",
	}, time.Minute)
	assert.Nil(err)
	//GetMulti
	data0, err := c.GetMulti([]string{"k1", "k2"})
	assert.Nil(err)
	assert.Equal("111", data0["k1"])
	assert.Equal("222", data0["k2"])

	//DelMulti
	err = c.DelMulti([]string{"k1", "k2"})
	assert.Nil(err)

	_data, err := c.GetMulti([]string{"k1", "k2"})

	_, ok := _data["k1"]
	assert.False(ok)
	_, ok = _data["k2"]
	assert.False(ok)

}

func TestMain(m *testing.M) {
	cfg := gmccachefile.NewFileCacheConfig()
	var e error
	c, e = gmccachefile.NewFileCache(cfg)
	if e != nil {
		panic(e)
	}
	code := m.Run()
	os.Exit(code)
}
