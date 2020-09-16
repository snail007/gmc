package memorystore

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/snail007/gmc/session"
	"github.com/stretchr/testify/assert"
)

var (
	store session.Store
)

func TestNew(t *testing.T) {
	assert := assert.New(t)
	sid := "testaaaa"
	_, ok := store.Load(sid)
	assert.False(ok)

	sess := session.NewSession()
	sess.Touch()
	assert.Nil(store.Save(sess))
	time.Sleep(time.Second * 3)
	_, ok = store.Load(sess.SessionID())
	assert.False(ok)
}
func TestLoad(t *testing.T) {
	assert := assert.New(t)
	cfg := NewConfig()
	cfg.GCtime = 0
	store, _ := New(cfg)
	k := "testbbb"
	_, ok := store.Load(k)
	assert.False(ok)
}
func TestDelete(t *testing.T) {
	assert := assert.New(t)
	cfg := NewConfig()
	store, err := New(cfg)
	assert.Nil(err)
	sess0 := session.NewSession()
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
	cfg := NewConfig()
	cfg.TTL = 1
	cfg.GCtime = 100
	store, _ := New(cfg)
	sess0 := session.NewSession()
	sid := sess0.SessionID()
	sess0.Touch()
	store.Save(sess0)
	time.Sleep(time.Second * 2)
	_, ok := store.Load(sid)
	assert.False(ok)
}
func TestMain(m *testing.M) {
	cfg := NewConfig()
	cfg.GCtime = 1
	cfg.TTL = 1
	var err error
	store, err = New(cfg)
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(m.Run())
}
