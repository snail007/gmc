// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More infomation at https://github.com/snail007/gmc

package gpool

import (
	"fmt"
	gcore "github.com/snail007/gmc/core"
	gerror "github.com/snail007/gmc/module/error"
	glog "github.com/snail007/gmc/module/log"
	assert2 "github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

//testing here

func TestNewGPool(t *testing.T) {
	p := NewGPool(3)
	if p != nil {
		t.Log("NewGPool is okay")
	} else {
		t.Fatalf("NewGPool is failed")
	}
	p.Stop()
}

func TestSubmit(t *testing.T) {
	p := NewGPool(3)
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
	p.Stop()
}
func TestStop(t *testing.T) {
	p := NewGPool(3)
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
	p := NewGPool(3)
	p.SetLogger(nil)
	p.Stop()
}
func TestRunning(t *testing.T) {
	p := NewGPool(3)
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
	if p.Running() == 3 {
		t.Log("Running is okay")
	} else {
		t.Fatalf("Running is failed")
	}
	p.Stop()
}

func TestIncrease(t *testing.T) {
	assert := assert2.New(t)
	p := NewGPool(3)
	for i := 0; i < 3; i++ {
		p.Submit(func() {
			fmt.Println(">>>>")
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
	t.Log(p.WorkerCount())
	assert.Equal(6, p.Running())
	p.Stop()
}

func TestAwaiting(t *testing.T) {
	p := NewGPool(3)
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
	if p.Awaiting() == 1 {
		t.Log("Awaiting is okay")
	} else {
		t.Fatalf("Awaiting is failed")
	}
	p.Stop()
}

func TestMain(m *testing.M) {
	providers := gcore.Providers
	providers.RegisterLogger("", func(ctx gcore.Ctx, prefix string) gcore.Logger {
		if ctx == nil {
			return glog.NewGMCLog(prefix)
		}
		return glog.NewFromConfig(ctx.Config(), prefix)
	})
	providers.RegisterError("", func() gcore.Error {
		return gerror.New()
	})
	os.Exit(m.Run())
}
