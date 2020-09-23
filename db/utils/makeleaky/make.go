package makeleaky

import (
	"reflect"
	"sync"
)

var (
	cacheX    = map[string]map[uint64](chan interface{}){}
	lock      = sync.Mutex{}
	cacheSize = 1024
)

func GetX(x interface{}, length uint64, factory func() interface{}) (item interface{}) {
	t := reflect.TypeOf(x)
	typ := t.String()
	lock.Lock()
	defer lock.Unlock()
	if _, ok := cacheX[typ]; !ok {
		cacheX[typ] = map[uint64]chan interface{}{}
	}
	if _, ok := cacheX[typ][length]; !ok {
		cacheX[typ][length] = make(chan interface{}, length)
	}
	select {
	case item = <-cacheX[typ][length]:
		// fmt.Printf(">>> form cache , typ  %s, len %d\n", typ, length)
	default:
		item = factory()
		// fmt.Printf(">>> new cache , typ  %s, len %d\n", typ, length)
	}
	return
}
func PutX(x interface{}, length uint64) {
	t := reflect.TypeOf(x)
	typ := t.String()
	lock.Lock()
	defer lock.Unlock()
	if _, ok := cacheX[typ]; !ok {
		return
	}
	if _, ok := cacheX[typ][length]; !ok {
		return
	}
	select {
	case cacheX[typ][length] <- x:
		// fmt.Printf(">>> put cache ok, typ  %s, len %d\n", typ, length)
	default:
		// fmt.Printf(">>> put drop cache , typ  %s, len %d\n", typ, length)
	}
	return
}
