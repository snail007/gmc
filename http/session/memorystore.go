// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More infomation at https://github.com/snail007/gmc

package gsession

import (
	"fmt"
	gcore "github.com/snail007/gmc/core"
	"sync"
	"time"
)

type MemoryStoreConfig struct {
	GCtime int //seconds
	Logger gcore.Logger
	TTL    int64 //seconds
}

func NewMemoryStoreConfig() MemoryStoreConfig {
	return MemoryStoreConfig{
		GCtime: 300,
		TTL:    15 * 60,
		Logger: gcore.Providers.Logger("")(nil, "[memorystore]"),
	}
}

type MemoryStore struct {
	gcore.SessionStorage
	cfg  MemoryStoreConfig
	data sync.Map
}

func NewMemoryStore(config interface{}) (st gcore.SessionStorage, err error) {
	cfg := config.(MemoryStoreConfig)
	s := &MemoryStore{
		cfg:  cfg,
		data: sync.Map{},
	}
	go s.gc()
	st = s
	return
}

func (s *MemoryStore) Load(sessionID string) (sess gcore.Session, isExists bool) {
	sess0, ok := s.data.Load(sessionID)
	if !ok {
		return
	}
	sess = sess0.(*Session)
	if time.Now().Unix()-sess.TouchTime() > s.cfg.TTL {
		sess = nil
		s.data.Delete(sessionID)
		return
	}
	isExists = true
	return
}
func (s *MemoryStore) Save(sess gcore.Session) (err error) {
	s.data.Store(sess.SessionID(), sess)
	return
}

func (s *MemoryStore) Delete(sessionID string) (err error) {
	s.data.Delete(sessionID)
	return
}

func (s *MemoryStore) gc() {
	defer gcore.Providers.Error("")().Recover(func(e interface{}) {
		fmt.Printf("memorystore gc error: %s", gcore.Providers.Error("")().StackError(e))
	})
	first := true
	for {
		if first {
			first = false
		} else {
			time.Sleep(time.Second * time.Duration(s.cfg.GCtime))
		}
		s.data.Range(func(k, v interface{}) bool {
			sess := v.(*Session)
			if time.Now().Unix()-sess.TouchTime() > s.cfg.TTL {
				s.data.Delete(k)
			}
			return true
		})
	}
}
