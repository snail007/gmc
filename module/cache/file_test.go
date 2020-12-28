package gcache

import (
	"fmt"
	"github.com/snail007/gmc/core"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFileCache_Get(t *testing.T) {
	assert := assert.New(t)
	_, err := cFile.Get("test")
	fmt.Println(err)
	assert.True(gcore.IsNotExits(err))
}
func TestFileCache_Set(t *testing.T) {
	assert := assert.New(t)
	err := cFile.Set("test", "aaa", time.Second)
	assert.Nil(err)
	v, err := cFile.Get("test")
	assert.Nil(err)
	assert.Equal("aaa", v)
}
func TestFileCache_Expire(t *testing.T) {
	assert := assert.New(t)
	err := cFile.Set("test", "aaa", time.Second)
	assert.Nil(err)
	time.Sleep(time.Second * 2)
	_, err = cFile.Get("test")
	assert.True(gcore.IsNotExits(err))
}
func TestFileCache_Delete(t *testing.T) {
	assert := assert.New(t)
	err := cFile.Set("test", "aaa", time.Second)
	assert.Nil(err)
	cFile.Del("test")
	_, err = cFile.Get("test")
	assert.True(gcore.IsNotExits(err))
}
func TestFileCache_Has(t *testing.T) {
	assert := assert.New(t)
	err := cFile.Set("test", "aaa", time.Second)
	assert.Nil(err)
	assert.True(cFile.Has("test"))
}

func TestFileCache_Clean(t *testing.T) {
	assert := assert.New(t)
	err := cFile.Set("test", "aaa", time.Second)
	assert.Nil(err)
	cFile.Clear()
	_, err = cFile.Get("test")
	assert.True(gcore.IsNotExits(err))
}
func TestFileCache_String(t *testing.T) {
	assert := assert.New(t)
	assert.Contains(cFile.String(), "gmc file cache, gc: 1s, dir: ")
}
func TestFileIncr(t *testing.T) {
	assert := assert.New(t)
	// Set
	err := cFile.Set("k3", "1", time.Minute)
	assert.Nil(err)

	// incr
	data, err := cFile.Incr("k3")
	assert.Nil(err)
	assert.EqualValues(2, data)

	// decr
	data, err = cFile.Decr("k3")
	assert.Nil(err)
	assert.EqualValues(1, data)

	// incr N
	data, err = cFile.IncrN("k3", 3)
	assert.Nil(err)
	assert.EqualValues(4, data)

	// decr N
	data, err = cFile.DecrN("k3", 3)
	assert.Nil(err)
	assert.EqualValues(1, data)

	//Get
	d, err := cFile.Get("k3")
	assert.Nil(err)
	assert.Equal("1", d)
}
func Test_FileMulti(t *testing.T) {
	assert := assert.New(t)
	//SetMulti
	err := cFile.SetMulti(map[string]string{
		"k1": "111",
		"k2": "222",
	}, time.Minute)
	assert.Nil(err)
	//GetMulti
	data0, err := cFile.GetMulti([]string{"k1", "k2"})
	assert.Nil(err)
	assert.Equal("111", data0["k1"])
	assert.Equal("222", data0["k2"])

	//DelMulti
	err = cFile.DelMulti([]string{"k1", "k2"})
	assert.Nil(err)

	_data, err := cFile.GetMulti([]string{"k1", "k2"})

	_, ok := _data["k1"]
	assert.False(ok)
	_, ok = _data["k2"]
	assert.False(ok)

}
