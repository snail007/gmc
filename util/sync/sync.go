package gsync

import "sync"

func WaitGroupToChan(g *sync.WaitGroup) <-chan bool {
	ch := make(chan bool)
	go func() {
		g.Wait()
		ch <- true
	}()
	return ch
}

type Mutex struct {
	*sync.Mutex
	key string
}

var (
	mutexGetLock sync.Mutex
	lockMap      = map[string]*Mutex{}
)

func GetLock(key string) *Mutex {
	mutexGetLock.Lock()
	defer mutexGetLock.Unlock()
	if v, ok := lockMap[key]; ok {
		return v
	}
	v := newMutex(key)
	lockMap[key] = v
	return v
}

func newMutex(Key string) *Mutex {
	return &Mutex{
		key:   Key,
		Mutex: &sync.Mutex{},
	}
}

func (r *Mutex) Unlock() {
	r.Mutex.Unlock()
	r.removeFromMap()
}

func (r *Mutex) removeFromMap() {
	mutexGetLock.Lock()
	defer mutexGetLock.Unlock()
	if lockMap[r.key] == nil {
		return
	}
	delete(lockMap, r.key)
}
