// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gpool_test

import (
	_ "github.com/snail007/gmc"
	gcore "github.com/snail007/gmc/core"
	gerror "github.com/snail007/gmc/module/error"
	glog "github.com/snail007/gmc/module/log"
	"github.com/snail007/gmc/util/gpool"
	gloop "github.com/snail007/gmc/util/loop"
	gatomic "github.com/snail007/gmc/util/sync/atomic"
	assert2 "github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestPreAlloc(t *testing.T) {
	p := gpool.NewWithPreAlloc(3)
	time.Sleep(time.Millisecond * 50)
	assert2.Equal(t, 3, p.IdleWorkerCount())
}

func TestBlocking(t *testing.T) {
	p := gpool.New(1)
	p.SetMaxJobCount(1)
	p.SetBlocking(true)
	assert2.True(t, p.Blocking())
	p.SetDebug(true)
	// this task will be run
	p.Submit(func() {
		time.Sleep(time.Second)
	})
	// this task will be queued
	p.Submit(func() {
		time.Sleep(time.Second)
	})
	// this submit will be blocking because of queue size is 1
	start := time.Now()
	p.Submit(func() {
		time.Sleep(time.Second)
	})
	dur := time.Since(start)
	t.Log(dur)
	assert2.Greater(t, dur, time.Second)
	p.Stop()
}

func TestIdle(t *testing.T) {
	p := gpool.New(3)
	p.SetDebug(true)
	p.SetIdleDuration(time.Second)
	assert2.Equal(t, time.Second, p.IdleDuration())
	cnt := gatomic.NewInt(0)
	gloop.For(3, func(loopIndex int) {
		p.Submit(func() {
			cnt.Increase(loopIndex)
		})
	})
	time.Sleep(time.Millisecond * 500)
	assert2.Equal(t, 3, cnt.Val())
	time.Sleep(time.Second * 2)
	assert2.Equal(t, 0, p.WorkerCount())
	p.Stop()
}

func TestNewGPool(t *testing.T) {
	p := gpool.New(3)
	if p != nil {
		t.Log("New is okay")
	} else {
		t.Fatalf("New is failed")
	}
	p.Stop()
}

func TestSubmit(t *testing.T) {
	p := gpool.New(3)
	a := make(chan bool)
	p.Submit(func() {
		a <- true
	})
	p.Submit(func() {})
	p.Submit(func() {})
	p.Submit(func() {
		x := 0
		_ = 10 / x
	})
	if <-a {
		t.Log("Submit is okay")
	} else {
		t.Fatalf("Submit is failed")
	}
	time.Sleep(time.Second)
}
func TestStop(t *testing.T) {
	p := gpool.New(3)
	a := make(chan bool)
	p.Submit(func() {
		time.Sleep(time.Second)
		a <- true
	})
	p.Stop()
	select {
	case <-a:
		t.Fatalf("Stop is failed")
	case <-time.After(time.Millisecond):
		t.Log("Stop is okay")
	}
	p.Stop()
}
func TestSetLogger(t *testing.T) {
	p := gpool.New(3)
	p.SetLogger(nil)
	p.Stop()
}
func TestWaitDone(t *testing.T) {
	start := time.Now()
	p := gpool.NewWithPreAlloc(10)
	gloop.For(10, func(loopIndex int) {
		p.Submit(func() {
			time.Sleep(time.Millisecond * 100)
		})
	})
	p.WaitDone()
	assert2.Greater(t, time.Since(start), time.Millisecond*100)
	assert2.Less(t, time.Since(start), time.Millisecond*150)
}

func TestRunning(t *testing.T) {
	p := gpool.NewWithPreAlloc(3)
	p.Submit(func() {
		time.Sleep(time.Second)
	})
	p.Submit(func() {
		time.Sleep(time.Second)
	})
	p.Submit(func() {
		time.Sleep(time.Second)
	})
	p.Submit(func() {
		time.Sleep(time.Second)
	})
	time.Sleep(time.Millisecond * 40)
	if p.RunningWorkCount() == 3 {
		t.Log("Running is okay")
	} else {
		t.Fatalf("Running is failed")
	}
	p.Stop()
}

func TestIncrease(t *testing.T) {
	assert := assert2.New(t)
	p := gpool.NewWithPreAlloc(3)
	for i := 0; i < 3; i++ {
		p.Submit(func() {
			time.Sleep(time.Second * 5)
		})
	}
	p.Increase(3)
	for i := 0; i < 3; i++ {
		p.Submit(func() {
			time.Sleep(time.Second * 5)
		})
	}
	time.Sleep(time.Second)
	assert.Equal(6, p.RunningWorkCount())
	p.Stop()
}

func TestDecrease(t *testing.T) {
	assert := assert2.New(t)
	p := gpool.NewWithPreAlloc(2)
	for i := 0; i < 6; i++ {
		p.Submit(func() {
			time.Sleep(time.Second)
		})
	}
	time.Sleep(time.Millisecond * 30)
	assert.Equal(2, p.RunningWorkCount())
	p.Decrease(1)
	time.Sleep(time.Second)
	assert.Equal(1, p.RunningWorkCount())
	p.Stop()
}

func TestQueuedJobCount(t *testing.T) {
	p := gpool.NewWithPreAlloc(3)
	p.Submit(func() {
		time.Sleep(time.Second)
	})
	p.Submit(func() {
		time.Sleep(time.Second)
	})
	p.Submit(func() {
		time.Sleep(time.Second)
	})
	p.Submit(func() {
		time.Sleep(time.Second)
	})
	time.Sleep(time.Millisecond * 40)
	if p.QueuedJobCount() == 1 {
		t.Log("QueuedJobCount is okay")
	} else {
		t.Fatalf("QueuedJobCount is failed")
	}
	p.Stop()
}

func TestGPool_MaxWaitCount(t *testing.T) {
	assert := assert2.New(t)
	p := gpool.New(1)
	//p.SetLogger(glog.New())
	p.SetDebug(true)
	assert.True(p.IsDebug())

	p.SetMaxJobCount(1)
	assert.Equal(1, p.MaxJobCount())

	//check reset
	p.ResetTo(2)
	assert.Equal(2, p.WorkerCount())
	//wait worker
	time.Sleep(time.Millisecond * 500)
	assert2.Equal(t, 2, p.IdleWorkerCount())

	p.ResetTo(1)
	assert.Equal(1, p.WorkerCount())
	//wait worker
	time.Sleep(time.Millisecond * 500)
	assert2.Equal(t, 1, p.IdleWorkerCount())

	assert.Nil(p.Submit(func() {
		time.Sleep(time.Second)
	}))

	//wait worker to fetch task
	time.Sleep(time.Millisecond * 100)

	assert.Nil(p.Submit(func() {
		time.Sleep(time.Second)
	}))

	assert.NotNil(p.Submit(func() {
		time.Sleep(time.Second)
	}))

	time.Sleep(time.Millisecond * 40)
	assert.Equal(0, p.IdleWorkerCount())
	assert.Equal(1, p.RunningWorkCount())
	p.Stop()
}

func TestMain(m *testing.M) {
	gcore.RegisterLogger(gcore.DefaultProviderKey, func(ctx gcore.Ctx, prefix string) gcore.Logger {
		if ctx == nil {
			return glog.New(prefix)
		}
		return glog.NewFromConfig(ctx.Config(), prefix)
	})
	gcore.RegisterError(gcore.DefaultProviderKey, func() gcore.Error {
		return gerror.New()
	})
	os.Exit(m.Run())
}
