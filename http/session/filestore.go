// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gsession

import (
	"crypto/md5"
	"fmt"
	gcore "github.com/snail007/gmc/core"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type FileStoreConfig struct {
	Dir    string
	GCtime int //seconds
	Prefix string
	Logger gcore.Logger
	TTL    int64 //seconds
}

var (
	folder = ".gmcsessions"
)

func NewFileStoreConfig() FileStoreConfig {
	return FileStoreConfig{
		Dir:    os.TempDir(),
		GCtime: 300,
		TTL:    15 * 60,
		Prefix: ".gmcsession_",
		Logger: gcore.Providers.Logger("")(nil, "[filestore]"),
	}
}

type FileStore struct {
	gcore.SessionStorage
	cfg  FileStoreConfig
	lock *sync.RWMutex
}

func NewFileStore(config interface{}) (st gcore.SessionStorage, err error) {
	cfg := config.(FileStoreConfig)
	cfg.Dir = strings.Replace(cfg.Dir, "{tmp}", os.TempDir(), 1)
	if cfg.Dir == "" {
		cfg.Dir = "."
	}
	cfg.Dir, err = filepath.Abs(cfg.Dir)
	if err != nil {
		return
	}
	cfg.Dir = filepath.Join(cfg.Dir, folder)
	if !ExistsDir(cfg.Dir) {
		err = os.MkdirAll(cfg.Dir, 0700)
		if err != nil && !strings.Contains(err.Error(), "file exists") {
			return
		}
		err = nil
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

func (s *FileStore) Load(sessionID string) (sess gcore.Session, isExists bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	f := s.file(sessionID)
	if !ExistsFile(f) {
		// s.cfg.Logger.Printf("filestore file not found: %s", f)
		return
	}
	str, err := ioutil.ReadFile(f)
	if err != nil {
		s.cfg.Logger.Warnf("filestore read file error: %s", err)
		return
	}
	sess = NewSession()
	err = sess.Unserialize(string(str))
	if err != nil {
		sess = nil
		s.cfg.Logger.Warnf("filestore unserialize error: %s", err)
		return
	}
	if time.Now().Unix()-sess.TouchTime() > s.cfg.TTL {
		sess = nil
		os.Remove(s.file(sessionID))
		return
	}
	isExists = true
	return
}
func (s *FileStore) Save(sess gcore.Session) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	str, err := sess.Serialize()
	if err != nil {
		return
	}
	f := s.file(sess.SessionID())
	dir := filepath.Dir(f)
	if !ExistsDir(dir) {
		err = os.MkdirAll(dir, 0700)
		if err != nil {
			return
		}
	}
	err = ioutil.WriteFile(f, []byte(str), 0700)
	return
}

func (s *FileStore) Delete(sessionID string) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	err = os.Remove(s.file(sessionID))
	return
}
func (s *FileStore) file(sessionID string) string {
	f := s.cfg.Prefix + sessionID
	m := fmt.Sprintf("%x", md5.Sum([]byte(f)))
	subDir := fmt.Sprintf("%s/%s/%s", string(m[0]), string(m[1]), string(m[2]))
	path := filepath.Join(s.cfg.Dir, subDir, f)
	return path
}
func (s *FileStore) gc() {
	defer gcore.Providers.Error("")().Recover(func(e interface{}) {
		fmt.Printf("filestore gc error: %s", gcore.Providers.Error("")().StackError(e))
	})
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
		names := []string{}
		err = s.tree(s.cfg.Dir, &names)
		if err != nil {
			s.cfg.Logger.Warnf("filestore gc error: %s", err)
			continue
		}
		files = names
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
func (s *FileStore) tree(folder string, names *[]string) (err error) {
	f, err := os.Open(folder)
	if err != nil {
		return
	}
	defer f.Close()
	finfo, err := f.Stat()
	if err != nil {
		return
	}
	if !finfo.IsDir() {
		return
	}
	files, err := filepath.Glob(folder + "/*")
	if err != nil {
		return
	}
	var file *os.File
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
			return err
		}
		if fileInfo.IsDir() {
			err = s.tree(v, names)
			if err != nil {
				return err
			}
		} else {
			v, err = filepath.Abs(v)
			if err != nil {
				return err
			}
			name := filepath.Base(v)
			if !strings.HasPrefix(name, s.cfg.Prefix) {
				continue
			}
			v0 := strings.Replace(v, "\\", "/", -1)
			*names = append(*names, v0)
		}
	}
	if file != nil {
		file.Close()
	}
	return
}
func ExistsFile(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	stat, err := f.Stat()
	if err != nil {
		return false
	}
	if stat.IsDir() {
		return false
	}
	return true
}

func ExistsDir(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	stat, err := f.Stat()
	if err != nil {
		return false
	}
	if !stat.IsDir() {
		return false
	}
	return true
}
