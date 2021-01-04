// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gsession

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewFileStore(t *testing.T) {
	assert := assert.New(t)
	sid := "testaaaa"
	_, ok := fileStore.Load(sid)
	assert.False(ok)
	sess := NewSession()
	sess.Touch()
	err := fileStore.Save(sess)
	assert.Nil(err)
	_, ok = fileStore.Load(sess.SessionID())
	assert.True(ok)
	for i := 0; i < 10; i++ {
		s := NewSession()
		s.Touch()
		err := fileStore.Save(s)
		assert.Nil(err)
	}
	time.Sleep(time.Second * 3)
	_, ok = fileStore.Load(sess.SessionID())
	assert.False(ok)

}
func TestMkdir(t *testing.T) {
	assert := assert.New(t)
	cfg := NewFileStoreConfig()
	cfg.GCtime = 0
	cfg.Dir = ".gmcsess"
	os.RemoveAll(cfg.Dir)
	ioutil.WriteFile(cfg.Dir, []byte("."), 0700)
	_, err := NewFileStore(cfg)
	assert.NotNil(err)
	os.Remove(cfg.Dir)
	store, err := NewFileStore(cfg)
	assert.Nil(err)
	s0 := store.(*FileStore)
	assert.DirExists(cfg.Dir)
	k := "testbbb"
	f := s0.file(k)
	dir := filepath.Dir(f)
	if !ExistsDir(dir) {
		os.MkdirAll(dir, 0700)
	}
	err = ioutil.WriteFile(f, []byte("\n"), 0700)
	assert.Nil(err)
	_, ok := store.Load(k)
	assert.False(ok)
	os.Remove(f)
	os.RemoveAll(cfg.Dir)
}
func TestDelete(t *testing.T) {
	assert := assert.New(t)
	cfg := NewFileStoreConfig()
	store, err := NewFileStore(cfg)
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
func TestDelete_2(t *testing.T) {
	assert := assert.New(t)
	cfg := NewFileStoreConfig()
	cfg.TTL = 1
	store, err := NewFileStore(cfg)
	assert.Nil(err)
	sess0 := NewSession()
	sess0.Touch()
	store.Save(sess0)
	_, ok := store.Load(sess0.SessionID())
	assert.True(ok)
	time.Sleep(time.Second * 2)
	_, ok = store.Load(sess0.SessionID())
	assert.False(ok)
}
