package router

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
