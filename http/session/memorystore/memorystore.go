// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More infomation at https://github.com/snail007/gmc

package gmemorystore

import (
	"fmt"
	"github.com/snail007/gmc"
	gcore "github.com/snail007/gmc/core"
	"github.com/snail007/gmc/gmc/log"
	"sync"
	"time"

	gerr "github.com/snail007/gmc/gmc/error"
	gsession "github.com/snail007/gmc/http/session"
)

type MemoryStoreConfig struct {
	GCtime int //seconds
	Logger gcore.Logger
	TTL    int64 //seconds
}

func NewConfig() MemoryStoreConfig {
	return MemoryStoreConfig{
		GCtime: 300,
		TTL:    15 * 60,
		Logger: log.NewGMCLog("[memorystore]"),
	}
}

type MemoryStore struct {
	gsession.Store
	cfg  MemoryStoreConfig
	data sync.Map
}

func New(config interface{}) (st gsession.Store, err error) {
	cfg := config.(MemoryStoreConfig)
	s := &MemoryStore{
		cfg:  cfg,
		data: sync.Map{},
	}
	go s.gc()
	st = s
	return
}

func (s *MemoryStore) Load(sessionID string) (sess *gsession.Session, isExists bool) {
	sess0, ok := s.data.Load(sessionID)
	if !ok {
		return
	}
	sess = sess0.(*gsession.Session)
	if time.Now().Unix()-sess.Touchtime() > s.cfg.TTL {
		sess = nil
		s.data.Delete(sessionID)
		return
	}
	isExists = true
	return
}
func (s *MemoryStore) Save(sess *gsession.Session) (err error) {
	s.data.Store(sess.SessionID(), sess)
	return
}

func (s *MemoryStore) Delete(sessionID string) (err error) {
	s.data.Delete(sessionID)
	return
}

func (s *MemoryStore) gc() {
	defer gmc.Recover(func(e interface{}) {
		fmt.Printf("memorystore gc error: %s", gerr.Stack(e))
	})
	first := true
	for {
		if first {
			first = false
		} else {
			time.Sleep(time.Second * time.Duration(s.cfg.GCtime))
		}
		s.data.Range(func(k, v interface{}) bool {
			sess := v.(*gsession.Session)
			if time.Now().Unix()-sess.Touchtime() > s.cfg.TTL {
				s.data.Delete(k)
			}
			return true
		})
	}
}
