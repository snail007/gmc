package session

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewSession(t *testing.T) {
	assert := assert.New(t)
	sess := NewSession()
	sess.Set("a", "b")
	assert.Equal(sess.Get("a"), "b")
	sess.Delete("a")
	assert.Nil(sess.Get("a"))
	sess.Set("a", "c")
	sess.Set("a", "d")
	assert.Equal(sess.Get("a"), "d")
	sess.Destory()
	assert.Nil(sess.Get("a"))
	//renew
	sess = NewSession()
	t1 := sess.Touchtime()
	time.Sleep(time.Second)
	sess.Touch()
	assert.NotEqual(sess.Touchtime(), t1)
	sess.Set("a", "c")
	v := sess.Values()
	assert.Equal(len(v), 1)
	//Destory
	sess.Destory()
	sess.Destory()
	sess.Set("a", "b")
	sess.Delete("a")
	sess.Touch()
	sess.Serialize()
	sess.Unserialize("")
	assert.Nil(sess.Values())
	assert.True(sess.IsDestory())
	assert.Len(sess.SessionID(), 32)
}
func TestSerialize(t *testing.T) {
	assert := assert.New(t)
	sess := NewSession()
	assert.NotNil(sess.Serialize())
}

func TestUnserialize(t *testing.T) {
	assert := assert.New(t)
	// NewSession()
	sess := NewSession()
	sess.Set("a", "b")
	str, err := sess.Serialize()
	assert.Nil(err)
	//fmt.Println(sess.Touchtime())
	sess2 := NewSession()
	err = sess2.Unserialize(str)
	assert.Nil(err)
	//fmt.Println(sess.Touchtime())
	assert.Equal(sess2.Touchtime(), sess.Touchtime())
	assert.Equal(sess2.SessionID(), sess.SessionID())
	assert.Equal(sess.Get("a").(string), "b")
	assert.Equal(sess2.Get("a"), sess.Get("a"))
}
