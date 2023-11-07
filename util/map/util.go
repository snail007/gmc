package gmap

import (
	"fmt"
	gcast "github.com/snail007/gmc/util/cast"
	"sort"
)

// ToAny converts map[string]string to map[string]interface{}
func ToAny(a map[string]string) (b map[string]interface{}) {
	b = map[string]interface{}{}
	for k, v := range a {
		b[k] = v
	}
	return
}

// ToString converts map[string]interface{} to map[string]string
func ToString(a map[string]interface{}) (b map[string]string) {
	b = map[string]string{}
	for k, v := range a {
		switch vv := v.(type) {
		case string:
			b[k] = vv
		default:
			b[k] = fmt.Sprintf("%v", vv)
		}
	}
	return
}

func SortStr(s []map[string]string, key string, aes bool) {
	sort.Slice(s, func(i, j int) bool {
		if aes {
			return Less(s[i][key], s[j][key])
		}
		return !Less(s[i][key], s[j][key])
	})
	return
}

func Sort(s []map[string]interface{}, key string, aes bool) {
	sort.Slice(s, func(i, j int) bool {
		if aes {
			return Less(s[i][key], s[j][key])
		}
		return !Less(s[i][key], s[j][key])
	})
	return
}

func Less(v1, v2 interface{}) bool {
	s1 := gcast.ToString(v1)
	a := []string{s1, gcast.ToString(v2)}
	sort.Strings(a)
	return a[0] == s1
}

func SortMap(m map[string]interface{}, aes bool) *Map {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if aes {
			return Less(keys[i], keys[j])
		}
		return !Less(keys[i], keys[j])
	})
	newMap := New()
	for _, k := range keys {
		newMap.Store(k, m[k])
	}
	return newMap
}

func SortMapStr(m map[string]string, aes bool) *Map {
	return SortMap(ToAny(m), aes)
}
