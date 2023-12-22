package gsync

import (
	"sync"
	"testing"
	"time"
)

func TestMutex(t *testing.T) {
	t.Run("Single Goroutine", func(t *testing.T) {
		key := "testKey"
		lock := GetLock(key)

		lock.Lock()
		defer lock.Unlock()

		// 在锁内执行一些操作
		counter := incrementCounter(key)
		assertCounter(t, key, counter, 1)
	})

	t.Run("Concurrent Goroutines", func(t *testing.T) {
		key := "concurrentKey"

		var wg sync.WaitGroup
		numGoroutines := 5

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				lock := GetLock(key)
				lock.Lock()
				defer lock.Unlock()

				incrementCounter(key)
			}()
		}

		wg.Wait()
		assertCounter(t, key, counterMap[key], 5)
	})

	t.Run("Lock Removed from Map", func(t *testing.T) {
		key := "removeKey"

		lock := GetLock(key)
		lock.Lock()
		counter := incrementCounter(key)
		assertCounter(t, key, counter, 1)
		lock.Unlock()
		time.Sleep(time.Millisecond * 50)

		lock = GetLock(key)
		lock.Lock()
		counter = incrementCounter(key)
		lock.Unlock()
		assertCounter(t, key, counter, 2)
	})

	t.Run("Concurrent GetLock", func(t *testing.T) {
		key := "concurrentGetLockKey"
		var wg sync.WaitGroup
		numGetLocks := 5
		for i := 0; i < numGetLocks; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				lock := GetLock(key)
				lock.Lock()
				defer lock.Unlock()
				time.Sleep(time.Millisecond * 100)
				incrementCounter(key)
			}()
		}

		wg.Wait()
		assertCounter(t, key, counterMap[key], 5)
	})
}

var (
	counterMap = map[string]int{}
)

func incrementCounter(key string) int {
	val, loaded := counterMap[key]
	if !loaded {
		val = 1
	} else {
		val++
	}
	counterMap[key] = val
	return val
}

func assertCounter(t *testing.T, key string, actual, expected int) {
	t.Helper()
	if actual != expected {
		t.Errorf("Counter for key %s: Expected %d, but got %d", key, expected, actual)
	}
}

func TestWait(t *testing.T) {
	t.Run("Wait for WaitGroup", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(1)

		ch := WaitGroupToChan(&wg)

		go func() {
			time.Sleep(time.Millisecond * 50)
			wg.Done()
		}()

		select {
		case <-ch:
		case <-time.After(time.Second):
			t.Error("Timeout waiting for WaitGroup")
		}
	})

	t.Run("Wait for WaitGroup with multiple goroutines", func(t *testing.T) {
		var wg sync.WaitGroup
		numGoroutines := 3
		wg.Add(numGoroutines)

		ch := WaitGroupToChan(&wg)

		for i := 0; i < numGoroutines; i++ {
			go func() {
				time.Sleep(time.Millisecond * 50)
				wg.Done()
			}()
		}

		select {
		case <-ch:
		case <-time.After(time.Second):
			t.Error("Timeout waiting for WaitGroup")
		}
	})
}
