// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More infomation at https://github.com/snail007/gmc

package gsession

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
	sess.Destroy()
	assert.Nil(sess.Get("a"))
	//renew
	sess = NewSession()
	t1 := sess.TouchTime()
	time.Sleep(time.Second)
	sess.Touch()
	assert.NotEqual(sess.TouchTime(), t1)
	sess.Set("a", "c")
	v := sess.Values()
	assert.Equal(len(v), 1)
	//Destroy
	sess.Destroy()
	sess.Destroy()
	sess.Set("a", "b")
	sess.Delete("a")
	sess.Touch()
	_,e:=sess.Serialize()
	assert.Error(e,"session is destroy")
	assert.Error(sess.Unserialize(""),"session is destroy")
	assert.Nil(sess.Values())
	assert.True(sess.IsDestroy())
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
	assert.Equal(sess2.TouchTime(), sess.TouchTime())
	assert.Equal(sess2.SessionID(), sess.SessionID())
	assert.Equal(sess.Get("a").(string), "b")
	assert.Equal(sess2.Get("a"), sess.Get("a"))
}
