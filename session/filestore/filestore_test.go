package filestore

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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

	sess := session.NewSession(1)
	sess.Touch()
	err := store.Save(sess)
	assert.Nil(err)
	_, ok = store.Load(sess.SessionID())
	assert.True(ok)
	time.Sleep(time.Second * 2)
	_, ok = store.Load(sess.SessionID())
	assert.False(ok)
}
func TestMkdir(t *testing.T) {
	assert := assert.New(t)
	cfg := NewConfig()
	cfg.GCtime = 0
	cfg.Dir = ".gmcsess"
	os.RemoveAll(cfg.Dir)
	ioutil.WriteFile(cfg.Dir, []byte("."), 0700)
	New(cfg)
	os.Remove(cfg.Dir)
	store, _ := New(cfg)
	assert.DirExists(cfg.Dir)
	k := "testbbb"
	f := filepath.Join(cfg.Dir, cfg.Prefix+k)
	ioutil.WriteFile(f, []byte("\n"), 0700)
	_, ok := store.Load(k)
	assert.False(ok)
	os.Remove(f)
	os.RemoveAll(cfg.Dir)
}
func TestDelete(t *testing.T) {
	assert := assert.New(t)
	cfg := NewConfig()
	store, err := New(cfg)
	assert.Nil(err)
	sess0 := session.NewSession(10)
	sess0.Touch()
	store.Save(sess0)
	_, ok := store.Load(sess0.SessionID())
	assert.True(ok)
	store.Delete(sess0.SessionID())
	_, ok = store.Load(sess0.SessionID())
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
