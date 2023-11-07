// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gmap

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestInterfaceToStr(t *testing.T) {
	type args struct {
		a map[string]interface{}
	}
	tests := []struct {
		name  string
		args  args
		wantB map[string]string
	}{
		{"map", args{a: map[string]interface{}{
			"a": 1, "b": 2,
		}}, map[string]string{
			"a": "1", "b": "2",
		}},
		{"map", args{a: map[string]interface{}{
			"a": "1", "b": "2",
		}}, map[string]string{
			"a": "1", "b": "2",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotB := ToString(tt.args.a); !reflect.DeepEqual(gotB, tt.wantB) {
				t.Errorf("ToString() = %v, want %v", gotB, tt.wantB)
			}
		})
	}
}

func Test_strToInterface(t *testing.T) {
	type args struct {
		a map[string]string
	}
	tests := []struct {
		name  string
		args  args
		wantB map[string]interface{}
	}{
		{"map", args{a: map[string]string{
			"a": "1", "b": "2",
		}}, map[string]interface{}{
			"a": "1", "b": "2",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotB := ToAny(tt.args.a); !reflect.DeepEqual(gotB, tt.wantB) {
				t.Errorf("ToAny() = %v, want %v", gotB, tt.wantB)
			}
		})
	}
}

func TestCompare(t *testing.T) {
	assert.True(t, Less("a", "aab"))
	assert.False(t, Less("c", "aab"))
	assert.True(t, Less("123", "234"))
}

func TestSortAnyKey(t *testing.T) {
	s := []map[string]interface{}{
		{"name": "c"}, {"name": "a"}, {"name": "b"},
	}
	Sort(s, "name", true)
	assert.Equal(t, "a", s[0]["name"])
	assert.Equal(t, "b", s[1]["name"])
	assert.Equal(t, "c", s[2]["name"])
}

func TestSortAnyKey_2(t *testing.T) {
	s := []map[string]interface{}{
		{"name": "c"}, {"name": "a"}, {"name": "b"},
	}
	Sort(s, "name", false)
	assert.Equal(t, "c", s[0]["name"])
	assert.Equal(t, "b", s[1]["name"])
	assert.Equal(t, "a", s[2]["name"])
}

func TestSortStr_1(t *testing.T) {
	s := []map[string]string{
		{"name": "c"}, {"name": "a"}, {"name": "b"},
	}
	SortStr(s, "name", true)
	assert.Equal(t, "a", s[0]["name"])
	assert.Equal(t, "b", s[1]["name"])
	assert.Equal(t, "c", s[2]["name"])
}

func TestSortStr(t *testing.T) {
	s := []map[string]string{
		{"name": "c"}, {"name": "a"}, {"name": "b"},
	}
	SortStr(s, "name", false)
	assert.Equal(t, "c", s[0]["name"])
	assert.Equal(t, "b", s[1]["name"])
	assert.Equal(t, "a", s[2]["name"])
}

func TestSortMap(t *testing.T) {
	m := map[string]string{
		"2": "6",
		"3": "7",
		"1": "5",
	}
	v := []string{}
	SortMapStr(m, true).RangeFast(func(key, value interface{}) bool {
		v = append(v, value.(string))
		return true
	})
	assert.Equal(t, []string{"5", "6", "7"}, v)
	v = []string{}
	SortMapStr(m, false).RangeFast(func(key, value interface{}) bool {
		v = append(v, value.(string))
		return true
	})
	assert.Equal(t, []string{"7", "6", "5"}, v)
}
