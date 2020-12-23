package gmap

import "fmt"

// StrStrToStrI converts map[string]string to map[string]interface{}
func StrStrToStrI(a map[string]string) (b map[string]interface{}) {
	b = map[string]interface{}{}
	for k, v := range a {
		b[k] = v
	}
	return
}

// StrIToStrStr converts map[string]interface{} to map[string]string
func StrIToStrStr(a map[string]interface{}) (b map[string]string) {
	b = map[string]string{}
	for k, v := range a {
		switch vv:=v.(type) {
		case string:
			b[k] = vv
		default:
			b[k] = fmt.Sprintf("%v",vv)
		}
	}
	return
}