// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More infomation at https://github.com/snail007/gmc

package grouter

import "reflect"

func methods(obj interface{}) (m []string) {
	t := reflect.TypeOf(obj)
	for i := 0; i < t.NumMethod(); i++ {
		m = append(m, t.Method(i).Name)
	}
	return
}
func invoke(obj interface{}, name string, args ...interface{}) []reflect.Value {
	inputs := make([]reflect.Value, len(args))
	for i := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	obj0, ok := obj.(reflect.Value)
	if ok {
		return obj0.MethodByName(name).Call(inputs)
	}
	return reflect.ValueOf(obj).MethodByName(name).Call(inputs)
}
