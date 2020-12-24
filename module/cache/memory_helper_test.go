package gcache

import (
	"github.com/snail007/gmc/core"
	assert "github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMemCache_Get(t *testing.T) {
	assert := assert.New(t)
	_, err := cMem.Get("test")
	assert.True(gcore.IsNotExits(err))
}
func TestMemCache_Set(t *testing.T) {
	assert := assert.New(t)
	err := cMem.Set("test", "aaa", time.Second)
	assert.Nil(err)
	v, err := cMem.Get("test")
	assert.Nil(err)
	assert.Equal("aaa", v)
}
func TestMemCache_Expire(t *testing.T) {
	assert := assert.New(t)
	err := cMem.Set("test", "aaa", time.Millisecond*500)
	assert.Nil(err)
	time.Sleep(time.Second)
	_, err = cMem.Get("test")
	assert.True(gcore.IsNotExits(err))
}
func TestMemCache_Delete(t *testing.T) {
	assert := assert.New(t)
	err := cMem.Set("test", "aaa", time.Millisecond*500)
	assert.Nil(err)
	cMem.Del("test")
	_, err = cMem.Get("test")
	assert.True(gcore.IsNotExits(err))
}
func TestMemCache_Has(t *testing.T) {
	assert := assert.New(t)
	err := cMem.Set("test", "aaa", time.Millisecond*500)
	assert.Nil(err)
	assert.True(cMem.Has("test"))
}

func TestMemCache_Clean(t *testing.T) {
	assert := assert.New(t)
	err := cMem.Set("test", "aaa", time.Millisecond*500)
	assert.Nil(err)
	cMem.Clear()
	_, err = cMem.Get("test")
	assert.True(gcore.IsNotExits(err))
}
func TestMemCache_String(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("gmc memory cache: 1", cMem.String())
}
func TestIncr(t *testing.T) {
	assert := assert.New(t)
	// Set
	err := cMem.Set("k3", "1", time.Minute)
	assert.Nil(err)

	// incr
	data, err := cMem.Incr("k3")
	assert.Nil(err)
	assert.EqualValues(2,data)

	// decr
	data, err = cMem.Decr("k3")
	assert.Nil(err)
	assert.EqualValues(1,data)

	// incr N
	data, err = cMem.IncrN("k3",3)
	assert.Nil(err)
	assert.EqualValues(4,data)

	// decr N
	data, err = cMem.DecrN("k3",3)
	assert.Nil(err)
	assert.EqualValues(1,data)

	//Get
	d, err := cMem.Get("k3")
	assert.Nil(err)
	assert.Equal("1", d)
}
func Test_Multi(t *testing.T) {
	assert := assert.New(t)
	//SetMulti
	err := cMem.SetMulti(map[string]string{
		"k1": "111",
		"k2": "222",
	}, time.Minute)
	assert.Nil(err)

	//GetMulti
	data0, err := cMem.GetMulti([]string{"k1", "k2"})
	assert.Equal("111", data0["k1"])
	assert.Equal("222", data0["k2"])

	//DelMulti
	err = cMem.DelMulti([]string{"k1", "k2"})
	assert.Nil(err)

	_data, err := cMem.GetMulti([]string{"k1", "k2"})

	_,ok:=_data["k1"]
	assert.False(ok)
	_,ok=_data["k2"]
	assert.False(ok)

}
