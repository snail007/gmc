// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gmap

import (
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
			if gotB := InterfaceToStr(tt.args.a); !reflect.DeepEqual(gotB, tt.wantB) {
				t.Errorf("InterfaceToStr() = %v, want %v", gotB, tt.wantB)
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
			if gotB := strToInterface(tt.args.a); !reflect.DeepEqual(gotB, tt.wantB) {
				t.Errorf("strToInterface() = %v, want %v", gotB, tt.wantB)
			}
		})
	}
}
