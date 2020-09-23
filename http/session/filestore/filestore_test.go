// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More infomation at https://github.com/snail007/gmc

package gmcfilestore

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	gmcsession "github.com/snail007/gmc/http/session"
	"github.com/snail007/gmc/util/fileutil"
	"github.com/stretchr/testify/assert"
)

var (
	store gmcsession.Store
)

func TestNew(t *testing.T) {
	assert := assert.New(t)
	sid := "testaaaa"
	_, ok := store.Load(sid)
	assert.False(ok)
	sess := gmcsession.NewSession().Touch()
	err := store.Save(sess)
	assert.Nil(err)
	_, ok = store.Load(sess.SessionID())
	assert.True(ok)
	for i := 0; i < 10; i++ {
		err := store.Save(gmcsession.NewSession().Touch())
		assert.Nil(err)
	}
	time.Sleep(time.Second * 3)
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
	s0 := store.(*FileStore)
	assert.DirExists(cfg.Dir)
	k := "testbbb"
	f := s0.file(k)
	dir := filepath.Dir(f)
	if !fileutil.ExistsDir(dir) {
		os.MkdirAll(dir, 0700)
	}
	err := ioutil.WriteFile(f, []byte("\n"), 0700)
	assert.Nil(err)
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
	sess0 := gmcsession.NewSession()
	sess0.Touch()
	store.Save(sess0)
	_, ok := store.Load(sess0.SessionID())
	assert.True(ok)
	store.Delete(sess0.SessionID())
	_, ok = store.Load(sess0.SessionID())
	assert.False(ok)
}
func TestDelete_2(t *testing.T) {
	assert := assert.New(t)
	cfg := NewConfig()
	cfg.TTL = 1
	store, err := New(cfg)
	assert.Nil(err)
	sess0 := gmcsession.NewSession()
	sess0.Touch()
	store.Save(sess0)
	_, ok := store.Load(sess0.SessionID())
	assert.True(ok)
	time.Sleep(time.Second * 2)
	_, ok = store.Load(sess0.SessionID())
	assert.False(ok)
}
func TestMain(m *testing.M) {
	cfg := NewConfig()
	cfg.GCtime = 1
	cfg.TTL = 1
	// fmt.Println(">>>", cfg)
	var err error
	store, err = New(cfg)
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(m.Run())
}
