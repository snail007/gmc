// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gatomic

import (
	"net"
	"sync"
	"sync/atomic"
)

type Value struct {
	val atomic.Value
	l   sync.RWMutex
}

func NewValue(x ...interface{}) *Value {
	var val atomic.Value
	if len(x) == 1 && x[0] != nil {
		val.Store(x[0])
	}
	return &Value{val: val}
}

func (s *Value) SetVal(x interface{}) {
	s.l.Lock()
	defer s.l.Unlock()
	s.val.Store(x)
}

func (s *Value) Val() interface{} {
	s.l.RLock()
	defer s.l.RUnlock()
	return s.val.Load()
}

func (s *Value) SetAfterGet(f func(oldVal interface{}) (newVal interface{})) {
	s.l.Lock()
	defer s.l.Unlock()
	s.val.Store(f(s.val.Load()))
}

type String struct {
	val string
	l   sync.RWMutex
}

func (a *String) IsEmpty() bool {
	return len(a.Val()) == 0
}

func (a *String) Val() string {
	a.l.RLock()
	defer a.l.RUnlock()
	return a.val
}

func (a *String) String() string {
	return a.Val()
}

func (a *String) SetVal(val string) {
	a.l.Lock()
	defer a.l.Unlock()
	a.val = val
}

func NewString(val ...string) *String {
	v := ""
	if len(val) == 1 {
		v = val[0]
	}
	return &String{val: v}
}

func NewBool() *Bool {
	a := new(int32)
	return &Bool{val: a}
}

type Bool struct {
	val *int32
}

func (i *Bool) IsTrue() bool {
	return atomic.LoadInt32(i.val) == 1
}

func (i *Bool) IsFalse() bool {
	return atomic.LoadInt32(i.val) == 0
}

func (i *Bool) SetTrue() {
	atomic.StoreInt32(i.val, 1)
}

func (i *Bool) SetFalse() {
	atomic.StoreInt32(i.val, 0)
}

type Bytes struct {
	bytes []byte
	l     sync.RWMutex
}

func (a *Bytes) Bytes() []byte {
	a.l.RLock()
	defer a.l.RUnlock()
	return a.bytes
}

func (a *Bytes) Len() int {
	return len(a.Bytes())
}

func (a *Bytes) SetBytes(data []byte) {
	a.l.Lock()
	defer a.l.Unlock()
	a.bytes = data
}

func (a *Bytes) Append(bytes []byte) {
	a.l.Lock()
	defer a.l.Unlock()
	a.bytes = append(a.bytes, bytes...)
}

func NewBytes(val ...[]byte) *Bytes {
	var v []byte
	if len(val) == 1 {
		v = val[0]
	}
	return &Bytes{bytes: v}
}

type Conn struct {
	conn net.Conn
	l    sync.RWMutex
}

func (c *Conn) Val() net.Conn {
	c.l.RLock()
	defer c.l.RUnlock()
	return c.conn
}

func (c *Conn) SetVal(conn net.Conn) {
	c.l.Lock()
	defer c.l.Unlock()
	c.conn = conn
}

func NewConn(conn ...net.Conn) *Conn {
	var v net.Conn
	if len(conn) == 1 {
		v = conn[0]
	}
	return &Conn{conn: v}
}
