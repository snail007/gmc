// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More infomation at https://github.com/snail007/gmc

package gmcfilestore

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	gmcerr "github.com/snail007/gmc/error"
	gmcsession "github.com/snail007/gmc/http/session"
	"github.com/snail007/gmc/util/fileutil"
)

type FileStoreConfig struct {
	Dir    string
	GCtime int //seconds
	Prefix string
	Logger *log.Logger
	TTL    int64 //seconds
}

var (
	folder = ".gmcsessions"
)

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
	gmcsession.Store
	cfg  FileStoreConfig
	lock *sync.RWMutex
}

func New(config interface{}) (st gmcsession.Store, err error) {
	cfg := config.(FileStoreConfig)
	cfg.Dir = strings.Replace(cfg.Dir, "{tmp}", os.TempDir(), 1)
	if cfg.Dir==""{
		cfg.Dir="."
	}
	if !fileutil.ExistsDir(cfg.Dir) {
		err = os.Mkdir(cfg.Dir, 0700)
		if err != nil {
			return
		}
	}
	cfg.Dir, err = filepath.Abs(cfg.Dir)
	if err != nil {
		return
	}
	cfg.Dir = filepath.Join(cfg.Dir, folder)
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

func (s *FileStore) Load(sessionID string) (sess *gmcsession.Session, isExists bool) {
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
	sess = gmcsession.NewSession()
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
func (s *FileStore) Save(sess *gmcsession.Session) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	str, err := sess.Serialize()
	if err != nil {
		return
	}
	f := s.file(sess.SessionID())
	dir := filepath.Dir(f)
	if !fileutil.ExistsDir(dir) {
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
	defer func() {
		e := recover()
		if e != nil {
			fmt.Printf("filestore gc error: %s", gmcerr.Stack(e))
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
		names := []string{}
		err = s.tree(s.cfg.Dir, &names)
		if err != nil {
			s.cfg.Logger.Printf("filestore gc error: %s", err)
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
