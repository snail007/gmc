// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More infomation at https://github.com/snail007/gmc

package gsession

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io"
	"sync"
	"time"
)

type sData struct {
	ID        string
	Values    map[interface{}]interface{}
	Touchtime int64
}
type Session struct {
	id        string
	values    map[interface{}]interface{}
	lock      *sync.Mutex
	isDestroy bool
	touchtime int64
}

func init()  {
	gob.Register([]interface{}{})
	gob.Register(map[interface{}]interface{}{})
	gob.Register(map[string]string{})
	gob.Register(map[string]interface{}{})
	gob.Register(map[int]interface{}{})
	gob.Register(map[int]string{})
	gob.Register(map[int]int{})
	gob.Register(map[int]int64{})
	gob.Register(sData{})
}
func NewSession() *Session {
	s := &Session{
		id:     newSessionID(),
		lock:   &sync.Mutex{},
		values: map[interface{}]interface{}{},
	}
	return s
}
func (s *Session) Set(k interface{}, v interface{}) {
	if s.isDestroy {
		return
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	s.values[k] = v
	gob.Register(v)
	s.touch()
	return
}
func (s *Session) Get(k interface{}) (value interface{}) {
	if s.isDestroy {
		return
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	v, ok := s.values[k]
	s.touch()
	if ok {
		return v
	}
	return nil
}
func (s *Session) Delete(k interface{}) (err error) {
	if s.isDestroy {
		return
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.values, k)
	s.touch()
	return
}
func (s *Session) Destroy() (err error) {
	if s.isDestroy {
		return
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	s.values = map[interface{}]interface{}{}
	s.isDestroy = true
	s.touch()
	return
}
func (s *Session) Values() (data map[interface{}]interface{}) {
	if s.isDestroy {
		return
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	data = map[interface{}]interface{}{}
	for k, v := range s.values {
		data[k] = v
	}
	s.touch()
	return
}
func (s *Session) IsDestroy() bool {
	return s.isDestroy
}
func (s *Session) SessionID() (sessionID string) {
	return s.id
}

//Touchtime return the last access unix time seconds of session
//
//time: unix time seconds
func (s *Session) Touchtime() (time int64) {
	return s.touchtime
}
func (s *Session) Touch() *Session {
	if s.isDestroy {
		return s
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	s.touch()
	return s
}
func (s *Session) touch() {
	s.touchtime = time.Now().Unix()
	return
}
func newSessionID() (sessionID string) {
	k := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, k); err != nil {
		return ""
	}
	return hex.EncodeToString(k)
}
func (s *Session) Serialize() (str string, err error) {
	if s.isDestroy {
		err = fmt.Errorf("session is destroy")
		return
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	item := sData{
		ID:        s.id,
		Values:    s.values,
		Touchtime: s.touchtime,
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(item)
	if err != nil {
		return
	}
	str = base64.StdEncoding.EncodeToString(buf.Bytes())
	return
}
func (s *Session) Unserialize(data string) (err error) {
	if s.isDestroy {
		err = fmt.Errorf("session is destroy")
		return
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	bs, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return
	}
	d := bytes.NewBuffer(bs)
	dec := gob.NewDecoder(d)
	var q sData
	err = dec.Decode(&q)
	if err != nil {
		return
	}
	s.touchtime = q.Touchtime
	s.values = q.Values
	s.id = q.ID
	return
}
