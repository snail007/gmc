// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More infomation at https://github.com/snail007/gmc

package filestore

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/snail007/gmc/util/file"
	"github.com/snail007/gmc/http/session"
)

type FileStoreConfig struct {
	Dir    string
	GCtime int //seconds
	Prefix string
	Logger *log.Logger
	TTL    int64 //seconds
}

func NewConfig() FileStoreConfig {
	return FileStoreConfig{
		Dir:    os.TempDir(),
		GCtime: 300,
		TTL:    15 * 60,
		Prefix: ".gmcsession_",
		Logger: log.New(os.Stdout, "[filestore]", log.LstdFlags),
	}
}

type FileStore struct {
	session.Store
	cfg  FileStoreConfig
	lock *sync.RWMutex
}

func New(config interface{}) (st session.Store, err error) {
	cfg := config.(FileStoreConfig)
	if !fileutil.ExistsDir(cfg.Dir) {
		err = os.Mkdir(cfg.Dir, 0700)
		if err != nil {
			return
		}
	}
	if cfg.GCtime <= 0 {
		cfg.GCtime = 300
	}
	s := &FileStore{
		cfg:  cfg,
		lock: &sync.RWMutex{},
	}
	go s.gc()
	st = s
	return
}

func (s *FileStore) Load(sessionID string) (sess *session.Session, isExists bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	f := s.file(sessionID)
	if !fileutil.ExistsFile(f) {
		// s.cfg.Logger.Printf("filestore file not found: %s", f)
		return
	}
	str, err := ioutil.ReadFile(f)
	if err != nil {
		s.cfg.Logger.Printf("filestore read file error: %s", err)
		return
	}
	sess = session.NewSession()
	err = sess.Unserialize(string(str))
	if err != nil {
		sess = nil
		s.cfg.Logger.Printf("filestore unserialize error: %s", err)
		return
	}
	if time.Now().Unix()-sess.Touchtime() > s.cfg.TTL {
		sess = nil
		os.Remove(s.file(sessionID))
		return
	}
	isExists = true
	return
}
func (s *FileStore) Save(sess *session.Session) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	str, err := sess.Serialize()
	if err != nil {
		return
	}
	err = ioutil.WriteFile(s.file(sess.SessionID()), []byte(str), 0700)
	return
}

func (s *FileStore) Delete(sessionID string) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	err = os.Remove(s.file(sessionID))
	return
}
func (s *FileStore) file(sessionID string) string {
	return filepath.Join(s.cfg.Dir, s.cfg.Prefix+sessionID)
}
func (s *FileStore) gc() {
	defer func() {
		e := recover()
		if e != nil {
			fmt.Printf("filestore gc error: %s", e)
		}
	}()
	var files []string
	var err error
	var file *os.File
	first := true
	for {
		if first {
			first = false
		} else {
			time.Sleep(time.Second * time.Duration(s.cfg.GCtime))
		}
		files, err = filepath.Glob(s.file("*"))
		if err != nil {
			s.cfg.Logger.Printf("filestore gc error: %s", err)
			continue
		}
		for _, v := range files {
			if file != nil {
				file.Close()
			}
			file, err = os.Open(v)
			if err != nil {
				return
			}
			fileInfo, err := file.Stat()
			if err != nil {
				// s.cfg.Logger.Printf("filestore gc error: %s", err)
				continue
			}
			if time.Now().Unix()-fileInfo.ModTime().Unix() > s.cfg.TTL {
				os.Remove(v)
			}
		}
		if file != nil {
			file.Close()
		}
	}
}
