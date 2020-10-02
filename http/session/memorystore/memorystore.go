// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More infomation at https://github.com/snail007/gmc

package gmcmemorystore

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	gmcerr "github.com/snail007/gmc/error"
	gmcsession "github.com/snail007/gmc/http/session"
)

type MemoryStoreConfig struct {
	GCtime int //seconds
	Logger *log.Logger
	TTL    int64 //seconds
}

func NewConfig() MemoryStoreConfig {
	return MemoryStoreConfig{
		GCtime: 300,
		TTL:    15 * 60,
		Logger: log.New(os.Stdout, "[memorystore]", log.LstdFlags),
	}
}

type MemoryStore struct {
	gmcsession.Store
	cfg  MemoryStoreConfig
	data sync.Map
}

func New(config interface{}) (st gmcsession.Store, err error) {
	cfg := config.(MemoryStoreConfig)
	s := &MemoryStore{
		cfg:  cfg,
		data: sync.Map{},
	}
	go s.gc()
	st = s
	return
}

func (s *MemoryStore) Load(sessionID string) (sess *gmcsession.Session, isExists bool) {
	sess0, ok := s.data.Load(sessionID)
	if !ok {
		return
	}
	sess = sess0.(*gmcsession.Session)
	if time.Now().Unix()-sess.Touchtime() > s.cfg.TTL {
		sess = nil
		s.data.Delete(sessionID)
		return
	}
	isExists = true
	return
}
func (s *MemoryStore) Save(sess *gmcsession.Session) (err error) {
	s.data.Store(sess.SessionID(), sess)
	return
}

func (s *MemoryStore) Delete(sessionID string) (err error) {
	s.data.Delete(sessionID)
	return
}

func (s *MemoryStore) gc() {
	defer func() {
		e := recover()
		if e != nil {
			fmt.Printf("memorystore gc error: %s", gmcerr.Stack(e))
		}
	}()
	first := true
	for {
		if first {
			first = false
		} else {
			time.Sleep(time.Second * time.Duration(s.cfg.GCtime))
		}
		s.data.Range(func(k, v interface{}) bool {
			sess := v.(*gmcsession.Session)
			if time.Now().Unix()-sess.Touchtime() > s.cfg.TTL {
				s.data.Delete(k)
			}
			return true
		})
	}
}
