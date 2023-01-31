// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gatomic

import (
	assert2 "github.com/stretchr/testify/assert"
	"io"
	"sync"
	"testing"
)

func TestNewValue(t *testing.T) {
	assert := assert2.New(t)
	type data struct {
		cnt int
	}
	g := sync.WaitGroup{}
	g.Add(6)
	value := NewValue(data{1})
	for i := 0; i < 2; i++ {
		go func() {
			defer g.Done()
			value.SetVal(data{1})
		}()
	}

	for i := 0; i < 2; i++ {
		go func() {
			defer g.Done()
			assert.Equal(1, value.Val().(data).cnt)
		}()
	}

	for i := 0; i < 2; i++ {
		go func() {
			defer g.Done()
			value.GetAndSet(func(x interface{}) interface{} {
				assert.Equal(1, x.(data).cnt)
				d := x.(data)
				d.cnt = 1
				return d
			})
		}()
	}
	g.Wait()
}

func TestString_IsEmpty(t *testing.T) {
	type fields struct {
		val string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"empty", fields{
			val: "",
		}, true},
		{"not_empty", fields{
			val: "abc",
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &String{
				val: tt.fields.val,
			}
			assert2.Equalf(t, tt.want, a.IsEmpty(), "IsEmpty()")
		})
	}
}

func TestBool_IsFalse(t *testing.T) {
	a := int32(1)
	b := int32(0)
	type fields struct {
		val *int32
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"false", fields{val: &a}, false},
		{"true", fields{val: &b}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Bool{
				val: tt.fields.val,
			}
			assert2.Equalf(t, tt.want, i.IsFalse(), "IsFalse()")
		})
	}
}

func TestBool_IsTrue(t *testing.T) {
	a := int32(1)
	b := int32(0)
	type fields struct {
		val *int32
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"false", fields{val: &b}, false},
		{"true", fields{val: &a}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Bool{
				val: tt.fields.val,
			}
			assert2.Equalf(t, tt.want, i.IsTrue(), "IsTrue()")
		})
	}
}

func TestBool_SetFalse(t *testing.T) {
	a := NewBool()
	a.SetFalse()
	assert2.True(t, a.IsFalse())
}

func TestBool_SetTrue(t *testing.T) {
	a := NewBool()
	a.SetTrue()
	assert2.True(t, a.IsTrue())
}

func TestBytes_Append(t *testing.T) {
	a := NewBytes(nil)
	a.Append([]byte("abc"))
	assert2.Equal(t, []byte("abc"), a.Bytes())
}

func TestBytes_Len(t *testing.T) {
	a := NewBytes(nil)
	a.Append([]byte("abc"))
	assert2.Equal(t, 3, a.Len())
}

func TestBytes_SetBytes(t *testing.T) {
	a := NewBytes([]byte("abc"))
	a.Append([]byte("abc"))
	assert2.Equal(t, []byte("abcabc"), a.Bytes())
}

func TestString_SetVal(t *testing.T) {
	a := NewString()
	a.SetVal("abc")
	assert2.Equal(t, "abc", a.Val())
}

func TestString_String(t *testing.T) {
	a := NewString("abc")
	assert2.Equal(t, "abc", a.String())
}

func TestBytes_SetBytes1(t *testing.T) {
	a := NewBytes()
	a.SetBytes([]byte("abc"))
	assert2.Equal(t, []byte("abc"), a.Bytes())
}

func TestAny_SetVal(t *testing.T) {
	a := NewAny("string")
	a.SetVal(io.Reader(nil))
	a.SetVal(io.ReadCloser(nil))
	assert2.IsType(t, io.Reader(nil), a.Val())
	assert2.IsType(t, io.ReadCloser(nil), a.Val())
	a.GetAndSet(func(oldVal interface{}) (newVal interface{}) {
		assert2.IsType(t, io.ReadCloser(nil), oldVal)
		return oldVal
	})
}
