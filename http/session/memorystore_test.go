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


func TestNewMemoryStore(t *testing.T) {
	assert := assert.New(t)
	sid := "testaaaa"
	_, ok := memoryStore.Load(sid)
	assert.False(ok)

	sess := NewSession()
	sess.Touch()
	assert.Nil(memoryStore.Save(sess))
	time.Sleep(time.Second * 3)
	_, ok = memoryStore.Load(sess.SessionID())
	assert.False(ok)
}
func TestLoad(t *testing.T) {
	assert := assert.New(t)
	cfg := NewMemoryStoreConfig()
	cfg.GCtime = 0
	store, _ := NewMemoryStore(cfg)
	k := "testbbb"
	_, ok := store.Load(k)
	assert.False(ok)
}
func TestMemoryStoreConfigDelete(t *testing.T) {
	assert := assert.New(t)
	cfg := NewMemoryStoreConfig()
	store, err := NewMemoryStore(cfg)
	assert.Nil(err)
	sess0 := NewSession()
	sess0.Touch()
	store.Save(sess0)
	_, ok := store.Load(sess0.SessionID())
	assert.True(ok)
	store.Delete(sess0.SessionID())
	_, ok = store.Load(sess0.SessionID())
	assert.False(ok)
}
func TestNoGC(t *testing.T) {
	assert := assert.New(t)
	cfg := NewMemoryStoreConfig()
	cfg.TTL = 1
	cfg.GCtime = 100
	store, _ := NewMemoryStore(cfg)
	sess0 := NewSession()
	sid := sess0.SessionID()
	sess0.Touch()
	store.Save(sess0)
	time.Sleep(time.Second * 2)
	_, ok := store.Load(sid)
	assert.False(ok)
}
